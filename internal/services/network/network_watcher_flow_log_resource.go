// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/flowlogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networkwatchers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkWatcherFlowLog() *pluginsdk.Resource {
	return &pluginsdk.Resource{
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

			"target_resource_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					networksecuritygroups.ValidateNetworkSecurityGroupID,
					commonids.ValidateVirtualNetworkID,
					commonids.ValidateSubnetID,
					commonids.ValidateNetworkInterfaceID,
				),
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

		CustomizeDiff: func(_ context.Context, d *pluginsdk.ResourceDiff, _ any) error {
			if d.Id() == "" {
				targetResourceId := d.Get("target_resource_id").(string)

				if _, err := networksecuritygroups.ParseNetworkSecurityGroupID(targetResourceId); err == nil {
					return errors.New("creation of new NSG flow logs is no longer supported by Azure as of June 30, 2025. NSG flow logs will be retired on September 30, 2027. For more information, see https://learn.microsoft.com/azure/network-watcher/nsg-flow-logs-migrate")
				}
			}

			return nil
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

func resourceNetworkWatcherFlowLogCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.FlowLogs
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := flowlogs.NewFlowLogID(subscriptionId, d.Get("resource_group_name").(string), d.Get("network_watcher_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_watcher_flow_log", id.ID())
	}

	targetResourceId := d.Get("target_resource_id").(string)

	locks.ByID(targetResourceId)
	defer locks.UnlockByID(targetResourceId)

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
		Location: pointer.To(location.Normalize(loc)),
		Properties: &flowlogs.FlowLogPropertiesFormat{
			TargetResourceId: targetResourceId,
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
		parameters.Properties.Format = &flowlogs.FlowLogFormatParameters{
			Version: pointer.To(int64(version.(int))),
		}
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

	targetResourceId := d.Get("target_resource_id").(string)

	locks.ByID(targetResourceId)
	defer locks.UnlockByID(targetResourceId)

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

	if d.HasChange("target_resource_id") {
		payload.Properties.TargetResourceId = targetResourceId
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

			targetResourceId := props.TargetResourceId
			if nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = nsgId.ID()
			} else if vnetId, err := commonids.ParseVirtualNetworkIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = vnetId.ID()
			} else if subnetId, err := commonids.ParseSubnetIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = subnetId.ID()
			} else if nicId, err := commonids.ParseNetworkInterfaceIDInsensitively(props.TargetResourceId); err == nil {
				targetResourceId = nicId.ID()
			}

			d.Set("target_resource_id", targetResourceId)

			if err := d.Set("retention_policy", flattenNetworkWatcherFlowLogRetentionPolicy(props.RetentionPolicy)); err != nil {
				return fmt.Errorf("setting `retention_policy`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
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
		return fmt.Errorf("retrieving %s: `properties` or `properties.TargetResourceID` was nil", id)
	}

	targetResourceId := resp.Model.Properties.TargetResourceId
	if nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupIDInsensitively(resp.Model.Properties.TargetResourceId); err == nil {
		targetResourceId = nsgId.ID()
	} else if vnetId, err := commonids.ParseVirtualNetworkIDInsensitively(resp.Model.Properties.TargetResourceId); err == nil {
		targetResourceId = vnetId.ID()
	}

	locks.ByID(targetResourceId)
	defer locks.UnlockByID(targetResourceId)

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
		days := 0
		if input.Days != nil {
			days = int(*input.Days)
		}
		output = append(output, map[string]interface{}{
			"days":    days,
			"enabled": pointer.From(input.Enabled),
		})
	}

	return output
}

func flattenNetworkWatcherFlowLogTrafficAnalytics(input *flowlogs.TrafficAnalyticsProperties) []interface{} {
	output := make([]interface{}, 0)
	if input != nil {
		if cfg := input.NetworkWatcherFlowAnalyticsConfiguration; cfg != nil {
			intervalInMinutes := 0
			if cfg.TrafficAnalyticsInterval != nil {
				intervalInMinutes = int(*cfg.TrafficAnalyticsInterval)
			}
			output = append(output, map[string]interface{}{
				"enabled":               pointer.From(cfg.Enabled),
				"interval_in_minutes":   intervalInMinutes,
				"workspace_id":          pointer.From(cfg.WorkspaceId),
				"workspace_region":      pointer.From(cfg.WorkspaceRegion),
				"workspace_resource_id": pointer.From(cfg.WorkspaceResourceId),
			})
		}
	}

	return output
}

func expandNetworkWatcherFlowLogTrafficAnalytics(d *pluginsdk.ResourceData) *flowlogs.TrafficAnalyticsProperties {
	vs := d.Get("traffic_analytics").([]interface{})

	if len(vs) == 0 {
		return nil
	}

	v := vs[0].(map[string]interface{})
	enabled := v["enabled"].(bool)
	workspaceID := v["workspace_id"].(string)
	workspaceRegion := v["workspace_region"].(string)
	workspaceResourceID := v["workspace_resource_id"].(string)
	interval := v["interval_in_minutes"].(int)

	return &flowlogs.TrafficAnalyticsProperties{
		NetworkWatcherFlowAnalyticsConfiguration: &flowlogs.TrafficAnalyticsConfigurationProperties{
			Enabled:                  pointer.To(enabled),
			WorkspaceId:              pointer.To(workspaceID),
			WorkspaceRegion:          pointer.To(workspaceRegion),
			WorkspaceResourceId:      pointer.To(workspaceResourceID),
			TrafficAnalyticsInterval: pointer.To(int64(interval)),
		},
	}
}
