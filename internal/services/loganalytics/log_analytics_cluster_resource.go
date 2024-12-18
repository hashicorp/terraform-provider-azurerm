// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LogAnalyticsClusterResource struct{}

type LogAnalyticsClusterModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	SizeGB            int64                                      `tfschema:"size_gb"`
	Tags              map[string]string                          `tfschema:"tags"`
	ClusterId         string                                     `tfschema:"cluster_id"`
}

var _ sdk.ResourceWithUpdate = LogAnalyticsClusterResource{}

func (l LogAnalyticsClusterResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LogAnalyticsClusterName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityRequiredForceNew(),

		"size_gb": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default: func() int {
				if !features.FourPointOh() {
					return 1000
				}
				return 100
			}(),
			ValidateFunc: validation.IntInSlice([]int{100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 25000, 50000}),
		},

		"tags": tags.Schema(),
	}
}

func (l LogAnalyticsClusterResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cluster_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r LogAnalyticsClusterResource) ModelObject() interface{} {
	return &LogAnalyticsClusterModel{}
}

func (r LogAnalyticsClusterResource) ResourceType() string {
	return "azurerm_log_analytics_cluster"
}

func (r LogAnalyticsClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clusters.ValidateClusterID
}

func (r LogAnalyticsClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.ClusterClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config LogAnalyticsClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := clusters.NewClusterID(subscriptionId, config.ResourceGroupName, config.Name)

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(config.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			capacityReservation := clusters.ClusterSkuNameEnumCapacityReservation
			parameters := clusters.Cluster{
				Location: location.Normalize(config.Location),
				Identity: expandedIdentity,
				Sku: &clusters.ClusterSku{
					Capacity: pointer.To(clusters.Capacity(config.SizeGB)),
					Name:     &capacityReservation,
				},
				Tags: pointer.To(config.Tags),
			}

			err = client.CreateOrUpdateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LogAnalyticsClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.ClusterClient

			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := LogAnalyticsClusterModel{
				Name:              id.ClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.NormalizeNilable(&model.Location)

				flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = flattenedIdentity

				if props := model.Properties; props != nil {
					state.ClusterId = pointer.From(props.ClusterId)
				}

				capacity := 0
				if sku := model.Sku; sku != nil {
					if sku.Capacity != nil {
						capacity = int(*sku.Capacity)
					}
				}
				state.SizeGB = int64(capacity)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogAnalyticsClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 6 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.ClusterClient

			var config LogAnalyticsClusterModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			payload := clusters.ClusterPatch{}

			if metadata.ResourceData.HasChange("size_gb") {
				payload.Sku = &clusters.ClusterSku{
					Capacity: pointer.To(clusters.Capacity(config.SizeGB)),
					Name:     pointer.To(clusters.ClusterSkuNameEnumCapacityReservation),
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err = client.UpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LogAnalyticsClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.ClusterClient

			id, err := clusters.ParseClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
