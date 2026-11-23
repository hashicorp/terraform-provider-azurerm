// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_project_connection_entra_id

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/projectconnectionresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Cognitive Services Account Project Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountProjectConnectionName(),
		},

		"cognitive_account_project_id": commonschema.ResourceIDReferenceRequiredForceNew(&projectconnectionresource.ProjectId{}),

		"category": {
			Description: "The category of the connection.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
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
	return map[string]*pluginsdk.Schema{
		"authentication_type": {
			Description: "The authentication type of the connection.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
