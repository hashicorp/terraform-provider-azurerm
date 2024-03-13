// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type KeyVaultManagedHSMKeyResource struct{}

// real test nested in TestAccKeyVaultManagedHardwareSecurityModule, only provide Exists logic here
func (k KeyVaultManagedHSMKeyResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ParseNestedItemID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.ManagedHSMs.DataPlaneManagedHSMClient.GetKey(ctx, id.HSMBaseUrl, id.Name, id.Version)
	if err != nil {
		return nil, fmt.Errorf("retrieving Managed HSM Key %s: %+v", id, err)
	}
	return utils.Bool(resp.Attributes != nil), nil
}

func (k KeyVaultManagedHSMKeyResource) template(data acceptance.TestData, purge bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_deleted_keys_on_destroy = "%[3]t"
      recover_soft_deleted_key_vaults    = true
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-KV-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acc%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = [
      "Create",
      "Delete",
      "DeleteIssuers",
      "Get",
      "Purge",
      "Update"
    ]
  }
}

resource "azurerm_key_vault_certificate" "cert" {
  count        = 3
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
  name                                      = "kvHsm%[1]d"
  resource_group_name                       = azurerm_resource_group.test.name
  location                                  = azurerm_resource_group.test.location
  sku_name                                  = "Standard_B1"
  tenant_id                                 = data.azurerm_client_config.current.tenant_id
  admin_object_ids                          = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled                  = false
  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 2
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "role_user" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "role_user" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name               = "706c03c7-69ad-33e5-2796-b3380d3a6e1a"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.role_user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "role_officer" {
  vault_base_url = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
}

// need crypto officer role to purge keys
resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "role_officer" {
  vault_base_url     = azurerm_key_vault_managed_hardware_security_module.test.hsm_uri
  name               = "d1a3242a-d521-11ee-9880-00155d316070"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.role_officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}
	`, data.RandomInteger, data.Locations.Primary, purge)
}

func (k KeyVaultManagedHSMKeyResource) basic(data acceptance.TestData, purge bool) string {

	return fmt.Sprintf(`

%s

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "key-%s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "EC-HSM"

  key_options = [
    "sign",
    "verify",
  ]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_user,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_officer
  ]
}
`, k.template(data, purge), data.RandomString)
}

func (k KeyVaultManagedHSMKeyResource) update(data acceptance.TestData, purge bool) string {

	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "key-%s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "EC-HSM"

  key_options = [
    "sign",
    "verify",
  ]

  expiration_date = "2037-12-31T00:00:00Z"
  tags = {
    Env = "test"
  }

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_user,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_officer
  ]
}
`, k.template(data, purge), data.RandomString)
}

func (k KeyVaultManagedHSMKeyResource) rotationPolicy(data acceptance.TestData, purge bool) string {

	return fmt.Sprintf(`


%s

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "key-%s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048

  key_options = [
    "sign",
    "verify",
    "encrypt",
    "decrypt",
  ]

  rotation_policy {
    automatic {
      duration_before_expiry = "P30D"
    }

    expire_after_duration = "P60D"
  }
  tags = {
    Env = "test"
  }

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_user,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.role_officer
  ]
}
`, k.template(data, purge), data.RandomString)
}
