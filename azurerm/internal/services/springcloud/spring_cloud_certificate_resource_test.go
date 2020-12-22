package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SpringCloudCertificateResource struct {
}

func TestAccSpringCloudCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")
	r := SpringCloudCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		data.ImportStep("key_vault_certificate_id"),
	})
}

func TestAccSpringCloudCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")
	r := SpringCloudCertificateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t SpringCloudCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudCertificateID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.CertificatesClient.Get(ctx, id.ResourceGroup, id.SpringName, id.CertificateName)
	if err != nil {
		return nil, fmt.Errorf("reading Spring Cloud Certificate %q (Spring Cloud Service %q / Resource Group %q): %+v", id.CertificateName, id.SpringName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r SpringCloudCertificateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_certificate" "test" {
  name                     = "acctest-scc-%d"
  resource_group_name      = azurerm_spring_cloud_service.test.resource_group_name
  service_name             = azurerm_spring_cloud_service.test.name
  key_vault_certificate_id = azurerm_key_vault_certificate.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r SpringCloudCertificateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_certificate" "import" {
  name                     = azurerm_spring_cloud_certificate.test.name
  resource_group_name      = azurerm_spring_cloud_certificate.test.resource_group_name
  service_name             = azurerm_spring_cloud_certificate.test.service_name
  key_vault_certificate_id = azurerm_spring_cloud_certificate.test.key_vault_certificate_id
}
`, r.basic(data))
}

func (r SpringCloudCertificateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Domain-Management"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
  soft_delete_enabled = true

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "set",
    ]

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "purge",
      "update",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.test.object_id

    secret_permissions = [
      "get",
      "list",
    ]

    certificate_permissions = [
      "get",
      "list",
    ]
  }
}

resource "azurerm_key_vault_certificate" "test" {
  name         = "acctestcert%s"
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
      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=contoso.com"
      validity_in_months = 12
    }
  }
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString, data.RandomInteger)
}
