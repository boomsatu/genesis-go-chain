
# 📊 Metrics & Monitoring Guide

Panduan lengkap untuk monitoring dan metrics collection pada blockchain node untuk production deployment.

## 📋 Daftar Isi

- [Overview](#overview)
- [Metrics Configuration](#metrics-configuration)
- [Available Metrics](#available-metrics)
- [Prometheus Integration](#prometheus-integration)
- [Grafana Dashboard](#grafana-dashboard)
- [Alerting](#alerting)
- [Performance Monitoring](#performance-monitoring)
- [Log Analysis](#log-analysis)
- [Troubleshooting](#troubleshooting)

## 🔍 Overview

Blockchain node menyediakan sistem monitoring komprehensif yang mencakup:
- **Real-time Metrics**: Performance, network, dan blockchain statistics
- **Prometheus Integration**: Time-series metrics collection
- **Structured Logging**: JSON-formatted logs untuk analysis
- **Health Checks**: Automated system health monitoring
- **Custom Alerts**: Configurable alerting untuk critical events

### Monitoring Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Blockchain    │    │    Metrics     │    │   Prometheus    │
│      Node       │───▶│   Collector    │───▶│     Server      │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                 │
                                 ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │     Grafana     │    │   Alerting      │
                       │    Dashboard    │    │     System      │
                       │                 │    │                 │
                       └─────────────────┘    └─────────────────┘
```

## ⚙️ Metrics Configuration

### Basic Configuration

```yaml
# .blockchain-node.yaml
metrics:
  enabled: true
  port: 8080
  path: "/metrics"
  interval: 30               # Collection interval in seconds
  retention: "24h"           # Metrics retention period
  
  # Prometheus settings
  prometheus:
    enabled: true
    namespace: "blockchain"
    subsystem: "node"
    
  # Custom metrics
  custom_metrics:
    - name: "custom_counter"
      type: "counter"
      help: "Custom counter metric"
```

### Advanced Configuration

```yaml
metrics:
  enabled: true
  port: 8080
  path: "/metrics"
  
  # Security
  auth_required: false
  allowed_ips: ["127.0.0.1", "10.0.0.0/8"]
  
  # Collection settings
  collection:
    system_metrics: true      # CPU, Memory, Disk
    network_metrics: true     # P2P, RPC stats
    blockchain_metrics: true  # Blocks, transactions
    mining_metrics: true      # Mining performance
    
  # Export settings
  exporters:
    prometheus: true
    statsd: false
    influxdb: false
    
  # Performance
  buffer_size: 1000
  batch_size: 100
  flush_interval: 10
```

## 📈 Available Metrics

### Blockchain Metrics

#### Block Metrics
```prometheus
# Total number of blocks processed
blockchain_blocks_total{status="confirmed|orphaned"} counter

# Current blockchain height
blockchain_height gauge

# Block processing time
blockchain_block_process_duration_seconds histogram

# Block size distribution
blockchain_block_size_bytes histogram

# Blocks per second
blockchain_blocks_per_second gauge
```

#### Transaction Metrics
```prometheus
# Total transactions processed
blockchain_transactions_total{type="transfer|contract_creation|contract_call"} counter

# Transaction pool size
blockchain_mempool_size gauge

# Transaction validation time
blockchain_transaction_validation_duration_seconds histogram

# Transaction fees
blockchain_transaction_fees_total{currency="wei"} counter

# Gas usage
blockchain_gas_used_total counter
blockchain_gas_limit_total counter
```

### Network Metrics

#### P2P Network
```prometheus
# Connected peers
blockchain_peers_connected gauge

# Network traffic
blockchain_network_bytes_total{direction="in|out"} counter

# Peer connection duration
blockchain_peer_connection_duration_seconds histogram

# Network latency
blockchain_network_latency_seconds histogram

# Message counts
blockchain_messages_total{type="block|transaction|peer"} counter
```

#### RPC Metrics
```prometheus
# RPC requests
blockchain_rpc_requests_total{method="eth_getBalance|eth_blockNumber"} counter

# RPC response time
blockchain_rpc_duration_seconds{method="eth_getBalance"} histogram

# RPC errors
blockchain_rpc_errors_total{method="eth_getBalance",error="timeout|invalid_params"} counter

# Active RPC connections
blockchain_rpc_connections_active gauge
```

### Mining Metrics

#### Mining Performance
```prometheus
# Blocks mined
blockchain_blocks_mined_total counter

# Mining hash rate
blockchain_mining_hashrate gauge

# Mining difficulty
blockchain_mining_difficulty gauge

# Mining efficiency
blockchain_mining_efficiency_ratio gauge

# Mining time per block
blockchain_mining_block_time_seconds histogram
```

### System Metrics

#### Resource Usage
```prometheus
# CPU usage
blockchain_cpu_usage_percent gauge

# Memory usage
blockchain_memory_usage_bytes gauge
blockchain_memory_available_bytes gauge

# Disk usage
blockchain_disk_usage_bytes{device="/dev/sda1"} gauge
blockchain_disk_io_bytes_total{device="/dev/sda1",direction="read|write"} counter

# File descriptor usage
blockchain_file_descriptors_used gauge
blockchain_file_descriptors_limit gauge
```

#### Go Runtime Metrics
```prometheus
# Garbage collection
go_gc_duration_seconds histogram
go_memstats_gc_cpu_fraction gauge

# Goroutines
go_goroutines gauge

# Memory allocations
go_memstats_alloc_bytes gauge
go_memstats_heap_objects gauge
```

## 🎯 Prometheus Integration

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "blockchain_rules.yml"

scrape_configs:
  - job_name: 'blockchain-node'
    static_configs:
      - targets: ['localhost:8080']
    scrape_interval: 10s
    metrics_path: /metrics
    
alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - localhost:9093
```

### Installation & Setup

```bash
# Download Prometheus
wget https://github.com/prometheus/prometheus/releases/download/v2.40.0/prometheus-2.40.0.linux-amd64.tar.gz
tar xvfz prometheus-*.tar.gz
cd prometheus-*

# Configure Prometheus
cp prometheus.yml.example prometheus.yml
# Edit prometheus.yml dengan config di atas

# Start Prometheus
./prometheus --config.file=prometheus.yml --storage.tsdb.path=./data
```

### Custom Recording Rules

```yaml
# blockchain_rules.yml
groups:
  - name: blockchain.rules
    rules:
      # Block rate (blocks per minute)
      - record: blockchain:block_rate_5m
        expr: rate(blockchain_blocks_total[5m]) * 60
        
      # Transaction throughput
      - record: blockchain:tx_throughput_5m
        expr: rate(blockchain_transactions_total[5m])
        
      # Average block time
      - record: blockchain:avg_block_time_5m
        expr: rate(blockchain_block_process_duration_seconds_sum[5m]) / rate(blockchain_block_process_duration_seconds_count[5m])
        
      # Network efficiency
      - record: blockchain:network_efficiency_5m
        expr: blockchain_peers_connected / 100 * rate(blockchain_blocks_total[5m])
```

## 📊 Grafana Dashboard

### Installation

```bash
# Install Grafana
sudo apt-get install -y software-properties-common
sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"
sudo apt-get update
sudo apt-get install grafana

# Start Grafana
sudo systemctl start grafana-server
sudo systemctl enable grafana-server
```

### Dashboard Configuration

```json
{
  "dashboard": {
    "title": "Blockchain Node Monitoring",
    "panels": [
      {
        "title": "Block Height",
        "type": "stat",
        "targets": [
          {
            "expr": "blockchain_height",
            "legendFormat": "Current Height"
          }
        ]
      },
      {
        "title": "Transaction Pool",
        "type": "graph",
        "targets": [
          {
            "expr": "blockchain_mempool_size",
            "legendFormat": "Mempool Size"
          }
        ]
      },
      {
        "title": "Network Peers",
        "type": "stat",
        "targets": [
          {
            "expr": "blockchain_peers_connected",
            "legendFormat": "Connected Peers"
          }
        ]
      },
      {
        "title": "Mining Performance",
        "type": "graph",
        "targets": [
          {
            "expr": "blockchain_mining_hashrate",
            "legendFormat": "Hash Rate"
          },
          {
            "expr": "blockchain_mining_difficulty",
            "legendFormat": "Difficulty"
          }
        ]
      }
    ]
  }
}
```

### Key Dashboard Panels

#### 1. Blockchain Overview
- Current block height
- Total transactions
- Network peers
- Mining status
- Sync progress

#### 2. Performance Metrics
- Transaction throughput
- Block processing time
- RPC response times
- System resource usage

#### 3. Network Health
- Peer connections
- Network latency
- Message propagation
- Connection stability

#### 4. Mining Dashboard
- Hash rate trends
- Difficulty adjustments
- Blocks mined
- Mining efficiency
- Power consumption estimates

## 🚨 Alerting

### Alertmanager Configuration

```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'localhost:587'
  smtp_from: 'alerts@blockchain.com'

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
  - name: 'web.hook'
    email_configs:
      - to: 'admin@blockchain.com'
        subject: 'Blockchain Alert: {{ .GroupLabels.alertname }}'
        body: |
          {{ range .Alerts }}
          Alert: {{ .Annotations.summary }}
          Description: {{ .Annotations.description }}
          {{ end }}

  - name: 'slack'
    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK_URL'
        channel: '#blockchain-alerts'
        title: 'Blockchain Alert'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

### Alert Rules

```yaml
# alert_rules.yml
groups:
  - name: blockchain.alerts
    rules:
      # Node down
      - alert: BlockchainNodeDown
        expr: up{job="blockchain-node"} == 0
        for: 30s
        labels:
          severity: critical
        annotations:
          summary: "Blockchain node is down"
          description: "Node {{ $labels.instance }} has been down for more than 30 seconds"

      # High memory usage
      - alert: HighMemoryUsage
        expr: blockchain_memory_usage_bytes / blockchain_memory_available_bytes > 0.9
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage detected"
          description: "Memory usage is above 90% for more than 5 minutes"

      # Low peer count
      - alert: LowPeerCount
        expr: blockchain_peers_connected < 3
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Low peer count"
          description: "Connected peers ({{ $value }}) is below minimum threshold"

      # Mempool congestion
      - alert: MempoolCongestion
        expr: blockchain_mempool_size > 10000
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Mempool congestion detected"
          description: "Mempool size ({{ $value }}) is above normal levels"

      # Mining stopped
      - alert: MiningStopped
        expr: increase(blockchain_blocks_mined_total[10m]) == 0
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "Mining has stopped"
          description: "No blocks mined in the last 10 minutes"

      # Sync lag
      - alert: SyncLag
        expr: blockchain_height < blockchain_network_height - 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Node falling behind network"
          description: "Node is {{ $value }} blocks behind network"

      # High RPC error rate
      - alert: HighRPCErrorRate
        expr: rate(blockchain_rpc_errors_total[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "High RPC error rate"
          description: "RPC error rate is {{ $value }} errors/second"
```

## 📈 Performance Monitoring

### Key Performance Indicators (KPIs)

#### Blockchain KPIs
```bash
# Block processing time
curl http://localhost:8080/metrics | grep blockchain_block_process_duration

# Transaction throughput
curl http://localhost:8080/metrics | grep blockchain_transactions_total

# Sync status
curl http://localhost:8080/metrics | grep blockchain_height
```

#### System KPIs
```bash
# CPU usage
curl http://localhost:8080/metrics | grep blockchain_cpu_usage

# Memory usage
curl http://localhost:8080/metrics | grep blockchain_memory_usage

# Disk I/O
curl http://localhost:8080/metrics | grep blockchain_disk_io
```

### Performance Benchmarking

```bash
#!/bin/bash
# performance_benchmark.sh

echo "=== Blockchain Node Performance Benchmark ==="

# Collect baseline metrics
echo "Collecting baseline metrics..."
START_TIME=$(date +%s)
START_BLOCKS=$(curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result')
START_TXS=$(curl -s http://localhost:8080/metrics | grep "blockchain_transactions_total" | awk '{print $2}')

# Wait for measurement period
echo "Measuring performance for 5 minutes..."
sleep 300

# Collect end metrics
END_TIME=$(date +%s)
END_BLOCKS=$(curl -s http://localhost:8545 -X POST -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' | jq -r '.result')
END_TXS=$(curl -s http://localhost:8080/metrics | grep "blockchain_transactions_total" | awk '{print $2}')

# Calculate performance
DURATION=$((END_TIME - START_TIME))
BLOCKS_PROCESSED=$((END_BLOCKS - START_BLOCKS))
TXS_PROCESSED=$((END_TXS - START_TXS))

echo "=== Performance Results ==="
echo "Duration: ${DURATION} seconds"
echo "Blocks processed: ${BLOCKS_PROCESSED}"
echo "Transactions processed: ${TXS_PROCESSED}"
echo "Block rate: $(echo "scale=2; $BLOCKS_PROCESSED * 60 / $DURATION" | bc) blocks/minute"
echo "Transaction rate: $(echo "scale=2; $TXS_PROCESSED / $DURATION" | bc) tx/second"
```

### Resource Monitoring Scripts

```bash
#!/bin/bash
# resource_monitor.sh

LOG_FILE="resource_usage.log"
INTERVAL=30  # seconds

while true; do
    TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
    
    # Get process ID
    PID=$(pgrep blockchain-node)
    
    if [ -n "$PID" ]; then
        # CPU and Memory usage
        CPU=$(ps -p $PID -o %cpu --no-headers)
        MEM=$(ps -p $PID -o %mem --no-headers)
        VSZ=$(ps -p $PID -o vsz --no-headers)
        RSS=$(ps -p $PID -o rss --no-headers)
        
        # Network connections
        CONNECTIONS=$(netstat -an | grep :8080 | wc -l)
        
        # File descriptors
        FDS=$(ls /proc/$PID/fd | wc -l)
        
        echo "$TIMESTAMP,$CPU,$MEM,$VSZ,$RSS,$CONNECTIONS,$FDS" >> $LOG_FILE
    else
        echo "$TIMESTAMP,NODE_DOWN" >> $LOG_FILE
    fi
    
    sleep $INTERVAL
done
```

## 📋 Log Analysis

### Log Configuration

```yaml
# .blockchain-node.yaml
logging:
  level: "info"
  output: "both"              # console, file, both
  file_path: "./logs/blockchain.log"
  max_size: 100              # MB
  max_backups: 10
  max_age: 30                # days
  compress: true
  
  # Structured logging
  format: "json"             # json, text
  component: "blockchain-node"
  
  # Log sampling
  sampling:
    enabled: true
    threshold: 100           # logs per second
    thereafter: 10           # 1 in 10 after threshold
```

### Log Analysis Tools

#### ELK Stack Integration

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  paths:
    - /path/to/blockchain/logs/*.log
  fields:
    service: blockchain-node
  fields_under_root: true

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "blockchain-logs-%{+yyyy.MM.dd}"

setup.kibana:
  host: "localhost:5601"
```

#### Log Parsing Scripts

```bash
#!/bin/bash
# log_analyzer.sh

LOG_FILE="./logs/blockchain.log"

echo "=== Blockchain Log Analysis ==="

# Error analysis
echo "Errors in last hour:"
grep -i "error" $LOG_FILE | grep "$(date -d '1 hour ago' '+%Y-%m-%d %H')" | wc -l

# Block mining success rate
echo "Mining analysis:"
BLOCKS_ATTEMPTED=$(grep "Mining..." $LOG_FILE | wc -l)
BLOCKS_MINED=$(grep "Block mined!" $LOG_FILE | wc -l)
if [ $BLOCKS_ATTEMPTED -gt 0 ]; then
    SUCCESS_RATE=$(echo "scale=2; $BLOCKS_MINED * 100 / $BLOCKS_ATTEMPTED" | bc)
    echo "Success rate: ${SUCCESS_RATE}%"
fi

# Network events
echo "Network events:"
grep "peer connected" $LOG_FILE | tail -5
grep "peer disconnected" $LOG_FILE | tail -5

# Performance metrics from logs
echo "Average block time:"
grep "Block mined!" $LOG_FILE | tail -100 | awk '{print $NF}' | sed 's/s$//' | \
awk '{sum+=$1; count++} END {if(count>0) print sum/count "s"}'
```

### Real-time Log Monitoring

```bash
# Monitor critical events
tail -f logs/blockchain.log | grep -E "(ERROR|CRITICAL|Block mined|peer connected)"

# Monitor performance issues
tail -f logs/blockchain.log | grep -E "(timeout|slow|congestion)"

# Monitor mining activity
tail -f logs/blockchain.log | grep -E "(Mining|difficulty|nonce)"
```

## 🔧 Troubleshooting

### Common Metrics Issues

#### 1. Metrics Not Available

**Problem**: `/metrics` endpoint returns 404 or connection refused

**Solutions**:
```bash
# Check if metrics are enabled
grep "metrics:" .blockchain-node.yaml

# Check if port is open
netstat -tlpn | grep :8080

# Check logs for errors
grep -i "metrics" logs/blockchain.log
```

#### 2. Missing Metrics

**Problem**: Some metrics are not being collected

**Solutions**:
```bash
# Verify configuration
./blockchain-node config validate

# Check metric collection components
curl http://localhost:8080/metrics | grep "blockchain_"

# Enable debug logging
./blockchain-node startnode --log-level debug
```

#### 3. High Memory Usage

**Problem**: Metrics collection consuming too much memory

**Solutions**:
```yaml
# Reduce collection frequency
metrics:
  interval: 60              # Increase from 30 to 60 seconds
  
# Reduce retention
  retention: "12h"          # Reduce from 24h to 12h
  
# Limit metric cardinality
  max_series: 10000
```

### Performance Troubleshooting

```bash
# Check for metrics bottlenecks
curl http://localhost:8080/metrics | grep "collection_duration"

# Monitor metrics memory usage
ps -p $(pgrep blockchain-node) -o pid,vsz,rss,comm

# Check for stuck metrics collectors
grep "metrics collection timeout" logs/blockchain.log
```

### Alerting Troubleshooting

```bash
# Test alert rules
promtool test rules alert_rules.yml

# Check alertmanager status
curl http://localhost:9093/api/v1/status

# Verify alert delivery
curl http://localhost:9093/api/v1/alerts
```

## 📚 Best Practices

### Metrics Collection

1. **Sample Rate**: Use appropriate sampling for high-frequency events
2. **Label Cardinality**: Avoid high-cardinality labels
3. **Retention**: Balance storage costs with observability needs
4. **Aggregation**: Use recording rules for complex calculations

### Monitoring Strategy

1. **RED Method**: Rate, Errors, Duration for services
2. **USE Method**: Utilization, Saturation, Errors for resources
3. **Four Golden Signals**: Latency, traffic, errors, saturation
4. **SLI/SLO**: Define Service Level Indicators and Objectives

### Alerting Guidelines

1. **Alert Fatigue**: Avoid too many noisy alerts
2. **Actionable Alerts**: Every alert should require action
3. **Severity Levels**: Use appropriate severity classification
4. **Escalation**: Implement proper escalation procedures

### Security Considerations

1. **Metrics Access**: Restrict access to metrics endpoints
2. **Sensitive Data**: Avoid exposing sensitive information in metrics
3. **Authentication**: Use authentication for production deployments
4. **Network Security**: Use TLS for metrics transmission

---

**Happy Monitoring!** 📊

Untuk bantuan lebih lanjut atau konfigurasi custom, silakan hubungi tim support atau buka issue di GitHub repository.
