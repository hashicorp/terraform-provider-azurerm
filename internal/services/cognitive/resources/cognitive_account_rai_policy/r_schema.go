// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: raipolicies.ValidateAccountID,
		},

		"base_policy_name": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"block_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"filter_enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
					"source": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForRaiPolicyContentSource(), false),
					},
					"severity_threshold": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(raipolicies.PossibleValuesForContentLevel(), false),
					},
				},
			},
		},

		"mode": {
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
