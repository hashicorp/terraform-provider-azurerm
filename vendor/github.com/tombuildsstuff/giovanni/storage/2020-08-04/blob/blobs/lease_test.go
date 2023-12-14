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

func TestLeaseLifecycle(t *testing.T) {
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
	defer blobClient.Delete(ctx, containerName, fileName, DeleteInput{})

	// Test begins here
	t.Logf("[DEBUG] Acquiring Lease..")
	leaseInput := AcquireLeaseInput{
		LeaseDuration: -1,
	}
	leaseInfo, err := blobClient.AcquireLease(ctx, containerName, fileName, leaseInput)
	if err != nil {
		t.Fatalf("Error acquiring lease: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %q", leaseInfo.LeaseID)

	t.Logf("[DEBUG] Changing Lease..")
	changeLeaseInput := ChangeLeaseInput{
		ExistingLeaseID: leaseInfo.LeaseID,
		ProposedLeaseID: "31f5bb01-cdd9-4166-bcdc-95186076bde0",
	}
	changeLeaseResult, err := blobClient.ChangeLease(ctx, containerName, fileName, changeLeaseInput)
	if err != nil {
		t.Fatalf("Error changing lease: %s", err)
	}
	t.Logf("[DEBUG] New Lease ID: %q", changeLeaseResult.LeaseID)

	t.Logf("[DEBUG] Releasing Lease..")
	if _, err := blobClient.ReleaseLease(ctx, containerName, fileName, ReleaseLeaseInput{LeaseID: changeLeaseResult.LeaseID}); err != nil {
		t.Fatalf("Error releasing lease: %s", err)
	}

	t.Logf("[DEBUG] Acquiring a new lease..")
	leaseInput = AcquireLeaseInput{
		LeaseDuration: 30,
	}
	leaseInfo, err = blobClient.AcquireLease(ctx, containerName, fileName, leaseInput)
	if err != nil {
		t.Fatalf("Error acquiring lease: %s", err)
	}
	t.Logf("[DEBUG] Lease ID: %q", leaseInfo.LeaseID)

	t.Logf("[DEBUG] Renewing lease..")
	if _, err := blobClient.RenewLease(ctx, containerName, fileName, RenewLeaseInput{LeaseID: leaseInfo.LeaseID}); err != nil {
		t.Fatalf("Error renewing lease: %s", err)
	}

	t.Logf("[DEBUG] Breaking lease..")
	breakLeaseInput := BreakLeaseInput{
		LeaseID: leaseInfo.LeaseID,
	}
	if _, err := blobClient.BreakLease(ctx, containerName, fileName, breakLeaseInput); err != nil {
		t.Fatalf("Error breaking lease: %s", err)
	}
}
