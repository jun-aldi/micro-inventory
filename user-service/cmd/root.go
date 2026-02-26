package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "user-service",
	Short: "User service CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Run(startCmd, nil)
	},
}

// Ngecek Error
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// INisialisasi config
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .env)")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// inisisalisasi config
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(".env")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading config file:", err)
	} else {
		fmt.Fprintln(os.Stdout, "Using config file:", viper.ConfigFileUsed())
	}
}
