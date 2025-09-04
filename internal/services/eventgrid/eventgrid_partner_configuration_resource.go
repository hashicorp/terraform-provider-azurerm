// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerconfigurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = EventGridPartnerConfigurationResource{}

type EventGridPartnerConfigurationResource struct{}

type EventGridPartnerConfigurationResourceModel struct {
	ResourceGroup                      string                 `tfschema:"resource_group_name"`
	DefaultMaximumExpirationTimeInDays int64                  `tfschema:"default_maximum_expiration_time_in_days"`
	PartnerAuthorizations              []PartnerAuthorization `tfschema:"partner_authorization"`
	Tags                               map[string]string      `tfschema:"tags"`
}

type PartnerAuthorization struct {
	PartnerRegistrationId            string `tfschema:"partner_registration_id"`
	PartnerName                      string `tfschema:"partner_name"`
	AuthorizationExpirationTimeInUtc string `tfschema:"authorization_expiration_time_in_utc"`
}

func (EventGridPartnerConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"resource_group_name": commonschema.ResourceGroupName(),
		"default_maximum_expiration_time_in_days": &schema.Schema{
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 365),
			Default:      7,
		},
		"partner_authorization": &schema.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"partner_registration_id": &schema.Schema{
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsUUID,
					},
					"partner_name": &schema.Schema{
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"authorization_expiration_time_in_utc": &schema.Schema{
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
		"tags": tags.Schema(),
	}
}

func (EventGridPartnerConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (EventGridPartnerConfigurationResource) ModelObject() interface{} {
	return &EventGridPartnerConfigurationResourceModel{}
}

func (EventGridPartnerConfigurationResource) ResourceType() string {
	return "azurerm_eventgrid_partner_configuration"
}

func (r EventGridPartnerConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			subscriptionId := metadata.Client.Account.SubscriptionId

			var config EventGridPartnerConfigurationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Uses resource group ID as it has 1-1 mapping with resource group
			// See the SDK sample usage: https://github.com/hashicorp/go-azure-sdk/tree/main/resource-manager/eventgrid/2022-06-15/partnerconfigurations
			id := commonids.NewResourceGroupID(subscriptionId, config.ResourceGroup)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := partnerconfigurations.PartnerConfiguration{
				Location: pointer.To("global"),
				Properties: &partnerconfigurations.PartnerConfigurationProperties{
					PartnerAuthorization: &partnerconfigurations.PartnerAuthorization{
						DefaultMaximumExpirationTimeInDays: pointer.To(config.DefaultMaximumExpirationTimeInDays),
						AuthorizedPartnersList:             expandAuthorizedPartnersList(config.PartnerAuthorizations),
					},
				},
				Tags: pointer.To(config.Tags),
			}
			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r EventGridPartnerConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			id, err := commonids.ParseResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config EventGridPartnerConfigurationResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", *id)
			}

			payload := existing.Model

			if metadata.ResourceData.HasChange("default_maximum_expiration_time_in_days") {
				payload.Properties.PartnerAuthorization.DefaultMaximumExpirationTimeInDays = pointer.To(config.DefaultMaximumExpirationTimeInDays)
			}

			if metadata.ResourceData.HasChange("partner_authorization") {
				if payload.Properties.PartnerAuthorization == nil {
					payload.Properties.PartnerAuthorization = &partnerconfigurations.PartnerAuthorization{}
				}
				payload.Properties.PartnerAuthorization.AuthorizedPartnersList = expandAuthorizedPartnersList(config.PartnerAuthorizations)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r EventGridPartnerConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			id, err := commonids.ParseResourceGroupID(metadata.ResourceData.Id())
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

			state := EventGridPartnerConfigurationResourceModel{
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil && props.PartnerAuthorization != nil {
					state.DefaultMaximumExpirationTimeInDays = pointer.From(props.PartnerAuthorization.DefaultMaximumExpirationTimeInDays)
					state.PartnerAuthorizations = flattenAuthorizedPartnersList(props.PartnerAuthorization.AuthorizedPartnersList)
				}

				if model.Tags != nil {
					state.Tags = pointer.From(model.Tags)
				}
			}
			return metadata.Encode(&state)
		},
	}
}

func (r EventGridPartnerConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.EventGrid.PartnerConfigurations

			id, err := commonids.ParseResourceGroupID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (EventGridPartnerConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateResourceGroupID
}

func expandAuthorizedPartnersList(partnerAuthorization []PartnerAuthorization) *[]partnerconfigurations.Partner {
	if len(partnerAuthorization) == 0 {
		return nil
	}

	partners := []partnerconfigurations.Partner{}

	for _, partnerAuth := range partnerAuthorization {
		partners = append(partners, partnerconfigurations.Partner{
			PartnerName:                      pointer.To(partnerAuth.PartnerName),
			PartnerRegistrationImmutableId:   pointer.To(partnerAuth.PartnerRegistrationId),
			AuthorizationExpirationTimeInUtc: pointer.To(partnerAuth.AuthorizationExpirationTimeInUtc),
		})
	}

	return &partners
}

func flattenAuthorizedPartnersList(partner *[]partnerconfigurations.Partner) []PartnerAuthorization {
	partnerAuthorizations := make([]PartnerAuthorization, 0)

	if partner == nil {
		return partnerAuthorizations
	}

	for _, partnerAuth := range *partner {
		partnerAuthorizations = append(partnerAuthorizations, PartnerAuthorization{
			PartnerName:                      pointer.From(partnerAuth.PartnerName),
			PartnerRegistrationId:            pointer.From(partnerAuth.PartnerRegistrationImmutableId),
			AuthorizationExpirationTimeInUtc: pointer.From(partnerAuth.AuthorizationExpirationTimeInUtc),
		})
	}

	return partnerAuthorizations
}
