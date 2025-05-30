
# 🌟 Lumina Blockchain Node

A production-ready, professional blockchain node implementation written in Go with custom consensus mechanism, execution environment, and comprehensive tooling.

## ✨ Features

### 🔗 Core Blockchain Features
- **Custom Proof-of-Work Consensus**: Complete PoW implementation from scratch
- **World State Management**: Custom StateDB with Patricia Merkle Trie-inspired structure
- **Transaction Execution Engine**: Custom execution environment with gas mechanics
- **Block Processing**: Full block validation, mining, and propagation
- **Account Model**: Ethereum-compatible account structure with nonce, balance, and code storage

### 🚀 Production-Ready Architecture
- **Modular Design**: Clean separation of concerns across packages
- **Comprehensive Logging**: Structured logging with multiple output formats
- **Metrics & Monitoring**: Prometheus-compatible metrics endpoint
- **Configuration Management**: Flexible YAML-based configuration
- **Error Handling**: Robust error handling throughout the codebase
- **Database Layer**: LevelDB integration with abstraction layer

### 🌐 Network & Communication
- **P2P Networking**: Full peer-to-peer protocol implementation
- **Message Broadcasting**: Efficient block and transaction propagation
- **Peer Discovery**: Automatic peer discovery and connection management
- **Network Protocol**: Custom message types with version negotiation

### 📡 JSON-RPC API
- **Ethereum Compatibility**: Standard eth_* RPC methods
- **Custom Methods**: Lumina-specific RPC endpoints
- **CORS Support**: Configurable CORS for web applications
- **Health Checks**: Built-in health and statistics endpoints

### ⛏️ Mining Capabilities
- **Multi-threaded Mining**: Configurable number of mining threads
- **Dynamic Difficulty**: Automatic difficulty adjustment
- **Mining Metrics**: Real-time hash rate and performance monitoring
- **Reward System**: Configurable mining rewards

### 🛠 Developer Tools
- **CLI Interface**: Comprehensive command-line interface
- **Wallet Management**: Built-in wallet creation and management
- **Transaction Tools**: Send transactions with custom data
- **Status Monitoring**: Real-time node status and metrics

### 🏦 Wallet Extension
- **Chrome Extension**: Professional MetaMask-compatible wallet extension
- **Web3 Provider**: Full EIP-1193 provider implementation
- **DApp Integration**: Seamless integration with decentralized applications
- **Security Features**: Secure key management and transaction signing

## 🏗 Architecture

