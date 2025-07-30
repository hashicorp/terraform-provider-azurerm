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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnernamespaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = EventGridPartnerNamespaceDataSource{}

type EventGridPartnerNamespaceDataSource struct{}

type EventGridPartnerNamespaceDataSourceModel struct {
	Name                                string            `tfschema:"name"`
	ResourceGroup                       string            `tfschema:"resource_group_name"`
	Location                            string            `tfschema:"location"`
	Endpoint                            string            `tfschema:"endpoint"`
	PartnerRegistrationFullyQualifiedID string            `tfschema:"partner_registration_id"`
	PartnerTopicRoutingMode             string            `tfschema:"partner_topic_routing_mode"`
	Tags                                map[string]string `tfschema:"tags"`
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

		"partner_registration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"partner_topic_routing_mode": {
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

			id := partnernamespaces.NewPartnerNamespaceID(subscriptionId, state.ResourceGroup, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				if props := model.Properties; props != nil {
					state.PartnerRegistrationFullyQualifiedID = pointer.From(props.PartnerRegistrationFullyQualifiedId)
					state.PartnerTopicRoutingMode = string(pointer.From(props.PartnerTopicRoutingMode))
					state.Endpoint = pointer.From(props.Endpoint)
				}
			}

			metadata.SetID(id)

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
