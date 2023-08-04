// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationsmanagement/2015-11-01-preview/solution"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = SolutionV0ToV1{}

type SolutionV0ToV1 struct{}

func (SolutionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId, err := solution.ParseSolutionIDInsensitively(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}

func (SolutionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"solution_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workspace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LogAnalyticsWorkspaceName,
		},

		"workspace_resource_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"location": commonschema.Location(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"plan": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"publisher": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
					"promotion_code": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},
					"product": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}
