
# 🌐 JSON-RPC API Documentation

Dokumentasi lengkap untuk JSON-RPC API yang kompatibel dengan Ethereum dan endpoint khusus blockchain node.

## 📋 Daftar Isi

- [Overview](#overview)
- [Endpoint Configuration](#endpoint-configuration)
- [Authentication](#authentication)
- [Standard Ethereum RPC Methods](#standard-ethereum-rpc-methods)
- [Custom Blockchain Methods](#custom-blockchain-methods)
- [WebSocket API](#websocket-api)
- [Error Handling](#error-handling)
- [Rate Limiting](#rate-limiting)
- [Examples](#examples)

## 🔍 Overview

Blockchain node menyediakan JSON-RPC API yang kompatibel dengan standar Ethereum, memungkinkan interaksi dengan blockchain menggunakan tools dan libraries yang sudah familiar.

### Supported Protocols
- **HTTP/HTTPS**: Standard JSON-RPC over HTTP
- **WebSocket**: Real-time subscriptions dan event streaming
- **IPC**: Unix socket untuk local connections (planned)

### API Versions
- **JSON-RPC 2.0**: Primary protocol
- **Ethereum JSON-RPC**: Full compatibility dengan eth_* methods

## ⚙️ Endpoint Configuration

### Default Configuration

```yaml
# .blockchain-node.yaml
rpc:
  enabled: true
  host: "localhost"
  port: 8545
  cors_origins: ["*"]
  max_connections: 100
  timeout: 30
  rate_limit: 1000        # requests per minute
  
  # Security settings
  auth_required: false
  api_keys: []
  
  # Feature flags
  enable_websocket: true
  enable_debug_api: false
  enable_admin_api: false
```

### Custom Configuration

```bash
# Start dengan custom RPC config
./blockchain-node startnode \
  --rpc-host 0.0.0.0 \
  --rpc-port 8545 \
  --rpc-cors-origins "https://mydapp.com,https://remix.ethereum.org"
```

### TLS/HTTPS Setup

```yaml
rpc:
  tls_enabled: true
  tls_cert: "/path/to/cert.pem"
  tls_key: "/path/to/key.pem"
```

## 🔐 Authentication

### API Key Authentication

```yaml
rpc:
  auth_required: true
  api_keys:
    - "your-secret-api-key-1"
    - "your-secret-api-key-2"
```

Request dengan API key:
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-secret-api-key-1" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

### JWT Token Authentication

```bash
# Get JWT token
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"username":"admin","password":"secret"}' \
  http://localhost:8545/auth/login

# Use token in requests
curl -X POST \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..." \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

## 📊 Standard Ethereum RPC Methods

### Block Methods

#### eth_blockNumber
Returns the current block number.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": "0x4b7",
  "id": 1
}
```

#### eth_getBlockByNumber
Returns block information by block number.

**Parameters:**
- `blockNumber` (string): Block number in hex, or "earliest", "latest", "pending"
- `fullTransactions` (boolean): If true, returns full transaction objects

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "number": "0x1b4",
    "hash": "0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae",
    "parentHash": "0xe99e022112df268ce40efe44c8dd1b2c1c2d71b1ec3fb2b0dd1ca3bd10b72d8c",
    "nonce": "0x9bd90d1fb4d9c0",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "stateRoot": "0xd5855eb08b3387c0af375e9cdb6acfc05eb8f519e419b874b6ff2ffda7ed1dff",
    "miner": "0x4e65fda2159562a496f9f3522f89122a3088497a",
    "difficulty": "0x027f07",
    "totalDifficulty": "0x027f07",
    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "size": "0x027f07",
    "gasLimit": "0x9f759",
    "gasUsed": "0x9f759",
    "timestamp": "0x54e34e8e",
    "transactions": [...],
    "uncles": []
  },
  "id": 1
}
```

#### eth_getBlockByHash
Returns block information by block hash.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0xdc0818cf78f21a8e70579cb46a43643f78291264dda342ae31049421c82d21ae", true],"id":1}' \
  http://localhost:8545
```

### Account Methods

#### eth_getBalance
Returns the balance of an account.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": "0x0234c8a3397aab58",
  "id": 1
}
```

#### eth_getTransactionCount
Returns the number of transactions sent from an address (nonce).

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1","latest"],"id":1}' \
  http://localhost:8545
```

#### eth_getCode
Returns the code at a given address.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getCode","params":["0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b", "0x2"],"id":1}' \
  http://localhost:8545
```

### Transaction Methods

#### eth_sendRawTransaction
Sends a signed transaction to the network.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": "0xe670ec64341771606e55d6b4ca35a1a6b75ee3d5145a99d05921026d1527331",
  "id": 1
}
```

#### eth_getTransactionByHash
Returns transaction information by hash.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"],"id":1}' \
  http://localhost:8545
```

#### eth_getTransactionReceipt
Returns the receipt of a transaction.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238"],"id":1}' \
  http://localhost:8545
```

### Contract Methods

#### eth_call
Executes a message call (read-only) against the blockchain.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_call","params":[{"to":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","data":"0x0000000000000000000000000000000000000000000000000000000000000001"},"latest"],"id":1}' \
  http://localhost:8545
```

#### eth_estimateGas
Returns an estimate of gas for a transaction.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{"to":"0xb60e8dd61c5d32be8058bb8eb970870f07233155","data":"0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"}],"id":1}' \
  http://localhost:8545
```

### Network Methods

#### eth_chainId
Returns the chain ID of the blockchain.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
  http://localhost:8545
```

#### eth_gasPrice
Returns the current gas price.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_gasPrice","params":[],"id":1}' \
  http://localhost:8545
```

#### net_version
Returns the network ID.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"net_version","params":[],"id":1}' \
  http://localhost:8545
```

## 🔧 Custom Blockchain Methods

### Node Information

#### blockchain_getNodeInfo
Returns comprehensive node information.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_getNodeInfo","params":[],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "version": "1.0.0",
    "network": "mainnet",
    "chainId": 1337,
    "currentBlock": 12345,
    "highestBlock": 12345,
    "peerCount": 8,
    "mining": true,
    "hashRate": "1000000",
    "uptime": "2d 3h 45m"
  },
  "id": 1
}
```

#### blockchain_getPeers
Returns list of connected peers.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_getPeers","params":[],"id":1}' \
  http://localhost:8545
```

### Mining Methods

#### blockchain_getMiningInfo
Returns current mining information.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_getMiningInfo","params":[],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "mining": true,
    "hashRate": 1000000,
    "difficulty": "0x4",
    "blocksMined": 145,
    "minerAddress": "0x1234567890abcdef...",
    "threads": 4,
    "efficiency": 98.5
  },
  "id": 1
}
```

#### blockchain_startMining
Starts mining operation.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_startMining","params":["0x1234567890abcdef...", 4],"id":1}' \
  http://localhost:8545
```

#### blockchain_stopMining
Stops mining operation.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_stopMining","params":[],"id":1}' \
  http://localhost:8545
```

### Mempool Methods

#### blockchain_getMempoolInfo
Returns mempool information.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_getMempoolInfo","params":[],"id":1}' \
  http://localhost:8545
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "size": 1024,
    "bytes": 256000,
    "usage": 512000,
    "maxMempool": 300000000,
    "mempoolMinFee": 1000000000
  },
  "id": 1
}
```

#### blockchain_getMempoolTransactions
Returns list of transactions in mempool.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"blockchain_getMempoolTransactions","params":[10],"id":1}' \
  http://localhost:8545
```

