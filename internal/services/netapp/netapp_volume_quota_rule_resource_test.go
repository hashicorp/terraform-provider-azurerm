// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumequotarules"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeQuotaRuleResource struct{}

func TestAccNetAppVolumeQuotaRule_individualGroupQuotaType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_quota_rule", "test")
	r := NetAppVolumeQuotaRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.individualGroupQuotaType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeQuotaRule_individualUserQuotaType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_quota_rule", "test")
	r := NetAppVolumeQuotaRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.individualUserQuotaType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeQuotaRule_defaultGroupQuotaType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_quota_rule", "test")
	r := NetAppVolumeQuotaRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultGroupQuotaType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeQuotaRule_defaultUserQuotaType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_quota_rule", "test")
	r := NetAppVolumeQuotaRuleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultUserQuotaType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppVolumeQuotaRuleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := volumequotarules.ParseVolumeQuotaRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.VolumeQuotaRules.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse.StatusCode == http.StatusNotFound {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (NetAppVolumeQuotaRuleResource) individualGroupQuotaType(data acceptance.TestData) string {
	template := NetAppVolumeQuotaRuleResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_quota_rule" "test" {
  name              = "acctest-NetAppVolumeQuotaRule-%[2]d"
  location          = azurerm_resource_group.test.location
  volume_id         = azurerm_netapp_volume.test.id
  quota_target      = "3001"
  quota_size_in_kib = 1024
  quota_type        = "IndividualGroupQuota"
}
`, template, data.RandomInteger)
}

func (NetAppVolumeQuotaRuleResource) individualUserQuotaType(data acceptance.TestData) string {
	template := NetAppVolumeQuotaRuleResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_quota_rule" "test" {
  name              = "acctest-NetAppVolumeQuotaRule-%[2]d"
  location          = azurerm_resource_group.test.location
  volume_id         = azurerm_netapp_volume.test.id
  quota_target      = "3001"
  quota_size_in_kib = 1024
  quota_type        = "IndividualUserQuota"
}
`, template, data.RandomInteger)
}

func (NetAppVolumeQuotaRuleResource) defaultUserQuotaType(data acceptance.TestData) string {
	template := NetAppVolumeQuotaRuleResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_quota_rule" "test" {
  name              = "acctest-NetAppVolumeQuotaRule-Default-Usr-%[2]d"
  location          = azurerm_resource_group.test.location
  volume_id         = azurerm_netapp_volume.test.id
  quota_size_in_kib = 2048
  quota_type        = "DefaultUserQuota"
}
`, template, data.RandomInteger)
}

func (NetAppVolumeQuotaRuleResource) defaultGroupQuotaType(data acceptance.TestData) string {
	template := NetAppVolumeQuotaRuleResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_quota_rule" "test" {
  name              = "acctest-NetAppVolumeQuotaRule-Default-Grp-%[2]d"
  location          = azurerm_resource_group.test.location
  volume_id         = azurerm_netapp_volume.test.id
  quota_size_in_kib = 2048
  quota_type        = "DefaultGroupQuota"
}
`, template, data.RandomInteger)
}

func (NetAppVolumeQuotaRuleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      prevent_volume_destruction             = false
      delete_backups_on_backup_vault_destroy = true
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-DelegatedSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.88.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Auto"

  tags = {
    "CreatedOnDate"    = "2023-08-17T08:01:00Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[1]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  protocols           = ["NFSv3"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
