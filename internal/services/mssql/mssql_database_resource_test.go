// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MsSqlDatabaseResource struct{}

func TestAccMsSqlDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_free(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.freeTier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

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

func TestAccMsSqlDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	maintenance_configuration_name := "SQL_Default"

	switch data.Locations.Primary {
	case "westeurope":
		maintenance_configuration_name = "SQL_WestEurope_DB_2"
	case "francecentral":
		maintenance_configuration_name = "SQL_FranceCentral_DB_1"
	}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("maintenance_configuration_name").HasValue(maintenance_configuration_name),
				check.That(data.ResourceName).Key("max_size_gb").HasValue("10"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("Local"),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Test"),
			),
		},
		data.ImportStep("sample_name"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("license_type").HasValue("LicenseIncluded"),
				check.That(data.ResourceName).Key("max_size_gb").HasValue("2"),
				check.That(data.ResourceName).Key("enclave_type").HasValue("Default"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("Staging"),
			),
		},
		data.ImportStep("sample_name"),
	})
}

func TestAccMsSqlDatabase_elasticPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_pool_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("ElasticPool"),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPoolDisassociation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_gp(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_gpServerless(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gpServerless(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_pause_delay_in_minutes").HasValue("42"),
				check.That(data.ResourceName).Key("min_capacity").HasValue("0.75"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_S_Gen5_2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.gpServerlessUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("auto_pause_delay_in_minutes").HasValue("90"),
				check.That(data.ResourceName).Key("min_capacity").HasValue("1.25"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_S_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_updateLicenseType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.gpWithLicenseType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.gpServerlessWithNullLicenseType(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("license_type"),
	})
}

func TestAccMsSqlDatabase_bc(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	// Limited regional availability for BC
	data.Locations.Primary = "westeurope"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bc(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_scale").HasValue("true"),
				check.That(data.ResourceName).Key("sku_name").HasValue("BC_Gen5_2"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.bcUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_scale").HasValue("false"),
				check.That(data.ResourceName).Key("sku_name").HasValue("BC_Gen5_2"),
				check.That(data.ResourceName).Key("zone_redundant").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_hs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("2"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.hsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("4"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_hsWithRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hsWithRetentionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("2"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.hsWithRetentionPolicyUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("4"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_hsWithLongRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hsWithLongRetentionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("read_replica_count").HasValue("2"),
				check.That(data.ResourceName).Key("sku_name").HasValue("HS_Gen5_2"),
			),
		},
	})
}

func TestAccMsSqlDatabase_s0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.s0WithRetentionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_createCopyMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "copy")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createCopyMode(data, `enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep("create_mode", "creation_source_database_id"),
	})
}

func TestAccMsSqlDatabase_createCopyModeError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "copy")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.createCopyMode(data, ""),
			ExpectError: regexp.MustCompile("specifying different 'enclave_type' properties for 'create_mode'"),
		},
	})
}

func TestAccMsSqlDatabase_createPITRMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),

		{
			PreConfig: func() { time.Sleep(11 * time.Minute) },
			Config:    r.createPITRMode(data, time.Now().Add(time.Duration(13)*time.Minute).UTC().Format(time.RFC3339)),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_mssql_database.pitr").ExistsInAzure(r),
			),
		},

		data.ImportStep("creation_source_database_id", "restore_point_in_time"),
	})
}

func TestAccMsSqlDatabase_createSecondaryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "secondary")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createSecondaryMode(data, "test1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		// SQL changes 'create_mode' to Default after creating the secondary and
		// clears the 'creation_source_database_id' which will cause a diff...
		data.ImportStep("sample_name", "create_mode", "creation_source_database_id"),
		{
			Config: r.createSecondaryMode(data, "test2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("sample_name", "create_mode", "creation_source_database_id"),
	})
}

func TestAccMsSqlDatabase_createOnlineSecondaryMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "secondary")

	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createOnlineSecondaryMode(data, "test1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("sample_name", "create_mode", "creation_source_database_id"),
		{
			Config: r.createOnlineSecondaryMode(data, "test2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep("sample_name", "create_mode", "creation_source_database_id"),
	})
}

func TestAccMsSqlDatabase_scaleReplicaSet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "primary")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scaleReplicaSet(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "P2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "BC_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "GP_Gen5_2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "S2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "Basic"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
		{
			Config: r.scaleReplicaSet(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("sample_name", "license_type"),
	})
}

func TestAccMsSqlDatabase_scaleReplicaSetWithFailovergroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "secondary")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.scaleReplicaSetWithFailovergroup(data, "GP_Gen5_2", 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.scaleReplicaSetWithFailovergroup(data, "GP_Gen5_8", 25),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_8"),
			),
		},
		data.ImportStep(),
		{
			Config: r.scaleReplicaSetWithFailovergroup(data, "GP_Gen5_2", 5),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("collation").HasValue("SQL_AltDiction_CP850_CI_AI"),
				check.That(data.ResourceName).Key("license_type").HasValue("BasePrice"),
				check.That(data.ResourceName).Key("sku_name").HasValue("GP_Gen5_2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_createRestoreMode(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.createRestoreMode(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("creation_source_database_id"),

		{
			PreConfig: func() { time.Sleep(8 * time.Minute) },
			Config:    r.createRestoreModeDBDeleted(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},

		data.ImportStep(),

		{
			PreConfig: func() { time.Sleep(8 * time.Minute) },
			Config:    r.createRestoreModeDBRestored(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That("azurerm_mssql_database.restore").ExistsInAzure(r),
			),
		},

		data.ImportStep("restore_dropped_database_id"),
	})
}

func TestAccMsSqlDatabase_storageAccountType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccountTypeLocal(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account_type").HasValue("Local"),
			),
		},
		data.ImportStep("sample_name"),
	})
}

func TestAccMsSqlDatabase_threatDetectionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threatDetectionPolicy(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Enabled"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.retention_days").HasValue("15"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.disabled_alerts.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.email_account_admins").HasValue("Enabled"),
			),
		},
		data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.threatDetectionPolicy(data, "Disabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.state").HasValue("Disabled"),
			),
		},
		data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
	})
}

func TestAccMsSqlDatabase_threatDetectionPolicyNoStorage(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.threatDetectionPolicyNoStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.storage_account_access_key").IsEmpty(),
				check.That(data.ResourceName).Key("threat_detection_policy.0.storage_endpoint").IsEmpty(),
			),
		},
		data.ImportStep("sample_name", "threat_detection_policy.0.storage_account_access_key"),
		{
			Config: r.threatDetectionPolicy(data, "Enabled"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.storage_account_access_key").IsSet(),
				check.That(data.ResourceName).Key("threat_detection_policy.0.storage_endpoint").IsSet(),
			),
		},
		{
			Config: r.threatDetectionPolicyNoStorage(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("threat_detection_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("threat_detection_policy.0.storage_account_access_key").IsEmpty(),
				check.That(data.ResourceName).Key("threat_detection_policy.0.storage_endpoint").IsEmpty(),
			),
		},
	})
}

func TestAccMsSqlDatabase_updateSku(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSku(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateSku2(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_minCapacity0(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.minCapacity0(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_withLongTermRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLongTermRetentionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLongTermRetentionPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withLongTermRetentionPolicyNoWeekOfYear(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_withShortTermRetentionPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withShortTermRetentionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withShortTermRetentionPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMsSqlDatabase_geoBackupPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withGeoBackupPoliciesDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.withGeoBackupPoliciesEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("geo_backup_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_transparentDataEncryptionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	// NOTE: You can only update TDE on DW SKU's...
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.transparentDataEncryptionUpdate(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("transparent_data_encryption_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.transparentDataEncryptionUpdate(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("transparent_data_encryption_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
		{
			Config: r.transparentDataEncryptionUpdate(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("transparent_data_encryption_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_errorOnDisabledEncryption(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config:      r.errorOnDisabledEncryption(data),
			ExpectError: regexp.MustCompile("transparent data encryption can only be disabled on Data Warehouse SKUs"),
		},
	})
}

func TestAccMsSqlDatabase_ledgerEnabled(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ledgerEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_bacpac(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.bacpac(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccMsSqlDatabase_enclaveType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enclaveType(data, `  enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_enclaveTypeUpdate(t *testing.T) {
	// NOTE: Once the enclave_type field has be set it cannot be changed...
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.enclaveType(data, `  enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enclaveType(data, `  enclave_type = "Default"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("Default"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enclaveType(data, `  enclave_type = "VBS"`),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").HasValue("VBS"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_elasticPoolEnclaveTypeError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPool(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("elastic_pool_id").Exists(),
				check.That(data.ResourceName).Key("sku_name").HasValue("ElasticPool"),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config:      r.elasticPoolEnclaveTypeError(data),
			ExpectError: regexp.MustCompile("Before updating a database that belongs to an elastic pool please ensure"),
		},
	})
}

func TestAccMsSqlDatabase_transparentDataEncryptionKey(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.transparentDataEncryptionKey(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_namedReplication(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.hs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
		{
			Config: r.namedReplication(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_namedReplicationZoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.namedReplicationZoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enclave_type").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMsSqlDatabase_elasticPoolHS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_mssql_database", "test")
	r := MsSqlDatabaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.elasticPoolHS(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPoolHSWithRetentionPolicy(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.elasticPoolHSWithRetentionPolicyUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MsSqlDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseSqlDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.MSSQL.DatabasesClient.Get(ctx, *id, databases.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("SQL %s does not exist", id)
		}

		return nil, fmt.Errorf("reading SQL %s: %v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (MsSqlDatabaseResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlDatabaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) freeTier(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "Free"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "import" {
  name      = azurerm_mssql_database.test.name
  server_id = azurerm_mssql_server.test.id
}
`, r.basic(data))
}

func (r MsSqlDatabaseResource) complete(data acceptance.TestData) string {
	configName := "SQL_Default"

	switch data.Locations.Primary {
	case "eastus": // Added due to subscription quota policies...
		configName = "SQL_EastUS_DB_2"
	case "westeurope":
		configName = "SQL_WestEurope_DB_2"
	case "francecentral":
		configName = "SQL_FranceCentral_DB_1"
	}

	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_mssql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 10
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"
  enclave_type = "VBS"

  maintenance_configuration_name = "%[3]s"
  storage_account_type           = "Local"

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger, configName)
}

func (r MsSqlDatabaseResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_mssql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "LicenseIncluded"
  max_size_gb  = 2
  sku_name     = "GP_Gen5_2"
  enclave_type = "Default"

  storage_account_type = "Zone"

  tags = {
    ENV = "Staging"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPool(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_mssql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPoolEnclaveTypeError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_mssql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"
  enclave_type    = "VBS"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPoolDisassociation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name
  max_size_gb         = 5

  sku {
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
    capacity = 4
    family   = "Gen5"
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "GP_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "GP_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gpServerless(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name                        = "acctest-db-%[2]d"
  server_id                   = azurerm_mssql_server.test.id
  auto_pause_delay_in_minutes = 42
  min_capacity                = 0.75
  sku_name                    = "GP_S_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gpServerlessUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name                        = "acctest-db-%[2]d"
  server_id                   = azurerm_mssql_server.test.id
  auto_pause_delay_in_minutes = 90
  min_capacity                = 1.25
  sku_name                    = "GP_S_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gpWithLicenseType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_mssql_server.test.id
  sku_name     = "GP_Gen5_2"
  license_type = "LicenseIncluded"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) gpServerlessWithNullLicenseType(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "GP_S_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hs(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_mssql_server.test.id
  read_replica_count = 2
  sku_name           = "HS_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_mssql_server.test.id
  read_replica_count = 4
  sku_name           = "HS_Gen5_2"


}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hsWithRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_mssql_server.test.id
  read_replica_count = 2
  sku_name           = "HS_Gen5_2"

  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }

  short_term_retention_policy {
    retention_days = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hsWithLongRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_mssql_server.test.id
  read_replica_count = 2
  sku_name           = "HS_Gen5_2"

  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 2
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) s0WithRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "S0"

  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }

  short_term_retention_policy {
    retention_days = 10
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) hsWithRetentionPolicyUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_mssql_server.test.id
  read_replica_count = 4
  sku_name           = "HS_Gen5_2"

  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 2
  }

  short_term_retention_policy {
    retention_days = 12
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) bc(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%[2]d"
  server_id      = azurerm_mssql_server.test.id
  read_scale     = true
  sku_name       = "BC_Gen5_2"
  zone_redundant = true
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) bcUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%[2]d"
  server_id      = azurerm_mssql_server.test.id
  read_scale     = false
  sku_name       = "BC_Gen5_2"
  zone_redundant = false
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) createCopyMode(data acceptance.TestData, enclaveType string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "copy" {
  name                        = "acctest-dbc-%[2]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Copy"
  creation_source_database_id = azurerm_mssql_database.test.id
  %[3]s
}
`, r.complete(data), data.RandomInteger, enclaveType)
}

func (r MsSqlDatabaseResource) createPITRMode(data acceptance.TestData, restorePointInTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "pitr" {
  name                        = "acctest-dbp-%[2]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "PointInTimeRestore"
  restore_point_in_time       = "%[3]s"
  creation_source_database_id = azurerm_mssql_database.test.id

}
`, r.basic(data), data.RandomInteger, restorePointInTime)
}

func (r MsSqlDatabaseResource) createSecondaryMode(data acceptance.TestData, tag string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "second" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_mssql_server" "second" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.second.name
  location                     = azurerm_resource_group.second.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs-%[2]d"
  server_id                   = azurerm_mssql_server.second.id
  create_mode                 = "Secondary"
  creation_source_database_id = azurerm_mssql_database.test.id
  enclave_type                = "VBS"

  tags = {
    tag = "%[4]s"
  }
}
`, r.complete(data), data.RandomInteger, data.Locations.Secondary, tag)
}

func (r MsSqlDatabaseResource) createOnlineSecondaryMode(data acceptance.TestData, tag string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_resource_group" "second" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_mssql_server" "second" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.second.name
  location                     = azurerm_resource_group.second.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs-%[2]d"
  server_id                   = azurerm_mssql_server.second.id
  create_mode                 = "OnlineSecondary"
  creation_source_database_id = azurerm_mssql_database.test.id
  enclave_type                = "VBS"

  tags = {
    tag = "%[4]s"
  }
}
`, r.complete(data), data.RandomInteger, data.Locations.Secondary, tag)
}

func (r MsSqlDatabaseResource) scaleReplicaSet(data acceptance.TestData, sku string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "primary" {
  name        = "acctest-db-%[2]d"
  server_id   = azurerm_mssql_server.test.id
  sample_name = "AdventureWorksLT"

  max_size_gb = "2"
  sku_name    = "%[4]s"
}

resource "azurerm_resource_group" "secondary" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_mssql_server" "secondary" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.secondary.name
  location                     = azurerm_resource_group.secondary.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog12"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-db-%[2]d"
  server_id                   = azurerm_mssql_server.secondary.id
  create_mode                 = "Secondary"
  creation_source_database_id = azurerm_mssql_database.primary.id

  sku_name = "%[4]s"
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, sku)
}

func (r MsSqlDatabaseResource) scaleReplicaSetWithFailovergroup(data acceptance.TestData, sku string, size int) string {
	return fmt.Sprintf(`
	%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_mssql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = %[5]d
  sample_name  = "AdventureWorksLT"
  sku_name     = "%[4]s"

  tags = {
    ENV = "Test"
  }
}

resource "azurerm_resource_group" "second" {
  name     = "acctestRG-mssql2-%[2]d"
  location = "%[3]s"
}

resource "azurerm_mssql_server" "second" {
  name                         = "acctest-sqlserver2-%[2]d"
  resource_group_name          = azurerm_resource_group.second.name
  location                     = azurerm_resource_group.second.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-db-%[2]d"
  server_id                   = azurerm_mssql_server.second.id
  create_mode                 = "Secondary"
  creation_source_database_id = azurerm_mssql_database.test.id
  sku_name                    = "%[4]s"
}

resource "azurerm_mssql_failover_group" "failover_group" {
  name      = "acctest-fog-%[2]d"
  server_id = azurerm_mssql_server.test.id
  databases = [azurerm_mssql_database.test.id]

  partner_server {
    id = azurerm_mssql_server.second.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }

  depends_on = [
    azurerm_mssql_database.test,
    azurerm_mssql_database.secondary
  ]
}
`, r.template(data), data.RandomInteger, data.Locations.Secondary, sku, size)
}

func (MsSqlDatabaseResource) createRestoreMode(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_database" "copy" {
  name                        = "acctest-dbc-%[1]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Copy"
  creation_source_database_id = azurerm_mssql_database.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlDatabaseResource) createRestoreModeDBDeleted(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MsSqlDatabaseResource) createRestoreModeDBRestored(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-mssql-%[1]d"
  location = "%[2]s"
}

resource "azurerm_mssql_server" "test" {
  name                         = "acctest-sqlserver-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}


resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[1]d"
  server_id = azurerm_mssql_server.test.id
}

resource "azurerm_mssql_database" "restore" {
  name                        = "acctest-dbr-%[1]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Restore"
  restore_dropped_database_id = azurerm_mssql_server.test.restorable_dropped_database_ids[0]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r MsSqlDatabaseResource) storageAccountTypeLocal(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id

  storage_account_type = "Local"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) threatDetectionPolicy(data acceptance.TestData, state string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "test%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_mssql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 1
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"

  threat_detection_policy {
    retention_days             = 15
    state                      = "%[3]s"
    disabled_alerts            = ["Sql_Injection"]
    email_account_admins       = "Enabled"
    storage_account_access_key = azurerm_storage_account.test.primary_access_key
    storage_endpoint           = azurerm_storage_account.test.primary_blob_endpoint
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger, state)
}

func (r MsSqlDatabaseResource) threatDetectionPolicyNoStorage(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name         = "acctest-db-%[2]d"
  server_id    = azurerm_mssql_server.test.id
  collation    = "SQL_AltDiction_CP850_CI_AI"
  license_type = "BasePrice"
  max_size_gb  = 1
  sample_name  = "AdventureWorksLT"
  sku_name     = "GP_Gen5_2"

  threat_detection_policy {
    retention_days       = 15
    state                = "Enabled"
    disabled_alerts      = ["Sql_Injection"]
    email_account_admins = "Enabled"
  }

  tags = {
    ENV = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) updateSku(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "HS_Gen5_2"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) updateSku2(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id
  sku_name  = "HS_Gen5_4"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) minCapacity0(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id

  min_capacity = 0
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withLongTermRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_mssql_server.test.id
  long_term_retention_policy {
    weekly_retention  = "P1W"
    monthly_retention = "P1M"
    yearly_retention  = "P1Y"
    week_of_year      = 1
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withLongTermRetentionPolicyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_mssql_server.test.id
  long_term_retention_policy {
    weekly_retention = "P1W"
    yearly_retention = "P1Y"
    week_of_year     = 2
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withLongTermRetentionPolicyNoWeekOfYear(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_mssql_server.test.id
  long_term_retention_policy {
    weekly_retention = "P10D"
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withShortTermRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_mssql_server.test.id
  short_term_retention_policy {
    retention_days           = 8
    backup_interval_in_hours = 12
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withShortTermRetentionPolicyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account" "test2" {
  name                     = "acctest2%[2]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[3]d"
  server_id = azurerm_mssql_server.test.id
  short_term_retention_policy {
    retention_days           = 10
    backup_interval_in_hours = 24
  }
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withGeoBackupPoliciesEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[3]d"
  server_id          = azurerm_mssql_server.test.id
  sku_name           = "DW100c"
  geo_backup_enabled = true
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) withGeoBackupPoliciesDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[3]d"
  server_id          = azurerm_mssql_server.test.id
  sku_name           = "DW100c"
  geo_backup_enabled = false
}
`, r.template(data), data.RandomIntOfLength(15), data.RandomInteger)
}

func (r MsSqlDatabaseResource) transparentDataEncryptionUpdate(data acceptance.TestData, state bool) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name                                = "acctest-db-%d"
  server_id                           = azurerm_mssql_server.test.id
  sku_name                            = "DW100c"
  transparent_data_encryption_enabled = %t
}
`, r.template(data), data.RandomInteger, state)
}

func (r MsSqlDatabaseResource) errorOnDisabledEncryption(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_mssql_database" "test" {
  name                                = "acctest-db-%d"
  server_id                           = azurerm_mssql_server.test.id
  transparent_data_encryption_enabled = false
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) ledgerEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%[2]d"
  server_id      = azurerm_mssql_server.test.id
  ledger_enabled = true
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) bacpac(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "accsa%d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "bacpac"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "test" {
  name                   = "test.bacpac"
  storage_account_name   = azurerm_storage_account.test.name
  storage_container_name = azurerm_storage_container.test.name
  type                   = "Block"
  source                 = "testdata/sql_import.bacpac"
}

resource "azurerm_mssql_firewall_rule" "test" {
  name             = "allowazure"
  server_id        = azurerm_mssql_server.test.id
  start_ip_address = "0.0.0.0"
  end_ip_address   = "0.0.0.0"
}

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id

  import {
    storage_uri                  = azurerm_storage_blob.test.url
    storage_key                  = azurerm_storage_account.test.primary_access_key
    storage_key_type             = "StorageAccessKey"
    administrator_login          = azurerm_mssql_server.test.administrator_login
    administrator_login_password = azurerm_mssql_server.test.administrator_login_password
    authentication_type          = "Sql"
  }

  timeouts {
    create = "10h"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) enclaveType(data acceptance.TestData, enclaveType string) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name      = "acctest-db-%[2]d"
  server_id = azurerm_mssql_server.test.id

  %[3]s
}
`, r.template(data), data.RandomInteger, enclaveType)
}

func (r MsSqlDatabaseResource) transparentDataEncryptionKey(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_client_config" "test" {}

resource "azurerm_user_assigned_identity" "test" {
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  name                = "test_identity"
}

resource "azurerm_key_vault" "test" {
  name                        = "vault%[2]d"
  location                    = azurerm_resource_group.test.location
  resource_group_name         = azurerm_resource_group.test.name
  enabled_for_disk_encryption = true
  tenant_id                   = azurerm_user_assigned_identity.test.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.test.tenant_id
    object_id = data.azurerm_client_config.test.object_id

    key_permissions = ["Get", "List", "Create", "Delete", "Update", "Recover", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.test.tenant_id
    object_id = azurerm_user_assigned_identity.test.principal_id

    key_permissions = ["Get", "WrapKey", "UnwrapKey"]
  }
}

resource "azurerm_key_vault_key" "test" {
  depends_on = [azurerm_key_vault.test]

  name         = "key-%[3]s"
  key_vault_id = azurerm_key_vault.test.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["unwrapKey", "wrapKey"]
}

resource "azurerm_mssql_database" "test" {
  name                                                       = "acctest-db-%[2]d"
  server_id                                                  = azurerm_mssql_server.test.id
  sku_name                                                   = "S0"
  transparent_data_encryption_enabled                        = true
  transparent_data_encryption_key_vault_key_id               = azurerm_key_vault_key.test.id
  transparent_data_encryption_key_automatic_rotation_enabled = true

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r MsSqlDatabaseResource) namedReplication(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs2-%[2]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Secondary"
  secondary_type              = "Named"
  creation_source_database_id = azurerm_mssql_database.test.id
}
`, r.hs(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) namedReplicationZoneRedundant(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_database" "test" {
  name               = "acctest-db-%[2]d"
  server_id          = azurerm_mssql_server.test.id
  read_replica_count = 2
  sku_name           = "HS_Gen5_2"

  storage_account_type = "Zone"
  zone_redundant       = true
}

resource "azurerm_mssql_database" "secondary" {
  name                        = "acctest-dbs2-%[2]d"
  server_id                   = azurerm_mssql_server.test.id
  create_mode                 = "Secondary"
  secondary_type              = "Named"
  creation_source_database_id = azurerm_mssql_database.test.id

  zone_redundant     = true
  read_replica_count = 1
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPoolHS(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name

  sku {
    name     = "HS_Gen5"
    tier     = "Hyperscale"
    family   = "Gen5"
    capacity = 4
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_mssql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"
}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPoolHSWithRetentionPolicy(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name

  sku {
    name     = "HS_Gen5"
    tier     = "Hyperscale"
    family   = "Gen5"
    capacity = 4
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_mssql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"

  short_term_retention_policy {
    retention_days = 10
  }

}
`, r.template(data), data.RandomInteger)
}

func (r MsSqlDatabaseResource) elasticPoolHSWithRetentionPolicyUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_mssql_elasticpool" "test" {
  name                = "acctest-pool-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  server_name         = azurerm_mssql_server.test.name

  sku {
    name     = "HS_Gen5"
    tier     = "Hyperscale"
    family   = "Gen5"
    capacity = 4
  }

  per_database_settings {
    min_capacity = 0.25
    max_capacity = 4
  }
}

resource "azurerm_mssql_database" "test" {
  name            = "acctest-db-%[2]d"
  server_id       = azurerm_mssql_server.test.id
  elastic_pool_id = azurerm_mssql_elasticpool.test.id
  sku_name        = "ElasticPool"

  short_term_retention_policy {
    retention_days = 12
  }

}
`, r.template(data), data.RandomInteger)
}
