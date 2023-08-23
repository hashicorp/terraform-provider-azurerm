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

type SpringCloudActiveDeploymentResource struct{}

func TestAccSpringCloudActiveDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_active_deployment", "test")
	r := SpringCloudActiveDeploymentResource{}

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

func TestAccSpringCloudActiveDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_active_deployment", "test")
	r := SpringCloudActiveDeploymentResource{}

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

func TestAccSpringCloudActiveDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_active_deployment", "test")
	r := SpringCloudActiveDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudActiveDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudAppID(state.ID)
	if err != nil {
		return nil, err
	}

	it, err := clients.AppPlatform.DeploymentsClient.ListComplete(ctx, id.ResourceGroup, id.SpringName, id.AppName, nil)
	if err != nil {
		return nil, fmt.Errorf("listing active deployment for %q: %+v", id, err)
	}
	deployments := make([]string, 0)
	for it.NotDone() {
		value := it.Value()
		if value.Properties != nil && value.Properties.Active != nil && *value.Properties.Active {
			deployments = append(deployments, *value.Name)
		}
		if err := it.NextWithContext(ctx); err != nil {
			return nil, fmt.Errorf("listing active deployment for %q: %+v", id, err)
		}
	}

	return utils.Bool(len(deployments) != 0), nil
}

func (r SpringCloudActiveDeploymentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_active_deployment" "test" {
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  deployment_name     = azurerm_spring_cloud_java_deployment.test.name
}
`, r.template(data))
}

func (r SpringCloudActiveDeploymentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_active_deployment" "import" {
  spring_cloud_app_id = azurerm_spring_cloud_active_deployment.test.spring_cloud_app_id
  deployment_name     = azurerm_spring_cloud_active_deployment.test.deployment_name
}
`, r.basic(data))
}

func (r SpringCloudActiveDeploymentResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_java_deployment" "test2" {
  name                = "acctest-scjd2%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
}

resource "azurerm_spring_cloud_active_deployment" "test" {
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
  deployment_name     = azurerm_spring_cloud_java_deployment.test2.name
}
`, r.template(data), data.RandomString)
}

func (SpringCloudActiveDeploymentResource) template(data acceptance.TestData) string {
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

resource "azurerm_spring_cloud_java_deployment" "test" {
  name                = "acctest-scjd%s"
  spring_cloud_app_id = azurerm_spring_cloud_app.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomString)
}
