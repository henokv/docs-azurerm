/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/henokv/docs-azurerm/internal"
	"github.com/spf13/cobra"
)

// vnetCmd represents the vnet command
var vnetCmd = &cobra.Command{
	Use:     "vnet",
	Short:   "This command will generate the docs for the networking components in azure",
	Run:     vnetCMDRun,
	Version: rootCmd.Version,
}

func vnetCMDRun(cmd *cobra.Command, args []string) {
	client, err := internal.NewDocumentationClient("docs")
	cobra.CheckErr(err)
	err = client.GenerateMarkdown(false)
	cobra.CheckErr(err)
}

func init() {
	rootCmd.AddCommand(vnetCmd)
}
