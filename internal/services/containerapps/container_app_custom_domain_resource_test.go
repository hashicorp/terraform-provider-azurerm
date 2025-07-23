// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppCustomDomainResource struct{}

func (r ContainerAppCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ContainerAppCustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	containerAppId := containerapps.NewContainerAppID(id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName)

	resp, err := client.ContainerApps.ContainerAppClient.Get(ctx, containerAppId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model
	if model == nil || model.Properties == nil || model.Properties.Configuration == nil || model.Properties.Configuration.Ingress == nil {
		return pointer.To(false), nil
	}

	ingress := *model.Properties.Configuration.Ingress

	if customDomains := ingress.CustomDomains; customDomains != nil {
		for _, v := range *customDomains {
			if strings.EqualFold(v.Name, id.CustomDomainName) {
				return pointer.To(true), nil
			}
		}
	}

	return pointer.To(false), nil
}

func TestAccContainerAppCustomDomainResource_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_custom_domain", "test")
	r := ContainerAppCustomDomainResource{}

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

func TestAccContainerAppCustomDomainResource_managedCertificate(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_custom_domain", "test")
	r := ContainerAppCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.managedCertificate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppCustomDomainResource_multiple(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_custom_domain", "test")
	r := ContainerAppCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerAppCustomDomainResource_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_custom_domain", "test")
	r := ContainerAppCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiple(data),
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

func (r ContainerAppCustomDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider azurerm {
  features {}
}

%s

resource "azurerm_container_app_custom_domain" "test" {
  name                                     = trimprefix(azurerm_dns_txt_record.test.fqdn, "asuid.")
  container_app_id                         = azurerm_container_app.test.id
  container_app_environment_certificate_id = azurerm_container_app_environment_certificate.test.id
  certificate_binding_type                 = "SniEnabled"
}

`, r.template(data))
}

func (r ContainerAppCustomDomainResource) managedCertificate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider azurerm {
  features {}
}

%s

resource "azurerm_container_app_custom_domain" "test" {
  name             = trimprefix(azurerm_dns_txt_record.test.fqdn, "asuid.")
  container_app_id = azurerm_container_app.test.id

  lifecycle {
    ignore_changes = [certificate_binding_type, container_app_environment_certificate_id]
  }
}


`, r.template(data))
}

func (r ContainerAppCustomDomainResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider azurerm {
  features {}
}

%s

resource "azurerm_container_app_custom_domain" "test" {
  name                                     = trimprefix(azurerm_dns_txt_record.test.fqdn, "asuid.")
  container_app_id                         = azurerm_container_app.test.id
  container_app_environment_certificate_id = azurerm_container_app_environment_certificate.test.id
  certificate_binding_type                 = "SniEnabled"
}

resource "azurerm_dns_txt_record" "test2" {
  name                = "asuid.containerapp%[2]d-2"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 300

  record {
    value = azurerm_container_app.test.custom_domain_verification_id
  }
}

resource "azurerm_container_app_environment_certificate" "test2" {
  name                         = "acctest-cacert%[2]d-2"
  container_app_environment_id = azurerm_container_app_environment.test.id
  certificate_blob_base64      = filebase64("testdata/testacc.pfx")
  certificate_password         = "TestAcc"
}

resource "azurerm_container_app_custom_domain" "test2" {
  name                                     = trimprefix(azurerm_dns_txt_record.test2.fqdn, "asuid.")
  container_app_id                         = azurerm_container_app.test.id
  container_app_environment_certificate_id = azurerm_container_app_environment_certificate.test.id
  certificate_binding_type                 = "SniEnabled"
}
`, r.template(data), data.RandomInteger)
}

func (r ContainerAppCustomDomainResource) template(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAE-%[1]d"
  location = "%[2]s"
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_dns_txt_record" "test" {
  name                = "asuid.containerapp%[1]d"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 300

  record {
    value = azurerm_container_app.test.custom_domain_verification_id
  }
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestCAEnv-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
}

resource "azurerm_container_app_environment_certificate" "test" {
  name                         = "acctest-cacert%[1]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  certificate_blob_base64      = filebase64("testdata/testacc.pfx")
  certificate_password         = "TestAcc"
}

resource "azurerm_container_app" "test" {
  name                         = "acctest-capp-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  container_app_environment_id = azurerm_container_app_environment.test.id
  revision_mode                = "Single"

  template {
    container {
      name   = "acctest-cont-%[1]d"
      image  = "jackofallops/azure-containerapps-python-acctest:v0.0.1"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }

  secret {
    name  = "rick"
    value = "morty"
  }

  ingress {
    allow_insecure_connections = false
    external_enabled           = true
    target_port                = 5000
    transport                  = "http"
    traffic_weight {
      latest_revision = true
      percentage      = 100
    }
  }
}

`, data.RandomInteger, data.Locations.Primary, dnsZone, dataResourceGroup)
}
