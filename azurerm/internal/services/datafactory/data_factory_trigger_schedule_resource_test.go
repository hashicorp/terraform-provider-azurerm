package datafactory_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataFactoryTriggerSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerDestroy("azurerm_data_factory_trigger_schedule"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerSchedule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryTriggerSchedule_complete(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	endTime := time.Now().UTC().Add(time.Hour * 7).In(loc).Format("2006-01-02T15:04:00Z07:00")
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerDestroy("azurerm_data_factory_trigger_schedule"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerSchedule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataFactoryTriggerSchedule_update(data, endTime),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMDataFactoryTriggerSchedule_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name                = "acctestdf%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  pipeline_name       = azurerm_data_factory_pipeline.test.name

  annotations = ["test1", "test2", "test3"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryTriggerSchedule_update(data acceptance.TestData, endTime string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-datafactory-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_schedule" "test" {
  name                = "acctestDFTS%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  pipeline_name       = azurerm_data_factory_pipeline.test.name

  pipeline_parameters = azurerm_data_factory_pipeline.test.parameters
  annotations         = ["test5"]
  frequency           = "Day"
  interval            = 5
  end_time            = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, endTime)
}
