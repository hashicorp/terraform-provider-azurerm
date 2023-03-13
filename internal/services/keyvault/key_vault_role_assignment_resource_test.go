package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultRoleAssignmentResource struct{}

func (a KeyVaultRoleAssignmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.MHSMNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.KeyVault.MHSMRoleAssignClient.Get(ctx, id.VaultBaseUrl, id.Scope, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Type %s: %+v", id, err)
	}
	return utils.Bool(resp.Properties != nil), nil
}

func (a KeyVaultRoleAssignmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                     = "kvHSM-%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = "%[2]s"
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false
}

data "azurerm_key_vault_role_definition" "user" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  role_definition_id = "21dbd100-6940-42c2-9190-5d6cb909625b"
  scope              = "/"
  principal_id       = data.azurerm_client_config.current.object_id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (a KeyVaultRoleAssignmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_key_vault_role_assignment" "test" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  scope              = "${data.azurerm_key_vault_role_definition.user.scope}"
  role_definition_id = "${data.azurerm_key_vault_role_definition.user.id}"
  principal_id       = "${data.azurerm_client_config.current.object_id}"
}
`, a.template(data))
}

func (a KeyVaultRoleAssignmentResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`




%s

resource "azurerm_key_vault_role_assignment" "test" {

  vault_base_url     = "example"
  name               = "acctest-%[2]d"
  scope              = "example"
  role_definition_id = "example"
  principal_id       = "example"
}
`, a.template(data), data.RandomInteger)
}

// TODO we cannot test it right now, because we cannot activate mnagedhsm keyvault from terraform
func TestAccKeyVaultRoleAssignment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, keyvault.KeyVaultRoleAssignmentResource{}.ResourceType(), "test")
	r := KeyVaultRoleAssignmentResource{}
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
