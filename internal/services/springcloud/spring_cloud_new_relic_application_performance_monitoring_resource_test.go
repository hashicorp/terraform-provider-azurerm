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

type SpringCloudNewRelicApplicationPerformanceMonitoringResource struct{}

func TestAccSpringCloudNewRelicApplicationPerformanceMonitoring_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_new_relic_application_performance_monitoring", "test")
	r := SpringCloudNewRelicApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("license_key"),
	})
}

func TestAccSpringCloudNewRelicApplicationPerformanceMonitoring_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_new_relic_application_performance_monitoring", "test")
	r := SpringCloudNewRelicApplicationPerformanceMonitoringResource{}
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

func TestAccSpringCloudNewRelicApplicationPerformanceMonitoring_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_new_relic_application_performance_monitoring", "test")
	r := SpringCloudNewRelicApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("license_key"),
	})
}

func TestAccSpringCloudNewRelicApplicationPerformanceMonitoring_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_new_relic_application_performance_monitoring", "test")
	r := SpringCloudNewRelicApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("license_key"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("license_key"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("license_key"),
	})
}

func (r SpringCloudNewRelicApplicationPerformanceMonitoringResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appplatform.ParseApmID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.AppPlatformClient.ApmsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r SpringCloudNewRelicApplicationPerformanceMonitoringResource) template(data acceptance.TestData) string {
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

func (r SpringCloudNewRelicApplicationPerformanceMonitoringResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_new_relic_application_performance_monitoring" "test" {
  name                    = "acctest-apm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  app_name                = "test-app-name"
  license_key             = "test-license-key"
}
`, template, data.RandomInteger)
}

func (r SpringCloudNewRelicApplicationPerformanceMonitoringResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_new_relic_application_performance_monitoring" "import" {
  name                    = azurerm_spring_cloud_new_relic_application_performance_monitoring.test.name
  spring_cloud_service_id = azurerm_spring_cloud_new_relic_application_performance_monitoring.test.spring_cloud_service_id
  app_name                = "test-app-name"
  license_key             = "test-license-key"
}
`, config)
}

func (r SpringCloudNewRelicApplicationPerformanceMonitoringResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_new_relic_application_performance_monitoring" "test" {
  name                            = "acctest-apm-%[2]d"
  spring_cloud_service_id         = azurerm_spring_cloud_service.test.id
  app_name                        = "updated-app-name"
  license_key                     = "updated-license-key"
  agent_enabled                   = false
  app_server_port                 = 8080
  audit_mode_enabled              = true
  auto_app_naming_enabled         = true
  auto_transaction_naming_enabled = false
  custom_tracing_enabled          = false
  labels = {
    tagName1 = "tagValue1"
    tagName2 = "tagValue2"
  }
  globally_enabled = true
}
`, template, data.RandomInteger)
}
