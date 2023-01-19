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
	id, err := parse.SpringCloudCustomizedAcceleratorID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.AppPlatform.CustomizedAcceleratorClient.Get(ctx, id.ResourceGroup, id.SpringName, id.ApplicationAcceleratorName, id.CustomizedAcceleratorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
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
    url                 = "https://github.com/Azure-Samples/piggymetrics"
    git_tag             = "spring.version.2.0.3"
    interval_in_seconds = 100
  }

  accelerator_tags = ["tag-a", "tag-b"]
  description      = "test description"
  display_name     = "test name"
  icon_url         = "https://images.freecreatives.com/wp-content/uploads/2015/05/smiley-559124_640.jpg"
}
`, template, data.RandomInteger)
}
