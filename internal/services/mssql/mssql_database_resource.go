package mssql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMsSqlDatabase() *pluginsdk.Resource {
	resourceData := &pluginsdk.Resource{
		Create: resourceMsSqlDatabaseCreateUpdate,
		Read:   resourceMsSqlDatabaseRead,
		Update: resourceMsSqlDatabaseCreateUpdate,
		Delete: resourceMsSqlDatabaseDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.DatabaseID(id)
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
				if !features.ThreePointOhBeta() {
					return nil
				}
				sku := d.Get("sku_name").(string)
				if !strings.HasPrefix(sku, "DW") && !d.Get("transparent_data_encryption_enabled").(bool) {
					return fmt.Errorf("transparent data encryption can only be disabled on Data Warehouse SKUs")
				}
				return nil
			}),
	}
	if features.ThreePointOhBeta() {
		// TODO: Update docs with the following text:
		//
		// * `transparent_data_encryption_enabled` - If set to true, Transparent Data Encryption will be enabled on the database.
		// -> **NOTE:** TDE cannot be disabled on servers with SKUs other than ones starting with DW.

		resourceData.Schema["transparent_data_encryption_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		}
	}
	return resourceData
}

func resourceMsSqlDatabaseImporter(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	replicationLinksClient := meta.(*clients.Client).MSSQL.ReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return nil, err
	}

	partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, client, replicationLinksClient, resourcesClient, *id, []sql.ReplicationRole{sql.ReplicationRolePrimary})
	if err != nil {
		return nil, err
	}

	if len(partnerDatabases) > 0 {
		partnerDatabase := partnerDatabases[0]

		partnerDatabaseId, err := parse.DatabaseID(*partnerDatabase.ID)
		if err != nil {
			return nil, fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
		}

		d.Set("create_mode", string(sql.CreateModeSecondary))
		d.Set("creation_source_database_id", partnerDatabaseId.ID())

		return []*pluginsdk.ResourceData{d}, nil
	}

	d.Set("create_mode", string(sql.CreateModeDefault))

	return []*pluginsdk.ResourceData{d}, nil
}

func resourceMsSqlDatabaseCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	auditingClient := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	serversClient := meta.(*clients.Client).MSSQL.ServersClient
	securityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	replicationLinksClient := meta.(*clients.Client).MSSQL.ReplicationLinksClient
	resourcesClient := meta.(*clients.Client).Resource.ResourcesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database creation.")

	if strings.HasPrefix(d.Get("sku_name").(string), "GP_S_") && d.Get("license_type").(string) != "" {
		return fmt.Errorf("serverless databases do not support license type")
	}

	name := d.Get("name").(string)

	serverId, err := parse.ServerID(d.Get("server_id").(string))
	if err != nil {
		return fmt.Errorf("parsing server ID: %+v", err)
	}

	id := parse.NewDatabaseID(serverId.SubscriptionId, serverId.ResourceGroup, serverId.Name, name)

	if d.IsNewResource() {
		if existing, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name); err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		} else {
			return tf.ImportAsExistsError("azurerm_mssql_database", id.ID())
		}
	}

	server, err := serversClient.Get(ctx, serverId.ResourceGroup, serverId.Name, "")
	if err != nil {
		return fmt.Errorf("making Read request on MsSql Server %q (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
	}

	if server.Location == nil || *server.Location == "" {
		return fmt.Errorf("reading %s: Location was nil/empoty", serverId)
	}
	location := *server.Location

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

	if skuName := d.Get("sku_name"); !d.IsNewResource() && d.HasChange("sku_name") && skuName != "" {
		partnerDatabases, err := helper.FindDatabaseReplicationPartners(ctx, client, replicationLinksClient, resourcesClient, id, []sql.ReplicationRole{sql.ReplicationRoleSecondary, sql.ReplicationRoleNonReadableSecondary})
		if err != nil {
			return err
		}

		// Place a lock for the partner databases, so they can't update themselves whilst we're poking their SKUs
		for _, partnerDatabase := range partnerDatabases {
			partnerDatabaseId, err := parse.DatabaseID(*partnerDatabase.ID)
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
			}

			locks.ByID(partnerDatabaseId.ID())
			defer locks.UnlockByID(partnerDatabaseId.ID())
		}

		// Update the SKUs of any partner databases where deemed necessary
		for _, partnerDatabase := range partnerDatabases {
			partnerDatabaseId, err := parse.DatabaseID(*partnerDatabase.ID)
			if err != nil {
				return fmt.Errorf("parsing ID for Replication Partner Database %q: %+v", *partnerDatabase.ID, err)
			}

			// See: https://docs.microsoft.com/en-us/azure/azure-sql/database/active-geo-replication-overview#configuring-secondary-database
			if partnerDatabase.Sku != nil && partnerDatabase.Sku.Name != nil && helper.CompareDatabaseSkuServiceTiers(skuName.(string), *partnerDatabase.Sku.Name) {
				future, err := client.Update(ctx, partnerDatabaseId.ResourceGroup, partnerDatabaseId.ServerName, partnerDatabaseId.Name, sql.DatabaseUpdate{
					Sku: &sql.Sku{
						Name: utils.String(skuName.(string)),
					},
				})
				if err != nil {
					return fmt.Errorf("updating SKU of Replication Partner %s: %+v", partnerDatabaseId, err)
				}

				if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
					return fmt.Errorf("waiting for SKU update for Replication Partner %s: %+v", partnerDatabaseId, err)
				}
			}
		}
	}

	params := sql.Database{
		Name:     &name,
		Location: &location,
		DatabaseProperties: &sql.DatabaseProperties{
			AutoPauseDelay:                   utils.Int32(int32(d.Get("auto_pause_delay_in_minutes").(int))),
			Collation:                        utils.String(d.Get("collation").(string)),
			ElasticPoolID:                    utils.String(d.Get("elastic_pool_id").(string)),
			LicenseType:                      sql.DatabaseLicenseType(d.Get("license_type").(string)),
			MinCapacity:                      utils.Float(d.Get("min_capacity").(float64)),
			HighAvailabilityReplicaCount:     utils.Int32(int32(d.Get("read_replica_count").(int))),
			SampleName:                       sql.SampleName(d.Get("sample_name").(string)),
			RequestedBackupStorageRedundancy: expandMsSqlBackupStorageRedundancy(d.Get("storage_account_type").(string)),
			ZoneRedundant:                    utils.Bool(d.Get("zone_redundant").(bool)),
		},

		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	createMode, ok := d.GetOk("create_mode")
	if _, dbok := d.GetOk("creation_source_database_id"); ok && (createMode.(string) == string(sql.CreateModeCopy) || createMode.(string) == string(sql.CreateModePointInTimeRestore) || createMode.(string) == string(sql.CreateModeSecondary)) && !dbok {
		return fmt.Errorf("'creation_source_database_id' is required for create_mode %s", createMode.(string))
	}
	if _, dbok := d.GetOk("recover_database_id"); ok && createMode.(string) == string(sql.CreateModeRecovery) && !dbok {
		return fmt.Errorf("'recover_database_id' is required for create_mode %s", createMode.(string))
	}
	if _, dbok := d.GetOk("restore_dropped_database_id"); ok && createMode.(string) == string(sql.CreateModeRestore) && !dbok {
		return fmt.Errorf("'restore_dropped_database_id' is required for create_mode %s", createMode.(string))
	}

	params.DatabaseProperties.CreateMode = sql.CreateMode(createMode.(string))

	auditingPolicies := d.Get("extended_auditing_policy").([]interface{})
	if (createMode == string(sql.CreateModeOnlineSecondary) || createMode == string(sql.CreateModeSecondary)) && len(auditingPolicies) > 0 {
		return fmt.Errorf("cannot configure `extended_auditing_policy` in secondary create mode for %s", id)
	}

	if v, ok := d.GetOk("max_size_gb"); ok {
		// `max_size_gb` is Computed, so has a value after the first run
		if createMode != string(sql.CreateModeOnlineSecondary) && createMode != string(sql.CreateModeSecondary) {
			params.DatabaseProperties.MaxSizeBytes = utils.Int64(int64(v.(int) * 1073741824))
		}
		// `max_size_gb` only has change if it is configured
		if d.HasChange("max_size_gb") && (createMode == string(sql.CreateModeOnlineSecondary) || createMode == string(sql.CreateModeSecondary)) {
			return fmt.Errorf("it is not possible to change maximum size nor advised to configure maximum size in secondary create mode for %s", id)
		}
	}

	readScale := sql.DatabaseReadScaleDisabled
	if v := d.Get("read_scale").(bool); v {
		readScale = sql.DatabaseReadScaleEnabled
	}
	params.DatabaseProperties.ReadScale = readScale

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		if cm, ok := d.GetOk("create_mode"); ok && cm.(string) != string(sql.CreateModePointInTimeRestore) {
			return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(sql.CreateModePointInTimeRestore))
		}
		restorePointInTime, err := time.Parse(time.RFC3339, v.(string))
		if err != nil {
			return fmt.Errorf("parsing `restore_point_in_time` value %q for %s: %+v", v, id, err)
		}
		params.DatabaseProperties.RestorePointInTime = &date.Time{Time: restorePointInTime}
	}

	skuName, ok := d.GetOk("sku_name")
	if ok {
		params.Sku = &sql.Sku{
			Name: utils.String(skuName.(string)),
		}
	}

	if v, ok := d.GetOk("creation_source_database_id"); ok {
		params.DatabaseProperties.SourceDatabaseID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("recover_database_id"); ok {
		params.DatabaseProperties.RecoverableDatabaseID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("restore_dropped_database_id"); ok {
		params.DatabaseProperties.RestorableDroppedDatabaseID = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, params)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of %s: %+v", id, err)
	}

	if features.ThreePointOhBeta() {
		statusProperty := sql.TransparentDataEncryptionStatusDisabled
		encryptionStatus := d.Get("transparent_data_encryption_enabled").(bool)
		if encryptionStatus {
			statusProperty = sql.TransparentDataEncryptionStatusEnabled
		}
		_, err := transparentEncryptionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, sql.TransparentDataEncryption{
			TransparentDataEncryptionProperties: &sql.TransparentDataEncryptionProperties{
				Status: statusProperty,
			},
		})
		if err != nil {
			return fmt.Errorf("while enabling Transparent Data Encryption for %q: %+v", id.String(), err)
		}

		if err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
			c, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("while polling cluster %s for status: %+v", id.String(), err))
			}
			if c.DatabaseProperties.Status == sql.DatabaseStatusScaling {
				return resource.RetryableError(fmt.Errorf("database %s is still scaling", id.String()))
			}

			return nil
		}); err != nil {
			return nil
		}
	}

	d.SetId(id.ID())

	// For datawarehouse SKUs only
	if strings.HasPrefix(skuName.(string), "DW") && (d.HasChange("geo_backup_enabled") || d.IsNewResource()) {
		isEnabled := d.Get("geo_backup_enabled").(bool)
		var geoBackupPolicyState sql.GeoBackupPolicyState

		// The default geo backup policy configuration for a new resource is 'enabled', so we don't need to set it in that scenario
		if !(d.IsNewResource() && isEnabled) {
			if isEnabled {
				geoBackupPolicyState = sql.GeoBackupPolicyStateEnabled
			} else {
				geoBackupPolicyState = sql.GeoBackupPolicyStateDisabled
			}

			geoBackupPolicy := sql.GeoBackupPolicy{
				GeoBackupPolicyProperties: &sql.GeoBackupPolicyProperties{
					State: geoBackupPolicyState,
				},
			}

			if _, err := geoBackupPoliciesClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, geoBackupPolicy); err != nil {
				return fmt.Errorf("setting Geo Backup Policies for %s: %+v", id, err)
			}
		}
	}

	if _, err = securityAlertPoliciesClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, expandMsSqlServerSecurityAlertPolicy(d)); err != nil {
		return fmt.Errorf("setting database threat detection policy for %s: %+v", id, err)
	}

	if createMode != string(sql.CreateModeOnlineSecondary) && createMode != string(sql.CreateModeSecondary) {
		auditingProps := sql.ExtendedDatabaseBlobAuditingPolicy{
			ExtendedDatabaseBlobAuditingPolicyProperties: helper.ExpandMsSqlDBBlobAuditingPolicies(auditingPolicies),
		}
		if _, err = auditingClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, auditingProps); err != nil {
			return fmt.Errorf("setting Blob Auditing Policies for %s: %+v", id, err)
		}
	}

	if d.HasChange("long_term_retention_policy") {
		v := d.Get("long_term_retention_policy")
		longTermRetentionProps := helper.ExpandLongTermRetentionPolicy(v.([]interface{}))
		if longTermRetentionProps != nil {
			longTermRetentionPolicy := sql.LongTermRetentionPolicy{}

			// hyper-scale SKU's do not support LRP currently
			if !strings.HasPrefix(skuName.(string), "HS") && !strings.HasPrefix(skuName.(string), "DW") {
				longTermRetentionPolicy.BaseLongTermRetentionPolicyProperties = longTermRetentionProps
			}

			longTermRetentionfuture, err := longTermRetentionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, longTermRetentionPolicy)
			if err != nil {
				return fmt.Errorf("setting Long Term Retention Policies for %s: %+v", id, err)
			}

			if err = longTermRetentionfuture.WaitForCompletionRef(ctx, longTermRetentionClient.Client); err != nil {
				return fmt.Errorf("waiting for update of Long Term Retention Policies for %s: %+v", id, err)
			}
		}
	}

	if d.HasChange("short_term_retention_policy") {
		v := d.Get("short_term_retention_policy")
		backupShortTermPolicyProps := helper.ExpandShortTermRetentionPolicy(v.([]interface{}))
		if backupShortTermPolicyProps != nil {
			backupShortTermPolicy := sql.BackupShortTermRetentionPolicy{}

			if !strings.HasPrefix(skuName.(string), "HS") && !strings.HasPrefix(skuName.(string), "DW") {
				backupShortTermPolicy.BackupShortTermRetentionPolicyProperties = backupShortTermPolicyProps
			}

			shortTermRetentionFuture, err := shortTermRetentionClient.CreateOrUpdate(ctx, id.ResourceGroup, id.ServerName, id.Name, backupShortTermPolicy)
			if err != nil {
				return fmt.Errorf("setting Short Term Retention Policies for %s: %+v", id, err)
			}

			if err = shortTermRetentionFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of Short Term Retention Policies for %s: %+v", id, err)
			}
		}
	}

	return resourceMsSqlDatabaseRead(d, meta)
}

func resourceMsSqlDatabaseRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	securityAlertPoliciesClient := meta.(*clients.Client).MSSQL.DatabaseSecurityAlertPoliciesClient
	auditingClient := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.LongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	geoBackupPoliciesClient := meta.(*clients.Client).MSSQL.GeoBackupPoliciesClient
	transparentEncryptionClient := meta.(*clients.Client).MSSQL.TransparentDataEncryptionsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	serverId := parse.NewServerID(id.SubscriptionId, id.ResourceGroup, id.ServerName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("server_id", serverId.ID())

	skuName := ""
	if props := resp.DatabaseProperties; props != nil {
		d.Set("auto_pause_delay_in_minutes", props.AutoPauseDelay)
		d.Set("collation", props.Collation)
		d.Set("elastic_pool_id", props.ElasticPoolID)
		d.Set("license_type", props.LicenseType)
		if props.MaxSizeBytes != nil {
			d.Set("max_size_gb", int32((*props.MaxSizeBytes)/int64(1073741824)))
		}
		d.Set("min_capacity", props.MinCapacity)
		d.Set("read_replica_count", props.HighAvailabilityReplicaCount)
		if props.ReadScale == sql.DatabaseReadScaleEnabled {
			d.Set("read_scale", true)
		} else if props.ReadScale == sql.DatabaseReadScaleDisabled {
			d.Set("read_scale", false)
		}
		if props.CurrentServiceObjectiveName != nil {
			skuName = *props.CurrentServiceObjectiveName
		}
		d.Set("sku_name", skuName)
		d.Set("storage_account_type", flattenMsSqlBackupStorageRedundancy(props.CurrentBackupStorageRedundancy))
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	securityAlertPolicy, err := securityAlertPoliciesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err == nil {
		if err := d.Set("threat_detection_policy", flattenMsSqlServerSecurityAlertPolicy(d, securityAlertPolicy)); err != nil {
			return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
		}
	}

	extendedAuditingPolicy := []interface{}{}
	if createMode, ok := d.GetOk("create_mode"); !ok || (createMode.(string) != "Secondary" && createMode.(string) != "OnlineSecondary") {
		auditingResp, err := auditingClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Blob Auditing Policies for %s: %+v", id, err)
		}

		extendedAuditingPolicy = helper.FlattenMsSqlDBBlobAuditingPolicies(&auditingResp, d)
	}
	d.Set("extended_auditing_policy", extendedAuditingPolicy)

	geoBackupPolicy := true

	// Hyper Scale SKU's do not currently support LRP and do not honour normal SRP operations
	if !strings.HasPrefix(skuName, "HS") && !strings.HasPrefix(skuName, "DW") {
		longTermPolicy, err := longTermRetentionClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Long Term Retention Policies for %s: %+v", id, err)
		}
		if err := d.Set("long_term_retention_policy", helper.FlattenLongTermRetentionPolicy(&longTermPolicy, d)); err != nil {
			return fmt.Errorf("setting `long_term_retention_policy`: %+v", err)
		}

		shortTermPolicy, err := shortTermRetentionClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			return fmt.Errorf("retrieving Short Term Retention Policies for %s: %+v", id, err)
		}

		if err := d.Set("short_term_retention_policy", helper.FlattenShortTermRetentionPolicy(&shortTermPolicy, d)); err != nil {
			return fmt.Errorf("setting `short_term_retention_policy`: %+v", err)
		}
	} else {
		// HS and DW SKUs need the retention policies zeroing for state consistency
		zero := make([]interface{}, 0)
		d.Set("long_term_retention_policy", zero)
		d.Set("short_term_retention_policy", zero)

		geoPoliciesResponse, err := geoBackupPoliciesClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				d.SetId("")
				return nil
			}
			return fmt.Errorf("retrieving Geo Backup Policies for %s: %+v", id, err)
		}

		// For Datawarehouse SKUs, set the geo-backup policy setting
		if strings.HasPrefix(skuName, "DW") && geoPoliciesResponse.GeoBackupPolicyProperties.State == sql.GeoBackupPolicyStateDisabled {
			geoBackupPolicy = false
		}
	}

	if err := d.Set("geo_backup_enabled", geoBackupPolicy); err != nil {
		return fmt.Errorf("setting `geo_backup_enabled`: %+v", err)
	}

	if features.ThreePointOhBeta() {
		tde, err := transparentEncryptionClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			return fmt.Errorf("while retrieving Transparent Data Encryption status of %q: %+v", id.String(), err)
		}
		tdeStatus := false
		if tde.TransparentDataEncryptionProperties != nil && tde.TransparentDataEncryptionProperties.Status == sql.TransparentDataEncryptionStatusEnabled {
			tdeStatus = true
		}
		d.Set("transparent_data_encryption_enabled", tdeStatus)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMsSqlDatabaseDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting deletion of %s: %+v", id, err)
	}

	return nil
}