## 🔌 WebSocket API

### Connection

```javascript
const ws = new WebSocket('ws://localhost:8545');

ws.onopen = function() {
    console.log('Connected to blockchain node');
    
    // Subscribe to new blocks
    ws.send(JSON.stringify({
        jsonrpc: "2.0",
        method: "eth_subscribe",
        params: ["newHeads"],
        id: 1
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('Received:', data);
};
```

### Subscriptions

#### eth_subscribe (newHeads)
Subscribe to new block headers.

```javascript
ws.send(JSON.stringify({
    jsonrpc: "2.0",
    method: "eth_subscribe",
    params: ["newHeads"],
    id: 1
}));
```

#### eth_subscribe (logs)
Subscribe to event logs.

```javascript
ws.send(JSON.stringify({
    jsonrpc: "2.0",
    method: "eth_subscribe",
    params: ["logs", {
        address: "0x1234567890abcdef...",
        topics: ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]
    }],
    id: 1
}));
```

#### eth_subscribe (newPendingTransactions)
Subscribe to pending transactions.

```javascript
ws.send(JSON.stringify({
    jsonrpc: "2.0",
    method: "eth_subscribe",
    params: ["newPendingTransactions"],
    id: 1
}));
```

### Unsubscribe

```javascript
ws.send(JSON.stringify({
    jsonrpc: "2.0",
    method: "eth_unsubscribe",
    params: ["0x1234..."], // subscription ID
    id: 1
}));
```

## ❌ Error Handling

### Standard JSON-RPC Errors

| Code | Message | Description |
|------|---------|-------------|
| -32700 | Parse error | Invalid JSON |
| -32600 | Invalid Request | JSON is not a valid request object |
| -32601 | Method not found | Method does not exist |
| -32602 | Invalid params | Invalid method parameters |
| -32603 | Internal error | Internal JSON-RPC error |

### Custom Errors

| Code | Message | Description |
|------|---------|-------------|
| -32000 | Server error | Generic server error |
| -32001 | Resource not found | Block/transaction not found |
| -32002 | Resource unavailable | Node syncing |
| -32003 | Transaction rejected | Invalid transaction |
| -32004 | Method not supported | Method disabled |
| -32005 | Limit exceeded | Request exceeds limit |

