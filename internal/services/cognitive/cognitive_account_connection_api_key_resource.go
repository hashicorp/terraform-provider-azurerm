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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name cognitive_account_connection_api_key -properties "name" -compare-values "subscription_id:cognitive_account_id,resource_group_name:cognitive_account_id,account_name:cognitive_account_id" -test-name "basic" -test-expect-non-empty

var (
	_ sdk.ResourceWithUpdate         = CognitiveAccountConnectionApiKeyResource{}
	_ sdk.ResourceWithIdentity       = CognitiveAccountConnectionApiKeyResource{}
	_ sdk.ResourceWithCustomImporter = CognitiveAccountConnectionApiKeyResource{}
)

type CognitiveAccountConnectionApiKeyResource struct{}

func (r CognitiveAccountConnectionApiKeyResource) CustomImporter() sdk.ResourceRunFunc {
	return cognitiveAccountConnectionImporter(accountconnectionresource.ConnectionAuthTypeApiKey, r.ResourceType())
}

func (r CognitiveAccountConnectionApiKeyResource) Identity() resourceids.ResourceId {
	return new(accountconnectionresource.ConnectionId)
}

func (r CognitiveAccountConnectionApiKeyResource) ResourceType() string {
	return "azurerm_cognitive_account_connection_api_key"
}

func (r CognitiveAccountConnectionApiKeyResource) ModelObject() interface{} {
	return &CognitiveAccountConnectionApiKeyModel{}
}

func (r CognitiveAccountConnectionApiKeyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return accountconnectionresource.ValidateConnectionID
}

type CognitiveAccountConnectionApiKeyModel struct {
	Name               string            `tfschema:"name"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	ApiKey             string            `tfschema:"api_key"`
	Category           string            `tfschema:"category"`
	Metadata           map[string]string `tfschema:"metadata"`
	Target             string            `tfschema:"target"`
}

func (r CognitiveAccountConnectionApiKeyResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountConnectionName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&accountconnectionresource.AccountId{}),

		"api_key": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"category": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(accountconnectionresource.ConnectionCategoryAIServices),
				string(accountconnectionresource.ConnectionCategoryApiManagement),
				string(accountconnectionresource.ConnectionCategoryApiKey),
				string(accountconnectionresource.ConnectionCategoryAppConfig),
				string(accountconnectionresource.ConnectionCategoryAppInsights),
				string(accountconnectionresource.ConnectionCategoryAzureOpenAI),
				string(accountconnectionresource.ConnectionCategoryBingLLMSearch),
				string(accountconnectionresource.ConnectionCategoryCognitiveService),
				string(accountconnectionresource.ConnectionCategoryCognitiveSearch),
				string(accountconnectionresource.ConnectionCategoryGroundingWithBingSearch),
				string(accountconnectionresource.ConnectionCategoryGroundingWithCustomSearch),
				string(accountconnectionresource.ConnectionCategoryModelGateway),
				string(accountconnectionresource.ConnectionCategoryOpenAI),
				string(accountconnectionresource.ConnectionCategoryPinecone),
				string(accountconnectionresource.ConnectionCategorySerp),
				string(accountconnectionresource.ConnectionCategoryServerless),
			}, false),
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"target": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r CognitiveAccountConnectionApiKeyResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model CognitiveAccountConnectionApiKeyModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding diff: %+v", err)
			}

			if model.Category != string(accountconnectionresource.ConnectionCategoryOpenAI) &&
				model.Category != string(accountconnectionresource.ConnectionCategoryPinecone) &&
				model.Category != string(accountconnectionresource.ConnectionCategorySerp) &&
				metadata.ResourceDiff.GetRawConfig().GetAttr("target").IsNull() {
				return fmt.Errorf("`target` must be specified when `category` is `%s`", model.Category)
			}

			return nil
		},
	}
}

func (r CognitiveAccountConnectionApiKeyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveAccountConnectionApiKeyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.AccountConnectionResourceClient

			var model CognitiveAccountConnectionApiKeyModel
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

			properties := accountconnectionresource.ApiKeyAuthConnectionProperties{
				AuthType: accountconnectionresource.ConnectionAuthTypeApiKey,
				Category: pointer.ToEnum[accountconnectionresource.ConnectionCategory](model.Category),
				Credentials: &accountconnectionresource.ConnectionApiKey{
					Key: pointer.To(model.ApiKey),
				},
			}
			if len(model.Metadata) > 0 {
				properties.Metadata = pointer.To(model.Metadata)
			}
			if model.Target != "" {
				properties.Target = pointer.To(model.Target)
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

func (r CognitiveAccountConnectionApiKeyResource) Read() sdk.ResourceFunc {
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

			var currentState CognitiveAccountConnectionApiKeyModel
			if err := metadata.Decode(&currentState); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			return r.flatten(metadata, id, resp.Model, currentState.Metadata, currentState.ApiKey)
		},
	}
}

func (r CognitiveAccountConnectionApiKeyResource) Update() sdk.ResourceFunc {
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

			var model CognitiveAccountConnectionApiKeyModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			props, ok := resp.Model.Properties.(accountconnectionresource.ApiKeyAuthConnectionProperties)
			if !ok {
				return fmt.Errorf("unexpected properties type for %s", *id)
			}

			if metadata.ResourceData.HasChange("api_key") {
				props.Credentials = &accountconnectionresource.ConnectionApiKey{
					Key: pointer.To(model.ApiKey),
				}
			}

			if metadata.ResourceData.HasChange("target") {
				if model.Target == "" {
					props.Target = nil
				} else {
					props.Target = pointer.To(model.Target)
				}
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

func (r CognitiveAccountConnectionApiKeyResource) Delete() sdk.ResourceFunc {
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

func (CognitiveAccountConnectionApiKeyResource) flatten(metadata sdk.ResourceMetaData, id *accountconnectionresource.ConnectionId, model *accountconnectionresource.ConnectionPropertiesV2BasicResource, priorMetadata map[string]string, priorApiKey string) error {
	state := CognitiveAccountConnectionApiKeyModel{
		CognitiveAccountId: accountconnectionresource.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
		Name:               id.ConnectionName,
		ApiKey:             priorApiKey,
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	if model != nil && model.Properties != nil {
		base := model.Properties.ConnectionPropertiesV2()
		state.Category = pointer.FromEnum(base.Category)
		state.Target = pointer.From(base.Target)
		state.Metadata = flattenAccountConnectionMetadata(priorMetadata, base.Metadata)
	}

	return metadata.Encode(&state)
}
