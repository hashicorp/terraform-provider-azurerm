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

type SpringCloudJavaDeploymentResource struct{}

func TestAccSpringCloudJavaDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_java_deployment", "test")
	r := SpringCloudJavaDeploymentResource{}

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

func TestAccSpringCloudJavaDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_java_deployment", "test")
	r := SpringCloudJavaDeploymentResource{}

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

func TestAccSpringCloudJavaDeployment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_java_deployment", "test")
	r := SpringCloudJavaDeploymentResource{}

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

func TestAccSpringCloudJavaDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_java_deployment", "test")
	r := SpringCloudJavaDeploymentResource{}

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

func TestAccSpringCloudJavaDeployment_updateHalfCpuMemory(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_java_deployment", "test")
	r := SpringCloudJavaDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.halfCpuMemory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.nonHalfCpuMemory(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudJavaDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r SpringCloudJavaDeploymentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_java_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
}
`, r.template(data), data.RandomString)
}

func (r SpringCloudJavaDeploymentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_java_deployment" "import" {
  name                = azurerm_spring_cloud_java_deployment.test.name
  spring_cloud_app_id = azurerm_spring_cloud_java_deployment.test.spring_cloud_app_id
}
`, r.basic(data))
}

func (r SpringCloudJavaDeploymentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_java_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  instance_count      = 2
  jvm_options         = "-XX:+PrintGC"
  runtime_version     = "Java_11"
  quota {
    cpu    = "2"
    memory = "2Gi"
  }
  environment_variables = {
    "Foo" : "Bar"
    "Env" : "Staging"
  }
}
`, r.template(data), data.RandomString)
}

func (r SpringCloudJavaDeploymentResource) halfCpuMemory(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_java_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  quota {
    cpu    = "500m"
    memory = "512Mi"
  }
}
`, r.template(data), data.RandomString)
}

func (r SpringCloudJavaDeploymentResource) nonHalfCpuMemory(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_java_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  quota {
    cpu    = "2"
    memory = "4Gi"
  }
}
`, r.template(data), data.RandomString)
}

func (SpringCloudJavaDeploymentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
