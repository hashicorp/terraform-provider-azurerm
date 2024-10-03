// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package managedhsm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MHSMKeyRotationPolicyTestResource struct{}

func testAccMHSMKeyRotationPolicy_all(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_key_vault_managed_hardware_security_module_key_rotation_policy", "test")
	r := MHSMKeyRotationPolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MHSMKeyRotationPolicyTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	domainSuffix, ok := clients.Account.Environment.ManagedHSM.DomainSuffix()
	if !ok {
		return nil, fmt.Errorf("could not determine Managed HSM domain suffix for environment %q", clients.Account.Environment.Name)
	}
	id, err := parse.ManagedHSMDataPlaneVersionlessKeyID(state.ID, domainSuffix)
	if err != nil {
		return nil, err
	}

	subscriptionId := commonids.NewSubscriptionID(clients.Account.SubscriptionId)
	resourceManagerId, err := clients.ManagedHSMs.ManagedHSMIDFromBaseUrl(ctx, subscriptionId, id.BaseUri(), domainSuffix)
	if err != nil {
		return nil, fmt.Errorf("determining Resource Manager ID for %q: %+v", id, err)
	}
	if resourceManagerId == nil {
		return nil, fmt.Errorf("unable to determine the Resource Manager ID for %s", id)
	}

	resp, err := clients.ManagedHSMs.DataPlaneKeysClient.GetKeyRotationPolicy(ctx, id.BaseUri(), id.KeyName)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %+v", *id, err)
	}

	return utils.Bool(resp.Attributes != nil), nil
}

func (r MHSMKeyRotationPolicyTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_hardware_security_module_key_rotation_policy" "test" {
  managed_hsm_key_id = azurerm_key_vault_managed_hardware_security_module_key.test.id
  time_before_expiry = "P30D"
  expire_after       = "P60D"
}
`, r.template(data))
}

func (r MHSMKeyRotationPolicyTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_key_vault_managed_hardware_security_module_key_rotation_policy" "test" {
  managed_hsm_key_id  = azurerm_key_vault_managed_hardware_security_module_key.test.id
  expire_after        = "P60D"
  time_after_creation = "P28D"
}
`, r.template(data))
}

func (r MHSMKeyRotationPolicyTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-KV-%[1]s"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                       = "acc%[3]d"
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
  name                     = "kvHsm%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  sku_name                 = "Standard_B1"
  tenant_id                = data.azurerm_client_config.current.tenant_id
  admin_object_ids         = [data.azurerm_client_config.current.object_id]
  purge_protection_enabled = false

  security_domain_key_vault_certificate_ids = [for cert in azurerm_key_vault_certificate.cert : cert.id]
  security_domain_quorum                    = 3
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "1e243909-064c-6ac3-84e9-1c8bf8d6ad22"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/21dbd100-6940-42c2-9190-5d6cb909625b"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_role_assignment" "test1" {
  managed_hsm_id     = azurerm_key_vault_managed_hardware_security_module.test.id
  name               = "1e243909-064c-6ac3-84e9-1c8bf8d6ad23"
  scope              = "/keys"
  role_definition_id = "/Microsoft.KeyVault/providers/Microsoft.Authorization/roleDefinitions/515eb02d-2335-4d2d-92f2-b1cbdf9c3778"
  principal_id       = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_managed_hardware_security_module_key" "test" {
  name           = "acctestHSMK-%[1]s"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.test.id
  key_type       = "EC-HSM"
  curve          = "P-521"
  key_opts       = ["sign"]

  depends_on = [
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test,
    azurerm_key_vault_managed_hardware_security_module_role_assignment.test1
  ]
}
`, data.RandomString, data.Locations.Primary, data.RandomInteger)
}
