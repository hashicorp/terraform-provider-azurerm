// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/volumes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// getOverriddenTestLocations returns the overridden test locations for the NetApp Volume tests, specifically for CRR
// that is not aligned with traditional region pairs.
func getOverriddenTestLocations() struct {
	Primary   string
	Secondary string
} {
	return struct {
		Primary   string
		Secondary string
	}{
		Primary:   "westus2",
		Secondary: "eastus2",
	}
}

type NetAppVolumeResource struct{}

func TestAccNetAppVolume_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

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

func TestAccNetAppVolume_backupPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_backupPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backupPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateBackupPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_availabilityZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.availabilityZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("zone").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_nfsv41(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nfsv41(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_features").HasValue(string(volumes.NetworkFeaturesBasic)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_standardNetworkFeature(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardNetworkFeature(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network_features").HasValue(string(volumes.NetworkFeaturesStandard)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_snapshotPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.snapshotPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("data_protection_snapshot_policy.0.snapshot_policy_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_snapshotPolicyUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.snapshotPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.snapshotPolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_crossRegionReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test_secondary")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.crossRegionReplication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("data_protection_replication.0.endpoint_type").HasValue("dst"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_nfsv3FromSnapshot(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test_snapshot_vol")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nfsv3FromSnapshot(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("create_from_snapshot_resource_id").MatchesRegex(regexp.MustCompile(fmt.Sprintf("(.)/snapshots/acctest-Snapshot-%d", data.RandomInteger))),
			),
		},
		data.ImportStep("create_from_snapshot_resource_id"),
	})
}

func TestAccNetAppVolume_nfsv3SnapshotDirectoryVisibleFalse(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test_snapshot_directory_visible_false")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nfsv3SnapshotDirectoryVisibleFalse(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("snapshot_directory_visible").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_nfsv3SnapshotDirectoryVisibleTrue(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test_snapshot_directory_visible_true")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.nfsv3SnapshotDirectoryVisibleTrue(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("snapshot_directory_visible").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_netapp_volume"),
		},
	})
}

