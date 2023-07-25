// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/clusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ClusterResource struct{}

type ClusterModel struct {
	Name              string                 `tfschema:"name"`
	ResourceGroup     string                 `tfschema:"resource_group_name"`
	Location          string                 `tfschema:"location"`
	StreamingCapacity int64                  `tfschema:"streaming_capacity"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

var (
	_ sdk.ResourceWithUpdate         = ClusterResource{}
	_ sdk.ResourceWithStateMigration = ClusterResource{}
)

func (r ClusterResource) ModelObject() interface{} {
	return &ClusterModel{}
}

func (r ClusterResource) ResourceType() string {
	return "azurerm_stream_analytics_cluster"
}

func (r ClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clusters.ValidateClusterID
}

func (r ClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"streaming_capacity": {
			Type:     pluginsdk.TypeInt,
			Required: true,
			ValidateFunc: validation.All(
				validation.IntBetween(36, 216),
				validation.IntDivisibleBy(36),
			),
		},

		"tags": commonschema.Tags(),
	}
}

func (r ClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ClusterModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.ClustersClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := clusters.NewClusterID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := clusters.Cluster{
				Name:     utils.String(model.Name),
				Location: utils.String(model.Location),
				Sku: &clusters.ClusterSku{
					Name:     utils.ToPtr(clusters.ClusterSkuNameDefault),
					Capacity: utils.ToPtr(model.StreamingCapacity),
				},
				Tags: tags.Expand(model.Tags),
			}

			var opts clusters.CreateOrUpdateOperationOptions
			if err := client.CreateOrUpdateThenPoll(ctx, id, props, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.ClustersClient
			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := ClusterModel{
				Name:              id.ClusterName,
				ResourceGroup:     id.ResourceGroupName,
				StreamingCapacity: *resp.Model.Sku.Capacity,
			}

			if model := resp.Model; model != nil {
				state.Location = *model.Location
				state.Tags = tags.Flatten(model.Tags)

				var capacity int64
				if v := model.Sku.Capacity; v != nil {
					capacity = *v
				}
				state.StreamingCapacity = capacity
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.ClustersClient
			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 90 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.ClustersClient

			var state ClusterModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("streaming_capacity") || metadata.ResourceData.HasChange("tags") {
				props := clusters.Cluster{
					Sku: &clusters.ClusterSku{
						Capacity: utils.ToPtr(state.StreamingCapacity),
					},
					Tags: tags.Expand(state.Tags),
				}

				var opts clusters.UpdateOperationOptions
				if err := client.UpdateThenPoll(ctx, *id, props, opts); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}
			return nil
		},
	}
}

func (r ClusterResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsClusterV0ToV1{},
		},
	}
}
