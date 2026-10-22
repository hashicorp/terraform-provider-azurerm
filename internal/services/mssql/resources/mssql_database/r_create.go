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
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/databasesecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/longtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/replicationlinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/transparentdataencryptions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMsSqlDatabaseCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	serversClient := meta.(*clients.Client).MSSQL.ServersClient
	elasticPoolClient := meta.(*clients.Client).MSSQL.ElasticPoolsClient
	databaseSecurityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	replicationLinksClient := meta.(*clients.Client).MSSQL.ReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if strings.HasPrefix(d.Get("sku_name").(string), "GP_S_") && d.Get("license_type").(string) != "" {
		return fmt.Errorf("serverless databases do not support license type")
	}

	name := d.Get("name").(string)

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID: %+v", err)
	}

	id := commonids.NewSqlDatabaseID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, name)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		if existing, err := client.Get(ctx, id, databases.DefaultGetOperationOptions()); err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		} else {
			return tf.ImportAsExistsError("azurerm_mssql_database", id.ID())
		}
	}

	server, err := serversClient.Get(ctx, *serverId, servers.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %q", serverId, err)
	}

	if server.Model == nil {
		return fmt.Errorf("server model was nil")
	}

	if server.Model.Location == "" {
		return fmt.Errorf("reading %s: Location was empty", serverId)
	}

	location := server.Model.Location
	ledgerEnabled := d.Get("ledger_enabled").(bool)

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
	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	// NOTE: The service default is actually nil/empty which indicates enclave is disabled. the value `Default` is NOT the default.
	var enclaveType databases.AlwaysEncryptedEnclaveType
	if v, ok := d.GetOk("enclave_type"); ok && v.(string) != "" {
		enclaveType = databases.AlwaysEncryptedEnclaveType(v.(string))
	}

	skuName := d.Get("sku_name").(string)

	if skuName != "" {
		partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, client, replicationLinksClient, resourcesClient, id, enclaveType, []replicationlinks.ReplicationRole{replicationlinks.ReplicationRoleSecondary, replicationlinks.ReplicationRoleNonReadableSecondary})
		if err != nil {
			return err
		}

		// Place a lock for the partner databases, so they can't update themselves whilst we're poking their SKUs
		for _, partnerDatabase := range partnerDatabases {
			partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.Id)
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.Id, err)
			}

			locks.ByID(partnerDatabaseId.ID())
			defer locks.UnlockByID(partnerDatabaseId.ID())
		}

		// Update the SKUs of any partner databases where deemed necessary
		for _, partnerDatabase := range partnerDatabases {
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

	// Determine whether the SKU is for SQL Data Warehouse
	isDwSku := strings.HasPrefix(strings.ToLower(skuName), "dw")

	// NOTE: If the database is being added to an elastic pool, we need to GET the elastic pool and check
	// if the 'enclave_type' matches. If they don't we need to raise an error stating that they must match.
	elasticPoolId := d.Get("elastic_pool_id").(string)
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

		if elasticPool.Model != nil {
			if elasticPool.Model.Properties != nil && elasticPool.Model.Properties.PreferredEnclaveType != nil {
				elasticEnclaveType := pointer.FromEnum(elasticPool.Model.Properties.PreferredEnclaveType)
				databaseEnclaveType := string(enclaveType)

				if !strings.EqualFold(elasticEnclaveType, databaseEnclaveType) {
					return fmt.Errorf("adding the %s with enclave type %q to the %s with enclave type %q is not supported. Before adding a database to an elastic pool please ensure that the 'enclave_type' is the same for both the database and the elastic pool", id, databaseEnclaveType, elasticId, elasticEnclaveType)
				}
			}

			if elasticPool.Model.Sku != nil {
				elasticPoolSku = elasticPool.Model.Sku.Name
			}
		}
	}

	input := databases.Database{
		Location: location,
		Properties: &databases.DatabaseProperties{
			AutoPauseDelay:                   pointer.To(int64(d.Get("auto_pause_delay_in_minutes").(int))),
			Collation:                        pointer.To(d.Get("collation").(string)),
			ElasticPoolId:                    pointer.To(elasticPoolId),
			LicenseType:                      pointer.ToEnum[databases.DatabaseLicenseType](d.Get("license_type").(string)),
			MinCapacity:                      pointer.To(d.Get("min_capacity").(float64)),
			HighAvailabilityReplicaCount:     pointer.To(int64(d.Get("read_replica_count").(int))),
			SampleName:                       pointer.ToEnum[databases.SampleName](d.Get("sample_name").(string)),
			RequestedBackupStorageRedundancy: pointer.ToEnum[databases.BackupStorageRedundancy](d.Get("storage_account_type").(string)),
			ZoneRedundant:                    pointer.To(d.Get("zone_redundant").(bool)),
			IsLedgerOn:                       pointer.To(ledgerEnabled),
			SecondaryType:                    pointer.ToEnum[databases.SecondaryType](d.Get("secondary_type").(string)),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	// NOTE: The 'PreferredEnclaveType' field cannot be passed to the APIs Create if the 'sku_name' is a DW or DC-series SKU...
	if !strings.HasPrefix(strings.ToLower(skuName), "dw") && !strings.Contains(strings.ToLower(skuName), "_dc_") && enclaveType != "" {
		input.Properties.PreferredEnclaveType = pointer.To(enclaveType)
	}

	v, ok := d.GetOk("transparent_data_encryption_key_automatic_rotation_enabled")
	if ok && !v.(bool) && isDwSku {
		input.Properties.EncryptionProtectorAutoRotation = nil
	} else if !isDwSku {
		input.Properties.EncryptionProtectorAutoRotation = pointer.To(v.(bool))
	}

	createMode := d.Get("create_mode").(string)

	switch databases.CreateMode(createMode) {
	case databases.CreateModeCopy, databases.CreateModePointInTimeRestore, databases.CreateModeSecondary, databases.CreateModeOnlineSecondary:
		if creationSourceDatabaseId, dbok := d.GetOk("creation_source_database_id"); !dbok {
			return fmt.Errorf("'creation_source_database_id' is required for 'create_mode' %q", createMode)
		} else {
			// We need to make sure the enclave types match...
			primaryDatabaseId, err := commonids.ParseSqlDatabaseID(creationSourceDatabaseId.(string))
			if err != nil {
				return fmt.Errorf("parsing creation source database ID: %+v", err)
			}

			primaryDatabase, err := client.Get(ctx, *primaryDatabaseId, databases.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving creation source %s: %+v", primaryDatabaseId, err)
			}

			if model := primaryDatabase.Model; model != nil && model.Properties != nil && model.Properties.PreferredEnclaveType != nil && enclaveType != *model.Properties.PreferredEnclaveType {
				return fmt.Errorf("specifying different 'enclave_type' properties for 'create_mode' %q is not supported, primary 'enclave_type' %q does not match current 'enclave_type' %q. please ensure that the 'enclave_type' is the same for both databases", createMode, string(*model.Properties.PreferredEnclaveType), string(enclaveType))
			}
		}
	case databases.CreateModeRecovery:
		if _, dbok := d.GetOk("recover_database_id"); !dbok {
			return fmt.Errorf("'recover_database_id' is required for create_mode %s", createMode)
		}
	case databases.CreateModeRestore:
		if _, dbok := d.GetOk("restore_dropped_database_id"); !dbok {
			return fmt.Errorf("'restore_dropped_database_id' is required for create_mode %s", createMode)
		}
	case databases.CreateModeRestoreLongTermRetentionBackup:
		if _, dbok := d.GetOk("restore_long_term_retention_backup_id"); !dbok {
			return fmt.Errorf("'restore_long_term_retention_backup_id' is required for create_mode %s", createMode)
		}
	}

	// we should not specify the value of `maintenance_configuration_name` when `elastic_pool_id` is set since its value depends on the elastic pool's `maintenance_configuration_name` value.
	if _, ok := d.GetOk("elastic_pool_id"); !ok {
		// set default value here because `elastic_pool_id` is not specified, API returns default value `SQL_Default` for `maintenance_configuration_name`
		maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, "SQL_Default")
		if v, ok := d.GetOk("maintenance_configuration_name"); ok {
			maintenanceConfigId = publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, v.(string))
		}
		input.Properties.MaintenanceConfigurationId = pointer.To(maintenanceConfigId.ID())
	}

	input.Properties.CreateMode = pointer.ToEnum[databases.CreateMode](createMode)

	if v, ok := d.GetOk("max_size_gb"); ok {
		// `max_size_gb` is Computed, so has a value after the first run
		if createMode != string(databases.CreateModeOnlineSecondary) && createMode != string(databases.CreateModeSecondary) {
			v, err := calculateMaxSizeBytes(v.(float64))
			if err != nil {
				return err
			}
			input.Properties.MaxSizeBytes = v
		}

		// `max_size_gb` only has change if it is configured
		if d.HasChange("max_size_gb") && (createMode == string(databases.CreateModeOnlineSecondary) || createMode == string(databases.CreateModeSecondary)) {
			return fmt.Errorf("it is not possible to change maximum size nor advised to configure maximum size in secondary create mode for %s", id)
		}
	}

	readScale := databases.DatabaseReadScaleDisabled
	if v := d.Get("read_scale").(bool); v {
		readScale = databases.DatabaseReadScaleEnabled
	}
	input.Properties.ReadScale = pointer.To(readScale)

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		if cm, ok := d.GetOk("create_mode"); ok && cm.(string) != string(databases.CreateModePointInTimeRestore) {
			return fmt.Errorf("'restore_point_in_time' is supported only for 'create_mode' %q", string(databases.CreateModePointInTimeRestore))
		}

		input.Properties.RestorePointInTime = pointer.To(v.(string))
	}

	if skuName != "" {
		input.Sku = pointer.To(databases.Sku{
			Name: skuName,
		})
	}

	if v, ok := d.GetOk("creation_source_database_id"); ok {
		input.Properties.SourceDatabaseId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("recover_database_id"); ok {
		input.Properties.RecoverableDatabaseId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("recovery_point_id"); ok {
		input.Properties.RecoveryServicesRecoveryPointId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("restore_dropped_database_id"); ok {
		input.Properties.RestorableDroppedDatabaseId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("restore_long_term_retention_backup_id"); ok {
		input.Properties.LongTermRetentionBackupResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("identity"); ok {
		expandedIdentity, err := identity.ExpandUserAssignedMap(v.([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		input.Identity = expandedIdentity
	}

	if v, ok := d.GetOk("transparent_data_encryption_key_vault_key_id"); ok {
		keyVaultKeyId := v.(string)

		keyId, err := keyvault.ParseNestedItemID(keyVaultKeyId, keyvault.VersionTypeVersioned, keyvault.NestedItemTypeKey)
		if err != nil {
			return err
		}

		input.Properties.EncryptionProtector = pointer.To(keyId.ID())
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, input, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
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

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Status != nil {
				return resp, pointer.FromEnum(resp.Model.Properties.Status), nil
			}

			return resp, "", nil
		},
		ContinuousTargetOccurence: 2,
		MinTimeout:                1 * time.Minute,
		Timeout:                   time.Until(deadline),
	}

	// NOTE: Internal x-ref, this is another case of hashicorp/go-azure-sdk#307 so this can be removed once that's fixed
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}

	// Cannot set transparent data encryption for secondary databases
	if createMode != string(databases.CreateModeOnlineSecondary) && createMode != string(databases.CreateModeSecondary) {
		state := transparentdataencryptions.TransparentDataEncryptionStateDisabled
		if v := d.Get("transparent_data_encryption_enabled").(bool); v {
			state = transparentdataencryptions.TransparentDataEncryptionStateEnabled
		}

		tde, retryErr := transparentEncryptionClient.Get(ctx, id)
		if retryErr != nil {
			return fmt.Errorf("while retrieving Transparent Data Encryption state for %s: %+v", id, retryErr)
		}

		currentState := transparentdataencryptions.TransparentDataEncryptionStateDisabled
		if model := tde.Model; model != nil {
			if props := model.Properties; props != nil {
				currentState = props.State
			}
		}

		// Submit TDE state only when state is being changed, otherwise it can cause unwanted detection of state changes from the cloud side
		if !strings.EqualFold(string(currentState), string(state)) {
			tdePayload := transparentdataencryptions.LogicalDatabaseTransparentDataEncryption{
				Properties: &transparentdataencryptions.TransparentDataEncryptionProperties{
					State: state,
				},
			}

			if err := transparentEncryptionClient.CreateOrUpdateThenPoll(ctx, id, tdePayload); err != nil {
				return fmt.Errorf("while enabling Transparent Data Encryption for %q: %+v", id.String(), err)
			}

			// NOTE: Internal x-ref, this is another case of hashicorp/go-azure-sdk#307 so this can be removed once that's fixed
			if retryErr = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
				c, err2 := client.Get(ctx, id, databases.DefaultGetOperationOptions())
				if err2 != nil {
					return pluginsdk.NonRetryableError(fmt.Errorf("while polling %s for status: %+v", id.String(), err2))
				}
				if c.Model != nil && c.Model.Properties != nil && c.Model.Properties.Status != nil {
					if c.Model.Properties.Status == pointer.To(databases.DatabaseStatusScaling) {
						return pluginsdk.RetryableError(fmt.Errorf("database %s is still scaling", id.String()))
					}
				} else {
					return pluginsdk.RetryableError(fmt.Errorf("retrieving database status %s: Model, Properties or Status is nil", id.String()))
				}

				return nil
			}); retryErr != nil {
				return retryErr
			}
		} else {
			log.Print("[DEBUG] Skipping re-writing of Transparent Data Encryption, since encryption state is not changing ...")
		}
	}

	if _, ok := d.GetOk("import"); ok {
		importParameters := expandMsSqlServerImport(d)

		if err := client.ImportThenPoll(ctx, id, importParameters); err != nil {
			return fmt.Errorf("while import bacpac into the new database %s: %+v", id, err)
		}
	}

	// For Data Warehouse SKUs only
	if isDwSku {
		enabled := d.Get("geo_backup_enabled").(bool)

		// The default geo backup policy configuration for a new resource is 'enabled', so we don't need to set it in that scenario
		if !enabled {
			input := geobackuppolicies.GeoBackupPolicy{
				Properties: pointer.To(geobackuppolicies.GeoBackupPolicyProperties{
					State: geobackuppolicies.GeoBackupPolicyStateDisabled,
				}),
			}

			if _, err := geoBackupPoliciesClient.CreateOrUpdate(ctx, id, input); err != nil {
				return fmt.Errorf("setting Geo Backup Policies %s: %+v", id, err)
			}
		}
	}

	if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		result, err := databaseSecurityAlertPoliciesClient.CreateOrUpdate(ctx, id, expandMsSqlDatabaseSecurityAlertPolicy(d))

		if response.WasNotFound(result.HttpResponse) {
			return pluginsdk.RetryableError(fmt.Errorf("database %s is still creating", id))
		}

		if err != nil {
			return pluginsdk.NonRetryableError(fmt.Errorf("setting database threat detection policy %s: %+v", id, err))
		}

		return nil
	}); err != nil {
		return err
	}

	longTermRetentionPolicyProps := helper.ExpandLongTermRetentionPolicy(d.Get("long_term_retention_policy").([]any))
	if longTermRetentionPolicyProps != nil {
		longTermRetentionPolicyPayload := longtermretentionpolicies.LongTermRetentionPolicy{}

		// DataWarehouse SKUs do not support LRP currently
		if !isDwSku {
			longTermRetentionPolicyPayload.Properties = longTermRetentionPolicyProps
		}

		if err := longTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, longTermRetentionPolicyPayload); err != nil {
			return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
		}
	}

	shortTermRetentionPolicyProps := helper.ExpandShortTermRetentionPolicy(d.Get("short_term_retention_policy").([]any))
	if shortTermRetentionPolicyProps != nil {
		shortTermRetentionPolicyPayload := backupshorttermretentionpolicies.BackupShortTermRetentionPolicy{}

		if !isDwSku {
			shortTermRetentionPolicyPayload.Properties = shortTermRetentionPolicyProps
		}

		if strings.HasPrefix(skuName, "HS") || strings.HasPrefix(elasticPoolSku, "HS") {
			shortTermRetentionPolicyPayload.Properties.DiffBackupIntervalInHours = nil
		} else if shortTermRetentionPolicyProps.DiffBackupIntervalInHours == nil || pointer.From(shortTermRetentionPolicyProps.DiffBackupIntervalInHours) == 0 {
			shortTermRetentionPolicyPayload.Properties.DiffBackupIntervalInHours = pointer.To(backupshorttermretentionpolicies.DiffBackupIntervalInHoursOneTwo)
		}

		if err := shortTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, shortTermRetentionPolicyPayload); err != nil {
			return fmt.Errorf("setting Short Term Retention Policies for %s: %+v", id, err)
		}
	}

	return resourceMsSqlDatabaseRead(d, meta)
}