func flattenMsSqlServerSecurityAlertPolicy(d *pluginsdk.ResourceData, policy sql.DatabaseSecurityAlertPolicy) []interface{} {
	// The SQL database security alert API always returns the default value even if never set.
	// If the values are on their default one, threat it as not set.
	properties := policy.SecurityAlertsPolicyProperties
	if properties == nil {
		return []interface{}{}
	}

	securityAlertPolicy := make(map[string]interface{})

	securityAlertPolicy["state"] = string(properties.State)
	if !features.ThreePointOhBeta() {
		securityAlertPolicy["use_server_default"] = "Disabled"
	}

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

func expandMsSqlServerSecurityAlertPolicy(d *pluginsdk.ResourceData) sql.DatabaseSecurityAlertPolicy {
	policy := sql.DatabaseSecurityAlertPolicy{
		SecurityAlertsPolicyProperties: &sql.SecurityAlertsPolicyProperties{
			State: sql.SecurityAlertsPolicyStateDisabled,
		},
	}
	properties := policy.SecurityAlertsPolicyProperties

	td, ok := d.GetOk("threat_detection_policy")
	if !ok {
		return policy
	}

	if tdl := td.([]interface{}); len(tdl) > 0 {
		securityAlert := tdl[0].(map[string]interface{})

		properties.State = sql.SecurityAlertsPolicyState(securityAlert["state"].(string))
		properties.EmailAccountAdmins = utils.Bool(securityAlert["email_account_admins"].(string) == "Enabled")

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
			properties.RetentionDays = utils.Int32(int32(v.(int)))
		}
		if v, ok := securityAlert["storage_account_access_key"]; ok {
			properties.StorageAccountAccessKey = utils.String(v.(string))
		}
		if v, ok := securityAlert["storage_endpoint"]; ok {
			properties.StorageEndpoint = utils.String(v.(string))
		}

		return policy
	}

	return policy
}

