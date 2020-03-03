package mssql

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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
			_, err := parse.MsSqlDatabaseID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"mssql_server_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"elastic_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BasePrice",
					"LicenseIncluded",
				}, false),
			},

			"sample_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AdventureWorksLT",
					"WideWorldImportersStd",
					"WideWorldImportersFull",
				}, false),
			},

			"create_copy_mode": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"create_pitr_mode": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"restore_point_in_time": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validate.RFC3339Time,
						},
					},
				},
			},

			"create_recovery_mode": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"restorable_dropped_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"create_restore_mode": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"source_database_deletion_date": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validate.RFC3339Time,
						},
					},
				},
			},

			"create_secondary_mode": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"general_purpose": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				ConfigMode: schema.SchemaConfigModeAttr,
				MaxItems:   1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40, 80}),
						},

						"family": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
						},

						"max_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1024),
						},
					},
				},
			},

			"hyper_scale": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 16, 24}),
						},

						"family": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
						},

						"read_replica_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 4),
						},
					},
				},
			},

			"business_critical": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 24, 32, 40, 80}),
						},
						"family": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
								"M-Series",
								"FSv2 Series",
							}, false),
						},

						"max_size_gb": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 1024),
						},

						"read_scale": {
							Type:     schema.TypeString,
							Optional: true,
							Computed:true,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, false),
						},

						"zone_redundant": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed:true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMsSqlDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database creation.")

	name := d.Get("name").(string)
	mssqlServerId := d.Get("mssql_server_id").(string)
	serverId, _ := parse.MsSqlServerID(mssqlServerId)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Database %q (MsSql Server %q / Resource Group %q): %s", name, serverId.Name, serverId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mssql_database ", *existing.ID)
		}
	}

	serverClient := meta.(*clients.Client).MSSQL.ServersClient
	serverResp, err := serverClient.Get(ctx, serverId.ResourceGroup, serverId.Name)
	if err != nil {
		return fmt.Errorf("Error getting existing Server %q (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
	}

	location := *serverResp.Location
	if location == "" {
		return fmt.Errorf("Error reading location of MsSql Server %q", serverId.Name)
	}

	t := d.Get("tags").(map[string]interface{})

	params := sql.Database{
		Name:               &name,
		Location:           &location,
		DatabaseProperties: &sql.DatabaseProperties{},

		Tags: tags.Expand(t),
	}

	expandAzureRmMsSqlDatabaseGP(d, &params)
	expandAzureRmMsSqlDatabaseHS(d, &params)
	expandAzureRmMsSqlDatabaseBC(d, &params)

	expandAzureRmMsSqlDatabaseCreateCopyMode(d, &params)
	expandAzureRmMsSqlDatabaseCreatePITRMode(d, &params)
	expandAzureRmMsSqlDatabaseCreateRecoveryMode(d, &params)
	expandAzureRmMsSqlDatabaseCreateRestoreMode(d, &params)
	expandAzureRmMsSqlDatabaseCreateSecondaryMode(d, &params)

	if v, ok := d.GetOkExists("collation"); ok {
		params.DatabaseProperties.Collation = utils.String(v.(string))
	}

	if v, ok := d.GetOkExists("elastic_pool_id"); ok {
		params.DatabaseProperties.ElasticPoolID = utils.String(v.(string))
	}

	if v, ok := d.GetOkExists("license_type"); ok {
		params.DatabaseProperties.LicenseType = sql.DatabaseLicenseType(v.(string))
	}

	if v, ok := d.GetOkExists("sample_name"); ok {
		params.DatabaseProperties.SampleName = sql.SampleName(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, params)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read MsSql Database %q (resource group %q) ID", name, serverId.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMsSqlDatabaseRead(d, meta)
}

func resourceArmMsSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseId, err := parse.MsSqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, databaseId.ResourceGroup, databaseId.MsSqlServer, databaseId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on MsSql Database %s: %s", databaseId.Name, err)
	}

	d.Set("name", resp.Name)

	serverClient := meta.(*clients.Client).MSSQL.ServersClient

	serverResp, err := serverClient.Get(ctx, databaseId.ResourceGroup, databaseId.MsSqlServer)
	if err != nil || *serverResp.ID == "" {
		return fmt.Errorf("Error getting existing Server %q (Resource Group %q): %s", databaseId.MsSqlServer, databaseId.ResourceGroup, err)
	}
	d.Set("mssql_server_id", serverResp.ID)

	d.Set("collation", resp.Collation)

	d.Set("elastic_pool_id", resp.ElasticPoolID)

	d.Set("license_type", resp.LicenseType)

	d.Set("sample_name", resp.SampleName)

	flattenedGP := flattenAzureRmMsSqlDatabaseGP(&resp)
	if err := d.Set("general_purpose", flattenedGP); err != nil {
		return fmt.Errorf("Error setting `general_purpose`: %+v", err)
	}

	flattenedHS := flattenAzureRmMsSqlDatabaseHS(&resp)
	if err := d.Set("hyper_scale", flattenedHS); err != nil {
		return fmt.Errorf("Error setting `hyper_scale`: %+v", err)
	}

	flattenedBC := flattenAzureRmMsSqlDatabaseBC(&resp)
	if err := d.Set("business_critical", flattenedBC); err != nil {
		return fmt.Errorf("Error setting `business_critical`: %+v", err)
	}

	flattenedCreateCopyMode := flattenAzureRmMsSqlDatabaseCreateCopyMode(&resp)
	if err := d.Set("create_copy_mode", flattenedCreateCopyMode); err != nil {
		return fmt.Errorf("Error setting `create_copy_mode`: %+v", err)
	}

	flattenedCreatePITRMode := flattenAzureRmMsSqlDatabaseCreatePITRMode(&resp)
	if err := d.Set("create_pitr_mode", flattenedCreatePITRMode); err != nil {
		return fmt.Errorf("Error setting `create_pitr_mode`: %+v", err)
	}

	flattenedCreateRecoveryMode := flattenAzureRmMsSqlDatabaseCreateRecoveryMode(&resp)
	if err := d.Set("create_recovery_mode", flattenedCreateRecoveryMode); err != nil {
		return fmt.Errorf("Error setting `create_recovery_mode`: %+v", err)
	}

	flattenedCreateRestoreMode := flattenAzureRmMsSqlDatabaseCreateRestoreMode(&resp)
	if err := d.Set("create_restore_mode", flattenedCreateRestoreMode); err != nil {
		return fmt.Errorf("Error setting `create_restore_mode`: %+v", err)
	}

	flattenedCreateSecondaryMode := flattenAzureRmMsSqlDatabaseCreateSecondaryMode(&resp)
	if err := d.Set("create_secondary_mode", flattenedCreateSecondaryMode); err != nil {
		return fmt.Errorf("Error setting `create_secondary_mode`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMsSqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MsSqlDatabaseID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.MsSqlServer, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting MsSql Database %q ( MsSql Server %q / Resource Group %q): %+v", id.Name, id.MsSqlServer, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error waiting for MsSql Database %q ( MsSql Server %q / Resource Group %q) to be deleted: %+v", id.Name, id.MsSqlServer, id.ResourceGroup, err)
	}

	return nil
}

