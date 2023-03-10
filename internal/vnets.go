package internal

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/c-robinson/iplib"
	"sort"
	"strings"
)

type VNETWrapper struct {
	*armnetwork.VirtualNetwork
	ResourceGroup string
	Subscription  SubscriptionWrapper
	IPSpaces      []*IPSpace
}

func NewVNETWrapper(vnet *armnetwork.VirtualNetwork, subscriptionWrapper SubscriptionWrapper) VNETWrapper {
	split := strings.Split(*vnet.ID, "/")
	rg := split[4]
	wrapper := VNETWrapper{vnet, rg, subscriptionWrapper, []*IPSpace{}}
	sort.Slice(wrapper.Properties.AddressSpace.AddressPrefixes, func(i, j int) bool {
		a := wrapper.Properties.AddressSpace.AddressPrefixes[i]
		b := wrapper.Properties.AddressSpace.AddressPrefixes[j]
		return iplib.CompareNets(iplib.Net4FromStr(*a), iplib.Net4FromStr(*b)) == -1
	})

	sort.Slice(wrapper.Properties.Subnets, func(i, j int) bool {
		a := wrapper.Properties.Subnets[i].Properties.AddressPrefix
		b := wrapper.Properties.Subnets[j].Properties.AddressPrefix
		return iplib.CompareNets(iplib.Net4FromStr(*a), iplib.Net4FromStr(*b)) == -1
	})
	wrapper.generateIPSpaces()
	return wrapper
}

func (vnet *VNETWrapper) generateIPSpaces() {
	vnet.IPSpaces = []*IPSpace{}
	for _, addressSpace := range vnet.Properties.AddressSpace.AddressPrefixes {
		ipSpace := NewIPSPace(vnet, *addressSpace)

		for _, subnet := range vnet.Properties.Subnets {
			if ipSpace.AddSubnet(subnet) {
				//Subnet is part of this IPspace
			}
		}
		ipSpace.generateFreeIPs()
		vnet.IPSpaces = append(vnet.IPSpaces, ipSpace)
	}

}

func (vnet *VNETWrapper) getFreeIPSPace() (freeSpaces []iplib.Net4) {
	for _, space := range vnet.IPSpaces {
		freeSpaces = append(freeSpaces, space.freeSpace...)
	}
	return freeSpaces
}

func (vnet VNETWrapper) GenerateMarkdown() string {
	var markdown string
	markdown += fmt.Sprintf("# %s  \n", *vnet.Name)
	markdown += fmt.Sprintf("#### Location: %s  \n", *vnet.Location)
	markdown += fmt.Sprintf("#### RG: %s  \n", vnet.ResourceGroup)
	markdown += fmt.Sprintf("#### Subscription: [%s](../Readme.md)  \n", *vnet.Subscription.DisplayName)
	markdown += fmt.Sprintf("#### Ranges  \n")
	for _, prefix := range vnet.Properties.AddressSpace.AddressPrefixes {
		markdown += fmt.Sprintf("- %s  \n", *prefix)
	}
	markdown += fmt.Sprintf("### Subnets  \n")
	markdown += fmt.Sprintf("| Prefix | Name | Route Table | NSG |\n")
	markdown += fmt.Sprintf("| --- | --- | --- | --- |\n")
	for _, subnet := range vnet.Properties.Subnets {
		markdown += fmt.Sprintf("| %s | %s | %s | %s |\n",
			*subnet.Properties.AddressPrefix,
			*subnet.Name,
			getRouteTableName(subnet),
			getNsgName(subnet))
	}
	markdown += fmt.Sprintf("### Free Space  \n")
	markdown += fmt.Sprintf("| Prefix | Size (Usable) |\n")
	markdown += fmt.Sprintf("| --- | --- |\n")
	for _, space := range vnet.getFreeIPSPace() {
		usableIPs := iplib.DeltaIP4(space.NetworkAddress(), iplib.IncrementIPBy(space.BroadcastAddress(), 1)) - 5
		markdown += fmt.Sprintf("| %s | %d |\n",
			space.String(), usableIPs,
		)
	}
	markdown += fmt.Sprintf("### Peerings  \n")
	markdown += fmt.Sprintf("| VNET | Spaces |\n")
	markdown += fmt.Sprintf("| --- | --- |\n")
	for _, peering := range vnet.Properties.VirtualNetworkPeerings {
		nameParts := strings.Split(*peering.Properties.RemoteVirtualNetwork.ID, "/")
		subscriptionId := nameParts[2]
		remoteRangesPtr := peering.Properties.RemoteVirtualNetworkAddressSpace.AddressPrefixes
		var remoteRanges []string
		for _, rrp := range remoteRangesPtr {
			remoteRanges = append(remoteRanges, *rrp)
		}
		vnetNameMarkdown := nameParts[8]
		for _, subscription := range subscriptionList {
			if *subscription.SubscriptionID == subscriptionId {
				vnetNameMarkdown = strings.ReplaceAll(fmt.Sprintf("[%s](./../../%s/%s/%s.md)", nameParts[8], *subscription.DisplayName, nameParts[4], nameParts[8]), " ", "%20")
				break
			}
		}
		markdown += fmt.Sprintf("| %s | %v |\n", vnetNameMarkdown, remoteRanges)
	}
	return markdown
}

func (vnet VNETWrapper) WriteMarkdown() error {
	markdown := vnet.GenerateMarkdown()
	err := WriteToFile(markdown, fmt.Sprintf("docs/%s/%s/%s.md", *vnet.Subscription.DisplayName, vnet.ResourceGroup, *vnet.Name))
	return err
}

func GetWrappedVNETsInSubscriptions(subscriptions []*SubscriptionWrapper) (vnets []*VNETWrapper, err error) {
	for _, subscription := range subscriptions {
		vnetsInSub, err := getWrappedVNETsInSubscription(subscription)
		if err != nil {
			return vnets, err
		}
		vnets = append(vnets, vnetsInSub...)
	}
	return vnets, nil
}

func getWrappedVNETsInSubscription(subscription *SubscriptionWrapper) (vnets []*VNETWrapper, err error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return vnets, err
	}
	client, err := armnetwork.NewVirtualNetworksClient(*subscription.SubscriptionID, cred, nil)
	if err != nil {
		return vnets, err
	}
	pager := client.NewListAllPager(nil)
	for pager.More() {
		vnetList, err := pager.NextPage(context.Background())
		if err != nil {
			return vnets, err
		}
		for _, vnet := range vnetList.VirtualNetworkListResult.Value {
			wrappedVNET := NewVNETWrapper(vnet, *subscription)
			vnets = append(vnets, &wrappedVNET)
			subscription.vnets = append(subscription.vnets, &wrappedVNET)
		}
	}
	return vnets, nil
}

func getRouteTableName(subnet *armnetwork.Subnet) string {
	if subnet.Properties.RouteTable == nil {
		return "-"
	} else {
		idSplit := strings.Split(*subnet.Properties.RouteTable.ID, "/")
		return idSplit[len(idSplit)-1]
	}
}

func getNsgName(subnet *armnetwork.Subnet) string {
	if subnet.Properties.NetworkSecurityGroup == nil {
		return "-"
	} else {
		idSplit := strings.Split(*subnet.Properties.NetworkSecurityGroup.ID, "/")
		return idSplit[len(idSplit)-1]
	}
}
