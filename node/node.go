
package node

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"blockchain-node/config"
	"blockchain-node/consensus"
	"blockchain-node/core"
	"blockchain-node/logger"
	"blockchain-node/mempool"
	"blockchain-node/metrics"
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
	metrics    *metrics.Metrics
	logger     *logger.Logger
	
	// Graceful shutdown
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	shutdownCh chan struct{}
}

// NewNode creates a new blockchain node
func NewNode(cfg *config.Config) (*Node, error) {
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	// Initialize logger
	if err := logger.Init(logger.Config{
		Level:     cfg.Logging.Level,
		Output:    cfg.Logging.Output,
		FilePath:  cfg.Logging.FilePath,
		MaxSize:   cfg.Logging.MaxSize,
		Component: cfg.Logging.Component,
	}); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	nodeLogger := logger.NewLogger("node")
	nodeLogger.Info("Initializing blockchain node...")

	// Initialize metrics
	metricsInstance := metrics.Init(&cfg.Metrics)

	// Initialize database with optimized settings
	db, err := storage.NewLevelDB(cfg.DB.Path, &storage.LevelDBOptions{
		CacheSize:    cfg.DB.CacheSize,
		MaxOpenFiles: cfg.DB.MaxOpenFiles,
		WriteBuffer:  cfg.DB.WriteBuffer,
	})
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

	// Initialize mempool with configuration
	mempool := mempool.NewMempool(&mempool.Config{
		MaxSize:     1000,
		MinGasPrice: cfg.EVM.MinGasPrice,
	})

	// Initialize consensus
	consensus := consensus.NewProofOfWork(big.NewInt(int64(cfg.Mining.Difficulty)))

	// Initialize P2P server
	p2pServer := p2p.NewServer(&cfg.Network)

	// Initialize RPC server
	var rpcServer *rpc.Server
	if cfg.RPC.Enabled {
		rpcServer = rpc.NewServer(&cfg.RPC, blockchain, mempool)
	}

	ctx, cancel := context.WithCancel(context.Background())

	node := &Node{
		config:     cfg,
		blockchain: blockchain,
		mempool:    mempool,
		consensus:  consensus,
		p2pServer:  p2pServer,
		rpcServer:  rpcServer,
		db:         db,
		metrics:    metricsInstance,
		logger:     nodeLogger,
		ctx:        ctx,
		cancel:     cancel,
		shutdownCh: make(chan struct{}),
	}

	nodeLogger.Info("Blockchain node initialized successfully")
	return node, nil
}

// Start starts the blockchain node
func (n *Node) Start() error {
	n.logger.Info("Starting blockchain node...")

	// Start P2P server
	if err := n.p2pServer.Start(); err != nil {
		return fmt.Errorf("failed to start P2P server: %v", err)
	}
	n.logger.Info("P2P server started on port %d", n.config.Network.Port)

	// Start RPC server
	if n.rpcServer != nil {
		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			if err := n.rpcServer.Start(); err != nil {
				n.logger.Error("RPC server error: %v", err)
			}
		}()
		n.logger.Info("RPC server started on %s:%d", n.config.RPC.Host, n.config.RPC.Port)
	}

	// Start mining if enabled
	if n.config.Mining.Enabled {
		n.wg.Add(1)
		go func() {
			defer n.wg.Done()
			n.startMining()
		}()
		n.logger.Info("Mining started with %d threads", n.config.Mining.Threads)
	}

	// Start metrics updater
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		n.updateMetrics()
	}()

	n.logger.Info("Node started successfully!")
	n.logger.Info("- Chain ID: %d", n.config.EVM.ChainID)
	n.logger.Info("- P2P listening on port %d", n.config.Network.Port)
	if n.config.RPC.Enabled {
		n.logger.Info("- RPC server on %s:%d", n.config.RPC.Host, n.config.RPC.Port)
	}
	n.logger.Info("- Mining enabled: %t", n.config.Mining.Enabled)
	if n.config.Metrics.Enabled {
		n.logger.Info("- Metrics server on port %d", n.config.Metrics.Port)
	}

	// Wait for shutdown signal
	n.waitForShutdown()

	return nil
}

