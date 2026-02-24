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
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2025-06-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = CognitiveAccountConnectionOAuth2Resource{}

type CognitiveAccountConnectionOAuth2Resource struct{}

func (r CognitiveAccountConnectionOAuth2Resource) ResourceType() string {
	return "azurerm_cognitive_account_connection_oauth2"
}

func (r CognitiveAccountConnectionOAuth2Resource) ModelObject() interface{} {
	return &CognitiveAccountConnectionOAuth2Model{}
}

func (r CognitiveAccountConnectionOAuth2Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}

type CognitiveAccountConnectionOAuth2Model struct {
	Category           string            `tfschema:"category"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	Metadata           map[string]string `tfschema:"metadata"`
	Name               string            `tfschema:"name"`
	OAuth2             []OAuth2AuthModel `tfschema:"oauth2"`
	Target             string            `tfschema:"target"`
}

type OAuth2AuthModel struct {
	AuthenticationURL string `tfschema:"authentication_url"`
	ClientId          string `tfschema:"client_id"`
	ClientSecret      string `tfschema:"client_secret"`
	DeveloperToken    string `tfschema:"developer_token"`
	Password          string `tfschema:"password"`
	RefreshToken      string `tfschema:"refresh_token"`
	TenantId          string `tfschema:"tenant_id"`
	Username          string `tfschema:"username"`
}

func (r CognitiveAccountConnectionOAuth2Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountConnectionName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&accountconnectionresource.AccountId{}),

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
	}
}

func (r CognitiveAccountConnectionOAuth2Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountConnectionOAuth2Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			var model CognitiveAccountConnectionOAuth2Model
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

			properties := accountconnectionresource.OAuth2AuthTypeConnectionProperties{
				AuthType: accountconnectionresource.ConnectionAuthTypeOAuthTwo,
				Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
				Metadata: pointer.To(model.Metadata),
				Target:   pointer.To(model.Target),
				Credentials: &accountconnectionresource.ConnectionOAuth2{
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

func (r CognitiveAccountConnectionOAuth2Resource) Read() sdk.ResourceFunc {
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

			var currentState CognitiveAccountConnectionOAuth2Model
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountConnectionOAuth2Model{
				CognitiveAccountId: accountconnectionresource.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
				Name:               id.ConnectionName,
			}

			if len(currentState.OAuth2) > 0 {
				state.OAuth2 = []OAuth2AuthModel{{
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

				// Only include metadata fields that were in the original config.
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

func (r CognitiveAccountConnectionOAuth2Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			id, err := accountconnectionresource.ParseConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.AccountConnectionsGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			var model CognitiveAccountConnectionOAuth2Model
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props, ok := resp.Model.Properties.(accountconnectionresource.OAuth2AuthTypeConnectionProperties)
			if !ok {
				return fmt.Errorf("unexpected properties type for %s", *id)
			}

			props.Credentials = &accountconnectionresource.ConnectionOAuth2{
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

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: props,
			}

			if _, err := client.AccountConnectionsCreate(ctx, *id, connection); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountConnectionOAuth2Resource) Delete() sdk.ResourceFunc {
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
