package mssql

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
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
				ValidateFunc: ValidateMsSqlServerID,
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
							}, false),
						},

						"max_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 1024),
						},

						"read_scale": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(sql.DatabaseReadScaleEnabled),
								string(sql.DatabaseReadScaleDisabled),
							}, false),
						},

						"zone_redundant": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
				ConflictsWith: []string{"elastic_pool_id", "general_purpose", "hyper_scale"},
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9_]+$`),
					`This collation is not valid.`,
				),
			},

			// source_database_id will not be returned
			"create_copy_mode": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateMsSqlDatabaseID,
						},
					},
				},
				ConflictsWith: []string{"create_pitr_mode", "create_secondary_mode"},
			},

			//source_database_id and restore_point_in_time will not be returned
			"create_pitr_mode": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateMsSqlDatabaseID,
						},
						"restore_point_in_time": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.RFC3339Time,
							ValidateFunc:     validation.IsRFC3339Time,
						},
					},
				},
				ConflictsWith: []string{"create_copy_mode", "create_secondary_mode"},
			},

			//source_database_id will not be returned
			"create_secondary_mode": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_database_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidateMsSqlDatabaseID,
						},
					},
				},
				ConflictsWith: []string{"create_copy_mode", "create_pitr_mode"},
			},

			"elastic_pool_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  ValidateMsSqlElasticPoolID,
				ConflictsWith: []string{"business_critical", "general_purpose", "hyper_scale"},
			},

			// By default, the db sku is general_purpose.
			"general_purpose": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
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
							}, true),
						},

						"max_size_gb": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 1024),
						},

						"serverless": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// validation in 1 hour and 7 days or -1
									"auto_pause_delay_in_minutes": {
										Type:         schema.TypeInt,
										Optional:     true,
										Computed:     true,
										ValidateFunc: ValidateMsSqlDatabaseAutoPauseDelay,
									},
									// validation in float slice {0.5,0.75,1,1.25,1.5,1.75,2}
									"min_capacity": {
										Type:         schema.TypeFloat,
										Optional:     true,
										Computed:     true,
										ValidateFunc: ValidateMsSqlDBMinCapacity,
									},
								},
							},
						},
					},
				},
				ConflictsWith: []string{"business_critical", "elastic_pool_id", "hyper_scale"},
			},

			// hyper_scale could not be changed to the other skus
			"hyper_scale": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
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
				ConflictsWith: []string{"business_critical", "elastic_pool_id", "general_purpose"},
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.BasePrice),
					string(sql.LicenseIncluded),
				}, false),
			},

			//sample_name doesn't return back
			"sample_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.AdventureWorksLT),
				}, false),
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
			return tf.ImportAsExistsError("azurerm_mssql_database", *existing.ID)
		}
	}

	serverClient := meta.(*clients.Client).MSSQL.ServersClient
	serverResp, err := serverClient.Get(ctx, serverId.ResourceGroup, serverId.Name)
	if err != nil {
		return fmt.Errorf("Error making Read request on MsSql Server  %q (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
	}

	location := *serverResp.Location
	if location == "" {
		return fmt.Errorf("Error location is empty from making Read request on MsSql Server %q", serverId.Name)
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
		return fmt.Errorf("Error creating MsSql Database %q (Sql Server %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of MsSql Database %q (MsSql Server Name %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MsSql Database %q (MsSql Server Name %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	if read.ID == nil || *read.ID == "" {
		return fmt.Errorf("Cannot read MsSql Database %q (MsSql Server Name %q / Resource Group %q) ID", name, serverId.Name, serverId.ResourceGroup)
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
		return fmt.Errorf("Error reading MsSql Database %s (MsSql Server Name %q / Resource Group %q): %s", databaseId.Name, databaseId.MsSqlServer, databaseId.ResourceGroup, err)
	}

	d.Set("name", resp.Name)

	serverClient := meta.(*clients.Client).MSSQL.ServersClient

	serverResp, err := serverClient.Get(ctx, databaseId.ResourceGroup, databaseId.MsSqlServer)
	if err != nil || *serverResp.ID == "" {
		return fmt.Errorf("Error making Read request on MsSql Server  %q (Resource Group %q): %s", databaseId.MsSqlServer, databaseId.ResourceGroup, err)
	}
	d.Set("mssql_server_id", serverResp.ID)

	flattenedBC := flattenAzureRmMsSqlDatabaseBC(&resp)
	if err := d.Set("business_critical", flattenedBC); err != nil {
		return fmt.Errorf("Error setting `business_critical`: %+v", err)
	}

	d.Set("collation", resp.Collation)

	d.Set("elastic_pool_id", resp.ElasticPoolID)

	flattenedGP := flattenAzureRmMsSqlDatabaseGP(&resp)
	if err := d.Set("general_purpose", flattenedGP); err != nil {
		return fmt.Errorf("Error setting `general_purpose`: %+v", err)
	}

	flattenedHS := flattenAzureRmMsSqlDatabaseHS(&resp)
	if err := d.Set("hyper_scale", flattenedHS); err != nil {
		return fmt.Errorf("Error setting `hyper_scale`: %+v", err)
	}

	d.Set("license_type", resp.LicenseType)

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

	expandAzureRmMsSqlDatabaseGPServerless(d, params)
}

func flattenAzureRmMsSqlDatabaseGP(input *sql.Database) []interface{} {
	if input.Sku == nil || input.Sku.Tier == nil || *input.Sku.Tier != "GeneralPurpose" {
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

	serverless := flattenAzureRmMsSqlDatabaseGPServerless(input)

	return []interface{}{
		map[string]interface{}{
			"capacity":    capacity,
			"family":      family,
			"max_size_gb": maxSizeGb,
			"serverless":  serverless,
		},
	}
}

func expandAzureRmMsSqlDatabaseGPServerless(d *schema.ResourceData, params *sql.Database) {
	serverlesslist := d.Get("general_purpose.0.serverless").([]interface{})
	if len(serverlesslist) == 0 {
		return
	}

	// if sku is general purpose serverless, the sku_name is "GP_S_family_capacity"
	skuNameList := strings.Split(*params.Sku.Name, "_")
	skuNameList = append(skuNameList, "0")
	copy(skuNameList[2:], skuNameList[1:])
	skuNameList[1] = "S"
	params.Sku.Name = utils.String(strings.Join(skuNameList, "_"))

	if serverlesslist[0] == nil {
		return
	}
	serverless := serverlesslist[0].(map[string]interface{})

	if v, ok := serverless["auto_pause_delay_in_minutes"]; ok {
		params.DatabaseProperties.AutoPauseDelay = utils.Int32(int32(v.(int)))
	}

	if v, ok := serverless["min_capacity"]; ok {
		params.DatabaseProperties.MinCapacity = utils.Float(v.(float64))
	}
}

func flattenAzureRmMsSqlDatabaseGPServerless(input *sql.Database) []interface{} {
	if input.Sku == nil || input.Sku.Name == nil || !strings.HasPrefix(*input.Sku.Name, "GP_S") {
		return []interface{}{}
	}

	var autoPauseDelay int32
	if input.DatabaseProperties.AutoPauseDelay != nil {
		autoPauseDelay = *input.DatabaseProperties.AutoPauseDelay
	}

	var minCapacity float64
	if input.DatabaseProperties.MinCapacity != nil {
		minCapacity = *input.DatabaseProperties.MinCapacity
	}

	return []interface{}{
		map[string]interface{}{
			"auto_pause_delay_in_minutes": autoPauseDelay,
			"min_capacity":                minCapacity,
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
}

func flattenAzureRmMsSqlDatabaseHS(input *sql.Database) []interface{} {
	if input.Sku == nil || input.Sku.Tier == nil || *input.Sku.Tier != "Hyperscale" {
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
}

func flattenAzureRmMsSqlDatabaseBC(input *sql.Database) []interface{} {
	if input.Sku == nil || input.Sku.Tier == nil || *input.Sku.Tier != "BusinessCritical" {
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
}

func expandAzureRmMsSqlDatabaseCreateSecondaryMode(d *schema.ResourceData, params *sql.Database) {
	secondaryModes := d.Get("create_secondary_mode").([]interface{})
	if len(secondaryModes) == 0 {
		return
	}
	secondaryMode := secondaryModes[0].(map[string]interface{})

	params.DatabaseProperties.CreateMode = sql.CreateModeSecondary
	params.DatabaseProperties.SourceDatabaseID = utils.String(secondaryMode["source_database_id"].(string))
}

func ValidateMsSqlServerID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.MsSqlServerID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Server resource id: %v", k, err))
	}

	return warnings, errors
}

func ValidateMsSqlElasticPoolID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, _, _, err := parseArmMSSqlElasticPoolId(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Elastic Pool resource id: %v", k, err))
	}

	return warnings, errors
}

func ValidateMsSqlDatabaseID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parse.MsSqlDatabaseID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a MsSql Database resource id: %v", k, err))
	}

	return warnings, errors
}

func ValidateMsSqlDatabaseAutoPauseDelay(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be integer", k))
		return warnings, errors
	}
	min := 60
	max := 10080
	if (v < min || v > max) && v != -1 {
		errors = append(errors, fmt.Errorf("expected %s to be in the range (%d - %d) or -1, got %d", k, min, max, v))
		return warnings, errors
	}

	return warnings, errors
}

func ValidateMsSqlDBMinCapacity(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be float", k))
		return warnings, errors
	}

	valid := []float64{0.5, 0.75, 1, 1.25, 1.5, 1.75, 2}

	for _, validValue := range valid {
		if v == validValue {
			return warnings, errors
		}
	}

	errors = append(errors, fmt.Errorf("expected %s to be one of %v, got %f", k, valid, v))
	return warnings, errors
}
