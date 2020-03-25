package mssql

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mssql/validate"
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

			//recovery is not support in version 2017-10-01-preview
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
				ForceNew:     true,
				ValidateFunc: validate.MsSqlElasticPoolID,
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
				ValidateFunc: azValidate.FloatInSlice([]float64{0.5, 0.75, 1, 1.25, 1.5, 1.75, 2}),
			},

			"restore_point_in_time": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
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

			// hyper_scale can not be changed into other sku
			"sku_name": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMsSqlDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.DatabasesClient
	serverClient := meta.(*clients.Client).MSSQL.ServersClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for MsSql Database creation.")

	name := d.Get("name").(string)
	sqlServerId := d.Get("server_id").(string)
	serverId, _ := parse.MsSqlServerID(sqlServerId)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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
		return fmt.Errorf("Failure in making Read request on MsSql Server %q (Resource Group %q): %s", serverId.Name, serverId.ResourceGroup, err)
	}

	location := *serverResp.Location
	if location == "" {
		return fmt.Errorf("Location is empty from making Read request on MsSql Server %q", serverId.Name)
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

	if v, ok := d.GetOk("create_mode"); ok {
		if _, ok := d.GetOk("creation_source_database_id"); (v.(string) == string(sql.CreateModeCopy) || v.(string) == string(sql.CreateModePointInTimeRestore) || v.(string) == string(sql.CreateModeRestore) || v.(string) == string(sql.CreateModeSecondary)) && !ok {
			return fmt.Errorf("'creation_source_database_id' is required for create_mode %s", v.(string))
		}
		params.DatabaseProperties.CreateMode = sql.CreateMode(v.(string))
	}

	if v, ok := d.GetOk("max_size_gb"); ok {
		params.DatabaseProperties.MaxSizeBytes = utils.Int64(int64(v.(int) * 1073741824))
	}

	if v, ok := d.GetOkExists("read_scale"); ok {
		if v.(bool) {
			params.DatabaseProperties.ReadScale = sql.DatabaseReadScaleEnabled
		} else {
			params.DatabaseProperties.ReadScale = sql.DatabaseReadScaleDisabled
		}
	}

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		if cm, ok := d.GetOk("create_mode"); ok && cm.(string) != string(sql.CreateModePointInTimeRestore) {
			return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(sql.CreateModePointInTimeRestore))
		}
		restorePointInTime, _ := time.Parse(time.RFC3339, v.(string))
		params.DatabaseProperties.RestorePointInTime = &date.Time{Time: restorePointInTime}
	}

	if v, ok := d.GetOk("sku_name"); ok {
		params.Sku = &sql.Sku{
			Name: utils.String(v.(string)),
		}
	}

	if v, ok := d.GetOk("creation_source_database_id"); ok {
		params.DatabaseProperties.SourceDatabaseID = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, serverId.ResourceGroup, serverId.Name, name, params)
	if err != nil {
		return fmt.Errorf("Failure in creating MsSql Database %q (Sql Server %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Failure in waiting for creation of MsSql Database %q (MsSql Server Name %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
	}

	read, err := client.Get(ctx, serverId.ResourceGroup, serverId.Name, name)
	if err != nil {
		return fmt.Errorf("Failure in retrieving MsSql Database %q (MsSql Server Name %q / Resource Group %q): %+v", name, serverId.Name, serverId.ResourceGroup, err)
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
		return fmt.Errorf("Failure in reading MsSql Database %s (MsSql Server Name %q / Resource Group %q): %s", databaseId.Name, databaseId.MsSqlServer, databaseId.ResourceGroup, err)
	}

	d.Set("name", resp.Name)

	serverClient := meta.(*clients.Client).MSSQL.ServersClient

	serverResp, err := serverClient.Get(ctx, databaseId.ResourceGroup, databaseId.MsSqlServer)
	if err != nil || *serverResp.ID == "" {
		return fmt.Errorf("Failure in making Read request on MsSql Server  %q (Resource Group %q): %s", databaseId.MsSqlServer, databaseId.ResourceGroup, err)
	}
	d.Set("server_id", serverResp.ID)

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
		d.Set("sku_name", props.CurrentServiceObjectiveName)
		d.Set("zone_redundant", props.ZoneRedundant)
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
		return fmt.Errorf("Failure in deleting MsSql Database %q ( MsSql Server %q / Resource Group %q): %+v", id.Name, id.MsSqlServer, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Failure in waiting for MsSql Database %q ( MsSql Server %q / Resource Group %q) to be deleted: %+v", id.Name, id.MsSqlServer, id.ResourceGroup, err)
	}

	return nil
}
