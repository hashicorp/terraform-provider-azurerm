package datafactory_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactoryTriggerSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_schedule", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerSchedule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerScheduleExists(data.ResourceName),
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
		CheckDestroy: testCheckAzureRMDataFactoryTriggerScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerSchedule_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerScheduleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataFactoryTriggerSchedule_update(data, endTime),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerScheduleExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataFactoryTriggerScheduleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.TriggersClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Data Factory: %s", name)
		}

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on dataFactory.TriggersClient: %+v", err)
		}

		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Bad: Data Factory Trigger Schdule %q (data factory name: %q / resource group: %q) does not exist", name, dataFactoryName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryTriggerScheduleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.TriggersClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_trigger_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Data Factory Trigger Schedule still exists:\n%#v", resp.Properties)
		}
	}

	return nil
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
