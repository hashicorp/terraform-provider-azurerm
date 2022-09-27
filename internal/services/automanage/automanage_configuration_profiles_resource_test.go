package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutomanageConfigurationProfileResource struct{}

func TestAccAutomanageConfigurationProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileResource{}
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

func TestAccAutomanageConfigurationProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileResource{}
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

func TestAccAutomanageConfigurationProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomanageConfigurationProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileResource{}
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
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationProfileID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automanage.ConfigurationProfileClient.Get(ctx, id.ConfigurationProfileName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Automanage ConfigurationProfile %q (Resource Group %q): %+v", id.ConfigurationProfileName, id.ResourceGroup, err)
	}
	return utils.Bool(true), nil
}

func (r AutomanageConfigurationProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-automanage-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AutomanageConfigurationProfileResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile" "test" {
  name = "acctest-acp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
}
`, template, data.RandomInteger)
}

func (r AutomanageConfigurationProfileResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile" "import" {
  name = azurerm_automanage_configuration_profile.test.name
  resource_group_name = azurerm_automanage_configuration_profile.test.resource_group_name
  location = azurerm_automanage_configuration_profile.test.location
}
`, config)
}

func (r AutomanageConfigurationProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile" "test" {
  name = "acctest-acp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
