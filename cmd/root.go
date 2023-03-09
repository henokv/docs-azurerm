/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var cfgFile string
var rootCmd = &cobra.Command{
	Use:     "docs-azurerm",
	Short:   "A tool to generate documentation for azure",
	Long:    `This tool will generate docs for azure based on the current resources deployed`,
	Version: "0.2.0",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//RunE: RootCmdRunE,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
	//err := rootCmd.Execute()
	//if err != nil {
	//	os.Exit(1)
	//}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "", "config file (default is $HOME/.cobra.yaml)")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfdocs-azurerm-networking.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//// Find home directory.
	//home, err := os.UserHomeDir()
	//cobra.CheckErr(err)

	// Search config in home directory with name ".cobra" (without extension).
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	//}

	//viper.AutomaticEnv()
	//
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