func expandAzureRmMsSqlDatabaseGP(d *schema.ResourceData, params *sql.Database) {
	gps := d.Get("general_purpose").([]interface{})
	if len(gps) == 0 {
		return
	}
	gp := gps[0].(map[string]interface{})

	params.Sku = &sql.Sku{
		Name: utils.String("GP" + "_" + gp["family"].(string) + "_" + strconv.Itoa(gp["capacity"].(int))),
	}

	if v, ok := gp["max_size_gb"]; ok {
		params.DatabaseProperties.MaxSizeBytes = utils.Int64(int64(v.(int) * 1073741824))
	}

	return
}

func flattenAzureRmMsSqlDatabaseGP(input *sql.Database) []interface{} {
	if *input.Sku.Tier != "GeneralPurpose" {
		return []interface{}{}
	}

	var capacity int32
	if input.Sku.Capacity != nil {
		capacity = *input.Sku.Capacity
	}

	var family string
	if input.Sku.Family != nil {
		family = *input.Sku.Family
	}

	var maxSizeGb int32
	if input.MaxSizeBytes != nil {
		maxSizeGb = int32((*input.MaxSizeBytes) / int64(1073741824))
	}

	return []interface{}{
		map[string]interface{}{
			"capacity":    capacity,
			"family":      family,
			"max_size_gb": maxSizeGb,
		},
	}
}

func expandAzureRmMsSqlDatabaseHS(d *schema.ResourceData, params *sql.Database) {
	hss := d.Get("hyper_scale").([]interface{})
	if len(hss) == 0 {
		return
	}
	hs := hss[0].(map[string]interface{})

	params.Sku = &sql.Sku{
		Name: utils.String("HS" + "_" + hs["family"].(string) + "_" + strconv.Itoa(hs["capacity"].(int))),
	}
	if readReplica, ok := hs["read_replica_count"]; ok {
		params.DatabaseProperties.ReadReplicaCount = utils.Int32(int32(readReplica.(int)))
	}

	return
}

