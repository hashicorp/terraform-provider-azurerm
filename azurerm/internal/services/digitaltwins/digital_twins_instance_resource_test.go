package digitaltwins_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DigitalTwinsInstanceResource struct {
}

func TestAccDigitalTwinsInstance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDigitalTwinsInstance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccDigitalTwinsInstance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDigitalTwinsInstance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_instance", "test")
	r := DigitalTwinsInstanceResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("host_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (DigitalTwinsInstanceResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DigitalTwinsInstanceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.DigitalTwins.InstanceClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Digital Twins Instance %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (DigitalTwinsInstanceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dtwin-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r DigitalTwinsInstanceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsInstanceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "import" {
  name                = azurerm_digital_twins_instance.test.name
  resource_group_name = azurerm_digital_twins_instance.test.resource_group_name
  location            = azurerm_digital_twins_instance.test.location
}
`, r.basic(data))
}

func (r DigitalTwinsInstanceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsInstanceResource) updateTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_instance" "test" {
  name                = "acctest-DT-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    ENV = "Stage"
  }
}
`, r.template(data), data.RandomInteger)
}
