// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagedPrivateEndpointResource struct{}

type ManagedPrivateEndpointModel struct {
	Name                   string `tfschema:"name"`
	ResourceGroup          string `tfschema:"resource_group_name"`
	StreamAnalyticsCluster string `tfschema:"stream_analytics_cluster_name"`
	TargetResourceId       string `tfschema:"target_resource_id"`
	SubResourceName        string `tfschema:"subresource_name"`
}

var _ sdk.ResourceWithStateMigration = ManagedPrivateEndpointResource{}

func (r ManagedPrivateEndpointResource) ModelObject() interface{} {
	return &ManagedPrivateEndpointModel{}
}

func (r ManagedPrivateEndpointResource) ResourceType() string {
	return "azurerm_stream_analytics_managed_private_endpoint"
}

func (r ManagedPrivateEndpointResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return privateendpoints.ValidatePrivateEndpointID
}

func (r ManagedPrivateEndpointResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"stream_analytics_cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"subresource_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ManagedPrivateEndpointResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagedPrivateEndpointResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagedPrivateEndpointModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.StreamAnalytics.EndpointsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := privateendpoints.NewPrivateEndpointID(subscriptionId, model.ResourceGroup, model.StreamAnalyticsCluster, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			props := privateendpoints.PrivateEndpoint{
				Properties: &privateendpoints.PrivateEndpointProperties{
					ManualPrivateLinkServiceConnections: &[]privateendpoints.PrivateLinkServiceConnection{
						{
							Properties: &privateendpoints.PrivateLinkServiceConnectionProperties{
								PrivateLinkServiceId: utils.String(model.TargetResourceId),
								GroupIds:             &[]string{model.SubResourceName},
							},
						},
					},
				},
			}

			var opts privateendpoints.CreateOrUpdateOperationOptions
			if _, err := client.CreateOrUpdate(ctx, id, props, opts); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagedPrivateEndpointResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.EndpointsClient
			id, err := privateendpoints.ParsePrivateEndpointID(metadata.ResourceData.Id())
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

			if resp.Model.Properties.ManualPrivateLinkServiceConnections == nil {
				return fmt.Errorf("no private link service connections available")
			}

			state := ManagedPrivateEndpointModel{
				Name:                   id.PrivateEndpointName,
				ResourceGroup:          id.ResourceGroupName,
				StreamAnalyticsCluster: id.ClusterName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					for _, mplsc := range *props.ManualPrivateLinkServiceConnections {
						state.TargetResourceId = *mplsc.Properties.PrivateLinkServiceId
						state.SubResourceName = strings.Join(*mplsc.Properties.GroupIds, "")
					}
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r ManagedPrivateEndpointResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.EndpointsClient
			id, err := privateendpoints.ParsePrivateEndpointID(metadata.ResourceData.Id())
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

func (r ManagedPrivateEndpointResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.StreamAnalyticsManagedPrivateEndpointV0ToV1{},
		},
	}
}
