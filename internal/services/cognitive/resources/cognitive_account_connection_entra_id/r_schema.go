// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_connection_entra_id

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/accountconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Cognitive Services Account Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountConnectionName(),
		},

		"cognitive_account_id": commonschema.ResourceIDReferenceRequiredForceNew(&accountconnectionresource.AccountId{}, "Cognitive Account"),

		"category": {
			Description: "The category of the connection.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(accountconnectionresource.ConnectionCategoryAIServices),
				string(accountconnectionresource.ConnectionCategoryApiManagement),
				string(accountconnectionresource.ConnectionCategoryAppConfig),
				string(accountconnectionresource.ConnectionCategoryAzureOpenAI),
				string(accountconnectionresource.ConnectionCategoryAzureStorageAccount),
				string(accountconnectionresource.ConnectionCategoryCognitiveService),
				string(accountconnectionresource.ConnectionCategoryCognitiveSearch),
				string(accountconnectionresource.ConnectionCategoryCosmosDb),
				string(accountconnectionresource.ConnectionCategoryDatabricks),
				string(accountconnectionresource.ConnectionCategoryManagedOnlineEndpoint),
				string(accountconnectionresource.ConnectionCategoryMicrosoftFabric),
				string(accountconnectionresource.ConnectionCategorySharepoint),
			}, false),
		},

		"target": {
			Description:  "The target endpoint or resource for the connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metadata": {
			Description: "A mapping of metadata key-value pairs for the connection.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
