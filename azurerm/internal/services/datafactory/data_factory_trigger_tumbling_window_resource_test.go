package datafactory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccAzureRMDataFactoryTriggerTumblingWindow_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerDestroy("azurerm_data_factory_trigger_tumbling_window"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerTumblingWindow_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryTriggerTumblingWindow_startstop(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerDestroy("azurerm_data_factory_trigger_tumbling_window"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerTumblingWindow_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
					testCheckAzureRMDataFactoryTriggerStarts(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:             testAccAzureRMDataFactoryTriggerTumblingWindow_basic(data),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerStops(data.ResourceName),
				),
			},
			{
				Config:             testAccAzureRMDataFactoryTriggerTumblingWindow_basic(data),
				ExpectNonEmptyPlan: false,
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryTriggerTumblingWindow_trigger_dependency(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerDestroy("azurerm_data_factory_trigger_tumbling_window"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerTumblingWindow_trigger_dependency(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
					testCheckAzureRMDataFactoryTriggerStarts(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:             testAccAzureRMDataFactoryTriggerTumblingWindow_trigger_dependency(data),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerStops(data.ResourceName),
				),
			},
			{
				Config:             testAccAzureRMDataFactoryTriggerTumblingWindow_trigger_dependency(data),
				ExpectNonEmptyPlan: false,
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryTriggerTumblingWindow_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_trigger_tumbling_window", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryTriggerDestroy("azurerm_data_factory_trigger_tumbling_window"),
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryTriggerTumblingWindow_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataFactoryTriggerTumblingWindow_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerExists(data.ResourceName),
					testCheckAzureRMDataFactoryTriggerStarts(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config:             testAccAzureRMDataFactoryTriggerTumblingWindow_update(data),
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryTriggerStops(data.ResourceName),
				),
			},
			{
				Config:             testAccAzureRMDataFactoryTriggerTumblingWindow_update(data),
				ExpectNonEmptyPlan: false,
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMDataFactoryTriggerTumblingWindow_basic(data acceptance.TestData) string {
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

resource "azurerm_data_factory_trigger_tumbling_window" "test" {
  name                = "acctestdf%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  pipeline_name       = azurerm_data_factory_pipeline.test.name

  start_time = "2020-09-21T00:00:00Z"
  interval   = 24
  frequency  = "Hour"

  annotations = ["test1", "test2", "test3"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryTriggerTumblingWindow_update(data acceptance.TestData) string {
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

resource "azurerm_data_factory_trigger_tumbling_window" "test" {
  name                = "acctestdf%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name

  pipeline_parameters = {
    test = "@{formatDateTime(trigger().outputs.windowStartTime,'yyyy-MM-dd')}"
  }

  pipeline_name = azurerm_data_factory_pipeline.test.name

  interval        = 24
  frequency       = "Hour"
  max_concurrency = 3
  start_time      = "2020-09-21T00:00:00Z"
  end_time        = "2020-10-21T00:00:00Z"
  delay           = "16:00:00"

  trigger_dependency {
    size   = "24:00:00"
    offset = "-24:00:00"
  }

  retry {
    count    = 3
    interval = 60
  }

  annotations = ["test1", "test2", "test3"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryTriggerTumblingWindow_trigger_dependency(data acceptance.TestData) string {
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

resource "azurerm_data_factory_pipeline" "test2" {
  name                = "acctest%d-2"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name

  parameters = {
    test = "testparameter"
  }
}

resource "azurerm_data_factory_trigger_tumbling_window" "test2" {
  name                = "acctesttr%d-2"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  pipeline_name       = azurerm_data_factory_pipeline.test2.name

  start_time = "2020-09-21T00:00:00Z"
  interval   = 24
  frequency  = "Hour"

}

resource "azurerm_data_factory_trigger_tumbling_window" "test" {
  name                = "acctesttr%d"
  data_factory_name   = azurerm_data_factory.test.name
  resource_group_name = azurerm_resource_group.test.name
  pipeline_name       = azurerm_data_factory_pipeline.test.name

  start_time = "2020-09-21T00:00:00Z"
  interval   = 24
  frequency  = "Hour"

  trigger_dependency {
    trigger = azurerm_data_factory_trigger_tumbling_window.test2.name
  }

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
