package internal

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"os"
)

func (sub SubscriptionWrapper) generateMarkdown() string {
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
		markdown := sub.generateMarkdown()
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

//func GetAllSubscriptionsAsString() (subscriptions []string, err error) {
//	subscriptionsResources, err := GetAllSubscriptions()
//	if err != nil {
//		return nil, err
//	}
//	for _, subscription := range subscriptionsResources {
//		subscriptions = append(subscriptions, *subscription.SubscriptionID)
//	}
//	return subscriptions, nil
//}

func subscriptionNeeded(subscription *armsubscription.Subscription) bool {
	value := *subscription.SubscriptionID
	if viper.GetString("subscriptions.key") == "name" {
		value = *subscription.DisplayName
	}
	include := viper.GetStringSlice("subscriptions.include")
	if len(include) > 0 {
		return slices.Contains(include, value)
	}
	if slices.Contains(viper.GetStringSlice("subscriptions.exclude"), value) {
		return false
	}
	return true
}

func CleanDocsDir() {
	os.RemoveAll("docs")
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
			if subscriptionNeeded(subscription) {
				wrappedSub := NewSubscriptionWrapper(*subscription)
				subscriptions = append(subscriptions, &wrappedSub)
			}
		}
	}
	return subscriptions, nil
}
