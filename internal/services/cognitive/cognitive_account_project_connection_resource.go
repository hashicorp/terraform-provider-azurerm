// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = CognitiveAccountProjectConnectionResource{}

type CognitiveAccountProjectConnectionResource struct{}

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
	ApiKey             string                   `tfschema:"api_key"`
	AuthType           string                   `tfschema:"auth_type"`
	Category           string                   `tfschema:"category"`
	CognitiveProjectId string                   `tfschema:"cognitive_project_id"`
	CustomKeys         map[string]string        `tfschema:"custom_keys"`
	Metadata           map[string]string        `tfschema:"metadata"`
	Name               string                   `tfschema:"name"`
	OAuth2             []ProjectOAuth2AuthModel `tfschema:"oauth2"`
	Target             string                   `tfschema:"target"`
}

type ProjectOAuth2AuthModel struct {
	AuthURL        string `tfschema:"auth_url"`
	ClientId       string `tfschema:"client_id"`
	ClientSecret   string `tfschema:"client_secret"`
	DeveloperToken string `tfschema:"developer_token"`
	Password       string `tfschema:"password"`
	RefreshToken   string `tfschema:"refresh_token"`
	TenantId       string `tfschema:"tenant_id"`
	Username       string `tfschema:"username"`
}

func (r CognitiveAccountProjectConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]{2,32}$"),
				"`name` must be between 3 and 33 characters long, start with an alphanumeric character, and contain only alphanumeric characters, dashes(-) or underscores(_).",
			),
		},

		"cognitive_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&projectconnectionresource.ProjectId{}),

		"auth_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(projectconnectionresource.ConnectionAuthTypeAAD),
				string(projectconnectionresource.ConnectionAuthTypeApiKey),
				string(projectconnectionresource.ConnectionAuthTypeCustomKeys),
				string(projectconnectionresource.ConnectionAuthTypeOAuthTwo),
			}, false),
		},

		"category": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				projectconnectionresource.PossibleValuesForConnectionCategory(),
				false,
			),
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"target": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"api_key": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"oauth2", "custom_keys"},
		},

		"custom_keys": {
			Type:          pluginsdk.TypeMap,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"api_key", "oauth2"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"oauth2": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"api_key", "custom_keys"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"auth_url": {
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
	}
}

func (r CognitiveAccountProjectConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountProjectConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			var model CognitiveAccountProjectConnectionModel
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

			properties, err := expandProjectConnectionProperties(model)
			if err != nil {
				return fmt.Errorf("expanding `properties`: %+v", err)
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
				CognitiveProjectId: projectconnectionresource.NewProjectID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ProjectName).ID(),
				Name:               id.ConnectionName,
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				base := props.ConnectionPropertiesV2()
				state.AuthType = string(base.AuthType)
				state.Category = pointer.FromEnum(base.Category)
				state.Target = pointer.From(base.Target)
				state.Metadata = map[string]string{}

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

				switch props.(type) {
				case projectconnectionresource.ApiKeyAuthConnectionProperties:
					state.ApiKey = currentState.ApiKey

				case projectconnectionresource.OAuth2AuthTypeConnectionProperties:
					if len(currentState.OAuth2) > 0 {
						state.OAuth2 = []ProjectOAuth2AuthModel{{
							AuthURL:        currentState.OAuth2[0].AuthURL,
							ClientId:       currentState.OAuth2[0].ClientId,
							ClientSecret:   currentState.OAuth2[0].ClientSecret,
							DeveloperToken: currentState.OAuth2[0].DeveloperToken,
							Password:       currentState.OAuth2[0].Password,
							RefreshToken:   currentState.OAuth2[0].RefreshToken,
							TenantId:       currentState.OAuth2[0].TenantId,
							Username:       currentState.OAuth2[0].Username,
						}}
					}

				case projectconnectionresource.CustomKeysConnectionProperties:
					state.CustomKeys = currentState.CustomKeys
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountProjectConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.ProjectConnectionResourceClient

			id, err := projectconnectionresource.ParseProjectConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CognitiveAccountProjectConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties, err := expandProjectConnectionPropertiesForUpdate(model, metadata.ResourceData)
			if err != nil {
				return fmt.Errorf("expanding `properties`: %+v", err)
			}

			updateContent := projectconnectionresource.ConnectionUpdateContent{
				Properties: properties,
			}

			if _, err := client.ProjectConnectionsUpdate(ctx, *id, updateContent); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
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

func expandProjectConnectionProperties(model CognitiveAccountProjectConnectionModel) (projectconnectionresource.ConnectionPropertiesV2, error) {
	switch authType := projectconnectionresource.ConnectionAuthType(model.AuthType); authType {
	case projectconnectionresource.ConnectionAuthTypeApiKey:
		if model.ApiKey == "" {
			return nil, errors.New("when `auth_type` is `ApiKey`, `api_key` must be specified")
		}

		return projectconnectionresource.ApiKeyAuthConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
			Credentials: &projectconnectionresource.ConnectionApiKey{
				Key: pointer.To(model.ApiKey),
			},
		}, nil

	case projectconnectionresource.ConnectionAuthTypeOAuthTwo:
		if len(model.OAuth2) == 0 {
			return nil, errors.New("when `auth_type` is `OAuth2`, `oauth2` block must be specified")
		}

		return projectconnectionresource.OAuth2AuthTypeConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
			Credentials: &projectconnectionresource.ConnectionOAuth2{
				AuthURL:        pointer.To(model.OAuth2[0].AuthURL),
				ClientId:       pointer.To(model.OAuth2[0].ClientId),
				ClientSecret:   pointer.To(model.OAuth2[0].ClientSecret),
				DeveloperToken: pointer.To(model.OAuth2[0].DeveloperToken),
				Password:       pointer.To(model.OAuth2[0].Password),
				RefreshToken:   pointer.To(model.OAuth2[0].RefreshToken),
				TenantId:       pointer.To(model.OAuth2[0].TenantId),
				Username:       pointer.To(model.OAuth2[0].Username),
			},
		}, nil

	case projectconnectionresource.ConnectionAuthTypeCustomKeys:
		if len(model.CustomKeys) == 0 {
			return nil, errors.New("when `auth_type` is `CustomKeys`, `custom_keys` must be specified")
		}

		return projectconnectionresource.CustomKeysConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
			Credentials: &projectconnectionresource.CustomKeys{
				Keys: &model.CustomKeys,
			},
		}, nil

	case projectconnectionresource.ConnectionAuthTypeAAD:
		if model.ApiKey != "" || len(model.OAuth2) > 0 || len(model.CustomKeys) > 0 {
			return nil, errors.New("when `auth_type` is `AAD`, no other auth configuration blocks should be specified")
		}

		return projectconnectionresource.AADAuthTypeConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[projectconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported `auth_type`: %q", model.AuthType)
	}
}

