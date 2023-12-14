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

func TestCreateHasNoTagsByDefault(t *testing.T) {
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
		Properties: map[string]string{},
	}
	if _, err = fileSystemsClient.Create(ctx, fileSystemName, input); err != nil {
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}

	t.Logf("[DEBUG] Retrieving the Properties..")
	props, err := fileSystemsClient.GetProperties(ctx, fileSystemName)
	if err != nil {
		t.Fatal(fmt.Errorf("Error getting properties: %s", err))
	}

	if len(props.Properties) != 0 {
		t.Fatalf("Expected 0 properties by default but got %d", len(props.Properties))
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, fileSystemName); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
