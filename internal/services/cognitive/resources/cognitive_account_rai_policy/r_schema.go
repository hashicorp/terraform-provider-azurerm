// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cognitive_account_rai_policy

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/raipolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			rawFilters, ok := metadata.ResourceDiff.GetOk("content_filter")
			if !ok {
				return nil
			}

			filters, ok := rawFilters.([]any)
			if !ok {
				return nil
			}

			for i, rawFilter := range filters {
				filter, ok := rawFilter.(map[string]any)
				if !ok {
					continue
				}

				name, _ := filter["name"].(string)
				severityThreshold, _ := filter["severity_threshold"].(string)

				if severityThreshold == "" {
					continue
				}

				if slices.Contains(severityThresholdNotApplicableFilterNames, name) {
					return fmt.Errorf("`severity_threshold` is not applicable for `content_filter[%d]` with name %q", i, name)
				}
			}

			return nil
		},
	}
}

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Cognitive Service Account RAI Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Description:  "The ID of the Cognitive Service Account to which this RAI Policy should be associated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: raipolicies.ValidateAccountID,
		},

		"base_policy_name": {
			Description:  "The name of the base policy to use for this RAI Policy.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"content_filter": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the content filter.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"block_enabled": {
						Description: "Whether the filter should block content.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"filter_enabled": {
						Description: "Whether the filter is enabled.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},
					"source": {
						Description:  "Content source to apply the content filter.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForRaiPolicyContentSource(), false),
					},
					"severity_threshold": {
						Description:  "The severity threshold for the filter.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForContentLevel(), false),
					},
				},
			},
		},

		"mode": {
			Description:  "The mode of the RAI Policy.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForRaiPolicyMode(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
