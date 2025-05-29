
package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Network NetworkConfig `mapstructure:"network"`
	RPC     RPCConfig     `mapstructure:"rpc"`
	Mining  MiningConfig  `mapstructure:"mining"`
	DB      DBConfig      `mapstructure:"db"`
	EVM     EVMConfig     `mapstructure:"evm"`
}

type NetworkConfig struct {
	Port       int      `mapstructure:"port"`
	SeedNodes  []string `mapstructure:"seed_nodes"`
	MaxPeers   int      `mapstructure:"max_peers"`
	ListenAddr string   `mapstructure:"listen_addr"`
}

type RPCConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Host    string `mapstructure:"host"`
}

type MiningConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Address    string `mapstructure:"address"`
	Threads    int    `mapstructure:"threads"`
	Difficulty uint64 `mapstructure:"difficulty"`
}

type DBConfig struct {
	Path string `mapstructure:"path"`
	Type string `mapstructure:"type"`
}

type EVMConfig struct {
	ChainID      uint64 `mapstructure:"chain_id"`
	BlockGasLimit uint64 `mapstructure:"block_gas_limit"`
}

func LoadConfig() *Config {
	// Set default values
	viper.SetDefault("network.port", 8080)
	viper.SetDefault("network.max_peers", 50)
	viper.SetDefault("network.listen_addr", "0.0.0.0")
	viper.SetDefault("rpc.enabled", true)
	viper.SetDefault("rpc.port", 8545)
	viper.SetDefault("rpc.host", "localhost")
	viper.SetDefault("mining.enabled", false)
	viper.SetDefault("mining.threads", 1)
	viper.SetDefault("mining.difficulty", 4)
	viper.SetDefault("db.path", "./data")
	viper.SetDefault("db.type", "leveldb")
	viper.SetDefault("evm.chain_id", 1337)
	viper.SetDefault("evm.block_gas_limit", 8000000)

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return &config
}
