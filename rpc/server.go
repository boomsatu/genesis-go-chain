
package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"blockchain-node/config"
	"blockchain-node/core"
	"blockchain-node/crypto"
	"blockchain-node/logger"
	"blockchain-node/mempool"

	"github.com/gorilla/mux"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      interface{} `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// RPC error codes
const (
	RPCErrorCodeParseError     = -32700
	RPCErrorCodeInvalidRequest = -32600
	RPCErrorCodeMethodNotFound = -32601
	RPCErrorCodeInvalidParams  = -32602
	RPCErrorCodeInternalError  = -32603
)

// Server represents the RPC server
type Server struct {
	config     *config.RPCConfig
	blockchain *core.Blockchain
	mempool    *mempool.Mempool
	server     *http.Server
	logger     *logger.Logger
	
	// Method handlers
	methods map[string]func(params interface{}) (interface{}, error)
}

// NewServer creates a new RPC server
func NewServer(config *config.RPCConfig, blockchain *core.Blockchain, mempool *mempool.Mempool) *Server {
	server := &Server{
		config:     config,
		blockchain: blockchain,
		mempool:    mempool,
		logger:     logger.NewLogger("rpc"),
		methods:    make(map[string]func(params interface{}) (interface{}, error)),
	}

	// Register RPC methods
	server.registerMethods()

	return server
}

// Start starts the RPC server
func (s *Server) Start() error {
	s.logger.Info("Starting RPC server", "host", s.config.Host, "port", s.config.Port)

	router := mux.NewRouter()
	
	// Add CORS middleware
	router.Use(s.corsMiddleware)
	
	// JSON-RPC endpoint
	router.HandleFunc("/", s.handleJSONRPC).Methods("POST", "OPTIONS")
	
	// Health check endpoint
	router.HandleFunc("/health", s.handleHealth).Methods("GET")
	
	// Stats endpoint
	router.HandleFunc("/stats", s.handleStats).Methods("GET")

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(s.config.Timeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Timeout) * time.Second,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("RPC server error", "error", err)
		}
	}()

	s.logger.Info("RPC server started successfully")
	return nil
}

// Stop stops the RPC server
func (s *Server) Stop() error {
	s.logger.Info("Stopping RPC server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to gracefully shutdown RPC server", "error", err)
		return err
	}

	s.logger.Info("RPC server stopped")
	return nil
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		for _, origin := range s.config.CORSOrigins {
			if origin == "*" || strings.Contains(r.Header.Get("Origin"), origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleJSONRPC handles JSON-RPC requests
func (s *Server) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, nil, RPCErrorCodeParseError, "Parse error", err.Error())
		return
	}

	// Validate JSON-RPC version
	if req.JSONRPC != "2.0" {
		s.sendError(w, req.ID, RPCErrorCodeInvalidRequest, "Invalid request", "JSON-RPC version must be 2.0")
		return
	}

	// Find method handler
	handler, exists := s.methods[req.Method]
	if !exists {
		s.sendError(w, req.ID, RPCErrorCodeMethodNotFound, "Method not found", req.Method)
		return
	}

	// Execute method
	result, err := handler(req.Params)
	if err != nil {
		s.sendError(w, req.ID, RPCErrorCodeInternalError, "Internal error", err.Error())
		return
	}

	// Send successful response
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      req.ID,
	}

	json.NewEncoder(w).Encode(response)
	s.logger.Debug("RPC method executed", "method", req.Method, "id", req.ID)
}

// sendError sends an error response
func (s *Server) sendError(w http.ResponseWriter, id interface{}, code int, message, data string) {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &RPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}

	w.WriteHeader(http.StatusOK) // JSON-RPC errors still return 200
	json.NewEncoder(w).Encode(response)
	s.logger.Warning("RPC error", "code", code, "message", message, "data", data)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	health := map[string]interface{}{
		"status":      "healthy",
		"timestamp":   time.Now().Unix(),
		"block_height": s.blockchain.GetBlockNumber().Uint64(),
		"peer_count":  0, // This would be updated with actual peer count
		"mempool_size": s.mempool.Size(),
	}

	json.NewEncoder(w).Encode(health)
}

// handleStats handles statistics requests
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	stats := map[string]interface{}{
		"block_height":    s.blockchain.GetBlockNumber().Uint64(),
		"mempool_size":    s.mempool.Size(),
		"mempool_stats":   s.mempool.GetStats(),
		"rpc_config": map[string]interface{}{
			"host":            s.config.Host,
			"port":            s.config.Port,
			"cors_origins":    s.config.CORSOrigins,
			"max_connections": s.config.MaxConnections,
		},
	}

	json.NewEncoder(w).Encode(stats)
}

// registerMethods registers all RPC methods
func (s *Server) registerMethods() {
	// Blockchain methods
	s.methods["eth_blockNumber"] = s.ethBlockNumber
	s.methods["eth_getBalance"] = s.ethGetBalance
	s.methods["eth_getTransactionCount"] = s.ethGetTransactionCount
	s.methods["eth_sendRawTransaction"] = s.ethSendRawTransaction
	s.methods["eth_getBlockByHash"] = s.ethGetBlockByHash
	s.methods["eth_getBlockByNumber"] = s.ethGetBlockByNumber
	s.methods["eth_getTransactionByHash"] = s.ethGetTransactionByHash
	s.methods["eth_getTransactionReceipt"] = s.ethGetTransactionReceipt
	s.methods["eth_call"] = s.ethCall
	s.methods["eth_estimateGas"] = s.ethEstimateGas
	s.methods["eth_gasPrice"] = s.ethGasPrice
	s.methods["eth_chainId"] = s.ethChainId
	
	// Network methods
	s.methods["net_version"] = s.netVersion
	s.methods["net_listening"] = s.netListening
	s.methods["net_peerCount"] = s.netPeerCount
	
	// Custom methods
	s.methods["lumina_getBlockNumber"] = s.ethBlockNumber
	s.methods["lumina_getBalance"] = s.ethGetBalance
	s.methods["lumina_sendRawTransaction"] = s.ethSendRawTransaction
	s.methods["lumina_getMempoolSize"] = s.luminaGetMempoolSize
	s.methods["lumina_getStats"] = s.luminaGetStats
}

// RPC method implementations

func (s *Server) ethBlockNumber(params interface{}) (interface{}, error) {
	blockNumber := s.blockchain.GetBlockNumber()
	return crypto.EncodeBig(blockNumber), nil
}

func (s *Server) ethGetBalance(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressStr, ok := paramList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	address := crypto.HexToAddress(addressStr)
	
	// For now, return zero balance (implement with state DB integration)
	balance := big.NewInt(0)
	
	return crypto.EncodeBig(balance), nil
}

func (s *Server) ethGetTransactionCount(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	addressStr, ok := paramList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid address parameter")
	}

	address := crypto.HexToAddress(addressStr)
	
	// For now, return zero nonce (implement with state DB integration)
	nonce := uint64(0)
	
	// Check mempool for pending transactions
	pendingTxs := s.mempool.GetTransactionsByFrom(address)
	nonce += uint64(len(pendingTxs))
	
	return crypto.EncodeUint64(nonce), nil
}

func (s *Server) ethSendRawTransaction(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	txDataStr, ok := paramList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid transaction data parameter")
	}

	// For now, return a mock transaction hash
	// In a real implementation, decode the transaction and add to mempool
	txHash := crypto.HexToHash(fmt.Sprintf("0x%x", time.Now().UnixNano()))
	
	s.logger.Info("Raw transaction received", "data", txDataStr, "hash", txHash.Hex())
	
	return txHash.Hex(), nil
}

func (s *Server) ethGetBlockByHash(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	hashStr, ok := paramList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hash parameter")
	}

	hash := crypto.HexToHash(hashStr)
	block, err := s.blockchain.GetBlockByHash(hash)
	if err != nil {
		return nil, nil // Return null for non-existent blocks
	}

	return s.formatBlock(block), nil
}

func (s *Server) ethGetBlockByNumber(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	var blockNumber *big.Int
	
	switch v := paramList[0].(type) {
	case string:
		if v == "latest" {
			blockNumber = s.blockchain.GetBlockNumber()
		} else if v == "earliest" {
			blockNumber = big.NewInt(0)
		} else if v == "pending" {
			blockNumber = s.blockchain.GetBlockNumber()
		} else {
			var err error
			blockNumber, err = crypto.DecodeBig(v)
			if err != nil {
				return nil, fmt.Errorf("invalid block number: %v", err)
			}
		}
	case float64:
		blockNumber = big.NewInt(int64(v))
	default:
		return nil, fmt.Errorf("invalid block number parameter")
	}

	block, err := s.blockchain.GetBlockByNumber(blockNumber)
	if err != nil {
		return nil, nil // Return null for non-existent blocks
	}

	return s.formatBlock(block), nil
}

func (s *Server) ethGetTransactionByHash(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	hashStr, ok := paramList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hash parameter")
	}

	hash := crypto.HexToHash(hashStr)
	
	// Check mempool first
	if tx := s.mempool.GetTransaction(hash); tx != nil {
		return s.formatTransaction(tx, nil, 0), nil
	}

	// TODO: Check blockchain for confirmed transactions
	
	return nil, nil // Return null for non-existent transactions
}

func (s *Server) ethGetTransactionReceipt(params interface{}) (interface{}, error) {
	paramList, ok := params.([]interface{})
	if !ok || len(paramList) < 1 {
		return nil, fmt.Errorf("invalid parameters")
	}

	hashStr, ok := paramList[0].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hash parameter")
	}

	// TODO: Implement transaction receipt lookup
	// For now, return null
	return nil, nil
}

func (s *Server) ethCall(params interface{}) (interface{}, error) {
	// TODO: Implement contract call simulation
	return "0x", nil
}

func (s *Server) ethEstimateGas(params interface{}) (interface{}, error) {
	// Return default gas estimate
	return crypto.EncodeUint64(21000), nil
}

func (s *Server) ethGasPrice(params interface{}) (interface{}, error) {
	gasPrice := big.NewInt(1000000000) // 1 Gwei
	return crypto.EncodeBig(gasPrice), nil
}

func (s *Server) ethChainId(params interface{}) (interface{}, error) {
	chainId := big.NewInt(1337) // Default chain ID
	return crypto.EncodeBig(chainId), nil
}

func (s *Server) netVersion(params interface{}) (interface{}, error) {
	return "1337", nil
}

func (s *Server) netListening(params interface{}) (interface{}, error) {
	return true, nil
}

func (s *Server) netPeerCount(params interface{}) (interface{}, error) {
	return crypto.EncodeUint64(0), nil // TODO: Get actual peer count
}

func (s *Server) luminaGetMempoolSize(params interface{}) (interface{}, error) {
	return s.mempool.Size(), nil
}

func (s *Server) luminaGetStats(params interface{}) (interface{}, error) {
	stats := map[string]interface{}{
		"block_height":  s.blockchain.GetBlockNumber().Uint64(),
		"mempool_size":  s.mempool.Size(),
		"mempool_stats": s.mempool.GetStats(),
	}
	return stats, nil
}

// Helper methods for formatting responses

func (s *Server) formatBlock(block *core.Block) map[string]interface{} {
	return map[string]interface{}{
		"number":           crypto.EncodeBig(block.Header.Number),
		"hash":             block.Hash.Hex(),
		"parentHash":       block.Header.PreviousHash.Hex(),
		"nonce":            crypto.EncodeUint64(block.Header.Nonce),
		"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000000",
		"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"logsBloom":        "0x" + string(block.Header.LogsBloom[:]),
		"transactionsRoot": block.Header.TransactionsRoot.Hex(),
		"stateRoot":        block.Header.StateRoot.Hex(),
		"receiptsRoot":     block.Header.ReceiptsRoot.Hex(),
		"miner":            block.Header.Coinbase.Hex(),
		"difficulty":       crypto.EncodeBig(block.Header.Difficulty),
		"totalDifficulty":  crypto.EncodeBig(block.Header.Difficulty), // Simplified
		"extraData":        crypto.Encode(block.Header.ExtraData),
		"size":             crypto.EncodeUint64(1000), // Estimated
		"gasLimit":         crypto.EncodeUint64(block.Header.GasLimit),
		"gasUsed":          crypto.EncodeUint64(block.Header.GasUsed),
		"timestamp":        crypto.EncodeUint64(block.Header.Timestamp),
		"transactions":     s.formatTransactions(block.Transactions, &block.Hash),
		"uncles":           []string{},
	}
}

func (s *Server) formatTransactions(txs []*core.Transaction, blockHash *crypto.Hash) []interface{} {
	result := make([]interface{}, len(txs))
	for i, tx := range txs {
		result[i] = s.formatTransaction(tx, blockHash, uint64(i))
	}
	return result
}

func (s *Server) formatTransaction(tx *core.Transaction, blockHash *crypto.Hash, index uint64) map[string]interface{} {
	result := map[string]interface{}{
		"hash":             tx.Hash.Hex(),
		"nonce":            crypto.EncodeUint64(tx.Nonce),
		"blockHash":        nil,
		"blockNumber":      nil,
		"transactionIndex": nil,
		"from":             tx.From.Hex(),
		"to":               nil,
		"value":            crypto.EncodeBig(tx.Value),
		"gasPrice":         crypto.EncodeBig(tx.GasPrice),
		"gas":              crypto.EncodeUint64(tx.GasLimit),
		"input":            crypto.Encode(tx.Data),
		"v":                crypto.EncodeBig(tx.V),
		"r":                crypto.EncodeBig(tx.R),
		"s":                crypto.EncodeBig(tx.S),
	}

	if tx.To != nil {
		result["to"] = tx.To.Hex()
	}

	if blockHash != nil {
		result["blockHash"] = blockHash.Hex()
		result["transactionIndex"] = crypto.EncodeUint64(index)
		// Block number would need to be looked up
	}

	return result
}
