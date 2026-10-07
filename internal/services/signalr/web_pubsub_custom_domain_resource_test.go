// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package signalr_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2024-03-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type WebPubsubCustomDomainResource struct{}

func TestAccWebPubsubCustomDomainResource_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_custom_domain", "test")
	r := WebPubsubCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccWebPubsubCustomDomainResource_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE and/or ARM_TEST_DATA_RESOURCE_GROUP are not specified")
		return
	}
	data := acceptance.BuildTestData(t, "azurerm_web_pubsub_custom_domain", "test")
	r := WebPubsubCustomDomainResource{}

	data.ResourceTestIgnoreRecreate(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r WebPubsubCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := webpubsub.ParseCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.SignalR.WebPubSubClient.WebPubSub.CustomDomainsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r WebPubsubCustomDomainResource) basic(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
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

resource "azurerm_web_pubsub" "test" {
  name                = "acctestwebPubsub-%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Premium_P1"

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_dns_zone" "test" {
  name                = "%s"
  resource_group_name = "%s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "wps-%s"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 3600
  record              = azurerm_web_pubsub.test.hostname
  depends_on = [
    azurerm_web_pubsub_custom_certificate.test
  ]
}

resource "azurerm_key_vault" "test" {
  name                       = "acctestkeyvault%s"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  rbac_authorization_enabled = false
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
    object_id = azurerm_web_pubsub.test.identity[0].principal_id
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
    contents = filebase64("testdata/custom-domain-cert-wps-acctest.pfx")
    password = ""
  }
}

resource "azurerm_web_pubsub_custom_certificate" "test" {
  name                  = "webPubsubcert-%s"
  web_pubsub_id         = azurerm_web_pubsub.test.id
  custom_certificate_id = azurerm_key_vault_certificate.test.id
  depends_on            = [azurerm_key_vault.test]
}

resource "azurerm_web_pubsub_custom_domain" "test" {
  name                             = "webpubsubcustom-domain-%s"
  web_pubsub_id                    = azurerm_web_pubsub.test.id
  domain_name                      = "wps-%s.${data.azurerm_dns_zone.test.name}"
  web_pubsub_custom_certificate_id = azurerm_web_pubsub_custom_certificate.test.id
  depends_on                       = [azurerm_dns_cname_record.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(3), dnsZone, dataResourceGroup, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString, data.RandomString)
}

func (r WebPubsubCustomDomainResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_web_pubsub_custom_domain" "import" {
  name                             = azurerm_web_pubsub_custom_domain.test.name
  web_pubsub_id                    = azurerm_web_pubsub_custom_domain.test.web_pubsub_id
  domain_name                      = azurerm_web_pubsub_custom_domain.test.domain_name
  web_pubsub_custom_certificate_id = azurerm_web_pubsub_custom_domain.test.web_pubsub_custom_certificate_id
}
`, r.basic(data))
}
