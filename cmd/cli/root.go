
package cli

import (
	"fmt"
	"os"

	"blockchain-node/config"
	"blockchain-node/node"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "blockchain-node",
	Short: "Blockchain Node dengan dukungan EVM",
	Long:  `Node blockchain yang mendukung Ethereum Virtual Machine dengan consensus Proof-of-Work`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.blockchain-node.yaml)")
	
	// Add subcommands
	rootCmd.AddCommand(startNodeCmd)
	rootCmd.AddCommand(createWalletCmd)
	rootCmd.AddCommand(getBalanceCmd)
	rootCmd.AddCommand(sendCmd)
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
		viper.SetConfigType("yaml")
		viper.SetConfigName(".blockchain-node")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	cfg = config.LoadConfig()
}

var startNodeCmd = &cobra.Command{
	Use:   "startnode",
	Short: "Start blockchain node",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting blockchain node...")
		
		nodeInstance, err := node.NewNode(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create node: %v\n", err)
			os.Exit(1)
		}

		if err := nodeInstance.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start node: %v\n", err)
			os.Exit(1)
		}
	},
}

var createWalletCmd = &cobra.Command{
	Use:   "createwallet",
	Short: "Create a new wallet",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement wallet creation
		fmt.Println("Creating new wallet...")
	},
}

var getBalanceCmd = &cobra.Command{
	Use:   "getbalance [address]",
	Short: "Get balance of an address",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		address := args[0]
		fmt.Printf("Getting balance for address: %s\n", address)
		// TODO: Implement balance query
	},
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send transaction",
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
	},
}

func init() {
	sendCmd.Flags().StringP("from", "f", "", "Sender address")
	sendCmd.Flags().StringP("to", "t", "", "Recipient address")
	sendCmd.Flags().StringP("amount", "a", "0", "Amount to send")
	sendCmd.Flags().StringP("data", "d", "", "Transaction data (hex)")
	sendCmd.Flags().Uint64P("gaslimit", "l", 21000, "Gas limit")
	sendCmd.Flags().Uint64P("gasprice", "p", 1000000000, "Gas price (wei)")
	
	sendCmd.MarkFlagRequired("from")
	sendCmd.MarkFlagRequired("to")
}
