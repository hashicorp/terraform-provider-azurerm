// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	cognitiveClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CognitiveAccountProjectDataSourceModel struct {
	CognitiveAccountName string                                     `tfschema:"cognitive_account_name"`
	Default              bool                                       `tfschema:"default"`
	Description          string                                     `tfschema:"description"`
	DisplayName          string                                     `tfschema:"display_name"`
	Endpoints            map[string]string                          `tfschema:"endpoints"`
	Identity             []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Location             string                                     `tfschema:"location"`
	Name                 string                                     `tfschema:"name"`
	ResourceGroupName    string                                     `tfschema:"resource_group_name"`
	Tags                 map[string]string                          `tfschema:"tags"`
}

var _ sdk.DataSource = DataSource{}

type DataSource struct{}

func (r DataSource) ResourceType() string {
	return "azurerm_cognitive_account_project"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about an existing Cognitive Services Account Project."
}

func (r DataSource) ModelObject() any {
	return &CognitiveAccountProjectDataSourceModel{}
}

func (r DataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cognitiveservicesprojects.ValidateProjectID
}

func (r DataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "The name of the Cognitive Services Account Project.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AccountProjectName(),
		},

		"cognitive_account_name": {
			Description:  "The name of the Cognitive Services Account in which the Project exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.AccountName(),
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location": commonschema.LocationComputed(),

		"default": {
			Description: "Whether this is the default project for the Cognitive Services Account.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"description": {
			Description: "The description of the Cognitive Services Account Project.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"display_name": {
			Description: "The display name of the Cognitive Services Account Project.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"endpoints": {
			Description: "A mapping of endpoint names to endpoint URLs for the Cognitive Services Account Project.",
			Type:        pluginsdk.TypeMap,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := cognitiveClient.From(metadata.Client).ProjectsClient
			subscriptionId := account.From(metadata.Client).SubscriptionId

			var model CognitiveAccountProjectDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := cognitiveservicesprojects.NewProjectID(subscriptionId, model.ResourceGroupName, model.CognitiveAccountName, model.Name)

			existing, err := client.ProjectsGet(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			metadata.SetID(id)

			state := CognitiveAccountProjectDataSourceModel{
				Name:                 id.ProjectName,
				CognitiveAccountName: id.AccountName,
				ResourceGroupName:    id.ResourceGroupName,
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
					state.Default = pointer.From(props.IsDefault)
					state.Description = pointer.From(props.Description)
					state.DisplayName = pointer.From(props.DisplayName)
					state.Endpoints = pointer.From(props.Endpoints)
				}
			}

			return metadata.Encode(&state)
		},
	}
}
