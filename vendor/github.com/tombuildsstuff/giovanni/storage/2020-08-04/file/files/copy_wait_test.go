package files

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestFilesCopyAndWaitFromURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	sharesClient, err := shares.NewWithBaseUri(fmt.Sprintf("https://%s.file.%s", accountName, *domainSuffix))
	if err := client.PrepareWithSharedKeyAuth(sharesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	input := shares.CreateInput{
		QuotaInGB: 10,
	}
	_, err = sharesClient.Create(ctx, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, shareName, shares.DeleteInput{DeleteSnapshots: false})

	filesClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.file.%s", accountName, *domainSuffix))
	if err := client.PrepareWithSharedKeyAuth(filesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	copiedFileName := "ubuntu.iso"
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}

	t.Logf("[DEBUG] Copy And Waiting..")
	if _, err := filesClient.CopyAndWait(ctx, shareName, "", copiedFileName, copyInput); err != nil {
		t.Fatalf("Error copy & waiting: %s", err)
	}

	t.Logf("[DEBUG] Asserting that the file's ready..")

	props, err := filesClient.GetProperties(ctx, shareName, "", copiedFileName)
	if err != nil {
		t.Fatalf("Error retrieving file: %s", err)
	}

	if !strings.EqualFold(props.CopyStatus, "success") {
		t.Fatalf("Expected the Copy Status to be `Success` but got %q", props.CopyStatus)
	}
}

func TestFilesCopyAndWaitFromBlob(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storageaccounts.KindStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}
	sharesClient, err := shares.NewWithBaseUri(fmt.Sprintf("https://%s.file.%s", accountName, *domainSuffix))
	if err := client.PrepareWithSharedKeyAuth(sharesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	input := shares.CreateInput{
		QuotaInGB: 10,
	}
	_, err = sharesClient.Create(ctx, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, shareName, shares.DeleteInput{DeleteSnapshots: false})

	baseUri := fmt.Sprintf("https://%s.file.%s", accountName, *domainSuffix)
	filesClient, err := NewWithBaseUri(baseUri)
	if err := client.PrepareWithSharedKeyAuth(filesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	originalFileName := "ubuntu.iso"
	copiedFileName := "ubuntu-copied.iso"
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/14.04/ubuntu-14.04.6-desktop-amd64.iso",
	}
	t.Logf("[DEBUG] Copy And Waiting the original file..")
	if _, err := filesClient.CopyAndWait(ctx, shareName, "", originalFileName, copyInput); err != nil {
		t.Fatalf("Error copy & waiting: %s", err)
	}

	t.Logf("[DEBUG] Now copying that blob..")
	duplicateInput := CopyInput{
		CopySource: fmt.Sprintf("%s/%s/%s", baseUri, shareName, originalFileName),
	}
	if _, err := filesClient.CopyAndWait(ctx, shareName, "", copiedFileName, duplicateInput); err != nil {
		t.Fatalf("Error copying duplicate: %s", err)
	}

	t.Logf("[DEBUG] Asserting that the file's ready..")
	props, err := filesClient.GetProperties(ctx, shareName, "", copiedFileName)
	if err != nil {
		t.Fatalf("Error retrieving file: %s", err)
	}

	if !strings.EqualFold(props.CopyStatus, "success") {
		t.Fatalf("Expected the Copy Status to be `Success` but got %q", props.CopyStatus)
	}
}
