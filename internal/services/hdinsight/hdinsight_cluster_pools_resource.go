// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package hdinsight

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2024-05-01/hdinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"
)

type ClusterPoolResource struct{}

type ClusterPoolModel struct {
	Name                     string                 `tfschema:"name"`
	ResourceGroup            string                 `tfschema:"resource_group_name"`
	Location                 string                 `tfschema:"location"`
	ManagedResourceGroupName string                 `tfschema:"managed_resource_group_name"`
	ClusterPoolProfile       []ClusterPoolProfile   `tfschema:"cluster_pool_profile"`
	ComputeProfile           []ComputeProfile       `tfschema:"compute_profile"`
	LogAnalyticsProfile      []LogAnalyticsProfile  `tfschema:"log_analytics_profile"`
	NetworkProfile           []NetworkProfile       `tfschema:"network_profile"`
	Tags                     map[string]interface{} `tfschema:"tags"`
}

type ClusterPoolProfile struct {
	ClusterPoolVersion string `tfschema:"cluster_pool_version"`
}

type ComputeProfile struct {
	VmSize string `tfschema:"vm_size"`
}

type LogAnalyticsProfile struct {
	LogAnalyticsProfileEnabled bool   `tfschema:"log_analytics_profile_enabled"`
	WorkspaceId                string `tfschema:"workspace_id"`
}

type NetworkProfile struct {
	PrivateApiServerEnabled bool   `tfschema:"private_api_server_enabled"`
	OutboundType            string `tfschema:"outbound_type"`
	SubnetId                string `tfschema:"subnet_id"`
}

var _ sdk.ResourceWithUpdate = ClusterPoolResource{}

func (r ClusterPoolResource) ModelObject() interface{} {
	return &ClusterPoolModel{}
}

func (r ClusterPoolResource) ResourceType() string {
	return "azurerm_hdinsight_cluster_pools"
}

func (r ClusterPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return hdinsights.ValidateClusterPoolID
}

