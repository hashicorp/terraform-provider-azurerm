// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package vmware_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAccVmwareNetappVolumeAttachment_list_basic(t *testing.T) {
	r := VmwareNetappVolumeAttachmentResource{}
	listResourceAddress := "azurerm_vmware_netapp_volume_attachment.list"

	data := acceptance.BuildTestData(t, "azurerm_vmware_netapp_volume_attachment", "test")

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basicList(data),
			},
			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast(listResourceAddress, 2),
				},
			},
		},
	})
}

func (r VmwareNetappVolumeAttachmentResource) basicList(data acceptance.TestData) string {
	return fmt.Sprintf(`

resource "azurerm_resource_group" "test" {
  name     = "acctestrg-vmware-nat-%d"
  location = "centralus"
}

%s

%s

resource "azurerm_vmware_netapp_volume_attachment" "test" {
  count = 2

  name              = "acctest-vmwareattachment${count.index}-%d"
  netapp_volume_id  = azurerm_netapp_volume.test[count.index].id
  vmware_cluster_id = "${azurerm_vmware_private_cloud.test.id}/clusters/Cluster-1"

  depends_on = [azurerm_virtual_network_gateway_connection.test]
}`, data.RandomInteger, r.templatePrivateCloud(data), r.templateNetappFileList(data), data.RandomInteger)
}

func (r VmwareNetappVolumeAttachmentResource) basicQuery() string {
	return `
list "azurerm_vmware_netapp_volume_attachment" "list" {
  provider = azurerm
  config {
    vmware_cluster_id = "${azurerm_vmware_private_cloud.test.id}/clusters/Cluster-1"
  }
}
`
}

func (r VmwareNetappVolumeAttachmentResource) templateNetappFileList(data acceptance.TestData) string {
	return fmt.Sprintf(`
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

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%d"
  location            = "central us"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%d"
  location            = "centralus"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4

  tags = {
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume" "test" {
  count = 2

  name                            = "acctest-NetAppVolume${count.index}-%d"
  location                        = "centralus"
  resource_group_name             = azurerm_resource_group.test.name
  account_name                    = azurerm_netapp_account.test.name
  pool_name                       = azurerm_netapp_pool.test.name
  volume_path                     = "my-unique-file-path${count.index}-%d"
  service_level                   = "Standard"
  subnet_id                       = azurerm_subnet.netappSubnet.id
  protocols                       = ["NFSv3"]
  storage_quota_in_gb             = 100
  azure_vmware_data_store_enabled = true
  snapshot_directory_visible      = true

  export_policy_rule {
    rule_index          = 1
    allowed_clients     = ["0.0.0.0/0"]
    protocols_enabled   = ["NFSv3"]
    unix_read_only      = false
    unix_read_write     = true
    root_access_enabled = true
  }

  tags = {
    "SkipASMAzSecPack" = "true"
  }
}`, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
