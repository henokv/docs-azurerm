package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/c-robinson/iplib"
	"testing"
)

func createTestSubnet(subnetPrefix string) *armnetwork.Subnet {
	properties := armnetwork.SubnetPropertiesFormat{}
	properties.AddressPrefix = &subnetPrefix
	subnet := armnetwork.Subnet{nil, nil, &properties, nil, nil}
	return &subnet
}

func TestIPSpace_GetFreeIPs(t *testing.T) {
	ipspace := IPSpace{addressSpace: iplib.Net4FromStr("10.0.0.0/16"), vnet: nil}
	subnets := []*armnetwork.Subnet{createTestSubnet("10.0.1.128/25")}
	for _, subnet := range subnets {
		ipspace.AddSubnet(subnet)
	}
	ipspace.generateFreeIPs()
	a := ipspace.freeSpace
	if len(a) != 9 {
		t.Fatalf("9 IPs expected but only got %d", len(a))
	}
}
