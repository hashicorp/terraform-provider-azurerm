package azurerm

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMHDInsightClusters_standardWorker(t *testing.T) {
	ri := acctest.RandInt()
	rs := strings.ToLower(acctest.RandString(11))
	config := testAccAzureRMHDInsightClusters_standardWorker(ri, rs, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMApplicationInsightsDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClustersExists("azurerm_hdinsight_cluster.test"),
				),
			},
		},
	})
}

func testCheckAzureRMHDInsightClustersExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for HDInsight: %s", name)
		}

		client := testAccProvider.Meta().(*ArmClient).hdInsightClustersClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on hdInsightClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: HDInsight '%q' (resource group: '%q') does not exist", name, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMHDInsightClusters_standardWorker(rInt int, rString string, location string) string {
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

resource "azurerm_storage_container" "test" {
  name                  = "hdistcontainer"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_hdinsight_cluster" "test" {
  name                = "acctesthdinsight-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  cluster_type        = "hadoop"
  cluster_version     = "3.6"
  os_type             = "Linux"
  login_username   = "super"
  login_password   = "Password1234!"

  storage_account {
    blob_endpoint = "${azurerm_storage_account.test.primary_blob_endpoint}"
    container     = "${azurerm_storage_container.test.name}"
    access_key    = "${azurerm_storage_account.test.primary_access_key}"
  }

  head_node {
    target_instance_count = 2
    vm_size               = "Large"

    linux_profile {
	  username = "super"
	  
      ssh_keys {
        key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAj4v5pRh3Lx0zs8pSYoRQiZdT7PlsTqQJn12J9zh8ebC/Prs2Xyyex8n34k9UC8Q293ALbyE4DZE66aphCU9Dqtb5+LTCK7b/DCGFSaGwDHC/jej2YP3UBGbiBKBFtVrOSFLzul8d9r+ssjdJw+u6wBKpF+fIt9O2eHlOjAHuhuEMQnTnqdQpNsMq5Jjo/XzAf/yxcDhVLUUN9kLTuHpbvW6UHxYT1ejx+f6+WTk8p5lfW2J7B/qAbdIF4823/lCcTd3RfRmmRY8MkK4RtDAWZxHfqtkct04ZVoaTVZh5qDFaYnhsgoTJ2rut7VsUF3Q+gMTlKNk6ES4XGTZIUJFfHQ=="
      }
    }
  }

  worker_node {
    target_instance_count = 2
    vm_size               = "Large"

    linux_profile {
      username = "super"

      ssh_keys {
        key_data = "ssh-rsa AAAAB3NzaC1yc2EAAAABJQAAAQEAj4v5pRh3Lx0zs8pSYoRQiZdT7PlsTqQJn12J9zh8ebC/Prs2Xyyex8n34k9UC8Q293ALbyE4DZE66aphCU9Dqtb5+LTCK7b/DCGFSaGwDHC/jej2YP3UBGbiBKBFtVrOSFLzul8d9r+ssjdJw+u6wBKpF+fIt9O2eHlOjAHuhuEMQnTnqdQpNsMq5Jjo/XzAf/yxcDhVLUUN9kLTuHpbvW6UHxYT1ejx+f6+WTk8p5lfW2J7B/qAbdIF4823/lCcTd3RfRmmRY8MkK4RtDAWZxHfqtkct04ZVoaTVZh5qDFaYnhsgoTJ2rut7VsUF3Q+gMTlKNk6ES4XGTZIUJFfHQ=="
      }
    }
  }
}
`, rInt, location, rString, rInt)
}
