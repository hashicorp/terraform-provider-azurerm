package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMSSqlDatabaseBlobExtendedAuditingPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesCreateUpdate,
		Read:   resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesRead,
		Update: resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesCreateUpdate,
		Delete: resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesDelete,

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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},

			"database_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateMsSqlDatabaseName,
			},

			"state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Enabled", "Disabled"}, false),
			},

			"storage_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_account_access_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"audit_actions_and_groups": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_account_subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := uuid.FromString(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q is not in correct format:%+v", key, err))
					}
					return
				},
			},
			"is_storage_secondary_key_in_use": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_azure_monitor_target_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"predicate_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ExtendedDatabaseBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM SQL Database Blob Extended Auditing Policies creation.")

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	databaseName := d.Get("database_name").(string)

	state := sql.BlobAuditingPolicyState(d.Get("state").(string))
	storageEndpoint := d.Get("storage_endpoint").(string)
	storageAccountAccessKey := d.Get("storage_account_access_key").(string)

	ExtendedDatabaseBlobAuditingPolicyProperties := sql.ExtendedDatabaseBlobAuditingPolicyProperties{
		State:                   state,
		StorageEndpoint:         &storageEndpoint,
		StorageAccountAccessKey: &storageAccountAccessKey,
	}
	//retention_days
	if retentionDays, ok := d.GetOk("retention_days"); ok {
		retentionDays := int32(retentionDays.(int))
		ExtendedDatabaseBlobAuditingPolicyProperties.RetentionDays = &retentionDays
	}
	//audit_actions_and_groups
	if auditActionsAndGroups, ok := d.GetOk("audit_actions_and_groups"); ok {
		auditActionsAndGroups := strings.Split(auditActionsAndGroups.(string), ",")
		ExtendedDatabaseBlobAuditingPolicyProperties.AuditActionsAndGroups = &auditActionsAndGroups
	}
	//storage_account_subscription_id
	if storageAccountSubscriptionID, ok := d.GetOk("storage_account_subscription_id"); ok {
		storageAccountSubscriptionID, _ := uuid.FromString(storageAccountSubscriptionID.(string))
		ExtendedDatabaseBlobAuditingPolicyProperties.StorageAccountSubscriptionID = &storageAccountSubscriptionID
	}
	//is_storage_secondary_key_in_use
	if isStorageSecondaryKeyInUse, ok := d.GetOk("is_storage_secondary_key_in_use"); ok {
		isStorageSecondaryKeyInUse := isStorageSecondaryKeyInUse.(bool)
		ExtendedDatabaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = &isStorageSecondaryKeyInUse
	}
	//is_azure_monitor_target_enabled
	if isAzureMonitorTargetEnabled, ok := d.GetOk("is_azure_monitor_target_enabled"); ok {
		isAzureMonitorTargetEnabled := isAzureMonitorTargetEnabled.(bool)
		ExtendedDatabaseBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled = &isAzureMonitorTargetEnabled
	}
	//predicate_expression
	if predictExpression, ok := d.GetOk("predicate_expression"); ok {
		predictExpression := predictExpression.(string)
		ExtendedDatabaseBlobAuditingPolicyProperties.PredicateExpression = &predictExpression
	}

	parameters := sql.ExtendedDatabaseBlobAuditingPolicy{
		ExtendedDatabaseBlobAuditingPolicyProperties: &ExtendedDatabaseBlobAuditingPolicyProperties,
	}
	resp, err := client.CreateOrUpdate(ctx, resGroup, serverName, databaseName, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q Database %q Blob Extended Auditing Policies(Resource Group %q): %+v", serverName, databaseName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read SQL Server '%s' Database %q Blob Extended Auditing Policies (resource group %s) ID", serverName, databaseName, resGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesRead(d, meta)
}

func resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ExtendedDatabaseBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	databaseName := id.Path["databases"]
	resp, err := client.Get(ctx, resGroup, serverName, databaseName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Server %q Database %q Blob Extended Auditing Policies - removing from state", serverName, databaseName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Server %s Database %s: %v Blob Extended Auditing Policies", serverName, databaseName, err)
	}

	d.Set("server_name", serverName)
	d.Set("database_name", databaseName)
	d.Set("resource_group_name", resGroup)
	if ExtendedDatabaseBlobAuditingPolicyProperties := resp.ExtendedDatabaseBlobAuditingPolicyProperties; ExtendedDatabaseBlobAuditingPolicyProperties != nil {
		d.Set("state", ExtendedDatabaseBlobAuditingPolicyProperties.State)
		d.Set("audit_actions_and_groups", strings.Join(*ExtendedDatabaseBlobAuditingPolicyProperties.AuditActionsAndGroups, ","))
		d.Set("is_azure_monitor_target_enabled", ExtendedDatabaseBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled)
		d.Set("is_storage_secondary_key_in_use", ExtendedDatabaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse)
		d.Set("retention_days", ExtendedDatabaseBlobAuditingPolicyProperties.RetentionDays)
		d.Set("storage_account_subscription_id", ExtendedDatabaseBlobAuditingPolicyProperties.StorageAccountSubscriptionID.String())
		d.Set("storage_endpoint", ExtendedDatabaseBlobAuditingPolicyProperties.StorageEndpoint)
		d.Set("predicate_expression", ExtendedDatabaseBlobAuditingPolicyProperties.PredicateExpression)
	}

	return nil
}

func resourceArmMSSqlDatabaseBlobExtendedAuditingPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ExtendedDatabaseBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	databaseName := id.Path["databases"]

	parameters := sql.ExtendedDatabaseBlobAuditingPolicy{
		ExtendedDatabaseBlobAuditingPolicyProperties: &sql.ExtendedDatabaseBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		},
	}
	_, err = client.CreateOrUpdate(ctx, resGroup, serverName, databaseName, parameters)
	if err != nil {
		return fmt.Errorf("Error deleting SQL Server %s Database %s Blob Extended Auditing Policies: %+v", serverName, databaseName, err)
	}

	return nil
}
