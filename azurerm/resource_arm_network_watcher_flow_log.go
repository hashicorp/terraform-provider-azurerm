package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNetworkWatcherFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkWatcherFlowLogCreateUpdate,
		Read:   resourceArmNetworkWatcherFlowLogRead,
		Update: resourceArmNetworkWatcherFlowLogCreateUpdate,
		Delete: resourceArmNetworkWatcherFlowLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"network_watcher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"network_security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"retention_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:             schema.TypeBool,
							Required:         true,
							DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyEnabledDiff,
						},
						"days": {
							Type:             schema.TypeInt,
							Required:         true,
							DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyDaysDiff,
						},
					},
				},
			},

			"traffic_analytics": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"workspace_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.UUID,
						},
						"workspace_region": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        azureRMNormalizeLocation,
							DiffSuppressFunc: azureRMSuppressLocationDiff,
						},
						"workspace_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},
					},
				},
			},
		},
	}
}

func azureRMSuppressFlowLogRetentionPolicyEnabledDiff(k, old, new string, d *schema.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `false` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func azureRMSuppressFlowLogRetentionPolicyDaysDiff(k, old, new string, d *schema.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `0` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func resourceArmNetworkWatcherFlowLogCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient
	ctx := meta.(*ArmClient).StopContext

	networkWatcherName := d.Get("network_watcher_name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	networkSecurityGroupID := d.Get("network_security_group_id").(string)
	storageAccountID := d.Get("storage_account_id").(string)
	enabled := d.Get("enabled").(bool)

	parameters := network.FlowLogInformation{
		TargetResourceID: &networkSecurityGroupID,
		FlowLogProperties: &network.FlowLogProperties{
			StorageID:       &storageAccountID,
			Enabled:         &enabled,
			RetentionPolicy: expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d),
		},
	}

	if _, ok := d.GetOk("traffic_analytics"); ok {
		parameters.FlowAnalyticsConfiguration = expandAzureRmNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	future, err := client.SetFlowLogConfiguration(ctx, resourceGroupName, networkWatcherName, parameters)
	if err != nil {
		return fmt.Errorf("Error setting Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of setting Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, networkWatcherName)
	if err != nil {
		return fmt.Errorf("Cannot read Network Watcher %q (Resource Group %q) ID: %+v", networkWatcherName, resourceGroupName, err)
	}
	d.SetId(*resp.ID)

	return resourceArmNetworkWatcherFlowLogRead(d, meta)
}

func resourceArmNetworkWatcherFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	networkWatcherName := id.Path["networkWatchers"]

	// Get current flow log status
	networkSecurityGroupID := d.Get("network_security_group_id").(string)
	statusParameters := network.FlowLogStatusParameters{
		TargetResourceID: &networkSecurityGroupID,
	}
	future, err := client.GetFlowLogStatus(ctx, resourceGroupName, networkWatcherName, statusParameters)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			// One of storage account, NSG, or flow log is missing
			log.Printf("[INFO] Error getting Flow Log Configuration %q for target %q - removing from state", d.Id(), networkSecurityGroupID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	fli, err := future.Result(client)
	if err != nil {
		return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	d.Set("network_watcher_name", networkWatcherName)
	d.Set("resource_group_name", resourceGroupName)

	d.Set("network_security_group_id", fli.TargetResourceID)
	if err := d.Set("traffic_analytics", flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(fli.FlowAnalyticsConfiguration)); err != nil {
		return fmt.Errorf("Error setting `traffic_analytics`: %+v", err)
	}

	if props := fli.FlowLogProperties; props != nil {
		d.Set("enabled", props.Enabled)

		// Azure API returns "" when flow log is disabled
		// Don't overwrite to prevent storage account ID diff when that is the case
		if *props.StorageID != "" {
			d.Set("storage_account_id", props.StorageID)
		}

		if err := d.Set("retention_policy", flattenAzureRmNetworkWatcherFlowLogRetentionPolicy(props.RetentionPolicy)); err != nil {
			return fmt.Errorf("Error setting `retention_policy`: %+v", err)
		}
	}

	return nil
}

func resourceArmNetworkWatcherFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).watcherClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	networkWatcherName := id.Path["networkWatchers"]

	// Get current flow log status
	networkSecurityGroupID := d.Get("network_security_group_id").(string)
	statusParameters := network.FlowLogStatusParameters{
		TargetResourceID: &networkSecurityGroupID,
	}
	future, err := client.GetFlowLogStatus(ctx, resourceGroupName, networkWatcherName, statusParameters)
	if err != nil {
		return fmt.Errorf("Error getting Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	fli, err := future.Result(client)
	if err != nil {
		return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
	}

	// There is no delete in Azure API. Disabling flow log is effectively a delete in Terraform.
	if *fli.FlowLogProperties.Enabled {
		fli.FlowLogProperties.Enabled = utils.Bool(false)

		if isDefaultDisabledFlowLogTrafficAnalytics(fli.FlowAnalyticsConfiguration) {
			fli.FlowAnalyticsConfiguration = nil
		}

		setFuture, err := client.SetFlowLogConfiguration(ctx, resourceGroupName, networkWatcherName, fli)
		if err != nil {
			return fmt.Errorf("Error disabling Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}

		if err = setFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for completion of disabling Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", networkSecurityGroupID, networkWatcherName, resourceGroupName, err)
		}
	}

	return nil
}

func expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d *schema.ResourceData) *network.RetentionPolicyParameters {
	vs := d.Get("retention_policy").([]interface{})
	v := vs[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	days := v["days"].(int)

	return &network.RetentionPolicyParameters{
		Enabled: utils.Bool(enabled),
		Days:    utils.Int32(int32(days)),
	}
}

func flattenAzureRmNetworkWatcherFlowLogRetentionPolicy(input *network.RetentionPolicyParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	if input != nil {
		if input.Enabled != nil {
			result["enabled"] = *input.Enabled
		}
		if input.Days != nil {
			result["days"] = *input.Days
		}
	}

	return []interface{}{result}
}

func flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(input *network.TrafficAnalyticsProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	} else if isDefaultDisabledFlowLogTrafficAnalytics(input) {
		return []interface{}{}
	}

	result := make(map[string]interface{})
	if input != nil {
		if cfg := input.NetworkWatcherFlowAnalyticsConfiguration; cfg != nil {
			if cfg.Enabled != nil {
				result["enabled"] = *cfg.Enabled
			}
			if cfg.WorkspaceID != nil {
				result["workspace_id"] = *cfg.WorkspaceID
			}
			if cfg.WorkspaceRegion != nil {
				result["workspace_region"] = *cfg.WorkspaceRegion
			}
			if cfg.WorkspaceResourceID != nil {
				result["workspace_resource_id"] = *cfg.WorkspaceResourceID
			}
		}
	}

	return []interface{}{result}
}

func expandAzureRmNetworkWatcherFlowLogTrafficAnalytics(d *schema.ResourceData) *network.TrafficAnalyticsProperties {
	vs := d.Get("traffic_analytics").([]interface{})

	v := vs[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	workspaceID := v["workspace_id"].(string)
	workspaceRegion := v["workspace_region"].(string)
	workspaceResourceID := v["workspace_resource_id"].(string)

	return &network.TrafficAnalyticsProperties{
		NetworkWatcherFlowAnalyticsConfiguration: &network.TrafficAnalyticsConfigurationProperties{
			Enabled:             utils.Bool(enabled),
			WorkspaceID:         utils.String(workspaceID),
			WorkspaceRegion:     utils.String(workspaceRegion),
			WorkspaceResourceID: utils.String(workspaceResourceID),
		},
	}
}

func isDefaultDisabledFlowLogTrafficAnalytics(input *network.TrafficAnalyticsProperties) bool {
	// Azure returns `/subscriptions//resourcegroups//providers/microsoft.operationalinsights/workspaces/` by default when traffic analytics is not set
	// along with empty strings for the rest of the values
	if cfg := input.NetworkWatcherFlowAnalyticsConfiguration; cfg != nil {
		return !(cfg.Enabled != nil && *cfg.Enabled &&
			cfg.WorkspaceID != nil && *cfg.WorkspaceID != "" &&
			cfg.WorkspaceRegion != nil && *cfg.WorkspaceRegion != "" &&
			cfg.WorkspaceResourceID != nil && *cfg.WorkspaceResourceID != "/subscriptions//resourcegroups//providers/microsoft.operationalinsights/workspaces/")
	}

	return true
}
