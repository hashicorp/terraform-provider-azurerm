package vmware_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type VmwarePrivateCloudDataSource struct {
}

func TestAccDataSourceVmwarePrivateCloud_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_vmware_private_cloud", "test")
	r := VmwarePrivateCloudDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("management_cluster.0.size").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.id").Exists(),
				check.That(data.ResourceName).Key("management_cluster.0.hosts.#").Exists(),
				check.That(data.ResourceName).Key("network_subnet_cidr").Exists(),
				check.That(data.ResourceName).Key("internet_connection_enabled").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.express_route_private_peering_id").Exists(),
				check.That(data.ResourceName).Key("circuit.0.primary_subnet_cidr").Exists(),
				check.That(data.ResourceName).Key("circuit.0.secondary_subnet_cidr").Exists(),
				check.That(data.ResourceName).Key("hcx_cloud_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("management_subnet_cidr").Exists(),
				check.That(data.ResourceName).Key("nsxt_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("nsxt_manager_endpoint").Exists(),
				check.That(data.ResourceName).Key("provisioning_subnet_cidr").Exists(),
				check.That(data.ResourceName).Key("vcenter_certificate_thumbprint").Exists(),
				check.That(data.ResourceName).Key("vcsa_endpoint").Exists(),
				check.That(data.ResourceName).Key("vmotion_subnet_cidr").Exists(),
			),
		},
	})
}

func (VmwarePrivateCloudDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_vmware_private_cloud" "test" {
  name                = azurerm_vmware_private_cloud.test.name
  resource_group_name = azurerm_vmware_private_cloud.test.resource_group_name
}
`, VmwarePrivateCloudResource{}.basic(data))
}
