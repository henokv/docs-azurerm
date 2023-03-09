/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/henokv/docs-azurerm/internal"
	"github.com/spf13/cobra"
	"log"
)

// vnetCmd represents the vnet command
var vnetCmd = &cobra.Command{
	Use:     "vnet",
	Short:   "This command will generate the docs for the networking components in azure",
	RunE:    vnetCMDRunE,
	Version: rootCmd.Version,
}

func vnetCMDRunE(cmd *cobra.Command, args []string) error {
	log.Println("Generating docs ...")
	internal.CleanDocsDir()
	subs, err := internal.GetAllSubscriptions()
	if err != nil {
		return err
	}
	vnets, err := internal.GetWrappedVNETsInSubscriptions(subs)
	if err != nil {
		return err
	}
	internal.WriteMarkdown(subs)
	for _, sub := range subs {
		sub.WriteMarkdown()
		if err != nil {
			return err
		}
	}
	for _, vnet := range vnets {
		vnet.WriteMarkdown()
		if err != nil {
			return err
		}
	}
	log.Println("Docs are generated")
	return nil
}

func init() {
	rootCmd.AddCommand(vnetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vnetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vnetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
