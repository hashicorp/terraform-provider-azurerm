package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SpringCloudConfigServerResource struct {
}

func TestAccSpringCloudConfigServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_config_server", "test")
	r := SpringCloudConfigServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudConfigServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_config_server", "test")
	r := SpringCloudConfigServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			// those field returned by api are "*"
			// import state verify ignore those fields
			"ssh_auth.0.private_key",
			"ssh_auth.0.host_key",
			"ssh_auth.0.host_key_algorithm",
			"repository.0.http_basic_auth.0.username",
			"repository.0.http_basic_auth.0.password",
		),
	})
}

func TestAccSpringCloudConfigServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_config_server", "test")
	r := SpringCloudConfigServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSpringCloudConfigServer_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_config_server", "test")
	r := SpringCloudConfigServerResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			// those field returned by api are "*"
			// import state verify ignore those fields
			"ssh_auth.0.private_key",
			"ssh_auth.0.host_key",
			"ssh_auth.0.host_key_algorithm",
			"repository.0.http_basic_auth.0.username",
			"repository.0.http_basic_auth.0.password",
		),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r SpringCloudConfigServerResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.ConfigServersClient.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return nil, fmt.Errorf("unable to read Spring Cloud Service %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil && resp.Properties != nil && resp.Properties.ConfigServer != nil && resp.Properties.ConfigServer.GitProperty != nil), nil
}

func (r SpringCloudConfigServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_config_server" "test" {
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  uri                     = "https://github.com/Azure-Samples/piggymetrics-config"
}
`, r.template(data))
}

func (r SpringCloudConfigServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_config_server" "import" {
  spring_cloud_service_id = azurerm_spring_cloud_config_server.test.spring_cloud_service_id
  uri                     = azurerm_spring_cloud_config_server.test.uri
}
`, r.basic(data))
}

func (r SpringCloudConfigServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_config_server" "test" {
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  uri                     = "git@bitbucket.org:Azure-Samples/piggymetrics.git"
  label                   = "config"
  search_paths            = ["dir1", "dir4"]

  ssh_auth {
    private_key                      = file("testdata/private_key")
    host_key                         = file("testdata/host_key")
    host_key_algorithm               = "ssh-rsa"
    strict_host_key_checking_enabled = false
  }

  repository {
    name         = "repo1"
    uri          = "https://github.com/Azure-Samples/piggymetrics"
    label        = "config"
    search_paths = ["dir1", "dir2"]
    http_basic_auth {
      username = "username"
      password = "password"
    }
  }

  repository {
    name         = "repo2"
    uri          = "https://github.com/Azure-Samples/piggymetrics"
    label        = "config"
    search_paths = ["dir1", "dir2"]
  }
}
`, r.template(data))
}

func (SpringCloudConfigServerResource) template(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
