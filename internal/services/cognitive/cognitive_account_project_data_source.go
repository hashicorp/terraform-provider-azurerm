// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/cognitiveservicesprojects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
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

var _ sdk.DataSource = CognitiveAccountProjectDataSource{}

type CognitiveAccountProjectDataSource struct{}

func (r CognitiveAccountProjectDataSource) ResourceType() string {
	return "azurerm_cognitive_account_project"
}

func (r CognitiveAccountProjectDataSource) ModelObject() interface{} {
	return &CognitiveAccountProjectDataSourceModel{}
}

func (r CognitiveAccountProjectDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return cognitiveservicesprojects.ValidateProjectID
}

func (r CognitiveAccountProjectDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cognitive_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r CognitiveAccountProjectDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"location": commonschema.LocationComputed(),

		"default": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"endpoints": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

		"tags": commonschema.TagsDataSource(),
	}
}

func (r CognitiveAccountProjectDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

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
