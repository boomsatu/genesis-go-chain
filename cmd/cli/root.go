
package cli

import (
	"fmt"
	"os"

	"blockchain-node/config"
	"blockchain-node/logger"
	"blockchain-node/node"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	cfg        *config.Config
	debugLevel string
	logOutput  string
	logFile    string
)

var rootCmd = &cobra.Command{
	Use:   "blockchain-node",
	Short: "Professional Blockchain Node with EVM Support",
	Long: `A production-ready blockchain node that supports Ethereum Virtual Machine 
with Proof-of-Work consensus, comprehensive logging, metrics, and monitoring capabilities.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blockchain-node.yaml)")
	rootCmd.PersistentFlags().StringVar(&debugLevel, "log-level", "", "log level (debug, info, warning, error)")
	rootCmd.PersistentFlags().StringVar(&logOutput, "log-output", "", "log output (console, file, both)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "log file path")
	
	// Add subcommands
	rootCmd.AddCommand(startNodeCmd)
	rootCmd.AddCommand(createWalletCmd)
	rootCmd.AddCommand(getBalanceCmd)
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(metricsCmd)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".blockchain-node")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	cfg = config.LoadConfig()

	// Override config with command line flags
	if debugLevel != "" {
		cfg.Logging.Level = debugLevel
	}
	if logOutput != "" {
		cfg.Logging.Output = logOutput
	}
	if logFile != "" {
		cfg.Logging.FilePath = logFile
	}
}

var startNodeCmd = &cobra.Command{
	Use:   "startnode",
	Short: "Start the blockchain node",
	Long:  `Start the blockchain node with all configured services including P2P, RPC, mining, and metrics.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting professional blockchain node...")
		
		// Initialize early logger for startup
		loggerConfig := logger.Config{
			Level:     cfg.Logging.Level,
			Output:    "console",
			Component: "startup",
		}
		if err := logger.Init(loggerConfig); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}

		nodeInstance, err := node.NewNode(cfg)
		if err != nil {
			logger.Fatal("Failed to create node: %v", err)
		}

		if err := nodeInstance.Start(); err != nil {
			logger.Fatal("Failed to start node: %v", err)
		}
	},
}

var createWalletCmd = &cobra.Command{
	Use:   "createwallet",
	Short: "Create a new wallet",
	Long:  `Generate a new cryptographic wallet with private/public key pair and address.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new wallet...")
		// TODO: Implement wallet creation
		fmt.Println("Wallet creation feature coming soon!")
	},
}

var getBalanceCmd = &cobra.Command{
	Use:   "getbalance [address]",
	Short: "Get balance of an address",
	Long:  `Query the blockchain for the balance of a specific address.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		fmt.Printf("Getting balance for address: %s\n", address)
		// TODO: Implement balance query
		fmt.Println("Balance query feature coming soon!")
	},
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a transaction",
	Long:  `Create and broadcast a transaction to the network.`,
	Run: func(cmd *cobra.Command, args []string) {
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		amount, _ := cmd.Flags().GetString("amount")
		data, _ := cmd.Flags().GetString("data")
		gasLimit, _ := cmd.Flags().GetUint64("gaslimit")
		gasPrice, _ := cmd.Flags().GetUint64("gasprice")

		fmt.Printf("Sending transaction from %s to %s, amount: %s\n", from, to, amount)
		if data != "" {
			fmt.Printf("Data: %s\n", data)
		}
		fmt.Printf("Gas Limit: %d, Gas Price: %d\n", gasLimit, gasPrice)
		// TODO: Implement transaction sending
		fmt.Println("Transaction sending feature coming soon!")
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show node status",
	Long:  `Display current status information about the running node.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Node Status:")
		fmt.Println("============")
		// TODO: Implement status display
		fmt.Println("Status display feature coming soon!")
		fmt.Printf("Config file: %s\n", viper.ConfigFileUsed())
		fmt.Printf("Log level: %s\n", cfg.Logging.Level)
		fmt.Printf("Log output: %s\n", cfg.Logging.Output)
	},
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show node metrics",
	Long:  `Display current metrics and statistics about the running node.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Node Metrics:")
		fmt.Println("=============")
		// TODO: Implement metrics display
		fmt.Println("Metrics display feature coming soon!")
	},
}

func init() {
	// Send command flags
	sendCmd.Flags().StringP("from", "f", "", "Sender address")
	sendCmd.Flags().StringP("to", "t", "", "Recipient address")
	sendCmd.Flags().StringP("amount", "a", "0", "Amount to send")
	sendCmd.Flags().StringP("data", "d", "", "Transaction data (hex)")
	sendCmd.Flags().Uint64P("gaslimit", "l", 21000, "Gas limit")
	sendCmd.Flags().Uint64P("gasprice", "p", 1000000000, "Gas price (wei)")
	
	sendCmd.MarkFlagRequired("from")
	sendCmd.MarkFlagRequired("to")

	// Start node command flags
	startNodeCmd.Flags().Bool("mining", false, "Enable mining")
	startNodeCmd.Flags().Bool("rpc", true, "Enable RPC server")
	startNodeCmd.Flags().Bool("metrics", false, "Enable metrics server")
}
