package internal

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
)

func (sub SubscriptionWrapper) GenerateMarkdown() string {
	var markdown string
	markdown += fmt.Sprintf("# %s  \n", *sub.DisplayName)
	markdown += fmt.Sprintf("#### ID: %s  \n", *sub.SubscriptionID)
	markdown += fmt.Sprintf("#### VNETs  \n")
	for _, vnet := range sub.vnets {
		markdown += fmt.Sprintf("- [%s](%s/%s.md).  \n", *vnet.Name, vnet.ResourceGroup, *vnet.Name)
	}
	return markdown
}

func (sub SubscriptionWrapper) WriteMarkdown() error {
	if len(sub.vnets) > 0 {
		markdown := sub.GenerateMarkdown()
		err := WriteToFile(markdown, fmt.Sprintf("docs/%s/Readme.md", *sub.DisplayName))
		return err
	}
	return nil
}

type SubscriptionWrapper struct {
	armsubscription.Subscription
	vnets []*VNETWrapper
}

func NewSubscriptionWrapper(sub armsubscription.Subscription) SubscriptionWrapper {
	return SubscriptionWrapper{sub, []*VNETWrapper{}}
}

func GetAllSubscriptionsAsString() (subscriptions []string, err error) {
	subscriptionsResources, err := GetAllSubscriptions()
	if err != nil {
		return nil, err
	}
	for _, subscription := range subscriptionsResources {
		subscriptions = append(subscriptions, *subscription.SubscriptionID)
	}
	return subscriptions, nil
}

func GetAllSubscriptions() (subscriptions []*SubscriptionWrapper, err error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return subscriptions, err
	}
	client, err := armsubscription.NewSubscriptionsClient(cred, nil)
	if err != nil {
		return subscriptions, err
	}
	pager := client.NewListPager(nil)
	for pager.More() {
		subList, err := pager.NextPage(context.Background())
		if err != nil {
			return subscriptions, err
		}
		for _, subscription := range subList.ListResult.Value {
			if true {
				wrappedSub := NewSubscriptionWrapper(*subscription)
				subscriptions = append(subscriptions, &wrappedSub)
			}
		}
	}
	return subscriptions, nil
}
