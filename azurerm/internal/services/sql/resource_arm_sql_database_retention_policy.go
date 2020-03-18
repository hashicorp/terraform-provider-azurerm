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
)

func resourceArmSQLServerLongTermRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSQLServerDatabaseRetentionPolicyCreateUpdate,
		Read:   resourceArmSQLServerDatabaseRetentionPolicyRead,
		Update: resourceArmSQLServerDatabaseRetentionPolicyCreateUpdate,
		Delete: resourceArmSQLServerDatabaseRetentionPolicyDelete,

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
			"backup_long_term_retention_policy":  helper.SQLLongTermRententionPolicy(),
			"backup_short_term_retention_policy": helper.SQLShortTermRetentionPolicy(),
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

func resourceArmSQLServerDatabaseRetentionPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	longTermPolicyClient := meta.(*clients.Client).Sql.BackupLongTermRetentionPoliciesClient
	shortTermPolicyClient := meta.(*clients.Client).Sql.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseName := d.Get("database_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)
	longTermPolicy := d.Get("backup_long_term_retention_policy").([]interface{})
	shortTermPolicy := d.Get("backup_short_term_retention_policy").([]interface{})

	backupLongTermPolicy := sql.BackupLongTermRetentionPolicy{
		LongTermRetentionPolicyProperties: helper.ExpandSQLLongTermRetentionPolicyProperties(longTermPolicy),
	}

	longTermFuture, err := longTermPolicyClient.CreateOrUpdate(ctx, resourceGroup, serverName, databaseName, backupLongTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	if err = longTermFuture.WaitForCompletionRef(ctx, longTermPolicyClient.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Create/Update for SQL Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	// Set Long term policy id

	backupShortTermPolicy := sql.BackupShortTermRetentionPolicy{
		BackupShortTermRetentionPolicyProperties: helper.ExpandSQLShortTermRetentionPolicyProperties(shortTermPolicy),
	}

	shortTermFuture, err := shortTermPolicyClient.CreateOrUpdate(ctx, resourceGroup, serverName, databaseName, backupShortTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	if err = shortTermFuture.WaitForCompletionRef(ctx, shortTermPolicyClient.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Create/Update for SQL Server %q (Database %q) Short Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	// set short term policy id

	return resourceArmSQLServerDatabaseRetentionPolicyRead(d, meta)
}

func resourceArmSQLServerDatabaseRetentionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	longTermPolicyClient := meta.(*clients.Client).Sql.BackupLongTermRetentionPoliciesClient
	shortTermPolicyClient := meta.(*clients.Client).Sql.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// Read from id
	databaseName := d.Get("database_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	backupLongTermPolicy, err := longTermPolicyClient.Get(ctx, resourceGroup, serverName, databaseName)
	if err != nil {
		return fmt.Errorf("Error retrieving Long Term Policies for Database %q (SQL Server %q ;Resource Group %q): %+v", databaseName, serverName, resourceGroup, err)
	}

	flattenedLongTermPolicy := helper.FlattenSQLLongTermRetentionPolicyProperties(&backupLongTermPolicy)
	if err := d.Set("backup_long_term_retention_policy", flattenedLongTermPolicy); err != nil {
		return fmt.Errorf("Error setting `backup_long_term_retention_policy`: %+v", err)
	}

	backupShortTermPolicy, err := shortTermPolicyClient.Get(ctx, resourceGroup, serverName, databaseName)
	if err != nil {
		return fmt.Errorf("Error retrieving Short Term Policies for Database %q (SQL Server %q ;Resource Group %q): %+v", databaseName, serverName, resourceGroup, err)
	}

	flattenedShortTermPolicy := helper.FlattenSQLShortTermRetentionPolicy(&backupShortTermPolicy)
	if err := d.Set("backup_short_term_retention_policy", flattenedShortTermPolicy); err != nil {
		return fmt.Errorf("Error setting `backup_short_term_retention_policy`: %+v", err)
	}

	return nil
}

func resourceArmSQLServerDatabaseRetentionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	longTermPolicyClient := meta.(*clients.Client).Sql.BackupLongTermRetentionPoliciesClient
	shortTermPolicyClient := meta.(*clients.Client).Sql.BackupShortTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	// Update both resources to default values for removal
	return nil
}
