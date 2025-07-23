// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mssql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/backupshorttermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/databasesecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/elasticpools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/geobackuppolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/longtermretentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/replicationlinks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/serversecurityalertpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/transparentdataencryptions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	keyVaultParser "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
				// hyperscale can not be changed to another sku
				return strings.HasPrefix(old.(string), "HS") && !strings.HasPrefix(new.(string), "HS")
			}),
			pluginsdk.ForceNewIfChange("enclave_type", func(ctx context.Context, old, new, _ interface{}) bool {
				// enclave_type cannot be removed once it has been set
				// but can be changed between VBS and Default...
				// this Diff will not work until 4.0 when we remove
				// the computed property from the field scheam.
				if old.(string) != "" && new.(string) == "" {
					return true
				}

				return false
			}),
			func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
				transparentDataEncryption := d.Get("transparent_data_encryption_enabled").(bool)
				skuName := d.Get("sku_name").(string)

				if !strings.HasPrefix(strings.ToLower(skuName), "dw") && !transparentDataEncryption {
					return fmt.Errorf("transparent data encryption can only be disabled on Data Warehouse SKUs")
				}

				// NOTE: VBS enclaves are not supported by DW or DC skus...
				if d.Get("enclave_type").(string) == string(databases.AlwaysEncryptedEnclaveTypeVBS) {
					if strings.HasPrefix(strings.ToLower(skuName), "dw") || strings.Contains(strings.ToLower(skuName), "_dc_") {
						return fmt.Errorf("virtualization based security (VBS) enclaves are not supported for the %q sku", skuName)
					}
				}

				if strings.HasPrefix(strings.ToLower(skuName), "dw") {
					// NOTE: Got `PerDatabaseCMKDWNotSupported` error from API when `sku_name` is set to `DW100c` and `transparent_data_encryption_key_vault_key_id` is specified
					keyVaultKeyId := d.Get("transparent_data_encryption_key_vault_key_id").(string)
					if keyVaultKeyId != "" {
						return fmt.Errorf("database-level CMK is not supported for Data Warehouse SKUs")
					}
					// NOTE: Got `InternalServerError` error from API when `sku_name` is set to `DW100c` and `transparent_data_encryption_key_automatic_rotation_enabled` is specified
					if d.Get("transparent_data_encryption_key_automatic_rotation_enabled").(bool) {
						return fmt.Errorf("transparent_data_encryption_key_automatic_rotation_enabled should not be specified when using Data Warehouse SKUs")
					}
				}

				return nil
			}),
	}
}

