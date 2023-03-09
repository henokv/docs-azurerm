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
}
