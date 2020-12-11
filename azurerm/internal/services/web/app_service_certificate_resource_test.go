package web_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/web/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AppServiceCertificateResource struct{}

func TestAccAzureRMAppServiceCertificate_Pfx(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate", "test")
	r := AppServiceCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.pfx(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("password").HasValue("terraform"),
				check.That(data.ResourceName).Key("thumbprint").HasValue("7B985BF42467791F23E52B364A3E8DEBAB9C606E"),
			),
		},
		data.ImportStep("pfx_blob", "password"),
	})
}

func TestAccAzureRMAppServiceCertificate_PfxNoPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate", "test")
	r := AppServiceCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.pfxNoPassword(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("thumbprint").HasValue("7B985BF42467791F23E52B364A3E8DEBAB9C606E"),
			),
		},
		data.ImportStep("pfx_blob"),
	})
}

func TestAccAzureRMAppServiceCertificate_KeyVault(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_certificate", "test")
	r := AppServiceCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVault(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("thumbprint").HasValue("7B985BF42467791F23E52B364A3E8DEBAB9C606E"),
			),
		},
		data.ImportStep("key_vault_secret_id"),
	})
}

func (r AppServiceCertificateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Web.CertificatesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving App Service Certificate %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r AppServiceCertificateResource) pfx(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestwebcert%d"
  location = "%s"
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  pfx_blob            = filebase64("testdata/app_service_certificate.pfx")
  password            = "terraform"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AppServiceCertificateResource) pfxNoPassword(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestwebcert%d"
  location = "%s"
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  pfx_blob            = filebase64("testdata/app_service_certificate_nopassword.pfx")
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AppServiceCertificateResource) keyVault(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "test" {}

data "azuread_service_principal" "test" {
  display_name = "Microsoft Azure App Service"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestwebcert%d"
  location = "%s"
}

resource "azurerm_key_vault" "test" {
  name                = "acct%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tenant_id = data.azurerm_client_config.test.tenant_id

  sku_name = "standard"

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azurerm_client_config.test.object_id
    secret_permissions      = ["delete", "get", "set"]
    certificate_permissions = ["create", "delete", "get", "import"]
  }

  access_policy {
    tenant_id               = data.azurerm_client_config.test.tenant_id
    object_id               = data.azuread_service_principal.test.object_id
    secret_permissions      = ["get"]
    certificate_permissions = ["get"]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctest%d"
  key_vault_id = azurerm_key_vault.test.id

  certificate {
    contents = filebase64("testdata/app_service_certificate.pfx")
    password = "terraform"
  }

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }
  }
}

resource "azurerm_app_service_certificate" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  key_vault_secret_id = azurerm_key_vault_certificate.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
