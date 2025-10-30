package cognitive

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CognitiveAccountProjectModel struct {
	CognitiveAccountId string                                     `tfschema:"cognitive_account_id"`
	Description        string                                     `tfschema:"description"`
	DisplayName        string                                     `tfschema:"display_name"`
	Endpoints          map[string]string                          `tfschema:"endpoints"`
	Identity           []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	IsDefault          bool                                       `tfschema:"is_default"`
	Location           string                                     `tfschema:"location"`
	Name               string                                     `tfschema:"name"`
	Tags               map[string]string                          `tfschema:"tags"`
}

type CognitiveAccountProjectResource struct{}

var _ sdk.ResourceWithUpdate = CognitiveAccountProjectResource{}

func (r CognitiveAccountProjectResource) ResourceType() string {
	return "azurerm_cognitive_account_project"
}

func (r CognitiveAccountProjectResource) ModelObject() interface{} {
	return &CognitiveAccountProjectModel{}
}

func (r CognitiveAccountProjectResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cognitiveservicesprojects.ValidateProjectID
}

func (r CognitiveAccountProjectResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}
			// Azure API issue: description and display_name cannot be updated to empty values, see: https://github.com/Azure/azure-rest-api-specs/issues/38422
			if oldDesc, newDesc := metadata.ResourceDiff.GetChange("description"); oldDesc.(string) != "" && newDesc.(string) == "" {
				metadata.ResourceDiff.ForceNew("description")
			}

			if oldDisplay, newDisplay := metadata.ResourceDiff.GetChange("display_name"); oldDisplay.(string) != "" && newDisplay.(string) == "" {
				metadata.ResourceDiff.ForceNew("display_name")
			}

			return nil
		},
	}
}

func (r CognitiveAccountProjectResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]{1,63}$"),
				"name must be between 2 and 64 characters long, start with an alphanumeric character, and contain only alphanumeric characters, dashes, periods or underscores.",
			),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&cognitiveservicesprojects.AccountId{}),

		"location": commonschema.Location(),

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"tags": commonschema.Tags(),
	}
}

func (r CognitiveAccountProjectResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"endpoints": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"is_default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r CognitiveAccountProjectResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			var model CognitiveAccountProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := cognitiveservicesprojects.NewProjectID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)

			existing, err := client.ProjectsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			payload := cognitiveservicesprojects.Project{
				Identity: expandedIdentity,
				Location: pointer.To(location.Normalize(model.Location)),
				Tags:     pointer.To(model.Tags),
				Properties: &cognitiveservicesprojects.ProjectProperties{
					Description: pointer.To(model.Description),
					DisplayName: pointer.To(model.DisplayName),
				},
			}

			if err := client.ProjectsCreateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CognitiveAccountProjectResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			id, err := cognitiveservicesprojects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CognitiveAccountProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			existing, err := client.ProjectsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				payload.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("description") {
				if payload.Properties == nil {
					payload.Properties = &cognitiveservicesprojects.ProjectProperties{}
				}
				payload.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("display_name") {
				if payload.Properties == nil {
					payload.Properties = &cognitiveservicesprojects.ProjectProperties{}
				}
				payload.Properties.DisplayName = pointer.To(model.DisplayName)
			}

			if err := client.ProjectsCreateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountProjectResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			id, err := cognitiveservicesprojects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.ProjectsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CognitiveAccountProjectModel{
				Name:               id.ProjectName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if model := existing.Model; model != nil {
				state.Location = location.NormalizeNilable(model.Location)
				state.Tags = pointer.From(model.Tags)

				flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = *flattenedIdentity

				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.DisplayName = pointer.From(props.DisplayName)
					state.Endpoints = pointer.From(props.Endpoints)
					state.IsDefault = pointer.From(props.IsDefault)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountProjectResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			id, err := cognitiveservicesprojects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			accountId := cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName)

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			if err := client.ProjectsDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
