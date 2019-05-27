package azurerm

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccAzureRMAnalysisServicesServer_basic(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAnalysisServicesServer_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
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

func TestAccAzureRMAnalysisServicesServer_withTags(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMAnalysisServicesServer_withTags(ri, testLocation())
	postConfig := testAccAzureRMAnalysisServicesServer_withTagsUpdate(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.label", "test1"),
					resource.TestCheckResourceAttr(resourceName, "tags.env", "prod"),
				),
			},
		},
	})
}

func TestAccAzureRMAnalysisServicesServer_querypoolConnectionMode(t *testing.T) {
	resourceName := "azurerm_analysis_services_server.test"
	ri := tf.AccRandTimeInt()
	preConfig := testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(ri, testLocation(), "All")
	postConfig := testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(ri, testLocation(), "ReadOnly")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAnalysisServicesServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "querypool_connection_mode", "All"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAnalysisServicesServerExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "querypool_connection_mode", "ReadOnly"),
				),
			},
		},
	})
}

func testAccAzureRMAnalysisServicesServer_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku 				  = "B1"
}
`, rInt, location, rInt)
}

func testAccAzureRMAnalysisServicesServer_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku 				  = "B1"

  tags = {
	label = "test"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAnalysisServicesServer_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                = "acctestass%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku 				  = "B1"

  tags = {
	label = "test1"
	env   = "prod"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMAnalysisServicesServer_querypoolConnectionMode(rInt int, location, connectionMode string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_analysis_services_server" "test" {
  name                		= "acctestass%d"
  location            		= "${azurerm_resource_group.test.location}"
  resource_group_name 		= "${azurerm_resource_group.test.name}"
  sku 				  		= "B1"
  querypool_connection_mode = "%s"
}
`, rInt, location, rInt, connectionMode)
}

func testCheckAzureRMAnalysisServicesServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).analysisServicesServerClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_analysis_services_server" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetDetails(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return nil
	}

	return nil
}

func testCheckAzureRMAnalysisServicesServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		analysisServicesServerName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Analysis Services Server: %s", analysisServicesServerName)
		}

		client := testAccProvider.Meta().(*ArmClient).analysisServicesServerClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.GetDetails(ctx, resourceGroup, analysisServicesServerName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Analysis Services Server %q (resource group: %q) does not exist", analysisServicesServerName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on analysisServicesServerClient: %+v", err)
		}

		return nil
	}
}
