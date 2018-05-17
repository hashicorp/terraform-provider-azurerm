package azurerm

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMKubernetesCluster_basic(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_basic(ri, clientId, clientSecret, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.client_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.client_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.cluster_ca_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.host"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.username"),
					resource.TestCheckResourceAttrSet(dataSourceName, "kube_config.0.password"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMKubernetesCluster_internalNetwork(t *testing.T) {
	dataSourceName := "data.azurerm_kubernetes_cluster.test"
	ri := acctest.RandInt()
	clientId := os.Getenv("ARM_CLIENT_ID")
	clientSecret := os.Getenv("ARM_CLIENT_SECRET")
	location := testLocation()
	config := testAccDataSourceAzureRMKubernetesCluster_internalNetwork(ri, clientId, clientSecret, location)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMKubernetesClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMKubernetesClusterExists(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_pool_profile.0.vnet_subnet_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMKubernetesCluster_basic(rInt int, clientId string, clientSecret string, location string) string {
	resource := testAccAzureRMKubernetesCluster_basic(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}

func testAccDataSourceAzureRMKubernetesCluster_internalNetwork(rInt int, clientId string, clientSecret string, location string) string {
	resource := testAccAzureRMKubernetesCluster_internalNetwork(rInt, clientId, clientSecret, location)
	return fmt.Sprintf(`
%s

data "azurerm_kubernetes_cluster" "test" {
  name                = "${azurerm_kubernetes_cluster.test.name}"
  resource_group_name = "${azurerm_kubernetes_cluster.test.resource_group_name}"
}
`, resource)
}
