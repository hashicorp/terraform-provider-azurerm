package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutoManageConfigurationProfileVersionResource struct{}

func TestAccAutoManageConfigurationProfileVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_version", "test")
	r := AutoManageConfigurationProfileVersionResource{}
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

func TestAccAutoManageConfigurationProfileVersion_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_version", "test")
	r := AutoManageConfigurationProfileVersionResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("antimalware.#").HasValue("0"),
				check.That(data.ResourceName).Key("automation_account_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("antimalware.#").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.#").HasValue("1"),
				check.That(data.ResourceName).Key("automation_account_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func (r AutoManageConfigurationProfileVersionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_automanage_configuration_version" "test" {
  name                       = "acctest-amcpv-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  configuration_profile_name = azurerm_automanage_configuration.test.name
}
`, template, data.RandomInteger)
}

func (r AutoManageConfigurationProfileVersionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_automanage_configuration_version" "test" {
  name                       = "acctest-amcpv-%d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  configuration_profile_name = azurerm_automanage_configuration.test.name

  antimalware {
    exclusions {
      extensions = "exe;dll"
    }
    real_time_protection_enabled = true
  }
  automation_account_enabled = true
}
`, template, data.RandomInteger)
}

func (r AutoManageConfigurationProfileVersionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AutoManageConfigurationProfileVersionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationVersionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Automanage.ConfigurationVersionClient
	resp, err := client.Get(ctx, id.ConfigurationProfileName, id.VersionName, id.ResourceGroup)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Response.Response != nil), nil
}
