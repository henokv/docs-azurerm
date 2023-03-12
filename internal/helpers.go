package internal

import (
	"github.com/c-robinson/iplib"
	"net"
)

func getFreeIPInfo(firstIP, lastIP net.IP) (nets []iplib.Net4) {
	largestBlock, done, err := iplib.NewNetBetween(firstIP, lastIP)
	if err != nil {
		panic(err)
	}
	block := iplib.Net4FromStr(largestBlock.String())
	nets = append(nets, block)
	if done {
		return nets
	} else {
		return append(nets, getFreeIPInfo(iplib.IncrementIPBy(block.BroadcastAddress(), 1), lastIP)...)
	}
}
