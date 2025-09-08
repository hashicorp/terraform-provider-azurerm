// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2025-02-15/channels"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = EventGridPartnerNamespaceChannelDataSource{}

type EventGridPartnerNamespaceChannelDataSource struct{}

type EventGridPartnerNamespaceChannelDataSourceModel struct {
	ChannelName                       string              `tfschema:"name"`
	PartnerNamespaceName              string              `tfschema:"partner_namespace_name"`
	ResourceGroupName                 string              `tfschema:"resource_group_name"`
	ChannelType                       string              `tfschema:"channel_type"`
	ExpirationTimeIfNotActivatedInUtc string              `tfschema:"expiration_time_if_not_activated_in_utc"`
	PartnerTopic                      []PartnerTopicModel `tfschema:"partner_topic"`
	ReadinessState                    string              `tfschema:"readiness_state"`
}

func (r EventGridPartnerNamespaceChannelDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"partner_namespace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r EventGridPartnerNamespaceChannelDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"channel_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"expiration_time_if_not_activated_in_utc": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"partner_topic": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"subscription_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"resource_group_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"source": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"event_type_definitions": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"inline_event_type": {
									Type:     pluginsdk.TypeSet,
									Computed: true,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"display_name": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"data_schema_url": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"description": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
											"documentation_url": {
												Type:     pluginsdk.TypeString,
												Computed: true,
											},
										},
									},
								},
								"kind": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
		"readiness_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r EventGridPartnerNamespaceChannelDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.Channels

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state EventGridPartnerNamespaceChannelDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := channels.NewChannelID(subscriptionId, state.ResourceGroupName, state.PartnerNamespaceName, state.ChannelName)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.ChannelType = pointer.FromEnum(props.ChannelType)
					state.ExpirationTimeIfNotActivatedInUtc = pointer.From(props.ExpirationTimeIfNotActivatedUtc)
					state.ReadinessState = pointer.FromEnum(props.ReadinessState)

					if partnerTopicInfo := props.PartnerTopicInfo; partnerTopicInfo != nil {
						state.PartnerTopic = flattenPartnerNamespaceChannelPartnerTopic(partnerTopicInfo)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (EventGridPartnerNamespaceChannelDataSource) ResourceType() string {
	return "azurerm_eventgrid_partner_namespace_channel"
}

func (EventGridPartnerNamespaceChannelDataSource) ModelObject() interface{} {
	return &EventGridPartnerNamespaceChannelDataSourceModel{}
}
