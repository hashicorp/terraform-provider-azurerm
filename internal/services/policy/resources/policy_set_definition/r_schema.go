// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package policy_set_definition

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Policy Set Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"policy_type": {
			Description:  "The Policy Set Definition type.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(policysetdefinitions.PossibleValuesForPolicyType(), false),
		},

		"display_name": {
			Description:  "The display name of this Policy Set Definition.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description: "The description of this Policy Set Definition.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"metadata": {
			Description:      "The metadata for the Policy Set Definition in JSON format.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: helpers.PolicySetDefinitionsMetadataDiffSuppressFunc,
		},

		"parameters": {
			Description:      "The parameters for the Policy Set Definition in JSON format. Reducing the number of parameters forces a new resource to be created.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"policy_definition_reference": helpers.PolicyDefinitionReferenceSchema(),

		"policy_definition_group": helpers.PolicyDefinitionGroupSchema(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff.HasChange("parameters") {
				oldParametersRaw, newParametersRaw := metadata.ResourceDiff.GetChange("parameters")
				if oldParametersString := oldParametersRaw.(string); oldParametersString != "" {
					newParametersString := newParametersRaw.(string)
					if newParametersString == "" {
						return metadata.ResourceDiff.ForceNew("parameters")
					}

					oldParameters, err := helpers.ExpandParameterDefinitionsValue(oldParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					newParameters, err := helpers.ExpandParameterDefinitionsValue(newParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					if len(*newParameters) < len(*oldParameters) {
						return metadata.ResourceDiff.ForceNew("parameters")
					}
				}
			}

			return nil
		},
	}
}
