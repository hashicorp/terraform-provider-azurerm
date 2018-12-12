package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMBatchPool_basicFixedScale(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()

	config := testaccAzureRMBatchPool_basicFixedScale(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchAccountDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "scale_mode", "Fixed"),
					resource.TestCheckResourceAttr(resourceName, "target_dedicated_nodes", "2"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
				),
			},
		},
	})
}

func TestAccAzureRMBatchPool_basicAutoScale(t *testing.T) {
	resourceName := "azurerm_batch_pool.test"
	ri := acctest.RandInt()
	rs := acctest.RandString(4)
	location := testLocation()

	config := testaccAzureRMBatchPool_basicAutoScale(ri, rs, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMBatchAccountDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMBatchPoolExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "vm_size", "STANDARD_A1"),
					resource.TestCheckResourceAttr(resourceName, "scale_mode", "Auto"),
					resource.TestCheckResourceAttr(resourceName, "node_agent_sku_id", "batch.node.ubuntu 16.04"),
					resource.TestCheckResourceAttr(resourceName, "account_name", fmt.Sprintf("testaccbatch%s", rs)),
				),
			},
		},
	})
}

func testCheckAzureRMBatchPoolExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		poolName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["account_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).batchPoolClient

		resp, err := conn.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			return fmt.Errorf("Bad: Get on batchPoolClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Batch pool %q (account: %q, resource group: %q) does not exist", poolName, accountName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMBatchPoolDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_batch_account" {
			continue
		}

		poolName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["account_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		conn := testAccProvider.Meta().(*ArmClient).batchPoolClient

		resp, err := conn.Get(ctx, resourceGroup, accountName, poolName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testaccAzureRMBatchPool_basicFixedScale(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "testaccsa%s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
  }

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
	name                   = "testaccpool%s"
	resource_group_name    = "${azurerm_resource_group.test.name}"
	account_name 		   = "${azurerm_batch_account.test.name}"
	display_name		   = "Test Acc Pool"
	vm_size				   = "Standard_A1"
	scale_mode			   = "Fixed"
	target_dedicated_nodes = 2
	node_agent_sku_id	= "batch.node.ubuntu 16.04"

	storage_image_reference {
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "16.04.0-LTS"
        version   = "latest"
	}
  }
`, rInt, location, rString, rString, rString)
}

func testaccAzureRMBatchPool_basicAutoScale(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "testaccbatch%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
	name                     = "testaccsa%s"
	resource_group_name      = "${azurerm_resource_group.test.name}"
	location                 = "${azurerm_resource_group.test.location}"
	account_tier             = "Standard"
	account_replication_type = "LRS"
  }

resource "azurerm_batch_account" "test" {
  name                 = "testaccbatch%s"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  pool_allocation_mode = "BatchService"
  storage_account_id   = "${azurerm_storage_account.test.id}"

  tags {
    env = "test"
  }
}

resource "azurerm_batch_pool" "test" {
	name                   		  = "testaccpool%s"
	resource_group_name           = "${azurerm_resource_group.test.name}"
	account_name 		   		  = "${azurerm_batch_account.test.name}"
	display_name		   		  = "Test Acc Pool Auto"
	vm_size				   		  = "Standard_A1"
	scale_mode			   		  = "Auto"
	autoscale_evaluation_interval = "PT15M"
	autoscale_formula			  = <<EOF
	startingNumberOfVMs = 1;
	maxNumberofVMs = 25;
	pendingTaskSamplePercent = $PendingTasks.GetSamplePercent(180 * TimeInterval_Second);
	pendingTaskSamples = pendingTaskSamplePercent < 70 ? startingNumberOfVMs : avg($PendingTasks.GetSample(180 * TimeInterval_Second));
	$TargetDedicatedNodes=min(maxNumberofVMs, pendingTaskSamples);
	EOF
	node_agent_sku_id			  = "batch.node.ubuntu 16.04"

	storage_image_reference {
        publisher = "Canonical"
        offer     = "UbuntuServer"
        sku       = "16.04.0-LTS"
        version   = "latest"
	}
  }
`, rInt, location, rString, rString, rString)
}
