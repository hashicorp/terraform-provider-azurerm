package accounts

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestContainerLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())

	_, err = client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorageVTwo)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	accountsClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", accountName, *domainSuffix))
	if err != nil {
		t.Fatal(fmt.Errorf("building client for environment: %+v", err))
	}
	client.PrepareWithResourceManagerAuth(accountsClient.Client)

	input := StorageServiceProperties{}
	_, err = accountsClient.SetServiceProperties(ctx, accountName, input)
	if err != nil {
		t.Fatal(fmt.Errorf("error setting properties: %s", err))
	}

	var index = "index.html"
	//var enabled = true
	var errorDocument = "404.html"

	input = StorageServiceProperties{
		StaticWebsite: &StaticWebsite{
			Enabled:              true,
			IndexDocument:        index,
			ErrorDocument404Path: errorDocument,
		},
		Logging: &Logging{
			Version: "2.0",
			Delete:  true,
			Read:    true,
			Write:   true,
			RetentionPolicy: DeleteRetentionPolicy{
				Enabled: true,
				Days:    7,
			},
		},
	}

	_, err = accountsClient.SetServiceProperties(ctx, accountName, input)
	if err != nil {
		t.Fatal(fmt.Errorf("error setting properties: %s", err))
	}

	t.Log("[DEBUG] Waiting 2 seconds..")
	time.Sleep(2 * time.Second)

	result, err := accountsClient.GetServiceProperties(ctx, accountName)
	if err != nil {
		t.Fatal(fmt.Errorf("error getting properties: %s", err))
	}

	website := result.Model.StaticWebsite
	if website.Enabled != true {
		t.Fatalf("Expected the StaticWebsite %t but got %t", true, website.Enabled)
	}

	logging := result.Model.Logging
	if logging.Version != "2.0" {
		t.Fatalf("Expected the Logging Version %s but got %s", "2.0", logging.Version)
	}
	if !logging.Read {
		t.Fatalf("Expected the Logging Read %t but got %t", true, logging.Read)
	}
	if !logging.Write {
		t.Fatalf("Expected the Logging Write %t but got %t", true, logging.Write)
	}
	if !logging.Delete {
		t.Fatalf("Expected the Logging Delete %t but got %t", true, logging.Delete)
	}
	if !logging.RetentionPolicy.Enabled {
		t.Fatalf("Expected the Logging RetentionPolicy.Enabled %t but got %t", true, logging.RetentionPolicy.Enabled)
	}
	if logging.RetentionPolicy.Days != 7 {
		t.Fatalf("Expected the Logging RetentionPolicy.Enabled %d but got %d", 7, logging.RetentionPolicy.Days)
	}
}
