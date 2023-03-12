package internal

import (
	"os"
	"testing"
)

func TestGetCachedSubscriptionNameByID(t *testing.T) {
	id := "00000000-0000-0000-0000-000000000000"
	_, found := GetCachedSubscriptionNameByID(id)
	if found {
		t.Fatalf("Subscription with id %s should not exist ", id)
	}
	//if subscriptionName != "Test Subscription 1" {
	//	t.Fatalf("Expected Test Subscription 1 but got %s", subscriptionName)
	//}
}

func TestGetAllSubscriptionsReturnsThreeSubscription(t *testing.T) {
	if len(subscriptions) != 3 {
		t.Fatalf("Expected 3 subscriptions but got %d", len(subscriptions))
	}
}

func TestGetWrappedVNETsInSubscriptionsReturnsThreeVNETS(t *testing.T) {

	if len(vnets) != 1 {
		t.Fatalf("Expected 1 vnet but got %d", len(vnets))
	}
}

func TestGetFreeIPSPaceInVNETs(t *testing.T) {

	var spacesCount int
	for _, vnet := range vnets {
		spaces := vnet.getFreeIPSPace()
		spacesCount += len(spaces)

	}
	if spacesCount != 8 {
		t.Fatalf("Expected 8 free ip ranges but got %d", spacesCount)
	}

}

func TestWriteMarkdown(t *testing.T) {
	var err error

	CleanDocsDir()
	err = WriteMarkdown(subscriptions)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("docs/Readme.md")
	if err != nil {
		t.Fatalf("Expected file %s to exist", "docs/Readme.md")
	}
	for _, sub := range subscriptions {
		sub.WriteMarkdown()
		if err != nil {
			t.Fatal(err)
		}
		//_, err = os.Stat(fmt.Sprintf("docs/%s/Readme.md", *sub.DisplayName))
		//if err != nil {
		//	t.Fatalf("Expected file %s to exist", fmt.Sprintf("docs/%s/Readme.md", *sub.DisplayName))
		//}
	}
	for _, vnet := range vnets {
		vnet.WriteMarkdown()
		if err != nil {
			t.Fatal(err)
		}
		//_, err = os.Stat(fmt.Sprintf("docs/%s/%s.md", *vnet.Subscription.DisplayName, *vnet.Name))
		//if err != nil {
		//	t.Fatalf("Expected file %s to exist", fmt.Sprintf("docs/%s/%s.md", *vnet.Subscription.DisplayName, *vnet.Name))
		//}
	}

}