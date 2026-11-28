// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_definition

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/resources/mgmt/2021-06-01-preview/policy" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourcePolicyDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the policy definition.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"policy_type": {
			Description:  "The policy type.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(policy.PossibleTypeValues(), false),
		},

		"mode": {
			Description: "The policy Resource Manager mode that allows you to specify which resource types will be evaluated.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice(
				[]string{
					"All",
					"Indexed",
					"Microsoft.ContainerService.Data",
					"Microsoft.CustomerLockbox.Data",
					"Microsoft.DataCatalog.Data",
					"Microsoft.KeyVault.Data",
					"Microsoft.Kubernetes.Data",
					"Microsoft.MachineLearningServices.Data",
					"Microsoft.Network.Data",
					"Microsoft.Synapse.Data",
				}, false,
			),
		},

		"management_group_id": {
			Description: "The ID of the Management Group where this policy should be defined.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
		},

		"display_name": {
			Description: "The display name of the policy definition.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"description": {
			Description: "The description of the policy definition.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"policy_rule": {
			Description:      "The policy rule for the policy definition. This is a JSON string representing the rule that contains an if and a then block.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"parameters": {
			Description:      "Parameters for the policy definition. This field is a JSON string that allows you to parameterize your policy definition. Reducing the number of parameters forces a new resource to be created.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"role_definition_ids": {
			Description: "A list of Role Definition ID extracted from `policy_rule` required for remediation.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"metadata": helpers.MetadataSchema(),
	}
}
