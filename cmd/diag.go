/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/azure"
	"github.com/henokv/docs-azurerm/internal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// diagCmd represents the diag command
var diagCmd = &cobra.Command{
	Use:   "diag",
	Short: "Generates a graphviz diagram",
	Run:   diagCmdRun,
}

func diagCmdRun(cmd *cobra.Command, args []string) {
	name := "Networks"
	subs, err := internal.GetAllSubscriptions()
	if err != nil {
		cobra.CheckErr(err)
	}
	vnets, err := internal.GetWrappedVNETsInSubscriptions(subs)
	if err != nil {
		cobra.CheckErr(err)
	}

	if _, err := os.Stat("go-diagrams"); !os.IsNotExist(err) {
		os.RemoveAll("go-diagrams")
	}

	d, err := diagram.New(diagram.Filename(strings.ToLower(name)), diagram.Label(name),
		diagram.Direction("TB"), diagram.WithAttribute("splines", "spline"),
	)
	diagram.DefaultOptions()
	if err != nil {
		log.Fatal(err)
	}
	groups := make(map[string]*diagram.Group)
	nodes := make(map[string]*diagram.Node)
	extraGroups := make(map[string]*diagram.Group)
	extraNodes := make(map[string]*diagram.Node)
	peerings := make(map[string]map[string]bool)
	for _, sub := range subs {
		g := diagram.NewGroup(*sub.DisplayName).Label(*sub.DisplayName)
		groups[*sub.SubscriptionID] = g
		d.Group(g)
	}
	for _, vnet := range vnets {
		node := azure.Network.VirtualNetworks(diagram.NodeLabel(*vnet.Name))
		nodes[*vnet.Name] = node
		groups[*vnet.Subscription.SubscriptionID].Add(node)
	}
	for _, vnet := range vnets {
		for _, peering := range vnet.Properties.VirtualNetworkPeerings {
			remoteVNETName := strings.Split(*peering.Properties.RemoteVirtualNetwork.ID, "/")[8]
			remoteVNETSubscriptionID := strings.Split(*peering.Properties.RemoteVirtualNetwork.ID, "/")[2]
			vnet1 := nodes[*vnet.Name]
			vnet2 := nodes[remoteVNETName]
			if vnet2 == nil { //Remote VNET not known to user
				vnet2 = extraNodes[remoteVNETName]
				if vnet2 == nil { // Remote VNET not in unknown vnet list
					vnet2 = azure.Network.VirtualNetworks(diagram.NodeLabel(remoteVNETName))
					extraNodes[remoteVNETName] = vnet2
					group := extraGroups[remoteVNETSubscriptionID]
					if group == nil {
						group = diagram.NewGroup(remoteVNETSubscriptionID).Label(remoteVNETSubscriptionID)
						extraGroups[remoteVNETSubscriptionID] = group
						d.Group(group)
					}
					group.Add(vnet2)

				}
			}
			nodeIDList := []string{vnet1.ID(), vnet2.ID()}
			fmt.Sprintf("%v", nodeIDList)
			sort.Strings(nodeIDList)
			_, ok := peerings[nodeIDList[0]][nodeIDList[1]]
			if !ok {
				_, ok := peerings[nodeIDList[0]][nodeIDList[1]]
				if !ok {
					peerings[nodeIDList[0]] = make(map[string]bool)
				}
				peerings[nodeIDList[0]][nodeIDList[1]] = true
				d.Connect(vnet1, vnet2, diagram.Bidirectional())
			}

		}
	}

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
	d.Render()
	a := exec.Command("dot", "-Tpng", fmt.Sprintf("%s.dot", strings.ToLower(name)))
	a.Stdout = os.Stdout
	a.Stderr = os.Stderr
	os.Stdin = os.Stdin
	a.Dir = "./go-diagrams"
	a.Run()
}

func init() {
	rootCmd.AddCommand(diagCmd)
}