func (r ClusterPoolResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"compute_profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vm_size": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"managed_resource_group_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cluster_pool_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cluster_pool_version": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"log_analytics_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"log_analytics_profile_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"workspace_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"network_profile": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"private_api_server_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"outbound_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r ClusterPoolResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ClusterPoolResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ClusterPoolModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.HDInsight2024.Hdinsights
			subscriptionId := metadata.Client.Account.SubscriptionId

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := hdinsights.NewClusterPoolID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.ClusterPoolsGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			clusterPool := hdinsights.ClusterPool{
				Location: location.Normalize(model.Location),
				Properties: &hdinsights.ClusterPoolResourceProperties{
					ClusterPoolProfile:       expandClusterPoolProfile(model.ClusterPoolProfile),
					ComputeProfile:           expandComputeProfile(model.ComputeProfile),
					LogAnalyticsProfile:      expandLogAnalyticsProfile(model.LogAnalyticsProfile),
					ManagedResourceGroupName: pointer.To(model.ManagedResourceGroupName),
					NetworkProfile:           expandNetworkProfile(model.NetworkProfile),
				},
				Tags: tags.Expand(model.Tags),
			}

			if _, err := client.ClusterPoolsCreateOrUpdate(ctx, id, clusterPool); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ClusterPoolResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HDInsight2024.Hdinsights

			id, err := hdinsights.ParseClusterPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state ClusterPoolModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding state: %+v", err)
			}

			existing, err := client.ClusterPoolsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			model := existing.Model

			if model.Properties == nil {
				return fmt.Errorf("retreiving properties for %s for update: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("cluster_pool_profile") {
				existing.Model.Properties.ClusterPoolProfile = expandClusterPoolProfile(state.ClusterPoolProfile)
			}

			if metadata.ResourceData.HasChange("compute_profile") {
				existing.Model.Properties.ComputeProfile = expandComputeProfile(state.ComputeProfile)
			}

			if metadata.ResourceData.HasChange("log_analytics_profile") {
				existing.Model.Properties.LogAnalyticsProfile = expandLogAnalyticsProfile(state.LogAnalyticsProfile)
			}

			if metadata.ResourceData.HasChange("network_profile") {
				existing.Model.Properties.NetworkProfile = expandNetworkProfile(state.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = tags.Expand(state.Tags)
			}

			if metadata.ResourceData.HasChange("managed_resource_group_name") {
				existing.Model.Properties.ManagedResourceGroupName = pointer.To(state.ManagedResourceGroupName)
			}

			if err := client.ClusterPoolsCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ClusterPoolResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HDInsight2024.Hdinsights

			id, err := hdinsights.ParseClusterPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.ClusterPoolsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			var state ClusterPoolModel

			state.Name = id.ClusterPoolName
			state.ResourceGroup = id.ResourceGroupName

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = tags.Flatten(model.Tags)

				if props := model.Properties; props != nil {
					state.ManagedResourceGroupName = pointer.From(props.ManagedResourceGroupName)
					state.ClusterPoolProfile = flattenClusterPoolProfile(props.ClusterPoolProfile)
					state.ComputeProfile = flattenComputeProfile(props.ComputeProfile)
					state.LogAnalyticsProfile = flattenLogAnalyticsProfile(props.LogAnalyticsProfile)
					state.NetworkProfile = flattenNetworkProfile(props.NetworkProfile)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r ClusterPoolResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HDInsight2024.Hdinsights

			id, err := hdinsights.ParseClusterPoolID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.ClusterPoolsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandClusterPoolProfile(profiles []ClusterPoolProfile) *hdinsights.ClusterPoolProfile {
	if len(profiles) == 0 {
		return nil
	}

	result := hdinsights.ClusterPoolProfile{
		ClusterPoolVersion: profiles[0].ClusterPoolVersion,
	}
	return &result
}

func expandComputeProfile(profiles []ComputeProfile) hdinsights.ClusterPoolComputeProfile {
	result := hdinsights.ClusterPoolComputeProfile{
		VMSize: profiles[0].VmSize,
	}
	return result
}

func expandLogAnalyticsProfile(profiles []LogAnalyticsProfile) *hdinsights.ClusterPoolLogAnalyticsProfile {
	if len(profiles) == 0 {
		return nil
	}

	result := hdinsights.ClusterPoolLogAnalyticsProfile{
		Enabled:     profiles[0].LogAnalyticsProfileEnabled,
		WorkspaceId: pointer.To(profiles[0].WorkspaceId),
	}
	return &result
}

func expandNetworkProfile(profiles []NetworkProfile) *hdinsights.ClusterPoolNetworkProfile {
	if len(profiles) == 0 {
		return nil
	}

	result := hdinsights.ClusterPoolNetworkProfile{
		SubnetId: profiles[0].SubnetId,
	}

	if profiles[0].PrivateApiServerEnabled {
		result.EnablePrivateApiServer = pointer.To(true)
	}

	if profiles[0].OutboundType != "" {
		result.OutboundType = pointer.To(hdinsights.OutboundType(profiles[0].OutboundType))
	}
	return &result
}

func flattenClusterPoolProfile(input *hdinsights.ClusterPoolProfile) []ClusterPoolProfile {
	result := make([]ClusterPoolProfile, 0)
	if input == nil {
		return result
	}

	profile := ClusterPoolProfile{
		ClusterPoolVersion: input.ClusterPoolVersion,
	}
	result = append(result, profile)

	return result
}

func flattenComputeProfile(input hdinsights.ClusterPoolComputeProfile) []ComputeProfile {
	result := make([]ComputeProfile, 0)
	profile := ComputeProfile{
		VmSize: input.VMSize,
	}
	result = append(result, profile)

	return result
}

func flattenLogAnalyticsProfile(input *hdinsights.ClusterPoolLogAnalyticsProfile) []LogAnalyticsProfile {
	result := make([]LogAnalyticsProfile, 0)
	if input == nil {
		return result
	}

	profile := LogAnalyticsProfile{
		LogAnalyticsProfileEnabled: input.Enabled,
		WorkspaceId:                pointer.From(input.WorkspaceId),
	}
	result = append(result, profile)

	return result
}

func flattenNetworkProfile(input *hdinsights.ClusterPoolNetworkProfile) []NetworkProfile {
	result := make([]NetworkProfile, 0)
	if input == nil {
		return result
	}

	profile := NetworkProfile{
		SubnetId: input.SubnetId,
	}

	if input.EnablePrivateApiServer != nil {
		profile.PrivateApiServerEnabled = *input.EnablePrivateApiServer
	}

	if input.OutboundType != nil {
		profile.OutboundType = string(*input.OutboundType)
	}
	result = append(result, profile)

	return result
}
