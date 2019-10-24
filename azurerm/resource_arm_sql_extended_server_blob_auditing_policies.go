package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSqlExtendedServerBlobAuditingPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlExtendedServerBlobAuditingPoliciesCreateUpdate,
		Read:   resourceArmSqlExtendedServerBlobAuditingPoliciesRead,
		Update: resourceArmSqlExtendedServerBlobAuditingPoliciesCreateUpdate,
		Delete: resourceArmSqlExtendedServerBlobAuditingPoliciesCreateUpdate,

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

			"state": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !(v == "Enabled" || v == "Disabled") {
						errs = append(errs, fmt.Errorf("%q can only be 'Enabled' or 'Disabled' ", key))
					}
					return
				},
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

func resourceArmSqlExtendedServerBlobAuditingPoliciesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ExtendedServerBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM SQL Server Extended Blob Auditing Policies creation.")

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, serverName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing blob auditing policies of SQL Server %q Extended Blob Auditing Policies(Resource Group %q): %+v", serverName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_server_blob_auditing_policies", *existing.ID)
		}
	}

	state := sql.BlobAuditingPolicyState(d.Get("state").(string))
	storageEndpoint := d.Get("storage_endpoint").(string)
	storageAccountAccessKey := d.Get("storage_account_access_key").(string)

	ExtendedServerBlobAuditingPolicyProperties := sql.ExtendedServerBlobAuditingPolicyProperties{
		State:                   state,
		StorageEndpoint:         &storageEndpoint,
		StorageAccountAccessKey: &storageAccountAccessKey,
	}
	//retention_days
	if retentionDays, ok := d.GetOk("retention_days"); ok {
		retentionDays := int32(retentionDays.(int))
		ExtendedServerBlobAuditingPolicyProperties.RetentionDays = &retentionDays
	}
	//audit_actions_and_groups
	if auditActionsAndGroups, ok := d.GetOk("audit_actions_and_groups"); ok {
		auditActionsAndGroups := strings.Split(auditActionsAndGroups.(string), ",")
		ExtendedServerBlobAuditingPolicyProperties.AuditActionsAndGroups = &auditActionsAndGroups
	}
	//storage_account_subscription_id
	if storageAccountSubscriptionID, ok := d.GetOk("storage_account_subscription_id"); ok {
		storageAccountSubscriptionID, _ := uuid.FromString(storageAccountSubscriptionID.(string))
		ExtendedServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID = &storageAccountSubscriptionID
	}
	//is_storage_secondary_key_in_use
	if isStorageSecondaryKeyInUse, ok := d.GetOk("is_storage_secondary_key_in_use"); ok {
		isStorageSecondaryKeyInUse := isStorageSecondaryKeyInUse.(bool)
		ExtendedServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = &isStorageSecondaryKeyInUse
	}
	//is_azure_monitor_target_enabled
	if isAzureMonitorTargetEnabled, ok := d.GetOk("is_azure_monitor_target_enabled"); ok {
		isAzureMonitorTargetEnabled := isAzureMonitorTargetEnabled.(bool)
		ExtendedServerBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled = &isAzureMonitorTargetEnabled
	}
	//predicate_expression
	if predictExpression, ok := d.GetOk("predicate_expression"); ok {
		predictExpression := predictExpression.(string)
		ExtendedServerBlobAuditingPolicyProperties.PredicateExpression = &predictExpression
	}

	parameters := sql.ExtendedServerBlobAuditingPolicy{
		ExtendedServerBlobAuditingPolicyProperties: &ExtendedServerBlobAuditingPolicyProperties,
	}
	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q Extended Blob Auditing Policies(Resource Group %q): %+v", serverName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use.", serverName)
		}

		return fmt.Errorf("Error waiting on create/update future for SQL Server %q Extended Blob Auditing Policies (Resource Group %q): %+v", serverName, resGroup, err)
	}

	read, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("Error issuing get request for SQL Server %q Extended Blob Auditing Policies (Resource Group %q): %+v", serverName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read SQL Server '%s' Extended Blob Auditing Policies (resource group %s) ID", serverName, resGroup)
	}
	d.SetId(*read.ID)

	return resourceArmSqlServerBlobAuditingPoliciesRead(d, meta)
}

func resourceArmSqlExtendedServerBlobAuditingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ExtendedServerBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["servers"]
	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Extended Server %q  Blob Auditing Policies - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Extended Server %s: %v  Blob Auditing Policies", name, err)
	}

	d.Set("server_name", name)
	d.Set("resource_group_name", resGroup)
	if serverProperties := resp.ExtendedServerBlobAuditingPolicyProperties; serverProperties != nil {
		d.Set("state", serverProperties.State)
		d.Set("audit_actions_and_groups", strings.Join(*serverProperties.AuditActionsAndGroups, ","))
		d.Set("is_azure_monitor_target_enabled", serverProperties.IsAzureMonitorTargetEnabled)
		d.Set("is_storage_secondary_key_in_use", serverProperties.IsStorageSecondaryKeyInUse)
		d.Set("retention_days", serverProperties.RetentionDays)
		d.Set("storage_account_subscription_id", serverProperties.StorageAccountSubscriptionID.String())
		d.Set("storage_endpoint", serverProperties.StorageEndpoint)
		d.Set("predicate_expression", serverProperties.PredicateExpression)
	}

	return nil
}
