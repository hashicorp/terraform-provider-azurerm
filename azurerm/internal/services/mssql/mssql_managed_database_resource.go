package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMSSQLManagedDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSQLManagedDatabaseCreateUpdate,
		Read:   resourceArmMSSQLManagedDatabaseRead,
		Update: resourceArmMSSQLManagedDatabaseCreateUpdate,
		Delete: resourceArmMSSQLManagedDatabaseDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"managed_instance_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "SQL_Latin1_General_CP1_CI_AS",
			},

			"restore_point_in_time": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
				ForceNew:         true,
			},

			"catalog_collation": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.DATABASEDEFAULT),
					string(sql.SQLLatin1GeneralCP1CIAS),
				}, false),
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ManagedDatabaseCreateModeDefault),
					string(sql.ManagedDatabaseCreateModePointInTimeRestore),
					string(sql.ManagedDatabaseCreateModeRecovery),
					string(sql.ManagedDatabaseCreateModeRestoreExternalBackup),
					string(sql.ManagedDatabaseCreateModeRestoreLongTermRetentionBackup),
				}, false),
			},

			"storage_container_uri": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"storage_container_uri", "storage_container_sas_token", "last_backup_name"},
			},

			"source_database_id": {
				Type:             schema.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"restorable_dropped_database_id": {
				Type:             schema.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"storage_container_sas_token": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"storage_container_uri", "storage_container_sas_token", "last_backup_name"},
			},

			"recoverable_database_id": {
				Type:             schema.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"last_backup_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"storage_container_uri", "storage_container_sas_token", "last_backup_name"},
			},

			"longterm_retention_backup_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"auto_complete_restore": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": tags.Schema(),

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"earliest_restore_point": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"failover_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmMSSQLManagedDatabaseCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedDatabasesClient
	managedInstanceClient := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	managedInstanceName := d.Get("managed_instance_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	instanceResponse, err := managedInstanceClient.Get(ctx, resourceGroup, managedInstanceName)
	if err != nil {
		return fmt.Errorf("making Read request on managed SQL instance %q (Resource Group %q): %s", managedInstanceName, resourceGroup, err)
	}

	location := *instanceResponse.Location

	if location == "" {
		return fmt.Errorf("Location is empty from making Read request on managed instance %q", managedInstanceName)
	}

	createMode := d.Get("create_mode").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, managedInstanceName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing managed database %q in managed instance instance %q (Resource Group %q): %+v", name, managedInstanceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_managed_database", *existing.ID)
		}
	}

	parameters := sql.ManagedDatabase{
		Location:                  utils.String(location),
		Tags:                      tags.Expand(t),
		ManagedDatabaseProperties: &sql.ManagedDatabaseProperties{},
	}

	if d.HasChange("create_mode") {
		if createMode == string(sql.ManagedDatabaseCreateModePointInTimeRestore) {
			_, sourceDatabaseExists := d.GetOk("source_database_id")
			_, recoverableDroppedDatabaseExists := d.GetOk("restorable_dropped_database_id")
			_, restorePointExists := d.GetOk("restore_point_in_time")
			if (!sourceDatabaseExists && !recoverableDroppedDatabaseExists) || !restorePointExists {
				return fmt.Errorf("could not create managed database %q in managed instance %q (Resource Group %q) in restore in point create mode. Restore point in time and either of source database id or restorable dropped database id values should be supplied.", name, managedInstanceName, resourceGroup)
			}
		}

		if createMode == string(sql.ManagedDatabaseCreateModeRestoreExternalBackup) {
			_, collationExists := d.GetOk("collation")
			_, storageContainerUriExists := d.GetOk("storage_container_uri")
			_, storageContainerSasExists := d.GetOk("storage_container_sas_token")
			_, lastBackupExists := d.GetOk("last_backup_name")

			if !collationExists || !storageContainerUriExists || !storageContainerSasExists || !lastBackupExists {
				return fmt.Errorf("could not create managed database %q in managed instance %q (Resource Group %q) in restore from external backup mode. storage_container_uri, storage_container_sas_token and last backup name values should be supplied.", name, managedInstanceName, resourceGroup)
			}
		}

		if createMode == string(sql.ManagedDatabaseCreateModeRecovery) {
			if _, recoverableDatabaseExists := d.GetOk("recoverable_database_id"); !recoverableDatabaseExists {
				return fmt.Errorf("could not create managed database %q in managed instance %q (Resource Group %q) in recovery mode. recoverable_database_id value should be supplied.", name, managedInstanceName, resourceGroup)
			}
		}

		if createMode == string(sql.ManagedDatabaseCreateModeRestoreLongTermRetentionBackup) {
			if _, longtimeRetentionBackupExists := d.GetOk("longterm_retention_backup_id"); !longtimeRetentionBackupExists {
				return fmt.Errorf("could not create managed database %q in managed instance %q (Resource Group %q) in long term retention backup mode. long_term_retention_backup_id value should be supplied.", name, managedInstanceName, resourceGroup)
			}
		}

		parameters.ManagedDatabaseProperties.CreateMode = sql.ManagedDatabaseCreateMode(createMode)
	}

	if v, exists := d.GetOk("collation"); exists {
		collation := v.(string)
		parameters.ManagedDatabaseProperties.Collation = utils.String(collation)
	}

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		if createMode != string(sql.ManagedDatabaseCreateModePointInTimeRestore) {
			return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModePointInTimeRestore))
		}
		restorePointInTime := v.(string)
		restorePointInTimeDate, err2 := date.ParseTime(time.RFC3339, restorePointInTime)
		if err2 != nil {
			return fmt.Errorf("`restore_point_in_time` wasn't a valid RFC3339 date %q: %+v", restorePointInTime, err2)
		}
		parameters.ManagedDatabaseProperties.RestorePointInTime = &date.Time{
			Time: restorePointInTimeDate,
		}
	}

	if v, exists := d.GetOk("catalog_collation"); exists {
		catalogCollation := v.(string)
		parameters.ManagedDatabaseProperties.CatalogCollation = sql.CatalogCollationType(catalogCollation)
	}

	if v, exists := d.GetOk("storage_container_uri"); exists {
		if createMode != string(sql.ManagedDatabaseCreateModeRestoreExternalBackup) {
			return fmt.Errorf("'storage_container_uri' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModeRestoreExternalBackup))
		}
		storageContainerUri := v.(string)
		parameters.ManagedDatabaseProperties.StorageContainerURI = utils.String(storageContainerUri)
	}

	if v, exists := d.GetOk("source_database_id"); exists {
		if createMode != string(sql.ManagedDatabaseCreateModePointInTimeRestore) {
			return fmt.Errorf("'source_database_id' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModePointInTimeRestore))
		}
		sourceDatabaseId := v.(string)
		parameters.ManagedDatabaseProperties.SourceDatabaseID = utils.String(sourceDatabaseId)
	}

	if v, exists := d.GetOk("restorable_dropped_database_id"); exists {
		restorableDroppedDatabaseId := v.(string)
		parameters.ManagedDatabaseProperties.RestorableDroppedDatabaseID = utils.String(restorableDroppedDatabaseId)
	}

	if v, exists := d.GetOk("last_backup_name"); exists {
		if createMode != string(sql.ManagedDatabaseCreateModeRestoreExternalBackup) {
			return fmt.Errorf("'last_backup_name' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModeRestoreExternalBackup))
		}
		parameters.ManagedDatabaseProperties.LastBackupName = utils.String(v.(string))
	}

	if v, exists := d.GetOk("auto_complete_restore"); exists {
		parameters.ManagedDatabaseProperties.AutoCompleteRestore = utils.Bool(v.(bool))
	}

	if v, exists := d.GetOk("storage_container_sas_token"); exists {
		if createMode != string(sql.ManagedDatabaseCreateModeRestoreExternalBackup) {
			return fmt.Errorf("'storage_container_sas_token' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModeRestoreExternalBackup))
		}
		storageAccountSasToken := v.(string)
		parameters.ManagedDatabaseProperties.StorageContainerSasToken = utils.String(storageAccountSasToken)
	}

	if v, exists := d.GetOk("recoverable_database_id"); exists {
		if createMode != string(sql.ManagedDatabaseCreateModeRecovery) {
			return fmt.Errorf("'recoverable_database_id' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModeRecovery))
		}
		recoverableDatabaseId := v.(string)
		parameters.ManagedDatabaseProperties.RecoverableDatabaseID = utils.String(recoverableDatabaseId)
	}

	if v, exists := d.GetOk("longterm_retention_backup_id"); exists {
		if createMode != string(sql.ManagedDatabaseCreateModeRestoreLongTermRetentionBackup) {
			return fmt.Errorf("'longterm_retention_backup_id' is supported only for create_mode %s", string(sql.ManagedDatabaseCreateModeRestoreLongTermRetentionBackup))
		}
		longtermRetentionBackupId := v.(string)
		parameters.ManagedDatabaseProperties.LongTermRetentionBackupResourceID = utils.String(longtermRetentionBackupId)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, managedInstanceName, name, parameters)
	if err != nil {
		return fmt.Errorf("while making create/update request for Managed Database %q (Managed instance %q, Resource group: %q): %+v", name, managedInstanceName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("while waiting on create/update request for Managed Database %q (Managed instance %q, Resource group: %q): %+v", name, managedInstanceName, resourceGroup, err)
	}

	time.Sleep(20 * time.Second)
	result, err := client.Get(ctx, resourceGroup, managedInstanceName, name)
	if err != nil {
		return fmt.Errorf("while making get request for Managed Database %q (Managed instance %q, Resource group: %q): %+v", name, managedInstanceName, resourceGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("while getting ID from Managed Database %q (Managed instance %q, Resource group: %q): %+v", name, managedInstanceName, resourceGroup, err)
	}

	d.SetId(*result.ID)

	return resourceArmMSSQLManagedDatabaseRead(d, meta)
}

func resourceArmMSSQLManagedDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedDatabasesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["databases"]
	managedInstanceName := id.Path["managedInstances"]

	resp, err := client.Get(ctx, resourceGroup, managedInstanceName, name)
	if err != nil {
		return fmt.Errorf("Error reading managed SQL Database %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("type", resp.Type)
	d.Set("managed_instance_name", managedInstanceName)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ManagedDatabaseProperties; props != nil {
		d.Set("collation", props.Collation)
		d.Set("status", props.Status)
		if props.CreationDate != nil && props.CreationDate.String() != "" {
			d.Set("creation_date", props.CreationDate.String())
		}
		if props.EarliestRestorePoint != nil && props.EarliestRestorePoint.String() != "" {
			d.Set("earliest_restore_point", props.EarliestRestorePoint.String())
		}
		if props.RestorePointInTime != nil && props.RestorePointInTime.String() != "" {
			d.Set("restore_point_in_time", props.RestorePointInTime.String())
		}
		d.Set("default_secondary_location", props.DefaultSecondaryLocation)
		d.Set("storage_container_uri", props.StorageContainerURI)
		d.Set("failover_group_id", props.FailoverGroupID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMSSQLManagedDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedDatabasesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	name := id.Path["databases"]
	managedInstanceName := id.Path["managedInstances"]

	future, err := client.Delete(ctx, resourceGroup, managedInstanceName, name)
	if err != nil {
		return fmt.Errorf("Error deleting managed SQL database %s: %+v", name, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
