package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

// var mockClient *DocumentationClient
var subscriptions []*SubscriptionWrapper
var vnets []*VNETWrapper

func dummySubscriptions() []*SubscriptionWrapper {
	sub1 := NewSubscriptionWrapper(armsubscription.Subscription{
		SubscriptionID: to.Ptr("00000000-0000-0000-0000-000000000001"),
		DisplayName:    to.Ptr("Test Subscription 1"),
	})
	sub2 := NewSubscriptionWrapper(armsubscription.Subscription{
		SubscriptionID: to.Ptr("00000000-0000-0000-0000-000000000002"),
		DisplayName:    to.Ptr("Test Subscription 1"),
	})
	sub3 := NewSubscriptionWrapper(armsubscription.Subscription{
		SubscriptionID: to.Ptr("00000000-0000-0000-0000-000000000003"),
		DisplayName:    to.Ptr("Test Subscription 1"),
	})
	return []*SubscriptionWrapper{&sub1, &sub2, &sub3}
}

func dummyVnets() []*VNETWrapper {
	vnet1 := NewVNETWrapper(&armnetwork.VirtualNetwork{
		Name:     to.Ptr("vnet1"),
		ID:       to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"),
		Location: to.Ptr("westeurope"),
		Properties: &armnetwork.VirtualNetworkPropertiesFormat{
			AddressSpace: &armnetwork.AddressSpace{
				AddressPrefixes: []*string{
					to.Ptr("10.0.0.0/16"),
				},
			},
			Subnets: []*armnetwork.Subnet{
				&armnetwork.Subnet{
					Name: to.Ptr("subnet1"),
					Properties: &armnetwork.SubnetPropertiesFormat{
						AddressPrefix: to.Ptr("10.0.0.0/24"),
						RouteTable: &armnetwork.RouteTable{
							ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/rg1/providers/Microsoft.Network/routeTables/rt1"),
						},
						NetworkSecurityGroup: &armnetwork.SecurityGroup{
							ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000001/resourceGroups/rg1/providers/Microsoft.Network/networkSecurityGroups/nsg1"),
						},
					},
				},
			},
			VirtualNetworkPeerings: []*armnetwork.VirtualNetworkPeering{
				&armnetwork.VirtualNetworkPeering{
					Name: to.Ptr("peering1"),
					ID:   to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000002/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet2/virtualNetworkPeerings/peering1"),
					Properties: &armnetwork.VirtualNetworkPeeringPropertiesFormat{
						RemoteAddressSpace: &armnetwork.AddressSpace{
							AddressPrefixes: []*string{to.Ptr("10.1.0.0/16")},
						},
						RemoteVirtualNetworkAddressSpace: &armnetwork.AddressSpace{
							AddressPrefixes: []*string{to.Ptr("10.1.0.0/16")},
						},
						RemoteVirtualNetwork: &armnetwork.SubResource{
							ID: to.Ptr("/subscriptions/00000000-0000-0000-0000-000000000002/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet2"),
						},
					},
				},
			},
		},
	}, *subscriptions[0])
	return []*VNETWrapper{&vnet1}
}

func init() {

	if clientSingleton == nil {
		subscriptions = dummySubscriptions()
		vnets = dummyVnets()
		subscriptions[0].vnets = append(subscriptions[0].vnets, vnets...)

		clientSingleton = newDocumentationClientWithData("docs", subscriptions, vnets)
	}
}
