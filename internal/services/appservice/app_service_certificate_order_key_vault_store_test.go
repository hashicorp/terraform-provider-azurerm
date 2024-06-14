package appservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-01-01/appservicecertificateorders"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CertificateOrderCertificateResource struct{}

func TestAccAppServiceCertificateOrderKeyVaultStore_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order_key_vault_store", "test")
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

func TestAccAppServiceCertificateOrderKeyVaultStore_updateKeyVaultId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order_key_vault_store", "test")
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
func TestAccAppServiceCertificateOrderKeyVaultStore_updateKeyVaultName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate_order_key_vault_store", "test")
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
	id, err := appservicecertificateorders.ParseCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.AppService.AppServiceCertificatesOrderClient.GetCertificate(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retreiving %s: %v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r CertificateOrderCertificateResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s
resource "azurerm_app_service_certificate_order_key_vault_store" "test" {
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
resource "azurerm_app_service_certificate_order_key_vault_store" "test" {
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
resource "azurerm_app_service_certificate_order_key_vault_store" "test" {
  name                  = "acctestcokv-%[2]s"
  certificate_order_id  = azurerm_app_service_certificate_order.test.id
  key_vault_id          = azurerm_key_vault.test.id
  key_vault_secret_name = "kvsec1%[2]s"
}
`, template, data.RandomStringOfLength(5))
}

func (r CertificateOrderCertificateResource) template(data acceptance.TestData) string {
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

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
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

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azuread_service_principal.cert-spn.object_id

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

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
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

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azuread_service_principal.cert-spn.object_id

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

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%[1]d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_app_service_certificate_order" "test" {
  name                = "tftestASCO-cert-%[3]s"
  location            = "global"
  resource_group_name = azurerm_resource_group.test.name
  distinguished_name  = "CN=${azurerm_dns_zone.test.name}"
  product_type        = "Standard"
}

resource "azurerm_dns_txt_record" "test" {
  name                = "@"
  zone_name           = azurerm_dns_zone.test.name
  resource_group_name = azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600

  record {
    value = azurerm_app_service_certificate_order.test.domain_verification_token
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomStringOfLength(5), data.RandomInteger)
}