func expandProjectConnectionPropertiesForUpdate(model CognitiveAccountProjectConnectionModel, d *pluginsdk.ResourceData) (projectconnectionresource.ConnectionPropertiesV2, error) {
	switch authType := projectconnectionresource.ConnectionAuthType(model.AuthType); authType {
	case projectconnectionresource.ConnectionAuthTypeApiKey:
		props := projectconnectionresource.ApiKeyAuthConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("api_key") {
			if model.ApiKey == "" {
				return nil, errors.New("when `auth_type` is `ApiKey`, `api_key` must be specified")
			}
			props.Credentials = &projectconnectionresource.ConnectionApiKey{
				Key: pointer.To(model.ApiKey),
			}
		}

		if d.HasChange("target") {
			props.Target = pointer.To(model.Target)
		}

		if d.HasChange("metadata") {
			props.Metadata = pointer.To(model.Metadata)
		}

		return props, nil

	case projectconnectionresource.ConnectionAuthTypeOAuthTwo:
		props := projectconnectionresource.OAuth2AuthTypeConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("oauth2") {
			if len(model.OAuth2) == 0 {
				return nil, errors.New("when `auth_type` is `OAuth2`, `oauth2` block must be specified")
			}
			props.Credentials = &projectconnectionresource.ConnectionOAuth2{
				AuthURL:        pointer.To(model.OAuth2[0].AuthURL),
				ClientId:       pointer.To(model.OAuth2[0].ClientId),
				ClientSecret:   pointer.To(model.OAuth2[0].ClientSecret),
				DeveloperToken: pointer.To(model.OAuth2[0].DeveloperToken),
				Password:       pointer.To(model.OAuth2[0].Password),
				RefreshToken:   pointer.To(model.OAuth2[0].RefreshToken),
				TenantId:       pointer.To(model.OAuth2[0].TenantId),
				Username:       pointer.To(model.OAuth2[0].Username),
			}
		}

		if d.HasChange("target") {
			props.Target = pointer.To(model.Target)
		}

		if d.HasChange("metadata") {
			props.Metadata = pointer.To(model.Metadata)
		}

		return props, nil

	case projectconnectionresource.ConnectionAuthTypeCustomKeys:
		props := projectconnectionresource.CustomKeysConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("custom_keys") {
			if len(model.CustomKeys) == 0 {
				return nil, errors.New("when `auth_type` is `CustomKeys`, `custom_keys` must be specified")
			}
			props.Credentials = &projectconnectionresource.CustomKeys{
				Keys: &model.CustomKeys,
			}
		}

		if d.HasChange("target") {
			props.Target = pointer.To(model.Target)
		}

		if d.HasChange("metadata") {
			props.Metadata = pointer.To(model.Metadata)
		}

		return props, nil

	case projectconnectionresource.ConnectionAuthTypeAAD:
		props := projectconnectionresource.AADAuthTypeConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("target") {
			props.Target = pointer.To(model.Target)
		}

		if d.HasChange("metadata") {
			props.Metadata = pointer.To(model.Metadata)
		}

		return props, nil

	default:
		return nil, fmt.Errorf("unsupported `auth_type`: %q", model.AuthType)
	}
}
