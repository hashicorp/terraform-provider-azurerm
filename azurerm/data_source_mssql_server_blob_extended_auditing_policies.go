package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlServerBlobExtendedAuditingPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlServerBlobExtendedAuditingPoliciesRead,

		Schema: map[string]*schema.Schema{

			"resource_group_name": azure.SchemaResourceGroupName(),

			"server_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateMsSqlServerName,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_account_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"audit_actions_and_groups": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_account_subscription_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_storage_secondary_key_in_use": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_azure_monitor_target_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"predicate_expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmMsSqlServerBlobExtendedAuditingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Sql.ExtendedServerBlobAuditingPoliciesClient
	ctx := meta.(*ArmClient).StopContext

	serverName := d.Get("server_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, serverName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading SQL Server %q Blob Extended Auditing Policies - removing from state", serverName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading SQL Server %s: %v Blob Extended Auditing Policies", serverName, err)
	}

	d.SetId(*resp.ID)

	d.Set("server_name", serverName)
	d.Set("resource_group_name", resourceGroup)

	d.Set("state", resp.State)
	if auditActionsAndGroups := resp.AuditActionsAndGroups; auditActionsAndGroups != nil {
		d.Set("audit_actions_and_groups", strings.Join(*auditActionsAndGroups, ","))
	}
	d.Set("is_azure_monitor_target_enabled", resp.IsAzureMonitorTargetEnabled)
	d.Set("is_storage_secondary_key_in_use", resp.IsStorageSecondaryKeyInUse)
	d.Set("retention_days", resp.RetentionDays)
	if StorageAccountSubscriptionID := resp.StorageAccountSubscriptionID; StorageAccountSubscriptionID != nil {
		d.Set("storage_account_subscription_id", StorageAccountSubscriptionID.String())
	}
	if storageEndpoint := resp.StorageEndpoint; storageEndpoint != nil {
		d.Set("storage_endpoint", storageEndpoint)
	}
	if predicateExpression := resp.PredicateExpression; predicateExpression != nil {
		d.Set("predicate_expression", predicateExpression)
	}

	return nil
}
