package paths

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestCreateDirectory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())
	path := "test"

	testData, err := client.BuildTestResourcesWithHns(ctx, resourceGroup, accountName, storageaccounts.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}

	baseUri := fmt.Sprintf("https://%s.%s.%s", accountName, "dfs", *domainSuffix)

	fileSystemsClient, err := filesystems.NewWithBaseUri(baseUri)
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(fileSystemsClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	t.Logf("[DEBUG] Creating an empty File System..")
	fileSystemInput := filesystems.CreateInput{
		Properties: map[string]string{},
	}
	if _, err = fileSystemsClient.Create(ctx, fileSystemName, fileSystemInput); err != nil {
		t.Fatal(fmt.Errorf("error creating: %s", err))
	}

	t.Logf("[DEBUG] Creating path..")

	pathsClient, err := NewWithBaseUri(baseUri)
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(pathsClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	input := CreateInput{
		Resource: PathResourceDirectory,
	}

	if _, err = pathsClient.Create(ctx, fileSystemName, path, input); err != nil {
		t.Fatal(fmt.Errorf("error creating path: %s", err))
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, fileSystemName); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
