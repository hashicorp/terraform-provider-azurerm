// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_self_hosted

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryIntegrationRuntimeSelfHostedSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[A-Za-z0-9]+(?:-[A-Za-z0-9]+)*$`),
				`Invalid name for Self-Hosted Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
			),
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"rbac_authorization": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Set:      resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: integrationruntimes.ValidateIntegrationRuntimeID,
					},
				},
			},
		},

		"self_contained_interactive_authoring_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"primary_authorization_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_authorization_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
