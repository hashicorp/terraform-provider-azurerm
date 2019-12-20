package azurerm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func resourceAzureRMContainerRegistryMigrateState(
	v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found AzureRM Container Registry State v0; migrating to v1")
		return migrateAzureRMContainerRegistryStateV0toV1(is)
	case 1:
		log.Println("[INFO] Found AzureRM Container Registry State v1; migrating to v2")
		return migrateAzureRMContainerRegistryStateV1toV2(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateAzureRMContainerRegistryStateV0toV1(is *terraform.InstanceState) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Container Registry Attributes before Migration: %#v", is.Attributes)

	is.Attributes["sku"] = "Basic"

	log.Printf("[DEBUG] ARM Container Registry Attributes after State Migration: %#v", is.Attributes)

	return is, nil
}

func migrateAzureRMContainerRegistryStateV1toV2(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] ARM Container Registry Attributes before Migration: %#v", is.Attributes)

	// Basic's been renamed Classic to allow for "ManagedBasic" ¯\_(ツ)_/¯
	is.Attributes["sku"] = "Classic"

	err := updateV1ToV2StorageAccountName(is, meta)

	// we have to look this up, since we don't have the resource group name

	log.Printf("[DEBUG] ARM Container Registry Attributes after State Migration: %#v", is.Attributes)

	return is, err
}

func updateV1ToV2StorageAccountName(is *terraform.InstanceState, meta interface{}) error {
	reader := &schema.MapFieldReader{
		Schema: resourceArmContainerRegistry().Schema,
		Map:    schema.BasicMapReader(is.Attributes),
	}

	result, err := reader.ReadField([]string{"storage_account"})
	if err != nil {
		return err
	}

	if result.Value == nil {
		return nil
	}

	inputAccounts := result.Value.([]interface{})
	inputAccount := inputAccounts[0]
	if inputAccount == nil {
		return nil
	}

	account := inputAccount.(map[string]interface{})
	name := account["name"].(string)
	storageAccountId, err := findAzureStorageAccountIdFromName(name, meta)
	if err != nil {
		return err
	}

	is.Attributes["storage_account_id"] = storageAccountId
	return nil
}

func findAzureStorageAccountIdFromName(name string, meta interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, time.Minute*5)
	defer cancel()

	client := meta.(*clients.Client).Storage.AccountsClient
	accounts, err := client.List(ctx)
	if err != nil {
		return "", err
	}

	for _, account := range accounts.Values() {
		if strings.EqualFold(*account.Name, name) {
			return *account.ID, nil
		}
	}

	return "", fmt.Errorf("Unable to find ID for Storage Account %q", name)
}
