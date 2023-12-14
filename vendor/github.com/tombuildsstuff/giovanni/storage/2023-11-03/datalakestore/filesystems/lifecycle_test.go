package filesystems

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	fileSystemsClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.%s.%s", accountName, "dfs", *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(fileSystemsClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	t.Logf("[DEBUG] Creating an empty File System..")
	input := CreateInput{
		Properties: map[string]string{
			"hello": "aGVsbG8=",
		},
	}
	if _, err = fileSystemsClient.Create(ctx, fileSystemName, input); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Retrieving the Properties..")
	props, err := fileSystemsClient.GetProperties(ctx, fileSystemName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	}

	if len(props.Properties) != 1 {
		t.Fatalf("Expected 1 properties by default but got %d", len(props.Properties))
	}
	if props.Properties["hello"] != "aGVsbG8=" {
		t.Fatalf("Expected `hello` to be `aGVsbG8=` but got %q", props.Properties["hello"])
	}

	t.Logf("[DEBUG] Updating the properties..")
	setInput := SetPropertiesInput{
		Properties: map[string]string{
			"hello":   "d29uZGVybGFuZA==",
			"private": "ZXll",
		},
	}
	if _, err := fileSystemsClient.SetProperties(ctx, fileSystemName, setInput); err != nil {
		t.Fatalf("Error setting properties: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving the Properties..")
	props, err = fileSystemsClient.GetProperties(ctx, fileSystemName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	}
	if len(props.Properties) != 2 {
		t.Fatalf("Expected 2 properties by default but got %d", len(props.Properties))
	}
	if props.Properties["hello"] != "d29uZGVybGFuZA==" {
		t.Fatalf("Expected `hello` to be `d29uZGVybGFuZA==` but got %q", props.Properties["hello"])
	}
	if props.Properties["private"] != "ZXll" {
		t.Fatalf("Expected `private` to be `ZXll` but got %q", props.Properties["private"])
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, fileSystemName); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
