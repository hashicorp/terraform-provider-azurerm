// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudCertificateResource struct{}

func TestAccSpringCloudCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")
	r := SpringCloudCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		data.ImportStep("key_vault_certificate_id"),
	})
}

func TestAccSpringCloudCertificate_content(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")
	r := SpringCloudCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.content(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("thumbprint").Exists(),
			),
		},
		data.ImportStep("certificate_content"),
	})
}

func TestAccSpringCloudCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_certificate", "test")
	r := SpringCloudCertificateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t SpringCloudCertificateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
  exclude_private_key      = true
}
`, r.template(data), data.RandomInteger)
}

func (r SpringCloudCertificateResource) content(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[1]d"
  location = "%[2]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_certificate" "test" {
  name                = "acctest-scc-%[1]d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
  certificate_content = filebase64("testdata/cer")
}
`, data.RandomInteger, data.Locations.Primary)
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

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Resource Provider"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkeyvault%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Set",
    ]

    certificate_permissions = [
      "Create",
      "Delete",
      "Get",
      "Purge",
      "Update",
    ]
  }

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azuread_service_principal.test.object_id

    secret_permissions = [
      "Get",
      "List",
    ]

    certificate_permissions = [
      "Get",
      "List",
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