// Stop stops the blockchain node gracefully
func (n *Node) Stop() error {
	n.logger.Info("Stopping blockchain node...")

	// Signal shutdown
	close(n.shutdownCh)
	n.cancel()

	// Stop P2P server
	if err := n.p2pServer.Stop(); err != nil {
		n.logger.Error("Error stopping P2P server: %v", err)
	}

	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		n.wg.Wait()
		close(done)
	}()

	// Wait with timeout
	select {
	case <-done:
		n.logger.Info("All services stopped")
	case <-time.After(30 * time.Second):
		n.logger.Warning("Shutdown timeout reached, forcing exit")
	}

	// Close database
	if err := n.db.Close(); err != nil {
		n.logger.Error("Error closing database: %v", err)
	}

	// Close logger
	if err := logger.Close(); err != nil {
		n.logger.Error("Error closing logger: %v", err)
	}

	n.logger.Info("Node stopped successfully")
	return nil
}

// startMining starts the mining process with enhanced logging
func (n *Node) startMining() {
	n.logger.Info("Starting mining with %d threads, difficulty %s", 
		n.config.Mining.Threads, n.consensus.GetDifficulty().String())

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	hashCount := uint64(0)
	lastTime := time.Now()

	for {
		select {
		case <-n.ctx.Done():
			n.logger.Info("Mining stopped")
			return
		case <-ticker.C:
			// Calculate hash rate
			now := time.Now()
			elapsed := now.Sub(lastTime).Seconds()
			if elapsed > 0 {
				hashRate := float64(hashCount) / elapsed
				n.metrics.UpdateMiningHashRate(hashRate)
				n.logger.Debug("Current hash rate: %.2f H/s", hashRate)
				hashCount = 0
				lastTime = now
			}
		default:
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
			start := time.Now()
			if err := n.consensus.Mine(newBlock); err != nil {
				n.logger.Error("Mining error: %v", err)
				continue
			}
			miningTime := time.Since(start)

			// Add block to blockchain
			if err := n.blockchain.AddBlock(newBlock); err != nil {
				n.logger.Error("Failed to add block: %v", err)
				continue
			}

			// Remove mined transactions from mempool
			for _, tx := range pendingTxs {
				n.mempool.RemoveTransaction(tx.Hash)
				n.metrics.IncrementTransactions()
			}

			n.logger.Info("New block mined: #%s, Hash: %x, Transactions: %d, Time: %v",
				newBlock.Header.Number.String(), newBlock.Hash, len(newBlock.Transactions), miningTime)

			// Update metrics
			n.metrics.UpdateBlockHeight(newBlock.Header.Number.Uint64())

			// Broadcast block to peers
			n.p2pServer.BroadcastMessage([]byte(fmt.Sprintf("NEW_BLOCK:%x", newBlock.Hash)))

			hashCount += newBlock.Header.Nonce
		}
	}
}

// updateMetrics updates various metrics periodically
func (n *Node) updateMetrics() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-n.ctx.Done():
			return
		case <-ticker.C:
			// Update peer count
			peerCount := n.p2pServer.GetPeerCount()
			n.metrics.UpdatePeerCount(peerCount)

			// Update mempool size
			mempoolSize := n.mempool.Size()
			n.metrics.UpdateMempoolSize(mempoolSize)

			// Update block height
			blockHeight := n.blockchain.GetBlockNumber().Uint64()
			n.metrics.UpdateBlockHeight(blockHeight)

			n.logger.Debug("Metrics updated - Peers: %d, Mempool: %d, Block: %d", 
				peerCount, mempoolSize, blockHeight)
		}
	}
}

// waitForShutdown waits for shutdown signal
func (n *Node) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		n.logger.Info("Received signal: %v", sig)
	case <-n.shutdownCh:
		n.logger.Info("Shutdown requested")
	}

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

// GetMetrics returns the metrics instance
func (n *Node) GetMetrics() *metrics.Metrics {
	return n.metrics
}
