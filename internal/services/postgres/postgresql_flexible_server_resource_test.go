// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package postgres_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2024-08-01/servers"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PostgresqlFlexibleServerResource struct{}

func TestAccPostgresqlFlexibleServer_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_mb").HasValue("32768"),
				check.That(data.ResourceName).Key("storage_tier").HasValue("P4"),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
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

func TestAccPostgresqlFlexibleServer_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_completeUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_tier").HasValue("P4"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("32768"),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.completeUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
				check.That(data.ResourceName).Key("storage_tier").HasValue("P6"),
				check.That(data.ResourceName).Key("storage_mb").HasValue("65536"),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_updateMaintenanceWindow(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateMaintenanceWindow(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateMaintenanceWindowUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_updateVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withVersion(data, 12, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.withVersion(data, 13, ""),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.withVersion(data, 13, "Update"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.withVersion(data, 14, "Default"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_geoRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoRestoreSource(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			PreConfig: func() { time.Sleep(30 * time.Minute) },
			Config:    r.geoRestore(data),
			Check:     acceptance.ComposeTestCheckFunc(),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_pointInTimeRestore(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.pointInTimeRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_postgresql_flexible_server.pitr").ExistsInAzure(r),
				check.That("azurerm_postgresql_flexible_server.pitr").Key("fqdn").Exists(),
				check.That("azurerm_postgresql_flexible_server.pitr").Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode", "point_in_time_restore_time_in_utc"),
	})
}

func TestAccPostgresqlFlexibleServer_failover(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.failover(data, "1", "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.failover(data, "2", "1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.failoverRemoveHA(data, "2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.failover(data, "2", "1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_geoRedundantBackupEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.geoRedundantBackupEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_authConfig(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.authConfig(data, false, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.authConfig(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_disablePwdAuth(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			// starts from pwdEnabled set to `false` to test add `admininistrator_login`
			Config: r.authConfig(data, true, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.authConfig(data, true, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_createWithCustomerManagedKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withCustomerManagedKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_postgresql_flexible_server.test").Key("customer_managed_key.0.key_vault_key_id").Exists(),
				check.That("azurerm_postgresql_flexible_server.test").Key("customer_managed_key.0.primary_user_assigned_identity_id").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_replica(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			PreConfig: func() { time.Sleep(15 * time.Minute) },
			Config:    r.replica(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_postgresql_flexible_server.replica").ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateReplicationRole(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_postgresql_flexible_server.replica").ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_upgradeVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.upgradeVersion(data, "13"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.upgradeVersion(data, "14"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.upgradeVersion(data, "15"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_identitySystemAssigned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.identitySystemAssigned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_enableGeoRedundantBackup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enableGeoRedundantBackup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_autoGrowEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.autoGrowEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.autoGrowEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
				check.That(data.ResourceName).Key("public_network_access_enabled").Exists(),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_invalidStorageTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.invalidStorageTier(data),
			ExpectError: regexp.MustCompile("invalid 'storage_tier'"),
		},
	})
}

func TestAccPostgresqlFlexibleServer_invalidStorageTierScalingStorageMb(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invalidStorageTierScaling(data, "P4", "32768"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config:      r.invalidStorageTierScaling(data, "P4", "262144"),
			ExpectError: regexp.MustCompile("invalid 'storage_tier'"),
		},
	})
}

func TestAccPostgresqlFlexibleServer_invalidStorageTierScalingStorageTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invalidStorageTierScaling(data, "P4", "32768"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config:      r.invalidStorageTierScaling(data, "P80", "32768"),
			ExpectError: regexp.MustCompile("invalid 'storage_tier'"),
		},
	})
}

func TestAccPostgresqlFlexibleServer_invalidStorageTierScalingStorageMbStorageTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invalidStorageTierScaling(data, "P4", "32768"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config:      r.invalidStorageTierScaling(data, "P6", "131072"),
			ExpectError: regexp.MustCompile("invalid 'storage_tier'"),
		},
	})
}

func TestAccPostgresqlFlexibleServer_updateOnlyWithStorageMb(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invalidStorageTierScaling(data, "P4", "32768"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateOnlyWithStorageMb(data, "65536"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccPostgresqlFlexibleServer_updateOnlyWithStorageTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.invalidStorageTierScaling(data, "P4", "32768"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateOnlyWithStorageTier(data, "P10"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateStorageTierWithoutProperty(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				// the storage tier is not changed to the default value, because p10 is still valid.
				check.That(data.ResourceName).Key("storage_tier").HasValue("P10"),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.updateOnlyWithStorageMb(data, "262144"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_tier").HasValue("P15"),
			),
		},
	})
}

func TestAccPostgresqlFlexibleServer_publicNetworkAccessEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.publicNetworkAccessEnabled(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
		{
			Config: r.publicNetworkAccessEnabled(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("administrator_password", "create_mode"),
	})
}

func TestAccPostgresqlFlexibleServer_writeOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.writeOnlyPassword(data, "QAZwsx123", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password_wo_version"),
			{
				Config: r.writeOnlyPassword(data, "QAZwsx123updated", 2),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password_wo_version"),
		},
	})
}

func TestAccPostgresqlFlexibleServer_updateToWriteOnlyPassword(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_postgresql_flexible_server", "test")
	r := PostgresqlFlexibleServerResource{}

	resource.ParallelTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password"),
			{
				Config: r.writeOnlyPassword(data, "QAZwsx123", 1),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password", "administrator_password_wo_version"),
			{
				Config: r.basic(data),
				Check:  check.That(data.ResourceName).ExistsInAzure(r),
			},
			data.ImportStep("administrator_password"),
		},
	})
}

func (PostgresqlFlexibleServerResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := servers.ParseFlexibleServerID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Postgres.FlexibleServersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (PostgresqlFlexibleServerResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r PostgresqlFlexibleServerResource) updateOnlyWithStorageTier(data acceptance.TestData, storageTier string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 65536
  storage_tier           = "%s"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger, storageTier)
}

func (r PostgresqlFlexibleServerResource) updateStorageTierWithoutProperty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 65536
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateOnlyWithStorageMb(data acceptance.TestData, storageMb string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = %s
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger, storageMb)
}

func (r PostgresqlFlexibleServerResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) withVersion(data acceptance.TestData, versionNum int, creatMode string) string {
	createModeProp := ""
	if creatMode != "" {
		createModeProp = fmt.Sprintf("create_mode = \"%s\"", creatMode)
	}

	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "%d"
  %s
  sku_name = "GP_Standard_D2s_v3"
  zone     = "2"
}
`, r.template(data), data.RandomInteger, versionNum, createModeProp)
}

func (r PostgresqlFlexibleServerResource) geoRestoreSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%[1]d"
  location = "eastus"
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  storage_mb                   = 32768
  version                      = "12"
  sku_name                     = "GP_Standard_D2s_v3"
  zone                         = "1"
  geo_redundant_backup_enabled = true
}
`, data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) geoRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "geo_restore" {
  name                              = "acctest-fs-restore-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = "westus"
  create_mode                       = "GeoRestore"
  source_server_id                  = azurerm_postgresql_flexible_server.test.id
  point_in_time_restore_time_in_utc = "%s"
}
`, r.geoRestoreSource(data), data.RandomInteger, time.Now().Add(time.Duration(15)*time.Minute).UTC().Format(time.RFC3339))
}

func (r PostgresqlFlexibleServerResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "import" {
  name                   = azurerm_postgresql_flexible_server.test.name
  resource_group_name    = azurerm_postgresql_flexible_server.test.resource_group_name
  location               = azurerm_postgresql_flexible_server.test.location
  administrator_login    = azurerm_postgresql_flexible_server.test.administrator_login
  administrator_password = azurerm_postgresql_flexible_server.test.administrator_password
  version                = azurerm_postgresql_flexible_server.test.version
  storage_mb             = azurerm_postgresql_flexible_server.test.storage_mb
  sku_name               = azurerm_postgresql_flexible_server.test.sku_name
  zone                   = azurerm_postgresql_flexible_server.test.zone
}
`, r.basic(data))
}

func (r PostgresqlFlexibleServerResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acc%[2]d.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%[2]d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name

  depends_on = [azurerm_subnet.test]
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-fs-%[2]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  administrator_login           = "adminTerraform"
  administrator_password        = "QAZwsx123"
  version                       = "13"
  backup_retention_days         = 7
  storage_mb                    = 32768
  delegated_subnet_id           = azurerm_subnet.test.id
  private_dns_zone_id           = azurerm_private_dns_zone.test.id
  public_network_access_enabled = false
  sku_name                      = "GP_Standard_D2s_v3"
  zone                          = "1"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "2"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Test"
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.test]
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) completeUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "acctest-sn-%[2]d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.2.0/24"]
  service_endpoints    = ["Microsoft.Storage"]
  delegation {
    name = "fs"
    service_delegation {
      name = "Microsoft.DBforPostgreSQL/flexibleServers"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acc%[2]d.postgres.database.azure.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "acctestVnetZone%[2]d.com"
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
  resource_group_name   = azurerm_resource_group.test.name

  depends_on = [azurerm_subnet.test]
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-fs-%[2]d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  administrator_login           = "adminTerraform"
  administrator_password        = "123wsxQAZ"
  version                       = "13"
  backup_retention_days         = 10
  storage_mb                    = 65536
  storage_tier                  = "P6"
  delegated_subnet_id           = azurerm_subnet.test.id
  private_dns_zone_id           = azurerm_private_dns_zone.test.id
  public_network_access_enabled = false
  sku_name                      = "GP_Standard_D2s_v3"
  zone                          = "2"

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "1"
  }

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }

  tags = {
    ENV = "Stage"
  }

  depends_on = [azurerm_private_dns_zone_virtual_network_link.test]
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateMaintenanceWindow(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  storage_mb             = 32768
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 8
    start_minute = 0
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateMaintenanceWindowUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  storage_mb             = 32768
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"

  maintenance_window {
    day_of_week  = 3
    start_hour   = 7
    start_minute = 15
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  storage_mb             = 32768
  sku_name               = "MO_Standard_E2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) pointInTimeRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "pitr" {
  name                              = "acctest-fs-pitr-%d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  create_mode                       = "PointInTimeRestore"
  source_server_id                  = azurerm_postgresql_flexible_server.test.id
  zone                              = "1"
  point_in_time_restore_time_in_utc = "%s"
}
`, r.basic(data), data.RandomInteger, time.Now().Add(time.Duration(15)*time.Minute).UTC().Format(time.RFC3339))
}

func (r PostgresqlFlexibleServerResource) failover(data acceptance.TestData, primaryZone string, standbyZone string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  version                = "12"
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  zone                   = "%s"
  backup_retention_days  = 10
  storage_mb             = 131072
  sku_name               = "GP_Standard_D2s_v3"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 0
    start_minute = 0
  }

  high_availability {
    mode                      = "ZoneRedundant"
    standby_availability_zone = "%s"
  }
}
`, r.template(data), data.RandomInteger, primaryZone, standbyZone)
}

