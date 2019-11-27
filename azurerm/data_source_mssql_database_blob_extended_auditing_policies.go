package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMsSqlDatabaseBlobExtendedAuditingPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmMsSqlDatabaseBlobExtendedAuditingPoliciesRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
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

func dataSourceArmMsSqlDatabaseBlobExtendedAuditingPoliciesRead(d *schema.ResourceData, meta interface{}) error {
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
