package files

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/storage/internal/endpoints"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestFilesCopyAndWaitFromURL(t *testing.T) {
	client, err := testhelpers.Build()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	sharesClient := shares.NewWithEnvironment(client.Environment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := shares.CreateInput{
		QuotaInGB: 10,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, accountName, shareName, false)

	filesClient := NewWithEnvironment(client.Environment)
	filesClient.Client = client.PrepareWithAuthorizer(filesClient.Client, storageAuth)

	copiedFileName := "ubuntu.iso"
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/18.04.2/ubuntu-18.04.2-desktop-amd64.iso",
	}

	t.Logf("[DEBUG] Copy And Waiting..")
	if _, err := filesClient.CopyAndWait(ctx, accountName, shareName, "", copiedFileName, copyInput, DefaultCopyPollDuration); err != nil {
		t.Fatalf("Error copy & waiting: %s", err)
	}

	t.Logf("[DEBUG] Asserting that the file's ready..")

	props, err := filesClient.GetProperties(ctx, accountName, shareName, "", copiedFileName)
	if err != nil {
		t.Fatalf("Error retrieving file: %s", err)
	}

	if !strings.EqualFold(props.CopyStatus, "success") {
		t.Fatalf("Expected the Copy Status to be `Success` but got %q", props.CopyStatus)
	}
}

func TestFilesCopyAndWaitFromBlob(t *testing.T) {
	client, err := testhelpers.Build()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.TODO()
	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	shareName := fmt.Sprintf("share-%d", testhelpers.RandomInt())

	testData, err := client.BuildTestResources(ctx, resourceGroup, accountName, storage.Storage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)

	storageAuth := auth.NewSharedKeyLiteAuthorizer(accountName, testData.StorageAccountKey)
	sharesClient := shares.NewWithEnvironment(client.Environment)
	sharesClient.Client = client.PrepareWithAuthorizer(sharesClient.Client, storageAuth)

	input := shares.CreateInput{
		QuotaInGB: 10,
	}
	_, err = sharesClient.Create(ctx, accountName, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, accountName, shareName, false)

	filesClient := NewWithEnvironment(client.Environment)
	filesClient.Client = client.PrepareWithAuthorizer(filesClient.Client, storageAuth)

	originalFileName := "ubuntu.iso"
	copiedFileName := "ubuntu-copied.iso"
	copyInput := CopyInput{
		CopySource: "http://releases.ubuntu.com/18.04.2/ubuntu-18.04.2-desktop-amd64.iso",
	}
	t.Logf("[DEBUG] Copy And Waiting the original file..")
	if _, err := filesClient.CopyAndWait(ctx, accountName, shareName, "", originalFileName, copyInput, DefaultCopyPollDuration); err != nil {
		t.Fatalf("Error copy & waiting: %s", err)
	}

	t.Logf("[DEBUG] Now copying that blob..")
	duplicateInput := CopyInput{
		CopySource: fmt.Sprintf("%s/%s/%s", endpoints.GetFileEndpoint(filesClient.BaseURI, accountName), shareName, originalFileName),
	}
	if _, err := filesClient.CopyAndWait(ctx, accountName, shareName, "", copiedFileName, duplicateInput, DefaultCopyPollDuration); err != nil {
		t.Fatalf("Error copying duplicate: %s", err)
	}

	t.Logf("[DEBUG] Asserting that the file's ready..")
	props, err := filesClient.GetProperties(ctx, accountName, shareName, "", copiedFileName)
	if err != nil {
		t.Fatalf("Error retrieving file: %s", err)
	}

	if !strings.EqualFold(props.CopyStatus, "success") {
		t.Fatalf("Expected the Copy Status to be `Success` but got %q", props.CopyStatus)
	}
}
