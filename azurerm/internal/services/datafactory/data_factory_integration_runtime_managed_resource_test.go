package datafactory_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IntegrationRuntimeManagedResource struct {
}

func TestAccDataFactoryIntegrationRuntimeManaged_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("compute_type").HasValue("General"),
				check.That(data.ResourceName).Key("core_count").HasValue("8"),
				check.That(data.ResourceName).Key("time_to_live").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_computeType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.computeType(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("compute_type").HasValue("ComputeOptimized"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_coreCount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.coreCount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("core_count").HasValue("16"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDataFactoryIntegrationRuntimeManaged_timeToLive(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_integration_runtime_managed", "test")
	r := IntegrationRuntimeManagedResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.coreCount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("time_to_live").HasValue("10"),
			),
		},
		data.ImportStep(),
	})
}

func (IntegrationRuntimeManagedResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) computeType(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	compute_type 				= "ComputeOptimized"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) coreCount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	core_count 					= 16
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (IntegrationRuntimeManagedResource) timeToLive(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfirm%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_integration_runtime_managed" "test" {
  name                = "managed-integration-runtime"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
	location            = azurerm_resource_group.test.location
	time_to_live 				= 10
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (t IntegrationRuntimeManagedResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	dataFactoryName := id.Path["factories"]
	name := id.Path["integrationruntimes"]

	resp, err := clients.DataFactory.IntegrationRuntimesClient.Get(ctx, resourceGroup, dataFactoryName, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Data Factory Integration Runtime Managed (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
