// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-03/galleryapplications"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type GalleryApplicationResource struct{}

var (
	_ sdk.ResourceWithUpdate        = GalleryApplicationResource{}
	_ sdk.ResourceWithCustomizeDiff = GalleryApplicationResource{}
)

type GalleryApplicationModel struct {
	Name                string            `tfschema:"name"`
	GalleryId           string            `tfschema:"gallery_id"`
	Location            string            `tfschema:"location"`
	SupportedOSType     string            `tfschema:"supported_os_type"`
	Description         string            `tfschema:"description"`
	EndOfLifeDate       string            `tfschema:"end_of_life_date"`
	Eula                string            `tfschema:"eula"`
	PrivacyStatementURI string            `tfschema:"privacy_statement_uri"`
	ReleaseNoteURI      string            `tfschema:"release_note_uri"`
	Tags                map[string]string `tfschema:"tags"`
}

func (r GalleryApplicationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.GalleryApplicationName,
		},

		"gallery_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSharedImageGalleryID,
		},

		"location": commonschema.Location(),

		"supported_os_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(galleryapplications.PossibleValuesForOperatingSystemTypes(), false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"end_of_life_date": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"eula": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"privacy_statement_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"release_note_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r GalleryApplicationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r GalleryApplicationResource) ResourceType() string {
	return "azurerm_gallery_application"
}

func (r GalleryApplicationResource) ModelObject() interface{} {
	return &GalleryApplicationModel{}
}

func (r GalleryApplicationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return galleryapplications.ValidateApplicationID
}

func (r GalleryApplicationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state GalleryApplicationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Compute.GalleryApplicationsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			galleryId, err := commonids.ParseSharedImageGalleryID(state.GalleryId)
			if err != nil {
				return err
			}

			id := galleryapplications.NewApplicationID(subscriptionId, galleryId.ResourceGroupName, galleryId.GalleryName, state.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := galleryapplications.GalleryApplication{
				Location: location.Normalize(state.Location),
				Properties: &galleryapplications.GalleryApplicationProperties{
					SupportedOSType: galleryapplications.OperatingSystemTypes(state.SupportedOSType),
				},
				Tags: pointer.To(state.Tags),
			}

			if state.Description != "" {
				payload.Properties.Description = utils.String(state.Description)
			}

			if state.EndOfLifeDate != "" {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				payload.Properties.SetEndOfLifeDateAsTime(endOfLifeDate)
			}

			if state.Eula != "" {
				payload.Properties.Eula = utils.String(state.Eula)
			}

			if state.PrivacyStatementURI != "" {
				payload.Properties.PrivacyStatementUri = utils.String(state.PrivacyStatementURI)
			}

			if state.ReleaseNoteURI != "" {
				payload.Properties.ReleaseNoteUri = utils.String(state.ReleaseNoteURI)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r GalleryApplicationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationsClient
			id, err := galleryapplications.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%q was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := &GalleryApplicationModel{
				Name:      id.ApplicationName,
				GalleryId: commonids.NewSharedImageGalleryID(id.SubscriptionId, id.ResourceGroupName, id.GalleryName).ID(),
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				if model.Tags != nil {
					state.Tags = *model.Tags
				}

				if props := model.Properties; props != nil {
					if v := props.Description; v != nil {
						state.Description = *props.Description
					}

					if v := props.EndOfLifeDate; v != nil {
						d, err := props.GetEndOfLifeDateAsTime()
						if err != nil {
							return fmt.Errorf("parsing `end_of_life_date` from API Response: %+v", err)
						}
						if d != nil {
							state.EndOfLifeDate = d.Format(time.RFC3339)
						}
					}

					if v := props.Eula; v != nil {
						state.Eula = *props.Eula
					}

					if v := props.PrivacyStatementUri; v != nil {
						state.PrivacyStatementURI = *props.PrivacyStatementUri
					}

					if v := props.ReleaseNoteUri; v != nil {
						state.ReleaseNoteURI = *props.ReleaseNoteUri
					}

					state.SupportedOSType = string(props.SupportedOSType)
				}
			}

			return metadata.Encode(state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r GalleryApplicationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationsClient

			id, err := galleryapplications.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state GalleryApplicationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			payload := *existing.Model
			if metadata.ResourceData.HasChanges("description", "end_of_life_date", "eula", "privacy_statement_uri", "release_note_uri") {
				if payload.Properties == nil {
					payload.Properties = &galleryapplications.GalleryApplicationProperties{}
				}

				if metadata.ResourceData.HasChange("description") {
					payload.Properties.Description = utils.String(state.Description)
				}

				if metadata.ResourceData.HasChange("end_of_life_date") {
					endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
					payload.Properties.SetEndOfLifeDateAsTime(endOfLifeDate)
				}

				if metadata.ResourceData.HasChange("eula") {
					payload.Properties.Eula = utils.String(state.Eula)
				}

				if metadata.ResourceData.HasChange("privacy_statement_uri") {
					payload.Properties.PrivacyStatementUri = utils.String(state.PrivacyStatementURI)
				}

				if metadata.ResourceData.HasChange("release_note_uri") {
					payload.Properties.ReleaseNoteUri = utils.String(state.ReleaseNoteURI)
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(state.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r GalleryApplicationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.GalleryApplicationsClient
			id, err := galleryapplications.ParseApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r GalleryApplicationResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if oldVal, newVal := metadata.ResourceDiff.GetChange("end_of_life_date"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("end_of_life_date"); err != nil {
					return err
				}
			}

			if oldVal, newVal := metadata.ResourceDiff.GetChange("privacy_statement_uri"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("privacy_statement_uri"); err != nil {
					return err
				}
			}

			if oldVal, newVal := metadata.ResourceDiff.GetChange("release_note_uri"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := metadata.ResourceDiff.ForceNew("release_note_uri"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
