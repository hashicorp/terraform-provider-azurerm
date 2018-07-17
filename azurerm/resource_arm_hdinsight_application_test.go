package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMHDInsightApplication_basic(t *testing.T) {
	resourceName := "azurerm_hdinsight_application.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(5)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightApplication_basic(rInt, rString, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightApplicationExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "emptynodeapp"),
					resource.TestCheckResourceAttr(resourceName, "marketplace_identifier", "EmptyNode"),
					resource.TestCheckResourceAttr(resourceName, "edge_node.0.target_instance_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "install_script_action.0.name", "emptynode-sayhello"),
				),
			},
		},
	})
}
func testCheckAzureRMHDInsightApplicationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).hdinsightApplicationsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hdinsight_application" {
			continue
		}

		applicationName := rs.Primary.Attributes["name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, clusterName, applicationName)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("HDInsight Application %q (Cluster %q) still exists in resource group %q", applicationName, clusterName, resourceGroup)
			}

			return nil
		}
	}
	return nil
}

func testCheckAzureRMHDInsightApplicationExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		applicationName := rs.Primary.Attributes["name"]
		clusterName := rs.Primary.Attributes["cluster_name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for HDInsight Application: %q", applicationName)
		}

		client := testAccProvider.Meta().(*ArmClient).hdinsightApplicationsClient
		resp, err := client.Get(ctx, resourceGroup, clusterName, applicationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: HDInsight Application %q (Cluster %q / Resource Group: %q) does not exist", applicationName, clusterName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on hdinsightApplicationsClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMHDInsightApplication_basic(rInt int, rString string, location string) string {
	template := testAccAzureRMHDInsightCluster_basicConfig(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_application" "test" {
  name                   = "new-edgenode"
  cluster_name           = "${azurerm_hdinsight_cluster.test.name}"
  resource_group_name    = "${azurerm_resource_group.test.name}"
  marketplace_identifier = "EmptyNode"

  edge_node {
    hardware_profile {
      vm_size = "Standard_D3_v2"
    }
  }

  install_script_action {
    name  = "emptynode-sayhello"
    uri   = "https://gist.githubusercontent.com/tombuildsstuff/74ff75620a83cf2a737843920185dbc2/raw/8217fbbcf9728e23807c19a35f65136351e6da7a/hello.sh"
    roles = [ "edgenode" ]
  }
}
`, template)
}
