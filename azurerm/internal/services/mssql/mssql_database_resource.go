package mssql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlDatabaseCreateUpdate,
		Read:   resourceArmMsSqlDatabaseRead,
		Update: resourceArmMsSqlDatabaseCreateUpdate,
		Delete: resourceArmMsSqlDatabaseDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DatabaseID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MsSqlServerID,
			},

			"auto_pause_delay_in_minutes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.MsSqlDatabaseAutoPauseDelay,
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.MsSqlDBCollation(),
			},

			"elastic_pool_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.MsSqlElasticPoolID,
			},

			"extended_auditing_policy": helper.ExtendedAuditingSchema(),

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.BasePrice),
					string(sql.LicenseIncluded),
				}, false),
			},

			"long_term_retention_policy": helper.LongTermRetentionPolicySchema(),

			"short_term_retention_policy": helper.ShortTermRetentionPolicySchema(),

			"max_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 4096),
			},

			"min_capacity": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Computed:     true,
				ValidateFunc: azValidate.FloatInSlice([]float64{0, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2, 2.25, 2.5, 3, 4, 5, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40}),
			},

			"restore_point_in_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
			},

			"recover_database_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.MsSqlRecoverableDatabaseID,
			},

			"restore_dropped_database_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.MsSqlRestorableDatabaseID,
			},

			"read_replica_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 4),
			},

			"read_scale": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"sample_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.AdventureWorksLT),
				}, false),
			},

			"sku_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validate.MsSqlDBSkuName(),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"creation_source_database_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validate.MsSqlDatabaseID,
			},

			"zone_redundant": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"threat_detection_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disabled_alerts": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Sql_Injection",
									"Sql_Injection_Vulnerability",
									"Access_Anomaly",
								}, true),
							},
						},

						"email_account_admins": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							Default:          string(sql.SecurityAlertPolicyEmailAccountAdminsDisabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.SecurityAlertPolicyEmailAccountAdminsDisabled),
								string(sql.SecurityAlertPolicyEmailAccountAdminsEnabled),
							}, true),
						},

						"email_addresses": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"retention_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"state": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							Default:          string(sql.SecurityAlertPolicyStateDisabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.SecurityAlertPolicyStateDisabled),
								string(sql.SecurityAlertPolicyStateEnabled),
								string(sql.SecurityAlertPolicyStateNew),
							}, true),
						},

						"storage_account_access_key": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"storage_endpoint": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"use_server_default": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							Default:          string(sql.SecurityAlertPolicyUseServerDefaultDisabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.SecurityAlertPolicyUseServerDefaultDisabled),
								string(sql.SecurityAlertPolicyUseServerDefaultEnabled),
							}, true),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("sku_name", func(old, new, meta interface{}) bool {
				// "hyperscale can not change to other sku
				return strings.HasPrefix(old.(string), "HS") && !strings.HasPrefix(new.(string), "HS")
			}),
		),
	}
}

func resourceArmMsSqlDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	auditingClient := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	serverClient := meta.(*clients.Client).MSSQL.ServersClient
	threatClient := meta.(*clients.Client).MSSQL.DatabaseThreatDetectionPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.BackupLongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database creation.")

	name := d.Get("name").(string)
	sqlServerId := d.Get("server_id").(string)
	serverId, _ := parse.ServerID(sqlServerId)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Failed to check for presence of existing Database %q (MsSql Server %q / Resource Group %q): %s", name, serverId.Name, serverId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_database", *existing.ID)
		}
	}

	serverResp, err := serverClient.Get(ctx, serverId.ResourceGroup, serverId.Name)
	if err != nil {
		return fmt.Errorf("making Read request on MsSql Server %q (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
	}

	location := *serverResp.Location
	if location == "" {
		return fmt.Errorf("Location is empty from making Read request on MsSql Server %q", serverId.Name)
	}

	// when disassociating mssql db from elastic pool, the sku_name must be specific
	if d.HasChange("elastic_pool_id") {
		if old, new := d.GetChange("elastic_pool_id"); old.(string) != "" && new.(string) == "" {
			if v, ok := d.GetOk("sku_name"); !ok || (ok && v.(string) == "ElasticPool") {
				return fmt.Errorf("`sku_name` must be assigned and not be `ElasticPool` when disassociating MsSql Database %q from MsSql Elastic Pool", name)
			}
		}
	}
	params := sql.Database{
		Name:     &name,
		Location: &location,
		DatabaseProperties: &sql.DatabaseProperties{
			AutoPauseDelay:   utils.Int32(int32(d.Get("auto_pause_delay_in_minutes").(int))),
			Collation:        utils.String(d.Get("collation").(string)),
			ElasticPoolID:    utils.String(d.Get("elastic_pool_id").(string)),
			LicenseType:      sql.DatabaseLicenseType(d.Get("license_type").(string)),
			MinCapacity:      utils.Float(d.Get("min_capacity").(float64)),
			ReadReplicaCount: utils.Int32(int32(d.Get("read_replica_count").(int))),
			SampleName:       sql.SampleName(d.Get("sample_name").(string)),
			ZoneRedundant:    utils.Bool(d.Get("zone_redundant").(bool)),
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
	if (createMode == string(sql.CreateModeOnlineSecondary) || createMode == string(sql.Secondary)) && len(auditingPolicies) > 0 {
		return fmt.Errorf("could not configure auditing policies on SQL Database %q (Resource Group %q, Server %q) in secondary create mode", name, serverId.ResourceGroup, serverId.Name)
	}

	if v, ok := d.GetOk("max_size_gb"); ok {
		params.DatabaseProperties.MaxSizeBytes = utils.Int64(int64(v.(int) * 1073741824))
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
		restorePointInTime, _ := time.Parse(time.RFC3339, v.(string))
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

	future, err := client.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, params)
	if err != nil {
		return fmt.Errorf("creating MsSql Database %q (Sql Server %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of MsSql Database %q (MsSql Server Name %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving MsSql Database %q (MsSql Server Name %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("Cannot read MsSql Database %q (MsSql Server Name %q / Resource Group %q) ID", name, serverId.Name, serverId.ResourceGroup)
	}

	d.SetId(*read.ID)

	if _, err = threatClient.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, *expandArmMsSqlServerThreatDetectionPolicy(d, location)); err != nil {
		return fmt.Errorf("setting database threat detection policy: %+v", err)
	}

	if createMode != string(sql.CreateModeOnlineSecondary) && createMode != string(sql.CreateModeSecondary) {
		auditingProps := sql.ExtendedDatabaseBlobAuditingPolicy{
			ExtendedDatabaseBlobAuditingPolicyProperties: helper.ExpandAzureRmMsSqlDBBlobAuditingPolicies(auditingPolicies),
		}
		if _, err = auditingClient.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, auditingProps); err != nil {
			return fmt.Errorf("failure in issuing create/update request for SQL Database %q Blob Auditing Policies(SQL Server %q/ Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
		}
	}

	// hyper-scale SKU's do not support LRP currently
	if d.HasChange("long_term_retention_policy") {
		v := d.Get("long_term_retention_policy")
		longTermRetentionProps := helper.ExpandLongTermRetentionPolicy(v.([]interface{}))
		if longTermRetentionProps != nil {
			longTermRetentionPolicy := sql.BackupLongTermRetentionPolicy{}

			if !strings.HasPrefix(skuName.(string), "HS") && !strings.HasPrefix(skuName.(string), "DW") {
				longTermRetentionPolicy.LongTermRetentionPolicyProperties = longTermRetentionProps
			}

			longTermRetentionfuture, err := longTermRetentionClient.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, longTermRetentionPolicy)
			if err != nil {
				return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverId.Name, name, serverId.ResourceGroup, err)
			}

			if err = longTermRetentionfuture.WaitForCompletionRef(ctx, longTermRetentionClient.Client); err != nil {
				return fmt.Errorf("Error waiting for completion of Create/Update for Sql Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverId.Name, name, serverId.ResourceGroup, err)
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

			shortTermRetentionFuture, err := shortTermRetentionClient.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, backupShortTermPolicy)
			if err != nil {
				return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverId.Name, name, serverId.ResourceGroup, err)
			}

			if err = shortTermRetentionFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("Error waiting for completion of Create/Update for Sql Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverId.Name, name, serverId.ResourceGroup, err)
			}
		}
	}

	return resourceArmMsSqlDatabaseRead(d, meta)
}

func resourceArmMsSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	threatClient := meta.(*clients.Client).MSSQL.DatabaseThreatDetectionPoliciesClient
	auditingClient := meta.(*clients.Client).MSSQL.DatabaseExtendedBlobAuditingPoliciesClient
	longTermRetentionClient := meta.(*clients.Client).MSSQL.BackupLongTermRetentionPoliciesClient
	shortTermRetentionClient := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading MsSql Database %s (MsSql Server Name %q / Resource Group %q): %s", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)

	serverClient := meta.(*clients.Client).MSSQL.ServersClient

	serverResp, err := serverClient.Get(ctx, id.ResourceGroup, id.ServerName)
	if err != nil || *serverResp.ID == "" {
		return fmt.Errorf("making Read request on MsSql Server  %q (Resource Group %q): %s", id.ServerName, id.ResourceGroup, err)
	}
	d.Set("server_id", serverResp.ID)

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
		d.Set("read_replica_count", props.ReadReplicaCount)
		if props.ReadScale == sql.DatabaseReadScaleEnabled {
			d.Set("read_scale", true)
		} else if props.ReadScale == sql.DatabaseReadScaleDisabled {
			d.Set("read_scale", false)
		}
		if props.CurrentServiceObjectiveName != nil {
			skuName = *props.CurrentServiceObjectiveName
		}
		d.Set("sku_name", props.CurrentServiceObjectiveName)
		d.Set("zone_redundant", props.ZoneRedundant)
	}

	threat, err := threatClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err == nil {
		if err := d.Set("threat_detection_policy", flattenArmMsSqlServerThreatDetectionPolicy(d, threat)); err != nil {
			return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
		}
	}

	auditingResp, err := auditingClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("failure in reading SQL Database %q: %v Blob Auditing Policies", id.Name, err)
	}

	flattenBlobAuditing := helper.FlattenAzureRmMsSqlDBBlobAuditingPolicies(&auditingResp, d)
	if err := d.Set("extended_auditing_policy", flattenBlobAuditing); err != nil {
		return fmt.Errorf("failure in setting `extended_auditing_policy`: %+v", err)
	}

	// Hyper Scale SKU's do not currently support LRP and do not honour normal SRP operations
	if !strings.HasPrefix(skuName, "HS") && !strings.HasPrefix(skuName, "DW") {
		longTermPolicy, err := longTermRetentionClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			return fmt.Errorf("Error retrieving Long Term Policies for Database %q (Sql Server %q ;Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
		}
		flattenlongTermPolicy := helper.FlattenLongTermRetentionPolicy(&longTermPolicy, d)
		if err := d.Set("long_term_retention_policy", flattenlongTermPolicy); err != nil {
			return fmt.Errorf("failure in setting `long_term_retention_policy`: %+v", err)
		}

		shortTermPolicy, err := shortTermRetentionClient.Get(ctx, id.ResourceGroup, id.ServerName, id.Name)
		if err != nil {
			return fmt.Errorf("Error retrieving Short Term Policies for Database %q (Sql Server %q ;Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
		}

		flattenShortTermPolicy := helper.FlattenShortTermRetentionPolicy(&shortTermPolicy, d)
		if err := d.Set("short_term_retention_policy", flattenShortTermPolicy); err != nil {
			return fmt.Errorf("failure in setting `short_term_retention_policy`: %+v", err)
		}
	} else {
		// HS and DW SKUs need the retention policies zeroing for state consistency
		zero := make([]interface{}, 0)
		d.Set("long_term_retention_policy", zero)
		d.Set("short_term_retention_policy", zero)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMsSqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ServerName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting MsSql Database %q ( MsSql Server %q / Resource Group %q): %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("waiting for MsSql Database %q ( MsSql Server %q / Resource Group %q) to be deleted: %+v", id.Name, id.ServerName, id.ResourceGroup, err)
	}

	return nil
}

