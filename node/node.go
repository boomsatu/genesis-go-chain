
package node

import (
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blockchain-node/config"
	"blockchain-node/consensus"
	"blockchain-node/core"
	"blockchain-node/mempool"
	"blockchain-node/p2p"
	"blockchain-node/rpc"
	"blockchain-node/storage"
)

// Node represents the blockchain node
type Node struct {
	config     *config.Config
	blockchain *core.Blockchain
	mempool    *mempool.Mempool
	consensus  *consensus.ProofOfWork
	p2pServer  *p2p.Server
	rpcServer  *rpc.Server
	db         storage.Database
}

// NewNode creates a new blockchain node
func NewNode(cfg *config.Config) (*Node, error) {
	// Initialize database
	db, err := storage.NewLevelDB(cfg.DB.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	// Initialize blockchain
	genesis := core.DefaultGenesis()
	genesis.Config.ChainID = big.NewInt(int64(cfg.EVM.ChainID))
	genesis.GasLimit = cfg.EVM.BlockGasLimit

	blockchain, err := core.NewBlockchain(db, genesis)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize blockchain: %v", err)
	}

	// Initialize mempool
	mempool := mempool.NewMempool()

	// Initialize consensus
	consensus := consensus.NewProofOfWork(big.NewInt(int64(cfg.Mining.Difficulty)))

	// Initialize P2P server
	p2pServer := p2p.NewServer(&cfg.Network)

	// Initialize RPC server
	var rpcServer *rpc.Server
	if cfg.RPC.Enabled {
		rpcServer = rpc.NewServer(&cfg.RPC, blockchain, mempool)
	}

	return &Node{
		config:     cfg,
		blockchain: blockchain,
		mempool:    mempool,
		consensus:  consensus,
		p2pServer:  p2pServer,
		rpcServer:  rpcServer,
		db:         db,
	}, nil
}

// Start starts the blockchain node
func (n *Node) Start() error {
	fmt.Println("Starting blockchain node...")

	// Start P2P server
	if err := n.p2pServer.Start(); err != nil {
		return fmt.Errorf("failed to start P2P server: %v", err)
	}

	// Start RPC server
	if n.rpcServer != nil {
		go func() {
			if err := n.rpcServer.Start(); err != nil {
				fmt.Printf("RPC server error: %v\n", err)
			}
		}()
	}

	// Start mining if enabled
	if n.config.Mining.Enabled {
		go n.startMining()
	}

	fmt.Printf("Node started successfully!\n")
	fmt.Printf("- P2P listening on port %d\n", n.config.Network.Port)
	if n.config.RPC.Enabled {
		fmt.Printf("- RPC server on %s:%d\n", n.config.RPC.Host, n.config.RPC.Port)
	}
	fmt.Printf("- Mining enabled: %t\n", n.config.Mining.Enabled)
	fmt.Printf("- Chain ID: %d\n", n.config.EVM.ChainID)

	// Wait for shutdown signal
	n.waitForShutdown()

	return nil
}

// Stop stops the blockchain node
func (n *Node) Stop() error {
	fmt.Println("Stopping blockchain node...")

	// Stop P2P server
	if err := n.p2pServer.Stop(); err != nil {
		fmt.Printf("Error stopping P2P server: %v\n", err)
	}

	// Close database
	if err := n.db.Close(); err != nil {
		fmt.Printf("Error closing database: %v\n", err)
	}

	fmt.Println("Node stopped.")
	return nil
}

// startMining starts the mining process
func (n *Node) startMining() {
	fmt.Printf("Starting mining with %d threads...\n", n.config.Mining.Threads)

	for {
		// Get pending transactions
		pendingTxs := n.mempool.GetPendingTransactionsForMining(1000)

		// Create new block
		currentBlock := n.blockchain.GetCurrentBlock()
		newBlockNumber := new(big.Int).Add(currentBlock.Header.Number, big.NewInt(1))

		header := &core.BlockHeader{
			PreviousHash: currentBlock.Hash,
			Number:       newBlockNumber,
			GasLimit:     n.config.EVM.BlockGasLimit,
			GasUsed:      0,
			Timestamp:    uint64(time.Now().Unix()),
			Difficulty:   n.consensus.GetDifficulty(),
		}

		newBlock := core.NewBlock(header, pendingTxs)

		// Mine the block
		if err := n.consensus.Mine(newBlock); err != nil {
			fmt.Printf("Mining error: %v\n", err)
			continue
		}

		// Add block to blockchain
		if err := n.blockchain.AddBlock(newBlock); err != nil {
			fmt.Printf("Failed to add block: %v\n", err)
			continue
		}

		// Remove mined transactions from mempool
		for _, tx := range pendingTxs {
			n.mempool.RemoveTransaction(tx.Hash)
		}

		fmt.Printf("New block mined: #%s, Hash: %x, Transactions: %d\n",
			newBlock.Header.Number.String(), newBlock.Hash, len(newBlock.Transactions))

		// Broadcast block to peers
		n.p2pServer.BroadcastMessage([]byte(fmt.Sprintf("NEW_BLOCK:%x", newBlock.Hash)))
	}
}

// waitForShutdown waits for shutdown signal
func (n *Node) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutdown signal received...")
	n.Stop()
}

// GetBlockchain returns the blockchain instance
func (n *Node) GetBlockchain() *core.Blockchain {
	return n.blockchain
}

// GetMempool returns the mempool instance
func (n *Node) GetMempool() *mempool.Mempool {
	return n.mempool
}

// GetP2PServer returns the P2P server instance
func (n *Node) GetP2PServer() *p2p.Server {
	return n.p2pServer
}
