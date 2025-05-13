// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumegroups"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupOracleResource struct{}

func TestAccNetAppVolumeGroupOracle_basicAvailabilityZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

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

func TestAccNetAppVolumeGroupOracle_basicProximityPlacementGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicProximityPlacementGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupOracle_nfsv3(t *testing.T) {
	// Adjust test configurations to exclude data_protection_replication
	// Use the new NetAppVolumeGroupOracleVolume model in test data
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nfsv3(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolumeGroupOracle_snapshotPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

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

func TestAccNetAppVolumeGroupOracle_snapshotPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

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

func TestAccNetAppVolumeGroupOracle_volumeUpdates(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicAvailabilityZone(data),
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

func TestAccNetAppVolumeGroupOracle_volCustomerManagedKeyEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.volEncryptionCmkOracle(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppVolumeGroupOracleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := volumegroups.ParseVolumeGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.VolumeGroupClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (NetAppVolumeGroupOracleResource) basicAvailabilityZone(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZoneOracle(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroupOracle-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                = "my-unique-file-ora-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-data1"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-oralog-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-log"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) basicProximityPlacementGroup(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templatePpgOracle(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroupOracle-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                  = "my-unique-file-ora-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "ora-data1"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false
    network_features             = "Basic"

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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                  = "my-unique-file-oralog-path-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    proximity_placement_group_id = azurerm_proximity_placement_group.test.id
    volume_spec_name             = "ora-log"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv4.1"]
    security_style               = "unix"
    snapshot_directory_visible   = false
    network_features             = "Basic"

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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
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

func (NetAppVolumeGroupOracleResource) nfsv3(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZoneOracle(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroupOracle-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                = "my-unique-file-ora-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-data1"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-oralog-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-log"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) avgSnapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZoneOracle(data)
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                = "my-unique-file-ora-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-data1"
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

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-ora-path-2-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-log"
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

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_netapp_snapshot_policy.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) updateAvgSnapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZoneOracle(data)
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                = "my-unique-file-ora-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-data1"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-ora-path-2-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-log"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  depends_on = [
    azurerm_netapp_snapshot_policy.test
  ]
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) updateVolumes(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZoneOracle(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroupOracle-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                       = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                = "my-unique-file-ora-path-1-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-data1"
    storage_quota_in_gb        = 1200
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-oralog-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "2"
    volume_spec_name           = "ora-log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) volEncryptionCmkOracle(data acceptance.TestData, tenantID string) string {
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault" "test" {
  name                            = "anfakv%[1]d"
  location                        = azurerm_resource_group.test.location
  resource_group_name             = azurerm_resource_group.test.name
  enabled_for_disk_encryption     = true
  enabled_for_deployment          = true
  enabled_for_template_deployment = true
  purge_protection_enabled        = true
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
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_key_vault_key" "test" {
  name         = "anfenckey%[1]d"
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
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
  size_in_tb          = 4
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.test
  ]
}

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroupOracle-%[1]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                          = "acctest-NetAppVolume-Ora1-%[1]d"
    volume_path                   = "my-unique-file-ora-path-1-%[1]d"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.test.id
    subnet_id                     = azurerm_subnet.test-delegated.id
    zone                          = "1"
    volume_spec_name              = "ora-data1"
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
    name                          = "acctest-NetAppVolume-OraLog-%[1]d"
    volume_path                   = "my-unique-file-oralog-path-%[1]d"
    service_level                 = "Standard"
    capacity_pool_id              = azurerm_netapp_pool.test.id
    subnet_id                     = azurerm_subnet.test-delegated.id
    zone                          = "1"
    volume_spec_name              = "ora-log"
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
`, data.RandomInteger, tenantID, data.Locations.Primary)
}

func (NetAppVolumeGroupOracleResource) templatePpgOracle(data acceptance.TestData) string {
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

locals {
  admin_username = "testadmin%[1]d"
  admin_password = "Password1234!%[1]d"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "user-assigned-identity-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate" = "2022-07-08T23:50:21Z"
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-DelegatedSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

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
  address_prefixes     = ["10.0.1.0/24"]
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
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
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  patch_assessment_mode = "AutomaticByPlatform"

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
    "CreatedOnDate"                                                            = "2022-07-08T23:50:21Z",
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
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
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
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (NetAppVolumeGroupOracleResource) templateAvailabilityZoneOracle(data acceptance.TestData) string {
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

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-DelegatedSubnet-%[1]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "test" {
  name                = "acctest-NetAppAccount-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  depends_on = [
    azurerm_subnet.test
  ]
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"
}
`, data.RandomInteger, data.Locations.Primary)
}
