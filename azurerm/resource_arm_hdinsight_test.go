package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMHDInsightCluster_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMHDInsightCluster_basic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists("azurerm_hdinsight_cluster.test"),
				),
			},
		},
	})
}

func testAccAzureRMHDInsightCluster_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
	name	= "acctesthdinsight%d"
	resource_group_name = "${azurerm_resource_group.test.name}"
	location = "${azurerm_resource_group.test.location}"
	account_tier = "Standard"
	account_replication_type = "GRS"
}

resource "azurerm_storage_container" "test" {
	name                  = "acctesthdinsight%d"
	resource_group_name   = "${azurerm_resource_group.test.name}"
	storage_account_name  = "${azurerm_storage_account.test.name}"
	container_access_type = "public"
}

resource "azurerm_hdinsight_cluster" "test" {
  name                = "acctesthdinsight%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  cluster_version  = "3.6"
  os_type  = "Linux"
  tier  = "Standard"

  cluster_definition {
		kind = "hbase"
  }

  "storage_profile": {
    storage_accounts: [
			{
				"name": "${azurerm_storage_account.test.name}",
				"is_default": true,
				"container": "${azurerm_storage_container.test.name}",
				"key": "${azurerm_storage_account.test.primary_access_key}"
			}
		]
  }

  "compute_profile": {
	  "roles": [
		  {
			  "name": "headnode",
			  "target_instance_count": 2,
			  "hardware_profile": {
				  "vm_size": "Small"
			  },
			  "os_profile": {
				  linux_operating_system_profile": {
					  "username": "testhbaselogin",
					  "password": "testhbasepass123"
				  }
			  }
		  },
		  {
			  "name": "workernode",
			  "target_instance_count": 2,
			  "hardware_profile": {
				  "vm_size": "Small"
			  },
			  "os_profile": {
				  "linux_operation_system_profile": {
					"username": "testhbaselogin",
					"password": "testhbasepass123"
				  }
			  }
		  },
		  {
			  "name": "zookeepernode",
			  "target_instance_count": 3,
			  "hardware_profile": {
				  "vm_size": "Small"
			  },
			  "os_profile": {
				  "linux_operation_system_profile": {
					"username": "testhbaselogin",
					"password": "testhbasepass123"
				  }
			  }
		  }
	  ]
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testCheckAzureRMHDInsightClusterExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		resourceName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for HDInsight cluster instance: %s", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).hdinsightClustersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		hdinsight, err := client.Get(ctx, resourceGroup, resourceName)
		if err != nil {
			return fmt.Errorf("Bad: Get on hdinsightClustersClient: %+v", err)
		}

		if hdinsight.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: HDInsight cluster instance %q (resource group: %q) does not exist", resourceName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMHDInsightClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).hdinsightClustersClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hdinsight_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := conn.Get(ctx, resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("HDInsight cluster instance still exists:\n%#v", resp)
		}
	}

	return nil
}
