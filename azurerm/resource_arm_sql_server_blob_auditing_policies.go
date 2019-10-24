package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
	"log"
	"strings"
	"time"
)

func resourceArmSqlServerBlobAuditingPolicies() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlServerBlobAuditingPoliciesCreateUpdate,
		Read:   resourceArmSqlServerBlobAuditingPoliciesRead,
		Update: resourceArmSqlServerBlobAuditingPoliciesCreateUpdate,
		Delete: resourceArmSqlServerBlobAuditingPoliciesCreateUpdate,

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
				Type:     schema.TypeString,
				Required: true,
			},

			"state": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceArmSqlServerBlobAuditingPoliciesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ServerBlobAuditingPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM SQL Server Blob Auditing Policies creation.")

	serverName := d.Get("server_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, serverName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing blob auditing policies of SQL Server %q Blob Auditing Policies(Resource Group %q): %+v", serverName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_server_blob_auditing_policies", *existing.ID)
		}
	}

	var state sql.BlobAuditingPolicyState
	switch d.Get("state").(string) {
	case "Enabled":
		state = sql.BlobAuditingPolicyStateEnabled
	case "Disabled":
		state = sql.BlobAuditingPolicyStateDisabled
	}
	storageEndpoint := d.Get("storage_endpoint").(string)

	//waitForStorageAccountAccessKeyAvail(ctx, d)
	storageAccountAccessKey := d.Get("storage_account_access_key").(string)

	ServerBlobAuditingPolicyProperties := sql.ServerBlobAuditingPolicyProperties{
		State:                   state,
		StorageEndpoint:         &storageEndpoint,
		StorageAccountAccessKey: &storageAccountAccessKey,
	}
	//retention_days
	if retentionDays, ok := d.GetOk("retention_days"); ok {
		retentionDays := int32(retentionDays.(int))
		ServerBlobAuditingPolicyProperties.RetentionDays = &retentionDays
	}
	//audit_actions_and_groups
	if auditActionsAndGroups, ok := d.GetOk("audit_actions_and_groups"); ok {
		auditActionsAndGroups := strings.Split(auditActionsAndGroups.(string), ",")
		ServerBlobAuditingPolicyProperties.AuditActionsAndGroups = &auditActionsAndGroups
	}
	//storage_account_subscription_id
	if storageAccountSubscriptionID, ok := d.GetOk("storage_account_subscription_id"); ok {
		storageAccountSubscriptionID, err := uuid.FromString(storageAccountSubscriptionID.(string))
		if err != nil {
			return fmt.Errorf("Error transforming storage account subscription id from string to uuid:%s", err)
		}
		ServerBlobAuditingPolicyProperties.StorageAccountSubscriptionID = &storageAccountSubscriptionID
	}
	//is_storage_secondary_key_in_use
	if isStorageSecondaryKeyInUse, ok := d.GetOk("is_storage_secondary_key_in_use"); ok {
		isStorageSecondaryKeyInUse := isStorageSecondaryKeyInUse.(bool)
		ServerBlobAuditingPolicyProperties.IsStorageSecondaryKeyInUse = &isStorageSecondaryKeyInUse
	}
	//is_azure_monitor_target_enabled
	if isAzureMonitorTargetEnabled, ok := d.GetOk("is_azure_monitor_target_enabled"); ok {
		isAzureMonitorTargetEnabled := isAzureMonitorTargetEnabled.(bool)
		ServerBlobAuditingPolicyProperties.IsAzureMonitorTargetEnabled = &isAzureMonitorTargetEnabled
	}

	parameters := sql.ServerBlobAuditingPolicy{
		ServerBlobAuditingPolicyProperties: &ServerBlobAuditingPolicyProperties,
	}
	future, err := client.CreateOrUpdate(ctx, resGroup, serverName, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for SQL Server %q Blob Auditing Policies(Resource Group %q): %+v", serverName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasConflict(future.Response()) {
			return fmt.Errorf("SQL Server names need to be globally unique and %q is already in use.", serverName)
		}

		return fmt.Errorf("Error waiting on create/update future for SQL Server %q Blob Auditing Policies (Resource Group %q): %+v", serverName, resGroup, err)
	}

	read, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("Error issuing get request for SQL Server %q Blob Auditing Policies (Resource Group %q): %+v", serverName, resGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read SQL Server '%s' Blob Auditing Policies (resource group %s) ID", serverName, resGroup)
	}
	d.SetId(*read.ID)

	return resourceArmSqlServerBlobAuditingPoliciesRead(d, meta)
}

func resourceArmSqlServerBlobAuditingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ServerBlobAuditingPoliciesClient
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
			log.Printf("[INFO] Error reading SQL Server %q Blob Auditing Policies - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Server %s: %v Blob Auditing Policies", name, err)
	}

	d.Set("server_name", name)
	d.Set("resource_group_name", resGroup)
	if serverProperties := resp.ServerBlobAuditingPolicyProperties; serverProperties != nil {
		d.Set("state", serverProperties.State)
		d.Set("audit_actions_and_groups", strings.Join(*serverProperties.AuditActionsAndGroups, ","))
		d.Set("is_azure_monitor_target_enabled", serverProperties.IsAzureMonitorTargetEnabled)
		d.Set("is_storage_secondary_key_in_use", serverProperties.IsStorageSecondaryKeyInUse)
		d.Set("retention_days", serverProperties.RetentionDays)
		d.Set("storage_account_subscription_id", serverProperties.StorageAccountSubscriptionID.String())
		d.Set("storage_endpoint", serverProperties.StorageEndpoint)
	}

	return nil
}
