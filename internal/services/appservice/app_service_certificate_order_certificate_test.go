package appservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CertificateOrderCertificateResource struct{}

func TestAccAppServiceCertificateOrderCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order_certificate", "test")
	r := CertificateOrderCertificateResource{}

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

func TestAccAppServiceCertificateOrderCertificate_updateKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order_certificate", "test")
	r := CertificateOrderCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultIdUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}
func TestAccAppServiceCertificateOrderCertificate_updateKeyVaultName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order_certificate", "test")
	r := CertificateOrderCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.keyVaultNameUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CertificateOrderCertificateResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appservicecertificateorders.ParseCertificateIDInsensitively(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.AppServiceCertificatesOrderClient.GetCertificate(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retreiving %s: %v", id, err)
	}

	if response.WasNotFound(resp.HttpResponse) {
		return utils.Bool(false), nil
	}
	return utils.Bool(true), nil
}

// Configs

func (r CertificateOrderCertificateResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_certificate_order_certificate" "test" {
  name                  = "acctestcokv-%[2]s"
  certificate_order_id  = azurerm_app_service_certificate_order.test.id
  key_vault_id          = azurerm_key_vault.test.id
  key_vault_secret_name = "kvsec%[2]s"
}
`, template, data.RandomStringOfLength(5))
}

func (r CertificateOrderCertificateResource) keyVaultIdUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_certificate_order_certificate" "test" {
  name                  = "acctestcokv-%[2]s"
  certificate_order_id  = azurerm_app_service_certificate_order.test.id
  key_vault_id          = azurerm_key_vault.test1.id
  key_vault_secret_name = "kvsec%[2]s"
}
`, template, data.RandomStringOfLength(5))
}

func (r CertificateOrderCertificateResource) keyVaultNameUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_certificate_order_certificate" "test" {
  name                  = "acctestcokv-%[2]s"
  certificate_order_id  = azurerm_app_service_certificate_order.test.id
  key_vault_id          = azurerm_key_vault.test.id
  key_vault_secret_name = "kvsec1%[2]s"
}
`, template, data.RandomStringOfLength(5))
}

func (r CertificateOrderCertificateResource) template(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dnsZoneRG := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

provider "azuread" {}

data "azurerm_client_config" "test" {}

data "azuread_service_principal" "cert-spn" {
  display_name = "Microsoft.Azure.CertificateRegistration"
}

data "azuread_service_principal" "app-service-spn" {
  display_name = "Microsoft Azure App Service"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tenant_id = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"


  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }

  // app service object ID
  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    //object_id = "f8daea97-62e7-4026-becf-13c2ea98e8b4"
    object_id = data.azuread_service_principal.app-service-spn.object_id

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }

  // Microsoft.Azure.CertificateRegistration 
  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azuread_service_principal.cert-spn.object_id
    //object_id = "ed47c2a1-bd23-4341-b39c-f4fd69138dd3"

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }
}

resource "azurerm_key_vault" "test1" {
  name                = "acctestkv1-%[3]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tenant_id = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"


  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }

  // app service
  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azuread_service_principal.app-service-spn.object_id
    //object_id = "f8daea97-62e7-4026-becf-13c2ea98e8b4"

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }

  // Microsoft.Azure.CertificateRegistration
  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azuread_service_principal.cert-spn.object_id
    //object_id = "ed47c2a1-bd23-4341-b39c-f4fd69138dd3"

    secret_permissions = [
      "Delete",
      "Get",
      "Purge",
      "Set",
      "List"
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Import",
      "List"
    ]
  }
}

data "azurerm_dns_zone" "test" {
  name                = "%[4]s"
  resource_group_name = "%[5]s"
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "tftestASCO-cert-%[3]s"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  distinguished_name  = "CN=${data.azurerm_dns_zone.test.name}"
  product_type        = "Standard"
}

resource "azurerm_dns_txt_record" "test" {
  name                = "@"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600

  record {
    value = azurerm_app_service_certificate_order.test.domain_verification_token
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5), dnsZone, dnsZoneRG)
}
