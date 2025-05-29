
# Blockchain Node with EVM Support

Implementasi blockchain node dengan dukungan Ethereum Virtual Machine (EVM) menggunakan bahasa pemrograman Go.

## Fitur Utama

- **Consensus**: Proof-of-Work (PoW)
- **EVM Support**: Integrasi dengan go-ethereum EVM
- **P2P Network**: Jaringan peer-to-peer untuk komunikasi antar node
- **JSON-RPC API**: Kompatibel dengan Ethereum JSON-RPC API
- **Mempool**: Transaction pool untuk transaksi pending
- **CLI Interface**: Command line interface untuk operasi node

## Struktur Proyek

```
blockchain-node/
├── cmd/                    # Command line interface
│   ├── main.go
│   └── cli/
├── config/                 # Konfigurasi sistem
├── core/                   # Tipe data dan struktur blockchain inti
├── consensus/              # Algoritma konsensus (PoW)
├── evm/                    # Integrasi Ethereum Virtual Machine
├── storage/                # Layer database dan penyimpanan
├── mempool/                # Transaction pool
├── p2p/                    # Jaringan peer-to-peer
├── rpc/                    # JSON-RPC API server
├── crypto/                 # Kriptografi dan manajemen kunci
└── node/                   # Node utama yang mengintegrasikan semua komponen
```

## Instalasi dan Penggunaan

### Prerequisites
- Go 1.21 atau lebih baru
- Git

### Instalasi
```bash
git clone <repository-url>
cd blockchain-node
go mod tidy
go build -o blockchain-node cmd/main.go
```

### Menjalankan Node

1. **Start Node**:
```bash
./blockchain-node startnode
```

2. **Create Wallet**:
```bash
./blockchain-node createwallet
```

3. **Get Balance**:
```bash
./blockchain-node getbalance [address]
```

4. **Send Transaction**:
```bash
./blockchain-node send -from [sender] -to [recipient] -amount [amount] -gaslimit [limit] -gasprice [price]
```

## Konfigurasi

Node dapat dikonfigurasi melalui file konfigurasi YAML. Contoh konfigurasi:

```yaml
network:
  port: 8080
  listen_addr: "0.0.0.0"
  max_peers: 50
  seed_nodes: []

rpc:
  enabled: true
  port: 8545
  host: "localhost"

mining:
  enabled: false
  threads: 1
  difficulty: 4

db:
  path: "./data"
  type: "leveldb"

evm:
  chain_id: 1337
  block_gas_limit: 8000000
```

## API Endpoints

### JSON-RPC API (Ethereum Compatible)
- `eth_blockNumber` - Mendapatkan nomor blok terbaru
- `eth_getBalance` - Mendapatkan balance alamat
- `eth_getBlockByNumber` - Mendapatkan blok berdasarkan nomor
- `eth_getBlockByHash` - Mendapatkan blok berdasarkan hash
- `eth_sendRawTransaction` - Mengirim raw transaction
- `eth_chainId` - Mendapatkan Chain ID
- `eth_gasPrice` - Mendapatkan gas price

### RESTful API
- `GET /eth/blockNumber` - Nomor blok terbaru
- `GET /eth/getBalance/:address` - Balance alamat
- `GET /eth/getBlockByNumber/:number` - Blok berdasarkan nomor
- `GET /eth/getBlockByHash/:hash` - Blok berdasarkan hash

## Arsitektur

### Komponen Utama

1. **Core**: Berisi struktur data inti seperti Block, Transaction, Account
2. **Consensus**: Implementasi Proof-of-Work untuk mining
3. **EVM**: Integrasi dengan go-ethereum untuk eksekusi smart contract
4. **Storage**: Layer abstraksi database menggunakan LevelDB
5. **P2P**: Server untuk komunikasi antar node
6. **RPC**: JSON-RPC API server
7. **Mempool**: Pool untuk transaksi pending

### Flow Eksekusi

1. Node menerima transaksi via RPC
2. Transaksi divalidasi dan ditambahkan ke mempool
3. Miner mengambil transaksi dari mempool untuk mining
4. Block baru di-mine menggunakan Proof-of-Work
5. Block valid ditambahkan ke blockchain
6. Block di-broadcast ke peer lain

## Development

### TODO Items
- Implementasi serialization/deserialization yang robust (RLP encoding)
- Implementasi Patricia Merkle Trie untuk state storage
- Integrasi penuh dengan go-ethereum EVM
- Implementasi access list dan transient storage
- Mekanisme suicide/selfdestruct untuk contract
- Implementasi signature verification yang lengkap
- Protocol P2P yang lebih sophisticated
- Wallet management system

### Testing
```bash
go test ./...
```

### Build
```bash
go build -o blockchain-node cmd/main.go
```

## Kontribusi

1. Fork repository
2. Buat branch untuk fitur baru
3. Commit perubahan
4. Push ke branch
5. Buat Pull Request

## Lisensi

MIT License - lihat file LICENSE untuk detail.
