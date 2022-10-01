package keyvault_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultKeyRotationPolicyResource struct{}

func TestAccKeyVaultKeyRotationPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key_rotation_policy", "test")
	r := KeyVaultKeyRotationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").MatchesRegex(regexp.MustCompile(`^/subscriptions/[\w-]+/resourceGroups/[\w-]+/providers/Microsoft.KeyVault/vaults/[\w-]+/keys/[\w-]+/rotationpolicy`)),
			),
		},
		data.ImportStep("key_vault_id"),
	})
}

func TestAccKeyVaultKeyRotationPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key_rotation_policy", "test")
	r := KeyVaultKeyRotationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_key_vault_key_rotation_policy"),
		},
	})
}

func TestAccKeyVaultKeyRotationPolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_key_rotation_policy", "test")
	r := KeyVaultKeyRotationPolicyResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basicUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r KeyVaultKeyRotationPolicyResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.KeyVault.ManagementClient
	keyVaultsClient := clients.KeyVault

	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}

	keyVaultIdRaw, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, clients.Resource, id.KeyVaultBaseUrl)
	if err != nil || keyVaultIdRaw == nil {
		return nil, fmt.Errorf("retrieving the Resource ID the Key Vault at URL %q: %s", id.KeyVaultBaseUrl, err)
	}
	keyVaultId, err := parse.VaultID(*keyVaultIdRaw)
	if err != nil {
		return nil, err
	}

	ok, err := keyVaultsClient.Exists(ctx, *keyVaultId)
	if err != nil || !ok {
		return nil, fmt.Errorf("checking if key vault %q for Certificate %q in Vault at url %q exists: %v", *keyVaultId, id.Name, id.KeyVaultBaseUrl, err)
	}

	resp, err := client.GetKeyRotationPolicy(ctx, id.KeyVaultBaseUrl, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Key Vault Key Rotation Policy %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r KeyVaultKeyRotationPolicyResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_key_vault_key_rotation_policy" "test" {
  key_name     = azurerm_key_vault_key.test.name
  key_vault_id = azurerm_key_vault.test.id
  auto_rotation {
    time_after_create = "P30D"
  }

  expiry_time       = "P60D"
  notification_time = "P7D"
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultKeyRotationPolicyResource) basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_key" "test" {
  name         = "key-%s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "EC"
  key_size     = 2048

  key_opts = [
    "sign",
    "verify",
  ]
}

resource "azurerm_key_vault_key_rotation_policy" "test" {
  key_name     = azurerm_key_vault_key.test.name
  key_vault_id = azurerm_key_vault.test.id
  auto_rotation {
    time_after_create = "P31D"
  }

  expiry_time       = "P61D"
  notification_time = "P8D"
}
`, r.template(data), data.RandomString)
}

func (r KeyVaultKeyRotationPolicyResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_key_rotation_policy" "import" {
  key_name     = azurerm_key_vault_key_rotation_policy.test.key_name
  key_vault_id = azurerm_key_vault.test.id

  auto_rotation {
    time_after_create = "P30D"
  }

  expiry_time       = "P60D"
  notification_time = "P2D"
}
`, r.basic(data))
}

func (KeyVaultKeyRotationPolicyResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkv-%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Recover",
      "Update",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]

    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
  }

  tags = {
    environment = "Production"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