func (r PostgresqlFlexibleServerResource) failoverRemoveHA(data acceptance.TestData, primaryZone string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  version                = "12"
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  zone                   = "%s"
  backup_retention_days  = 10
  storage_mb             = 131072
  sku_name               = "GP_Standard_D2s_v3"

  maintenance_window {
    day_of_week  = 0
    start_hour   = 0
    start_minute = 0
  }
}
`, r.template(data), data.RandomInteger, primaryZone)
}

func (r PostgresqlFlexibleServerResource) geoRedundantBackupEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  storage_mb                   = 32768
  version                      = "12"
  sku_name                     = "GP_Standard_D2s_v3"
  zone                         = "2"
  backup_retention_days        = 7
  geo_redundant_backup_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) authConfig(data acceptance.TestData, aadEnabled bool, pwdEnabled bool) string {
	tenantIdBlock := ""
	if aadEnabled {
		tenantIdBlock = "tenant_id = data.azurerm_client_config.current.tenant_id"
	}

	return fmt.Sprintf(`
%s

data "azurerm_client_config" "current" {
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "1"

  authentication {
    active_directory_auth_enabled = %[3]t
    password_auth_enabled         = %[4]t
   %[5]s
  }

}
`, r.template(data), data.RandomInteger, aadEnabled, pwdEnabled, tenantIdBlock)
}

func (r PostgresqlFlexibleServerResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) cmkTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-postgresql-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctestmi%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_key_vault" "test" {
  name                     = "acctestkv%s"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.test.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test" {
  name         = "test"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomString)
}

func (r PostgresqlFlexibleServerResource) withCustomerManagedKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  version                = "12"
  sku_name               = "B_Standard_B1ms"
  zone                   = "1"

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  customer_managed_key {
    key_vault_key_id                  = azurerm_key_vault_key.test.id
    primary_user_assigned_identity_id = azurerm_user_assigned_identity.test.id
  }
}
`, r.cmkTemplate(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) replica(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "replica" {
  name                = "acctest-fs-replica-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zone                = "2"
  create_mode         = "Replica"
  source_server_id    = azurerm_postgresql_flexible_server.test.id
}
`, r.basic(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) updateReplicationRole(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "replica" {
  name                = "acctest-fs-replica-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  zone                = "2"
  create_mode         = "Replica"
  source_server_id    = azurerm_postgresql_flexible_server.test.id
  replication_role    = "None"
}
`, r.basic(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) upgradeVersion(data acceptance.TestData, version string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  create_mode            = "Update"
  version                = "%s"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger, version)
}

