package sql

import (
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
	client := meta.(*clients.Client).Sql.BackupLongTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	client.CreateOrUpdate(ctx, "rgname", "serverName", "databaseName", sql.BackupLongTermRetentionPolicy{
		LongTermRetentionPolicyProperties: &sql.LongTermRetentionPolicyProperties{
			MonthlyRetention: utils.String("3"),
			WeekOfYear:       utils.Int32(3),
			WeeklyRetention:  utils.String("3"),
			YearlyRetention:  utils.String("3"),
		},
	})
	return nil
}

func resourceArmSQLServerDatabaseRetentionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceArmSQLServerDatabaseRetentionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
