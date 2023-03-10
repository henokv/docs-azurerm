/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/spf13/cobra"
)

// routeCheckCmd represents the routeCheck command
var routeCheckCmd = &cobra.Command{
	Use:   "routeCheck",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: routeCheckRun,
}

func routeCheckRun(cmd *cobra.Command, args []string) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		cobra.CheckErr(err)
	}
	client, err := armnetwork.NewRoutesClient("32595ab2-1344-49c6-af39-cd6fe41334b3", cred, nil)
	if err != nil {
		cobra.CheckErr(err)
	}
	pager := client.NewListPager("eruza-lab-we-rg-01", "a-rt", nil)
	for pager.More() {
		set, err := pager.NextPage(context.Background())
		if err != nil {
			cobra.CheckErr(err)
		}
		fmt.Sprintf("%v", set)
	}
}

func init() {
	rootCmd.AddCommand(routeCheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// routeCheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// routeCheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
