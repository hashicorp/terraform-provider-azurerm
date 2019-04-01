package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDataFactoryPipeline_basic(t *testing.T) {
	resourceName := "azurerm_data_factory_pipeline.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMDataFactoryPipeline_basic(ri, testLocation())
	config2 := testAccAzureRMDataFactoryPipeline_update(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDataFactoryPipelineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.test", "testparameter"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDataFactoryPipelineExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "parameters.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameters.test", "testparameter"),
					resource.TestCheckResourceAttr(resourceName, "parameters.test2", "testparameter2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMDataFactoryPipelineDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).dataFactoryPipelineClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_data_factory_pipeline" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		dataFactoryName := rs.Primary.Attributes["data_factory_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).dataFactoryPipelineClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, dataFactoryName, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: API Management Property %q (Resource Group %q / Data Factory %q) does not exist", name, resourceGroup, dataFactoryName)
			}
			return fmt.Errorf("Bad: Get on DataFactoryPipelineClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMDataFactoryPipeline_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_data_factory_v2" "test" {
  name                = "acctestdfv2%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory_v2.test.name}"

  parameters = {
	test = "testparameter"
  }
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDataFactoryPipeline_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_data_factory_v2" "test" {
  name                = "acctestdfv2%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_data_factory_pipeline" "test" {
  name                = "acctest%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  data_factory_name   = "${azurerm_data_factory_v2.test.name}"

  parameters = {
	test = "testparameter"
	test2 = "testparameter2"
  }
}
`, rInt, location, rInt, rInt)
}
