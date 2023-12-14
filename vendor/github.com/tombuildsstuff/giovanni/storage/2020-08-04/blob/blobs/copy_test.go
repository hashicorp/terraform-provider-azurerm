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

func TestCopyFromExistingFile(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "ubuntu.iso"
	copiedFileName := "copied.iso"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindBlobStorage)
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

	baseUri := fmt.Sprintf("https://%s.blob.%s", testData.StorageAccountName, *domainSuffix)
	blobClient, err := NewWithBaseUri(baseUri)
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

	t.Logf("[DEBUG] Duplicating that file..")
	copiedInput := CopyInput{
		CopySource: fmt.Sprintf("%s/%s/%s", baseUri, containerName, fileName),
	}
	if err := blobClient.CopyAndWait(ctx, containerName, copiedFileName, copiedInput); err != nil {
		t.Fatalf("Error duplicating file: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties for the Original File..")
	props, err := blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties for the original file: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties for the Copied File..")
	copiedProps, err := blobClient.GetProperties(ctx, containerName, copiedFileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties for the copied file: %s", err)
	}

	if props.ContentLength != copiedProps.ContentLength {
		t.Fatalf("Expected the content length to be %d but it was %d", props.ContentLength, copiedProps.ContentLength)
	}

	t.Logf("[DEBUG] Deleting copied file..")
	if _, err := blobClient.Delete(ctx, containerName, copiedFileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting file: %s", err)
	}

	t.Logf("[DEBUG] Deleting original file..")
	if _, err := blobClient.Delete(ctx, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting file: %s", err)
	}
}

func TestCopyFromURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	containerName := fmt.Sprintf("cont-%d", testhelpers.RandomInt())
	fileName := "ubuntu.iso"

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindBlobStorage)
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

	t.Logf("[DEBUG] Copying file to Blob Storage..")
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}

	if err := blobClient.CopyAndWait(ctx, containerName, fileName, copyInput); err != nil {
		t.Fatalf("Error copying: %s", err)
	}

	t.Logf("[DEBUG] Retrieving Properties..")
	props, err := blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties: %s", err)
	}

	if props.ContentLength == 0 {
		t.Fatalf("Expected the file to be there but looks like it isn't: %d", props.ContentLength)
	}

	t.Logf("[DEBUG] Deleting file..")
	if _, err := blobClient.Delete(ctx, containerName, fileName, DeleteInput{}); err != nil {
		t.Fatalf("Error deleting file: %s", err)
	}
}
