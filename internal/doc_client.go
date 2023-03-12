package internal

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/subscription/armsubscription"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"strings"
)

var clientSingleton *DocumentationClient

type DocumentationClient struct {
	*Markdown
	subscriptions []*SubscriptionWrapper
	vnets         []*VNETWrapper
	docsDir       string
}

func GetSingletonDocumentationClient() (*DocumentationClient, error) {
	if clientSingleton == nil {
		docsDir := viper.GetString("docsDir")
		if docsDir == "" {
			docsDir = "docs"
		}
		var err error
		clientSingleton, err = newDocumentationClient(docsDir)
		if err != nil {
			return nil, err
		}
	}
	return clientSingleton, nil
}

func newDocumentationClient(docsDir string) (*DocumentationClient, error) {
	clientSingleton = &DocumentationClient{
		Markdown: NewMarkdown(),
		docsDir:  docsDir,
	}
	subscriptions, err := getAllSubscriptions()
	if err != nil {
		return nil, err
	}
	clientSingleton.subscriptions = subscriptions
	vnets, err := GetWrappedVNETsInSubscriptions(subscriptions)
	if err != nil {
		return nil, err
	}
	clientSingleton.vnets = vnets
	return clientSingleton, nil
}

func (client *DocumentationClient) GetSubscriptions() []*SubscriptionWrapper {
	return client.subscriptions
}

func (client *DocumentationClient) GetSubscriptionById(id string) *SubscriptionWrapper {
	for _, subscription := range client.subscriptions {
		if *subscription.SubscriptionID == id {
			return subscription
		}
	}
	return nil
}

func (client *DocumentationClient) GetVNETs() []*VNETWrapper {
	return client.vnets
}

//func (client *DocumentationClient) AddVNET(vnet *VNETWrapper) {
//	client.vnets = append(client.vnets, vnet)
//}

func (client *DocumentationClient) GenerateMarkdown(verbose bool) error {
	if verbose {
		log.Println("Generating markdown files")
	}
	client.CleanDocsDir()
	err := client.WriteDocumentation()
	if err != nil {
		return err
	}
	for _, sub := range client.GetSubscriptions() {
		err = sub.WriteMarkdown()
		if err != nil {
			return err
		}
	}
	for _, vnet := range client.GetVNETs() {
		vnet.WriteMarkdown()
		if err != nil {
			return err
		}
	}
	if verbose {
		log.Println("Markdown files generated")
	}
	return nil
}

func (client *DocumentationClient) generateDocumentation() string {
	var markdown string
	markdown += client.GenerateTitle("Subscriptions", 1)
	for _, sub := range client.GetSubscriptions() {
		if len(sub.vnets) > 0 {
			markdown += client.GenerateListItem(
				client.GenerateLink(*sub.DisplayName, fmt.Sprintf("%s/Readme.md", strings.ReplaceAll(*sub.DisplayName, " ", "%20"))),
			)
		}
	}
	return markdown
}

func (client *DocumentationClient) WriteDocumentation() error {
	markdown := client.generateDocumentation()
	err := client.writeToFile(markdown, fmt.Sprintf("%s/Readme.md", client.GetDocsDir()))
	return err
}

func (client *DocumentationClient) GetDocsDir() string {
	return client.docsDir
}

func (client *DocumentationClient) CleanDocsDir() {
	os.RemoveAll(client.docsDir)
}

func (client *DocumentationClient) GetSubscriptionNameByID(subscriptionId string) (name string, found bool) {
	for _, subscription := range client.GetSubscriptions() {
		if *subscription.SubscriptionID == subscriptionId {
			name = *subscription.DisplayName
			found = true
			break
		}
	}
	return name, found
}

func getAllSubscriptions() (subscriptions []*SubscriptionWrapper, err error) {
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
	//subscriptionList = subscriptions
	return subscriptions, nil
}

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
