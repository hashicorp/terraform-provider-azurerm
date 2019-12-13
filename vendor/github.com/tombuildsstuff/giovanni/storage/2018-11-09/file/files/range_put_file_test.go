package files

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/tombuildsstuff/giovanni/storage/2018-11-09/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/auth"
	"github.com/tombuildsstuff/giovanni/testhelpers"
)

func TestPutSmallFile(t *testing.T) {
	// the purpose of this test is to ensure that a small file (< 4MB) is a single chunk
	testPutFile(t, "small-file.png", "image/png")
}

func TestPutLargeFile(t *testing.T) {
	// the purpose of this test is to ensure that large files (> 4MB) are chunked
	testPutFile(t, "blank-large-file.dmg", "application/x-apple-diskimage")
}

func testPutFile(t *testing.T, fileName string, contentType string) {
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

	// store files outside of this directory, since they're reused
	file, err := os.Open("../../../testdata/" + fileName)
	if err != nil {
		t.Fatalf("Error opening: %s", err)
	}

	info, err := file.Stat()
	if err != nil {
		t.Fatalf("Error 'stat'-ing: %s", err)
	}

	t.Logf("[DEBUG] Creating Top Level File..")
	createFileInput := CreateInput{
		ContentLength: info.Size(),
		ContentType:   &contentType,
	}
	if _, err := filesClient.Create(ctx, accountName, shareName, "", fileName, createFileInput); err != nil {
		t.Fatalf("Error creating Top-Level File: %s", err)
	}

	t.Logf("[DEBUG] Uploading File..")
	if err := filesClient.PutFile(ctx, accountName, shareName, "", fileName, file, 4); err != nil {
		t.Fatalf("Error uploading File: %s", err)
	}

	t.Logf("[DEBUG] Deleting Top Level File..")
	if _, err := filesClient.Delete(ctx, accountName, shareName, "", fileName); err != nil {
		t.Fatalf("Error deleting Top-Level File: %s", err)
	}
}
