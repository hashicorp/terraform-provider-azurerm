// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_api_key

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
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

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
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

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
