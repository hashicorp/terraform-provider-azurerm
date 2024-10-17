// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkwatchers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkWatcherFlowLog() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceNetworkWatcherFlowLogCreate,
		Read:   resourceNetworkWatcherFlowLogRead,
		Update: resourceNetworkWatcherFlowLogUpdate,
		Delete: resourceNetworkWatcherFlowLogDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkWatcherFlowLogV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := flowlogs.ParseFlowLogID(id)
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

			"resource_group_name": commonschema.ResourceGroupName(),

			// lintignore: S013
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NetworkWatcherFlowLogName,
			},

			"network_security_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
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
							ValidateFunc: azure.ValidateResourceIDOrEmpty, // nolint: staticcheck
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
				Default:      1,
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

			"tags": commonschema.Tags(),
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["version"] = &pluginsdk.Schema{
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 2),
		}
	}

	return resource
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

func resourceNetworkWatcherFlowLogCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := flowlogs.NewFlowLogID(subscriptionId, d.Get("resource_group_name").(string), d.Get("network_watcher_name").(string), d.Get("name").(string))
	nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Get("network_security_group_id").(string))
	if err != nil {
		return err
	}

	// For newly created resources, the "name" is required, it is set as Optional and Computed is merely for the existing ones for the sake of backward compatibility.
	if id.NetworkWatcherName == "" {
		return fmt.Errorf("`name` is required for Network Watcher Flow Log")
	}

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_watcher_flow_log", id.ID())
	}

	locks.ByID(nsgId.ID())
	defer locks.UnlockByID(nsgId.ID())

	loc := d.Get("location").(string)
	if loc == "" {
		// Get the containing network watcher in order to reuse its location if the "location" is not specified.
		watcherClient := meta.(*clients.Client).Network.NetworkWatchers
		watcherId := networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName)
		resp, err := watcherClient.Get(ctx, watcherId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %v", watcherId, err)
		}
		if model := resp.Model; model != nil && model.Location != nil {
			loc = *model.Location
		}
	}

	parameters := flowlogs.FlowLog{
		Location: utils.String(location.Normalize(loc)),
		Properties: &flowlogs.FlowLogPropertiesFormat{
			TargetResourceId: nsgId.ID(),
			StorageId:        d.Get("storage_account_id").(string),
			Enabled:          pointer.To(d.Get("enabled").(bool)),
			RetentionPolicy:  expandNetworkWatcherFlowLogRetentionPolicy(d.Get("retention_policy").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("traffic_analytics"); ok {
		parameters.Properties.FlowAnalyticsConfiguration = expandNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	if version, ok := d.GetOk("version"); ok {
		format := &flowlogs.FlowLogFormatParameters{
			Version: pointer.To(int64(version.(int))),
		}

		parameters.Properties.Format = format
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkWatcherFlowLogRead(d, meta)
}

func resourceNetworkWatcherFlowLogUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := flowlogs.ParseFlowLogID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(d.Get("network_security_group_id").(string))
	if err != nil {
		return err
	}
	locks.ByID(nsgId.ID())
	defer locks.UnlockByID(nsgId.ID())

	if d.HasChange("storage_account_id") {
		payload.Properties.StorageId = d.Get("storage_account_id").(string)
	}

	if d.HasChange("enabled") {
		payload.Properties.Enabled = pointer.To(d.Get("enabled").(bool))
	}

	if d.HasChange("retention_policy") {
		payload.Properties.RetentionPolicy = expandNetworkWatcherFlowLogRetentionPolicy(d.Get("retention_policy").([]interface{}))
	}

	if d.HasChange("traffic_analytics") {
		payload.Properties.FlowAnalyticsConfiguration = expandNetworkWatcherFlowLogTrafficAnalytics(d)
	}

	if d.HasChange("version") {
		if version, ok := d.GetOk("version"); ok && version != 0 {
			payload.Properties.Format = &flowlogs.FlowLogFormatParameters{
				Version: pointer.To(int64(d.Get("version").(int))),
			}
		} else {
			payload.Properties.Format = &flowlogs.FlowLogFormatParameters{}
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceNetworkWatcherFlowLogRead(d, meta)
}

func resourceNetworkWatcherFlowLogRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := flowlogs.ParseFlowLogID(d.Id())
	if err != nil {
		return err
	}

	// Get current flow log status
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("name", id.FlowLogName)
	d.Set("network_watcher_name", id.NetworkWatcherName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("traffic_analytics", flattenNetworkWatcherFlowLogTrafficAnalytics(props.FlowAnalyticsConfiguration)); err != nil {
				return fmt.Errorf("setting `traffic_analytics`: %+v", err)
			}

			d.Set("enabled", props.Enabled)

			version := 0
			if format := props.Format; format != nil {
				version = int(*format.Version)
			}
			d.Set("version", version)

			// Azure API returns "" when flow log is disabled
			// Don't overwrite to prevent storage account ID diff when that is the case
			if props.StorageId != "" {
				d.Set("storage_account_id", props.StorageId)
			}

			networkSecurityGroupId := ""
			nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(props.TargetResourceId)
			if err == nil {
				networkSecurityGroupId = nsgId.ID()
			}
			d.Set("network_security_group_id", networkSecurityGroupId)

			if err := d.Set("retention_policy", flattenNetworkWatcherFlowLogRetentionPolicy(props.RetentionPolicy)); err != nil {
				return fmt.Errorf("setting `retention_policy`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetworkWatcherFlowLogDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := flowlogs.ParseFlowLogID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.TargetResourceId == "" {
		return fmt.Errorf("retreiving %s: `properties` or `properties.TargetResourceID` was nil", id)
	}

	networkSecurityGroupId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(resp.Model.Properties.TargetResourceId)
	if err != nil {
		return fmt.Errorf("parsing %q as a Network Security Group ID: %+v", resp.Model.Properties.TargetResourceId, err)
	}

	locks.ByID(networkSecurityGroupId.ID())
	defer locks.UnlockByID(networkSecurityGroupId.ID())

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %v", id, err)
	}

	return nil
}

func expandNetworkWatcherFlowLogRetentionPolicy(input []interface{}) *flowlogs.RetentionPolicyParameters {
	if len(input) < 1 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	days := v["days"].(int)

	return &flowlogs.RetentionPolicyParameters{
		Enabled: pointer.To(enabled),
		Days:    pointer.To(int64(days)),
	}
}

func flattenNetworkWatcherFlowLogRetentionPolicy(input *flowlogs.RetentionPolicyParameters) []interface{} {
	output := make([]interface{}, 0)

	if input != nil {
		enabled := false
		if input.Enabled != nil {
			enabled = *input.Enabled
		}
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}
		output = append(output, map[string]interface{}{
			"days":    days,
			"enabled": enabled,
		})
	}

	return output
}

func flattenNetworkWatcherFlowLogTrafficAnalytics(input *flowlogs.TrafficAnalyticsProperties) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		if cfg := input.NetworkWatcherFlowAnalyticsConfiguration; cfg != nil {
			enabled := false
			if cfg.Enabled != nil {
				enabled = *cfg.Enabled
			}
			workspaceId := ""
			if cfg.WorkspaceId != nil {
				workspaceId = *cfg.WorkspaceId
			}
			workspaceRegion := ""
			if cfg.WorkspaceRegion != nil {
				workspaceRegion = *cfg.WorkspaceRegion
			}
			workspaceResourceId := ""
			if cfg.WorkspaceResourceId != nil {
				workspaceResourceId = *cfg.WorkspaceResourceId
			}

			intervalInMinutes := 0
			if cfg.TrafficAnalyticsInterval != nil {
				intervalInMinutes = int(*cfg.TrafficAnalyticsInterval)
			}
			output = append(output, map[string]interface{}{
				"enabled":               enabled,
				"interval_in_minutes":   intervalInMinutes,
				"workspace_id":          workspaceId,
				"workspace_region":      workspaceRegion,
				"workspace_resource_id": workspaceResourceId,
			})
		}
	}

	return output
}

func expandNetworkWatcherFlowLogTrafficAnalytics(d *pluginsdk.ResourceData) *flowlogs.TrafficAnalyticsProperties {
	vs := d.Get("traffic_analytics").([]interface{})

	v := vs[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	workspaceID := v["workspace_id"].(string)
	workspaceRegion := v["workspace_region"].(string)
	workspaceResourceID := v["workspace_resource_id"].(string)
	interval := v["interval_in_minutes"].(int)

	return &flowlogs.TrafficAnalyticsProperties{
		NetworkWatcherFlowAnalyticsConfiguration: &flowlogs.TrafficAnalyticsConfigurationProperties{
			Enabled:                  pointer.To(enabled),
			WorkspaceId:              utils.String(workspaceID),
			WorkspaceRegion:          utils.String(workspaceRegion),
			WorkspaceResourceId:      utils.String(workspaceResourceID),
			TrafficAnalyticsInterval: pointer.To(int64(interval)),
		},
	}
}
