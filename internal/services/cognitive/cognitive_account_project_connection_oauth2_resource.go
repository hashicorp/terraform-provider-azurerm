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
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = CognitiveAccountProjectConnectionOAuth2Resource{}

type CognitiveAccountProjectConnectionOAuth2Resource struct{}

func (r CognitiveAccountProjectConnectionOAuth2Resource) ResourceType() string {
	return "azurerm_cognitive_account_project_connection_oauth2"
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) ModelObject() interface{} {
	return &CognitiveAccountProjectConnectionOAuth2Model{}
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return projectconnectionresource.ValidateProjectConnectionID
}

type CognitiveAccountProjectConnectionOAuth2Model struct {
	Category           string                        `tfschema:"category"`
	CognitiveProjectId string                        `tfschema:"cognitive_project_id"`
	Metadata           map[string]string             `tfschema:"metadata"`
	Name               string                        `tfschema:"name"`
	OAuth2             []ProjectConnectionOAuth2Auth `tfschema:"oauth2"`
	Target             string                        `tfschema:"target"`
}

type ProjectConnectionOAuth2Auth struct {
	AuthenticationURL string `tfschema:"authentication_url"`
	ClientId          string `tfschema:"client_id"`
	ClientSecret      string `tfschema:"client_secret"`
	DeveloperToken    string `tfschema:"developer_token"`
	Password          string `tfschema:"password"`
	RefreshToken      string `tfschema:"refresh_token"`
	TenantId          string `tfschema:"tenant_id"`
	Username          string `tfschema:"username"`
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ProjectConnectionName(),
		},

		"cognitive_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&projectconnectionresource.ProjectId{}),

		"category": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				projectconnectionresource.PossibleValuesForConnectionCategory(),
				false,
			),
		},

		"target": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"oauth2": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"authentication_url": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPorHTTPS,
					},
					"client_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
					"client_secret": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"developer_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"refresh_token": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
					"username": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
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

func (r CognitiveAccountProjectConnectionOAuth2Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			var model CognitiveAccountProjectConnectionOAuth2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			projectId, err := projectconnectionresource.ParseProjectID(model.CognitiveProjectId)
			if err != nil {
				return err
			}

			id := projectconnectionresource.NewProjectConnectionID(projectId.SubscriptionId, projectId.ResourceGroupName, projectId.AccountName, projectId.ProjectName, model.Name)
			existing, err := client.ProjectConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := projectconnectionresource.OAuth2AuthTypeConnectionProperties{
				AuthType: projectconnectionresource.ConnectionAuthTypeOAuthTwo,
				Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
				Metadata: pointer.To(model.Metadata),
				Target:   pointer.To(model.Target),
				Credentials: &projectconnectionresource.ConnectionOAuth2{
					AuthURL:        pointer.To(model.OAuth2[0].AuthenticationURL),
					ClientId:       pointer.To(model.OAuth2[0].ClientId),
					ClientSecret:   pointer.To(model.OAuth2[0].ClientSecret),
					DeveloperToken: pointer.To(model.OAuth2[0].DeveloperToken),
					Password:       pointer.To(model.OAuth2[0].Password),
					RefreshToken:   pointer.To(model.OAuth2[0].RefreshToken),
					TenantId:       pointer.To(model.OAuth2[0].TenantId),
					Username:       pointer.To(model.OAuth2[0].Username),
				},
			}

			connection := projectconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.ProjectConnectionsCreate(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) Read() sdk.ResourceFunc {
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

			var currentState CognitiveAccountProjectConnectionOAuth2Model
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountProjectConnectionOAuth2Model{
				CognitiveProjectId: projectconnectionresource.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName).ID(),
				Name:               id.ConnectionName,
			}

			if len(currentState.OAuth2) > 0 {
				state.OAuth2 = []ProjectConnectionOAuth2Auth{{
					AuthenticationURL: currentState.OAuth2[0].AuthenticationURL,
					ClientId:          currentState.OAuth2[0].ClientId,
					ClientSecret:      currentState.OAuth2[0].ClientSecret,
					DeveloperToken:    currentState.OAuth2[0].DeveloperToken,
					Password:          currentState.OAuth2[0].Password,
					RefreshToken:      currentState.OAuth2[0].RefreshToken,
					TenantId:          currentState.OAuth2[0].TenantId,
					Username:          currentState.OAuth2[0].Username,
				}}
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				base := props.ConnectionPropertiesV2()
				state.Category = pointer.FromEnum(base.Category)
				state.Target = pointer.From(base.Target)
				state.Metadata = map[string]string{}

				// The API returns additional metadata fields beyond what was configured (e.g., `ApiVersion`,
				// `DeploymentApiVersion`), which would cause unwanted diffs.
				if len(currentState.Metadata) > 0 {
					apiMetadata := pointer.From(base.Metadata)

					for configKey := range currentState.Metadata {
						for apiKey, apiValue := range apiMetadata {
							if strings.EqualFold(configKey, apiKey) {
								state.Metadata[configKey] = apiValue
								break
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountProjectConnectionOAuth2Resource) Update() sdk.ResourceFunc {
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

			var model CognitiveAccountProjectConnectionOAuth2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props := resp.Model.Properties.(projectconnectionresource.OAuth2AuthTypeConnectionProperties)

			props.Credentials = &projectconnectionresource.ConnectionOAuth2{
				AuthURL:        pointer.To(model.OAuth2[0].AuthenticationURL),
				ClientId:       pointer.To(model.OAuth2[0].ClientId),
				ClientSecret:   pointer.To(model.OAuth2[0].ClientSecret),
				DeveloperToken: pointer.To(model.OAuth2[0].DeveloperToken),
				Password:       pointer.To(model.OAuth2[0].Password),
				RefreshToken:   pointer.To(model.OAuth2[0].RefreshToken),
				TenantId:       pointer.To(model.OAuth2[0].TenantId),
				Username:       pointer.To(model.OAuth2[0].Username),
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

func (r CognitiveAccountProjectConnectionOAuth2Resource) Delete() sdk.ResourceFunc {
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
