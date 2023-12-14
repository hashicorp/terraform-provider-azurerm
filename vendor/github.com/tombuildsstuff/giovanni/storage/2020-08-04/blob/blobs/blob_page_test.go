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

func TestPageBlobLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "append-blob.txt"

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
		t.Fatal(fmt.Errorf("Error creating: %s", err))
	}
	defer containersClient.Delete(ctx, containerName)

	blobClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.blob.%s", testData.StorageAccountName, *domainSuffix))
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}

	if err := client.PrepareWithSharedKeyAuth(blobClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	t.Logf("[DEBUG] Putting Page Blob..")
	fileSize := int64(10240000)
	if _, err := blobClient.PutPageBlob(ctx, containerName, fileName, PutPageBlobInput{
		BlobContentLengthBytes: fileSize,
	}); err != nil {
		t.Fatalf("Error putting page blob: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties..")
	props, err := blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error retrieving properties: %s", err)
	}
	if props.ContentLength != fileSize {
		t.Fatalf("Expected Content-Length to be %d but it was %d", fileSize, props.ContentLength)
	}

	for iteration := 1; iteration <= 3; iteration++ {
		t.Logf("[DEBUG] Putting Page %d of 3..", iteration)
		byteArray := func() []byte {
			o := make([]byte, 0)

			for i := 0; i < 512; i++ {
				o = append(o, byte(i))
			}

			return o
		}()
		startByte := int64(512 * iteration)
		endByte := int64(startByte + 511)
		putPageInput := PutPageUpdateInput{
			StartByte: startByte,
			EndByte:   endByte,
			Content:   byteArray,
		}
		if _, err := blobClient.PutPageUpdate(ctx, containerName, fileName, putPageInput); err != nil {
			t.Fatalf("Error putting page: %s", err)
		}
	}

	t.Logf("[DEBUG] Deleting..")
	if _, err := blobClient.Delete(ctx, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting: %s", err)
	}
}
