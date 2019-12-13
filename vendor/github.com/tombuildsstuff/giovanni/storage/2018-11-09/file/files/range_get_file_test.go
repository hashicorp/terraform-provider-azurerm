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

func TestGetSmallFile(t *testing.T) {
	// the purpose of this test is to verify that the small, single-chunked file gets downloaded correctly
	testGetFile(t, "small-file.png", "image/png")
}

func TestGetLargeFile(t *testing.T) {
	// the purpose of this test is to verify that the large, multi-chunked file gets downloaded correctly
	testGetFile(t, "blank-large-file.dmg", "application/x-apple-diskimage")
}

func testGetFile(t *testing.T, fileName string, contentType string) {
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

	t.Logf("[DEBUG] Downloading file..")
	_, downloadedBytes, err := filesClient.GetFile(ctx, accountName, shareName, "", fileName, 4)
	if err != nil {
		t.Fatalf("Error downloading file: %s", err)
	}

	t.Logf("[DEBUG] Asserting the files are the same size..")
	expectedBytes := make([]byte, info.Size())
	file.Read(expectedBytes)
	if len(expectedBytes) != len(downloadedBytes) {
		t.Fatalf("Expected %d bytes but got %d", len(expectedBytes), len(downloadedBytes))
	}

	t.Logf("[DEBUG] Asserting the files are the same content-wise..")
	// overkill, but it's this or shasum-ing
	for i := int64(0); i < info.Size(); i++ {
		if expectedBytes[i] != downloadedBytes[i] {
			t.Fatalf("Expected byte %d to be %q but got %q", i, expectedBytes[i], downloadedBytes[i])
		}
	}

	t.Logf("[DEBUG] Deleting Top Level File..")
	if _, err := filesClient.Delete(ctx, accountName, shareName, "", fileName); err != nil {
		t.Fatalf("Error deleting Top-Level File: %s", err)
	}

}
