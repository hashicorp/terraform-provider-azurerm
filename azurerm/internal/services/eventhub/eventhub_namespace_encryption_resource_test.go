package eventhub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type EventHubNamespaceEncryptionResource struct {
}

func TestAccEventHubNamespaceEncryption_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventhub_namespace_encryption", "test")
	r := EventHubNamespaceEncryptionResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (EventHubNamespaceEncryptionResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.NamespaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Eventhub.NamespacesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", id.String(), err)
	}

	return utils.Bool(resp.EHNamespaceProperties != nil), nil
}

func (EventHubNamespaceEncryptionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}


data "azurerm_resource_group" "test" {
  name = "acctestRG-210303215248990504"
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv-%s"
  location                 = data.azurerm_resource_group.test.location
  resource_group_name      = data.azurerm_resource_group.test.name
  tenant_id                = azurerm_eventhub_namespace.test.identity[0].tenant_id
  sku_name                 = "premium"
  purge_protection_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }

  access_policy {
    tenant_id = azurerm_eventhub_namespace.test.identity[0].tenant_id
    object_id = azurerm_eventhub_namespace.test.identity[0].principal_id

    key_permissions = ["get", "unwrapkey", "wrapkey"]
  }

  tags = {
    environment = "Production"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey"
  ]
}

data "azurerm_eventhub_cluster" "test" {
  name                = "acctesteventhubcluster-210303215248990504"
  resource_group_name = data.azurerm_resource_group.test.name
}

resource "azurerm_eventhub_namespace" "test" {
  name                 = "acctesteventhubnamespace-%d"
  location             = data.azurerm_resource_group.test.location
  resource_group_name  = data.azurerm_resource_group.test.name
  sku                  = "Standard"
  capacity             = "2"
  dedicated_cluster_id = data.azurerm_eventhub_cluster.test.id
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_eventhub_namespace_encryption" "test" {
  namespace_id  = azurerm_eventhub_namespace.test.id
  key_vault_uri = azurerm_key_vault.test.vault_uri
  key_name      = azurerm_key_vault_key.test.name
}
`, data.RandomString, data.RandomString, data.RandomInteger)
}
