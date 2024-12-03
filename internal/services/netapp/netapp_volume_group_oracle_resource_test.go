// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2024-03-01/volumegroups"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type NetAppVolumeGroupOracleResource struct{}

func TestAccNetAppVolumeGroupOracle_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume_group_oracle", "test")
	r := NetAppVolumeGroupOracleResource{}

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
			Config: r.basic(data),
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

func (NetAppVolumeGroupOracleResource) basic(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZone(data)
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
    zone                         = "1"
    volume_spec_name             = "ora-data1"
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
  }

  volume {
    name                         = "acctest-NetAppVolume-Ora2-%[2]d"
    volume_path                  = "my-unique-file-ora-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-data2"
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
      "foo"    = "BAR",
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-Ora3-%[2]d"
    volume_path                  = "my-unique-file-ora-path-3-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-data3"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora4-%[2]d"
    volume_path                = "my-unique-file-ora-path-4-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data4"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora5-%[2]d"
    volume_path                = "my-unique-file-ora-path-5-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data5"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora6-%[2]d"
    volume_path                = "my-unique-file-ora-path-6-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data6"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora7-%[2]d"
    volume_path                = "my-unique-file-ora-path-7-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data7"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora8-%[2]d"
    volume_path                = "my-unique-file-ora-path-8-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data8"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-oralog-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLogMirror-%[2]d"
    volume_path                = "my-unique-file-oralogmirror-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-log-mirror"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraBinary-%[2]d"
    volume_path                = "my-unique-file-orabinary-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-binary"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraBackup-%[2]d"
    volume_path                = "my-unique-file-orabackup-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-backup"
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
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) nfsv3(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZone(data)
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
    zone                         = "1"
    volume_spec_name             = "ora-data1"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv3"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = true
      nfsv41_enabled      = false
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-Ora2-%[2]d"
    volume_path                  = "my-unique-file-ora-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-data2"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv3"]
    security_style               = "unix"
    snapshot_directory_visible   = false

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
      "foo"    = "BAR",
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-Ora3-%[2]d"
    volume_path                  = "my-unique-file-ora-path-3-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-data3"
    storage_quota_in_gb          = 1024
    throughput_in_mibps          = 24
    protocols                    = ["NFSv3"]
    security_style               = "unix"
    snapshot_directory_visible   = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = true
      nfsv41_enabled      = false
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora4-%[2]d"
    volume_path                = "my-unique-file-ora-path-4-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data4"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora5-%[2]d"
    volume_path                = "my-unique-file-ora-path-5-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data5"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora6-%[2]d"
    volume_path                = "my-unique-file-ora-path-6-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data6"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora7-%[2]d"
    volume_path                = "my-unique-file-ora-path-7-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data7"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-Ora8-%[2]d"
    volume_path                = "my-unique-file-ora-path-8-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-data8"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                = "my-unique-file-oralog-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraLogMirror-%[2]d"
    volume_path                = "my-unique-file-oralogmirror-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-log-mirror"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraBinary-%[2]d"
    volume_path                = "my-unique-file-orabinary-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-binary"
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
  }

  volume {
    name                       = "acctest-NetAppVolume-OraBackup-%[2]d"
    volume_path                = "my-unique-file-orabackup-path-%[2]d"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.test.id
    subnet_id                  = azurerm_subnet.test.id
    zone                       = "1"
    volume_spec_name           = "ora-backup"
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
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) avgSnapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZone(data)
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
    name                         = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                  = "my-unique-file-ora-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-data1"
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

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                  = "my-unique-file-ora-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-log"
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

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) updateAvgSnapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZone(data)
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
    name                         = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                  = "my-unique-file-ora-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "2"
    volume_spec_name             = "ora-data1"
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

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                  = "my-unique-file-ora-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "2"
    volume_spec_name             = "ora-log"
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

    data_protection_snapshot_policy {
      snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
    }

    tags = {
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeGroupOracleResource) updateVolumes(data acceptance.TestData) string {
	template := NetAppVolumeGroupOracleResource{}.templateAvailabilityZone(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume_group_oracle" "test" {
  name                   = "acctest-NetAppVolumeGroup-%[2]d"
  location               = azurerm_resource_group.test.location
  resource_group_name    = azurerm_resource_group.test.name
  account_name           = azurerm_netapp_account.test.name
  group_description      = "Test volume group"
  application_identifier = "TST"

  volume {
    name                         = "acctest-NetAppVolume-Ora1-%[2]d"
    volume_path                  = "my-unique-file-ora-path-1-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-data1"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }

  volume {
    name                         = "acctest-NetAppVolume-OraLog-%[2]d"
    volume_path                  = "my-unique-file-ora-path-2-%[2]d"
    service_level                = "Standard"
    capacity_pool_id             = azurerm_netapp_pool.test.id
    subnet_id                    = azurerm_subnet.test.id
    zone                         = "1"
    volume_spec_name             = "ora-log"
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
      "CreatedOnDate"    = "2022-07-08T23:50:21Z",
      "SkipASMAzSecPack" = "true"
    }
  }
}
`, template, data.RandomInteger)
}

// func (NetAppVolumeGroupSAPHanaResource) templatePPG(data acceptance.TestData) string {
// 	return fmt.Sprintf(`
// provider "azurerm" {
//   alias = "all2"
//   features {
//     resource_group {
//       prevent_deletion_if_contains_resources = false
//     }
//     netapp {
//       prevent_volume_destruction             = false
//       delete_backups_on_backup_vault_destroy = true
//     }
//   }
// }

