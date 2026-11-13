// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package powerbi_embedded

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/powerbidedicated/2021-01-01/capacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func powerbiEmbeddedSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.EmbeddedName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"A1",
				"A2",
				"A3",
				"A4",
				"A5",
				"A6",
				"A7",
				"A8",
			}, false),
		},

		"administrators": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validate.EmbeddedAdministratorName,
			},
		},

		"mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(capacities.ModeGenTwo),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(capacities.PossibleValuesForMode(), false),
		},

		"tags": commonschema.Tags(),
	}
}
