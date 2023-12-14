package paths

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2023-01-01/storageaccounts"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/tombuildsstuff/giovanni/storage/2020-08-04/datalakestore/filesystems"
	"github.com/tombuildsstuff/giovanni/storage/internal/testhelpers"
)

func TestLifecycle(t *testing.T) {
	const defaultACLString = "user::rwx,group::r-x,other::---"

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	client, err := testhelpers.Build(ctx, t)
	if err != nil {
		t.Fatal(err)
	}

	resourceGroup := fmt.Sprintf("acctestrg-%d", testhelpers.RandomInt())
	accountName := fmt.Sprintf("acctestsa%s", testhelpers.RandomString())
	fileSystemName := fmt.Sprintf("acctestfs-%s", testhelpers.RandomString())
	path := "test"

	testData, err := client.BuildTestResourcesWithHns(ctx, resourceGroup, accountName, storageaccounts.KindBlobStorage)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DestroyTestResources(ctx, resourceGroup, accountName)
	domainSuffix, ok := client.Environment.Storage.DomainSuffix()
	if !ok {
		t.Fatalf("storage didn't return a domain suffix for this environment")
	}

	baseUri := fmt.Sprintf("https://%s.%s.%s", accountName, "dfs", *domainSuffix)

	fileSystemsClient, err := filesystems.NewWithBaseUri(baseUri)
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}
	if err := client.PrepareWithSharedKeyAuth(fileSystemsClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	pathsClient, err := NewWithBaseUri(baseUri)
	if err != nil {
		t.Fatalf("building client for environment: %+v", err)
	}
	if err := client.PrepareWithSharedKeyAuth(pathsClient.Client, testData, auth.SharedKey); err != nil {
		t.Fatalf("adding authorizer to client: %+v", err)
	}

	t.Logf("[DEBUG] Creating an empty File System..")
	fileSystemInput := filesystems.CreateInput{}
	if _, err = fileSystemsClient.Create(ctx, fileSystemName, fileSystemInput); err != nil {
		t.Fatal(fmt.Errorf("error creating: %s", err))
	}

	t.Logf("[DEBUG] Creating folder 'test' ..")
	input := CreateInput{
		Resource: PathResourceDirectory,
	}
	if _, err = pathsClient.Create(ctx, fileSystemName, path, input); err != nil {
		t.Fatal(fmt.Errorf("error creating: %s", err))
	}

	t.Logf("[DEBUG] Getting properties for folder 'test' ..")
	props, err := pathsClient.GetProperties(ctx, fileSystemName, path, GetPropertiesInput{action: GetPropertiesActionGetAccessControl})
	if err != nil {
		t.Fatal(fmt.Errorf("error getting properties: %s", err))
	}
	t.Logf("[DEBUG] Props.Owner: %q", props.Owner)
	t.Logf("[DEBUG] Props.Group: %q", props.Group)
	t.Logf("[DEBUG] Props.ACL: %q", props.ACL)
	t.Logf("[DEBUG] Props.ETag: %q", props.ETag)
	t.Logf("[DEBUG] Props.LastModified: %q", props.LastModified)
	if props.ACL != defaultACLString {
		t.Fatal(fmt.Errorf("Expected Default ACL %q, got %q", defaultACLString, props.ACL))
	}

	newACL := "user::rwx,group::r-x,other::r-x,default:user::rwx,default:group::r-x,default:other::---"
	accessControlInput := SetAccessControlInput{
		ACL: &newACL,
	}
	t.Logf("[DEBUG] Setting Access Control for folder 'test' ..")
	if _, err = pathsClient.SetAccessControl(ctx, fileSystemName, path, accessControlInput); err != nil {
		t.Fatal(fmt.Errorf("error setting Access Control %s", err))
	}

	t.Logf("[DEBUG] Getting properties for folder 'test' (2) ..")
	props, err = pathsClient.GetProperties(ctx, fileSystemName, path, GetPropertiesInput{action: GetPropertiesActionGetAccessControl})
	if err != nil {
		t.Fatal(fmt.Errorf("error getting properties (2): %s", err))
	}
	if props.ACL != newACL {
		t.Fatal(fmt.Errorf("expected new ACL %q, got %q", newACL, props.ACL))
	}

	t.Logf("[DEBUG] Deleting path 'test' ..")
	if _, err = pathsClient.Delete(ctx, fileSystemName, path); err != nil {
		t.Fatal(fmt.Errorf("error deleting path: %s", err))
	}

	t.Logf("[DEBUG] Getting properties for folder 'test' (3) ..")
	props, err = pathsClient.GetProperties(ctx, fileSystemName, path, GetPropertiesInput{action: GetPropertiesActionGetAccessControl})
	if err == nil {
		t.Fatal(fmt.Errorf("didn't get error getting properties after deleting path (3)"))
	}

	t.Logf("[DEBUG] Deleting File System..")
	if _, err := fileSystemsClient.Delete(ctx, fileSystemName); err != nil {
		t.Fatalf("Error deleting filesystem: %s", err)
	}
}