### Directory Structure
```
blockchain-node/
├── cmd/                     # Command-line interface
│   ├── cli/                 # CLI commands and parsing
│   └── main.go             # Application entry point
├── core/                    # Core blockchain components
│   ├── blockchain.go        # Blockchain implementation
│   ├── types.go            # Core data structures
│   ├── execution.go        # Transaction execution engine
│   └── statedb.go          # World state management
├── consensus/               # Consensus mechanisms
│   └── pow.go              # Proof-of-Work implementation
├── p2p/                    # Peer-to-peer networking
│   └── server.go           # P2P server implementation
├── rpc/                    # JSON-RPC server
│   └── server.go           # RPC server with Ethereum compatibility
├── mempool/                # Transaction pool
│   └── mempool.go          # Mempool implementation with prioritization
├── storage/                # Database layer
│   └── database.go         # LevelDB wrapper with abstraction
├── crypto/                 # Cryptographic functions
│   └── keys.go             # Key generation and signing
├── config/                 # Configuration management
│   └── config.go           # YAML configuration parsing
├── logger/                 # Logging system
│   └── logger.go           # Structured logging implementation
├── metrics/                # Metrics and monitoring
│   └── metrics.go          # Prometheus-compatible metrics
├── node/                   # Node orchestration
│   └── node.go             # Main node coordinator
├── evm/                    # EVM compatibility layer
│   └── statedb.go          # EVM StateDB adapter
├── wallet-extension/       # Chrome wallet extension
│   ├── manifest.json       # Extension manifest
│   ├── popup.html          # Wallet interface
│   ├── background.js       # Service worker
│   ├── content.js          # Content script
│   ├── inpage.js           # Web3 provider
│   └── ...                 # Additional extension files
└── docs/                   # Documentation
    ├── mining.md           # Mining guide
    ├── rpc.md              # RPC API documentation
    └── metrics.md          # Metrics guide
```

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git
- Make (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd blockchain-node
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Build the application**
   ```bash
   go build -o lumina-node cmd/main.go
   ```

### Basic Usage

1. **Start a node**
   ```bash
   ./lumina-node startnode
   ```

2. **Start with mining enabled**
   ```bash
   ./lumina-node startnode --mining --rpc
   ```

3. **Create a wallet**
   ```bash
   ./lumina-node createwallet
   ```

4. **Check balance**
   ```bash
   ./lumina-node getbalance 0x1234...
   ```

5. **Send transaction**
   ```bash
   ./lumina-node send --from 0x1234... --to 0x5678... --amount 1.5
   ```

### Configuration

The node uses YAML configuration files. Create a `.blockchain-node.yaml` file in your home directory or working directory:

```yaml
network:
  port: 8080
  max_peers: 50
  
rpc:
  enabled: true
  port: 8545
  host: "localhost"
  
mining:
  enabled: true
  threads: 4
  difficulty: 4
  
logging:
  level: "info"
  output: "both"
  file_path: "./logs/node.log"
```

## 📖 Documentation

### Core Components

#### Blockchain Engine
The blockchain engine manages the chain state, validates blocks, and processes transactions. It uses a custom execution environment that interprets transaction data and modifies account states accordingly.

#### Consensus Mechanism
Implements Proof-of-Work consensus with SHA256 hashing. The difficulty adjusts automatically based on block time targets, ensuring consistent block production.

#### State Management
Uses a custom StateDB that maintains account states, contract storage, and transaction logs. The state root is calculated using merkle tree structures for integrity verification.

#### Transaction Processing
Transactions are validated for signature correctness, nonce sequencing, and sufficient balance before execution. The execution engine handles value transfers, contract creation, and contract calls.

### API Reference

#### JSON-RPC Methods

**Ethereum Compatible Methods:**
- `eth_blockNumber` - Get current block number
- `eth_getBalance` - Get account balance
- `eth_getTransactionCount` - Get account nonce
- `eth_sendRawTransaction` - Submit raw transaction
- `eth_getBlockByHash` - Get block by hash
- `eth_getBlockByNumber` - Get block by number
- `eth_getTransactionByHash` - Get transaction by hash
- `eth_getTransactionReceipt` - Get transaction receipt
- `eth_call` - Simulate transaction call
- `eth_estimateGas` - Estimate gas for transaction
- `eth_gasPrice` - Get current gas price
- `eth_chainId` - Get chain ID

**Custom Methods:**
- `lumina_getStats` - Get node statistics
- `lumina_getMempoolSize` - Get mempool size

### Mining Guide

1. **Enable Mining**
   ```bash
   ./lumina-node startnode --mining
   ```

2. **Configure Mining**
   ```yaml
   mining:
     enabled: true
     threads: 4
     difficulty: 4
     address: "0x..." # Optional mining reward address
   ```

3. **Monitor Mining**
   - Check logs for mining progress
   - Use metrics endpoint for hash rate
   - Monitor block production rate

### Wallet Extension Setup

1. **Load Extension**
   - Open Chrome and go to `chrome://extensions/`
   - Enable "Developer mode"
   - Click "Load unpacked" and select `wallet-extension/` folder

2. **Connect to DApps**
   The extension provides a MetaMask-compatible Web3 provider that DApps can use to interact with the Lumina blockchain.

## 🔧 Development

### Building from Source

```bash
# Clone repository
git clone <repository-url>
cd blockchain-node

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o lumina-node cmd/main.go

# Run with development settings
./lumina-node startnode --log-level debug --mining
```

### Adding Custom Features

1. **Custom Transaction Types**
   Extend the execution engine in `core/execution.go` to handle new transaction data formats.

2. **New RPC Methods**
   Add methods to `rpc/server.go` and register them in the `registerMethods()` function.

3. **Enhanced Consensus**
   Modify `consensus/pow.go` to implement different difficulty algorithms or consensus mechanisms.

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./core
```

## 📊 Monitoring & Metrics

### Prometheus Metrics

The node exposes Prometheus-compatible metrics at `/metrics`:

- `lumina_block_height` - Current block height
- `lumina_total_transactions` - Total transactions processed
- `lumina_mempool_size` - Current mempool size
- `lumina_peer_count` - Number of connected peers
- `lumina_hash_rate` - Current mining hash rate
- `lumina_uptime_seconds` - Node uptime

### Health Checks

Health endpoint at `/health` provides:
- Node status
- Current block height
- Peer count
- Mempool size

### Logging

Structured logging with configurable levels and outputs:
- Console output with colors
- File output with rotation
- JSON format for log aggregation

## 🛡 Security

### Cryptographic Security
- ECDSA signatures using secp256k1 curve
- Keccak256 hashing for Ethereum compatibility
- SHA256 for Proof-of-Work mining
- Secure random number generation

### Network Security
- Peer authentication and verification
- Message validation and sanitization
- Rate limiting for RPC endpoints
- CORS configuration for web security

### Wallet Security
- Private key encryption in wallet extension
- Secure key storage in browser extension
- Transaction signing with user confirmation
- Network validation for transaction safety

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards
- Follow Go best practices and conventions
- Write comprehensive tests for new features
- Update documentation for API changes
- Use structured logging for debugging
- Handle errors gracefully with proper context

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- [Mining Guide](docs/mining.md)
- [RPC API Documentation](docs/rpc.md)
- [Metrics Guide](docs/metrics.md)
- [Wallet Extension Guide](wallet-extension/README.md)

## 🎯 Roadmap

### Phase 1: Core Infrastructure ✅
- [x] Basic blockchain implementation
- [x] Proof-of-Work consensus
- [x] P2P networking
- [x] JSON-RPC API
- [x] Wallet extension

### Phase 2: Enhanced Features
- [ ] Smart contract virtual machine
- [ ] Advanced transaction types
- [ ] Stake-based consensus option
- [ ] Cross-chain bridges

### Phase 3: Ecosystem Tools
- [ ] Block explorer web interface
- [ ] Development framework
- [ ] Testing suite
- [ ] Deployment tools

### Phase 4: Enterprise Features
- [ ] Permissioned networks
- [ ] Advanced monitoring
- [ ] High availability setup
- [ ] Performance optimization

---

**Built with ❤️ for the blockchain community**

For questions, issues, or contributions, please visit our [GitHub repository](https://github.com/lumina-blockchain/node).