### Error Response Format

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32602,
    "message": "Invalid params",
    "data": "Missing required parameter: address"
  },
  "id": 1
}
```

## 🚦 Rate Limiting

### Default Limits

```yaml
rpc:
  rate_limit: 1000        # requests per minute
  burst_limit: 100        # burst requests
  per_ip_limit: 100       # per IP address
```

### Custom Rate Limits per Method

```yaml
rpc:
  method_limits:
    eth_getBalance: 60      # 60 requests per minute
    eth_call: 120          # 120 requests per minute
    eth_sendRawTransaction: 10  # 10 requests per minute
```

### Rate Limit Headers

Response headers include rate limit information:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## 📝 Examples

### JavaScript/Node.js

```javascript
const axios = require('axios');

const rpcUrl = 'http://localhost:8545';

async function getBlockNumber() {
    const response = await axios.post(rpcUrl, {
        jsonrpc: '2.0',
        method: 'eth_blockNumber',
        params: [],
        id: 1
    });
    
    return parseInt(response.data.result, 16);
}

async function getBalance(address) {
    const response = await axios.post(rpcUrl, {
        jsonrpc: '2.0',
        method: 'eth_getBalance',
        params: [address, 'latest'],
        id: 1
    });
    
    return parseInt(response.data.result, 16);
}

// Usage
getBlockNumber().then(blockNumber => {
    console.log('Current block:', blockNumber);
});

getBalance('0x1234567890abcdef...').then(balance => {
    console.log('Balance:', balance, 'wei');
});
```

### Python

```python
import requests
import json

rpc_url = 'http://localhost:8545'

def rpc_call(method, params=[]):
    payload = {
        'jsonrpc': '2.0',
        'method': method,
        'params': params,
        'id': 1
    }
    
    response = requests.post(rpc_url, json=payload)
    return response.json()

# Get current block number
result = rpc_call('eth_blockNumber')
block_number = int(result['result'], 16)
print(f'Current block: {block_number}')

# Get balance
result = rpc_call('eth_getBalance', ['0x1234567890abcdef...', 'latest'])
balance = int(result['result'], 16)
print(f'Balance: {balance} wei')
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type JSONRPCRequest struct {
    JSONRPC string        `json:"jsonrpc"`
    Method  string        `json:"method"`
    Params  []interface{} `json:"params"`
    ID      int           `json:"id"`
}

type JSONRPCResponse struct {
    JSONRPC string      `json:"jsonrpc"`
    Result  interface{} `json:"result"`
    Error   *RPCError   `json:"error"`
    ID      int         `json:"id"`
}

type RPCError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func rpcCall(method string, params []interface{}) (*JSONRPCResponse, error) {
    request := JSONRPCRequest{
        JSONRPC: "2.0",
        Method:  method,
        Params:  params,
        ID:      1,
    }
    
    jsonData, err := json.Marshal(request)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.Post("http://localhost:8545", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var response JSONRPCResponse
    err = json.NewDecoder(resp.Body).Decode(&response)
    return &response, err
}

func main() {
    // Get block number
    response, err := rpcCall("eth_blockNumber", []interface{}{})
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Current block: %s\n", response.Result)
}
```

### curl Examples

```bash
# Get current block number
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# Get balance
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x1234567890abcdef...","latest"],"id":1}' \
  http://localhost:8545

# Send transaction
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf86c..."],"id":1}' \
  http://localhost:8545

# Call contract method
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_call","params":[{"to":"0x1234...","data":"0xabcd..."},"latest"],"id":1}' \
  http://localhost:8545
```

## 🔧 Debugging & Testing

### Enable Debug API

```yaml
rpc:
  enable_debug_api: true
```

### Debug Methods

#### debug_traceTransaction
Traces transaction execution.

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"debug_traceTransaction","params":["0x1234..."],"id":1}' \
  http://localhost:8545
```

### Testing Tools

```bash
# Test RPC connectivity
curl -f http://localhost:8545 || echo "RPC server not responding"

# Benchmark RPC performance
ab -n 1000 -c 10 -T application/json -p request.json http://localhost:8545/

# Monitor RPC logs
tail -f logs/rpc.log | grep -E "(ERROR|WARN)"
```

## 📚 Best Practices

### Client Implementation

1. **Connection Pooling**: Reuse HTTP connections
2. **Timeout Handling**: Set appropriate timeouts
3. **Error Retry**: Implement exponential backoff
4. **Request Batching**: Batch multiple requests when possible

### Performance Optimization

1. **Use Latest Block**: Specify "latest" instead of specific block numbers when possible
2. **Minimal Data**: Request only necessary fields
3. **WebSocket for Real-time**: Use WebSocket for frequent updates
4. **Cache Results**: Cache block data and static contract calls

### Security

1. **Input Validation**: Always validate input parameters
2. **Rate Limiting**: Respect rate limits
3. **HTTPS**: Use HTTPS in production
4. **API Keys**: Use API keys for authentication

---

**Happy Coding!** 🚀

Untuk bantuan lebih lanjut atau melaporkan issues, silakan gunakan GitHub Issues atau bergabung dengan komunitas Discord.
