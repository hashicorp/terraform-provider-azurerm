// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = EventGridPartnerNamespaceDataSource{}

type EventGridPartnerNamespaceDataSource struct{}

type EventGridPartnerNamespaceDataSourceModel struct {
	PartnerNamespaceName                string                      `tfschema:"name"`
	ResourceGroup                       string                      `tfschema:"resource_group_name"`
	Endpoint                            string                      `tfschema:"endpoint"`
	Location                            string                      `tfschema:"location"`
	InboundIPRules                      []PartnerInboundIpRuleModel `tfschema:"inbound_ip_rule"`
	LocalAuthEnabled                    bool                        `tfschema:"local_authentication_enabled"`
	PartnerRegistrationFullyQualifiedID string                      `tfschema:"partner_registration_id"`
	PartnerTopicRoutingMode             string                      `tfschema:"partner_topic_routing_mode"`
	PublicNetworkAccess                 string                      `tfschema:"public_network_access"`
	Tags                                map[string]string           `tfschema:"tags"`
}

func (r EventGridPartnerNamespaceDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r EventGridPartnerNamespaceDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"endpoint": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"inbound_ip_rule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"ip_mask": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"action": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"partner_registration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"partner_topic_routing_mode": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_network_access": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r EventGridPartnerNamespaceDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerNamespaces

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state EventGridPartnerNamespaceDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := partnernamespaces.NewPartnerNamespaceID(subscriptionId, state.ResourceGroup, state.PartnerNamespaceName)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.PartnerRegistrationFullyQualifiedID = pointer.From(props.PartnerRegistrationFullyQualifiedId)
					state.PartnerTopicRoutingMode = pointer.FromEnum(props.PartnerTopicRoutingMode)
					state.Endpoint = pointer.From(props.Endpoint)
					state.InboundIPRules = flattenPartnerInboundIPRules(props.InboundIPRules)
					state.LocalAuthEnabled = !pointer.From(props.DisableLocalAuth)
					state.PublicNetworkAccess = pointer.FromEnum(props.PublicNetworkAccess)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (EventGridPartnerNamespaceDataSource) ResourceType() string {
	return "azurerm_eventgrid_partner_namespace"
}

func (EventGridPartnerNamespaceDataSource) ModelObject() interface{} {
	return &EventGridPartnerNamespaceDataSourceModel{}
}
