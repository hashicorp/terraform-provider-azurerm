package iotcentral_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IoTCentralApplicationResource struct {
}

func TestAccIoTCentralApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("ST1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("template").HasValue("iotc-default@1.0.0"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku").HasValue("ST1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIoTCentralApplication_requiresImportErrorStep(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")
	r := IoTCentralApplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_iotcentral_application"),
		},
	})
}

func (IoTCentralApplicationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ApplicationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.IoTCentral.AppsClient.Get(ctx, id.ResourceGroup, id.IoTAppName)
	if err != nil {
		return nil, fmt.Errorf("retrieving Analysis Services Server %q (resource group: %q): %+v", id.IoTAppName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.AppProperties != nil), nil
}

func (IoTCentralApplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "acctest-iotcentralapp-%[1]d"
  sku                 = "ST1"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (IoTCentralApplicationResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "acctest-iotcentralapp-%[2]d"
  display_name        = "acctest-iotcentralapp-%[2]d"
  sku                 = "ST1"
  template            = "iotc-default@1.0.0"
  tags = {
    ENV = "Test"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (IoTCentralApplicationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sub_domain          = "acctest-iotcentralapp-%[2]d"
  display_name        = "acctest-iotcentralapp-%[2]d"
  sku                 = "ST1"
  tags = {
    ENV = "Test"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func (IoTCentralApplicationResource) requiresImport(data acceptance.TestData) string {
	template := IoTCentralApplicationResource{}.basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_iotcentral_application" "import" {
  name                = azurerm_iotcentral_application.test.name
  resource_group_name = azurerm_iotcentral_application.test.resource_group_name
  location            = azurerm_iotcentral_application.test.location
  sub_domain          = azurerm_iotcentral_application.test.sub_domain
  display_name        = azurerm_iotcentral_application.test.display_name
  sku                 = "ST1"
}
`, template)
}
