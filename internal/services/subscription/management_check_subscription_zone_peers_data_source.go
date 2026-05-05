// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subscription

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	resourcesSubscription "github.com/hashicorp/go-azure-sdk/resource-manager/resources/2022-12-01/subscriptions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = ManagementCheckSubscriptionZonePeersDataSource{}

type ManagementCheckSubscriptionZonePeersDataSource struct{}

type ManagementCheckSubscriptionZonePeersDataSourceModel struct {
	Location              string                                            `tfschema:"location"`
	PeerSubscriptionID    string                                            `tfschema:"peer_subscription_id"`
	SubscriptionID        string                                            `tfschema:"subscription_id"`
	AvailabilityZonePeers []ManagementCheckSubscriptionZonePeerMappingModel `tfschema:"availability_zone_peers"`
}

type ManagementCheckSubscriptionZonePeerMappingModel struct {
	AvailabilityZone string                                              `tfschema:"availability_zone"`
	Peers            []ManagementCheckSubscriptionZonePeerPeerEntryModel `tfschema:"peers"`
}

type ManagementCheckSubscriptionZonePeerPeerEntryModel struct {
	SubscriptionID   string `tfschema:"subscription_id"`
	AvailabilityZone string `tfschema:"availability_zone"`
}

func (d ManagementCheckSubscriptionZonePeersDataSource) ResourceType() string {
	return "azurerm_management_check_subscription_zone_peers"
}

func (d ManagementCheckSubscriptionZonePeersDataSource) ModelObject() interface{} {
	return &ManagementCheckSubscriptionZonePeersDataSourceModel{}
}

func (d ManagementCheckSubscriptionZonePeersDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationWithoutForceNew(),

		"peer_subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.Any(validation.IsUUID, commonids.ValidateSubscriptionID),
		},
	}
}

func (d ManagementCheckSubscriptionZonePeersDataSource) Attributes() map[string]*pluginsdk.Schema {
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

func (d ManagementCheckSubscriptionZonePeersDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Subscription.SubscriptionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagementCheckSubscriptionZonePeersDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			normalizedLocation := location.Normalize(model.Location)

			peerGuid, err := normalizePeerSubscriptionID(model.PeerSubscriptionID)
			if err != nil {
				return err
			}

			subscriptionIds := []string{fmt.Sprintf("subscriptions/%s", peerGuid)}
			req := resourcesSubscription.CheckZonePeersRequest{
				Location:        &normalizedLocation,
				SubscriptionIds: &subscriptionIds,
			}

			id := commonids.NewSubscriptionID(subscriptionId)
			resp, err := client.CheckZonePeers(ctx, id, req)
			if err != nil {
				return fmt.Errorf("retrieving zone peers for %s: %+v", id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving zone peers for %s: model was nil", id)
			}

			state := ManagementCheckSubscriptionZonePeersDataSourceModel{
				Location:              normalizedLocation,
				PeerSubscriptionID:    model.PeerSubscriptionID,
				SubscriptionID:        pointer.From(resp.Model.SubscriptionId),
				AvailabilityZonePeers: flattenAvailabilityZonePeers(resp.Model.AvailabilityZonePeers),
			}

			metadata.ResourceData.SetId(fmt.Sprintf("%s/checkZonePeers/%s/%s", id.ID(), normalizedLocation, peerGuid))

			return metadata.Encode(&state)
		},
	}
}

func normalizePeerSubscriptionID(input string) (string, error) {
	s := strings.TrimSpace(input)

	if parsed, err := commonids.ParseSubscriptionIDInsensitively(s); err == nil {
		return strings.ToLower(parsed.SubscriptionId), nil
	}

	s = strings.ToLower(s)
	if s == "" {
		return "", fmt.Errorf("`peer_subscription_id` is empty")
	}

	return s, nil
}

func flattenAvailabilityZonePeers(input *[]resourcesSubscription.AvailabilityZonePeers) []ManagementCheckSubscriptionZonePeerMappingModel {
	results := make([]ManagementCheckSubscriptionZonePeerMappingModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		mapping := ManagementCheckSubscriptionZonePeerMappingModel{
			AvailabilityZone: pointer.From(item.AvailabilityZone),
			Peers:            flattenPeers(item.Peers),
		}
		results = append(results, mapping)
	}

	return results
}

func flattenPeers(input *[]resourcesSubscription.Peers) []ManagementCheckSubscriptionZonePeerPeerEntryModel {
	results := make([]ManagementCheckSubscriptionZonePeerPeerEntryModel, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		results = append(results, ManagementCheckSubscriptionZonePeerPeerEntryModel{
			SubscriptionID:   pointer.From(item.SubscriptionId),
			AvailabilityZone: pointer.From(item.AvailabilityZone),
		})
	}

	return results
}
