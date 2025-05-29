
# Blockchain Node with EVM Support

Implementasi blockchain node enterprise-grade dengan dukungan Ethereum Virtual Machine (EVM) menggunakan bahasa pemrograman Go. Node ini dirancang untuk production dengan fitur logging komprehensif, monitoring, dan optimisasi performa.

## 🚀 Fitur Utama

### Core Blockchain Features
- **Consensus Algorithm**: Proof-of-Work (PoW) dengan difficulty adjustment
- **EVM Compatibility**: Integrasi penuh dengan go-ethereum EVM
- **Transaction Pool**: Mempool dengan validasi dan prioritas transaksi
- **World State Management**: Patricia Merkle Trie untuk state storage
- **Cryptography**: ECDSA signatures dengan Keccak256 hashing

### Network & Communication
- **P2P Network**: Jaringan peer-to-peer untuk komunikasi antar node
- **JSON-RPC API**: Kompatibel dengan Ethereum JSON-RPC API standard
- **RESTful Endpoints**: Additional REST API untuk kemudahan integrasi

### Production Features
- **Comprehensive Logging**: Multi-level logging (debug, info, warning, error)
- **Metrics & Monitoring**: Prometheus-compatible metrics
- **Graceful Shutdown**: Proper resource cleanup dan shutdown handling
- **Configuration Management**: YAML-based configuration dengan validation
- **CLI Interface**: Command line interface yang powerful

### Performance & Reliability
- **Database Optimization**: LevelDB dengan tuning untuk performa optimal
- **Memory Management**: Efficient memory usage dengan garbage collection optimization
- **Error Handling**: Comprehensive error handling di seluruh codebase
- **Resource Management**: Proper connection pooling dan resource cleanup

## 📁 Struktur Proyek

```
blockchain-node/
├── cmd/                    # Command line interface dan entry point
│   ├── main.go            # Main application entry point
│   └── cli/               # CLI commands dan flags
│       └── root.go        # Root command configuration
├── config/                 # Sistem konfigurasi
│   └── config.go          # Configuration structure dan loading
├── core/                   # Komponen blockchain inti
│   ├── types.go           # Data structures (Block, Transaction, etc.)
│   └── blockchain.go      # Blockchain logic dan state management
├── consensus/              # Algoritma konsensus
│   └── pow.go             # Proof-of-Work implementation
├── evm/                    # Ethereum Virtual Machine integration
│   └── statedb.go         # EVM StateDB interface implementation
├── storage/                # Database layer dan penyimpanan
│   └── database.go        # Database abstraction layer
├── mempool/                # Transaction pool management
│   └── mempool.go         # Transaction validation dan queuing
├── p2p/                    # Peer-to-peer networking
│   └── server.go          # P2P server dan message handling
├── rpc/                    # JSON-RPC API server
│   └── server.go          # RPC endpoint implementations
├── crypto/                 # Cryptographic operations
│   └── keys.go            # Key generation dan signature handling
├── logger/                 # Logging system
│   └── logger.go          # Structured logging implementation
├── metrics/                # Monitoring dan metrics
│   └── metrics.go         # Prometheus metrics collection
├── node/                   # Node orchestration
│   └── node.go            # Main node logic dan component integration
└── .blockchain-node.yaml  # Default configuration file
```

## 🛠 Instalasi dan Setup

### Prerequisites
- **Go**: Version 1.21 atau lebih baru
- **Git**: Untuk cloning repository
- **Make**: Untuk build automation (optional)

### Quick Start

1. **Clone Repository**:
```bash
git clone <repository-url>
cd blockchain-node
```

2. **Install Dependencies**:
```bash
go mod tidy
```

3. **Build Binary**:
```bash
go build -o blockchain-node cmd/main.go
```

4. **Run with Default Configuration**:
```bash
./blockchain-node startnode
```

### Advanced Installation

1. **Build dengan Optimizations**:
```bash
go build -ldflags="-s -w" -o blockchain-node cmd/main.go
```

2. **Cross-Platform Build**:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o blockchain-node-linux cmd/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o blockchain-node.exe cmd/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o blockchain-node-macos cmd/main.go
```

## 🔧 Konfigurasi

### File Konfigurasi

Node menggunakan file konfigurasi YAML (`.blockchain-node.yaml`) untuk mengatur semua parameter:

```yaml
# Network Configuration
network:
  port: 8080
  listen_addr: "0.0.0.0"
  max_peers: 50
  timeout: 30
  seed_nodes:
    - "127.0.0.1:8081"
    - "127.0.0.1:8082"

# JSON-RPC Configuration
rpc:
  enabled: true
  port: 8545
  host: "localhost"
  cors_origins: ["*"]
  max_connections: 100
  timeout: 30

