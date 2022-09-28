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

type AutomanageConfigurationProfilesVersionResource struct{}

func TestAccAutomanageConfigurationProfilesVersion_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionResource{}
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

func TestAccAutomanageConfigurationProfilesVersion_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionResource{}
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

func TestAccAutomanageConfigurationProfilesVersion_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionResource{}
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

func TestAccAutomanageConfigurationProfilesVersion_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionResource{}
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

func (r AutomanageConfigurationProfilesVersionResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.AutomanageConfigurationProfilesVersionID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automanage.ConfigurationProfilesVersionClient.Get(ctx, id.ConfigurationProfileName, id.ConfigurationProfileName, id.ResourceGroup)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Automanage ConfigurationProfilesVersion %q (Resource Group %q / configurationProfileName %q): %+v", id.ConfigurationProfileName, id.ResourceGroup, id.ConfigurationProfileName, err)
	}
	return utils.Bool(true), nil
}

func (r AutomanageConfigurationProfilesVersionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-automanage-%d"
  location = "%s"
}

resource "azurerm_automanage_configuration_profile" "test" {
  name = "acctest-acp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r AutomanageConfigurationProfilesVersionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profiles_version" "test" {
  name = "acctest-acpv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  configuration_profile_name = azurerm_automanage_configuration_profile.test.name
}
`, template, data.RandomInteger)
}

func (r AutomanageConfigurationProfilesVersionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profiles_version" "import" {
  name = azurerm_automanage_configuration_profiles_version.test.name
  resource_group_name = azurerm_automanage_configuration_profiles_version.test.resource_group_name
  location = azurerm_automanage_configuration_profiles_version.test.location
  configuration_profile_name = azurerm_automanage_configuration_profiles_version.test.configuration_profile_name
}
`, config)
}

func (r AutomanageConfigurationProfilesVersionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profiles_version" "test" {
  name = "acctest-acpv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration_profile_name = azurerm_automanage_configuration_profile.test.name
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":false,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}
