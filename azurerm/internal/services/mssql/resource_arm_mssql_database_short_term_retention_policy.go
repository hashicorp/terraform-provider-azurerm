package mssql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMsSqlDatabaseShortTermRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlDatabaseShortTermRetentionPolicyCreateUpdate,
		Read:   resourceArmMsSqlDatabaseShortTermRetentionPolicyRead,
		Update: resourceArmMsSqlDatabaseShortTermRetentionPolicyCreateUpdate,
		Delete: resourceArmMsSqlDatabaseShortTermRetentionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			// RetentionDays - The backup retention period in days. This is how many days Point-in-Time Restore will be supported.
			"retention_days": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(7, 35),
			},
			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},
		},
	}
}

func resourceArmMsSqlDatabaseShortTermRetentionPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseName := d.Get("database_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	backupShortTermPolicyProps := sql.BackupShortTermRetentionPolicyProperties{
		RetentionDays: utils.Int32(7),
	}

	if v, ok := d.GetOk("retention_days"); ok {
		backupShortTermPolicyProps.RetentionDays = utils.Int32(int32(v.(int)))
	}

	backupShortTermPolicy := sql.BackupShortTermRetentionPolicy{
		BackupShortTermRetentionPolicyProperties: &backupShortTermPolicyProps,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, databaseName, backupShortTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Create/Update for Sql Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	response, err := client.Get(ctx, resourceGroup, serverName, databaseName)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Database %q Short Term Policies (Sql Server %q ,Resource Group %q): %+v", databaseName, serverName, resourceGroup, err)
	}
	d.SetId(*response.ID)

	return resourceArmMsSqlDatabaseShortTermRetentionPolicyRead(d, meta)
}

func resourceArmMsSqlDatabaseShortTermRetentionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	databaseName := id.Path["databases"]
	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]

	backupShortTermPolicy, err := client.Get(ctx, resourceGroup, serverName, databaseName)
	if err != nil {
		return fmt.Errorf("Error retrieving Short Term Policies for Database %q (Sql Server %q ;Resource Group %q): %+v", databaseName, serverName, resourceGroup, err)
	}

	if backupShortTermPolicy.RetentionDays != nil {
		retentionDays := *backupShortTermPolicy.RetentionDays
		if err := d.Set("retention_days", retentionDays); err != nil {
			return fmt.Errorf("Error setting `retention_days`: %+v", err)
		}
	}

	d.Set("database_name", databaseName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	return nil
}

// Default value for PITR is 7 days, therefore on delete we just set the defaults back
func resourceArmMsSqlDatabaseShortTermRetentionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	databaseName := id.Path["databases"]
	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]

	// Update to default values for removal
	backupShortTermPolicy := sql.BackupShortTermRetentionPolicy{
		BackupShortTermRetentionPolicyProperties: &sql.BackupShortTermRetentionPolicyProperties{
			RetentionDays: utils.Int32(7),
		},
	}

	future, err := client.Update(ctx, resourceGroup, serverName, databaseName, backupShortTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