# Mining Configuration
mining:
  enabled: false
  address: ""
  threads: 1
  difficulty: 4

# Database Configuration
db:
  path: "./data"
  type: "leveldb"
  cache_size: 64      # MB
  max_open_files: 1000
  write_buffer: 4     # MB

# EVM Configuration
evm:
  chain_id: 1337
  block_gas_limit: 8000000
  min_gas_price: 1000000000  # 1 Gwei

# Logging Configuration
logging:
  level: "info"              # debug, info, warning, error
  output: "both"             # console, file, both
  file_path: "./logs/blockchain.log"
  max_size: 100             # MB
  component: "blockchain-node"

# Metrics Configuration
metrics:
  enabled: true
  port: 8080
  path: "/metrics"
```

### Environment Variables

Anda juga dapat menggunakan environment variables untuk override konfigurasi:

```bash
export BLOCKCHAIN_NETWORK_PORT=8080
export BLOCKCHAIN_RPC_PORT=8545
export BLOCKCHAIN_MINING_ENABLED=true
export BLOCKCHAIN_LOGGING_LEVEL=debug
```

## 🎯 Penggunaan

### Command Line Interface

#### Menjalankan Node

```bash
# Start node dengan konfigurasi default
./blockchain-node startnode

# Start node dengan konfigurasi custom
./blockchain-node startnode --config custom-config.yaml

# Start node dengan logging debug
./blockchain-node startnode --log-level debug

# Start node dengan mining enabled
./blockchain-node startnode --mining-enabled --mining-address 0x1234...
```

#### Wallet Management

```bash
# Create new wallet
./blockchain-node createwallet

# Create wallet dengan custom nama
./blockchain-node createwallet --name "my-wallet"

# Import existing private key
./blockchain-node importwallet --private-key "0xabcd..."
```

#### Balance dan Transactions

```bash
# Check balance
./blockchain-node getbalance 0x1234567890abcdef1234567890abcdef12345678

# Send transaction
./blockchain-node send \
  --from 0x1234567890abcdef1234567890abcdef12345678 \
  --to 0xabcdef1234567890abcdef1234567890abcdef12 \
  --amount 1000000000000000000 \
  --gaslimit 21000 \
  --gasprice 1000000000

# Send contract deployment
./blockchain-node send \
  --from 0x1234567890abcdef1234567890abcdef12345678 \
  --data 0x608060405234801561001057600080fd5b50... \
  --gaslimit 1000000 \
  --gasprice 1000000000
```

#### Network Commands

```bash
# Get network info
./blockchain-node networkinfo

# List connected peers
./blockchain-node peers

# Add peer manually
./blockchain-node addpeer 127.0.0.1:8081
```

## 🌐 API Documentation

### JSON-RPC API (Ethereum Compatible)

#### Block Methods
```bash
# Get latest block number
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  http://localhost:8545

# Get block by number
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1", true],"id":1}' \
  http://localhost:8545

# Get block by hash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBlockByHash","params":["0xabcd...", true],"id":1}' \
  http://localhost:8545
```

#### Account Methods
```bash
# Get account balance
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x1234...","latest"],"id":1}' \
  http://localhost:8545

# Get transaction count (nonce)
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionCount","params":["0x1234...","latest"],"id":1}' \
  http://localhost:8545
```

#### Transaction Methods
```bash
# Send raw transaction
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_sendRawTransaction","params":["0xf86c..."],"id":1}' \
  http://localhost:8545

# Get transaction by hash
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xabcd..."],"id":1}' \
  http://localhost:8545

# Get transaction receipt
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0xabcd..."],"id":1}' \
  http://localhost:8545
```

#### Contract Methods
```bash
# Call contract method
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_call","params":[{"to":"0x1234...","data":"0xabcd..."},"latest"],"id":1}' \
  http://localhost:8545

# Estimate gas
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_estimateGas","params":[{"to":"0x1234...","data":"0xabcd..."}],"id":1}' \
  http://localhost:8545

# Get contract code
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_getCode","params":["0x1234...","latest"],"id":1}' \
  http://localhost:8545
```

### RESTful API

```bash
# Get latest block number
curl http://localhost:8545/eth/blockNumber

# Get account balance
curl http://localhost:8545/eth/getBalance/0x1234567890abcdef1234567890abcdef12345678

# Get block by number
curl http://localhost:8545/eth/getBlockByNumber/1

# Get block by hash
curl http://localhost:8545/eth/getBlockByHash/0xabcd...
```

### Metrics API

```bash
# Get Prometheus metrics
curl http://localhost:8080/metrics
```

## 🔍 Monitoring dan Debugging

### Logging

Node menyediakan logging berlevel dengan format JSON untuk production:

```bash
# Set log level
./blockchain-node startnode --log-level debug

