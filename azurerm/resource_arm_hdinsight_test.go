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
	rStr := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	config := testAccAzureRMHDInsightCluster_basic(ri, rStr, testLocation())

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

func testAccAzureRMHDInsightCluster_basic(rInt int, rString, location string) string {
	return fmt.Sprintf(`
		resource "azurerm_resource_group" "test" {
			name     = "acctestrg-%d"
			location = "%s"
		  
			tags {
			  Source = "testAccAzureRMHDInsightCluster_basic"
			}
		  }
		  
		  resource "azurerm_storage_account" "test" {
			  name	= "acctestsa%s"
			  resource_group_name = "${azurerm_resource_group.test.name}"
			  location = "${azurerm_resource_group.test.location}"
			  account_tier = "Standard"
			  account_replication_type = "GRS"
		  }
		  
		  resource "azurerm_hdinsight_cluster" "test" {
			name                = "acctesthdi%d"
			resource_group_name = "${azurerm_resource_group.test.name}"
			location            = "${azurerm_resource_group.test.location}"
			cluster_version  = "3.6"
		  
			cluster_definition {
				  kind = "spark"
				  configurations {
					  gateway {
						  rest_auth_credential_is_enabled = true
						  rest_auth_credential_username = "http-user"
						  rest_auth_credential_password = "AbcAbc123123!"
					  }
				  }
			}
		  
			compute_profile {
				roles = [
					{
						name = "headnode"
						target_instance_count = 2,
						hardware_profile {
							vm_size = "Standard_D12_v2"
						},
						os_profile {
							linux_operating_system_profile {
								username = "username"
								password = "testPass123!"
								ssh_key {
									key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
								}
							}
						}
					},
					{
						name = "workernode"
						hardware_profile {
							vm_size = "Standard_D13_v2"
						},
						os_profile {
							linux_operating_system_profile {
								username = "username"
								ssh_key {
									key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCqaZoyiz1qbdOQ8xEf6uEu1cCwYowo5FHtsBhqLoDnnp7KUTEBN+L2NxRIfQ781rxV6Iq5jSav6b2Q8z5KiseOlvKA/RF2wqU0UPYqQviQhLmW6THTpmrv/YkUCuzxDpsH7DUDhZcwySLKVVe0Qm3+5N2Ta6UYH3lsDf9R9wTP2K/+vAnflKebuypNlmocIvakFWoZda18FOmsOoIVXQ8HWFNCuw9ZCunMSN62QGamCe3dL5cXlkgHYv7ekJE15IA9aOJcM7e90oeTqo+7HTcWfdu0qQqPWY5ujyMw/llas8tsXY85LFqRnr3gJ02bAscjc477+X+j/gkpFoN1QEmt terraform@demo.tld"
								}
							}
						}
					}
				]
			}

			storage_profile {
				storage_accounts = [
					{
						name = "${azurerm_storage_account.test.primary_blob_endpoint}"
						is_default = true
						container = "${azurerm_resource_group.test.name}"
						key = "${azurerm_storage_account.test.primary_access_key}"
					}
				]
			}
		}
`, rInt, location, rString, rInt)
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

		client := testAccProvider.Meta().(*ArmClient).hdInsightClustersClient
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
	conn := testAccProvider.Meta().(*ArmClient).hdInsightClustersClient

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
