// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource             = CognitiveAccountProjectConnectionResource{}
	_ sdk.ResourceWithIdentity = CognitiveAccountProjectConnectionResource{}
)

type CognitiveAccountProjectConnectionResource struct{}

func (r CognitiveAccountProjectConnectionResource) Identity() resourceids.ResourceId {
	return new(projectconnectionresource.ProjectConnectionId)
}

func (r CognitiveAccountProjectConnectionResource) ResourceType() string {
	return "azurerm_cognitive_account_project_connection"
}

func (r CognitiveAccountProjectConnectionResource) ModelObject() interface{} {
	return &CognitiveAccountProjectConnectionModel{}
}

func (r CognitiveAccountProjectConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projectconnectionresource.ValidateProjectConnectionID
}

type CognitiveAccountProjectConnectionModel struct {
	Name                      string            `tfschema:"name"`
	CognitiveAccountProjectId string            `tfschema:"cognitive_account_project_id"`
	AuthType                  string            `tfschema:"auth_type"`
	Category                  string            `tfschema:"category"`
	Target                    string            `tfschema:"target"`
	Metadata                  map[string]string `tfschema:"metadata"`
}

func (r CognitiveAccountProjectConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountProjectConnectionName(),
		},

		"cognitive_account_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&projectconnectionresource.ProjectId{}),

		"category": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(projectconnectionresource.ConnectionCategoryAIServices),
				string(projectconnectionresource.ConnectionCategoryApiManagement),
				string(projectconnectionresource.ConnectionCategoryAppConfig),
				string(projectconnectionresource.ConnectionCategoryAzureOpenAI),
				string(projectconnectionresource.ConnectionCategoryAzureStorageAccount),
				string(projectconnectionresource.ConnectionCategoryCognitiveService),
				string(projectconnectionresource.ConnectionCategoryCognitiveSearch),
				string(projectconnectionresource.ConnectionCategoryCosmosDb),
				string(projectconnectionresource.ConnectionCategoryDatabricks),
				string(projectconnectionresource.ConnectionCategoryManagedOnlineEndpoint),
				string(projectconnectionresource.ConnectionCategoryMicrosoftFabric),
				string(projectconnectionresource.ConnectionCategorySharepoint),
			}, false),
		},
	}
}

func (r CognitiveAccountProjectConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"auth_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"target": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r CognitiveAccountProjectConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			return errors.New("creating connections with `azurerm_cognitive_account_project_connection` is not supported; use an auth-type-specific resource such as `azurerm_cognitive_account_project_connection_entra_id`")
		},
	}
}

func (r CognitiveAccountProjectConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var currentState CognitiveAccountProjectConnectionModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountProjectConnectionModel{
				CognitiveAccountProjectId: projectconnectionresource.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName).ID(),
				Name:                      id.ConnectionName,
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			if model := resp.Model; model != nil && model.Properties != nil {
				base := model.Properties.ConnectionPropertiesV2()
				state.AuthType = string(base.AuthType)
				state.Category = pointer.FromEnum(base.Category)
				state.Target = pointer.From(base.Target)

				if len(currentState.Metadata) > 0 {
					state.Metadata = map[string]string{}
					apiMetadata := pointer.From(base.Metadata)

					for configKey := range currentState.Metadata {
						for apiKey, apiValue := range apiMetadata {
							if strings.EqualFold(configKey, apiKey) {
								state.Metadata[configKey] = apiValue
								break
							}
						}
					}
				} else {
					state.Metadata = pointer.From(base.Metadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountProjectConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.ProjectConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
