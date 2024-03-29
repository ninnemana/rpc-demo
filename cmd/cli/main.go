package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ninnemana/rpc-demo/pkg/vinyltap"
)

var cfgFile string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vinyltap",
	Short: "A brief description of your application",
}

func main() {
	cobra.OnInitialize(func(){
		if cfgFile != "" { // enable ability to specify config file via flag
			viper.SetConfigFile(cfgFile)
		}

		viper.SetConfigName(".example") // name of config file (without extension)
		viper.AddConfigPath("$HOME")  // adding home directory as first search path
		viper.AutomaticEnv()          // read in environment variables that match

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	})

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.example.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.AddCommand(vinyltap.TapClientCommand)

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
