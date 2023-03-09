package internal

import (
	"fmt"
	"github.com/c-robinson/iplib"
	"net"
	"os"
	"strings"
)

func WriteToFile(content string, destination string) (err error) {
	path := destination[:strings.LastIndex(destination, "/")]
	os.MkdirAll(path, 0770)
	file, err := os.Create(destination)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
		file, err = os.Open(destination)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func coalesceString(value string, fallback string) string {
	if len(value) == 0 {
		return fallback
	} else {
		return ""
	}
}

func GenerateMarkdown(subs []*SubscriptionWrapper) string {
	var markdown string
	markdown += fmt.Sprintf("# Subscriptions  \n")
	for _, sub := range subs {
		if len(sub.vnets) > 0 {
			markdown += fmt.Sprintf("- [%s](%s/Readme.md).  \n", *sub.DisplayName, strings.ReplaceAll(*sub.DisplayName, " ", "%20"))
		}
	}
	return markdown
}

func WriteMarkdown(subs []*SubscriptionWrapper) error {
	markdown := GenerateMarkdown(subs)
	err := WriteToFile(markdown, fmt.Sprintf("docs/Readme.md"))
	return err
}

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
