// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cognitive_account_connection_account_managed_identity -properties "name" -compare-values "subscription_id:cognitive_account_id,resource_group_name:cognitive_account_id,account_name:cognitive_account_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.ResourceWithUpdate   = CognitiveAccountConnectionAccountManagedIdentityResource{}
	_ sdk.ResourceWithIdentity = CognitiveAccountConnectionAccountManagedIdentityResource{}
)

type CognitiveAccountConnectionAccountManagedIdentityResource struct{}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Identity() resourceids.ResourceId {
	return new(accountconnectionresource.ConnectionId)
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) ResourceType() string {
	return "azurerm_cognitive_account_connection_account_managed_identity"
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) ModelObject() interface{} {
	return &CognitiveAccountConnectionAccountManagedIdentityModel{}
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}

type CognitiveAccountConnectionAccountManagedIdentityModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	Category           string            `tfschema:"category"`
	Metadata           map[string]string `tfschema:"metadata"`
	Target             string            `tfschema:"target"`
}

type accountManagedIdentityConnectionProperties struct {
	Category *accountconnectionresource.ConnectionCategory `json:"category,omitempty"`
	Metadata *map[string]string                            `json:"metadata,omitempty"`
	Target   *string                                       `json:"target,omitempty"`
}

func (s accountManagedIdentityConnectionProperties) ConnectionPropertiesV2() accountconnectionresource.BaseConnectionPropertiesV2Impl {
	return accountconnectionresource.BaseConnectionPropertiesV2Impl{
		AuthType: accountconnectionresource.ConnectionAuthTypeAccountManagedIdentity,
		Category: s.Category,
		Metadata: s.Metadata,
		Target:   s.Target,
	}
}

func (s accountManagedIdentityConnectionProperties) MarshalJSON() ([]byte, error) {
	type alias accountManagedIdentityConnectionProperties
	wrapper := struct {
		alias
		AuthType string `json:"authType"`
	}{
		alias:    alias(s),
		AuthType: string(accountconnectionresource.ConnectionAuthTypeAccountManagedIdentity),
	}

	encoded, err := json.Marshal(wrapper)
	if err != nil {
		return nil, fmt.Errorf("marshaling accountManagedIdentityConnectionProperties: %+v", err)
	}

	return encoded, nil
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountConnectionName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&accountconnectionresource.AccountId{}),

		"category": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{string(accountconnectionresource.ConnectionCategoryAzureKeyVault)}, false),
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
	}
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			var model CognitiveAccountConnectionAccountManagedIdentityModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := accountconnectionresource.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			id := accountconnectionresource.NewConnectionID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.AccountConnectionsGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := accountManagedIdentityConnectionProperties{
				Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
				Metadata: pointer.To(model.Metadata),
				Target:   pointer.To(model.Target),
			}

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.AccountConnectionsCreate(ctx, id, connection); err != nil {
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

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Read() sdk.ResourceFunc {
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

			var currentState CognitiveAccountConnectionAccountManagedIdentityModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := CognitiveAccountConnectionAccountManagedIdentityModel{
				CognitiveAccountId: accountconnectionresource.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
				Name:               id.ConnectionName,
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			if model := resp.Model; model != nil && model.Properties != nil {
				base := model.Properties.ConnectionPropertiesV2()
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
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Update() sdk.ResourceFunc {
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

			var model CognitiveAccountConnectionAccountManagedIdentityModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			properties := accountManagedIdentityConnectionProperties{
				Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
				Metadata: pointer.To(model.Metadata),
				Target:   pointer.To(model.Target),
			}

			connection := accountconnectionresource.ConnectionPropertiesV2BasicResource{
				Properties: properties,
			}

			if _, err := client.AccountConnectionsCreate(ctx, *id, connection); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveAccountConnectionAccountManagedIdentityResource) Delete() sdk.ResourceFunc {
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
