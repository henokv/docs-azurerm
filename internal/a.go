package internal

import (
	"fmt"
	"github.com/c-robinson/iplib"
	"log"
	"net"
	"sort"
)

func avc() {
	space := iplib.Net4FromStr("10.10.0.0/16")
	sn1 := iplib.Net4FromStr("10.10.10.0/24")
	sn2 := iplib.Net4FromStr("10.10.20.0/23")
	sn3 := iplib.Net4FromStr("10.10.8.0/27")
	sn4 := iplib.Net4FromStr("10.20.8.0/27")
	subnets := []iplib.Net4{sn1, sn2, sn3, sn4}

	b, _ := sn1.Supernet(16)
	c := fmt.Sprintf("%v %v %v", iplib.IncrementIPBy(sn1.BroadcastAddress(), 1),
		sn1.NetworkAddress(),
		b,
	)
	getFreeIPs(space, subnets)

	log.Printf("%v", c)
}

func getFreeIPs(space iplib.Net4, subnets []iplib.Net4) (ips []string) {
	var contains []iplib.Net4
	for _, subnet := range subnets {
		if space.ContainsNet(subnet) {
			contains = append(contains, subnet)
		}
	}
	sort.Slice(contains, func(i, j int) bool {
		return iplib.CompareNets(contains[i], contains[j]) == -1
	})
	if len(contains) == 0 { // If no subnets present the entire range is free
		ips = append(ips, getFreeIPInfo(space.NetworkAddress(), space.BroadcastAddress()))
		return ips
	} else if iplib.CompareNets(space, contains[0]) == 0 { // If the first subnet contains the entire range no ips are free
		return ips
	} else if iplib.CompareIPs(space.NetworkAddress(), contains[0].NetworkAddress()) != 0 { //Compare the first ip of the first subnet against the first ip of the address space
		ips = append(ips, getFreeIPInfo(space.NetworkAddress(), iplib.DecrementIPBy(contains[0].NetworkAddress(), 1)))
	}
	for i := 0; i < len(contains); i++ {
		current := contains[i]
		if i == len(contains)-1 {
			if iplib.CompareIPs(space.BroadcastAddress(), current.BroadcastAddress()) != 0 {
				ips = append(ips, getFreeIPInfo(iplib.IncrementIPBy(current.BroadcastAddress(), 1), space.BroadcastAddress()))
				break
			}
		} else {
			next := contains[i+1]
			if iplib.CompareIPs(iplib.IncrementIPBy(current.BroadcastAddress(), 1), next.NetworkAddress()) != 0 {
				ips = append(ips, getFreeIPInfo(iplib.IncrementIPBy(current.BroadcastAddress(), 1), iplib.DecrementIPBy(next.NetworkAddress(), 1)))
			}
		}
	}
	return ips
}

func getFreeIPInfo(firstIP, secondIP net.IP) string {
	largestBlock, _, err := iplib.NewNetBetween(firstIP, secondIP)
	info := ""
	if err == nil {
		info = fmt.Sprintf("largest: %s (%d ips)", largestBlock.String(), iplib.DeltaIP4(firstIP, secondIP)+1)
	}
	return fmt.Sprintf("%s-%s %s", firstIP, secondIP, info)
}
