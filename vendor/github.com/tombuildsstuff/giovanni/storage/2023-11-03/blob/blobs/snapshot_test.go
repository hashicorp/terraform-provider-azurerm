package blobs

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/blob/containers"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestSnapshotLifecycle(t *testing.T) {
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

	t.Logf("[DEBUG] First Snapshot..")
	firstSnapshot, err := blobClient.Snapshot(ctx, containerName, fileName, SnapshotInput{})
	if err != nil {
		t.Fatalf("Error taking first snapshot: %s", err)
	}
	t.Logf("[DEBUG] First Snapshot ID: %q", firstSnapshot.SnapshotDateTime)

	t.Log("[DEBUG] Waiting 2 seconds..")
	time.Sleep(2 * time.Second)

	t.Logf("[DEBUG] Second Snapshot..")
	secondSnapshot, err := blobClient.Snapshot(ctx, containerName, fileName, SnapshotInput{
		MetaData: map[string]string{
			"hello": "world",
		},
	})
	if err != nil {
		t.Fatalf("Error taking Second snapshot: %s", err)
	}
	t.Logf("[DEBUG] Second Snapshot ID: %q", secondSnapshot.SnapshotDateTime)

	t.Logf("[DEBUG] Leasing the Blob..")
	leaseDetails, err := blobClient.AcquireLease(ctx, containerName, fileName, AcquireLeaseInput{
		// infinite
		LeaseDuration: -1,
	})
	if err != nil {
		t.Fatalf("Error leasing Blob: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %q", leaseDetails.LeaseID)

	t.Logf("[DEBUG] Third Snapshot..")
	thirdSnapshot, err := blobClient.Snapshot(ctx, containerName, fileName, SnapshotInput{
		LeaseID: &leaseDetails.LeaseID,
	})
	if err != nil {
		t.Fatalf("Error taking Third snapshot: %s", err)
	}
	t.Logf("[DEBUG] Third Snapshot ID: %q", thirdSnapshot.SnapshotDateTime)

	t.Logf("[DEBUG] Releasing Lease..")
	if _, err := blobClient.ReleaseLease(ctx, containerName, fileName, ReleaseLeaseInput{leaseDetails.LeaseID}); err != nil {
		t.Fatalf("Error releasing Lease: %s", err)
	}

	// get the properties from the blob, which should include the LastModifiedDate
	t.Logf("[DEBUG] Retrieving Properties for Blob")
	props, err := blobClient.GetProperties(ctx, containerName, fileName, GetPropertiesInput{})
	if err != nil {
		t.Fatalf("Error getting properties: %s", err)
	}

	// confirm that the If-Modified-None returns an error
	t.Logf("[DEBUG] Third Snapshot..")
	fourthSnapshot, err := blobClient.Snapshot(ctx, containerName, fileName, SnapshotInput{
		LeaseID:         &leaseDetails.LeaseID,
		IfModifiedSince: &props.LastModified,
	})
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
	if fourthSnapshot.HttpResponse.StatusCode != http.StatusPreconditionFailed {
		t.Fatalf("Expected the status code to be Precondition Failed but got: %d", fourthSnapshot.HttpResponse.StatusCode)
	}

	t.Logf("[DEBUG] Retrieving the Second Snapshot Properties..")
	getSecondSnapshotInput := GetSnapshotPropertiesInput{
		SnapshotID: secondSnapshot.SnapshotDateTime,
	}
	if _, err := blobClient.GetSnapshotProperties(ctx, containerName, fileName, getSecondSnapshotInput); err != nil {
		t.Fatalf("Error retrieving properties for the second snapshot: %s", err)
	}

	t.Logf("[DEBUG] Deleting the Second Snapshot..")
	deleteSnapshotInput := DeleteSnapshotInput{
		SnapshotDateTime: secondSnapshot.SnapshotDateTime,
	}
	if _, err := blobClient.DeleteSnapshot(ctx, containerName, fileName, deleteSnapshotInput); err != nil {
		t.Fatalf("Error deleting snapshot: %s", err)
	}

	t.Logf("[DEBUG] Re-Retrieving the Second Snapshot Properties..")
	secondSnapshotProps, err := blobClient.GetSnapshotProperties(ctx, containerName, fileName, getSecondSnapshotInput)
	if err == nil {
		t.Fatalf("Expected an error retrieving the snapshot but got none")
	}
	if secondSnapshotProps.HttpResponse.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected the status code to be %d but got %q", http.StatusNoContent, secondSnapshotProps.HttpResponse.StatusCode)
	}

	t.Logf("[DEBUG] Deleting all the snapshots..")
	if _, err := blobClient.DeleteSnapshots(ctx, containerName, fileName, DeleteSnapshotsInput{}); err != nil {
		t.Fatalf("Error deleting snapshots: %s", err)
	}

	t.Logf("[DEBUG] Deleting the Blob..")
	deleteInput := DeleteInput{
		DeleteSnapshots: false,
	}
	if _, err := blobClient.Delete(ctx, containerName, fileName, deleteInput); err != nil {
		t.Fatalf("Error deleting Blob: %s", err)
	}
}
