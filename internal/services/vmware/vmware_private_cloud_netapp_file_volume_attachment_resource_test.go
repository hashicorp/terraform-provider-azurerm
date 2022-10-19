package vmware_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/vmware/2021-12-01/datastores"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"testing"
)

type VmwareNetappFileVolumeAttachmentResource struct{}

func (r VmwareNetappFileVolumeAttachmentResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := datastores.ParseDataStoreID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Vmware.DataStoreClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return utils.Bool(true), nil
}

func TestAccVmwarePrivateCloudNetappFileVolumeAttachment(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_vmware_netapp_volume_attachment", "test")
	r := VmwareNetappFileVolumeAttachmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func (r VmwareNetappFileVolumeAttachmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_vmware_netapp_volume_attachment" "test" {
  name              = "acctest-vmwareattachment-%d"
  netapp_volume_id  = azurerm_netapp_volume.test.id
  vmware_cluster_id = azurerm_vmware_cluster.test.id
}`, r.templatePrivateCloud(data), data.RandomInteger)
}

func (r VmwareNetappFileVolumeAttachmentResource) templatePrivateCloud(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
  # In Vmware acctest, please disable correlation request id, else the continuous operations like update or delete will not be triggered
  # issue https://github.com/Azure/azure-rest-api-specs/issues/14086 
  disable_correlation_request_id = true
}

%s
resource "azurerm_vmware_private_cloud" "test" {
  name                = "acctest-PC-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "av36"

  management_cluster {
    size = 3
  }
  network_subnet_cidr = "192.168.48.0/22"
}

resource "azurerm_vmware_cluster" "test" {
  name               = "acctest-Cluster-%d"
  vmware_cloud_id    = azurerm_vmware_private_cloud.test.id
  cluster_node_count = 3
  sku_name           = "av36"
}

resource "azurerm_vmware_express_route_authorization" "test" {
  name             = "acctest-VmwareAuthorization-%d"
  private_cloud_id = azurerm_vmware_private_cloud.test.id
}

resource "azurerm_virtual_network_gateway_connection" "test" {
  name                = "acctestvnetgwconn-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type                       = "ExpressRoute"
  virtual_network_gateway_id = azurerm_virtual_network_gateway.test.id
  express_route_circuit_id   = azurerm_vmware_private_cloud.test.express_route_id
  authorization_key          = azurerm_vmware_express_route_authorization.test.express_route_authorization_key
}
`, r.templateVnet(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VmwareNetappFileVolumeAttachmentResource) templateNetappFile(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume" "test" {
  name                            = "acctest-NetAppVolume-%d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  account_name                    = azurerm_netapp_account.test.name
  pool_name                       = azurerm_netapp_pool.test.name
  volume_path                     = "my-unique-file-path-%d"
  service_level                   = "Standard"
  subnet_id                       = azurerm_subnet.netappSubnet.id
  protocols                       = ["NFSv3"]
  storage_quota_in_gb             = 100
  azure_vmware_data_store_enabled = true

  export_policy_rule {
    rule_index          = 1
    allowed_clients     = ["0.0.0.0/0"]
    protocols_enabled   = ["NFSv3"]
    unix_read_only      = false
    unix_read_write     = true
    root_access_enabled = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}`, r.templateVnet(data), data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VmwareNetappFileVolumeAttachmentResource) templateVnet(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.6.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "netappSubnet" {
  name                 = "acctest-Subnet-%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.6.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "gatewaySubnet" {
  name                 = "GatewaySubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.6.1.0/24"]
}

resource "azurerm_virtual_network_gateway" "test" {
  name                = "acctestvnetgw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  type = "ExpressRoute"
  sku  = "Standard"

  ip_configuration {
    name                 = "vnetGatewayConfig"
    public_ip_address_id = azurerm_public_ip.test.id
    subnet_id            = azurerm_subnet.gatewaySubnet.id
  }
}
`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