func flattenMsSqlBackupStorageRedundancy(currentBackupStorageRedundancy sql.CurrentBackupStorageRedundancy) string {
	if !features.ThreePointOhBeta() {
		switch currentBackupStorageRedundancy {
		case sql.CurrentBackupStorageRedundancyLocal:
			return "LRS"
		case sql.CurrentBackupStorageRedundancyZone:
			return "ZRS"
		default:
			return "GRS"
		}
	}
	return string(currentBackupStorageRedundancy)
}

func expandMsSqlBackupStorageRedundancy(storageAccountType string) sql.RequestedBackupStorageRedundancy {
	if !features.ThreePointOhBeta() {
		switch storageAccountType {
		case "LRS":
			return sql.RequestedBackupStorageRedundancyLocal
		case "ZRS":
			return sql.RequestedBackupStorageRedundancyZone
		default:
			return sql.RequestedBackupStorageRedundancyGeo
		}
	}
	return sql.RequestedBackupStorageRedundancy(storageAccountType)
}

func resourceMsSqlDatabaseSchema() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
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
			Default:  string(sql.CreateModeDefault),
			ValidateFunc: validation.StringInSlice([]string{
				string(sql.CreateModeCopy),
				string(sql.CreateModeDefault),
				string(sql.CreateModeOnlineSecondary),
				string(sql.CreateModePointInTimeRestore),
				string(sql.CreateModeRestore),
				string(sql.CreateModeRecovery),
				string(sql.CreateModeRestoreExternalBackup),
				string(sql.CreateModeRestoreExternalBackupSecondary),
				string(sql.CreateModeRestoreLongTermRetentionBackup),
				string(sql.CreateModeSecondary),
			}, false),
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

		"extended_auditing_policy": helper.ExtendedAuditingSchema(),

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(sql.DatabaseLicenseTypeBasePrice),
				string(sql.DatabaseLicenseTypeLicenseIncluded),
			}, false),
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
			ValidateFunc: azValidate.FloatInSlice([]float64{0, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40}),
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
				string(sql.SampleNameAdventureWorksLT),
			}, false),
		},

		"sku_name": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validate.DatabaseSkuName(),
			DiffSuppressFunc: suppress.CaseDifferenceV2Only,
		},

		"creation_source_database_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validate.DatabaseID,
		},

		"storage_account_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default: func() string {
				if !features.ThreePointOhBeta() {
					return "GRS"
				}
				return string(sql.CurrentBackupStorageRedundancyGeo)
			}(),
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				if !features.ThreePointOhBeta() {
					return validation.StringInSlice([]string{
						"GRS",
						"LRS",
						"ZRS",
					}, false)
				}
				return validation.StringInSlice([]string{
					string(sql.CurrentBackupStorageRedundancyGeo),
					string(sql.CurrentBackupStorageRedundancyLocal),
					string(sql.CurrentBackupStorageRedundancyZone),
				}, false)
			}(),
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
							}, !features.ThreePointOh()),
							DiffSuppressFunc: suppress.CaseDifferenceV2Only,
						},
					},

					"email_account_admins": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifferenceV2Only,
						Default:          "Disabled",
						ValidateFunc: validation.StringInSlice([]string{
							"Disabled",
							"Enabled",
						}, !features.ThreePointOh()),
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

					"state": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						DiffSuppressFunc: suppress.CaseDifferenceV2Only,
						Default:          string(sql.SecurityAlertPolicyStateDisabled),
						ValidateFunc: validation.StringInSlice([]string{
							string(sql.SecurityAlertPolicyStateDisabled),
							string(sql.SecurityAlertPolicyStateEnabled),
							string(sql.SecurityAlertPolicyStateNew),
						}, !features.ThreePointOh()),
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
		"tags": tags.Schema(),
	}

	if !features.ThreePointOhBeta() {
		s := out["threat_detection_policy"].Elem.(*schema.Resource)
		s.Schema["use_server_default"] = &pluginsdk.Schema{
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			Default:          "Disabled",
			ValidateFunc: validation.StringInSlice([]string{
				"Disabled",
				"Enabled",
			}, true),
			Deprecated: "This field is now non-functional and thus will be removed in version 3.0 of the Azure Provider",
		}
	}

	return out
}
