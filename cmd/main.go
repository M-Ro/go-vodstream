package main

import (
	"fmt"
	"github.com/M-Ro/go-vodstream/cmd/streamingester"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{}

// init registers all the available commands to the cli
func init() {
	rootCmd.AddCommand(streamingester.NewCmd())
}

// initialises viper config library.
func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Loaded conf:", viper.ConfigFileUsed())
}

func main() {
	initConfig()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}