// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-01-01/buckets"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type NetAppVolumeBucketResource struct{}

func TestAccNetAppVolumeBucket_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucket_cifsUser(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cifsUser(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucket_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucket_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("server.0.certificate_pem"),
	})
}

func TestAccNetAppVolumeBucket_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_bucket", "test")
	r := NetAppVolumeBucketResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (t NetAppVolumeBucketResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := buckets.ParseBucketID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.BucketsClient.Get(ctx, *id)
	if err != nil {
		if resp.HttpResponse != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (NetAppVolumeBucketResource) basic(data acceptance.TestData) string {
	template := NetAppVolumeBucketResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "test" {
  name      = "acctest-bucket-%[2]d"
  volume_id = azurerm_netapp_volume.test.id

  file_system_user {
    nfs_user {
      group_id = 1000
      user_id  = 1000
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeBucketResource) cifsUser(data acceptance.TestData) string {
	template := NetAppVolumeBucketResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "test" {
  name      = "acctest-bucket-cifs-%[2]d"
  volume_id = azurerm_netapp_volume.test.id

  file_system_user {
    cifs_user {
      username = "anfuser"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeBucketResource) complete(data acceptance.TestData) string {
	template := NetAppVolumeBucketResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "test" {
  name        = "acctest-bucket-%[2]d"
  volume_id   = azurerm_netapp_volume.test.id
  path        = "/"
  permissions = "ReadWrite"

  file_system_user {
    nfs_user {
      group_id = 2000
      user_id  = 2000
    }
  }

  server {
    fqdn = "anf-bucket-%[2]d.example.com"
  }
}
`, template, data.RandomInteger)
}

func (r NetAppVolumeBucketResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_bucket" "import" {
  name      = azurerm_netapp_volume_bucket.test.name
  volume_id = azurerm_netapp_volume_bucket.test.volume_id

  file_system_user {
    nfs_user {
      group_id = 1000
      user_id  = 1000
    }
  }
}
`, r.basic(data))
}

func (NetAppVolumeBucketResource) template(data acceptance.TestData) string {
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
    "CreatedOnDate"    = "2026-01-15T00-00-00Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.99.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-DelegatedSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.99.2.0/24"]

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
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Auto"
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
}
`, data.RandomInteger, data.Locations.Primary)
}
