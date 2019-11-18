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

func resourceArmMSSqlDatabaseBlobAuditingPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSqlDataBaseBlobAuditingPoliciesCreateUpdate,
		Read:   resourceArmMSSqlDataBaseBlobAuditingPoliciesRead,
		Update: resourceArmMSSqlDataBaseBlobAuditingPoliciesCreateUpdate,
		Delete: resourceArmMSSqlDataBaseBlobAuditingPoliciesDelete,

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
		},
	}
}

func resourceArmMSSqlDataBaseBlobAuditingPoliciesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.DatabaseBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM SQL Database Blob Auditing Policies creation.")

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	databaseName := d.Get("database_name").(string)

	state := sql.BlobAuditingPolicyState(d.Get("state").(string))
	storageEndpoint := d.Get("storage_endpoint").(string)
	storageAccountAccessKey := d.Get("storage_account_access_key").(string)

	DatabaseBlobAuditingPolicyProperties := sql.DatabaseBlobAuditingPolicyProperties{
		State:                   state,
		StorageEndpoint:         &storageEndpoint,
		StorageAccountAccessKey: &storageAccountAccessKey,
	}
	//retention_days
	if retentionDays, ok := d.GetOk("retention_days"); ok {
		retentionDays := int32(retentionDays.(int))
		DatabaseBlobAuditingPolicyProperties.RetentionDays = &retentionDays
	}
	//audit_actions_and_groups
	if auditActionsAndGroups, ok := d.GetOk("audit_actions_and_groups"); ok {
		auditActionsAndGroups := strings.Split(auditActionsAndGroups.(string), ",")
		DatabaseBlobAuditingPolicyProperties.AuditActionsAndGroups = &auditActionsAndGroups
	}
	//storage_account_subscription_id
	if storageAccountSubscriptionID, ok := d.GetOk("storage_account_subscription_id"); ok {
		storageAccountSubscriptionID, _ := uuid.FromString(storageAccountSubscriptionID.(string))
		DatabaseBlobAuditingPolicyProperties.StorageAccountSubscriptionID = &storageAccountSubscriptionID
	}
	//is_storage_secondary_key_in_use
	if isStorageSecondaryKeyInUse, ok := d.GetOk("is_storage_secondary_key_in_use"); ok {
		isStorageSecondaryKeyInUse := isStorageSecondaryKeyInUse.(bool)
		DatabaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = &isStorageSecondaryKeyInUse
	}
	//is_azure_monitor_target_enabled
	if isAzureMonitorTargetEnabled, ok := d.GetOk("is_azure_monitor_target_enabled"); ok {
		isAzureMonitorTargetEnabled := isAzureMonitorTargetEnabled.(bool)
		DatabaseBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled = &isAzureMonitorTargetEnabled
	}

	parameters := sql.DatabaseBlobAuditingPolicy{
		DatabaseBlobAuditingPolicyProperties: &DatabaseBlobAuditingPolicyProperties,
	}
	resp, err := client.CreateOrUpdate(ctx, resGroup, serverName, databaseName, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q Database %q Blob Auditing Policies(Resource Group %q): %+v", serverName, databaseName, resGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read SQL Server '%s' Database %q Blob Auditing Policies (resource group %s) ID", serverName, databaseName, resGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmMSSqlDataBaseBlobAuditingPoliciesRead(d, meta)
}

func resourceArmMSSqlDataBaseBlobAuditingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.DatabaseBlobAuditingPoliciesClient
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
			log.Printf("[INFO] Error reading SQL Server %q Database %q Blob Auditing Policies - removing from state", serverName, databaseName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Server %s Database %s: %v Blob Auditing Policies", serverName, databaseName, err)
	}

	d.Set("server_name", serverName)
	d.Set("database_name", databaseName)
	d.Set("resource_group_name", resGroup)
	if databaseBlobAuditingPolicyProperties := resp.DatabaseBlobAuditingPolicyProperties; databaseBlobAuditingPolicyProperties != nil {
		d.Set("state", databaseBlobAuditingPolicyProperties.State)
		d.Set("audit_actions_and_groups", strings.Join(*databaseBlobAuditingPolicyProperties.AuditActionsAndGroups, ","))
		d.Set("is_azure_monitor_target_enabled", databaseBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled)
		d.Set("is_storage_secondary_key_in_use", databaseBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse)
		d.Set("retention_days", databaseBlobAuditingPolicyProperties.RetentionDays)
		d.Set("storage_account_subscription_id", databaseBlobAuditingPolicyProperties.StorageAccountSubscriptionID.String())
		d.Set("storage_endpoint", databaseBlobAuditingPolicyProperties.StorageEndpoint)
	}

	return nil
}

func resourceArmMSSqlDataBaseBlobAuditingPoliciesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.DatabaseBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	databaseName := id.Path["databases"]

	parameters := sql.DatabaseBlobAuditingPolicy{
		DatabaseBlobAuditingPolicyProperties: &sql.DatabaseBlobAuditingPolicyProperties{
			State: sql.BlobAuditingPolicyStateDisabled,
		},
	}
	_, err = client.CreateOrUpdate(ctx, resGroup, serverName, databaseName, parameters)
	if err != nil {
		return fmt.Errorf("Error deleting SQL Server %s Databse %s Blob Auditing Policies: %+v", serverName, databaseName, err)
	}

	return nil
}
