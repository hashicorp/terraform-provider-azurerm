// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_azure

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryIntegrationRuntimeAzureSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
				`Invalid name for Azure Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
			),
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				location.EnhancedValidate,
				validation.StringInSlice([]string{"AutoResolve"}, false),
			),
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"cleanup_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"compute_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(integrationruntimes.DataFlowComputeTypeGeneral),
			ValidateFunc: validation.StringInSlice(integrationruntimes.PossibleValuesForDataFlowComputeType(), false),
		},

		"core_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  8,
			ValidateFunc: validation.IntInSlice([]int{
				8, 16, 32, 48, 80, 144, 272,
			}),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"interactive_authoring_time_to_live_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntInSlice([]int{10, 30, 60, 120}),
		},

		"time_to_live_min": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  0,
		},

		"virtual_network_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},
	}
}
