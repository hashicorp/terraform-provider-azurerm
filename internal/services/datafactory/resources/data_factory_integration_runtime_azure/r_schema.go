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
			Description: "The name of the Managed Integration Runtime. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^([a-zA-Z0-9](-|-?[a-zA-Z0-9]+)+[a-zA-Z0-9])$`),
				`Invalid name for Azure Integration Runtime: minimum 3 characters, must start and end with a number or a letter, may only consist of letters, numbers and dashes and no consecutive dashes.`,
			),
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Linked Service with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"location": {
			Description: "The supported Azure location where the resource exists. Use `AutoResolve` to create an auto-resolve integration runtime.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				location.EnhancedValidate,
				validation.StringInSlice([]string{"AutoResolve"}, false),
			),
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
		},

		"cleanup_enabled": {
			Description: "Cluster will not be recycled and it will be used in next data flow activity run until TTL (time to live) is reached if this is set as `false`.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"compute_type": {
			Description:  "Compute type of the cluster which will execute data flow job.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(integrationruntimes.DataFlowComputeTypeGeneral),
			ValidateFunc: validation.StringInSlice(integrationruntimes.PossibleValuesForDataFlowComputeType(), false),
		},

		"core_count": {
			Description: "Core count of the cluster which will execute data flow job.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     8,
			ValidateFunc: validation.IntInSlice([]int{
				8, 16, 32, 48, 80, 144, 272,
			}),
		},

		"description": {
			Description: "Integration runtime description.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"interactive_authoring_time_to_live_in_minutes": {
			Description:  "The time to live in minutes for the interactive authoring feature, setting this enables interactive authoring.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntInSlice([]int{10, 30, 60, 120}),
		},

		"time_to_live_min": {
			Description: "Time to live (in minutes) setting of the cluster which will execute data flow job.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     0,
		},

		"virtual_network_enabled": {
			Description: "Whether Integration Runtime compute is provisioned within Managed Virtual Network.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},
	}
}
