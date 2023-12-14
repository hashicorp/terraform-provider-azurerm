package blobs

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var _ StorageBlob = Client{}

func TestLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "example.txt"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorageVTwo)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}

	containersClient, err := containers.NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", testData.StorageAccountName, *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(containersClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	_, err = containersClient.Create(ctx, containerName, containers.CreateInput{})
	if err != nil {
		t.Fatal(fmt.Errorf("error creating: %s", err))
	}
	defer containersClient.Delete(ctx, containerName)

	blobClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", testData.StorageAccountName, *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(blobClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	t.Logf("[DEBUG] Copying file to Blob Storage..")
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}

	if err := blobClient.CopyAndWait(ctx, containerName, fileName, copyInput); err != nil {
		t.Fatalf("Error copying: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Blob Properties..")
	details, err := blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}

	// default value
	if details.AccessTier != Hot {
		t.Fatalf("Expected the AccessTier to be %q but got %q", Hot, details.AccessTier)
	}
	if details.BlobType != BlockBlob {
		t.Fatalf("Expected BlobType to be %q but got %q", BlockBlob, details.BlobType)
	}
	if len(details.MetaData) != 0 {
		t.Fatalf("Expected there to be no items of metadata but got %d", len(details.MetaData))
	}

	t.Logf("[DEBUG] Checking it's returned in the List API..")
	listInput := containers.ListBlobsInput{}
	listResult, err := containersClient.ListBlobs(ctx, containerName, listInput)
	if err != nil {
		t.Fatalf("Error listing blobs: %s", err)
	}

	if model := listResult.Model; model != nil {
		if len(model.Blobs.Blobs) != 1 {
			t.Fatalf("Expected there to be 1 blob in the container but got %d", len(model.Blobs.Blobs))
		}
	}

	t.Logf("[DEBUG] Setting MetaData..")
	metaDataInput := SetMetaDataInput{
		MetaData: map[string]string{
			"hello": "there",
		},
	}
	if _, err := blobClient.SetMetaData(ctx, containerName, fileName, metaDataInput); err != nil {
		t.Fatalf("Error setting MetaData: %s", err)
	}

	t.Logf("[DEBUG] Re-retrieving Blob Properties..")
	details, err = blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error re-retrieving properties: %s", err)
	}

	// default value
	if details.AccessTier != Hot {
		t.Fatalf("Expected the AccessTier to be %q but got %q", Hot, details.AccessTier)
	}
	if details.BlobType != BlockBlob {
		t.Fatalf("Expected BlobType to be %q but got %q", BlockBlob, details.BlobType)
	}
	if len(details.MetaData) != 1 {
		t.Fatalf("Expected there to be 1 item of metadata but got %d", len(details.MetaData))
	}
	if details.MetaData["hello"] != "there" {
		t.Fatalf("Expected `hello` to be `there` but got %q", details.MetaData["there"])
	}

	t.Logf("[DEBUG] Retrieving the Block List..")
	getBlockListInput := GetBlockListInput{
		BlockListType: All,
	}
	blockList, err := blobClient.GetBlockList(ctx, containerName, fileName, getBlockListInput)
	if err != nil {
		t.Fatalf("Error retrieving Block List: %s", err)
	}

	// since this is a copy from an existing file, all blocks should be present
	if len(blockList.CommittedBlocks.Blocks) == 0 {
		t.Fatalf("Expected there to be committed blocks but there weren't!")
	}
	if len(blockList.UncommittedBlocks.Blocks) != 0 {
		t.Fatalf("Expected all blocks to be committed but got %d uncommitted blocks", len(blockList.UncommittedBlocks.Blocks))
	}

	t.Logf("[DEBUG] Changing the Access Tiers..")
	tiers := []AccessTier{
		Hot,
		Cool,
		Archive,
	}
	for _, tier := range tiers {
		t.Logf("[DEBUG] Updating the Access Tier to %q..", string(tier))
		if _, err := blobClient.SetTier(ctx, containerName, fileName, SetTierInput{Tier: tier}); err != nil {
			t.Fatalf("Error setting the Access Tier: %s", err)
		}

		t.Logf("[DEBUG] Re-retrieving Blob Properties..")
		details, err = blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
		if err != nil {
			t.Fatalf("Error re-retrieving properties: %s", err)
		}

		if details.AccessTier != tier {
			t.Fatalf("Expected the AccessTier to be %q but got %q", tier, details.AccessTier)
		}
	}

	t.Logf("[DEBUG] Deleting Blob")
	if _, err := blobClient.Delete(ctx, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting Blob: %s", err)
	}
}
