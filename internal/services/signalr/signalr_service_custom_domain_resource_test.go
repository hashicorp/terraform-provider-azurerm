// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SignalrServiceCustomDomainResource struct{}

func TestAccSignalrServiceCustomDomainResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_custom_domain", "test")
	r := SignalrServiceCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccSignalrServiceCustomDomainResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_signalr_service_custom_domain", "test")
	r := SignalrServiceCustomDomainResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r SignalrServiceCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := signalr.ParseCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.SignalR.SignalRClient.CustomDomainsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(resp.Model != nil), nil

}

func (r SignalrServiceCustomDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
data "azurerm_client_config" "current" {
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_signalr_service" "test" {
  name                = "acctestSignalR-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku {
    name     = "Premium_P1"
    capacity = 1
  }
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_dns_zone" "test" {
  name                = "wpstftestzone.com"
  resource_group_name = azurerm_resource_group.test.name
  depends_on = [
    azurerm_signalr_service.test,
    azurerm_signalr_service_custom_certificate.test
  ]
}

resource "azurerm_dns_cname_record" "test" {
  name                = "signalr"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_dns_zone.test.name
  ttl                 = 3600
  record              = azurerm_signalr_service.test.hostname
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkeyvault%s"
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
      "Get",
      "Import",
      "Purge",
      "Recover",
      "Update",
      "List",
    ]
    secret_permissions = [
      "Get",
      "Set",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = azurerm_signalr_service.test.identity.0.principal_id
    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Import",
      "Purge",
      "Recover",
      "Update",
      "List",
    ]

    secret_permissions = [
      "Get",
      "Set",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
  key_vault_id = azurerm_key_vault.test.id
  certificate {
    contents = filebase64("testdata/custom-domain-cert-signalr.pfx")
    password = ""
  }
}

resource "azurerm_signalr_service_custom_certificate" "test" {
  name                  = "signalr-cert-%s"
  signalr_service_id    = azurerm_signalr_service.test.id
  custom_certificate_id = azurerm_key_vault_certificate.test.id
  depends_on            = [azurerm_key_vault.test]
}

resource "azurerm_signalr_service_custom_domain" "test" {
  name                          = "signalrcustom-domain-%s"
  signalr_service_id            = azurerm_signalr_service.test.id
  domain_name                   = "signalr.${azurerm_dns_zone.test.name}"
  signalr_custom_certificate_id = azurerm_signalr_service_custom_certificate.test.id
  depends_on                    = [azurerm_dns_cname_record.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(3), data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (r SignalrServiceCustomDomainResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_signalr_service_custom_domain" "import" {
  name                          = azurerm_signalr_service_custom_domain.test.name
  signalr_service_id            = azurerm_signalr_service_custom_domain.test.signalr_service_id
  domain_name                   = azurerm_signalr_service_custom_domain.test.domain_name
  signalr_custom_certificate_id = azurerm_signalr_service_custom_domain.test.signalr_custom_certificate_id
}
`, r.basic(data))
}
