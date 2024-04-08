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

type SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource struct{}

func TestAccSpringCloudApplicationInsightsApplicationPerformanceMonitoringrformanceMonitoring_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_insights_application_performance_monitoring", "test")
	r := SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccSpringCloudApplicationInsightsApplicationPerformanceMonitoringrformanceMonitoring_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_insights_application_performance_monitoring", "test")
	r := SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource{}
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

func TestAccSpringCloudApplicationInsightsApplicationPerformanceMonitoringrformanceMonitoring_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_insights_application_performance_monitoring", "test")
	r := SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func TestAccSpringCloudApplicationInsightsApplicationPerformanceMonitoringrformanceMonitoring_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_insights_application_performance_monitoring", "test")
	r := SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("connection_string"),
	})
}

func (r SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appplatform.ParseApmID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.AppPlatformClient.ApmsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_application_insights_application_performance_monitoring" "test" {
  name                    = "acctest-apm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  connection_string       = azurerm_application_insights.test.instrumentation_key
}
`, template, data.RandomInteger)
}

func (r SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_application_insights_application_performance_monitoring" "import" {
  name                    = azurerm_spring_cloud_application_insights_application_performance_monitoring.test.name
  spring_cloud_service_id = azurerm_spring_cloud_application_insights_application_performance_monitoring.test.spring_cloud_service_id
  connection_string       = azurerm_spring_cloud_application_insights_application_performance_monitoring.test.connection_string
}
`, config)
}

func (r SpringCloudApplicationInsightsApplicationPerformanceMonitoringResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_application_insights_application_performance_monitoring" "test" {
  name                         = "acctest-apm-%[2]d"
  spring_cloud_service_id      = azurerm_spring_cloud_service.test.id
  connection_string            = azurerm_application_insights.test.instrumentation_key
  globally_enabled             = true
  role_name                    = "test-role"
  role_instance                = "test-instance"
  sampling_percentage          = 50
  sampling_requests_per_second = 10
}
`, template, data.RandomInteger)
}
