// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_policy

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/policies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func devTestPolicySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Dev Test Policy.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(policies.PolicyFactNameGalleryImage),
				string(policies.PolicyFactNameLabPremiumVMCount),
				string(policies.PolicyFactNameLabTargetCost),
				string(policies.PolicyFactNameLabVMCount),
				string(policies.PolicyFactNameLabVMSize),
				string(policies.PolicyFactNameUserOwnedLabPremiumVMCount),
				string(policies.PolicyFactNameUserOwnedLabVMCount),
				string(policies.PolicyFactNameUserOwnedLabVMCountInSubnet),
			}, false),
		},

		"policy_set_name": {
			Description: "The name of the Policy Set within the Dev Test Lab where this policy should be created.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"lab_name": {
			Description:  "The name of the Dev Test Lab in which the Policy should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DevTestLabName(),
		},

		// There's a bug in the Azure API where this is returned in lower-case
		// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"threshold": {
			Description:  "The Threshold for this Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"evaluator_type": {
			Description:  "The Evaluation Type used for this Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(policies.PossibleValuesForPolicyEvaluatorType(), false),
		},

		"description": {
			Description: "A description for the Policy.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"fact_data": {
			Description: "The Fact Data for this Policy.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"tags": commonschema.Tags(),
	}
}
