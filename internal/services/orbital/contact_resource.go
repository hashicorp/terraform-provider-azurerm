// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package orbital

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contact"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/contactprofile"
	"github.com/hashicorp/go-azure-sdk/resource-manager/orbital/2022-11-01/spacecraft"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var _ sdk.ResourceWithDeprecationAndNoReplacement = ContactResource{}

func (r ContactResource) DeprecationMessage() string {
	return "The `azurerm_orbital_contact` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider."
}

type ContactResource struct{}

type ContactResourceModel struct {
	Name                 string `tfschema:"name"`
	Spacecraft           string `tfschema:"spacecraft_id"`
	ReservationStartTime string `tfschema:"reservation_start_time"`
	ReservationEndTime   string `tfschema:"reservation_end_time"`
	GroundStationName    string `tfschema:"ground_station_name"`
	ContactProfileId     string `tfschema:"contact_profile_id"`
}

func (r ContactResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spacecraft_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: spacecraft.ValidateSpacecraftID,
		},

		"reservation_start_time": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"reservation_end_time": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"ground_station_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"contact_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: contactprofile.ValidateContactProfileID,
		},
	}
}

func (r ContactResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r ContactResource) ModelObject() interface{} {
	return &ContactResourceModel{}
}

func (r ContactResource) ResourceType() string {
	return "azurerm_orbital_contact"
}

func (r ContactResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ContactResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Orbital.ContactClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			spacecraftId, err := contact.ParseSpacecraftID(model.Spacecraft)
			if err != nil {
				return err
			}
			id := contact.NewContactID(subscriptionId, spacecraftId.ResourceGroupName, spacecraftId.SpacecraftName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			contactProfile := contact.ResourceReference{
				Id: model.ContactProfileId,
			}

			contactProperties := contact.ContactsProperties{
				ContactProfile:       contactProfile,
				GroundStationName:    model.GroundStationName,
				ReservationEndTime:   model.ReservationEndTime,
				ReservationStartTime: model.ReservationStartTime,
			}

			contact := contact.Contact{
				Id:         utils.String(id.ID()),
				Name:       utils.String(model.Name),
				Properties: contactProperties,
			}
			if _, err = client.Create(ctx, id, contact); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)
			return nil
		},
	}
}

func (r ContactResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.ContactClient
			id, err := contact.ParseContactID(metadata.ResourceData.Id())
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

			spacecraftId := contact.NewSpacecraftID(id.SubscriptionId, id.ResourceGroupName, id.SpacecraftName)
			if model := resp.Model; model != nil {
				props := model.Properties
				state := ContactResourceModel{
					Name:       id.ContactName,
					Spacecraft: spacecraftId.ID(),
				}

				state.ReservationStartTime = props.ReservationStartTime
				state.ReservationEndTime = props.ReservationEndTime
				state.GroundStationName = props.GroundStationName
				state.ContactProfileId = props.ContactProfile.Id

				return metadata.Encode(&state)
			}
			return nil
		},
	}
}

func (r ContactResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Orbital.ContactClient
			id, err := contact.ParseContactID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s", *id)

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r ContactResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return contact.ValidateContactID
}
