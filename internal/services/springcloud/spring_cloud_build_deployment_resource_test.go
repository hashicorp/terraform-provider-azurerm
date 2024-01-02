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

type SpringCloudBuildDeploymentResource struct{}

func TestAccSpringCloudBuildDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_deployment", "test")
	r := SpringCloudBuildDeploymentResource{}

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

func TestAccSpringCloudBuildDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_deployment", "test")
	r := SpringCloudBuildDeploymentResource{}

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

func TestAccSpringCloudBuildDeployment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_deployment", "test")
	r := SpringCloudBuildDeploymentResource{}

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

func TestAccSpringCloudBuildDeployment_addon(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_deployment", "test")
	r := SpringCloudBuildDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.addon(data, "auth-service"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.addon(data, "gateway"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudBuildDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_build_deployment", "test")
	r := SpringCloudBuildDeploymentResource{}

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

func (r SpringCloudBuildDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.DeploymentsClient.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.DeploymentName)
	if err != nil {
		return nil, fmt.Errorf("reading Spring Cloud Deployment %q (Spring Cloud Service %q / App %q / Resource Group %q): %+v", id.DeploymentName, id.SpringName, id.AppName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (r SpringCloudBuildDeploymentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  build_result_id     = "<default>"
}
`, r.template(data), data.RandomString)
}

func (r SpringCloudBuildDeploymentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_deployment" "import" {
  name                = azurerm_spring_cloud_build_deployment.test.name
  spring_cloud_app_id = azurerm_spring_cloud_build_deployment.test.spring_cloud_app_id
  build_result_id     = "<default>"
}
`, r.basic(data))
}

func (r SpringCloudBuildDeploymentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_deployment" "test" {
  name                                   = "acctest-scjd%s"
  spring_cloud_app_id                    = azurerm_spring_cloud_app.test.id
  build_result_id                        = "<default>"
  instance_count                         = 2
  application_performance_monitoring_ids = [azurerm_spring_cloud_application_insights_application_performance_monitoring.test.id]

  environment_variables = {
    "Foo" : "Bar"
    "Env" : "Staging"
  }
  quota {
    cpu    = "2"
    memory = "2Gi"
  }
}
`, r.template(data), data.RandomString)
}

func (r SpringCloudBuildDeploymentResource) addon(data acceptance.TestData, pattern string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_build_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  build_result_id     = "<default>"
  instance_count      = 2

  environment_variables = {
    "Foo" : "Bar"
    "Env" : "Staging"
  }
  quota {
    cpu    = "2"
    memory = "2Gi"
  }
  addon_json = jsonencode({
    applicationConfigurationService = {
      configFilePatterns = "%s"
    }
  })
}
`, SpringCloudAppResource{}.addon(data), data.RandomString, pattern)
}

func (SpringCloudBuildDeploymentResource) template(data acceptance.TestData) string {
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

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%[2]d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
}


resource "azurerm_application_insights" "test" {
  name                = "acctest-ai-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_spring_cloud_application_insights_application_performance_monitoring" "test" {
  name                    = "acctest-apm-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  connection_string       = azurerm_application_insights.test.instrumentation_key
}
`, data.Locations.Primary, data.RandomInteger)
}
