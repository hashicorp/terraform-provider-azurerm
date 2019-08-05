package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBatchApplication_basic(t *testing.T) {
	resourceName := "azurerm_batch_application.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBatchApplication_template(ri, rs, location, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchApplicationExists(resourceName),
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

func TestAccAzureRMBatchApplication_update(t *testing.T) {
	resourceName := "azurerm_batch_application.test"
	ri := tf.AccRandTimeInt()
	rs := strings.ToLower(acctest.RandString(11))
	location := testLocation()
	displayName := fmt.Sprintf("TestAccDisplayName-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMBatchApplication_template(ri, rs, location, ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchApplicationExists(resourceName),
				),
			},
			{
				Config: testAccAzureRMBatchApplication_template(ri, rs, location, fmt.Sprintf(`display_name = "%s"`, displayName)),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
				),
			},
		},
	})
}

func testCheckAzureRMBatchApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Batch Application not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["account_name"]

		client := testAccProvider.Meta().(*ArmClient).batch.ApplicationClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Batch Application %q (Account Name %q / Resource Group %q) does not exist", name, accountName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on batchApplicationClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMBatchApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).batch.ApplicationClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_application" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["account_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on batchApplicationClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMBatchApplication_template(rInt int, rString string, location string, displayName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_batch_account" "test" {
  name                 = "acctestba%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"
}

resource "azurerm_batch_application" "test" {
  name                = "acctestbatchapp-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
	account_name        = "${azurerm_batch_account.test.name}"
	%s
}
`, rInt, location, rString, rString, rInt, displayName)
}
