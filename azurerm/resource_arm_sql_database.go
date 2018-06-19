package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlDatabaseCreateUpdate,
		Read:   resourceArmSqlDatabaseRead,
		Update: resourceArmSqlDatabaseCreateUpdate,
		Delete: resourceArmSqlDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(sql.Default),
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.Copy),
					string(sql.Default),
					string(sql.NonReadableSecondary),
					string(sql.OnlineSecondary),
					string(sql.PointInTimeRestore),
					string(sql.Recovery),
					string(sql.Restore),
					string(sql.RestoreLongTermRetentionBackup),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"import": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_uri": {
							Type:     schema.TypeString,
							Required: true,
						},
						"storage_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"storage_key_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"StorageAccessKey",
								"SharedAccessKey",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"administrator_login": {
							Type:     schema.TypeString,
							Required: true,
						},
						"administrator_login_password": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"authentication_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADPassword",
								"SQL",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"operation_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Import",
							ValidateFunc: validation.StringInSlice([]string{
								"Import",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"source_database_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"restore_point_in_time": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRFC3339Date,
			},

			"edition": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.Basic),
					string(sql.Standard),
					string(sql.Premium),
					string(sql.DataWarehouse),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"max_size_bytes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"requested_service_objective_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateUUID,
			},

			"requested_service_objective_name": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				// TODO: add validation once the Enum's complete
				// https://github.com/Azure/azure-rest-api-specs/issues/1609
			},

			"source_database_deletion_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRFC3339Date,
			},

			"elastic_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSqlDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlDatabasesClient

	name := d.Get("name").(string)
	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	location := azureRMNormalizeLocation(d.Get("location").(string))
	createMode := d.Get("create_mode").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties := sql.Database{
		Location: utils.String(location),
		DatabaseProperties: &sql.DatabaseProperties{
			CreateMode: sql.CreateMode(createMode),
		},
		Tags: expandTags(tags),
	}

	if v, ok := d.GetOk("source_database_id"); ok {
		sourceDatabaseID := v.(string)
		properties.DatabaseProperties.SourceDatabaseID = utils.String(sourceDatabaseID)
	}

	if v, ok := d.GetOk("edition"); ok {
		edition := v.(string)
		properties.DatabaseProperties.Edition = sql.DatabaseEdition(edition)
	}

	if v, ok := d.GetOk("collation"); ok {
		collation := v.(string)
		properties.DatabaseProperties.Collation = utils.String(collation)
	}

	if v, ok := d.GetOk("max_size_bytes"); ok {
		maxSizeBytes := v.(string)
		properties.DatabaseProperties.MaxSizeBytes = utils.String(maxSizeBytes)
	}

	if v, ok := d.GetOk("source_database_deletion_date"); ok {
		sourceDatabaseDeletionString := v.(string)
		sourceDatabaseDeletionDate, err := date.ParseTime(time.RFC3339, sourceDatabaseDeletionString)
		if err != nil {
			return fmt.Errorf("`source_database_deletion_date` wasn't a valid RFC3339 date %q: %+v", sourceDatabaseDeletionString, err)
		}

		properties.DatabaseProperties.SourceDatabaseDeletionDate = &date.Time{
			Time: sourceDatabaseDeletionDate,
		}
	}

	if v, ok := d.GetOk("requested_service_objective_id"); ok {
		requestedServiceObjectiveID := v.(string)
		id, err := uuid.FromString(requestedServiceObjectiveID)
		if err != nil {
			return fmt.Errorf("`requested_service_objective_id` wasn't a valid UUID %q: %+v", requestedServiceObjectiveID, err)
		}
		properties.DatabaseProperties.RequestedServiceObjectiveID = &id
	}

	if v, ok := d.GetOk("elastic_pool_name"); ok {
		elasticPoolName := v.(string)
		properties.DatabaseProperties.ElasticPoolName = utils.String(elasticPoolName)
	}

	if v, ok := d.GetOk("requested_service_objective_name"); ok {
		requestedServiceObjectiveName := v.(string)
		properties.DatabaseProperties.RequestedServiceObjectiveName = sql.ServiceObjectiveName(requestedServiceObjectiveName)
	}

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		restorePointInTime := v.(string)
		restorePointInTimeDate, err := date.ParseTime(time.RFC3339, restorePointInTime)
		if err != nil {
			return fmt.Errorf("`restore_point_in_time` wasn't a valid RFC3339 date %q: %+v", restorePointInTime, err)
		}

		properties.DatabaseProperties.RestorePointInTime = &date.Time{
			Time: restorePointInTimeDate,
		}
	}

	ctx := meta.(*ArmClient).StopContext
	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, name, properties)
	if err != nil {
		return err
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return err
	}

	if _, ok := d.GetOk("import"); ok {
		if !strings.EqualFold(createMode, "default") {
			return fmt.Errorf("import can only be used when create_mode is Default")
		}
		importParameters := expandAzureRmSqlDatabaseImport(d)
		importFuture, err := client.CreateImportOperation(ctx, resourceGroup, serverName, name, importParameters)
		if err != nil {
			return err
		}

		// this is set in config.go, but something sets
		// it back to 15 minutes, which isn't long enough
		// for most imports
		client.Client.PollingDuration = 60 * time.Minute

		err = importFuture.WaitForCompletion(ctx, client.Client)
		if err != nil {
			return err
		}
	}

	resp, err := client.Get(ctx, resourceGroup, serverName, name, "")
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmSqlDatabaseRead(d, meta)
}

func resourceArmSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlDatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["databases"]

	resp, err := client.Get(ctx, resourceGroup, serverName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Database %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Sql Database %s: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("server_name", serverName)

	if props := resp.DatabaseProperties; props != nil {
		// TODO: set `create_mode` & `source_database_id` once this issue is fixed:
		// https://github.com/Azure/azure-rest-api-specs/issues/1604

		d.Set("collation", props.Collation)
		d.Set("default_secondary_location", props.DefaultSecondaryLocation)
		d.Set("edition", string(props.Edition))
		d.Set("elastic_pool_name", props.ElasticPoolName)
		d.Set("max_size_bytes", props.MaxSizeBytes)
		d.Set("requested_service_objective_name", string(props.RequestedServiceObjectiveName))

		if cd := props.CreationDate; cd != nil {
			d.Set("creation_date", cd.String())
		}

		if rsoid := props.RequestedServiceObjectiveID; rsoid != nil {
			d.Set("requested_service_objective_id", rsoid.String())
		}

		if rpit := props.RestorePointInTime; rpit != nil {
			d.Set("restore_point_in_time", rpit.String())
		}

		if sddd := props.SourceDatabaseDeletionDate; sddd != nil {
			d.Set("source_database_deletion_date", sddd.String())
		}

		d.Set("encryption", flattenEncryptionStatus(props.TransparentDataEncryption))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).sqlDatabasesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	name := id.Path["databases"]

	resp, err := client.Delete(ctx, resourceGroup, serverName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error making Read request on Sql Database %s: %+v", name, err)
	}

	if err != nil {
		return fmt.Errorf("Error deleting SQL Database: %+v", err)
	}

	return nil
}

func flattenEncryptionStatus(encryption *[]sql.TransparentDataEncryption) string {
	if encryption != nil {
		encrypted := *encryption
		if len(encrypted) > 0 {
			if props := encrypted[0].TransparentDataEncryptionProperties; props != nil {
				return string(props.Status)
			}
		}
	}

	return ""
}

func expandAzureRmSqlDatabaseImport(d *schema.ResourceData) sql.ImportExtensionRequest {
	v := d.Get("import")
	dbimportRefs := v.([]interface{})
	dbimportRef := dbimportRefs[0].(map[string]interface{})
	return sql.ImportExtensionRequest{
		Name: utils.String("terraform"),
		ImportExtensionProperties: &sql.ImportExtensionProperties{
			StorageKeyType:             sql.StorageKeyType(dbimportRef["storage_key_type"].(string)),
			StorageKey:                 utils.String(dbimportRef["storage_key"].(string)),
			StorageURI:                 utils.String(dbimportRef["storage_uri"].(string)),
			AdministratorLogin:         utils.String(dbimportRef["administrator_login"].(string)),
			AdministratorLoginPassword: utils.String(dbimportRef["administrator_login_password"].(string)),
			AuthenticationType:         sql.AuthenticationType(dbimportRef["authentication_type"].(string)),
			OperationMode:              utils.String(dbimportRef["operation_mode"].(string)),
		},
	}
}
