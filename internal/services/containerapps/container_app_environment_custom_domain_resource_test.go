// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containerapps_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ContainerAppEnvironmentCustomDomainResource struct{}

func (r ContainerAppEnvironmentCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := managedenvironments.ParseManagedEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ContainerApps.ManagedEnvironmentClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	model := resp.Model

	if *model.Properties.CustomDomainConfiguration.DnsSuffix == os.Getenv("ARM_TEST_DNS_ZONE") {
		return pointer.To(true), nil
	}

	return pointer.To(false), nil
}

func TestAccContainerAppEnvironmentCustomDomainResource_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_custom_domain", "test")
	r := ContainerAppEnvironmentCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_blob_base64", "certificate_password"),
	})
}

func TestAccContainerAppEnvironmentCustomDomainResource_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_custom_domain", "test")
	r := ContainerAppEnvironmentCustomDomainResource{}

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

func TestAccContainerAppEnvironmentCustomDomainResource_update(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skipf("Skipping as either ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_container_app_environment_custom_domain", "test")
	r := ContainerAppEnvironmentCustomDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_blob_base64", "certificate_password"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("certificate_blob_base64", "certificate_password"),
	})
}

func (r ContainerAppEnvironmentCustomDomainResource) basic(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")

	return fmt.Sprintf(`
provider azurerm {
  features {}
}

%s

resource "azurerm_container_app_environment_custom_domain" "test" {
  container_app_environment_id = azurerm_container_app_environment.test.id
  dns_suffix                   = "%[2]s"
  certificate_blob_base64      = filebase64("testdata/testacc.pfx")
  certificate_password         = "TestAcc"

  depends_on = [
    time_sleep.wait_60_seconds
  ]
}

`, r.template(data), dnsZone)
}

func (r ContainerAppEnvironmentCustomDomainResource) requiresImport(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")

	return fmt.Sprintf(`

%s

resource "azurerm_container_app_environment_custom_domain" "example" {
  container_app_environment_id = azurerm_container_app_environment.test.id
  dns_suffix                   = "%[2]s"
  certificate_blob_base64      = filebase64("testdata/testacc.pfx")
  certificate_password         = "TestAcc"
}

`, r.basic(data), dnsZone)
}

func (r ContainerAppEnvironmentCustomDomainResource) update(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")

	return fmt.Sprintf(`
provider azurerm {
  features {}
}

%s

resource "azurerm_container_app_environment_custom_domain" "test" {
  container_app_environment_id = azurerm_container_app_environment.test.id
  dns_suffix                   = "%[2]s"
  certificate_blob_base64      = filebase64("testdata/testacc_nopassword.pfx")
  certificate_password         = ""

  depends_on = [
    time_sleep.wait_60_seconds
  ]
}

`, r.template(data), dnsZone)
}

func (r ContainerAppEnvironmentCustomDomainResource) template(data acceptance.TestData) string {
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
  name                = "asuid"
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  zone_name           = data.azurerm_dns_zone.test.name
  ttl                 = 300

  record {
    value = azurerm_container_app_environment.test.custom_domain_verification_id
  }
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [azurerm_dns_txt_record.test]

  create_duration = "60s"
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

`, data.RandomInteger, data.Locations.Primary, dnsZone, dataResourceGroup)
}
