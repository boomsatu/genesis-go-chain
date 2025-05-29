
# 🔨 Blockchain Mining Guide

Panduan lengkap untuk memulai mining di blockchain node dengan algoritma Proof-of-Work (PoW).

## 📋 Daftar Isi

- [Persyaratan Sistem](#persyaratan-sistem)
- [Konfigurasi Mining](#konfigurasi-mining)
- [Memulai Mining](#memulai-mining)
- [Monitoring Mining](#monitoring-mining)
- [Troubleshooting](#troubleshooting)
- [Optimisasi Performa](#optimisasi-performa)

## 🖥 Persyaratan Sistem

### Hardware Minimum
- **CPU**: 4 cores (Intel i5 atau AMD Ryzen 5)
- **RAM**: 8 GB
- **Storage**: 100 GB SSD
- **Network**: Koneksi internet stabil 10 Mbps

### Hardware Recommended
- **CPU**: 8+ cores (Intel i7/i9 atau AMD Ryzen 7/9)
- **RAM**: 16+ GB
- **Storage**: 500+ GB NVMe SSD
- **Network**: Koneksi internet stabil 50+ Mbps

### Software Requirements
- **OS**: Linux (Ubuntu 20.04+), macOS, atau Windows 10+
- **Go**: Version 1.21+
- **Git**: Untuk cloning repository

## ⚙️ Konfigurasi Mining

### 1. Setup Wallet untuk Mining

Sebelum memulai mining, Anda perlu membuat wallet untuk menerima reward:

```bash
# Buat wallet baru
./blockchain-node createwallet --name "mining-wallet"

# Atau import private key yang sudah ada
./blockchain-node importwallet --private-key "your-private-key-here"
```

Catat alamat wallet yang dihasilkan, karena akan digunakan untuk menerima mining rewards.

### 2. Konfigurasi File Mining

Edit file `.blockchain-node.yaml` untuk mengaktifkan mining:

```yaml
# Mining Configuration
mining:
  enabled: true                              # Aktifkan mining
  address: "0x1234567890abcdef..."          # Alamat wallet untuk reward
  threads: 4                                # Jumlah thread mining (sesuai CPU)
  difficulty: 4                             # Target difficulty (4-20)
  block_time: 15                            # Target waktu antar block (detik)
  gas_price_minimum: 1000000000             # Minimum gas price (1 Gwei)

# Network Configuration untuk Mining
network:
  port: 8080
  max_peers: 20                             # Lebih banyak peer = info block lebih cepat
  timeout: 30

# Database Configuration (optimasi untuk mining)
db:
  cache_size: 128                           # Increase cache untuk performa
  max_open_files: 2000
  write_buffer: 8
```

### 3. Optimisasi Sistem Operasi

#### Linux Optimization:
```bash
# Increase file descriptor limit
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# Optimize TCP settings
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
sysctl -p
```

#### Resource Monitoring:
```bash
# Install monitoring tools
sudo apt install htop iotop nethogs

# Monitor CPU usage
htop

# Monitor disk I/O
iotop

# Monitor network usage
nethogs
```

## 🚀 Memulai Mining

### Step 1: Persiapan Environment

```bash
# 1. Clone dan build project
git clone <repository-url>
cd blockchain-node
go mod tidy
go build -o blockchain-node cmd/main.go

# 2. Buat direktori data
mkdir -p ./data
mkdir -p ./logs

# 3. Set permissions (Linux/macOS)
chmod +x blockchain-node
```

### Step 2: Validasi Konfigurasi

```bash
# Cek konfigurasi mining
./blockchain-node config validate

# Test koneksi network
./blockchain-node networkinfo
```

### Step 3: Start Mining Node

```bash
# Start dengan mining enabled
./blockchain-node startnode --mining-enabled \
  --mining-address "0x1234567890abcdef..." \
  --mining-threads 4 \
  --log-level info

# Atau menggunakan config file
./blockchain-node startnode --config .blockchain-node.yaml
```

### Step 4: Verifikasi Mining Berjalan

Cek log untuk memastikan mining berjalan:

```bash
# Monitor logs real-time
tail -f ./logs/blockchain.log

# Atau jika menggunakan console output
./blockchain-node startnode --log-output console --log-level debug
```

Output yang diharapkan:
```
[INFO] Mining engine started with 4 threads
[INFO] Mining address: 0x1234567890abcdef...
[INFO] Target difficulty: 4
[INFO] Mining... nonce: 100000
[INFO] Mining... nonce: 200000
[INFO] Block mined! Nonce: 245891, Hash: 0xabcd..., Time: 12.3s
[INFO] Block added to blockchain. Height: 1245
```

## 📊 Monitoring Mining

### 1. Mining Metrics

Akses metrics melalui HTTP endpoint:

```bash
# Get mining metrics
curl http://localhost:8080/metrics

# Specific mining data
curl http://localhost:8545/eth/blockNumber
```

Key metrics untuk monitoring:
- **Hash Rate**: Kecepatan hashing (hash/detik)
- **Block Time**: Waktu rata-rata untuk mine block
- **Difficulty**: Current mining difficulty
- **Rewards**: Total mining rewards earned

### 2. Performance Monitoring

```bash
# CPU usage monitoring
top -p $(pgrep blockchain-node)

# Memory usage
ps -p $(pgrep blockchain-node) -o pid,vsz,rss,comm

# Network connections
netstat -an | grep $(pgrep blockchain-node)
```

### 3. Dashboard Commands

```bash
# Get current mining status
./blockchain-node mining status

# Get wallet balance
./blockchain-node getbalance 0x1234567890abcdef...

# Get network info
./blockchain-node networkinfo

# Get latest blocks
./blockchain-node getblock latest
```

## 🔧 Troubleshooting

### Problem: Mining Tidak Dimulai

**Symptoms**: Log menunjukkan "Mining disabled" atau tidak ada aktivitas mining

**Solutions**:
```bash
# 1. Cek konfigurasi
./blockchain-node config show | grep mining

# 2. Pastikan wallet address valid
./blockchain-node validateaddress 0x1234567890abcdef...

# 3. Cek apakah port tidak conflict
netstat -tulpn | grep :8080
```

### Problem: Hash Rate Rendah

**Symptoms**: Hash rate di bawah ekspektasi berdasarkan hardware

**Solutions**:
```yaml
# Optimasi config untuk hash rate
mining:
  threads: 8          # Increase threads (max = CPU cores)
  
db:
  cache_size: 256     # Increase cache
  write_buffer: 16    # Increase buffer
```

### Problem: Memory Usage Tinggi

**Symptoms**: Node menggunakan RAM berlebihan

**Solutions**:
```bash
# 1. Restart node secara berkala
./blockchain-node stop
./blockchain-node startnode

# 2. Reduce cache size
# Edit .blockchain-node.yaml:
db:
  cache_size: 64      # Reduce dari 128
```

### Problem: Sering Disconnect dari Peers

**Symptoms**: Log menunjukkan "peer disconnected" berulang

**Solutions**:
```yaml
# Optimasi network config
network:
  timeout: 60         # Increase timeout
  max_peers: 10       # Reduce max peers
  
# Tambah seed nodes stabil
  seed_nodes:
    - "stable-node-1.blockchain.com:8080"
    - "stable-node-2.blockchain.com:8080"
```

## ⚡ Optimisasi Performa

### 1. CPU Optimization

```bash
# Set CPU affinity (Linux)
taskset -c 0-3 ./blockchain-node startnode

# Set process priority
nice -n -10 ./blockchain-node startnode
```

### 2. Disk I/O Optimization

```yaml
# Database optimizations
db:
  type: "leveldb"
  cache_size: 256        # 256MB cache
  max_open_files: 4000   # More file handles
  write_buffer: 16       # 16MB write buffer
  compression: true      # Enable compression
```

### 3. Network Optimization

```yaml
network:
  buffer_size: 65536     # 64KB buffer
  read_timeout: 30
  write_timeout: 30
  keepalive: true
```

### 4. Memory Management

```bash
# Set Go garbage collection
export GOGC=100        # Default GC percentage
export GOMEMLIMIT=8GB  # Limit memory usage

# Run with memory optimizations
./blockchain-node startnode
```

## 📈 Mining Economics

### Reward Calculation

```
Block Reward = Base Reward + Transaction Fees
Base Reward = 5 ETH (configurable)
Transaction Fees = Sum of (Gas Used * Gas Price) for all transactions
```

### Profitability Analysis

```bash
# Calculate mining profitability
./blockchain-node mining profitability \
  --electricity-cost 0.12 \     # USD per kWh
  --power-consumption 500 \      # Watts
  --hash-rate 1000000           # Hash per second
```

### Mining Pool Considerations

Untuk mining pool, tambahkan konfigurasi:

```yaml
mining:
  pool_enabled: true
  pool_address: "pool.blockchain.com:8080"
  pool_username: "your-username"
  pool_password: "your-password"
```

## 🔒 Security Best Practices

### 1. Wallet Security

```bash
# Encrypt wallet
./blockchain-node encryptwallet --passphrase "strong-password"

# Backup wallet
cp ./data/wallet.dat ./backup/wallet-$(date +%Y%m%d).dat
```

### 2. Network Security

```yaml
# Restrict RPC access
rpc:
  host: "127.0.0.1"    # Only localhost
  cors_origins: []      # No CORS
  
# Use firewall
# iptables -A INPUT -p tcp --dport 8080 -s trusted-ip -j ACCEPT
```

### 3. Monitoring & Alerts

```bash
# Setup log monitoring
tail -f logs/blockchain.log | grep -i "error\|warning"

# Monitor failed mining attempts
grep "mining failed" logs/blockchain.log | wc -l
```

## 📊 Mining Statistics

### Daily Mining Report

```bash
#!/bin/bash
# daily-report.sh

echo "=== Daily Mining Report ==="
echo "Date: $(date)"
echo "Blocks Mined: $(grep 'Block mined' logs/blockchain.log | wc -l)"
echo "Total Rewards: $(./blockchain-node getbalance $MINING_ADDRESS)"
echo "Average Block Time: $(./blockchain-node stats blocktime)"
echo "Current Difficulty: $(./blockchain-node stats difficulty)"
```

### Performance Metrics

```bash
# Get mining performance
./blockchain-node mining performance

# Expected output:
# Hash Rate: 1,000,000 H/s
# Blocks Found: 144 (last 24h)
# Success Rate: 12.5%
# Average Block Time: 600s
# Efficiency: 98.5%
```

## 🎯 Advanced Mining

### Custom Mining Algorithm

Untuk advanced users yang ingin modify mining algorithm:

```go
// consensus/custom_pow.go
func (pow *CustomProofOfWork) Mine(block *core.Block) error {
    // Implement custom mining logic
    // with optimized hashing or additional validation
}
```

### GPU Mining Integration

```bash
# Compile with GPU support (jika tersedia)
CGO_ENABLED=1 go build -tags gpu -o blockchain-node-gpu cmd/main.go

# Run with GPU acceleration
./blockchain-node-gpu startnode --mining-gpu --gpu-devices 0,1
```

---

**Happy Mining!** 🎉

Untuk bantuan lebih lanjut, bergabunglah dengan komunitas di Discord atau buka issue di GitHub repository.
