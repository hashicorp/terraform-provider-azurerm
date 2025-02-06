// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudConfigurationServiceResource struct{}

func TestAccSpringCloudConfigurationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
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

func TestAccSpringCloudConfigurationService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
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

func TestAccSpringCloudConfigurationService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("repository.0.password", "repository.0.username"),
	})
}

func TestAccSpringCloudConfigurationService_sshAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.sshAuth(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("repository.0.host_key", "repository.0.host_key_algorithm", "repository.0.private_key"),
	})
}

func TestAccSpringCloudConfigurationService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
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
		data.ImportStep("repository.0.password", "repository.0.username"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudConfigurationService_generation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.generation(data, "Gen1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.generation(data, "Gen2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudConfigurationService_caCertificateId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_configuration_service", "test")
	r := SpringCloudConfigurationServiceResource{}
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

func (r SpringCloudConfigurationServiceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appplatform.ParseConfigurationServiceID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.AppPlatformClient.ConfigurationServicesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r SpringCloudConfigurationServiceResource) template(data acceptance.TestData) string {
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
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudConfigurationServiceResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_configuration_service" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
}
`, template)
}

func (r SpringCloudConfigurationServiceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_configuration_service" "import" {
  name                    = azurerm_spring_cloud_configuration_service.test.name
  spring_cloud_service_id = azurerm_spring_cloud_configuration_service.test.spring_cloud_service_id
}
`, config)
}

func (r SpringCloudConfigurationServiceResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_configuration_service" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  repository {
    name                     = "fake"
    label                    = "master"
    patterns                 = ["app/dev"]
    uri                      = "https://github.com/Azure-Samples/piggymetrics"
    search_paths             = ["dir1", "dir2"]
    strict_host_key_checking = false
    username                 = "adminuser"
    password                 = "H@Sh1CoR3!"
  }
  refresh_interval_in_seconds = 10
}
`, template)
}

func (r SpringCloudConfigurationServiceResource) sshAuth(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_configuration_service" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  repository {
    name                     = "fake"
    label                    = "master"
    patterns                 = ["app/dev"]
    uri                      = "https://github.com/Azure-Samples/piggymetrics"
    host_key                 = file("testdata/host_key")
    host_key_algorithm       = "ssh-rsa"
    private_key              = file("testdata/private_key")
    search_paths             = ["dir1", "dir2"]
    strict_host_key_checking = false
  }
}
`, template)
}

func (r SpringCloudConfigurationServiceResource) generation(data acceptance.TestData, generation string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_configuration_service" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  generation              = "%s"
}
`, template, generation)
}

func (r SpringCloudConfigurationServiceResource) caCertificateId(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s


data "azurerm_client_config" "current" {
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Resource Provider"
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
  exclude_private_key      = true
}

resource "azurerm_spring_cloud_configuration_service" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  generation              = "Gen2"
  repository {
    name              = "fake"
    label             = "master"
    patterns          = ["app/dev"]
    uri               = "https://github.com/Azure-Samples/piggymetrics"
    ca_certificate_id = azurerm_spring_cloud_certificate.test.id
  }
}
`, template, data.RandomIntOfLength(10))
}
