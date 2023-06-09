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

type SpringCloudApplicationPerformanceMonitoringResource struct{}

func TestAccSpringCloudApplicationPerformanceMonitoring_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_performance_monitoring", "test")
	r := SpringCloudApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("secrets.%", "secrets.connection-string"),
	})
}

func TestAccSpringCloudApplicationPerformanceMonitoring_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_performance_monitoring", "test")
	r := SpringCloudApplicationPerformanceMonitoringResource{}
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

func TestAccSpringCloudApplicationPerformanceMonitoring_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_performance_monitoring", "test")
	r := SpringCloudApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("secrets.%", "secrets.connection-string"),
	})
}

func TestAccSpringCloudApplicationPerformanceMonitoring_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_application_performance_monitoring", "test")
	r := SpringCloudApplicationPerformanceMonitoringResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("secrets.%", "secrets.connection-string"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("secrets.%", "secrets.connection-string"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("secrets.%", "secrets.connection-string"),
	})
}

func (r SpringCloudApplicationPerformanceMonitoringResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudApplicationPerformanceMonitoringID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.ApmClient.Get(ctx, id.ResourceGroup, id.SpringName, id.ApmName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r SpringCloudApplicationPerformanceMonitoringResource) template(data acceptance.TestData) string {
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

func (r SpringCloudApplicationPerformanceMonitoringResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_application_performance_monitoring" "test" {
  name                    = "acctest-apm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  type                    = "ApplicationInsights"
  properties = {
    any-string    = "any-string"
    sampling-rate = "12.0"
  }
  secrets = {
    connection-string = "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXXXXXXXX"
  }
}
`, template, data.RandomInteger)
}

func (r SpringCloudApplicationPerformanceMonitoringResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_application_performance_monitoring" "import" {
  name                    = azurerm_spring_cloud_application_performance_monitoring.test.name
  spring_cloud_service_id = azurerm_spring_cloud_application_performance_monitoring.test.spring_cloud_service_id
  type                    = "ApplicationInsights"
  properties = {
    any-string    = "any-string"
    sampling-rate = "12.0"
  }
  secrets = {
    connection-string = "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXXXXXXXX"
  }
}
`, config)
}

func (r SpringCloudApplicationPerformanceMonitoringResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_application_performance_monitoring" "test" {
  name                    = "acctest-apm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  type                    = "ApplicationInsights"
  properties = {
    any-string    = "any-string"
    sampling-rate = "12.0"
  }
  secrets = {
    connection-string = "XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXX-XXXXXXXXXXXXXXXXXXX;XXXXXXXXXXXXXXXXX=XXXXXXXXXXXXXXXXXXX"
  }
  globally_enabled = true
}
`, template, data.RandomInteger)
}