func expandMsSqlDatabaseSecurityAlertPolicy(d *pluginsdk.ResourceData) databasesecurityalertpolicies.DatabaseSecurityAlertPolicy {
	policy := databasesecurityalertpolicies.DatabaseSecurityAlertPolicy{
		Properties: pointer.To(databasesecurityalertpolicies.SecurityAlertsPolicyProperties{
			State: databasesecurityalertpolicies.SecurityAlertsPolicyStateDisabled,
		}),
	}
	properties := policy.Properties

	td, ok := d.GetOk("threat_detection_policy")
	if !ok {
		return policy
	}

	if tdl := td.([]any); len(tdl) > 0 {
		securityAlert := tdl[0].(map[string]any)

		properties.State = databasesecurityalertpolicies.SecurityAlertsPolicyState(securityAlert["state"].(string))
		properties.EmailAccountAdmins = pointer.To(securityAlert["email_account_admins_enabled"].(bool))

		if v, ok := securityAlert["disabled_alerts"]; ok {
			alerts := v.(*pluginsdk.Set).List()
			expandedAlerts := make([]string, len(alerts))
			for i, a := range alerts {
				expandedAlerts[i] = a.(string)
			}
			properties.DisabledAlerts = &expandedAlerts
		}
		if v, ok := securityAlert["email_addresses"]; ok {
			emails := v.(*pluginsdk.Set).List()
			expandedEmails := make([]string, len(emails))
			for i, e := range emails {
				expandedEmails[i] = e.(string)
			}
			properties.EmailAddresses = &expandedEmails
		}
		if v, ok := securityAlert["retention_days"]; ok {
			properties.RetentionDays = pointer.To(int64(v.(int)))
		}
		if v, ok := securityAlert["storage_account_access_key"]; ok && v.(string) != "" {
			properties.StorageAccountAccessKey = pointer.To(v.(string))
		}
		if v, ok := securityAlert["storage_endpoint"]; ok && v.(string) != "" {
			properties.StorageEndpoint = pointer.To(v.(string))
		}

		return policy
	}

	return policy
}

func expandMsSqlServerImport(d *pluginsdk.ResourceData) (out databases.ImportExistingDatabaseDefinition) {
	v := d.Get("import")
	dbImportRefs := v.([]any)
	dbImportRef := dbImportRefs[0].(map[string]any)
	out = databases.ImportExistingDatabaseDefinition{
		StorageKeyType:             databases.StorageKeyType(dbImportRef["storage_key_type"].(string)),
		StorageKey:                 dbImportRef["storage_key"].(string),
		StorageUri:                 dbImportRef["storage_uri"].(string),
		AdministratorLogin:         dbImportRef["administrator_login"].(string),
		AdministratorLoginPassword: pointer.To(dbImportRef["administrator_login_password"].(string)),
		AuthenticationType:         pointer.To(dbImportRef["authentication_type"].(string)),
	}

	if storageAccountId, ok := d.GetOk("storage_account_id"); ok {
		out.NetworkIsolation = &databases.NetworkIsolationSettings{
			StorageAccountResourceId: pointer.To(storageAccountId.(string)),
			SqlServerResourceId:      pointer.To(d.Get("server_id").(string)),
		}
	}
	return
}
