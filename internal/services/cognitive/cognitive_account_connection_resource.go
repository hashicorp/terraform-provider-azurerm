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
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = CognitiveAccountConnectionResource{}

type CognitiveAccountConnectionResource struct{}

func (r CognitiveAccountConnectionResource) ResourceType() string {
	return "azurerm_cognitive_account_connection"
}

func (r CognitiveAccountConnectionResource) ModelObject() interface{} {
	return &CognitiveAccountConnectionModel{}
}

func (r CognitiveAccountConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}

type CognitiveAccountConnectionModel struct {
	ApiKey             string            `tfschema:"api_key"`
	AuthType           string            `tfschema:"auth_type"`
	Category           string            `tfschema:"category"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	CustomKeys         map[string]string `tfschema:"custom_keys"`
	Metadata           map[string]string `tfschema:"metadata"`
	Name               string            `tfschema:"name"`
	OAuth2             []OAuth2AuthModel `tfschema:"oauth2"`
	Target             string            `tfschema:"target"`
}

type OAuth2AuthModel struct {
	AuthURL        string `tfschema:"auth_url"`
	ClientId       string `tfschema:"client_id"`
	ClientSecret   string `tfschema:"client_secret"`
	DeveloperToken string `tfschema:"developer_token"`
	Password       string `tfschema:"password"`
	RefreshToken   string `tfschema:"refresh_token"`
	TenantId       string `tfschema:"tenant_id"`
	Username       string `tfschema:"username"`
}

func (r CognitiveAccountConnectionResource) Arguments() map[string]*pluginsdk.Schema {
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

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&accountconnectionresource.AccountId{}),

		"auth_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// From service team we only support these auth types for Nov 2025
			ValidateFunc: validation.StringInSlice([]string{
				string(accountconnectionresource.ConnectionAuthTypeAAD),
				string(accountconnectionresource.ConnectionAuthTypeApiKey),
				string(accountconnectionresource.ConnectionAuthTypeCustomKeys),
				string(accountconnectionresource.ConnectionAuthTypeOAuthTwo),
			}, false),
		},

		"category": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice(
				accountconnectionresource.PossibleValuesForConnectionCategory(),
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

func (r CognitiveAccountConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			var model CognitiveAccountConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := accountconnectionresource.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			id := accountconnectionresource.NewConnectionID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.AccountConnectionsGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties, err := expandConnectionProperties(model)
			if err != nil {
				return fmt.Errorf("expanding `properties`: %+v", err)
			}

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.AccountConnectionsCreate(ctx, id, connection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CognitiveAccountConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountConnectionsGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			// Get current state to preserve sensitive values not returned by API
			var currentState CognitiveAccountConnectionModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountConnectionModel{
				CognitiveAccountId: accountconnectionresource.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
				Name:               id.ConnectionName,
			}

			if model := resp.Model; model != nil {
				props := model.Properties

				base := props.ConnectionPropertiesV2()
				state.AuthType = string(base.AuthType)
				state.Category = pointer.FromEnum(base.Category)
				state.Target = pointer.From(base.Target)
				state.Metadata = map[string]string{}

				// Only include metadata fields that were in the original config
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

				// Handle auth-specific properties
				// Note: Sensitive credentials are not returned by API, so preserve from config
				switch props.(type) {
				case accountconnectionresource.ApiKeyAuthConnectionProperties:
					state.ApiKey = currentState.ApiKey

				case accountconnectionresource.OAuth2AuthTypeConnectionProperties:
					if len(currentState.OAuth2) > 0 {
						state.OAuth2 = []OAuth2AuthModel{{
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

				case accountconnectionresource.CustomKeysConnectionProperties:
					state.CustomKeys = currentState.CustomKeys
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CognitiveAccountConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties, err := expandConnectionPropertiesForUpdate(model, metadata.ResourceData)
			if err != nil {
				return fmt.Errorf("expanding `properties`: %+v", err)
			}

			updateContent := accountconnectionresource.ConnectionUpdateContent{
				Properties: properties,
			}

			if _, err := client.AccountConnectionsUpdate(ctx, *id, updateContent); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.AccountConnectionsDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandConnectionProperties(model CognitiveAccountConnectionModel) (accountconnectionresource.ConnectionPropertiesV2, error) {
	switch authType := accountconnectionresource.ConnectionAuthType(model.AuthType); authType {
	case accountconnectionresource.ConnectionAuthTypeApiKey:
		if model.ApiKey == "" {
			return nil, errors.New("when `auth_type` is `ApiKey`, `api_key` must be specified")
		}

		return accountconnectionresource.ApiKeyAuthConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
			Credentials: &accountconnectionresource.ConnectionApiKey{
				Key: pointer.To(model.ApiKey),
			},
		}, nil

	case accountconnectionresource.ConnectionAuthTypeOAuthTwo:
		if len(model.OAuth2) == 0 {
			return nil, errors.New("when `auth_type` is `OAuth2`, `oauth2` block must be specified")
		}

		return accountconnectionresource.OAuth2AuthTypeConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
			Credentials: &accountconnectionresource.ConnectionOAuth2{
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

	case accountconnectionresource.ConnectionAuthTypeCustomKeys:
		if len(model.CustomKeys) == 0 {
			return nil, errors.New("when `auth_type` is `CustomKeys`, `custom_keys` must be specified")
		}

		return accountconnectionresource.CustomKeysConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
			Credentials: &accountconnectionresource.CustomKeys{
				Keys: &model.CustomKeys,
			},
		}, nil

	case accountconnectionresource.ConnectionAuthTypeAAD:
		if model.ApiKey != "" || len(model.OAuth2) > 0 || len(model.CustomKeys) > 0 {
			return nil, errors.New("when `auth_type` is `AAD`, no other auth configuration blocks should be specified")
		}

		return accountconnectionresource.AADAuthTypeConnectionProperties{
			AuthType: authType,
			Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
			Metadata: pointer.To(model.Metadata),
			Target:   pointer.To(model.Target),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported `auth_type`: %q", model.AuthType)
	}
}

func expandConnectionPropertiesForUpdate(model CognitiveAccountConnectionModel, d *pluginsdk.ResourceData) (accountconnectionresource.ConnectionPropertiesV2, error) {
	// The auth_type will be used to decide which properties need updating, so we do not check hasChange for auth_type.
	switch authType := accountconnectionresource.ConnectionAuthType(model.AuthType); authType {
	case accountconnectionresource.ConnectionAuthTypeApiKey:
		props := accountconnectionresource.ApiKeyAuthConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("api_key") {
			if model.ApiKey == "" {
				return nil, errors.New("when `auth_type` is `ApiKey`, `api_key` must be specified")
			}
			props.Credentials = &accountconnectionresource.ConnectionApiKey{
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

	case accountconnectionresource.ConnectionAuthTypeOAuthTwo:
		props := accountconnectionresource.OAuth2AuthTypeConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("oauth2") {
			if len(model.OAuth2) == 0 {
				return nil, errors.New("when `auth_type` is `OAuth2`, `oauth2` block must be specified")
			}
			props.Credentials = &accountconnectionresource.ConnectionOAuth2{
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

	case accountconnectionresource.ConnectionAuthTypeCustomKeys:
		props := accountconnectionresource.CustomKeysConnectionProperties{
			AuthType: authType,
		}

		if d.HasChange("custom_keys") {
			if len(model.CustomKeys) == 0 {
				return nil, errors.New("when `auth_type` is `CustomKeys`, `custom_keys` must be specified")
			}
			props.Credentials = &accountconnectionresource.CustomKeys{
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

	case accountconnectionresource.ConnectionAuthTypeAAD:
		props := accountconnectionresource.AADAuthTypeConnectionProperties{
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