func resourceMsSqlDatabaseImporter(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
	replicationLinksClient := meta.(*clients.Client).MSSQL.ReplicationLinksClient
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient

	id, err := commonids.ParseSqlDatabaseID(d.Id())
	if err != nil {
		return nil, err
	}

	// NOTE: The service default is actually nil/empty which indicates enclave is disabled. the value `Default` is NOT the default.
	var enclaveType databases.AlwaysEncryptedEnclaveType
	if v, ok := d.GetOk("enclave_type"); ok && v.(string) != "" {
		enclaveType = databases.AlwaysEncryptedEnclaveType(v.(string))
	}
	d.Set("enclave_type", enclaveType)

	partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, client, replicationLinksClient, resourcesClient, *id, enclaveType, []replicationlinks.ReplicationRole{replicationlinks.ReplicationRolePrimary})
	if err != nil {
		return nil, err
	}

	if len(partnerDatabases) > 0 {
		partnerDatabase := partnerDatabases[0]

		partnerDatabaseId, err := commonids.ParseSqlDatabaseIDInsensitively(*partnerDatabase.Id)
		if err != nil {
			return nil, fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.Id, err)
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

	log.Printf("[INFO] preparing arguments for MsSql Database creation")

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

			// See: https://docs.microsoft.com/en-us/azure/azure-sql/database/active-geo-replication-overview#configuring-secondary-database
			if partnerDatabase.Sku != nil && partnerDatabase.Sku.Name != "" && helper.CompareDatabaseSkuServiceTiers(skuName, partnerDatabase.Sku.Name) {
				err := client.UpdateThenPoll(ctx, *partnerDatabaseId, databases.DatabaseUpdate{
					Sku: &databases.Sku{
						Name: skuName,
					},
				})
				if err != nil {
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
				elasticEnclaveType := string(pointer.From(elasticPool.Model.Properties.PreferredEnclaveType))
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
			LicenseType:                      pointer.To(databases.DatabaseLicenseType(d.Get("license_type").(string))),
			MinCapacity:                      pointer.To(d.Get("min_capacity").(float64)),
			HighAvailabilityReplicaCount:     pointer.To(int64(d.Get("read_replica_count").(int))),
			SampleName:                       pointer.To(databases.SampleName(d.Get("sample_name").(string))),
			RequestedBackupStorageRedundancy: pointer.To(databases.BackupStorageRedundancy(d.Get("storage_account_type").(string))),
			ZoneRedundant:                    pointer.To(d.Get("zone_redundant").(bool)),
			IsLedgerOn:                       pointer.To(ledgerEnabled),
			SecondaryType:                    pointer.To(databases.SecondaryType(d.Get("secondary_type").(string))),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

	input.Properties.CreateMode = pointer.To(databases.CreateMode(createMode))

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
		expandedIdentity, err := identity.ExpandUserAssignedMap(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		input.Identity = expandedIdentity
	}

	if v, ok := d.GetOk("transparent_data_encryption_key_vault_key_id"); ok {
		keyVaultKeyId := v.(string)

		keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)
		if err != nil {
			return fmt.Errorf("unable to parse key: %q: %+v", keyVaultKeyId, err)
		}

		input.Properties.EncryptionProtector = pointer.To(keyId.ID())
	}

	if err = client.CreateOrUpdateThenPoll(ctx, id, input); err != nil {
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

	d.SetId(id.ID())

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

	longTermRetentionPolicyProps := helper.ExpandLongTermRetentionPolicy(d.Get("long_term_retention_policy").([]interface{}))
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

	shortTermRetentionPolicyProps := helper.ExpandShortTermRetentionPolicy(d.Get("short_term_retention_policy").([]interface{}))
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

func resourceMsSqlDatabaseUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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

	log.Printf("[INFO] preparing arguments for MsSql Database update")

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

	_, err = serversClient.Get(ctx, pointer.From(serverId), servers.DefaultGetOperationOptions())
	if err != nil {
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
		if createMode != string(databases.CreateModeOnlineSecondary) && createMode != string(databases.CreateModeSecondary) {
			props.MaxSizeBytes = pointer.To(int64(v.(int)) * 1073741824)
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

				// See: https://docs.microsoft.com/en-us/azure/azure-sql/database/active-geo-replication-overview#configuring-secondary-database
				if partnerDatabase.Sku != nil && partnerDatabase.Sku.Name != "" && helper.CompareDatabaseSkuServiceTiers(skuName, partnerDatabase.Sku.Name) {
					log.Printf("[INFO] Updating SKU of Replication Partner Database from %q to %q", partnerDatabase.Sku.Name, skuName)
					err := client.UpdateThenPoll(ctx, *partnerDatabaseId, databases.DatabaseUpdate{
						Sku: &databases.Sku{
							Name: skuName,
						},
					})
					if err != nil {
						return fmt.Errorf("updating SKU of Replication Partner Database %s: %+v", partnerDatabaseId, err)
					}

					log.Printf("[INFO] SKU of Replication Partner Database updated successfully to %q", skuName)
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
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("identity") {
		expanded, err := identity.ExpandUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		payload.Identity = expanded
	}

	if d.HasChange("transparent_data_encryption_key_vault_key_id") {
		keyVaultKeyId := d.Get("transparent_data_encryption_key_vault_key_id").(string)

		keyId, err := keyVaultParser.ParseNestedItemID(keyVaultKeyId)
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
		longTermRetentionProps := helper.ExpandLongTermRetentionPolicy(v.([]interface{}))
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
		backupShortTermPolicyProps := helper.ExpandShortTermRetentionPolicy(v.([]interface{}))
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
	ledgerEnabled := false
	enclaveType := ""

	if model := resp.Model; model != nil {
		d.Set("name", id.DatabaseName)

		if props := model.Properties; props != nil {
			minCapacity := pointer.From(props.MinCapacity)

			requestedBackupStorageRedundancy := ""
			if props.RequestedBackupStorageRedundancy != nil {
				requestedBackupStorageRedundancy = string(*props.RequestedBackupStorageRedundancy)
			}

			// A named replica doesn't return props.RequestedBackupStorageRedundancy from the api but it is Geo in the portal regardless of what the parent database has
			// so we'll copy that here to get around a perpetual diff
			if props.SecondaryType != nil && *props.SecondaryType == "Named" {
				requestedBackupStorageRedundancy = string(databases.BackupStorageRedundancyGeo)
			}

			d.Set("auto_pause_delay_in_minutes", pointer.From(props.AutoPauseDelay))
			d.Set("collation", pointer.From(props.Collation))
			d.Set("read_replica_count", pointer.From(props.HighAvailabilityReplicaCount))
			d.Set("storage_account_type", requestedBackupStorageRedundancy)
			d.Set("zone_redundant", pointer.From(props.ZoneRedundant))
			d.Set("read_scale", pointer.From(props.ReadScale) == databases.DatabaseReadScaleEnabled)
			d.Set("secondary_type", pointer.From(props.SecondaryType))

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

			// NOTE: Always set the PreferredEnclaveType to an empty string
			// if not in the properties that were returned from Azure...
			if v := props.PreferredEnclaveType; v != nil {
				enclaveType = string(pointer.From(v))
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
			d.Set("enclave_type", enclaveType)
			d.Set("transparent_data_encryption_key_vault_key_id", props.EncryptionProtector)
			d.Set("transparent_data_encryption_key_automatic_rotation_enabled", pointer.From(props.EncryptionProtectorAutoRotation))

			identity, err := identity.FlattenUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
				return err
			}
		}

		// Determine whether the SKU is for SQL Data Warehouse
		isDwSku := strings.HasPrefix(strings.ToLower(skuName), "dw")

		// Determine whether the SKU is for SQL Database Free tier
		isFreeSku := strings.EqualFold(skuName, "free")

		// DW SKUs and SQL Database Free tier do not currently support LRP and do not honour normal SRP operations
		if !isDwSku && !isFreeSku {
			longTermPolicy, err := longTermRetentionClient.Get(ctx, pointer.From(id))
			if err != nil {
				return fmt.Errorf("retrieving Long Term Retention Policies for %s: %+v", id, err)
			}

			if longTermPolicyModel := longTermPolicy.Model; longTermPolicyModel != nil {
				if err := d.Set("long_term_retention_policy", helper.FlattenLongTermRetentionPolicy(longTermPolicyModel)); err != nil {
					return fmt.Errorf("setting `long_term_retention_policy`: %+v", err)
				}
			}

			shortTermPolicy, err := shortTermRetentionClient.Get(ctx, pointer.From(id))
			if err != nil {
				return fmt.Errorf("retrieving Short Term Retention Policies for %s: %+v", id, err)
			}

			if shortTermPolicyModel := shortTermPolicy.Model; shortTermPolicyModel != nil {
				if err := d.Set("short_term_retention_policy", helper.FlattenShortTermRetentionPolicy(shortTermPolicyModel)); err != nil {
					return fmt.Errorf("setting `short_term_retention_policy`: %+v", err)
				}
			}
		} else {
			// DW SKUs and SQL Database Free tier need the retention policies to be empty for state consistency
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

			// For Datawarehouse SKUs and SQL Database Free tier, set the geo-backup policy setting
			if geoPolicyModel := geoPoliciesResponse.Model; geoPolicyModel != nil {
				if (isDwSku || isFreeSku) && geoPolicyModel.Properties.State == geobackuppolicies.GeoBackupPolicyStateDisabled {
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
		properties.EmailAccountAdmins = pointer.To(securityAlert["email_account_admins"].(string) == string(EmailAccountAdminsStatusEnabled))

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

// The following data comes from the results of "az maintenance public-configuration list --query "[?contains(name, `SQL`) && contains(name, `DB`)]".name --output table"
func resourceMsSqlDatabaseMaintenanceNames() []string {
	return []string{
		"SQL_Default", "SQL_EastUS_DB_1", "SQL_EastUS2_DB_1", "SQL_SoutheastAsia_DB_1", "SQL_AustraliaEast_DB_1", "SQL_NorthEurope_DB_1", "SQL_SouthCentralUS_DB_1", "SQL_WestUS2_DB_1",
		"SQL_UKSouth_DB_1", "SQL_WestEurope_DB_1", "SQL_EastUS_DB_2", "SQL_EastUS2_DB_2", "SQL_WestUS2_DB_2", "SQL_SoutheastAsia_DB_2", "SQL_AustraliaEast_DB_2", "SQL_NorthEurope_DB_2", "SQL_SouthCentralUS_DB_2",
		"SQL_UKSouth_DB_2", "SQL_WestEurope_DB_2", "SQL_AustraliaSoutheast_DB_1", "SQL_BrazilSouth_DB_1", "SQL_CanadaCentral_DB_1", "SQL_CanadaEast_DB_1", "SQL_CentralUS_DB_1", "SQL_EastAsia_DB_1",
		"SQL_FranceCentral_DB_1", "SQL_GermanyWestCentral_DB_1", "SQL_CentralIndia_DB_1", "SQL_SouthIndia_DB_1", "SQL_JapanEast_DB_1", "SQL_JapanWest_DB_1", "SQL_NorthCentralUS_DB_1", "SQL_UKWest_DB_1",
		"SQL_WestUS_DB_1", "SQL_AustraliaSoutheast_DB_2", "SQL_BrazilSouth_DB_2", "SQL_CanadaCentral_DB_2", "SQL_CanadaEast_DB_2", "SQL_CentralUS_DB_2", "SQL_EastAsia_DB_2", "SQL_FranceCentral_DB_2",
		"SQL_GermanyWestCentral_DB_2", "SQL_CentralIndia_DB_2", "SQL_SouthIndia_DB_2", "SQL_JapanEast_DB_2", "SQL_JapanWest_DB_2", "SQL_NorthCentralUS_DB_2", "SQL_UKWest_DB_2", "SQL_WestUS_DB_2",
		"SQL_WestCentralUS_DB_1", "SQL_FranceSouth_DB_1", "SQL_WestCentralUS_DB_2", "SQL_FranceSouth_DB_2", "SQL_SwitzerlandNorth_DB_1", "SQL_SwitzerlandNorth_DB_2", "SQL_BrazilSoutheast_DB_1",
		"SQL_UAENorth_DB_1", "SQL_BrazilSoutheast_DB_2", "SQL_UAENorth_DB_2", "SQL_SouthAfricaNorth_DB_1", "SQL_SouthAfricaNorth_DB_2", "SQL_WestUS3_DB_1", "SQL_WestUS3_DB_2", "SQL_SwedenCentral_DB_1",
		"SQL_SwedenCentral_DB_2",
	}
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
	resource := map[string]*pluginsdk.Schema{
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

		"enclave_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // TODO: Remove Computed in 4.0
			ValidateFunc: validation.StringInSlice([]string{
				string(databases.AlwaysEncryptedEnclaveTypeVBS),
				string(databases.AlwaysEncryptedEnclaveTypeDefault),
			}, false),
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

		"recovery_point_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"restore_dropped_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.RestorableDatabaseID,
		},

		"restore_long_term_retention_backup_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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

		"identity": commonschema.UserAssignedIdentityOptional(),

		"transparent_data_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"transparent_data_encryption_key_vault_key_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: keyVaultValidate.NestedItemId,
		},

		"transparent_data_encryption_key_automatic_rotation_enabled": {
			Type:         pluginsdk.TypeBool,
			Optional:     true,
			Default:      false,
			RequiredWith: []string{"transparent_data_encryption_key_vault_key_id"},
		},

		"secondary_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// This must be Computed as it has defaulted to Geo for replicas but not all databases are replicas.
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(databases.PossibleValuesForSecondaryType(), false),
		},

		"tags": commonschema.Tags(),
	}

	if !features.FivePointOh() {
		atLeastOneOf := []string{
			"long_term_retention_policy.0.weekly_retention", "long_term_retention_policy.0.monthly_retention",
			"long_term_retention_policy.0.yearly_retention", "long_term_retention_policy.0.week_of_year",
		}
		resource["long_term_retention_policy"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
					"weekly_retention": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: helperValidate.ISO8601Duration,
						AtLeastOneOf: atLeastOneOf,
					},

					// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
					"monthly_retention": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: helperValidate.ISO8601Duration,
						AtLeastOneOf: atLeastOneOf,
					},

					// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
					"yearly_retention": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: helperValidate.ISO8601Duration,
						AtLeastOneOf: atLeastOneOf,
					},

					// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
					"week_of_year": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IntBetween(0, 52),
						AtLeastOneOf: atLeastOneOf,
					},

					"immutable_backups_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		}
	}

	return resource
}