func (r PostgresqlFlexibleServerResource) enableGeoRedundantBackup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG-postgresql2-%d"
  location = "%s"
}

resource "azurerm_user_assigned_identity" "test2" {
  name                = "acctestmi2%s"
  location            = azurerm_resource_group.test2.location
  resource_group_name = azurerm_resource_group.test2.name
}

resource "azurerm_key_vault" "test2" {
  name                     = "acctestkv2%s"
  location                 = azurerm_resource_group.test2.location
  resource_group_name      = azurerm_resource_group.test2.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "standard"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_user_assigned_identity.test2.principal_id

  key_permissions = ["Get", "List", "WrapKey", "UnwrapKey", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_access_policy" "client2" {
  key_vault_id = azurerm_key_vault.test2.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy", "SetRotationPolicy"]
}

resource "azurerm_key_vault_key" "test2" {
  name         = "test2"
  key_vault_id = azurerm_key_vault.test2.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client2,
    azurerm_key_vault_access_policy.server2,
  ]
}

resource "azurerm_postgresql_flexible_server" "test" {
  name                         = "acctest-fs-%d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  administrator_login          = "adminTerraform"
  administrator_password       = "QAZwsx123"
  storage_mb                   = 32768
  version                      = "12"
  sku_name                     = "B_Standard_B1ms"
  zone                         = "2"
  geo_redundant_backup_enabled = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id, azurerm_user_assigned_identity.test2.id]
  }

  customer_managed_key {
    key_vault_key_id                     = azurerm_key_vault_key.test.id
    primary_user_assigned_identity_id    = azurerm_user_assigned_identity.test.id
    geo_backup_key_vault_key_id          = azurerm_key_vault_key.test2.id
    geo_backup_user_assigned_identity_id = azurerm_user_assigned_identity.test2.id
  }
}
`, r.cmkTemplate(data), data.RandomInteger, data.Locations.Ternary, data.RandomString, data.RandomString, data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) autoGrowEnabled(data acceptance.TestData, autoGrowEnabled bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 32768
  auto_grow_enabled      = %t
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger, autoGrowEnabled)
}

func (r PostgresqlFlexibleServerResource) invalidStorageTier(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = 65536
  storage_tier           = "P4"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger)
}

func (r PostgresqlFlexibleServerResource) invalidStorageTierScaling(data acceptance.TestData, storageTier string, storageMb string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_postgresql_flexible_server" "test" {
  name                   = "acctest-fs-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  storage_mb             = %s
  storage_tier           = "%s"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}
`, r.template(data), data.RandomInteger, storageMb, storageTier)
}

func (r PostgresqlFlexibleServerResource) publicNetworkAccessEnabled(data acceptance.TestData, publicNetworkAccessEnabled bool) string {
	return fmt.Sprintf(`
%s
resource "azurerm_postgresql_flexible_server" "test" {
  name                          = "acctest-fs-%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  administrator_login           = "adminTerraform"
  administrator_password        = "QAZwsx123"
  version                       = "12"
  sku_name                      = "GP_Standard_D2s_v3"
  zone                          = "2"
  public_network_access_enabled = %t
}
`, r.template(data), data.RandomInteger, publicNetworkAccessEnabled)
}

func (r PostgresqlFlexibleServerResource) writeOnlyPassword(data acceptance.TestData, secret string, version int) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "azurerm_postgresql_flexible_server" "test" {
  name                              = "acctest-fs-%[3]d"
  resource_group_name               = azurerm_resource_group.test.name
  location                          = azurerm_resource_group.test.location
  administrator_login               = "adminTerraform"
  administrator_password_wo         = ephemeral.azurerm_key_vault_secret.test.value
  administrator_password_wo_version = %[4]d
  version                           = "12"
  sku_name                          = "GP_Standard_D2s_v3"
  zone                              = "2"
}
`, r.template(data), acceptance.WriteOnlyKeyVaultSecretTemplate(data, secret), data.RandomInteger, version)
}
