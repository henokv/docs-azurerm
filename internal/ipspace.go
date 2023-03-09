package internal

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/c-robinson/iplib"
	"sort"
)

type IPSpace struct {
	vnet         *VNETWrapper
	subnets      []*armnetwork.Subnet
	addressSpace iplib.Net4
	freeSpace    []iplib.Net4
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

func (ipSpace *IPSpace) generateFreeIPs() {
	var ips []iplib.Net4
	var contains []iplib.Net4
	for _, subnet := range ipSpace.subnets {
		subnetPrefix := iplib.Net4FromStr(*subnet.Properties.AddressPrefix)
		if ipSpace.addressSpace.ContainsNet(subnetPrefix) {
			contains = append(contains, subnetPrefix)
		}
	}
	sort.Slice(contains, func(i, j int) bool {
		return iplib.CompareNets(contains[i], contains[j]) == -1
	})
	if len(contains) == 0 { // If no subnets present the entire range is free
		ips = append(ips, getFreeIPInfo(ipSpace.addressSpace.NetworkAddress(), ipSpace.addressSpace.BroadcastAddress())...)
		ipSpace.freeSpace = ips
		return
	} else if iplib.CompareNets(ipSpace.addressSpace, contains[0]) == 0 { // If the first subnet contains the entire range no ips are free
		ipSpace.freeSpace = ips
		return
	} else if iplib.CompareIPs(ipSpace.addressSpace.NetworkAddress(), contains[0].NetworkAddress()) != 0 { //Compare the first ip of the first subnet against the first ip of the address space
		ips = append(ips, getFreeIPInfo(ipSpace.addressSpace.NetworkAddress(), iplib.DecrementIPBy(contains[0].NetworkAddress(), 1))...)
	}
	for i := 0; i < len(contains); i++ {
		current := contains[i]
		if i == len(contains)-1 {
			if iplib.CompareIPs(ipSpace.addressSpace.BroadcastAddress(), current.BroadcastAddress()) != 0 {
				ips = append(ips, getFreeIPInfo(iplib.IncrementIPBy(current.BroadcastAddress(), 1), ipSpace.addressSpace.BroadcastAddress())...)
				break
			}
		} else {
			next := contains[i+1]
			if iplib.CompareIPs(iplib.IncrementIPBy(current.BroadcastAddress(), 1), next.NetworkAddress()) != 0 {
				ips = append(ips, getFreeIPInfo(iplib.IncrementIPBy(current.BroadcastAddress(), 1), iplib.DecrementIPBy(next.NetworkAddress(), 1))...)
			}
		}
	}
	ipSpace.freeSpace = ips
	return
}
