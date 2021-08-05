package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceNetworkWatcherFlowLog() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkWatcherFlowLogCreateUpdate,
		Read:   resourceNetworkWatcherFlowLogRead,
		Update: resourceNetworkWatcherFlowLogCreateUpdate,
		Delete: resourceNetworkWatcherFlowLogDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FlowLogID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"network_watcher_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
				// TODO 3.0: Make this required, and remove computed.
				//Required: true,
				//ValidateFunc: validate.NetworkWatcherFlowLogName,
			},

			"network_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NetworkSecurityGroupID,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"retention_policy": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:             pluginsdk.TypeBool,
							Required:         true,
							DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyEnabledDiff,
						},

						"days": {
							Type:             pluginsdk.TypeInt,
							Required:         true,
							DiffSuppressFunc: azureRMSuppressFlowLogRetentionPolicyDaysDiff,
						},
					},
				},
			},

			"traffic_analytics": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Required: true,
						},

						"workspace_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},

						"workspace_region": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							StateFunc:        location.StateFunc,
							DiffSuppressFunc: location.DiffSuppressFunc,
						},

						"workspace_resource_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"interval_in_minutes": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntInSlice([]int{10, 60}),
							Default:      60,
						},
					},
				},
			},

			"version": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},

			"location": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateFunc:     location.EnhancedValidate,
				StateFunc:        location.StateFunc,
				DiffSuppressFunc: location.DiffSuppressFunc,
			},

			"tags": tags.Schema(),
		},
	}
}

func azureRMSuppressFlowLogRetentionPolicyEnabledDiff(_, old, _ string, d *pluginsdk.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `false` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func azureRMSuppressFlowLogRetentionPolicyDaysDiff(_, old, _ string, d *pluginsdk.ResourceData) bool {
	// Ignore if flow log is disabled as the returned flow log configuration
	// returns default value `0` which may differ from config
	return old != "" && !d.Get("enabled").(bool)
}

func resourceNetworkWatcherFlowLogCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroupName := d.Get("resource_group_name").(string)
	networkWatcherName := d.Get("network_watcher_name").(string)
	networkSecurityGroupID := d.Get("network_security_group_id").(string)

	// guaranteed via schema validation
	nsgId, _ := parse.NetworkSecurityGroupID(networkSecurityGroupID)
	id := parse.NewFlowLogID(subscriptionId, resourceGroupName, networkWatcherName, *nsgId)

	loc := d.Get("location").(string)
	if loc == "" {
		// Get the containing network watcher in order to reuse its location if the "location" is not specified.
		watcherClient := meta.(*clients.Client).Network.WatcherClient
		resp, err := watcherClient.Get(ctx, id.ResourceGroupName, id.NetworkWatcherName)
		if err != nil {
			return fmt.Errorf("retrieving %s: %v", parse.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName).ID(), err)
		}
		if resp.Location != nil {
			loc = *resp.Location
		}
	}

	parameters := network.FlowLog{
		Location: utils.String(location.Normalize(loc)),
		FlowLogPropertiesFormat: &network.FlowLogPropertiesFormat{
			TargetResourceID: utils.String(id.NetworkSecurityGroupID()),
			StorageID:        utils.String(d.Get("storage_account_id").(string)),
			Enabled:          utils.Bool(d.Get("enabled").(bool)),
			RetentionPolicy:  expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("traffic_analytics"); ok {
		parameters.FlowAnalyticsConfiguration = expandAzureRmNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	if version, ok := d.GetOk("version"); ok {
		format := &network.FlowLogFormatParameters{
			Version: utils.Int32(int32(version.(int))),
		}

		parameters.Format = format
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.NetworkWatcherName, id.Name(), parameters)
	if err != nil {
		return fmt.Errorf("Error creating %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkWatcherFlowLogRead(d, meta)
}

func resourceNetworkWatcherFlowLogRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlowLogID(d.Id())
	if err != nil {
		return err
	}

	// Get current flow log status
	resp, err := client.Get(ctx, id.ResourceGroupName, id.NetworkWatcherName, id.Name())
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving %q: %+v", id, err)
	}

	d.Set("network_watcher_name", id.NetworkWatcherName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("network_security_group_id", id.NetworkSecurityGroupID())
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("name", resp.Name)

	if prop := resp.FlowLogPropertiesFormat; prop != nil {
		if err := d.Set("traffic_analytics", flattenAzureRmNetworkWatcherFlowLogTrafficAnalytics(prop.FlowAnalyticsConfiguration)); err != nil {
			return fmt.Errorf("Error setting `traffic_analytics`: %+v", err)
		}

		d.Set("enabled", prop.Enabled)

		if format := prop.Format; format != nil {
			d.Set("version", format.Version)
		}

		// Azure API returns "" when flow log is disabled
		// Don't overwrite to prevent storage account ID diff when that is the case
		if prop.StorageID != nil && *prop.StorageID != "" {
			d.Set("storage_account_id", prop.StorageID)
		}

		if err := d.Set("retention_policy", flattenAzureRmNetworkWatcherFlowLogRetentionPolicy(prop.RetentionPolicy)); err != nil {
			return fmt.Errorf("Error setting `retention_policy`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNetworkWatcherFlowLogDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlowLogID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroupName, id.NetworkWatcherName, id.Name())
	if err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %v", id, err)
	}

	return nil
}

func expandAzureRmNetworkWatcherFlowLogRetentionPolicy(d *pluginsdk.ResourceData) *network.RetentionPolicyParameters {
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

func expandAzureRmNetworkWatcherFlowLogTrafficAnalytics(d *pluginsdk.ResourceData) *network.TrafficAnalyticsProperties {
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
