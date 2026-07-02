// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	resourcesSubscription "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = SubscriptionZonePeersDataSource{}

type SubscriptionZonePeersDataSource struct{}

type SubscriptionZonePeersDataSourceModel struct {
	Location              string                             `tfschema:"location"`
	PeerSubscriptionID    string                             `tfschema:"peer_subscription_id"`
	SubscriptionID        string                             `tfschema:"subscription_id"`
	AvailabilityZonePeers []SubscriptionZonePeerMappingModel `tfschema:"availability_zone_peers"`
}

type SubscriptionZonePeerMappingModel struct {
	AvailabilityZone string                               `tfschema:"availability_zone"`
	Peers            []SubscriptionZonePeerPeerEntryModel `tfschema:"peers"`
}

type SubscriptionZonePeerPeerEntryModel struct {
	SubscriptionID   string `tfschema:"subscription_id"`
	AvailabilityZone string `tfschema:"availability_zone"`
}

func (d SubscriptionZonePeersDataSource) ResourceType() string {
	return "azurerm_subscription_zone_peers"
}

func (d SubscriptionZonePeersDataSource) ModelObject() interface{} {
	return &SubscriptionZonePeersDataSourceModel{}
}

func (d SubscriptionZonePeersDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationWithoutForceNew(),

		"peer_subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (d SubscriptionZonePeersDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"subscription_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"availability_zone_peers": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"availability_zone": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"peers": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"subscription_id": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"availability_zone": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d SubscriptionZonePeersDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Subscription.SubscriptionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model SubscriptionZonePeersDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			normalizedLocation := location.Normalize(model.Location)

			req := resourcesSubscription.CheckZonePeersRequest{
				Location:        &normalizedLocation,
				SubscriptionIds: &[]string{fmt.Sprintf("subscriptions/%s", model.PeerSubscriptionID)},
			}

			id := commonids.NewSubscriptionID(subscriptionId)
			resp, err := client.CheckZonePeers(ctx, id, req)
			if err != nil {
				return fmt.Errorf("retrieving zone peers for %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving zone peers for %s: model was nil", id)
			}

			state := SubscriptionZonePeersDataSourceModel{
				Location:              normalizedLocation,
				PeerSubscriptionID:    model.PeerSubscriptionID,
				SubscriptionID:        pointer.From(resp.Model.SubscriptionId),
				AvailabilityZonePeers: flattenAvailabilityZonePeers(resp.Model.AvailabilityZonePeers),
			}

			metadata.ResourceData.SetId(id.ID())

			return metadata.Encode(&state)
		},
	}
}

func flattenAvailabilityZonePeers(input *[]resourcesSubscription.AvailabilityZonePeers) []SubscriptionZonePeerMappingModel {
	results := make([]SubscriptionZonePeerMappingModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		mapping := SubscriptionZonePeerMappingModel{
			AvailabilityZone: pointer.From(item.AvailabilityZone),
			Peers:            flattenPeers(item.Peers),
		}
		results = append(results, mapping)
	}

	return results
}

func flattenPeers(input *[]resourcesSubscription.Peers) []SubscriptionZonePeerPeerEntryModel {
	results := make([]SubscriptionZonePeerPeerEntryModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, SubscriptionZonePeerPeerEntryModel{
			SubscriptionID:   pointer.From(item.SubscriptionId),
			AvailabilityZone: pointer.From(item.AvailabilityZone),
		})
	}

	return results
}
