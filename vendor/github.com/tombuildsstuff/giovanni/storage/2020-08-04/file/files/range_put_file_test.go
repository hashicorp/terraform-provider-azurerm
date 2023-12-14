package files

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestPutSmallFile(t *testing.T) {
	// the purpose of this test is to ensure that a small file (< 4MB) is a single chunk
	testPutFile(t, "small-file.png", "image/png")
}

func TestPutLargeFile(t *testing.T) {
	// the purpose of this test is to ensure that large files (> 4MB) are chunked
	testPutFile(t, "blank-large-file.dmg", "application/x-apple-diskimage")
}

func TestPutVerySmallFile(t *testing.T) {
	// the purpose of this test is to ensure that a very small file (< 4KB) is a single chunk
	testPutFile(t, "very-small.json", "application/json")
}

func testPutFile(t *testing.T, fileName string, contentType string) {
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
	if _, err := filesClient.Create(ctx, shareName, "", fileName, createFileInput); err != nil {
		t.Fatalf("Error creating Top-Level File: %s", err)
	}

	t.Logf("[DEBUG] Uploading File..")
	if err := filesClient.PutFile(ctx, shareName, "", fileName, file, 4); err != nil {
		t.Fatalf("Error uploading File: %s", err)
	}

	t.Logf("[DEBUG] Deleting Top Level File..")
	if _, err := filesClient.Delete(ctx, shareName, "", fileName); err != nil {
		t.Fatalf("Error deleting Top-Level File: %s", err)
	}
}
