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

var subscriptionList []*SubscriptionWrapper

func (sub *SubscriptionWrapper) generateMarkdown() string {
	var markdown string
	markdown += fmt.Sprintf("# %s  \n", *sub.DisplayName)
	markdown += fmt.Sprintf("#### ID: %s  \n", *sub.SubscriptionID)
	markdown += fmt.Sprintf("#### VNETs  \n")
	for _, vnet := range sub.vnets {
		markdown += fmt.Sprintf("- [%s](%s/%s.md).  \n", *vnet.Name, vnet.ResourceGroup, *vnet.Name)
	}
	return markdown
}

func (sub *SubscriptionWrapper) WriteMarkdown() error {
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

// GetCachedSubscriptionNameByID
// Function to check against local subscription store if
// subscriptionId : the resource id of the subscription
// name : the display name
// found : if the subscription can be found in cache
func GetCachedSubscriptionNameByID(subscriptionId string) (name string, found bool) {
	for _, subscription := range subscriptionList {
		if *subscription.SubscriptionID == subscriptionId {
			name = *subscription.DisplayName
			found = true
			break
		}
	}
	return name, found
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

// Cleans the docs directory
func CleanDocsDir() {
	os.RemoveAll("docs")
}

// GetAllSubscriptions returns all the subscriptions on which the authenticated user has permissions
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
	subscriptionList = subscriptions
	return subscriptions, nil
}
