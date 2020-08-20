package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMavsPrivateCloud_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_avs_private_cloud", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMavsPrivateCloudDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceavsPrivateCloud_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMavsPrivateCloudExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.cluster_size"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.cluster_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_cluster.0.hosts.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_block"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "internet_connected"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.express_route_private_peering_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.primary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "circuit.0.secondary_subnet"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "hcx_cloud_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "nsxt_manager_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "provisioning_network"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcenter_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vcsa_endpoint"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "vmotion_network"),
				),
			},
		},
	})
}

func testAccDataSourceavsPrivateCloud_basic(data acceptance.TestData) string {
	config := testAccAzureRMavsPrivateCloud_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_avs_private_cloud" "test" {
  name                = azurerm_avs_private_cloud.test.name
  resource_group_name = azurerm_avs_private_cloud.test.resource_group_name
}
`, config)
}