func flattenAzureRmMsSqlDatabaseHS(input *sql.Database) []interface{} {
	if *input.Sku.Tier != "Hyperscale" {
		return []interface{}{}
	}

	var capacity, readReplica int32
	if input.Sku.Capacity != nil {
		capacity = *input.Sku.Capacity
	}

	var family string
	if input.Sku.Family != nil {
		family = *input.Sku.Family
	}

	if input.ReadReplicaCount != nil {
		readReplica = *input.ReadReplicaCount
	}

	return []interface{}{
		map[string]interface{}{
			"capacity":           capacity,
			"family":             family,
			"read_replica_count": readReplica,
		},
	}
}

func expandAzureRmMsSqlDatabaseBC(d *schema.ResourceData, params *sql.Database) {
	bcs := d.Get("business_critical").([]interface{})
	if len(bcs) == 0 {
		return
	}
	bc := bcs[0].(map[string]interface{})

	params.Sku = &sql.Sku{
		Name: utils.String("BC" + "_" + bc["family"].(string) + "_" + strconv.Itoa(bc["capacity"].(int))),
	}

	if v, ok := bc["max_size_gb"]; ok {
		params.DatabaseProperties.MaxSizeBytes = utils.Int64(int64(v.(int) * 1073741824))
	}

	if readScale, ok := bc["read_scale"]; ok {
		params.DatabaseProperties.ReadScale = sql.DatabaseReadScale(readScale.(string))
	}

	if zoneRedundant, ok := bc["zone_redundant"]; ok {
		params.DatabaseProperties.ZoneRedundant = utils.Bool(zoneRedundant.(bool))
	}

	return
}

func flattenAzureRmMsSqlDatabaseBC(input *sql.Database) []interface{} {
	if *input.Sku.Tier != "BusinessCritical" {
		return []interface{}{}
	}

	var capacity int32
	if input.Sku.Capacity != nil {
		capacity = *input.Sku.Capacity
	}

	var family, readScale string
	if input.Sku.Family != nil {
		family = *input.Sku.Family
	}

	var maxSizeGb int32
	if input.MaxSizeBytes != nil {
		maxSizeGb = int32((*input.MaxSizeBytes) / int64(1073741824))
	}

	if input.ReadScale != "" {
		readScale = string(input.ReadScale)
	}

	zoneRedundant := *input.ZoneRedundant

	return []interface{}{
		map[string]interface{}{
			"capacity":       capacity,
			"family":         family,
			"max_size_gb":    maxSizeGb,
			"read_scale":     readScale,
			"zone_redundant": zoneRedundant,
		},
	}
}

func expandAzureRmMsSqlDatabaseCreateCopyMode(d *schema.ResourceData, params *sql.Database) {
	copyModes := d.Get("create_copy_mode").([]interface{})
	if len(copyModes) == 0 {
		return
	}
	copyMode := copyModes[0].(map[string]interface{})

	params.DatabaseProperties.CreateMode = sql.CreateModeCopy
	params.DatabaseProperties.SourceDatabaseID = utils.String(copyMode["source_database_id"].(string))

	return
}

func flattenAzureRmMsSqlDatabaseCreateCopyMode(input *sql.Database) []interface{} {
	if input.CreateMode != sql.CreateModeCopy {
		return []interface{}{}
	}

	var sourceDBId string
	if input.SourceDatabaseID != nil {
		sourceDBId = *input.SourceDatabaseID
	}

	return []interface{}{
		map[string]interface{}{
			"source_database_id": sourceDBId,
		},
	}
}

func expandAzureRmMsSqlDatabaseCreatePITRMode(d *schema.ResourceData, params *sql.Database) {
	pitrModes := d.Get("create_pitr_mode").([]interface{})
	if len(pitrModes) == 0 {
		return
	}
	pitrMode := pitrModes[0].(map[string]interface{})

	params.DatabaseProperties.CreateMode = sql.CreateModePointInTimeRestore
	params.DatabaseProperties.SourceDatabaseID = utils.String(pitrMode["source_database_id"].(string))
	restorePointInTime, _ := time.Parse(time.RFC3339, pitrMode["restore_point_in_time"].(string))
	params.DatabaseProperties.RestorePointInTime = &date.Time{Time: restorePointInTime}

	return
}

