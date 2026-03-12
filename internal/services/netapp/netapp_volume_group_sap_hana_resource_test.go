// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-06-01/volumegroups"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

type NetAppVolumeGroupSAPHanaResource struct{}

func TestAccNetAppVolumeGroupSAPHana_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_basicAvailabilityZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAvailabilityZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_volCustomerManagedKeyEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.volEncryptionCmkSAPHana(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_backupVolumeSpecsNfsv3(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupVolumeSpecsNfsv3(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_snapshotPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.avgSnapshotPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_snapshotPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.avgSnapshotPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateAvgSnapshotPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_volumeUpdates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVolumes(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume.0.storage_quota_in_gb").HasValue("1200"),
				check.That(data.ResourceName).Key("volume.1.export_policy_rule.0.allowed_clients").HasValue("10.0.0.0/8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_protocolConversion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.protocolConversionNFSv3(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume.0.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.0.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.1.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.1.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.2.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.2.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.3.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.3.protocols.0").HasValue("NFSv3"),
				check.That(data.ResourceName).Key("volume.4.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.4.protocols.0").HasValue("NFSv3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.protocolConversionNFSv41(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("volume.0.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.0.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.1.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.1.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.2.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.2.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.3.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.3.protocols.0").HasValue("NFSv4.1"),
				check.That(data.ResourceName).Key("volume.4.protocols.#").HasValue("1"),
				check.That(data.ResourceName).Key("volume.4.protocols.0").HasValue("NFSv4.1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_crossRegionReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test_secondary")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionReplication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_crossRegionReplicationAZ(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test_secondary")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionReplicationAZ(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupSAPHana_crossZoneReplicationAZ(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_sap_hana", "test_secondary")
	r := NetAppVolumeGroupSAPHanaResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossZoneReplicationAZ(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppVolumeGroupSAPHanaResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := volumegroups.ParseVolumeGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.VolumeGroupClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(true), nil
}

func (NetAppVolumeGroupSAPHanaResource) basic(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templatePPG(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-1-%[2]d"
    volume_path                  = "my-unique-file-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-%[2]d"
    volume_path                  = "my-unique-file-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-3-%[2]d"
    volume_path                  = "my-unique-file-path-3-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_linux_virtual_machine.test,
    azurerm_proximity_placement_group.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) backupVolumeSpecsNfsv3(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templatePPG(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-1-%[2]d"
    volume_path                  = "my-unique-file-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-%[2]d"
    volume_path                  = "my-unique-file-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-4-%[2]d"
    volume_path                = "my-unique-file-path-4-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    volume_spec_name           = "data-backup"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv3"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = true
      nfsv41_enabled      = false
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-5-%[2]d"
    volume_path                = "my-unique-file-path-5-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    volume_spec_name           = "log-backup"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv3"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = true
      nfsv41_enabled      = false
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-6-%[2]d"
    volume_path                  = "my-unique-file-path-6-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_linux_virtual_machine.test,
    azurerm_proximity_placement_group.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) avgSnapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateAvailabilityZone(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [15, 30]
    hour              = 23
    minute            = 30
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-1-%[2]d"
    volume_path                = "my-unique-file-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-2-%[2]d"
    volume_path                = "my-unique-file-path-2-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-3-%[2]d"
    volume_path                = "my-unique-file-path-3-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_netapp_snapshot_policy.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) updateAvgSnapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateAvailabilityZone(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_snapshot_policy" "test" {
  name                = "acctest-NetAppSnapshotPolicy-New-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  enabled             = true

  monthly_schedule {
    snapshots_to_keep = 3
    days_of_month     = [10, 25]
    hour              = 23
    minute            = 30
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-1-%[2]d"
    volume_path                = "my-unique-file-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-2-%[2]d"
    volume_path                = "my-unique-file-path-2-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-3-%[2]d"
    volume_path                = "my-unique-file-path-3-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_netapp_snapshot_policy.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) updateVolumes(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templatePPG(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-1-%[2]d"
    volume_path                  = "my-unique-file-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1200
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-%[2]d"
    volume_path                  = "my-unique-file-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "10.0.0.0/8"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-3-%[2]d"
    volume_path                  = "my-unique-file-path-3-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_linux_virtual_machine.test,
    azurerm_proximity_placement_group.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) crossRegionReplication(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateForAvgCrossRegionReplication(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test_primary" {
  name                   = "acctest-NetAppVolumeGroup-Primary-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-1-Primary-%[2]d"
    volume_path                  = "my-unique-file-path-1-Primary-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-Primary-%[2]d"
    volume_path                  = "my-unique-file-path-2-Primary-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-3-Primary-%[2]d"
    volume_path                  = "my-unique-file-path-3-Primary-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_linux_virtual_machine.test,
    azurerm_proximity_placement_group.test
  ]
}

resource "azurerm_netapp_volume_group_sap_hana" "test_secondary" {
  name                   = "acctest-NetAppVolumeGroup-Secondary-%[2]d"
  location               = "%[3]s"
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test_secondary.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-1-Secondary-%[2]d"
    volume_path                  = "my-unique-file-path-1-Secondary-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test_secondary.id
    subnet_id                    = azurerm_subnet.test_secondary.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test_secondary.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_netapp_volume_group_sap_hana.test_primary.location
      remote_volume_resource_id = azurerm_netapp_volume_group_sap_hana.test_primary.volume[0].id
      replication_frequency     = "daily"
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-Secondary-%[2]d"
    volume_path                  = "my-unique-file-path-2-Secondary-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test_secondary.id
    subnet_id                    = azurerm_subnet.test_secondary.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test_secondary.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-3-Secondary-%[2]d"
    volume_path                  = "my-unique-file-path-3-Secondary-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test_secondary.id
    subnet_id                    = azurerm_subnet.test_secondary.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test_secondary.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_netapp_volume_group_sap_hana.test_primary.location
      remote_volume_resource_id = azurerm_netapp_volume_group_sap_hana.test_primary.volume[2].id
      replication_frequency     = "daily"
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_linux_virtual_machine.test_secondary,
    azurerm_proximity_placement_group.test_secondary,
  ]
}


`, template, data.RandomInteger, data.Locations.Secondary)
}

func (NetAppVolumeGroupSAPHanaResource) crossRegionReplicationAZ(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateForAvgCrossRegionReplicationAZ(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test_primary" {
  name                   = "acctest-NetAppVolumeGroup-Primary-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-1-Primary-%[2]d"
    volume_path                = "my-unique-file-path-1-Primary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-2-Primary-%[2]d"
    volume_path                = "my-unique-file-path-2-Primary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-3-Primary-%[2]d"
    volume_path                = "my-unique-file-path-3-Primary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}

resource "azurerm_netapp_volume_group_sap_hana" "test_secondary" {
  name                   = "acctest-NetAppVolumeGroup-Secondary-%[2]d"
  location               = azurerm_resource_group.test_secondary.location
  resource_group_name    = azurerm_resource_group.test_secondary.name
  account_name           = azurerm_netapp_account.test_secondary.name
  group_description      = "Test volume group secondary"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-1-Secondary-%[2]d"
    volume_path                = "my-unique-file-path-1-Secondary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test_secondary.id
    subnet_id                  = azurerm_subnet.test_secondary.id
    zone                       = "2"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_netapp_volume_group_sap_hana.test_primary.location
      remote_volume_resource_id = azurerm_netapp_volume_group_sap_hana.test_primary.volume[0].id
      replication_frequency     = "daily"
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-2-Secondary-%[2]d"
    volume_path                = "my-unique-file-path-2-Secondary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test_secondary.id
    subnet_id                  = azurerm_subnet.test_secondary.id
    zone                       = "2"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-3-Secondary-%[2]d"
    volume_path                = "my-unique-file-path-3-Secondary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test_secondary.id
    subnet_id                  = azurerm_subnet.test_secondary.id
    zone                       = "2"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_netapp_volume_group_sap_hana.test_primary.location
      remote_volume_resource_id = azurerm_netapp_volume_group_sap_hana.test_primary.volume[2].id
      replication_frequency     = "daily"
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_netapp_volume_group_sap_hana.test_primary
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) crossZoneReplicationAZ(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateAvailabilityZone(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test_primary" {
  name                   = "acctest-NetAppVolumeGroup-Primary-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group primary"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-1-Primary-%[2]d"
    volume_path                = "my-unique-file-path-1-Primary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-2-Primary-%[2]d"
    volume_path                = "my-unique-file-path-2-Primary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-3-Primary-%[2]d"
    volume_path                = "my-unique-file-path-3-Primary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}

resource "azurerm_netapp_volume_group_sap_hana" "test_secondary" {
  name                   = "acctest-NetAppVolumeGroup-Secondary-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group secondary"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-1-Secondary-%[2]d"
    volume_path                = "my-unique-file-path-1-Secondary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_netapp_volume_group_sap_hana.test_primary.location
      remote_volume_resource_id = azurerm_netapp_volume_group_sap_hana.test_primary.volume[0].id
      replication_frequency     = "daily"
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-2-Secondary-%[2]d"
    volume_path                = "my-unique-file-path-2-Secondary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-3-Secondary-%[2]d"
    volume_path                = "my-unique-file-path-3-Secondary-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_netapp_volume_group_sap_hana.test_primary.location
      remote_volume_resource_id = azurerm_netapp_volume_group_sap_hana.test_primary.volume[2].id
      replication_frequency     = "daily"
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_netapp_volume_group_sap_hana.test_primary
  ]
}
`, template, data.RandomInteger)
}

func (r NetAppVolumeGroupSAPHanaResource) templateForAvgCrossRegionReplicationAZ(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateAvailabilityZone(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "test_secondary" {
  name     = "acctestRG-netapp-secondary-%[2]d"
  location = "%[3]s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_virtual_network" "test_secondary" {
  name                = "acctest-VirtualNetwork-secondary-%[2]d"
  location            = azurerm_resource_group.test_secondary.location
  resource_group_name = azurerm_resource_group.test_secondary.name
  address_space       = ["10.89.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test_secondary" {
  name                 = "acctest-DelegatedSubnet-secondary-%[2]d"
  resource_group_name  = azurerm_resource_group.test_secondary.name
  virtual_network_name = azurerm_virtual_network.test_secondary.name
  address_prefixes     = ["10.89.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test_secondary" {
  name                = "acctest-NetAppAccount-secondary-%[2]d"
  location            = azurerm_resource_group.test_secondary.location
  resource_group_name = azurerm_resource_group.test_secondary.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test_secondary" {
  name                = "acctest-NetAppPool-secondary-%[2]d"
  location            = azurerm_resource_group.test_secondary.location
  resource_group_name = azurerm_resource_group.test_secondary.name
  account_name        = azurerm_netapp_account.test_secondary.name
  service_level       = "Standard"
  size_in_tb          = 10
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (r NetAppVolumeGroupSAPHanaResource) templateForAvgCrossRegionReplication(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templatePPG(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_user_assigned_identity" "test_secondary" {
  name                = "user-assigned-identity-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test_secondary" {
  name                = "acctest-NSG-Secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_virtual_network" "test_secondary" {
  name                = "acctest-VirtualNetwork-Secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test_secondary" {
  name                 = "acctest-DelegatedSubnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test_secondary.name
  address_prefixes     = ["10.88.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "test1_secondary" {
  name                 = "acctest-HostsSubnet-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test_secondary.name
  address_prefixes     = ["10.88.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "test_secondary" {
  subnet_id                 = azurerm_subnet.test1_secondary.id
  network_security_group_id = azurerm_network_security_group.test_secondary.id
}

resource "azurerm_proximity_placement_group" "test_secondary" {
  name                = "acctest-PPG-Secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_availability_set" "test_secondary" {
  name                = "acctest-avset-Secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  platform_update_domain_count = 2
  platform_fault_domain_count  = 2

  proximity_placement_group_id = azurerm_proximity_placement_group.test_secondary.id

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_network_interface" "test_secondary" {
  name                = "acctest-nic-Secondary-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%[3]s"

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test1_secondary.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_linux_virtual_machine" "test_secondary" {
  name                            = "acctest-vm-secondary-%[2]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = "%[3]s"
  size                            = "Standard_D2s_v4"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false
  proximity_placement_group_id    = azurerm_proximity_placement_group.test_secondary.id
  availability_set_id             = azurerm_availability_set.test_secondary.id
  network_interface_ids = [
    azurerm_network_interface.test_secondary.id
  ]

  source_image_reference {
    publisher = "MicrosoftCBLMariner"
    offer     = "azure-linux-3"
    sku       = "azure-linux-3"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test_secondary.id
    ]
  }

  tags = {
    "AzSecPackAutoConfigReady"                                                 = "true",
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true",
    "CreatedOnDate"                                                            = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack"                                                         = "true",
    "Owner"                                                                    = "pmarques"
  }
}

resource "azurerm_netapp_account" "test_secondary" {
  name                = "acctest-NetAppAccount-Secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  depends_on = [
    azurerm_subnet.test_secondary,
    azurerm_subnet.test1_secondary
  ]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test_secondary" {
  name                = "acctest-NetAppPool-Secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test_secondary.name
  service_level       = "Standard"
  size_in_tb          = 8
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func (NetAppVolumeGroupSAPHanaResource) templatePPG(data acceptance.TestData) string {
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

locals {
  admin_username = "testadmin%[1]d"
  admin_password = "Password1234!%[1]d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "user-assigned-identity-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
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

resource "azurerm_subnet" "test1" {
  name                 = "acctest-HostsSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.88.1.0/24"]
}

resource "azurerm_subnet_network_security_group_association" "public" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}

resource "azurerm_proximity_placement_group" "test" {
  name                = "acctest-PPG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_availability_set" "test" {
  name                = "acctest-avset-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  platform_update_domain_count = 2
  platform_fault_domain_count  = 2

  proximity_placement_group_id = azurerm_proximity_placement_group.test.id

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_network_interface" "test" {
  name                = "acctest-nic-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test1.id
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_linux_virtual_machine" "test" {
  name                            = "acctest-vm-%[1]d"
  resource_group_name             = azurerm_resource_group.test.name
  location                        = azurerm_resource_group.test.location
  size                            = "Standard_D2s_v4"
  admin_username                  = local.admin_username
  admin_password                  = local.admin_password
  disable_password_authentication = false
  proximity_placement_group_id    = azurerm_proximity_placement_group.test.id
  availability_set_id             = azurerm_availability_set.test.id
  network_interface_ids = [
    azurerm_network_interface.test.id
  ]

  source_image_reference {
    publisher = "MicrosoftCBLMariner"
    offer     = "azure-linux-3"
    sku       = "azure-linux-3"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }

  tags = {
    "AzSecPackAutoConfigReady"                                                 = "true",
    "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true",
    "CreatedOnDate"                                                            = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack"                                                         = "true",
    "Owner"                                                                    = "pmarques"
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  depends_on = [
    azurerm_subnet.test,
    azurerm_subnet.test1
  ]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 10
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetAppVolumeGroupSAPHanaResource) protocolConversionNFSv3(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templatePPG(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-AVG-HANA-ProtocolConversion-%[2]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for protocol conversion"
  application_identifier = "SH9"

  volume {
    name                         = "acctest-NetAppVolume-1-%[2]d"
    volume_path                  = "my-unique-file-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-%[2]d"
    volume_path                  = "my-unique-file-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-3-%[2]d"
    volume_path                  = "my-unique-file-path-3-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-4-%[2]d"
    volume_path                = "my-unique-file-path-4-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    volume_spec_name           = "data-backup"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv3"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = true
      nfsv41_enabled      = false
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-5-%[2]d"
    volume_path                = "my-unique-file-path-5-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    volume_spec_name           = "log-backup"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv3"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = true
      nfsv41_enabled      = false
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) protocolConversionNFSv41(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templatePPG(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-AVG-HANA-ProtocolConversion-%[2]d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for protocol conversion"
  application_identifier = "SH9"

  volume {
    name                         = "acctest-NetAppVolume-1-%[2]d"
    volume_path                  = "my-unique-file-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "data"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-2-%[2]d"
    volume_path                  = "my-unique-file-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-3-%[2]d"
    volume_path                  = "my-unique-file-path-3-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "shared"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-4-%[2]d"
    volume_path                = "my-unique-file-path-4-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    volume_spec_name           = "data-backup"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-5-%[2]d"
    volume_path                = "my-unique-file-path-5-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    volume_spec_name           = "log-backup"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) basicAvailabilityZone(data acceptance.TestData) string {
	template := NetAppVolumeGroupSAPHanaResource{}.templateAvailabilityZone(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroupSAPHana-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for SAP HANA"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-Data1-%[2]d"
    volume_path                = "my-unique-file-data-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "data"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-Log-%[2]d"
    volume_path                = "my-unique-file-log-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-Shared-%[2]d"
    volume_path                = "my-unique-file-shared-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "shared"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false
    network_features           = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23-50-21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupSAPHanaResource) templateAvailabilityZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      prevent_volume_destruction = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
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

  depends_on = [
    azurerm_subnet.test
  ]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 10
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetAppVolumeGroupSAPHanaResource) volEncryptionCmkSAPHana(data acceptance.TestData, tenantID string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      prevent_volume_destruction = false
    }
  }
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[3]s"

  tags = {
    "SkipNRMSNSG"   = "true",
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

resource "azurerm_key_vault" "test" {
  name                            = "anfkv%[4]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
  soft_delete_retention_days      = 7
  tenant_id                       = "%[2]s"
  sku_name                        = "standard"

  access_policy {
    tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    certificate_permissions = []
    secret_permissions      = []
    storage_permissions     = []
    key_permissions = [
      "Get",
      "Create",
      "Delete",
      "WrapKey",
      "UnwrapKey",
      "GetRotationPolicy",
      "SetRotationPolicy",
    ]
  }

  access_policy {
    tenant_id = azurerm_netapp_account.test.identity.0.tenant_id
    object_id = azurerm_netapp_account.test.identity.0.principal_id

    certificate_permissions = []
    secret_permissions      = []
    storage_permissions     = []
    key_permissions = [
      "Get",
      "Encrypt",
      "Decrypt"
    ]
  }

  tags = {
    "CreatedOnDate" = "2022-07-08T23-50-21Z"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfkey%[4]d"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_netapp_account_encryption" "test" {
  netapp_account_id                     = azurerm_netapp_account.test.id
  system_assigned_identity_principal_id = azurerm_netapp_account.test.identity.0.principal_id
  encryption_key                        = azurerm_key_vault_key.test.versionless_id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test-delegated" {
  name                 = "acctest-Delegated-Subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.88.1.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "test-non-delegated" {
  name                 = "acctest-Non-Delegated-Subnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.88.0.0/24"]
}

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-akv-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test-non-delegated.id

  private_service_connection {
    name                           = "acctest-pe-sc-akv-%[1]d"
    private_connection_resource_id = azurerm_key_vault.test.id
    is_manual_connection           = false
    subresource_names              = ["Vault"]
  }

  tags = {
    CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 8
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23-50-21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.test
  ]
}

resource "azurerm_netapp_volume_group_sap_hana" "test" {
  name                   = "acctest-NetAppVolumeGroupSAPHana-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for SAP HANA"
  application_identifier = "TST"

  volume {
    name                          = "acctest-NetAppVolume-Data1-%[1]d"
    volume_path                   = "my-unique-file-data-path-1-%[1]d"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.test.id
    subnet_id                     = azurerm_subnet.test-delegated.id
    zone                          = "1"
    volume_spec_name              = "data"
    storage_quota_in_gb           = 1024
    throughput_in_mibps           = 24
    protocols                     = ["NFSv4.1"]
    security_style                = "unix"
    snapshot_directory_visible    = false
    encryption_key_source         = "Microsoft.KeyVault"
    key_vault_private_endpoint_id = azurerm_private_endpoint.test.id
    network_features              = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                          = "acctest-NetAppVolume-Log-%[1]d"
    volume_path                   = "my-unique-file-log-path-%[1]d"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.test.id
    subnet_id                     = azurerm_subnet.test-delegated.id
    zone                          = "1"
    volume_spec_name              = "log"
    storage_quota_in_gb           = 1024
    throughput_in_mibps           = 24
    protocols                     = ["NFSv4.1"]
    security_style                = "unix"
    snapshot_directory_visible    = false
    encryption_key_source         = "Microsoft.KeyVault"
    key_vault_private_endpoint_id = azurerm_private_endpoint.test.id
    network_features              = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                          = "acctest-NetAppVolume-Shared-%[1]d"
    volume_path                   = "my-unique-file-shared-path-%[1]d"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.test.id
    subnet_id                     = azurerm_subnet.test-delegated.id
    zone                          = "1"
    volume_spec_name              = "shared"
    storage_quota_in_gb           = 1024
    throughput_in_mibps           = 24
    protocols                     = ["NFSv4.1"]
    security_style                = "unix"
    snapshot_directory_visible    = false
    encryption_key_source         = "Microsoft.KeyVault"
    key_vault_private_endpoint_id = azurerm_private_endpoint.test.id
    network_features              = "Standard"

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }
}
`, data.RandomInteger, tenantID, data.Locations.Primary, data.RandomIntOfLength(17))
}
