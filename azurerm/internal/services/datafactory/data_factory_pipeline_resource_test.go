package datafactory_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactoryPipeline_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryPipeline_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryPipeline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryPipeline_update1(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "annotations.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "test description"),
					resource.TestCheckResourceAttr(data.ResourceName, "variables.%", "2"),
				),
			},
			{
				Config: testAccAzureRMDataFactoryPipeline_update2(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "parameters.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "annotations.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "test description2"),
					resource.TestCheckResourceAttr(data.ResourceName, "variables.%", "3"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMDataFactoryPipeline_activities(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_data_factory_pipeline", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDataFactoryPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDataFactoryPipeline_activities(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "activities_json"),
					testCheckAzureRMDataFactoryPipelineHasAppenVarActivity(data.ResourceName, "Append variable1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataFactoryPipeline_activitiesUpdated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "activities_json"),
					testCheckAzureRMDataFactoryPipelineHasAppenVarActivity(data.ResourceName, "Append variable1"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMDataFactoryPipeline_activities(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "activities_json"),
					testCheckAzureRMDataFactoryPipelineHasAppenVarActivity(data.ResourceName, "Append variable1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMDataFactoryPipelineDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.PipelinesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_pipeline" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}
	return nil
}

func testCheckAzureRMDataFactoryPipelineExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.PipelinesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Data Factory Pipeline %q (Resource Group %q / Data Factory %q) does not exist", name, resourceGroup, dataFactoryName)
			}
			return fmt.Errorf("Bad: Get on DataFactoryPipelineClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDataFactoryPipelineHasAppenVarActivity(resourceName string, activityName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).DataFactory.PipelinesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Data Factory Pipeline %q (Resource Group %q / Data Factory %q) does not exist", name, resourceGroup, dataFactoryName)
			}
			return fmt.Errorf("Bad: Get on DataFactoryPipelineClient: %+v", err)
		}

		activities := *resp.Activities
		if len(activities) == 0 {
			return fmt.Errorf("Bad: No activities associated with Data Factory Pipeline %q (Resource Group %q / Data Factory %q)", name, resourceGroup, dataFactoryName)
		}
		appvarActivity, _ := activities[0].AsAppendVariableActivity()
		if *appvarActivity.Name != activityName {
			return fmt.Errorf("Bad: Data Factory Pipeline %q (Resource Group %q / Data Factory %q) could not cast as activity", name, resourceGroup, dataFactoryName)
		}

		return nil
	}
}

func testAccAzureRMDataFactoryPipeline_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryPipeline_update1(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  annotations         = ["test1", "test2", "test3"]
  description         = "test description"

  parameters = {
    test = "testparameter"
  }

  variables = {
    foo = "test1"
    bar = "test2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryPipeline_update2(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-df-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  annotations         = ["test1", "test2"]
  description         = "test description2"

  parameters = {
    test  = "testparameter"
    test2 = "testparameter2"
  }

  variables = {
    foo = "test1"
    bar = "test2"
    baz = "test3"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryPipeline_activities(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMDataFactoryPipeline_activitiesUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdfv2%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = azurerm_resource_group.test.name
  data_factory_name   = azurerm_data_factory.test.name
  variables = {
    "bob" = "item1"
  }
  activities_json = <<JSON
[
  {
    "name": "Append variable1",
    "type": "AppendVariable",
    "dependsOn": [],
    "userProperties": [],
    "typeProperties": {
      "variableName": "bob",
      "value": "something"
    }
  }
]
JSON
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
