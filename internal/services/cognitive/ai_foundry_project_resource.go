// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cognitiveValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type AIFoundryProjectModel struct {
	Name               string                                     `tfschema:"name"`
	CognitiveAccountId string                                     `tfschema:"ai_foundry_id"`
	Location           string                                     `tfschema:"location"`
	Identity           []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Description        string                                     `tfschema:"description"`
	DisplayName        string                                     `tfschema:"display_name"`
	Tags               map[string]string                          `tfschema:"tags"`
}

type AIFoundryProject struct{}

var _ sdk.Resource = AIFoundryProject{}

func (AIFoundryProject) ResourceType() string {
	return "azurerm_ai_foundry_project"
}

func (AIFoundryProject) ModelObject() interface{} {
	return &AIFoundryProjectModel{}
}

func (r AIFoundryProject) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cognitiveservicesprojects.ValidateProjectID
}

func (r AIFoundryProject) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveValidate.AccountName(),
		},

		"ai_foundry_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

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

		"tags": commonschema.Tags(),
	}
}

func (r AIFoundryProject) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r AIFoundryProject) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model AIFoundryProjectModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			client := metadata.Client.Cognitive.ProjectsClient
			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id := cognitiveservicesprojects.NewProjectID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.ProjectsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			props := cognitiveservicesprojects.Project{
				Identity:   expandIdentity,
				Location:   pointer.To(location.Normalize(model.Location)),
				Properties: &cognitiveservicesprojects.ProjectProperties{},
				Tags:       pointer.To(model.Tags),
			}

			if model.Description != "" {
				props.Properties.Description = &model.Description
			}

			if model.DisplayName != "" {
				props.Properties.DisplayName = &model.DisplayName
			}

			if err := client.ProjectsCreateThenPoll(ctx, id, props); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r AIFoundryProject) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			var model AIFoundryProjectModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}
			locks.ByID(accountId.ID())
			defer locks.UnlockByID(accountId.ID())

			id, err := cognitiveservicesprojects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			resp, err := client.ProjectsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			props := resp.Model
			if metadata.ResourceData.HasChange("description") {
				props.Properties.Description = pointer.To(model.Description)
			}
			if metadata.ResourceData.HasChange("display_name") {
				props.Properties.DisplayName = pointer.To(model.DisplayName)
			}
			if metadata.ResourceData.HasChange("tags") {
				props.Tags = pointer.To(model.Tags)
			}
			if metadata.ResourceData.HasChange("identity") {
				expandIdentity, err := identity.ExpandSystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				props.Identity = expandIdentity
			}

			if err := client.ProjectsUpdateThenPoll(ctx, *id, *props); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r AIFoundryProject) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient

			id, err := cognitiveservicesprojects.ParseProjectID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := AIFoundryProjectModel{
				Name:               id.ProjectName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			identityFlatten, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			state.Location = location.NormalizeNilable(model.Location)
			state.Identity = *identityFlatten
			if props := model.Properties; props != nil {
				state.Description = pointer.From(props.Description)
				state.DisplayName = pointer.From(props.DisplayName)
			}
			state.Tags = pointer.From(model.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r AIFoundryProject) Delete() sdk.ResourceFunc {
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
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
