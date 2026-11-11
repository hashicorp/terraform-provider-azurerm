// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_assembly

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppIntegrationAccountAssemblySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountAssemblyArtifactName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountName(),
		},

		"assembly_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.IntegrationAccountAssemblyName(),
		},

		"assembly_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "0.0.0.0",
			ValidateFunc: validate.IntegrationAccountAssemblyVersion(),
		},

		"content": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"content", "content_link_uri"},
		},

		"content_link_uri": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			AtLeastOneOf: []string{"content", "content_link_uri"},
		},

		"metadata": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
