package avs_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AvsPrivateCloudDataSource struct {
}

func TestAccDataSourceAvsPrivateCloud_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_avs_private_cloud", "test")
	r := AvsPrivateCloudDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("management_cluster.0.cluster_size").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.cluster_id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("network_block").Exists(),
				check.That(data.ResourceName).Key("internet_connected").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_network").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_network").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_network").Exists(),
			),
		},
	})
}

func (AvsPrivateCloudDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_avs_private_cloud" "test" {
  name                = azurerm_avs_private_cloud.test.name
  resource_group_name = azurerm_avs_private_cloud.test.resource_group_name
}
`, AvsPrivateCloudResource{}.basic(data))
}
