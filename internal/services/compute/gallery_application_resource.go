package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			ValidateFunc: validate.SharedImageGalleryID,
		},

		"location": commonschema.Location(),

		"supported_os_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(compute.OperatingSystemTypesWindows),
				string(compute.OperatingSystemTypesLinux),
			}, false),
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

		"tags": tags.Schema(),
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
	return validate.GalleryApplicationID
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

			galleryId, err := parse.SharedImageGalleryID(state.GalleryId)
			if err != nil {
				return err
			}

			id := parse.NewGalleryApplicationID(subscriptionId, galleryId.ResourceGroup, galleryId.GalleryName, state.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for the presence of existing %q: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			input := compute.GalleryApplication{
				Location: utils.String(location.Normalize(state.Location)),
				GalleryApplicationProperties: &compute.GalleryApplicationProperties{
					SupportedOSType: compute.OperatingSystemTypes(state.SupportedOSType),
				},
				Tags: tags.FromTypedObject(state.Tags),
			}

			if state.Description != "" {
				input.GalleryApplicationProperties.Description = utils.String(state.Description)
			}

			if state.EndOfLifeDate != "" {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				input.GalleryApplicationProperties.EndOfLifeDate = &date.Time{
					Time: endOfLifeDate,
				}
			}

			if state.Eula != "" {
				input.GalleryApplicationProperties.Eula = utils.String(state.Eula)
			}

			if state.PrivacyStatementURI != "" {
				input.GalleryApplicationProperties.PrivacyStatementURI = utils.String(state.PrivacyStatementURI)
			}

			if state.ReleaseNoteURI != "" {
				input.GalleryApplicationProperties.ReleaseNoteURI = utils.String(state.ReleaseNoteURI)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, input)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
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
			id, err := parse.GalleryApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					metadata.Logger.Infof("%q was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			galleryId := parse.NewSharedImageGalleryID(id.SubscriptionId, id.ResourceGroup, id.GalleryName)

			state := &GalleryApplicationModel{
				Name:      id.ApplicationName,
				GalleryId: galleryId.ID(),
				Location:  location.NormalizeNilable(resp.Location),
				Tags:      tags.ToTypedObject(resp.Tags),
			}

			if props := resp.GalleryApplicationProperties; props != nil {
				if v := props.Description; v != nil {
					state.Description = *props.Description
				}

				if v := props.EndOfLifeDate; v != nil {
					state.EndOfLifeDate = props.EndOfLifeDate.Format(time.RFC3339)
				}

				if v := props.Eula; v != nil {
					state.Eula = *props.Eula
				}

				if v := props.PrivacyStatementURI; v != nil {
					state.PrivacyStatementURI = *props.PrivacyStatementURI
				}

				if v := props.ReleaseNoteURI; v != nil {
					state.ReleaseNoteURI = *props.ReleaseNoteURI
				}

				state.SupportedOSType = string(props.SupportedOSType)
			}

			return metadata.Encode(state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r GalleryApplicationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.GalleryApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state GalleryApplicationModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Compute.GalleryApplicationsClient
			existing, err := client.Get(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if metadata.ResourceData.HasChange("description") {
				existing.GalleryApplicationProperties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("end_of_life_date") {
				endOfLifeDate, _ := time.Parse(time.RFC3339, state.EndOfLifeDate)
				existing.GalleryApplicationProperties.EndOfLifeDate = &date.Time{
					Time: endOfLifeDate,
				}
			}

			if metadata.ResourceData.HasChange("eula") {
				existing.GalleryApplicationProperties.Eula = utils.String(state.Eula)
			}

			if metadata.ResourceData.HasChange("privacy_statement_uri") {
				existing.GalleryApplicationProperties.PrivacyStatementURI = utils.String(state.PrivacyStatementURI)
			}

			if metadata.ResourceData.HasChange("release_note_uri") {
				existing.GalleryApplicationProperties.ReleaseNoteURI = utils.String(state.ReleaseNoteURI)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.FromTypedObject(state.Tags)
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName, existing)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
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
			id, err := parse.GalleryApplicationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.GalleryName, id.ApplicationName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
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
			rd := metadata.ResourceDiff

			if oldVal, newVal := rd.GetChange("end_of_life_date"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := rd.ForceNew("end_of_life_date"); err != nil {
					return err
				}
			}

			if oldVal, newVal := rd.GetChange("privacy_statement_uri"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := rd.ForceNew("privacy_statement_uri"); err != nil {
					return err
				}
			}

			if oldVal, newVal := rd.GetChange("release_note_uri"); oldVal.(string) != "" && newVal.(string) == "" {
				if err := rd.ForceNew("release_note_uri"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
