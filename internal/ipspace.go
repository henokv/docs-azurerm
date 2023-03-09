package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/c-robinson/iplib"
)

type IPSpace struct {
	vnet         *VNETWrapper
	subnets      []*armnetwork.Subnet
	addressSpace iplib.Net4
}

func NewIPSPace(vnet *VNETWrapper, addressSpace string) *IPSpace {
	ipSpace := IPSpace{vnet: vnet, addressSpace: iplib.Net4FromStr(addressSpace)}
	return &ipSpace
}

func (ipSpace *IPSpace) AddSubnet(subnet *armnetwork.Subnet) bool {
	subnetPrefix := iplib.Net4FromStr(*subnet.Properties.AddressPrefix)
	if ipSpace.addressSpace.ContainsNet(subnetPrefix) {
		ipSpace.subnets = append(ipSpace.subnets, subnet)
		return true
	}
	return false
}
