package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type KeyVaultRoleDefinitionResource struct{}

func (a KeyVaultRoleDefinitionResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	baseURL := state.Attributes["vault_base_url"]
	id, err := parse.MHSMNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.KeyVault.MHSMRoleClient.Get(ctx, baseURL, "/", id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.RoleDefinitionProperties != nil), nil
}

func (a KeyVaultRoleDefinitionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                     = "kvHsm%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a KeyVaultRoleDefinitionResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_role_definition" "test" {
  name              = "acctest-%[2]d"
  scope             = "/"
  vault_base_url    = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  description       = "desc foo"
  assignable_scopes = ["/keys"]
  permission {
    actions = ["Microsoft.KeyVault/managedHsm/keys/read/action"]
  }
}
`, a.template(data), data.RandomInteger, data.Locations.Primary)
}

func (a KeyVaultRoleDefinitionResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_role_definition" "test" {
  name              = "acctest-%[2]d"
  scope             = "/"
  vault_base_url    = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  description       = "desc foo"
  assignable_scopes = ["/keys"]
  permission {
    actions     = ["Microsoft.KeyVault/managedHsm/keys/read/action"]
    not_actions = ["Microsoft.KeyVault/managedHsm/keys/delete/action"]
  }
}
`, a.template(data), data.RandomInteger)
}

// We cannot run this Test for now, because we cannot activate mhsm keyvault by terraform now.
func TestAccKeyVaultRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, keyvault.KeyVaultRoleDefinitionResource{}.ResourceType(), "test")
	r := KeyVaultRoleDefinitionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_global").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("is_global").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}