func flattenAzureRmMsSqlDatabaseCreatePITRMode(input *sql.Database) []interface{} {
	if input.CreateMode != sql.CreateModePointInTimeRestore {
		return []interface{}{}
	}

	var sourceDBId, restorePointInTime string

	if input.SourceDatabaseID != nil {
		sourceDBId = *input.SourceDatabaseID
	}

	if input.RestorePointInTime != nil && !input.RestorePointInTime.IsZero() {
		restorePointInTime = input.RestorePointInTime.Format(time.RFC3339)
	}

	return []interface{}{
		map[string]interface{}{
			"source_database_id":    sourceDBId,
			"restore_point_in_time": restorePointInTime,
		},
	}
}

func expandAzureRmMsSqlDatabaseCreateRecoveryMode(d *schema.ResourceData, params *sql.Database) {
	recoveryModes := d.Get("create_recovery_mode").([]interface{})
	if len(recoveryModes) == 0 {
		return
	}

	recoveryMode := recoveryModes[0].(map[string]interface{})

	params.DatabaseProperties.CreateMode = sql.CreateModeRecovery
	params.DatabaseProperties.RestorableDroppedDatabaseID = utils.String(recoveryMode["restorable_dropped_database_id"].(string))

	return
}

func flattenAzureRmMsSqlDatabaseCreateRecoveryMode(input *sql.Database) []interface{} {
	if input.CreateMode != sql.CreateModeRecovery {
		return []interface{}{}
	}

	var restorableDroppedDatabaseID string
	if input.RestorableDroppedDatabaseID != nil {
		restorableDroppedDatabaseID = *input.RestorableDroppedDatabaseID
	}

	return []interface{}{
		map[string]interface{}{
			"restorable_dropped_database_id": restorableDroppedDatabaseID,
		},
	}
}

func expandAzureRmMsSqlDatabaseCreateRestoreMode(d *schema.ResourceData, params *sql.Database) {
	restoreModes := d.Get("create_restore_mode").([]interface{})
	if len(restoreModes) == 0 {
		return
	}
	restoreMode := restoreModes[0].(map[string]interface{})

	params.DatabaseProperties.CreateMode = sql.CreateModeRestore
	params.DatabaseProperties.SourceDatabaseID = utils.String(restoreMode["source_database_id"].(string))

	if v, ok := restoreMode["source_database_deletion_date"]; ok {
		sourceDatabaseDeletionDate, _ := time.Parse(time.RFC3339, v.(string))
		params.DatabaseProperties.SourceDatabaseDeletionDate = &date.Time{Time: sourceDatabaseDeletionDate}
	}

	return
}

func flattenAzureRmMsSqlDatabaseCreateRestoreMode(input *sql.Database) []interface{} {
	if input.CreateMode != sql.CreateModeRestore {
		return []interface{}{}
	}

	var sourceDatabaseID, sourceDatabaseDeletionDate string

	if input.SourceDatabaseID != nil {
		sourceDatabaseID = *input.SourceDatabaseID
	}

	if input.SourceDatabaseDeletionDate != nil && !input.SourceDatabaseDeletionDate.IsZero() {
		sourceDatabaseDeletionDate = input.SourceDatabaseDeletionDate.Format(time.RFC3339)
	}

	return []interface{}{
		map[string]interface{}{
			"source_database_id":            sourceDatabaseID,
			"source_database_deletion_date": sourceDatabaseDeletionDate,
		},
	}
}

func expandAzureRmMsSqlDatabaseCreateSecondaryMode(d *schema.ResourceData, params *sql.Database) {
	secondaryModes := d.Get("create_secondary_mode").([]interface{})
	if len(secondaryModes) == 0 {
		return
	}
	secondaryMode := secondaryModes[0].(map[string]interface{})

	params.DatabaseProperties.CreateMode = sql.CreateModeSecondary
	params.DatabaseProperties.SourceDatabaseID = utils.String(secondaryMode["source_database_id"].(string))

	return
}

func flattenAzureRmMsSqlDatabaseCreateSecondaryMode(input *sql.Database) []interface{} {
	if input.CreateMode != sql.CreateModeSecondary {
		return []interface{}{}
	}

	var sourceDatabaseID string
	if input.SourceDatabaseID != nil {
		sourceDatabaseID = *input.SourceDatabaseID
	}

	return []interface{}{
		map[string]interface{}{
			"source_database_id": sourceDatabaseID,
		},
	}
}
