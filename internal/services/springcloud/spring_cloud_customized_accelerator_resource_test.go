// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudCustomizedAcceleratorResource struct{}

func TestAccSpringCloudCustomizedAccelerator_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_customized_accelerator", "test")
	r := SpringCloudCustomizedAcceleratorResource{}
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

func TestAccSpringCloudCustomizedAccelerator_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_customized_accelerator", "test")
	r := SpringCloudCustomizedAcceleratorResource{}
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

func TestAccSpringCloudCustomizedAccelerator_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_customized_accelerator", "test")
	r := SpringCloudCustomizedAcceleratorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudCustomizedAccelerator_caCertificateId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_customized_accelerator", "test")
	r := SpringCloudCustomizedAcceleratorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.caCertificateId(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudCustomizedAccelerator_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_customized_accelerator", "test")
	r := SpringCloudCustomizedAcceleratorResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudCustomizedAcceleratorResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appplatform.ParseCustomizedAcceleratorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.AppPlatformClient.CustomizedAcceleratorsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudCustomizedAcceleratorResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_accelerator" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudCustomizedAcceleratorResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_customized_accelerator" "test" {
  name                        = "acctest-ca-%[2]d"
  spring_cloud_accelerator_id = azurerm_spring_cloud_accelerator.test.id
  git_repository {
    url    = "https://github.com/Azure-Samples/piggymetrics"
    branch = "master"
  }
}
`, template, data.RandomInteger)
}

func (r SpringCloudCustomizedAcceleratorResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_customized_accelerator" "import" {
  name                        = azurerm_spring_cloud_customized_accelerator.test.name
  spring_cloud_accelerator_id = azurerm_spring_cloud_customized_accelerator.test.spring_cloud_accelerator_id
  git_repository {
    url    = "https://github.com/Azure-Samples/piggymetrics"
    branch = "master"
  }
}
`, config)
}

func (r SpringCloudCustomizedAcceleratorResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_customized_accelerator" "test" {
  name                        = "acctest-ca-%[2]d"
  spring_cloud_accelerator_id = azurerm_spring_cloud_accelerator.test.id

  git_repository {
    url                 = "https://github.com/sample-accelerators/fragments.git"
    branch              = "main"
    interval_in_seconds = 100
    path                = "java-version"
  }

  accelerator_tags = ["tag-a", "tag-b"]
  accelerator_type = "Fragment"
  description      = "test description"
  display_name     = "test name"
  icon_url         = "https://images.freecreatives.com/wp-content/uploads/2015/05/smiley-559124_640.jpg"
}
`, template, data.RandomInteger)
}

func (r SpringCloudCustomizedAcceleratorResource) caCertificateId(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Domain-Management"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest-kv-%[2]d"
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
  name         = "acctest-cert-%[2]d"
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


resource "azurerm_spring_cloud_certificate" "test" {
  name                     = "acctest-scc-%[2]d"
  resource_group_name      = azurerm_spring_cloud_service.test.resource_group_name
  service_name             = azurerm_spring_cloud_service.test.name
  key_vault_certificate_id = azurerm_key_vault_certificate.test.id
}

resource "azurerm_spring_cloud_customized_accelerator" "test" {
  name                        = "acctest-ca-%[2]d"
  spring_cloud_accelerator_id = azurerm_spring_cloud_accelerator.test.id

  git_repository {
    url               = "https://github.com/Azure-Samples/piggymetrics"
    branch            = "Azure"
    ca_certificate_id = azurerm_spring_cloud_certificate.test.id
  }
}
`, template, data.RandomIntOfLength(10))
}