func flattenArmMsSqlServerThreatDetectionPolicy(d *schema.ResourceData, policy sql.DatabaseSecurityAlertPolicy) []interface{} {
	// The SQL database threat detection API always returns the default value even if never set.
	// If the values are on their default one, threat it as not set.
	properties := policy.DatabaseSecurityAlertPolicyProperties
	if properties == nil {
		return []interface{}{}
	}

	threatDetectionPolicy := make(map[string]interface{})

	threatDetectionPolicy["state"] = string(properties.State)
	threatDetectionPolicy["email_account_admins"] = string(properties.EmailAccountAdmins)
	threatDetectionPolicy["use_server_default"] = string(properties.UseServerDefault)

	if disabledAlerts := properties.DisabledAlerts; disabledAlerts != nil {
		flattenedAlerts := schema.NewSet(schema.HashString, []interface{}{})
		if v := *disabledAlerts; v != "" {
			parsedAlerts := strings.Split(v, ";")
			for _, a := range parsedAlerts {
				flattenedAlerts.Add(a)
			}
		}
		threatDetectionPolicy["disabled_alerts"] = flattenedAlerts
	}
	if emailAddresses := properties.EmailAddresses; emailAddresses != nil {
		flattenedEmails := schema.NewSet(schema.HashString, []interface{}{})
		if v := *emailAddresses; v != "" {
			parsedEmails := strings.Split(*emailAddresses, ";")
			for _, e := range parsedEmails {
				flattenedEmails.Add(e)
			}
		}
		threatDetectionPolicy["email_addresses"] = flattenedEmails
	}
	if properties.StorageEndpoint != nil {
		threatDetectionPolicy["storage_endpoint"] = *properties.StorageEndpoint
	}
	if properties.RetentionDays != nil {
		threatDetectionPolicy["retention_days"] = int(*properties.RetentionDays)
	}

	// If storage account access key is in state read it to the new state, as the API does not return it for security reasons
	if v, ok := d.GetOk("threat_detection_policy.0.storage_account_access_key"); ok {
		threatDetectionPolicy["storage_account_access_key"] = v.(string)
	}

	return []interface{}{threatDetectionPolicy}
}

func expandArmMsSqlServerThreatDetectionPolicy(d *schema.ResourceData, location string) *sql.DatabaseSecurityAlertPolicy {
	policy := sql.DatabaseSecurityAlertPolicy{
		Location: utils.String(location),
		DatabaseSecurityAlertPolicyProperties: &sql.DatabaseSecurityAlertPolicyProperties{
			State: sql.SecurityAlertPolicyStateDisabled,
		},
	}
	properties := policy.DatabaseSecurityAlertPolicyProperties

	td, ok := d.GetOk("threat_detection_policy")
	if !ok {
		return &policy
	}

	if tdl := td.([]interface{}); len(tdl) > 0 {
		threatDetection := tdl[0].(map[string]interface{})

		properties.State = sql.SecurityAlertPolicyState(threatDetection["state"].(string))
		properties.EmailAccountAdmins = sql.SecurityAlertPolicyEmailAccountAdmins(threatDetection["email_account_admins"].(string))
		properties.UseServerDefault = sql.SecurityAlertPolicyUseServerDefault(threatDetection["use_server_default"].(string))

		if v, ok := threatDetection["disabled_alerts"]; ok {
			alerts := v.(*schema.Set).List()
			expandedAlerts := make([]string, len(alerts))
			for i, a := range alerts {
				expandedAlerts[i] = a.(string)
			}
			properties.DisabledAlerts = utils.String(strings.Join(expandedAlerts, ";"))
		}
		if v, ok := threatDetection["email_addresses"]; ok {
			emails := v.(*schema.Set).List()
			expandedEmails := make([]string, len(emails))
			for i, e := range emails {
				expandedEmails[i] = e.(string)
			}
			properties.EmailAddresses = utils.String(strings.Join(expandedEmails, ";"))
		}
		if v, ok := threatDetection["retention_days"]; ok {
			properties.RetentionDays = utils.Int32(int32(v.(int)))
		}
		if v, ok := threatDetection["storage_account_access_key"]; ok {
			properties.StorageAccountAccessKey = utils.String(v.(string))
		}
		if v, ok := threatDetection["storage_endpoint"]; ok {
			properties.StorageEndpoint = utils.String(v.(string))
		}

		return &policy
	}

	return &policy
}
