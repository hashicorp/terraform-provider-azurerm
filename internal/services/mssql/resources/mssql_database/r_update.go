// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/backupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/elasticpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/longtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/replicationlinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/transparentdataencryptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	serversClient := meta.(*clients.Client).MSSQL.ServersClient
	securityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	elasticPoolClient := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	replicationLinksClient := meta.(*clients.Client).MSSQL.ReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	skuName := d.Get("sku_name").(string)
	elasticPoolId := d.Get("elastic_pool_id").(string)
	createMode := d.Get("create_mode").(string)
	restorePointInTime := d.Get("restore_point_in_time").(string)

	// Determine whether the SKU is for SQL Data Warehouse
	isDwSku := strings.HasPrefix(strings.ToLower(skuName), "dw")

	if strings.HasPrefix(skuName, "GP_S_") && !pluginsdk.IsExplicitlyNullInConfig(d, "license_type") {
		return fmt.Errorf("serverless databases do not support license type")
	}

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID: %+v", err)
	}

	id := commonids.NewSqlDatabaseID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, name)

	existing, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+q", id, err)
	}

	if _, err = serversClient.Get(ctx, pointer.From(serverId), servers.DefaultGetOperationOptions()); err != nil {
		return fmt.Errorf("retrieving %s: %q", serverId, err)
	}

	// when disassociating mssql db from elastic pool, the sku_name must be specific
	if d.HasChange("elastic_pool_id") {
		if old, new := d.GetChange("elastic_pool_id"); old.(string) != "" && new.(string) == "" {
			if skuName == "" || skuName == "ElasticPool" {
				return fmt.Errorf("`sku_name` must be assigned and not be %q when disassociating from Elastic Pool", "ElasticPool")
			}
		}
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	payload := databases.DatabaseUpdate{}
	props := databases.DatabaseUpdateProperties{}

	if d.HasChange("auto_pause_delay_in_minutes") {
		props.AutoPauseDelay = pointer.To(int64(d.Get("auto_pause_delay_in_minutes").(int)))
	}

	if d.HasChange("elastic_pool_id") {
		props.ElasticPoolId = pointer.To(d.Get("elastic_pool_id").(string))
	}

	if d.HasChange("license_type") {
		props.LicenseType = pointer.ToEnum[databases.DatabaseLicenseType](d.Get("license_type").(string))
	}

	if d.HasChange("min_capacity") {
		props.MinCapacity = pointer.To(d.Get("min_capacity").(float64))
	}

	if d.HasChange("read_replica_count") {
		props.HighAvailabilityReplicaCount = pointer.To(int64(d.Get("read_replica_count").(int)))
	}

	if d.HasChange("sample_name") {
		props.SampleName = pointer.ToEnum[databases.SampleName](d.Get("sample_name").(string))
	}

	if d.HasChange("storage_account_type") {
		props.RequestedBackupStorageRedundancy = pointer.ToEnum[databases.BackupStorageRedundancy](d.Get("storage_account_type").(string))
	}

	if d.HasChange("zone_redundant") {
		props.ZoneRedundant = pointer.To(d.Get("zone_redundant").(bool))
	}

	if d.HasChange("enclave_type") {
		var enclaveType databases.AlwaysEncryptedEnclaveType
		if v, ok := d.GetOk("enclave_type"); ok && v.(string) != "" {
			enclaveType = databases.AlwaysEncryptedEnclaveType(v.(string))
		}

		// The 'PreferredEnclaveType' field cannot be passed to the APIs Update if the
		// 'sku_name' is a DW or DC-series SKU...
		if !strings.HasPrefix(strings.ToLower(skuName), "dw") && !strings.Contains(strings.ToLower(skuName), "_dc_") && enclaveType != "" {
			props.PreferredEnclaveType = pointer.To(enclaveType)
		} else {
			props.PreferredEnclaveType = nil
		}

		// If the database belongs to an elastic pool, we need to GET the elastic pool and check
		// if the updated 'enclave_type' matches the existing elastic pools 'enclave_type'. If they don't
		// we need to raise an error stating that they must match.
		if elasticPoolId != "" {
			elasticId, err := commonids.ParseSqlElasticPoolID(elasticPoolId)
			if err != nil {
				return err
			}

			elasticPool, err := elasticPoolClient.Get(ctx, *elasticId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %s", elasticId, err)
			}

			var elasticEnclaveType elasticpools.AlwaysEncryptedEnclaveType
			if elasticPool.Model != nil && elasticPool.Model.Properties != nil && elasticPool.Model.Properties.PreferredEnclaveType != nil {
				elasticEnclaveType = pointer.From(elasticPool.Model.Properties.PreferredEnclaveType)
			}

			if elasticEnclaveType != "" || enclaveType != "" {
				if !strings.EqualFold(string(elasticEnclaveType), string(enclaveType)) {
					return fmt.Errorf("updating the %s with enclave type %q to the %s with enclave type %q is not supported. Before updating a database that belongs to an elastic pool please ensure that the 'enclave_type' is the same for both the database and the elastic pool", id, enclaveType, elasticId, elasticEnclaveType)
				}
			}
		}
	}

	// we should not specify the value of `maintenance_configuration_name` when `elastic_pool_id` is set since its value depends on the elastic pool's `maintenance_configuration_name` value.
	if elasticPoolId == "" && d.HasChange("maintenance_configuration_name") {
		// set default value here because `elastic_pool_id` is not specified, API returns default value `SQL_Default` for `maintenance_configuration_name`
		maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, "SQL_Default")
		if v, ok := d.GetOk("maintenance_configuration_name"); ok {
			maintenanceConfigId = publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, v.(string))
		}

		props.MaintenanceConfigurationId = pointer.To(maintenanceConfigId.ID())
	}

	if v, ok := d.GetOk("max_size_gb"); ok {
		// `max_size_gb` is Computed, so has a value after the first run
		v, err := calculateMaxSizeBytes(v.(float64))
		if err != nil {
			return err
		}
		props.MaxSizeBytes = v

		// `max_size_gb` only has change if it is configured
		if d.HasChange("max_size_gb") && (createMode == string(databases.CreateModeOnlineSecondary) || createMode == string(databases.CreateModeSecondary)) {
			return fmt.Errorf("it is not possible to change maximum size nor advised to configure maximum size in secondary create mode for %s", id)
		}
	}

	if d.HasChanges("read_scale") {
		readScale := databases.DatabaseReadScaleDisabled
		if v := d.Get("read_scale").(bool); v {
			readScale = databases.DatabaseReadScaleEnabled
		}
		props.ReadScale = pointer.To(readScale)
	}

	if d.HasChange("restore_point_in_time") {
		if restorePointInTime != "" {
			if createMode != string(databases.CreateModePointInTimeRestore) {
				return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(databases.CreateModePointInTimeRestore))
			}
			props.RestorePointInTime = pointer.To(restorePointInTime)
		}
	}

	if d.HasChange("sku_name") {
		// When databases are replicating, the primary cannot have a SKU belonging to a higher service tier than any of its
		// partner databases. To work around this, we'll try to identify any partner databases that are secondary to this
		// database, and where the new SKU tier for this database is going to be higher, first upgrade those databases to
		// the same sku_name as we'll be changing this database to. If that sku is different to the one configured for any
		// of the partner databases, that discrepancy will have to be corrected by the resource for that database. That
		// might happen as part of the same apply, if a change was already planned for it, else it will only be picked up
		// in a second plan/apply.
		//
		// TLDR: for the best experience, configs should use the same SKU for primary and partner databases and when
		// upgrading those SKUs, we'll try to upgrade the partner databases first.

		// Place a lock for the current database so any partner resources can't bump its SKU out of band
		if skuName != "" {
			var existingEnclaveType databases.AlwaysEncryptedEnclaveType
			if model := existing.Model; model != nil && model.Properties != nil && model.Properties.PreferredEnclaveType != nil {
				existingEnclaveType = *model.Properties.PreferredEnclaveType
			}

			partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, client, replicationLinksClient, resourcesClient, id, existingEnclaveType, []replicationlinks.ReplicationRole{replicationlinks.ReplicationRoleSecondary, replicationlinks.ReplicationRoleNonReadableSecondary})
			if err != nil {
				return err
			}

			log.Printf("[INFO] Found %d Partner Databases", len(partnerDatabases))

			// Place a lock for the partner databases, so they can't update themselves whilst we're poking their SKUs
			for _, v := range partnerDatabases {
				id, err := commonids.ParseSqlDatabaseIDInsensitively(pointer.From(v.Id))
				if err != nil {
					return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", id.ID(), err)
				}

				locks.ByID(id.ID())
				defer locks.UnlockByID(id.ID())
			}

			// Update the SKUs of any partner databases where deemed necessary
			for _, partnerDatabase := range partnerDatabases {
				log.Printf("[INFO] Parsing Replication Partner Database ID: %s", *partnerDatabase.Id)
				partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.Id)
				if err != nil {
					return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.Id, err)
				}

				// See: https://docs.microsoft.com/azure/azure-sql/database/active-geo-replication-overview#configuring-secondary-database
				if partnerDatabase.Sku != nil && partnerDatabase.Sku.Name != "" && helper.CompareDatabaseSkuServiceTiers(skuName, partnerDatabase.Sku.Name) {
					if err := client.UpdateThenPoll(ctx, *partnerDatabaseId, databases.DatabaseUpdate{
						Sku: &databases.Sku{
							Name: skuName,
						},
					}); err != nil {
						return fmt.Errorf("updating SKU of Replication Partner Database %s: %+v", partnerDatabaseId, err)
					}
				}
			}
		}

		payload.Sku = pointer.To(databases.Sku{
			Name: skuName,
		})
	}

	if d.HasChange("recover_database_id") {
		props.RecoverableDatabaseId = pointer.To(d.Get("recover_database_id").(string))
	}

	if d.HasChange("recovery_point_id") {
		props.RecoveryServicesRecoveryPointId = pointer.To(d.Get("recovery_point_id").(string))
	}

	if d.HasChange("restore_dropped_database_id") {
		props.RestorableDroppedDatabaseId = pointer.To(d.Get("restore_dropped_database_id").(string))
	}

	if d.HasChange("restore_long_term_retention_backup_id") {
		props.LongTermRetentionBackupResourceId = pointer.To(d.Get("restore_long_term_retention_backup_id").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("identity") {
		expanded, err := identity.ExpandUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expanded
	}

	if d.HasChange("transparent_data_encryption_key_vault_key_id") {
		keyVaultKeyId := d.Get("transparent_data_encryption_key_vault_key_id").(string)

		keyId, err := keyvault.ParseNestedItemID(keyVaultKeyId, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return fmt.Errorf("unable to parse key: %q: %+v", keyVaultKeyId, err)
		}

		props.EncryptionProtector = pointer.To(keyId.ID())
	}

	if d.HasChange("transparent_data_encryption_key_automatic_rotation_enabled") {
		v, ok := d.GetOk("transparent_data_encryption_key_automatic_rotation_enabled")
		if ok && !v.(bool) && isDwSku {
			props.EncryptionProtectorAutoRotation = nil
		} else if !isDwSku {
			props.EncryptionProtectorAutoRotation = pointer.To(v.(bool))
		}
	}

	payload.Properties = pointer.To(props)
	if err = client.UpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Wait for the ProvisioningState to become "Succeeded"
	log.Printf("[DEBUG] Waiting for %s to become ready", id)
	pendingStatuses := make([]string, 0)
	for _, s := range databases.PossibleValuesForDatabaseStatus() {
		if s != string(databases.DatabaseStatusOnline) {
			pendingStatuses = append(pendingStatuses, s)
		}
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	// NOTE: Internal x-ref, this is another case of hashicorp/go-azure-sdk#307 so this can be removed once that's fixed
	stateConf := &pluginsdk.StateChangeConf{
		Pending: pendingStatuses,
		Target:  []string{string(databases.DatabaseStatusOnline)},
		Refresh: func() (any, string, error) {
			log.Printf("[DEBUG] Checking to see if %s is online...", id)

			resp, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
			if err != nil {
				return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					return resp, pointer.FromEnum(props.Status), nil
				}
			}

			return resp.Model, "", nil
		},

		ContinuousTargetOccurence: 2,
		MinTimeout:                1 * time.Minute,
		Timeout:                   time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}

	// Cannot set transparent data encryption for secondary databases
	if createMode != string(databases.CreateModeOnlineSecondary) && createMode != string(databases.CreateModeSecondary) {
		state := transparentdataencryptions.TransparentDataEncryptionStateDisabled
		if d.HasChange("transparent_data_encryption_enabled") {
			if v := d.Get("transparent_data_encryption_enabled").(bool); v {
				state = transparentdataencryptions.TransparentDataEncryptionStateEnabled
			}

			input := transparentdataencryptions.LogicalDatabaseTransparentDataEncryption{
				Properties: pointer.To(transparentdataencryptions.TransparentDataEncryptionProperties{
					State: state,
				}),
			}

			if err := transparentEncryptionClient.CreateOrUpdateThenPoll(ctx, id, input); err != nil {
				return fmt.Errorf("while updating Transparent Data Encryption state for %s: %+v", id, err)
			}

			// NOTE: Internal x-ref, this is another case of hashicorp/go-azure-sdk#307 so this can be removed once that's fixed
			if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
				c, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
				if err != nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("while polling %s for status: %+v", id.String(), err))
				}

				if model := c.Model; model != nil && model.Properties != nil && model.Properties.Status != nil {
					if model.Properties.Status == pointer.To(databases.DatabaseStatusScaling) {
						return pluginsdk.RetryableError(fmt.Errorf("database %s is still scaling", id.String()))
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}
	}

	if d.HasChange("import") {
		if _, ok := d.GetOk("import"); ok {
			importParameters := expandMsSqlServerImport(d)

			if err := client.ImportThenPoll(ctx, id, importParameters); err != nil {
				return fmt.Errorf("while importing the BACPAC file into the new database %s: %+v", id.ID(), err)
			}
		}
	}

	// For datawarehouse SKUs only
	if isDwSku && d.HasChange("geo_backup_enabled") {
		isEnabled := d.Get("geo_backup_enabled").(bool)
		var geoBackupPolicyState geobackuppolicies.GeoBackupPolicyState

		geoBackupPolicyState = geobackuppolicies.GeoBackupPolicyStateDisabled
		if isEnabled {
			geoBackupPolicyState = geobackuppolicies.GeoBackupPolicyStateEnabled
		}

		geoBackupPolicy := geobackuppolicies.GeoBackupPolicy{
			Properties: &geobackuppolicies.GeoBackupPolicyProperties{
				State: geoBackupPolicyState,
			},
		}

		if _, err := geoBackupPoliciesClient.CreateOrUpdate(ctx, id, geoBackupPolicy); err != nil {
			return fmt.Errorf("setting Geo Backup Policies for %s: %+v", id, err)
		}
	}

	if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		result, err := securityAlertPoliciesClient.CreateOrUpdate(ctx, id, expandMsSqlDatabaseSecurityAlertPolicy(d))

		if response.WasNotFound(result.HttpResponse) {
			return pluginsdk.RetryableError(fmt.Errorf("database %s is still creating", id.String()))
		}

		if err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("setting database threat detection policy for %s: %+v", id, err))
		}

		return nil
	}); err != nil {
		return err
	}

	if d.HasChange("long_term_retention_policy") {
		v := d.Get("long_term_retention_policy")
		longTermRetentionProps := helper.ExpandLongTermRetentionPolicy(v.([]any))
		if longTermRetentionProps != nil {
			longTermRetentionPolicy := longtermretentionpolicies.LongTermRetentionPolicy{}

			// DataWarehouse SKUs do not support LRP currently
			if !isDwSku {
				longTermRetentionPolicy.Properties = longTermRetentionProps
			}

			if err := longTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, longTermRetentionPolicy); err != nil {
				return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("short_term_retention_policy") {
		v := d.Get("short_term_retention_policy")
		backupShortTermPolicyProps := helper.ExpandShortTermRetentionPolicy(v.([]any))
		if backupShortTermPolicyProps != nil {
			backupShortTermPolicy := backupshorttermretentionpolicies.BackupShortTermRetentionPolicy{}

			if !isDwSku {
				backupShortTermPolicy.Properties = backupShortTermPolicyProps
			}

			elasticPoolSku := ""
			if elasticPoolId != "" {
				elasticId, err := commonids.ParseSqlElasticPoolID(elasticPoolId)
				if err != nil {
					return err
				}

				elasticPool, err := elasticPoolClient.Get(ctx, *elasticId)
				if err != nil {
					return fmt.Errorf("retrieving %s: %v", elasticId, err)
				}

				if elasticPool.Model != nil && elasticPool.Model.Sku != nil {
					elasticPoolSku = elasticPool.Model.Sku.Name
				}
			}

			if strings.HasPrefix(skuName, "HS") || strings.HasPrefix(elasticPoolSku, "HS") {
				backupShortTermPolicy.Properties.DiffBackupIntervalInHours = nil
			}

			if err := shortTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, backupShortTermPolicy); err != nil {
				return fmt.Errorf("setting Short Term Retention Policies for %s: %+v", id, err)
			}
		}
	}

	return resourceMsSqlDatabaseRead(d, meta)
}
