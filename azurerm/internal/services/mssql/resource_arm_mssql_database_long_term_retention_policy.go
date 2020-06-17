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

func resourceArmMsSqlDatabaseLongTermRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMsSqlDatabaseLongTermRetentionPolicyCreateUpdate,
		Read:   resourceArmMsSqlDatabaseLongTermRetentionPolicyRead,
		Update: resourceArmMsSqlDatabaseLongTermRetentionPolicyCreateUpdate,
		Delete: resourceArmMsSqlDatabaseLongTermRetentionPolicyDelete,

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
			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},
			// WeeklyRetention - The weekly retention policy for an LTR backup in an ISO 8601 format.
			"weekly_retention": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT0S",
				ValidateFunc: azure.ValidateLongTermRetentionPoliciesIsoFormat,
			},
			// MonthlyRetention - The monthly retention policy for an LTR backup in an ISO 8601 format.
			"monthly_retention": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT0S",
				ValidateFunc: azure.ValidateLongTermRetentionPoliciesIsoFormat,
			},
			// YearlyRetention - The yearly retention policy for an LTR backup in an ISO 8601 format.
			"yearly_retention": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PT0S",
				ValidateFunc: azure.ValidateLongTermRetentionPoliciesIsoFormat,
			},
			// WeekOfYear - The week of year to take the yearly backup in an ISO 8601 format.
			"week_of_year": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(1, 52),
			},
		},
	}
}

func resourceArmMsSqlDatabaseLongTermRetentionPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BackupLongTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseName := d.Get("database_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	serverName := d.Get("server_name").(string)

	backupLongTermPolicyProps := sql.LongTermRetentionPolicyProperties{
		WeeklyRetention:  utils.String("PT0S"),
		MonthlyRetention: utils.String("PT0S"),
		YearlyRetention:  utils.String("PT0S"),
		WeekOfYear:       utils.Int32(0),
	}

	if v, ok := d.GetOk("weekly_retention"); ok {
		backupLongTermPolicyProps.WeeklyRetention = utils.String(v.(string))
	}

	if v, ok := d.GetOk("monthly_retention"); ok {
		backupLongTermPolicyProps.MonthlyRetention = utils.String(v.(string))
	}

	if v, ok := d.GetOk("yearly_retention"); ok {
		backupLongTermPolicyProps.YearlyRetention = utils.String(v.(string))
	}

	if v, ok := d.GetOk("week_of_year"); ok {
		backupLongTermPolicyProps.WeekOfYear = utils.Int32(int32(v.(int)))
	}

	backupLongTermPolicy := sql.BackupLongTermRetentionPolicy{
		LongTermRetentionPolicyProperties: &backupLongTermPolicyProps,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, databaseName, backupLongTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Create/Update for Sql Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	response, err := client.Get(ctx, resourceGroup, serverName, databaseName)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Database %q Long Term Policies (Sql Server %q ,Resource Group %q): %+v", databaseName, serverName, resourceGroup, err)
	}

	d.SetId(*response.ID)

	return resourceArmMsSqlDatabaseLongTermRetentionPolicyRead(d, meta)
}

func resourceArmMsSqlDatabaseLongTermRetentionPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BackupLongTermRetentionPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	databaseName := id.Path["databases"]
	resourceGroup := id.ResourceGroup
	serverName := id.Path["servers"]

	backupLongTermPolicy, err := client.Get(ctx, resourceGroup, serverName, databaseName)
	if err != nil {
		return fmt.Errorf("Error retrieving Long Term Policies for Database %q (Sql Server %q ;Resource Group %q): %+v", databaseName, serverName, resourceGroup, err)
	}

	if backupLongTermPolicy.WeeklyRetention != nil {
		weeklyRetention := *backupLongTermPolicy.WeeklyRetention
		if err := d.Set("weekly_retention", weeklyRetention); err != nil {
			return fmt.Errorf("Error setting `weekly_retention`: %+v", err)
		}
	}

	if backupLongTermPolicy.MonthlyRetention != nil {
		monthlyRetention := *backupLongTermPolicy.MonthlyRetention
		if err := d.Set("monthly_retention", monthlyRetention); err != nil {
			return fmt.Errorf("Error setting `monthly_retention`: %+v", err)
		}
	}

	if backupLongTermPolicy.YearlyRetention != nil {
		yearlyRetention := *backupLongTermPolicy.YearlyRetention
		if err := d.Set("yearly_retention", yearlyRetention); err != nil {
			return fmt.Errorf("Error setting `yearly_retention`: %+v", err)
		}
	}

	if backupLongTermPolicy.WeekOfYear != nil {
		weekOfYear := *backupLongTermPolicy.WeekOfYear
		if err := d.Set("week_of_year", weekOfYear); err != nil {
			return fmt.Errorf("Error setting `week_of_year`: %+v", err)
		}
	}

	d.Set("database_name", databaseName)
	d.Set("resource_group_name", resourceGroup)
	d.Set("server_name", serverName)

	return nil
}

func resourceArmMsSqlDatabaseLongTermRetentionPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.BackupLongTermRetentionPoliciesClient
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
	backupLongTermPolicy := sql.BackupLongTermRetentionPolicy{
		LongTermRetentionPolicyProperties: &sql.LongTermRetentionPolicyProperties{
			WeeklyRetention:  utils.String("PT0S"),
			MonthlyRetention: utils.String("PT0S"),
			YearlyRetention:  utils.String("PT0S"),
			WeekOfYear:       utils.Int32(1),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, serverName, databaseName, backupLongTermPolicy)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Sql Server %q (Database %q) Long Term Retention Policies (Resource Group %q): %+v", serverName, databaseName, resourceGroup, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}
