
package rpc

import (
	"fmt"
	"math/big"
	"net/http"

	"blockchain-node/config"
	"blockchain-node/core"
	"blockchain-node/mempool"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// Server represents the RPC server
type Server struct {
	config     *config.RPCConfig
	blockchain *core.Blockchain
	mempool    *mempool.Mempool
	router     *gin.Engine
}

// NewServer creates a new RPC server
func NewServer(cfg *config.RPCConfig, bc *core.Blockchain, mp *mempool.Mempool) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	server := &Server{
		config:     cfg,
		blockchain: bc,
		mempool:    mp,
		router:     router,
	}

	server.setupRoutes()
	return server
}

// Start starts the RPC server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	fmt.Printf("RPC server starting on %s\n", addr)
	return s.router.Run(addr)
}

// setupRoutes sets up the RPC routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/")
	
	// Ethereum-compatible JSON-RPC endpoints
	api.POST("/", s.handleJSONRPC)
	
	// RESTful endpoints for easier testing
	api.GET("/eth/blockNumber", s.getBlockNumber)
	api.GET("/eth/getBalance/:address", s.getBalance)
	api.GET("/eth/getBlockByNumber/:number", s.getBlockByNumber)
	api.GET("/eth/getBlockByHash/:hash", s.getBlockByHash)
	api.GET("/eth/getTransactionCount/:address", s.getTransactionCount)
	api.POST("/eth/sendRawTransaction", s.sendRawTransaction)
}

// JSONRPCRequest represents a JSON-RPC request
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents an RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// handleJSONRPC handles JSON-RPC requests
func (s *Server) handleJSONRPC(c *gin.Context) {
	var req JSONRPCRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32700, Message: "Parse error"},
			ID:      req.ID,
		})
		return
	}

	var result interface{}
	var err error

	switch req.Method {
	case "eth_blockNumber":
		result = fmt.Sprintf("0x%x", s.blockchain.GetBlockNumber().Uint64())
	case "eth_getBalance":
		if len(req.Params) < 1 {
			err = fmt.Errorf("missing address parameter")
			break
		}
		address := req.Params[0].(string)
		// TODO: Implement balance lookup
		result = "0x0"
		_ = address
	case "eth_getBlockByNumber":
		if len(req.Params) < 2 {
			err = fmt.Errorf("missing parameters")
			break
		}
		// TODO: Implement block lookup
		result = nil
	case "eth_chainId":
		result = "0x539" // 1337 in hex
	case "eth_gasPrice":
		result = "0x3b9aca00" // 1 Gwei in hex
	default:
		err = fmt.Errorf("method not found: %s", req.Method)
	}

	if err != nil {
		c.JSON(http.StatusOK, JSONRPCResponse{
			JSONRPC: "2.0",
			Error:   &RPCError{Code: -32601, Message: err.Error()},
			ID:      req.ID,
		})
		return
	}

	c.JSON(http.StatusOK, JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	})
}

// getBlockNumber returns the current block number
func (s *Server) getBlockNumber(c *gin.Context) {
	blockNumber := s.blockchain.GetBlockNumber()
	c.JSON(http.StatusOK, gin.H{
		"blockNumber": fmt.Sprintf("0x%x", blockNumber.Uint64()),
	})
}

// getBalance returns the balance of an address
func (s *Server) getBalance(c *gin.Context) {
	address := c.Param("address")
	// TODO: Implement actual balance lookup from state
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"balance": "0x0",
	})
}

// getBlockByNumber returns a block by number
func (s *Server) getBlockByNumber(c *gin.Context) {
	numberStr := c.Param("number")
	
	var number *big.Int
	if numberStr == "latest" {
		number = s.blockchain.GetBlockNumber()
	} else {
		number = new(big.Int)
		number.SetString(numberStr, 0)
	}

	block, err := s.blockchain.GetBlockByNumber(number)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number":       fmt.Sprintf("0x%x", block.Header.Number.Uint64()),
		"hash":         block.Hash.Hex(),
		"parentHash":   block.Header.PreviousHash.Hex(),
		"gasLimit":     fmt.Sprintf("0x%x", block.Header.GasLimit),
		"gasUsed":      fmt.Sprintf("0x%x", block.Header.GasUsed),
		"timestamp":    fmt.Sprintf("0x%x", block.Header.Timestamp),
		"transactions": len(block.Transactions),
	})
}

// getBlockByHash returns a block by hash
func (s *Server) getBlockByHash(c *gin.Context) {
	hashStr := c.Param("hash")
	hash := common.HexToHash(hashStr)

	block, err := s.blockchain.GetBlockByHash(hash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"number":       fmt.Sprintf("0x%x", block.Header.Number.Uint64()),
		"hash":         block.Hash.Hex(),
		"parentHash":   block.Header.PreviousHash.Hex(),
		"gasLimit":     fmt.Sprintf("0x%x", block.Header.GasLimit),
		"gasUsed":      fmt.Sprintf("0x%x", block.Header.GasUsed),
		"timestamp":    fmt.Sprintf("0x%x", block.Header.Timestamp),
		"transactions": len(block.Transactions),
	})
}

// getTransactionCount returns the transaction count for an address
func (s *Server) getTransactionCount(c *gin.Context) {
	address := c.Param("address")
	// TODO: Implement actual nonce lookup from state
	c.JSON(http.StatusOK, gin.H{
		"address": address,
		"nonce":   "0x0",
	})
}

// sendRawTransaction sends a raw transaction
func (s *Server) sendRawTransaction(c *gin.Context) {
	var req struct {
		Data string `json:"data"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: Decode and process raw transaction
	c.JSON(http.StatusOK, gin.H{
		"txHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
	})
}
