package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type NetworkWatcherFlowLogAccountID struct {
	azure.ResourceID
	NetworkWatcherName     string
	NetworkSecurityGroupID string
}

func ParseNetworkWatcherFlowLogID(id string) (*NetworkWatcherFlowLogAccountID, error) {
	parts := strings.Split(id, "/networkSecurityGroupId")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Error: Network Watcher Flow Log ID could not be split on `/networkSecurityGroupId`: %s", id)
	}

	watcherId, err := azure.ParseAzureResourceID(parts[0])
	if err != nil {
		return nil, err
	}

	watcherName, ok := watcherId.Path["networkWatchers"]
	if !ok {
		return nil, fmt.Errorf("Error: Unable to parse Network Watcher Flow Log ID: networkWatchers is missing from: %s", id)
	}

	return &NetworkWatcherFlowLogAccountID{
		ResourceID:             *watcherId,
		NetworkWatcherName:     watcherName,
		NetworkSecurityGroupID: parts[1],
	}, nil
}

func resourceNetworkWatcherFlowLog() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkWatcherFlowLogCreateUpdate,
		Read:   resourceNetworkWatcherFlowLogRead,
		Update: resourceNetworkWatcherFlowLogCreateUpdate,
		Delete: resourceNetworkWatcherFlowLogDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"network_watcher_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
							ValidateFunc: validation.IsUUID,
						},

						"workspace_region": {
							Type:             schema.TypeString,
							Required:         true,
							StateFunc:        location.StateFunc,
							DiffSuppressFunc: location.DiffSuppressFunc,
						},

						"workspace_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"interval_in_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{10, 60}),
							Default:      60,
						},
					},
				},
			},

			"version": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},
		},
	}
}

func azureRMSuppressFlowLogRetentionPolicyEnabledDiff(_, old, _ string, d *schema.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `false` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func azureRMSuppressFlowLogRetentionPolicyDaysDiff(_, old, _ string, d *schema.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `0` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func resourceNetworkWatcherFlowLogCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WatcherClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

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

	if version, ok := d.GetOk("version"); ok {
		format := &network.FlowLogFormatParameters{
			Version: utils.Int32(int32(version.(int))),
		}

		parameters.FlowLogProperties.Format = format
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
		return fmt.Errorf("Cannot read Network Watcher %q (Resource Group %q) err: %+v", networkWatcherName, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Network Watcher %q is nil (Resource Group %q)", networkWatcherName, resourceGroupName)
	}

	d.SetId(*resp.ID + "/networkSecurityGroupId" + networkSecurityGroupID)

	return resourceNetworkWatcherFlowLogRead(d, meta)
}

func resourceNetworkWatcherFlowLogRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WatcherClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseNetworkWatcherFlowLogID(d.Id())
	if err != nil {
		return err
	}

	// Get current flow log status
	statusParameters := network.FlowLogStatusParameters{
		TargetResourceID: &id.NetworkSecurityGroupID,
	}

	future, err := client.GetFlowLogStatus(ctx, id.ResourceGroup, id.NetworkWatcherName, statusParameters)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			// One of storage account, NSG, or flow log is missing
			log.Printf("[INFO] Error getting Flow Log Configuration %q for target %q - removing from state", d.Id(), id.NetworkSecurityGroupID)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	fli, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("Error retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	d.Set("network_watcher_name", id.NetworkWatcherName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("network_security_group_id", fli.TargetResourceID)
	if err := d.Set("traffic_analytics", flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(fli.FlowAnalyticsConfiguration)); err != nil {
		return fmt.Errorf("Error setting `traffic_analytics`: %+v", err)
	}

	if props := fli.FlowLogProperties; props != nil {
		d.Set("enabled", props.Enabled)

		if format := props.Format; format != nil {
			d.Set("version", format.Version)
		}

		// Azure API returns "" when flow log is disabled
		// Don't overwrite to prevent storage account ID diff when that is the case
		if props.StorageID != nil && *props.StorageID != "" {
			d.Set("storage_account_id", props.StorageID)
		}

		if err := d.Set("retention_policy", flattenAzureRmNetworkWatcherFlowLogRetentionPolicy(props.RetentionPolicy)); err != nil {
			return fmt.Errorf("Error setting `retention_policy`: %+v", err)
		}
	}

	return nil
}

func resourceNetworkWatcherFlowLogDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.WatcherClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParseNetworkWatcherFlowLogID(d.Id())
	if err != nil {
		return err
	}

	// Get current flow log status
	statusParameters := network.FlowLogStatusParameters{
		TargetResourceID: &id.NetworkSecurityGroupID,
	}
	future, err := client.GetFlowLogStatus(ctx, id.ResourceGroup, id.NetworkWatcherName, statusParameters)
	if err != nil {
		return fmt.Errorf("getting Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for retrieval of Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	fli, err := future.Result(*client)
	if err != nil {
		return fmt.Errorf("retrieving Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
	}

	// There is no delete in Azure API. Disabling flow log is effectively a delete in Terraform.
	if props := fli.FlowLogProperties; props != nil {
		if props.Enabled != nil && *props.Enabled {
			props.Enabled = utils.Bool(false)

			param := network.FlowLogInformation{
				TargetResourceID: &id.NetworkSecurityGroupID,
				FlowLogProperties: &network.FlowLogProperties{
					StorageID: utils.String(*fli.StorageID),
					Enabled:   utils.Bool(false),
				},
				FlowAnalyticsConfiguration: &network.TrafficAnalyticsProperties{
					NetworkWatcherFlowAnalyticsConfiguration: &network.TrafficAnalyticsConfigurationProperties{
						Enabled: utils.Bool(false),
					},
				},
			}
			setFuture, err := client.SetFlowLogConfiguration(ctx, id.ResourceGroup, id.NetworkWatcherName, param)
			if err != nil {
				return fmt.Errorf("disabling Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
			}

			if err = setFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for completion of disabling Flow Log Configuration for target %q (Network Watcher %q / Resource Group %q): %+v", id.NetworkSecurityGroupID, id.NetworkWatcherName, id.ResourceGroup, err)
			}
		}
	}

	return nil
}

func expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d *schema.ResourceData) *network.RetentionPolicyParameters {
	vs := d.Get("retention_policy").([]interface{})
	if len(vs) < 1 || vs[0] == nil {
		return nil
	}

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

	if input.Enabled != nil {
		result["enabled"] = *input.Enabled
	}
	if input.Days != nil {
		result["days"] = *input.Days
	}

	return []interface{}{result}
}

func flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(input *network.TrafficAnalyticsProperties) []interface{} {
	if input == nil || input.NetworkWatcherFlowAnalyticsConfiguration == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

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
		if cfg.TrafficAnalyticsInterval != nil {
			result["interval_in_minutes"] = int(*cfg.TrafficAnalyticsInterval)
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
	interval := v["interval_in_minutes"].(int)

	return &network.TrafficAnalyticsProperties{
		NetworkWatcherFlowAnalyticsConfiguration: &network.TrafficAnalyticsConfigurationProperties{
			Enabled:                  utils.Bool(enabled),
			WorkspaceID:              utils.String(workspaceID),
			WorkspaceRegion:          utils.String(workspaceRegion),
			WorkspaceResourceID:      utils.String(workspaceResourceID),
			TrafficAnalyticsInterval: utils.Int32(int32(interval)),
		},
	}
}
