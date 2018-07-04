package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestHDInsightsClusterValidation(t *testing.T) {
	data := map[string]bool{
		"":           false,
		"3":          false,
		"3.":         false,
		"4.6":        true,
		"3.2.1":      false,
		"whowhatnow": false,
	}

	for val, expected := range data {
		_, errors := validateHDInsightsClusterVersion(val, "test_example")
		result := len(errors) == 0
		if expected != result {
			t.Fatalf("Expected %q to return %t but returned %t", val, expected, result)
		}
	}
}

func TestAccAzureRMHDInsightCluster_basic(t *testing.T) {
	resourceName := "azurerm_hdinsight_cluster.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(5)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightCluster_basicConfig(rInt, rString, location, "hadoop", 3),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.kind", "hadoop"),
					resource.TestCheckResourceAttr(resourceName, "worker_node.0.target_instance_count", "3"),
				),
			},
		},
	})
}

func TestAccAzureRMHDInsightCluster_updateWorkerCount(t *testing.T) {
	resourceName := "azurerm_hdinsight_cluster.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(5)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightCluster_basicConfig(rInt, rString, location, "hadoop", 4),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.kind", "hadoop"),
					resource.TestCheckResourceAttr(resourceName, "worker_node.0.target_instance_count", "4"),
				),
			},
			{
				Config: testAccAzureRMHDInsightCluster_basicConfig(rInt, rString, location, "hadoop", 5),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.kind", "hadoop"),
					resource.TestCheckResourceAttr(resourceName, "worker_node.0.target_instance_count", "5"),
				),
			},
		},
	})
}

func TestAccAzureRMHDInsightCluster_basicNetwork(t *testing.T) {
	resourceName := "azurerm_hdinsight_cluster.test"
	rInt := acctest.RandInt()
	rString := acctest.RandString(5)
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMHDInsightClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMHDInsightCluster_basicNetworkConfig(rInt, rString, location, "spark"),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMHDInsightClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.Source", "AcceptanceTest"),
					resource.TestCheckResourceAttr(resourceName, "tier", "standard"),
					resource.TestCheckResourceAttr(resourceName, "cluster.0.kind", "spark"),
					resource.TestCheckResourceAttrSet(resourceName, "head_node.0.virtual_network_profile.0.virtual_network_id"),
					resource.TestCheckResourceAttrSet(resourceName, "head_node.0.virtual_network_profile.0.subnet_id"),
				),
			},
		},
	})
}

func testCheckAzureRMHDInsightClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).hdinsightClustersClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_hdinsight_cluster" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("HDInsight Cluster %q still exists in resource group %q", name, resourceGroup)
			}

			return nil
		}
	}
	return nil
}

func testCheckAzureRMHDInsightClusterExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		clusterName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for HDInsight Cluster: %q", clusterName)
		}

		client := testAccProvider.Meta().(*ArmClient).hdinsightClustersClient
		resp, err := client.Get(ctx, resourceGroup, clusterName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: HDInsight Cluster %q (resource group: %q) does not exist", clusterName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on hdinsightClustersClient: %+v", err)
		}

		return nil

	}
}

func testAccAzureRMHDInsightCluster_basicConfig(rInt int, rString string, location string, clusterType string, nodes int) string {
	template := testAccAzureRMHDInsightCluster_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_hdinsight_cluster" "test" {
  name                = "acctesthdic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tier                = "standard"

  cluster {
    kind    = "%s"
    version = "3.6"

    gateway {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  storage_profile {
    storage_account {
      storage_account_name = "${azurerm_storage_account.test.primary_blob_domain}"
      storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
      container_name       = "${azurerm_storage_container.test.name}"
      is_default           = true
    }
  }

  head_node {
    target_instance_count = 2

    hardware_profile {
      vm_size = "Standard_D3_V2"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  worker_node {
    target_instance_count = %d

    hardware_profile {
      vm_size = "Medium"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  zookeeper_node {
    target_instance_count = 3

    hardware_profile {
      vm_size = "A5"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }
  }
}
`, template, rInt, clusterType, nodes)
}

func testAccAzureRMHDInsightCluster_basicNetworkConfig(rInt int, rString string, location string, clusterType string) string {
	template := testAccAzureRMHDInsightCluster_template(rInt, rString, location)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_hdinsight_cluster" "test" {
  name                = "acctesthdic-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tier                = "standard"

  cluster {
    kind    = "%s"
    version = "3.6"

    gateway {
      username = "${var.username}"
      password = "${var.password}"
    }
  }

  storage_profile {
    storage_account {
      storage_account_name = "${azurerm_storage_account.test.primary_blob_domain}"
      storage_account_key  = "${azurerm_storage_account.test.primary_access_key}"
      container_name       = "${azurerm_storage_container.test.name}"
      is_default           = true
    }
  }

  head_node {
    target_instance_count = 2

    hardware_profile {
      vm_size = "Standard_D3_V2"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }

    virtual_network_profile {
      virtual_network_id = "${azurerm_virtual_network.test.id}"
      subnet_id          = "${azurerm_subnet.test.id}"
    }
  }

  worker_node {
    target_instance_count = 4

    hardware_profile {
      vm_size = "Medium"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }

    virtual_network_profile {
      virtual_network_id = "${azurerm_virtual_network.test.id}"
      subnet_id          = "${azurerm_subnet.test.id}"
    }
  }

  zookeeper_node {
    target_instance_count = 3

    hardware_profile {
      vm_size = "A5"
    }

    os_profile {
      username = "${var.username}"
      password = "${var.password}"
    }

    virtual_network_profile {
      virtual_network_id = "${azurerm_virtual_network.test.id}"
      subnet_id          = "${azurerm_subnet.test.id}"
    }
  }

  tags {
    "Source" = "AcceptanceTest"
  }
}
`, template, rInt, rInt, clusterType)
}

func testAccAzureRMHDInsightCluster_template(rInt int, rString string, location string) string {
	return fmt.Sprintf(`
variable "username" {
  default = "tfadmin"
}

variable "password" {
  default = "Password21!$"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  location                 = "${azurerm_resource_group.test.location}"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "data"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}
`, rInt, location, rString)
}
