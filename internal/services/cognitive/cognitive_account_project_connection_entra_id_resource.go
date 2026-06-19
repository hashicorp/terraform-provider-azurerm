// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
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

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cognitive_account_project_connection_entra_id -properties "name" -compare-values "subscription_id:cognitive_account_project_id,resource_group_name:cognitive_account_project_id,account_name:cognitive_account_project_id,project_name:cognitive_account_project_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.ResourceWithUpdate   = CognitiveAccountProjectConnectionEntraIDResource{}
	_ sdk.ResourceWithIdentity = CognitiveAccountProjectConnectionEntraIDResource{}
)

type CognitiveAccountProjectConnectionEntraIDResource struct{}

func (r CognitiveAccountProjectConnectionEntraIDResource) Identity() resourceids.ResourceId {
	return new(projectconnectionresource.ProjectConnectionId)
}

func (r CognitiveAccountProjectConnectionEntraIDResource) ResourceType() string {
	return "azurerm_cognitive_account_project_connection_entra_id"
}

func (r CognitiveAccountProjectConnectionEntraIDResource) ModelObject() interface{} {
	return &CognitiveAccountProjectConnectionEntraIDModel{}
}

func (r CognitiveAccountProjectConnectionEntraIDResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projectconnectionresource.ValidateProjectConnectionID
}

type CognitiveAccountProjectConnectionEntraIDModel struct {
	Name                      string            `tfschema:"name"`
	CognitiveAccountProjectId string            `tfschema:"cognitive_account_project_id"`
	AuthType                  string            `tfschema:"auth_type"`
	Category                  string            `tfschema:"category"`
	Target                    string            `tfschema:"target"`
	Metadata                  map[string]string `tfschema:"metadata"`
}

func (r CognitiveAccountProjectConnectionEntraIDResource) Arguments() map[string]*pluginsdk.Schema {
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

		"target": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r CognitiveAccountProjectConnectionEntraIDResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"auth_type": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r CognitiveAccountProjectConnectionEntraIDResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			var model CognitiveAccountProjectConnectionEntraIDModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			projectId, err := projectconnectionresource.ParseProjectID(model.CognitiveAccountProjectId)
			if err != nil {
				return err
			}

			id := projectconnectionresource.NewProjectConnectionID(projectId.SubscriptionId, projectId.ResourceGroupName, projectId.AccountName, projectId.ProjectName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ProjectConnectionsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := projectconnectionresource.AADAuthTypeConnectionProperties{
				AuthType: projectconnectionresource.ConnectionAuthTypeAAD,
				Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
				Target:   pointer.To(model.Target),
			}
			if len(model.Metadata) > 0 {
				properties.Metadata = pointer.To(model.Metadata)
			}

			connection := projectconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.ProjectConnectionsCreate(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r CognitiveAccountProjectConnectionEntraIDResource) Read() sdk.ResourceFunc {
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

			var currentState CognitiveAccountProjectConnectionEntraIDModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountProjectConnectionEntraIDModel{
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

				// Only include metadata fields that were in the original config.
				// The API returns additional metadata fields beyond what was configured (e.g., `ApiVersion`,
				// `DeploymentApiVersion`), which would cause unwanted diffs.
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
					// if metadata is empty in config (e.g., terraform import), read all metadata fields from API
					state.Metadata = pointer.From(base.Metadata)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountProjectConnectionEntraIDResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ProjectConnectionsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			var model CognitiveAccountProjectConnectionEntraIDModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props, ok := resp.Model.Properties.(projectconnectionresource.AADAuthTypeConnectionProperties)
			if !ok {
				return fmt.Errorf("unexpected properties type for %s", *id)
			}

			if metadata.ResourceData.HasChange("target") {
				props.Target = pointer.To(model.Target)
			}

			if metadata.ResourceData.HasChange("metadata") {
				props.Metadata = pointer.To(model.Metadata)
			}

			connection := projectconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: props,
			}

			if _, err := client.ProjectConnectionsCreate(ctx, *id, connection); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountProjectConnectionEntraIDResource) Delete() sdk.ResourceFunc {
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
