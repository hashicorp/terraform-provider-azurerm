package sql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/sql/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlDatabaseShortTermRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlDatabaseShortTermRetentionPolicyCreateUpdate,
		Read:   resourceArmSqlDatabaseShortTermRetentionPolicyRead,
		Update: resourceArmSqlDatabaseShortTermRetentionPolicyCreateUpdate,
		Delete: resourceArmSqlDatabaseShortTermRetentionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"backup_short_term_retention_policy": helper.SqlShortTermRetentionPolicy(),
			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},
		},
	}
}

func resourceArmSqlDatabaseShortTermRetentionPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseName := d.Get("database_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	shortTermPolicy := d.Get("backup_short_term_retention_policy").([]interface{})

	backupShortTermPolicy := sql.BackupShortTermRetentionPolicy{
		BackupShortTermRetentionPolicyProperties: helper.ExpandSqlShortTermRetentionPolicyProperties(shortTermPolicy),
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

	return resourceArmSqlDatabaseShortTermRetentionPolicyRead(d, meta)
}

func resourceArmSqlDatabaseShortTermRetentionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.BackupShortTermRetentionPoliciesClient
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

	flattenedShortTermPolicy := helper.FlattenSqlShortTermRetentionPolicy(&backupShortTermPolicy)
	if err := d.Set("backup_short_term_retention_policy", flattenedShortTermPolicy); err != nil {
		return fmt.Errorf("Error setting `backup_short_term_retention_policy`: %+v", err)
	}

	return nil
}

func resourceArmSqlDatabaseShortTermRetentionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sql.BackupShortTermRetentionPoliciesClient
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
			RetentionDays: utils.Int32(1),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, databaseName, backupShortTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
