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
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.DataSource = EventGridPartnerRegistrationDataSource{}

type EventGridPartnerRegistrationDataSource struct{}

type EventGridPartnerRegistrationDataSourceModel struct {
	Name                  string            `tfschema:"name"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	PartnerRegistrationID string            `tfschema:"partner_registration_id"`
	Tags                  map[string]string `tfschema:"tags"`
}

func (r EventGridPartnerRegistrationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r EventGridPartnerRegistrationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"partner_registration_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"tags": commonschema.TagsDataSource(),
	}
}

func (r EventGridPartnerRegistrationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerRegistrations

			subscriptionId := metadata.Client.Account.SubscriptionId

			var state EventGridPartnerRegistrationDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := partnerregistrations.NewPartnerRegistrationID(subscriptionId, state.ResourceGroup, state.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				state.Tags = pointer.From(model.Tags)
				if props := model.Properties; props != nil {
					state.PartnerRegistrationID = pointer.From(props.PartnerRegistrationImmutableId)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

func (EventGridPartnerRegistrationDataSource) ResourceType() string {
	return "azurerm_eventgrid_partner_registration"
}

func (EventGridPartnerRegistrationDataSource) ModelObject() interface{} {
	return &EventGridPartnerRegistrationDataSourceModel{}
}
