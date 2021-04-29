package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type KeyVaultManagedHardwareSecurityModuleResource struct {
}

func TestAccKeyVaultManagedHardwareSecurityModule(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being able provision against one instance at a time
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"data_source": {
			"basic": testAccDataSourceKeyVaultManagedHardwareSecurityModule_basic,
		},
		"resource": {
			"basic":    testAccKeyVaultManagedHardwareSecurityModule_basic,
			"update":   testAccKeyVaultManagedHardwareSecurityModule_requiresImport,
			"complete": testAccKeyVaultManagedHardwareSecurityModule_complete,
		},
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KeyVaultManagedHardwareSecurityModuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ManagedHSMID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.KeyVault.ManagedHsmClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r KeyVaultManagedHardwareSecurityModuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                = "kvHsm%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "Standard_B1"
  tenant_id           = data.azurerm_client_config.current.tenant_id
  admin_object_ids    = [data.azurerm_client_config.current.object_id]
}
`, template, data.RandomInteger)
}

func (r KeyVaultManagedHardwareSecurityModuleResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_key_vault_managed_hardware_security_module" "import" {
  name                = azurerm_key_vault_managed_hardware_security_module.test.name
  resource_group_name = azurerm_key_vault_managed_hardware_security_module.test.resource_group_name
  location            = azurerm_key_vault_managed_hardware_security_module.test.location
  sku_name            = azurerm_key_vault_managed_hardware_security_module.test.sku_name
  tenant_id           = azurerm_key_vault_managed_hardware_security_module.test.tenant_id
  admin_object_ids    = azurerm_key_vault_managed_hardware_security_module.test.admin_object_ids
}
`, template)
}

func (r KeyVaultManagedHardwareSecurityModuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  soft_delete_retention_days = 7
  purge_protection_enabled   = true
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]

  tags = {
    Env = "Test"
  }
}
`, template, data.RandomInteger)
}

func (KeyVaultManagedHardwareSecurityModuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-KV-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