func TestAccNetAppVolume_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_level").HasValue("Standard"),
				check.That(data.ResourceName).Key("storage_quota_in_gb").HasValue("101"),
				check.That(data.ResourceName).Key("export_policy_rule.#").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("3"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaR"),
				check.That(data.ResourceName).Key("mount_ip_addresses.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.completePoolQosManual(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_quota_in_gb").HasValue("100"),
				check.That(data.ResourceName).Key("export_policy_rule.#").HasValue("3"),
				check.That(data.ResourceName).Key("tags.%").HasValue("3"),
				check.That(data.ResourceName).Key("throughput_in_mibps").HasValue("64"),
			),
		},
		data.ImportStep(),
		{
			Config: r.completePoolQosManualNewThroughput(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_quota_in_gb").HasValue("101"),
				check.That(data.ResourceName).Key("export_policy_rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("tags.%").HasValue("4"),
				check.That(data.ResourceName).Key("tags.FoO").HasValue("BaR"),
				check.That(data.ResourceName).Key("tags.bAr").HasValue("fOo"),
				check.That(data.ResourceName).Key("throughput_in_mibps").HasValue("63"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_updateExportPolicyRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("export_policy_rule.#").HasValue("3"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateExportPolicyRule(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("export_policy_rule.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_volEncryptionCmkUserAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.volEncryptionCmkUserAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_volEncryptionCmkSystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	tenantID := os.Getenv("ARM_TENANT_ID")

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.volEncryptionCmkSystemAssigned(data, tenantID),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccNetAppVolume_serviceLevelUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_netapp_volume", "test")
	r := NetAppVolumeResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceLevelUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t NetAppVolumeResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := volumes.ParseVolumeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.NetApp.VolumeClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Netapp Volume (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (NetAppVolumeResource) volEncryptionCmkUserAssigned(data acceptance.TestData, tenantID string) string {
	cmkUserAssginedTemplate := NetAppAccountEncryptionResource{}.cmkUserAssigned(data, tenantID)
	networkTemplate := NetAppVolumeResource{}.networkTemplate(data)
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-akv-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test-non-delegated.id

  private_service_connection {
    name                           = "acctest-pe-sc-akv-%[3]d"
    private_connection_resource_id = azurerm_key_vault.test.id
    is_manual_connection           = false
    subresource_names              = ["Vault"]
  }

  tags = {
    CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.test
  ]
}

resource "azurerm_netapp_volume" "test" {
  name                          = "acctest-NetAppVolume-%[3]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  account_name                  = azurerm_netapp_account.test.name
  pool_name                     = azurerm_netapp_pool.test.name
  volume_path                   = "my-unique-file-path-%[3]d"
  service_level                 = "Standard"
  subnet_id                     = azurerm_subnet.test-delegated.id
  storage_quota_in_gb           = 100
  network_features              = "Standard"
  encryption_key_source         = "Microsoft.KeyVault"
  key_vault_private_endpoint_id = azurerm_private_endpoint.test.id

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.test,
    azurerm_private_endpoint.test
  ]
}
`, cmkUserAssginedTemplate, networkTemplate, data.RandomInteger)
}

func (NetAppVolumeResource) volEncryptionCmkSystemAssigned(data acceptance.TestData, tenantID string) string {
	cmkUserAssginedTemplate := NetAppAccountEncryptionResource{}.cmkSystemAssigned(data, tenantID)
	networkTemplate := NetAppVolumeResource{}.networkTemplate(data)
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_private_endpoint" "test" {
  name                = "acctest-pe-akv-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  subnet_id           = azurerm_subnet.test-non-delegated.id

  private_service_connection {
    name                           = "acctest-pe-sc-akv-%[3]d"
    private_connection_resource_id = azurerm_key_vault.test.id
    is_manual_connection           = false
    subresource_names              = ["Vault"]
  }

  tags = {
    CreatedOnDate = "2023-10-03T19:58:43.6509795Z"
  }
}

resource "azurerm_netapp_pool" "test" {
  name                = "acctest-NetAppPool-%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Standard"
  size_in_tb          = 4

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.test
  ]
}

resource "azurerm_netapp_volume" "test" {
  name                          = "acctest-NetAppVolume-%[3]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  account_name                  = azurerm_netapp_account.test.name
  pool_name                     = azurerm_netapp_pool.test.name
  volume_path                   = "my-unique-file-path-%[3]d"
  service_level                 = "Standard"
  subnet_id                     = azurerm_subnet.test-delegated.id
  storage_quota_in_gb           = 100
  network_features              = "Standard"
  encryption_key_source         = "Microsoft.KeyVault"
  key_vault_private_endpoint_id = azurerm_private_endpoint.test.id

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }

  depends_on = [
    azurerm_netapp_account_encryption.test,
    azurerm_private_endpoint.test
  ]
}
`, cmkUserAssginedTemplate, networkTemplate, data.RandomInteger)
}

func (NetAppVolumeResource) backupPolicy(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_vault" "test" {
  name                = "acctest-NetAppBackupVault-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  account_name        = azurerm_netapp_account.test.name

  tags = {
    "testTag" = "testTagValue"
  }
}

resource "azurerm_netapp_backup_policy" "test" {
  name                    = "acctest-NetAppBackupPolicy-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  account_name            = azurerm_netapp_account.test.name
  daily_backups_to_keep   = 2
  weekly_backups_to_keep  = 2
  monthly_backups_to_keep = 2
  enabled                 = true
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[2]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  throughput_in_mibps = 10

  data_protection_backup_policy {
    backup_vault_id  = azurerm_netapp_backup_vault.test.id
    backup_policy_id = azurerm_netapp_backup_policy.test.id
    policy_enabled   = true
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeResource) updateBackupPolicy(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_backup_vault" "test" {
  name                = "acctest-NetAppBackupVault-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  account_name        = azurerm_netapp_account.test.name

  tags = {
    "testTag" = "testTagValue"
  }
}

resource "azurerm_netapp_backup_policy" "test" {
  name                    = "acctest-NetAppBackupPolicy-%[2]d"
  resource_group_name     = azurerm_resource_group.test.name
  location                = azurerm_resource_group.test.location
  account_name            = azurerm_netapp_account.test.name
  daily_backups_to_keep   = 2
  weekly_backups_to_keep  = 2
  monthly_backups_to_keep = 2
  enabled                 = true
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[2]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  throughput_in_mibps = 10

  data_protection_backup_policy {
    backup_vault_id  = azurerm_netapp_backup_vault.test.id
    backup_policy_id = azurerm_netapp_backup_policy.test.id
    policy_enabled   = false
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeResource) basic(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (NetAppVolumeResource) availabilityZone(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  zone                = "1"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.0

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (NetAppVolumeResource) nfsv41(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv4.1"]
  security_style      = "unix"
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv4.1"]
    unix_read_only    = false
    unix_read_write   = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (NetAppVolumeResource) standardNetworkFeature(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  network_features    = "Standard"
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (NetAppVolumeResource) snapshotPolicy(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
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
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-WithSnapshotPolicy-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[2]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  security_style      = "unix"
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  data_protection_snapshot_policy {
    snapshot_policy_id = azurerm_netapp_snapshot_policy.test.id
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeResource) snapshotPolicyUpdate(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
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
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-WithSnapshotPolicy-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[2]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  security_style      = "unix"
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeResource) crossRegionReplication(data acceptance.TestData) string {
	overriddenlocations := getOverriddenTestLocations()
	template := NetAppVolumeResource{}.templateForCrossRegionReplication(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume" "test_primary" {
  name                       = "acctest-NetAppVolume-primary-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  account_name               = azurerm_netapp_account.test.name
  pool_name                  = azurerm_netapp_pool.test.name
  volume_path                = "my-unique-file-path-primary-%[2]d"
  service_level              = "Standard"
  subnet_id                  = azurerm_subnet.test.id
  protocols                  = ["NFSv3"]
  storage_quota_in_gb        = 100
  snapshot_directory_visible = true
  throughput_in_mibps        = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume" "test_secondary" {
  name                       = "acctest-NetAppVolume-secondary-%[2]d"
  location                   = "%[3]s"
  resource_group_name        = azurerm_resource_group.test.name
  account_name               = azurerm_netapp_account.test_secondary.name
  pool_name                  = azurerm_netapp_pool.test_secondary.name
  volume_path                = "my-unique-file-path-secondary-%[2]d"
  service_level              = "Standard"
  subnet_id                  = azurerm_subnet.test_secondary.id
  protocols                  = ["NFSv3"]
  storage_quota_in_gb        = 100
  snapshot_directory_visible = true
  throughput_in_mibps        = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  data_protection_replication {
    endpoint_type             = "dst"
    remote_volume_location    = azurerm_resource_group.test.location
    remote_volume_resource_id = azurerm_netapp_volume.test_primary.id
    replication_frequency     = "10minutes"
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger, overriddenlocations.Secondary)
}

func (NetAppVolumeResource) nfsv3FromSnapshot(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_path         = "my-unique-file-path-%[2]d"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_snapshot" "test" {
  name                = "acctest-Snapshot-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  volume_name         = azurerm_netapp_volume.test.name
}

resource "azurerm_netapp_volume" "test_snapshot_vol" {
  name                             = "acctest-NetAppVolume-NewFromSnapshot-%[2]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  account_name                     = azurerm_netapp_account.test.name
  pool_name                        = azurerm_netapp_pool.test.name
  volume_path                      = "my-unique-file-path-snapshot-%[2]d"
  service_level                    = "Standard"
  subnet_id                        = azurerm_subnet.test.id
  protocols                        = ["NFSv3"]
  storage_quota_in_gb              = 200
  create_from_snapshot_resource_id = azurerm_netapp_snapshot.test.id
  throughput_in_mibps              = 3.125

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_write   = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeResource) nfsv3SnapshotDirectoryVisibleTrue(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume" "test_snapshot_directory_visible_true" {
  name                       = "acctest-NetAppVolume-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  account_name               = azurerm_netapp_account.test.name
  pool_name                  = azurerm_netapp_pool.test.name
  volume_path                = "my-unique-file-path-%[2]d"
  service_level              = "Standard"
  subnet_id                  = azurerm_subnet.test.id
  protocols                  = ["NFSv3"]
  storage_quota_in_gb        = 100
  snapshot_directory_visible = true
  throughput_in_mibps        = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger)
}

func (NetAppVolumeResource) nfsv3SnapshotDirectoryVisibleFalse(data acceptance.TestData) string {
	template := NetAppVolumeResource{}.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_volume" "test_snapshot_directory_visible_false" {
  name                       = "acctest-NetAppVolume-%[2]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  account_name               = azurerm_netapp_account.test.name
  pool_name                  = azurerm_netapp_pool.test.name
  volume_path                = "my-unique-file-path-%[2]d"
  service_level              = "Standard"
  subnet_id                  = azurerm_subnet.test.id
  protocols                  = ["NFSv3"]
  storage_quota_in_gb        = 100
  snapshot_directory_visible = false
  throughput_in_mibps        = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, template, data.RandomInteger)
}

func (r NetAppVolumeResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "import" {
  name                = azurerm_netapp_volume.test.name
  location            = azurerm_netapp_volume.test.location
  resource_group_name = azurerm_netapp_volume.test.resource_group_name
  account_name        = azurerm_netapp_volume.test.account_name
  pool_name           = azurerm_netapp_volume.test.pool_name
  volume_path         = azurerm_netapp_volume.test.volume_path
  service_level       = azurerm_netapp_volume.test.service_level
  subnet_id           = azurerm_netapp_volume.test.subnet_id
  storage_quota_in_gb = azurerm_netapp_volume.test.storage_quota_in_gb
  throughput_in_mibps = azurerm_netapp_volume.test.throughput_in_mibps
}
`, r.basic(data))
}

func (r NetAppVolumeResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  service_level       = "Standard"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 101
  throughput_in_mibps = 1.562

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  export_policy_rule {
    rule_index        = 2
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = true
    unix_read_write   = false
  }

  export_policy_rule {
    rule_index        = 3
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = true
    unix_read_write   = false
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "FoO"              = "BaR",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r NetAppVolumeResource) completePoolQosManual(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  service_level       = "Standard"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 100
  throughput_in_mibps = 64

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  export_policy_rule {
    rule_index        = 2
    allowed_clients   = ["1.2.5.0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = true
    unix_read_write   = false
  }

  export_policy_rule {
    rule_index        = 3
    allowed_clients   = ["1.2.6.0/24"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = true
    unix_read_write   = false
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "FoO"              = "BaR",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.templatePoolQosManual(data), data.RandomInteger, data.RandomInteger)
}

func (r NetAppVolumeResource) completePoolQosManualNewThroughput(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  service_level       = "Standard"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 101
  throughput_in_mibps = 63

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["1.2.3.0/24"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = false
    unix_read_write   = true
  }

  export_policy_rule {
    rule_index        = 2
    allowed_clients   = ["1.2.5.0"]
    protocols_enabled = ["NFSv3"]
    unix_read_only    = true
    unix_read_write   = false
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "FoO"              = "BaR",
    "bAr"              = "fOo",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.templatePoolQosManual(data), data.RandomInteger, data.RandomInteger)
}

func (r NetAppVolumeResource) updateExportPolicyRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test.name
  service_level       = "Standard"
  volume_path         = "my-unique-file-path-%d"
  subnet_id           = azurerm_subnet.test.id
  protocols           = ["NFSv3"]
  storage_quota_in_gb = 101

  export_policy_rule {
    rule_index          = 1
    allowed_clients     = ["1.2.4.0/24", "1.3.4.0"]
    protocols_enabled   = ["NFSv3"]
    unix_read_only      = false
    unix_read_write     = true
    root_access_enabled = true
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "FoO"              = "BaR",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r NetAppVolumeResource) templateForCrossRegionReplication(data acceptance.TestData) string {
	overriddenlocations := getOverriddenTestLocations()
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test_secondary" {
  name                = "acctest-VirtualNetwork-secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test_secondary" {
  name                 = "acctest-Subnet-secondary-%[2]d"
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

resource "azurerm_netapp_account" "test_secondary" {
  name                = "acctest-NetAppAccount-secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_pool" "test_secondary" {
  name                = "acctest-NetAppPool-secondary-%[2]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test_secondary.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.template(data), data.RandomInteger, overriddenlocations.Secondary)
}

func (r NetAppVolumeResource) template(data acceptance.TestData) string {
	overriddenlocations := getOverriddenTestLocations()
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true",
    "SkipNRMSNSG"      = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
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
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.templateProviderFeatureFlags(), data.RandomInteger, overriddenlocations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r NetAppVolumeResource) templatePoolQosManual(data acceptance.TestData) string {
	overriddenlocations := getOverriddenTestLocations()
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-netapp-%d"
  location = "%s"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
    "SkipNRMSNSG"      = "true"
  }
}

resource "azurerm_network_security_group" "test" {
  name                = "acctest-NSG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    environment        = "Production",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-VirtualNetwork-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.88.0.0/16"]

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-Subnet-%d"
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
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.templateProviderFeatureFlags(), data.RandomInteger, overriddenlocations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r NetAppVolumeResource) networkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
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
`, data.RandomInteger)
}

func (NetAppVolumeResource) templateProviderFeatureFlags() string {
	return `
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }

    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }

	netapp {
      prevent_volume_destruction            = false
      delete_backups_on_backup_vault_destroy = true
    } 
  }
}
`
}

func (r NetAppVolumeResource) serviceLevelUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_netapp_pool" "test2" {
  name                = "acctest-NetAppPool2-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  service_level       = "Premium"
  size_in_tb          = 4
  qos_type            = "Manual"

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume" "test" {
  name                = "acctest-NetAppVolume-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  account_name        = azurerm_netapp_account.test.name
  pool_name           = azurerm_netapp_pool.test2.name
  volume_path         = "my-unique-file-path-%[2]d"
  service_level       = "Premium"
  subnet_id           = azurerm_subnet.test.id
  storage_quota_in_gb = 100
  throughput_in_mibps = 1.562

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}
`, r.template(data), data.RandomInteger)
}
