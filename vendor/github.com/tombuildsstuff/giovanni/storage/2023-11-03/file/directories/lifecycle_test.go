package directories

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/file/shares"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

var StorageFile = Client{}

func TestDirectoriesLifeCycle(t *testing.T) {
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

	directoriesClient, err := NewWithBaseUri(fmt.Sprintf("https://%s.file.%s", accountName, *domainSuffix))
	if err := client.PrepareWithSharedKeyAuth(directoriesClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	input := shares.CreateInput{
		QuotaInGB: 1,
	}
	_, err = sharesClient.Create(ctx, shareName, input)
	if err != nil {
		t.Fatalf("Error creating fileshare: %s", err)
	}
	defer sharesClient.Delete(ctx, shareName, shares.DeleteInput{DeleteSnapshots: true})

	metaData := map[string]string{
		"hello": "world",
	}

	log.Printf("[DEBUG] Creating Top Level..")
	createInput := CreateDirectoryInput{
		MetaData: metaData,
	}
	if _, err := directoriesClient.Create(ctx, shareName, "hello", createInput); err != nil {
		t.Fatalf("Error creating Top Level Directory: %s", err)
	}

	log.Printf("[DEBUG] Creating Inner..")
	if _, err := directoriesClient.Create(ctx, shareName, "hello/there", createInput); err != nil {
		t.Fatalf("Error creating Inner Directory: %s", err)
	}

	log.Printf("[DEBUG] Retrieving share")
	innerDir, err := directoriesClient.Get(ctx, shareName, "hello/there")
	if err != nil {
		t.Fatalf("Error retrieving Inner Directory: %s", err)
	}

	if innerDir.DirectoryMetaDataEncrypted != true {
		t.Fatalf("Expected MetaData to be encrypted but got: %t", innerDir.DirectoryMetaDataEncrypted)
	}

	if len(innerDir.MetaData) != 1 {
		t.Fatalf("Expected MetaData to contain 1 item but got %d", len(innerDir.MetaData))
	}
	if innerDir.MetaData["hello"] != "world" {
		t.Fatalf("Expected MetaData `hello` to be `world`: %s", innerDir.MetaData["hello"])
	}

	log.Printf("[DEBUG] Setting MetaData")
	updatedMetaData := map[string]string{
		"panda": "pops",
	}
	if _, err := directoriesClient.SetMetaData(ctx, shareName, "hello/there", SetMetaDataInput{MetaData: updatedMetaData}); err != nil {
		t.Fatalf("Error updating MetaData: %s", err)
	}

	log.Printf("[DEBUG] Retrieving MetaData")
	retrievedMetaData, err := directoriesClient.GetMetaData(ctx, shareName, "hello/there")
	if err != nil {
		t.Fatalf("Error retrieving the updated metadata: %s", err)
	}
	if len(retrievedMetaData.MetaData) != 1 {
		t.Fatalf("Expected the updated metadata to have 1 item but got %d", len(retrievedMetaData.MetaData))
	}
	if retrievedMetaData.MetaData["panda"] != "pops" {
		t.Fatalf("Expected the metadata `panda` to be `pops` but got %q", retrievedMetaData.MetaData["panda"])
	}

	t.Logf("[DEBUG] Deleting Inner..")
	if _, err := directoriesClient.Delete(ctx, shareName, "hello/there"); err != nil {
		t.Fatalf("Error deleting Inner Directory: %s", err)
	}

	t.Logf("[DEBUG] Deleting Top Level..")
	if _, err := directoriesClient.Delete(ctx, shareName, "hello"); err != nil {
		t.Fatalf("Error deleting Top Level Directory: %s", err)
	}
}