// locals {
//   admin_username = "testadmin%[1]d"
//   admin_password = "Password1234!%[1]d"
// }

// resource "azurerm_resource_group" "test" {
//   name     = "acctestRG-netapp-%[1]d"
//   location = "%[2]s"

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true",
//     "SkipNRMSNSG"      = "true"
//   }
// }

// resource "azurerm_network_security_group" "test" {
//   name                = "acctest-NSG-%[1]d"
//   location            = azurerm_resource_group.test.location
//   resource_group_name = azurerm_resource_group.test.name

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }

// resource "azurerm_virtual_network" "test" {
//   name                = "acctest-VirtualNetwork-%[1]d"
//   location            = azurerm_resource_group.test.location
//   resource_group_name = azurerm_resource_group.test.name
//   address_space       = ["10.6.0.0/16"]

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }

// resource "azurerm_subnet" "test" {
//   name                 = "acctest-DelegatedSubnet-%[1]d"
//   resource_group_name  = azurerm_resource_group.test.name
//   virtual_network_name = azurerm_virtual_network.test.name
//   address_prefixes     = ["10.6.2.0/24"]

//   delegation {
//     name = "testdelegation"

//     service_delegation {
//       name    = "Microsoft.Netapp/volumes"
//       actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
//     }
//   }
// }

// resource "azurerm_subnet" "test1" {
//   name                 = "acctest-HostsSubnet-%[1]d"
//   resource_group_name  = azurerm_resource_group.test.name
//   virtual_network_name = azurerm_virtual_network.test.name
//   address_prefixes     = ["10.6.1.0/24"]
// }

// resource "azurerm_subnet_network_security_group_association" "public" {
//   subnet_id                 = azurerm_subnet.test.id
//   network_security_group_id = azurerm_network_security_group.test.id
// }

// resource "azurerm_proximity_placement_group" "test" {
//   name                = "acctest-PPG-%[1]d"
//   location            = azurerm_resource_group.test.location
//   resource_group_name = azurerm_resource_group.test.name

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }

// resource "azurerm_availability_set" "test" {
//   name                = "acctest-avset-%[1]d"
//   location            = azurerm_resource_group.test.location
//   resource_group_name = azurerm_resource_group.test.name

//   proximity_placement_group_id = azurerm_proximity_placement_group.test.id

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }

// resource "azurerm_network_interface" "test" {
//   name                = "acctest-nic-%[1]d"
//   resource_group_name = azurerm_resource_group.test.name
//   location            = azurerm_resource_group.test.location

//   ip_configuration {
//     name                          = "internal"
//     subnet_id                     = azurerm_subnet.test1.id
//     private_ip_address_allocation = "Dynamic"
//   }

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }

// resource "azurerm_linux_virtual_machine" "test" {
//   name                            = "acctest-vm-%[1]d"
//   resource_group_name             = azurerm_resource_group.test.name
//   location                        = azurerm_resource_group.test.location
//   size                            = "Standard_M8ms"
//   admin_username                  = local.admin_username
//   admin_password                  = local.admin_password
//   disable_password_authentication = false
//   proximity_placement_group_id    = azurerm_proximity_placement_group.test.id
//   availability_set_id             = azurerm_availability_set.test.id
//   network_interface_ids = [
//     azurerm_network_interface.test.id
//   ]

//   source_image_reference {
//     publisher = "Canonical"
//     offer     = "0001-com-ubuntu-server-jammy"
//     sku       = "22_04-lts"
//     version   = "latest"
//   }

//   os_disk {
//     storage_account_type = "Standard_LRS"
//     caching              = "ReadWrite"
//   }

//   tags = {
//     "platformsettings.host_environment.service.platform_optedin_for_rootcerts" = "true",
//     "CreatedOnDate"                                                            = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack"                                                         = "true",
//     "Owner"                                                                    = "pmarques"
//   }
// }

// resource "azurerm_netapp_account" "test" {
//   name                = "acctest-NetAppAccount-%[1]d"
//   location            = azurerm_resource_group.test.location
//   resource_group_name = azurerm_resource_group.test.name

//   depends_on = [
//     azurerm_subnet.test,
//     azurerm_subnet.test1
//   ]

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }

// resource "azurerm_netapp_pool" "test" {
//   name                = "acctest-NetAppPool-%[1]d"
//   location            = azurerm_resource_group.test.location
//   resource_group_name = azurerm_resource_group.test.name
//   account_name        = azurerm_netapp_account.test.name
//   service_level       = "Standard"
//   size_in_tb          = 8
//   qos_type            = "Manual"

//   tags = {
//     "CreatedOnDate"    = "2022-07-08T23:50:21Z",
//     "SkipASMAzSecPack" = "true"
//   }
// }
// `, data.RandomInteger, "eastus")
// }

func (NetAppVolumeGroupOracleResource) templateAvailabilityZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    netapp {
      prevent_volume_destruction             = false
    }
  }
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%[1]d"
  location = "%[2]s"

  tags     = {
    "SkipNRMSNSG" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.6.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-DelegatedSubnet-%[1]d"
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
  size_in_tb          = 18
  qos_type            = "Manual"
}
`, data.RandomInteger, "eastus")
}
