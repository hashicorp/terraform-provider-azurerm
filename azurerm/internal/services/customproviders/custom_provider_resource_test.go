package customproviders_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/customproviders/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CustomProviderResource struct{}

func TestAccAzureRMCustomProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")
	r := CustomProviderResource{}

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

func TestAccAzureRMCustomProvider_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")
	r := CustomProviderResource{}

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

func TestAccAzureRMCustomProvider_action(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_provider", "test")
	r := CustomProviderResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.action(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.actionUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.action(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CustomProviderResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.CustomProviderID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.CustomProviders.CustomProviderClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Custom Provider %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(true), nil
}

func (r CustomProviderResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CustomProviderResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  resource_type {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }

  action {
    name     = "dEf2"
    endpoint = "https://example.com/"
  }

  validation {
    specification = "https://raw.githubusercontent.com/Azure/azure-custom-providers/master/CustomRPWithSwagger/Artifacts/Swagger/pingaction.json"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CustomProviderResource) action(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  action {
    name     = "dEf1"
    endpoint = "https://testendpoint.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r CustomProviderResource) actionUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-cp-%d"
  location = "%s"
}
resource "azurerm_custom_provider" "test" {
  name                = "accTEst_saa%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  action {
    name     = "dEf2"
    endpoint = "https://example.com/"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
