package azurerm

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2017-10-01/storage"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMContainerRegistryMigrateState(t *testing.T) {
	config := testGetAzureConfig(t)
	if config == nil {
		t.SkipNow()
		return
	}
	client, err := getArmClient(config)
	if err != nil {
		t.Fatal(fmt.Errorf("Error building ARM Client: %+v", err))
		return
	}

	client.StopContext = testAccProvider.StopContext()

	rs := acctest.RandString(4)
	resourceGroupName := fmt.Sprintf("acctestRG%s", rs)
	storageAccountName := fmt.Sprintf("acctestsa%s", rs)
	location := azureRMNormalizeLocation(testLocation())
	ctx := client.StopContext

	err = createResourceGroup(ctx, client, resourceGroupName, location)
	if err != nil {
		t.Fatal(err)
		return
	}

	storageAccount, err := createStorageAccount(client, resourceGroupName, storageAccountName, location)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer destroyStorageAccountAndResourceGroup(client, resourceGroupName, storageAccountName)

	cases := map[string]struct {
		StateVersion int
		ID           string
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"v0_1_without_value": {
			StateVersion: 0,
			ID:           "some_id",
			Attributes:   map[string]string{},
			Expected: map[string]string{
				"sku": "Basic",
			},
		},
		"v1_2_with_value": {
			StateVersion: 1,
			ID:           "some_id",
			Attributes: map[string]string{
				// TODO: storage_account also needs to become a List
				"sku":                    "Basic",
				"storage_account.#":      "1",
				"storage_account.0.name": storageAccountName,
			},
			Expected: map[string]string{
				"sku":                    "Classic",
				"storage_account.#":      "1",
				"storage_account.0.name": storageAccountName,
				"storage_account_id":     *storageAccount.ID,
			},
			Meta: client,
		},
	}

	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         tc.ID,
			Attributes: tc.Attributes,
		}
		is, err := resourceAzureRMContainerRegistryMigrateState(tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %q, err: %#v", tn, err)
		}

		if !reflect.DeepEqual(tc.Expected, is.Attributes) {
			t.Fatalf("Bad Container Registry Migrate\n\n. Got: %+v\n\n expected: %+v", is.Attributes, tc.Expected)
		}
	}
}

func createResourceGroup(ctx context.Context, client *ArmClient, resourceGroupName string, location string) error {
	group := resources.Group{
		Location: &location,
	}
	_, err := client.resourceGroupsClient.CreateOrUpdate(ctx, resourceGroupName, group)
	if err != nil {
		return fmt.Errorf("Error creating Resource Group %q: %+v", resourceGroupName, err)
	}
	return nil
}

func createStorageAccount(client *ArmClient, resourceGroupName, storageAccountName, location string) (*storage.Account, error) {
	storageClient := client.storageServiceClient
	createParams := storage.AccountCreateParameters{
		Location: &location,
		Kind:     storage.Storage,
		Sku: &storage.Sku{
			Name: storage.StandardLRS,
			Tier: storage.Standard,
		},
	}
	ctx := client.StopContext
	future, err := storageClient.Create(ctx, resourceGroupName, storageAccountName, createParams)
	if err != nil {
		return nil, fmt.Errorf("Error creating Storage Account %q: %+v", resourceGroupName, err)
	}

	err = future.WaitForCompletion(ctx, storageClient.Client)
	if err != nil {
		return nil, fmt.Errorf("Error waiting for creation of Storage Account %q: %+v", resourceGroupName, err)
	}

	account, err := storageClient.GetProperties(ctx, resourceGroupName, storageAccountName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Storage Account %q: %+v", resourceGroupName, err)
	}

	return &account, nil
}

func destroyStorageAccountAndResourceGroup(client *ArmClient, resourceGroupName, storageAccountName string) {
	ctx := client.StopContext
	client.storageServiceClient.Delete(ctx, resourceGroupName, storageAccountName)
	client.resourceGroupsClient.Delete(ctx, resourceGroupName)
}
