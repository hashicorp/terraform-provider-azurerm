// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package keyvault_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/keyvault/2023-02-01/managedhsms"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultManagedHardwareSecurityModuleResource struct{}

func TestAccKeyVaultManagedHardwareSecurityModule(t *testing.T) {
	// NOTE: this is a combined test rather than separate split out tests due to
	// Azure only being able provision against one instance at a time
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"resource": {
			"data_source": testAccDataSourceKeyVaultManagedHardwareSecurityModule_basic,
			"basic":       testAccKeyVaultManagedHardwareSecurityModule_basic,
			"update":      testAccKeyVaultManagedHardwareSecurityModule_requiresImport,
			"complete":    testAccKeyVaultManagedHardwareSecurityModule_complete,
			"download":    testAccKeyVaultManagedHardwareSecurityModule_download,
		},
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_download(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.download(data, 3),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("security_domain_quorum", "security_domain_key_vault_certificate_ids", "security_domain_encrypted_data"),
		{
			Config: r.download(data, 4),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("security_domain_quorum", "security_domain_key_vault_certificate_ids", "security_domain_encrypted_data"),
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func testAccKeyVaultManagedHardwareSecurityModule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module", "test")
	r := KeyVaultManagedHardwareSecurityModuleResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (KeyVaultManagedHardwareSecurityModuleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedhsms.ParseManagedHSMID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.KeyVault.ManagedHsmClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r KeyVaultManagedHardwareSecurityModuleResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                     = "kvHsm%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false
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

func (r KeyVaultManagedHardwareSecurityModuleResource) download(data acceptance.TestData, certCount int) string {
	template := r.template(data)
	activateConfig := ""
	if certCount > 0 {
		activateConfig = `
  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum 				    = 2
`
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
%[1]s
resource "azurerm_key_vault" "test" {
  name                       = "acc%[2]d"
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
    ]
    secret_permissions = [
      "Delete",
      "Get",
      "Set",
    ]
    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
  tags = {
    environment = "Production"
  }
}
resource "azurerm_key_vault_certificate" "cert" {
  count        = %[3]d
  name         = "acchsmcert${count.index}"
  key_vault_id = azurerm_key_vault.test.id
  certificate_policy {
    issuer_parameters {
      name = "Self"
    }
    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }
    lifetime_action {
      action {
        action_type = "AutoRenew"
      }
      trigger {
        days_before_expiry = 30
      }
    }
    secret_properties {
      content_type = "application/x-pkcs12"
    }
    x509_certificate_properties {
      extended_key_usage = []
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]
      subject            = "CN=hello-world"
      validity_in_months = 12
    }
  }
}
resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                     = "kvHsm%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false
  %[4]s
}
`, template, data.RandomInteger, certCount, activateConfig)
}
func (r KeyVaultManagedHardwareSecurityModuleResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test_a" {
  name                 = "acctestsubneta%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.KeyVault"]
}

resource "azurerm_key_vault_managed_hardware_security_module" "test" {
  name                       = "kvHsm%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]

  network_acls {
    default_action = "Deny"
    bypass         = "None"
  }

  public_network_access_enabled = true

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
