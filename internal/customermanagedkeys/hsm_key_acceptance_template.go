package customermanagedkeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
)

// ManagedHSMKeyTempalte: Helper function to generate a template for HSM key acceptance tests
// Ensure `azurerm_client_config.current` datasource is defined before using this template.
// Verify there are no resource address conflicts in the caller of this template.
func ManagedHSMKeyTempalte(randomInteger int, randomString string, purgeProtectionEnabled bool, principalRefs []string) string {
	roleAssignes := []string{}
	for idx, principal := range principalRefs {
		randomUUID, _ := uuid.GenerateUUID()
		roleAssignes = append(roleAssignes, fmt.Sprintf(`
resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "ra%[1]d" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[3]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.encrypt-user.resource_manager_id
  principal_id       = %[2]s

  depends_on = [azurerm_key_vault_managed_hardware_security_module_key.test]
}
		`, idx, principal, randomUUID))
	}
	roleAssigneName1, _ := uuid.GenerateUUID()
	roleAssigneName2, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
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
  name                       = "kvHsm%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  sku_name                   = "Standard_B1"
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  admin_object_ids           = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled   = %[6]t
  soft_delete_retention_days = 7

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 3
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "crypto-officer" {
  name           = "515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "crypto-user" {
  name           = "21dbd100-6940-42c2-9190-5d6cb909625b"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

data "azurerm_key_vault_managed_hardware_security_module_role_definition" "encrypt-user" {
  name           = "33413926-3206-4cdd-b39a-83574fe37a17"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "client1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[4]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-officer.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "client2" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "%[5]s"
  scope              = "/keys"
  role_definition_id = data.azurerm_key_vault_managed_hardware_security_module_role_definition.crypto-user.resource_manager_id
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[1]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "RSA-HSM"
  key_size       = 2048
  key_opts       = ["unwrapKey", "wrapKey"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.client1,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.client2
  ]
}

%[3]s
`, randomString, randomInteger, strings.Join(roleAssignes, "\n\n"), roleAssigneName1, roleAssigneName2, purgeProtectionEnabled)
}
