package internal

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

// var subscriptionList []*SubscriptionWrapper

func (sub *SubscriptionWrapper) generateMarkdown() string {
	var markdown string
	markdown += sub.GenerateTitle(*sub.DisplayName, 1)
	markdown += sub.GenerateTitle(fmt.Sprintf("ID: %s", *sub.SubscriptionID), 4)
	markdown += sub.GenerateTitle("VNETs", 4)
	for _, vnet := range sub.vnets {
		markdown += sub.GenerateListItem(sub.GenerateLink(*vnet.Name, fmt.Sprintf("%s/%s.md", vnet.ResourceGroup, *vnet.Name)))
	}
	return markdown
}

func (sub *SubscriptionWrapper) WriteMarkdown() error {
	if len(sub.vnets) > 0 {
		markdown := sub.generateMarkdown()
		err := sub.writeToFile(markdown, fmt.Sprintf("docs/%s/Readme.md", *sub.DisplayName))
		return err
	}
	return nil
}

type SubscriptionWrapper struct {
	*Markdown
	*armsubscription.Subscription
	vnets []*VNETWrapper
}

func NewSubscriptionWrapper(sub armsubscription.Subscription) SubscriptionWrapper {
	return SubscriptionWrapper{NewMarkdown(), &sub, []*VNETWrapper{}}
}

func (sub *SubscriptionWrapper) getWrappedVNETsInSubscription() (vnets []*VNETWrapper, err error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return vnets, err
	}
	client, err := armnetwork.NewVirtualNetworksClient(*sub.SubscriptionID, cred, nil)
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
			wrappedVNET := NewVNETWrapper(vnet, *sub)
			vnets = append(vnets, &wrappedVNET)
			sub.vnets = append(sub.vnets, &wrappedVNET)
		}
	}
	return vnets, nil
}

func GetCachedSubscriptionNameByID(subscriptionId string) (name string, found bool) {
	for _, subscription := range client.GetSubscriptions() {
		if *subscription.SubscriptionID == subscriptionId {
			name = *subscription.DisplayName
			found = true
			break
		}
	}
	return name, found
}
