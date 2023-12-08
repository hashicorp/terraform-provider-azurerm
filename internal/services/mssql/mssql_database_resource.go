// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/backupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/databasesecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/longtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/serversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-02-01-preview/transparentdataencryptions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlDatabase() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMsSqlDatabaseCreate,
		Read:   resourceMsSqlDatabaseRead,
		Update: resourceMsSqlDatabaseUpdate,
		Delete: resourceMsSqlDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := commonids.ParseSqlDatabaseID(id)
			return err
		}, resourceMsSqlDatabaseImporter),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DatabaseV0ToV1{},
		}),

		Schema: resourceMsSqlDatabaseSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("sku_name", func(ctx context.Context, old, new, _ interface{}) bool {
				// "hyperscale can not change to other sku
				return strings.HasPrefix(old.(string), "HS") && !strings.HasPrefix(new.(string), "HS")
			}),
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
				transparentDataEncryption := d.Get("transparent_data_encryption_enabled").(bool)
				sku := d.Get("sku_name").(string)
				if !strings.HasPrefix(sku, "DW") && !transparentDataEncryption {
					return fmt.Errorf("transparent data encryption can only be disabled on Data Warehouse SKUs")
				}

				return nil
			}),
	}
}

func resourceMsSqlDatabaseImporter(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
	legacyClient := meta.(*clients.Client).MSSQL.LegacyDatabasesClient
	legacyreplicationLinksClient := meta.(*clients.Client).MSSQL.LegacyReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient

	id, err := commonids.ParseSqlDatabaseID(d.Id())
	if err != nil {
		return nil, err
	}

	partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, legacyClient, legacyreplicationLinksClient, resourcesClient, *id, []sql.ReplicationRole{sql.ReplicationRolePrimary})
	if err != nil {
		return nil, err
	}

	if len(partnerDatabases) > 0 {
		partnerDatabase := partnerDatabases[0]

		partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
		}

		d.Set("create_mode", string(databases.CreateModeSecondary))
		d.Set("creation_source_database_id", partnerDatabaseId.ID())

		return []*pluginsdk.ResourceData{d}, nil
	}

	d.Set("create_mode", string(databases.CreateModeDefault))

	return []*pluginsdk.ResourceData{d}, nil
}

func resourceMsSqlDatabaseCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	legacyClient := meta.(*clients.Client).MSSQL.LegacyDatabasesClient
	serversClient := meta.(*clients.Client).MSSQL.ServersClient
	databaseSecurityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	legacyReplicationLinksClient := meta.(*clients.Client).MSSQL.LegacyReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database creation.")

	if strings.HasPrefix(d.Get("sku_name").(string), "GP_S_") && d.Get("license_type").(string) != "" {
		return fmt.Errorf("serverless databases do not support license type")
	}

	name := d.Get("name").(string)

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID: %+v", err)
	}

	id := commonids.NewSqlDatabaseID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, name)

	if existing, err := client.Get(ctx, id, databases.DefaultGetOperationOptions()); err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	} else {
		return tf.ImportAsExistsError("azurerm_mssql_database", id.ID())
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

	if skuName := d.Get("sku_name"); skuName != "" {
		partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, legacyClient, legacyReplicationLinksClient, resourcesClient, id, []sql.ReplicationRole{sql.ReplicationRoleSecondary, sql.ReplicationRoleNonReadableSecondary})
		if err != nil {
			return err
		}

		// Place a lock for the partner databases, so they can't update themselves whilst we're poking their SKUs
		for _, partnerDatabase := range partnerDatabases {
			partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.ID)
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
			}

			locks.ByID(partnerDatabaseId.ID())
			defer locks.UnlockByID(partnerDatabaseId.ID())
		}

		// Update the SKUs of any partner databases where deemed necessary
		for _, partnerDatabase := range partnerDatabases {
			partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.ID)
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
			}

			// See: https://docs.microsoft.com/en-us/azure/azure-sql/database/active-geo-replication-overview#configuring-secondary-database
			if partnerDatabase.Sku != nil && partnerDatabase.Sku.Name != nil && helper.CompareDatabaseSkuServiceTiers(skuName.(string), *partnerDatabase.Sku.Name) {
				future, err := legacyClient.Update(ctx, partnerDatabaseId.ResourceGroupName, partnerDatabaseId.ServerName, partnerDatabaseId.DatabaseName, sql.DatabaseUpdate{
					Sku: &sql.Sku{
						Name: utils.String(skuName.(string)),
					},
				})
				if err != nil {
					return fmt.Errorf("updating SKU of Replication Partner %s: %+v", partnerDatabaseId, err)
				}

				if err = future.WaitForCompletionRef(ctx, legacyClient.Client); err != nil {
					return fmt.Errorf("waiting for SKU update for Replication Partner %s: %+v", partnerDatabaseId, err)
				}
			}
		}
	}

	input := databases.Database{
		Location: location,
		Properties: &databases.DatabaseProperties{
			AutoPauseDelay:                   pointer.To(int64(d.Get("auto_pause_delay_in_minutes").(int))),
			Collation:                        pointer.To(d.Get("collation").(string)),
			ElasticPoolId:                    pointer.To(d.Get("elastic_pool_id").(string)),
			LicenseType:                      pointer.To(databases.DatabaseLicenseType(d.Get("license_type").(string))),
			MinCapacity:                      utils.Float(d.Get("min_capacity").(float64)),
			HighAvailabilityReplicaCount:     pointer.To(int64(d.Get("read_replica_count").(int))),
			SampleName:                       pointer.To(databases.SampleName(d.Get("sample_name").(string))),
			RequestedBackupStorageRedundancy: pointer.To(databases.BackupStorageRedundancy(d.Get("storage_account_type").(string))),
			ZoneRedundant:                    pointer.To(d.Get("zone_redundant").(bool)),
			IsLedgerOn:                       pointer.To(ledgerEnabled),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	createMode, ok := d.GetOk("create_mode")
	if _, dbok := d.GetOk("creation_source_database_id"); ok && (createMode.(string) == string(databases.CreateModeCopy) || createMode.(string) == string(databases.CreateModePointInTimeRestore) || createMode.(string) == string(databases.CreateModeSecondary)) && !dbok {
		return fmt.Errorf("'creation_source_database_id' is required for create_mode %s", createMode.(string))
	}
	if _, dbok := d.GetOk("recover_database_id"); ok && createMode.(string) == string(databases.CreateModeRecovery) && !dbok {
		return fmt.Errorf("'recover_database_id' is required for create_mode %s", createMode.(string))
	}
	if _, dbok := d.GetOk("restore_dropped_database_id"); ok && createMode.(string) == string(databases.CreateModeRestore) && !dbok {
		return fmt.Errorf("'restore_dropped_database_id' is required for create_mode %s", createMode.(string))
	}

	// we should not specify the value of `maintenance_configuration_name` when `elastic_pool_id` is set since its value depends on the elastic pool's `maintenance_configuration_name` value.
	if _, ok := d.GetOk("elastic_pool_id"); !ok {
		// set default value here because `elastic_pool_id` is not specified, API returns default value `SQL_Default` for `maintenance_configuration_name`
		maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, "SQL_Default")
		if v, ok := d.GetOk("maintenance_configuration_name"); ok {
			maintenanceConfigId = publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, v.(string))
		}
		input.Properties.MaintenanceConfigurationId = utils.String(maintenanceConfigId.ID())
	}

	input.Properties.CreateMode = pointer.To(databases.CreateMode(createMode.(string)))

	if v, ok := d.GetOk("max_size_gb"); ok {
		// `max_size_gb` is Computed, so has a value after the first run
		if createMode != string(databases.CreateModeOnlineSecondary) && createMode != string(databases.CreateModeSecondary) {
			input.Properties.MaxSizeBytes = pointer.To(int64(v.(int)) * 1073741824)
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
			return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(databases.CreateModePointInTimeRestore))
		}

		input.Properties.RestorePointInTime = pointer.To(v.(string))
	}

	skuName, ok := d.GetOk("sku_name")
	if ok {
		input.Sku = pointer.To(databases.Sku{
			Name: skuName.(string),
		})
	}

	if v, ok := d.GetOk("creation_source_database_id"); ok {
		input.Properties.SourceDatabaseId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("recover_database_id"); ok {
		input.Properties.RecoverableDatabaseId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("restore_dropped_database_id"); ok {
		input.Properties.RestorableDroppedDatabaseId = pointer.To(v.(string))
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, input)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
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
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Checking to see if %s is online...", id)

			resp, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
			if err != nil {
				return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
			}

			if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.Status != nil {
				return resp, string(pointer.From(resp.Model.Properties.Status)), nil
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

		input := transparentdataencryptions.LogicalDatabaseTransparentDataEncryption{
			Properties: &transparentdataencryptions.TransparentDataEncryptionProperties{
				State: state,
			},
		}

		err := transparentEncryptionClient.CreateOrUpdateThenPoll(ctx, id, input)
		if err != nil {
			return fmt.Errorf("while enabling Transparent Data Encryption for %q: %+v", id.String(), err)
		}

		// NOTE: Internal x-ref, this is another case of hashicorp/go-azure-sdk#307 so this can be removed once that's fixed
		if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
			c, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
			if err != nil {
				return pluginsdk.NonRetryableError(fmt.Errorf("while polling %s for status: %+v", id.String(), err))
			}
			if c.Model != nil && c.Model.Properties != nil && c.Model.Properties.Status != nil {
				if c.Model.Properties.Status == pointer.To(databases.DatabaseStatusScaling) {
					return pluginsdk.RetryableError(fmt.Errorf("database %s is still scaling", id.String()))
				}
			} else {
				return pluginsdk.RetryableError(fmt.Errorf("retrieving database status %s: Model, Properties or Status is nil", id.String()))
			}

			return nil
		}); err != nil {
			return nil
		}
	}

	if _, ok := d.GetOk("import"); ok {
		importParameters := expandMsSqlServerImport(d)
		err := client.ImportThenPoll(ctx, id, importParameters)
		if err != nil {
			return fmt.Errorf("while import bacpac into the new database %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	// For datawarehouse SKUs only
	if strings.HasPrefix(skuName.(string), "DW") {
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
		return nil
	}

	securityAlertPolicyProps := helper.ExpandLongTermRetentionPolicy(d.Get("long_term_retention_policy").([]interface{}))
	if securityAlertPolicyProps != nil {
		securityAlertPolicyPayload := longtermretentionpolicies.LongTermRetentionPolicy{}

		// DataWarehouse SKU's do not support LRP currently
		if !strings.HasPrefix(skuName.(string), "DW") {
			securityAlertPolicyPayload.Properties = securityAlertPolicyProps
		}

		err := longTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, securityAlertPolicyPayload)
		if err != nil {
			return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
		}
	}

	shortTermSecurityAlertPolicyProps := helper.ExpandShortTermRetentionPolicy(d.Get("short_term_retention_policy").([]interface{}))
	if securityAlertPolicyProps != nil {
		securityAlertPolicyPayload := backupshorttermretentionpolicies.BackupShortTermRetentionPolicy{}

		if !strings.HasPrefix(skuName.(string), "DW") {
			securityAlertPolicyPayload.Properties = shortTermSecurityAlertPolicyProps
		}

		if strings.HasPrefix(skuName.(string), "HS") {
			securityAlertPolicyPayload.Properties.DiffBackupIntervalInHours = nil
		}

		err := shortTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, securityAlertPolicyPayload)
		if err != nil {
			return fmt.Errorf("setting Short Term Retention Policies for %s: %+v", id, err)
		}
	}

	return resourceMsSqlDatabaseRead(d, meta)
}

func resourceMsSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	securityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient

	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	serverId := commonids.NewSqlServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
	d.Set("server_id", serverId.ID())

	resp, err := client.Get(ctx, pointer.From(id), databases.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	geoBackupPolicy := true
	skuName := ""
	elasticPoolId := ""
	minCapacity := float64(0)
	ledgerEnabled := false

	if model := resp.Model; model != nil {
		d.Set("name", id.DatabaseName)

		if props := model.Properties; props != nil {
			minCapacity = pointer.From(props.MinCapacity)

			requestedBackupStorageRedundancy := ""
			if props.RequestedBackupStorageRedundancy != nil {
				requestedBackupStorageRedundancy = string(*props.RequestedBackupStorageRedundancy)
			}

			d.Set("auto_pause_delay_in_minutes", pointer.From(props.AutoPauseDelay))
			d.Set("collation", pointer.From(props.Collation))
			d.Set("read_replica_count", pointer.From(props.HighAvailabilityReplicaCount))
			d.Set("storage_account_type", requestedBackupStorageRedundancy)
			d.Set("zone_redundant", pointer.From(props.ZoneRedundant))
			d.Set("read_scale", pointer.From(props.ReadScale) == databases.DatabaseReadScaleEnabled)

			if props.ElasticPoolId != nil {
				elasticPoolId = pointer.From(props.ElasticPoolId)
			}

			if props.LicenseType != nil {
				d.Set("license_type", string(pointer.From(props.LicenseType)))
			} else {
				// value not returned, try to set from state
				d.Set("license_type", d.Get("license_type").(string))
			}

			if props.MaxSizeBytes != nil {
				d.Set("max_size_gb", int32((*props.MaxSizeBytes)/int64(1073741824)))
			}

			if props.CurrentServiceObjectiveName != nil {
				skuName = *props.CurrentServiceObjectiveName
			}

			if props.IsLedgerOn != nil {
				ledgerEnabled = *props.IsLedgerOn
			}

			configurationName := ""
			if v := props.MaintenanceConfigurationId; v != nil {
				maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(pointer.From(v))
				if err != nil {
					return err
				}
				configurationName = maintenanceConfigId.PublicMaintenanceConfigurationName
			}

			d.Set("elastic_pool_id", elasticPoolId)
			d.Set("min_capacity", minCapacity)
			d.Set("sku_name", skuName)
			d.Set("maintenance_configuration_name", configurationName)
			d.Set("ledger_enabled", ledgerEnabled)

			if err := tags.FlattenAndSet(d, resp.Model.Tags); err != nil {
				return err
			}
		}

		// DW SKU's do not currently support LRP and do not honour normal SRP operations
		if !strings.HasPrefix(skuName, "DW") {
			longTermPolicy, err := longTermRetentionClient.Get(ctx, pointer.From(id))
			if err != nil {
				return fmt.Errorf("retrieving Long Term Retention Policies for %s: %+v", id, err)
			}

			if model := longTermPolicy.Model; model != nil {
				if err := d.Set("long_term_retention_policy", helper.FlattenLongTermRetentionPolicy(model)); err != nil {
					return fmt.Errorf("setting `long_term_retention_policy`: %+v", err)
				}
			}

			shortTermPolicy, err := shortTermRetentionClient.Get(ctx, pointer.From(id))

			if model := shortTermPolicy.Model; model != nil {
				if err != nil {
					return fmt.Errorf("retrieving Short Term Retention Policies for %s: %+v", id, err)
				}

				if err := d.Set("short_term_retention_policy", helper.FlattenShortTermRetentionPolicy(model)); err != nil {
					return fmt.Errorf("setting `short_term_retention_policy`: %+v", err)
				}
			}
		} else {
			// DW SKUs need the retention policies to be empty for state consistency
			emptySlice := make([]interface{}, 0)
			d.Set("long_term_retention_policy", emptySlice)
			d.Set("short_term_retention_policy", emptySlice)

			geoPoliciesResponse, err := geoBackupPoliciesClient.Get(ctx, pointer.From(id))
			if err != nil {
				if response.WasNotFound(geoPoliciesResponse.HttpResponse) {
					d.SetId("")
					return nil
				}

				return fmt.Errorf("retrieving Geo Backup Policies for %s: %+v", id, err)
			}

			// For Datawarehouse SKUs, set the geo-backup policy setting
			if model := geoPoliciesResponse.Model; model != nil {
				if strings.HasPrefix(skuName, "DW") && model.Properties.State == geobackuppolicies.GeoBackupPolicyStateDisabled {
					geoBackupPolicy = false
				}
			}
		}
	}

	if err := d.Set("geo_backup_enabled", geoBackupPolicy); err != nil {
		return fmt.Errorf("setting `geo_backup_enabled`: %+v", err)
	}

	securityAlertPolicy, err := securityAlertPoliciesClient.Get(ctx, pointer.From(id))
	if err == nil && securityAlertPolicy.Model != nil {
		if err := d.Set("threat_detection_policy", flattenMsSqlServerSecurityAlertPolicy(d, pointer.From(securityAlertPolicy.Model))); err != nil {
			return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
		}
	}

	tde, err := transparentEncryptionClient.Get(ctx, pointer.From(id))
	if err != nil {
		return fmt.Errorf("while retrieving Transparent Data Encryption state for %s: %+v", id, err)
	}

	tdeState := false
	if model := tde.Model; model != nil {
		if props := model.Properties; props != nil {
			tdeState = (props.State == transparentdataencryptions.TransparentDataEncryptionStateEnabled)
		}
	}
	d.Set("transparent_data_encryption_enabled", tdeState)

	return nil
}

func resourceMsSqlDatabaseUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	legacyClient := meta.(*clients.Client).MSSQL.LegacyDatabasesClient
	serversClient := meta.(*clients.Client).MSSQL.ServersClient
	securityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	legacyReplicationLinksClient := meta.(*clients.Client).MSSQL.LegacyReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database creation.")

	skuName := d.Get("sku_name").(string)
	if strings.HasPrefix(skuName, "GP_S_") && d.Get("license_type").(string) != "" {
		return fmt.Errorf("serverless databases do not support license type")
	}

	name := d.Get("name").(string)

	serverId, err := commonids.ParseSqlServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID: %+v", err)
	}

	id := commonids.NewSqlDatabaseID(serverId.SubscriptionId, serverId.ResourceGroupName, serverId.ServerName, name)

	_, err = client.Get(ctx, id, databases.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+q", id, err)
	}

	_, err = serversClient.Get(ctx, pointer.From(serverId), servers.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %q", serverId, err)
	}

	// when disassociating mssql db from elastic pool, the sku_name must be specific
	if d.HasChange("elastic_pool_id") {
		if old, new := d.GetChange("elastic_pool_id"); old.(string) != "" && new.(string) == "" {
			if v, ok := d.GetOk("sku_name"); !ok || (ok && v.(string) == "ElasticPool") {
				return fmt.Errorf("`sku_name` must be assigned and not be %q when disassociating from Elastic Pool", "ElasticPool")
			}
		}
	}

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

	if d.HasChange("sku_name") && skuName != "" {
		partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, legacyClient, legacyReplicationLinksClient, resourcesClient, id, []sql.ReplicationRole{sql.ReplicationRoleSecondary, sql.ReplicationRoleNonReadableSecondary})
		if err != nil {
			return err
		}

		// Place a lock for the partner databases, so they can't update themselves whilst we're poking their SKUs
		for _, v := range partnerDatabases {
			id, err := commonids.ParseSqlDatabaseIDInsensitively(pointer.From(v.ID))
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", id.ID(), err)
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())
		}

		// Update the SKUs of any partner databases where deemed necessary
		for _, partnerDatabase := range partnerDatabases {
			partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.ID)
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
			}

			// See: https://docs.microsoft.com/en-us/azure/azure-sql/database/active-geo-replication-overview#configuring-secondary-database
			if partnerDatabase.Sku != nil && partnerDatabase.Sku.Name != nil && helper.CompareDatabaseSkuServiceTiers(skuName, *partnerDatabase.Sku.Name) {
				future, err := legacyClient.Update(ctx, partnerDatabaseId.ResourceGroupName, partnerDatabaseId.ServerName, partnerDatabaseId.DatabaseName, sql.DatabaseUpdate{
					Sku: &sql.Sku{
						Name: pointer.To(skuName),
					},
				})
				if err != nil {
					return fmt.Errorf("updating SKU of Replication Partner %s: %+v", partnerDatabaseId, err)
				}

				if err = future.WaitForCompletionRef(ctx, legacyClient.Client); err != nil {
					return fmt.Errorf("waiting for SKU update for Replication Partner %s: %+v", partnerDatabaseId, err)
				}
			}
		}
	}

	payload := databases.DatabaseUpdate{}
	props := databases.DatabaseUpdateProperties{}

	if d.HasChange("auto_pause_delay_in_minutes") {
		props.AutoPauseDelay = pointer.To(int64(d.Get("auto_pause_delay_in_minutes").(int)))
	}

	if d.HasChange("elastic_pool_id") {
		props.ElasticPoolId = pointer.To(d.Get("elastic_pool_id").(string))
	}

	if d.HasChange("license_type") {
		props.LicenseType = pointer.To(databases.DatabaseLicenseType(d.Get("license_type").(string)))
	}

	if d.HasChange("min_capacity") {
		props.MinCapacity = pointer.To(d.Get("min_capacity").(float64))
	}

	if d.HasChange("read_replica_count") {
		props.HighAvailabilityReplicaCount = pointer.To(int64(d.Get("read_replica_count").(int)))
	}

	if d.HasChange("sample_name") {
		props.SampleName = pointer.To(databases.SampleName(d.Get("sample_name").(string)))
	}

	if d.HasChange("storage_account_type") {
		props.RequestedBackupStorageRedundancy = pointer.To(databases.BackupStorageRedundancy(d.Get("storage_account_type").(string)))
	}

	if d.HasChange("zone_redundant") {
		props.ZoneRedundant = pointer.To(d.Get("zone_redundant").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	// we should not specify the value of `maintenance_configuration_name` when `elastic_pool_id` is set since its value depends on the elastic pool's `maintenance_configuration_name` value.
	if _, ok := d.GetOk("elastic_pool_id"); !ok && d.HasChange("maintenance_configuration_name") {
		// set default value here because `elastic_pool_id` is not specified, API returns default value `SQL_Default` for `maintenance_configuration_name`
		maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, "SQL_Default")
		if v, ok := d.GetOk("maintenance_configuration_name"); ok {
			maintenanceConfigId = publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(serverId.SubscriptionId, v.(string))
		}
		props.MaintenanceConfigurationId = pointer.To(maintenanceConfigId.ID())
	}

	createMode := d.Get("create_mode").(string)
	if v, ok := d.GetOk("max_size_gb"); ok {
		// `max_size_gb` is Computed, so has a value after the first run
		if createMode != string(databases.CreateModeOnlineSecondary) && createMode != string(databases.CreateModeSecondary) {
			props.MaxSizeBytes = utils.Int64(int64(v.(int)) * 1073741824)
		}
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
		if v, ok := d.GetOk("restore_point_in_time"); ok {
			if createMode != string(databases.CreateModePointInTimeRestore) {
				return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(databases.CreateModePointInTimeRestore))
			}
			props.RestorePointInTime = pointer.To(v.(string))
		}
	}

	if d.HasChange("sku_name") {
		payload.Sku = pointer.To(databases.Sku{
			Name: skuName,
		})
	}

	if d.HasChange("recover_database_id") {
		props.RecoverableDatabaseId = pointer.To(d.Get("recover_database_id").(string))
	}

	if d.HasChange("restore_dropped_database_id") {
		props.RestorableDroppedDatabaseId = pointer.To(d.Get("restore_dropped_database_id").(string))
	}

	payload.Properties = pointer.To(props)

	err = client.UpdateThenPoll(ctx, id, payload)
	if err != nil {
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
		Refresh: func() (interface{}, string, error) {
			log.Printf("[DEBUG] Checking to see if %s is online...", id)

			resp, err := client.Get(ctx, id, databases.DefaultGetOperationOptions())
			if err != nil {
				return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					return resp, string(pointer.From(props.Status)), nil
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

			err := transparentEncryptionClient.CreateOrUpdateThenPoll(ctx, id, input)
			if err != nil {
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
				return nil
			}
		}

	}

	if d.HasChange("import") {
		if _, ok := d.GetOk("import"); ok {
			importParameters := expandMsSqlServerImport(d)

			err := client.ImportThenPoll(ctx, id, importParameters)
			if err != nil {
				return fmt.Errorf("while importing the BACPAC file into the new database %s: %+v", id.ID(), err)
			}
		}
	}

	d.SetId(id.ID())

	// For datawarehouse SKUs only
	if strings.HasPrefix(skuName, "DW") && d.HasChange("geo_backup_enabled") {
		isEnabled := d.Get("geo_backup_enabled").(bool)
		var geoBackupPolicyState geobackuppolicies.GeoBackupPolicyState

		if isEnabled {
			geoBackupPolicyState = geobackuppolicies.GeoBackupPolicyStateEnabled
		} else {
			geoBackupPolicyState = geobackuppolicies.GeoBackupPolicyStateDisabled
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
		return nil
	}

	if d.HasChange("long_term_retention_policy") {
		v := d.Get("long_term_retention_policy")
		longTermRetentionProps := helper.ExpandLongTermRetentionPolicy(v.([]interface{}))
		if longTermRetentionProps != nil {
			longTermRetentionPolicy := longtermretentionpolicies.LongTermRetentionPolicy{}

			// DataWarehouse SKU's do not support LRP currently
			if !strings.HasPrefix(skuName, "DW") {
				longTermRetentionPolicy.Properties = longTermRetentionProps
			}

			err := longTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, longTermRetentionPolicy)
			if err != nil {
				return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("short_term_retention_policy") {
		v := d.Get("short_term_retention_policy")
		backupShortTermPolicyProps := helper.ExpandShortTermRetentionPolicy(v.([]interface{}))
		if backupShortTermPolicyProps != nil {
			backupShortTermPolicy := backupshorttermretentionpolicies.BackupShortTermRetentionPolicy{}

			if !strings.HasPrefix(skuName, "DW") {
				backupShortTermPolicy.Properties = backupShortTermPolicyProps
			}

			if strings.HasPrefix(skuName, "HS") {
				backupShortTermPolicy.Properties.DiffBackupIntervalInHours = nil
			}

			err := shortTermRetentionClient.CreateOrUpdateThenPoll(ctx, id, backupShortTermPolicy)
			if err != nil {
				return fmt.Errorf("setting Short Term Retention Policies for %s: %+v", id, err)
			}
		}
	}

	return resourceMsSqlDatabaseRead(d, meta)
}

func resourceMsSqlDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func flattenMsSqlServerSecurityAlertPolicy(d *pluginsdk.ResourceData, policy databasesecurityalertpolicies.DatabaseSecurityAlertPolicy) []interface{} {
	// The SQL database security alert API always returns the default value even if never set.
	// If the values are on their default one, threat it as not set.
	properties := policy.Properties
	if properties == nil {
		return []interface{}{}
	}

	securityAlertPolicy := make(map[string]interface{})

	securityAlertPolicy["state"] = string(properties.State)

	securityAlertPolicy["email_account_admins"] = "Disabled"
	if properties.EmailAccountAdmins != nil && *properties.EmailAccountAdmins {
		securityAlertPolicy["email_account_admins"] = "Enabled"
	}

	if disabledAlerts := properties.DisabledAlerts; disabledAlerts != nil {
		flattenedAlerts := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
		for _, a := range *disabledAlerts {
			if a != "" {
				flattenedAlerts.Add(a)
			}
		}
		securityAlertPolicy["disabled_alerts"] = flattenedAlerts
	}
	if emailAddresses := properties.EmailAddresses; emailAddresses != nil {
		flattenedEmails := pluginsdk.NewSet(pluginsdk.HashString, []interface{}{})
		for _, e := range *emailAddresses {
			if e != "" {
				flattenedEmails.Add(e)
			}
		}
		securityAlertPolicy["email_addresses"] = flattenedEmails
	}
	if properties.StorageEndpoint != nil {
		securityAlertPolicy["storage_endpoint"] = *properties.StorageEndpoint
	}
	if properties.RetentionDays != nil {
		securityAlertPolicy["retention_days"] = int(*properties.RetentionDays)
	}

	// If storage account access key is in state read it to the new state, as the API does not return it for security reasons
	if v, ok := d.GetOk("threat_detection_policy.0.storage_account_access_key"); ok {
		securityAlertPolicy["storage_account_access_key"] = v.(string)
	}

	return []interface{}{securityAlertPolicy}
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

	if tdl := td.([]interface{}); len(tdl) > 0 {
		securityAlert := tdl[0].(map[string]interface{})

		properties.State = databasesecurityalertpolicies.SecurityAlertsPolicyState(securityAlert["state"].(string))
		properties.EmailAccountAdmins = utils.Bool(securityAlert["email_account_admins"].(string) == string(EmailAccountAdminsStatusEnabled))

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
			properties.StorageAccountAccessKey = utils.String(v.(string))
		}
		if v, ok := securityAlert["storage_endpoint"]; ok && v.(string) != "" {
			properties.StorageEndpoint = utils.String(v.(string))
		}

		return policy
	}

	return policy
}

func expandMsSqlServerImport(d *pluginsdk.ResourceData) (out databases.ImportExistingDatabaseDefinition) {
	v := d.Get("import")
	dbImportRefs := v.([]interface{})
	dbImportRef := dbImportRefs[0].(map[string]interface{})
	out = databases.ImportExistingDatabaseDefinition{
		StorageKeyType:             databases.StorageKeyType(dbImportRef["storage_key_type"].(string)),
		StorageKey:                 dbImportRef["storage_key"].(string),
		StorageUri:                 dbImportRef["storage_uri"].(string),
		AdministratorLogin:         dbImportRef["administrator_login"].(string),
		AdministratorLoginPassword: dbImportRef["administrator_login_password"].(string),
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

func resourceMsSqlDatabaseMaintenanceNames() []string {
	return []string{"SQL_Default", "SQL_EastUS_DB_1", "SQL_EastUS2_DB_1", "SQL_SoutheastAsia_DB_1", "SQL_AustraliaEast_DB_1", "SQL_NorthEurope_DB_1", "SQL_SouthCentralUS_DB_1", "SQL_WestUS2_DB_1",
		"SQL_UKSouth_DB_1", "SQL_WestEurope_DB_1", "SQL_EastUS_DB_2", "SQL_EastUS2_DB_2", "SQL_WestUS2_DB_2", "SQL_SoutheastAsia_DB_2", "SQL_AustraliaEast_DB_2", "SQL_NorthEurope_DB_2", "SQL_SouthCentralUS_DB_2",
		"SQL_UKSouth_DB_2", "SQL_WestEurope_DB_2", "SQL_AustraliaSoutheast_DB_1", "SQL_BrazilSouth_DB_1", "SQL_CanadaCentral_DB_1", "SQL_CanadaEast_DB_1", "SQL_CentralUS_DB_1", "SQL_EastAsia_DB_1",
		"SQL_FranceCentral_DB_1", "SQL_GermanyWestCentral_DB_1", "SQL_CentralIndia_DB_1", "SQL_SouthIndia_DB_1", "SQL_JapanEast_DB_1", "SQL_JapanWest_DB_1", "SQL_NorthCentralUS_DB_1", "SQL_UKWest_DB_1",
		"SQL_WestUS_DB_1", "SQL_AustraliaSoutheast_DB_2", "SQL_BrazilSouth_DB_2", "SQL_CanadaCentral_DB_2", "SQL_CanadaEast_DB_2", "SQL_CentralUS_DB_2", "SQL_EastAsia_DB_2", "SQL_FranceCentral_DB_2",
		"SQL_GermanyWestCentral_DB_2", "SQL_CentralIndia_DB_2", "SQL_SouthIndia_DB_2", "SQL_JapanEast_DB_2", "SQL_JapanWest_DB_2", "SQL_NorthCentralUS_DB_2", "SQL_UKWest_DB_2", "SQL_WestUS_DB_2",
		"SQL_WestCentralUS_DB_1", "SQL_FranceSouth_DB_1", "SQL_WestCentralUS_DB_2", "SQL_FranceSouth_DB_2", "SQL_SwitzerlandNorth_DB_1", "SQL_SwitzerlandNorth_DB_2", "SQL_BrazilSoutheast_DB_1",
		"SQL_UAENorth_DB_1", "SQL_BrazilSoutheast_DB_2", "SQL_UAENorth_DB_2"}
}

type EmailAccountAdminsStatus string

const (
	EmailAccountAdminsStatusDisabled EmailAccountAdminsStatus = "Disabled"
	EmailAccountAdminsStatusEnabled  EmailAccountAdminsStatus = "Enabled"
)

func PossibleValuesForEmailAccountAdminsStatus() []string {
	return []string{
		string(EmailAccountAdminsStatusDisabled),
		string(EmailAccountAdminsStatusEnabled),
	}
}

func resourceMsSqlDatabaseSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ValidateMsSqlDatabaseName,
		},

		"server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ServerID,
		},

		"auto_pause_delay_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.DatabaseAutoPauseDelay,
		},

		"create_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(databases.CreateModeDefault),
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForCreateMode(),
				false),
			ConflictsWith: []string{"import"},
		},
		"import": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"storage_uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"storage_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"storage_key_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(databases.PossibleValuesForStorageKeyType(),
							false),
					},
					"administrator_login": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"administrator_login_password": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"authentication_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"ADPassword",
							"Sql",
						}, false),
					},
					"storage_account_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateStorageAccountID,
					},
				},
			},
			ConflictsWith: []string{"create_mode"}, // it needs `create_mode` to be `Default` to work, so make them conflict.
		},

		"collation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validate.DatabaseCollation(),
		},

		"elastic_pool_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ElasticPoolID,
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForDatabaseLicenseType(),
				false),
		},

		"long_term_retention_policy": helper.LongTermRetentionPolicySchema(),

		"short_term_retention_policy": helper.ShortTermRetentionPolicySchema(),

		"max_size_gb": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 4096),
		},

		"min_capacity": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.FloatInSlice([]float64{0, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40}),
		},

		"restore_point_in_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"recover_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RecoverableDatabaseID,
		},

		"restore_dropped_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RestorableDatabaseID,
		},

		"read_replica_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"read_scale": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"sample_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(databases.SampleNameAdventureWorksLT),
			}, false),
		},

		"sku_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validate.DatabaseSkuName(),
		},

		"creation_source_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: commonids.ValidateSqlDatabaseID,
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(databases.BackupStorageRedundancyGeo),
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForBackupStorageRedundancy(),
				false),
		},

		"zone_redundant": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"threat_detection_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disabled_alerts": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      pluginsdk.HashString,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Sql_Injection",
								"Sql_Injection_Vulnerability",
								"Access_Anomaly",
							}, false),
						},
					},

					// NOTE: this is a Boolean in SDK rather than a String
					// TODO: update this to be `email_account_admins_enabled` in 4.0
					"email_account_admins": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  EmailAccountAdminsStatusDisabled,
						ValidateFunc: validation.StringInSlice(PossibleValuesForEmailAccountAdminsStatus(),
							false),
					},

					"email_addresses": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"retention_days": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					// NOTE: I believe that this was originally implemented incorrect, where it exposed the
					// 'serveradvancedthreatprotectionsettings.PossibleValuesForAdvancedThreatProtectionState'
					// which contains the values of 'Enabled', 'Disabled', and 'New'
					// where 'serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState'
					// only contains 'Enabled' and 'Disabled'
					"state": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(serversecurityalertpolicies.SecurityAlertsPolicyStateDisabled),
						ValidateFunc: validation.StringInSlice(serversecurityalertpolicies.PossibleValuesForSecurityAlertsPolicyState(),
							false),
					},

					"storage_account_access_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"storage_endpoint": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"geo_backup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"maintenance_configuration_name": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ConflictsWith: []string{"elastic_pool_id"},
			ValidateFunc:  validation.StringInSlice(resourceMsSqlDatabaseMaintenanceNames(), false),
		},

		"ledger_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"transparent_data_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}
