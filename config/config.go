package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Network NetworkConfig `mapstructure:"network"`
	RPC     RPCConfig     `mapstructure:"rpc"`
	Mining  MiningConfig  `mapstructure:"mining"`
	DB      DBConfig      `mapstructure:"db"`
	EVM     EVMConfig     `mapstructure:"evm"`
	Logging LoggingConfig `mapstructure:"logging"`
	Metrics MetricsConfig `mapstructure:"metrics"`
}

type NetworkConfig struct {
	Port       int      `mapstructure:"port"`
	SeedNodes  []string `mapstructure:"seed_nodes"`
	MaxPeers   int      `mapstructure:"max_peers"`
	ListenAddr string   `mapstructure:"listen_addr"`
	Timeout    int      `mapstructure:"timeout"`
}

type RPCConfig struct {
	Enabled        bool     `mapstructure:"enabled"`
	Port           int      `mapstructure:"port"`
	Host           string   `mapstructure:"host"`
	CORSOrigins    []string `mapstructure:"cors_origins"`
	MaxConnections int      `mapstructure:"max_connections"`
	Timeout        int      `mapstructure:"timeout"`
}

type MiningConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Address    string `mapstructure:"address"`
	Threads    int    `mapstructure:"threads"`
	Difficulty uint64 `mapstructure:"difficulty"`
}

type DBConfig struct {
	Path          string `mapstructure:"path"`
	Type          string `mapstructure:"type"`
	CacheSize     int    `mapstructure:"cache_size"`
	MaxOpenFiles  int    `mapstructure:"max_open_files"`
	WriteBuffer   int    `mapstructure:"write_buffer"`
}

type EVMConfig struct {
	ChainID       uint64 `mapstructure:"chain_id"`
	BlockGasLimit uint64 `mapstructure:"block_gas_limit"`
	MinGasPrice   uint64 `mapstructure:"min_gas_price"`
}

type LoggingConfig struct {
	Level     string `mapstructure:"level"`
	Output    string `mapstructure:"output"`
	FilePath  string `mapstructure:"file_path"`
	MaxSize   int64  `mapstructure:"max_size"`
	Component string `mapstructure:"component"`
}

type MetricsConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Port       int    `mapstructure:"port"`
	Path       string `mapstructure:"path"`
}

func LoadConfig() *Config {
	// Set default values
	viper.SetDefault("network.port", 8080)
	viper.SetDefault("network.max_peers", 50)
	viper.SetDefault("network.listen_addr", "0.0.0.0")
	viper.SetDefault("network.timeout", 30)
	
	viper.SetDefault("rpc.enabled", true)
	viper.SetDefault("rpc.port", 8545)
	viper.SetDefault("rpc.host", "localhost")
	viper.SetDefault("rpc.cors_origins", []string{"*"})
	viper.SetDefault("rpc.max_connections", 100)
	viper.SetDefault("rpc.timeout", 30)
	
	viper.SetDefault("mining.enabled", false)
	viper.SetDefault("mining.threads", 1)
	viper.SetDefault("mining.difficulty", 4)
	
	viper.SetDefault("db.path", "./data")
	viper.SetDefault("db.type", "leveldb")
	viper.SetDefault("db.cache_size", 64)
	viper.SetDefault("db.max_open_files", 1000)
	viper.SetDefault("db.write_buffer", 4)
	
	viper.SetDefault("evm.chain_id", 1337)
	viper.SetDefault("evm.block_gas_limit", 8000000)
	viper.SetDefault("evm.min_gas_price", 1000000000)
	
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.output", "console")
	viper.SetDefault("logging.file_path", "./logs/blockchain.log")
	viper.SetDefault("logging.max_size", 100)
	viper.SetDefault("logging.component", "blockchain-node")
	
	viper.SetDefault("metrics.enabled", false)
	viper.SetDefault("metrics.port", 8080)
	viper.SetDefault("metrics.path", "/metrics")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Network.Port <= 0 || c.Network.Port > 65535 {
		return fmt.Errorf("invalid network port: %d", c.Network.Port)
	}
	
	if c.RPC.Enabled && (c.RPC.Port <= 0 || c.RPC.Port > 65535) {
		return fmt.Errorf("invalid RPC port: %d", c.RPC.Port)
	}
	
	if c.Mining.Threads <= 0 {
		return fmt.Errorf("mining threads must be positive: %d", c.Mining.Threads)
	}
	
	if c.EVM.ChainID == 0 {
		return fmt.Errorf("chain ID cannot be zero")
	}
	
	return nil
}
