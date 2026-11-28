// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_integration_runtime_azure

import (
	"regexp"

	"github.com/Azure/azure-sdk-for-go/services/preview/synapse/mgmt/v2.0/synapse" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseIntegrationRuntimeAzureSchema() map[string]*pluginsdk.Schema {
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

		"synapse_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
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

		"compute_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(synapse.DataFlowComputeTypeGeneral),
			ValidateFunc: validation.StringInEnumSlice(synapse.PossibleDataFlowComputeTypeValues(), false),
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
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"time_to_live_min": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  0,
		},
	}
}