# Log ke file
./blockchain-node startnode --log-output file --log-file ./custom.log

# Log ke console dan file
./blockchain-node startnode --log-output both
```

### Metrics

Metrics tersedia dalam format Prometheus di `/metrics` endpoint:

- `blockchain_blocks_total`: Total blocks processed
- `blockchain_transactions_total`: Total transactions processed
- `blockchain_peers_connected`: Number of connected peers
- `blockchain_mempool_size`: Current mempool size
- `blockchain_mining_difficulty`: Current mining difficulty

### Debug Mode

```bash
# Start dengan debug mode
./blockchain-node startnode --debug

# Enable verbose P2P logging
./blockchain-node startnode --debug --log-level debug
```

## 🏗 Development

### Struktur Development

```bash
# Install development dependencies
go mod download

# Run tests
go test ./...

# Run tests dengan coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linter
golangci-lint run

# Format code
go fmt ./...
```

### Build Automation

```bash
# Build untuk semua platform
make build-all

# Run tests
make test

# Clean build artifacts
make clean

# Install dependencies
make deps
```

### Testing Environment

```bash
# Start test network dengan 3 nodes
./scripts/start-testnet.sh

# Deploy test contracts
./scripts/deploy-contracts.sh

# Run integration tests
make test-integration
```

## 🔐 Security Considerations

### Production Deployment

1. **Network Security**:
   - Gunakan firewall untuk membatasi akses ke port RPC
   - Implementasikan TLS untuk komunikasi antar node
   - Whitelist peer nodes yang trusted

2. **Key Management**:
   - Simpan private keys dalam encrypted storage
   - Gunakan hardware security modules (HSM) untuk production
   - Implement key rotation policies

3. **Monitoring**:
   - Monitor metrics untuk anomaly detection
   - Setup alerting untuk critical events
   - Log audit trail untuk semua transactions

4. **Configuration**:
   - Gunakan environment variables untuk sensitive data
   - Validate semua input parameters
   - Implement rate limiting untuk API endpoints

## 📈 Performance Tuning

### Database Optimization

```yaml
db:
  cache_size: 128        # Increase untuk better read performance
  max_open_files: 2000   # Increase untuk high-load scenarios
  write_buffer: 8        # Increase untuk better write performance
```

### Network Optimization

```yaml
network:
  max_peers: 100         # Increase untuk better network connectivity
  timeout: 60            # Increase untuk slow networks
```

### EVM Optimization

```yaml
evm:
  block_gas_limit: 12000000  # Increase untuk more transactions per block
```

## 🚦 Roadmap

### Planned Features

- [ ] **Smart Contract Debugging**: Debug interface untuk contract execution
- [ ] **Sharding Support**: Horizontal scaling melalui sharding
- [ ] **Light Client Support**: SPV-like light client implementation
- [ ] **Cross-Chain Bridge**: Bridge ke Ethereum mainnet
- [ ] **Consensus Upgrade**: Migrasi ke Proof-of-Stake
- [ ] **Advanced P2P**: DHT-based peer discovery
- [ ] **State Pruning**: Historical state cleanup untuk disk efficiency
- [ ] **WebSocket API**: Real-time event subscriptions

### Current Limitations

- Single-threaded EVM execution
- In-memory state (tidak persistent untuk restarts)
- Basic P2P protocol (belum ada advanced features seperti fast sync)
- Limited smart contract debugging tools

## 🤝 Contributing

### Development Guidelines

1. **Code Style**: Ikuti Go standard formatting dengan `gofmt`
2. **Testing**: Semua code baru harus memiliki unit tests
3. **Documentation**: Update documentation untuk API changes
4. **Commit Messages**: Gunakan conventional commit format

### Pull Request Process

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Issue Reporting

Gunakan GitHub Issues dengan template yang tersedia:
- Bug reports
- Feature requests  
- Performance issues
- Documentation improvements

## 📄 License

Proyek ini dilisensikan di bawah MIT License - lihat file [LICENSE](LICENSE) untuk detail lengkap.

## 🙏 Acknowledgments

- **go-ethereum team**: Untuk EVM implementation yang excellent
- **LevelDB**: Untuk database storage yang reliable
- **Gin framework**: Untuk HTTP routing yang powerful
- **Cobra CLI**: Untuk command-line interface yang elegant

## 📞 Support

- **Documentation**: [Wiki](../../wiki)
- **Issues**: [GitHub Issues](../../issues)
- **Discussions**: [GitHub Discussions](../../discussions)
- **Email**: support@blockchain-node.dev

---

**Blockchain Node** - Enterprise-grade blockchain implementation dengan dukungan EVM penuh.
