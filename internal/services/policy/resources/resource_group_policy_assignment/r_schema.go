// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resource_group_policy_assignment

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringIsNotWhiteSpace,
				validation.StringDoesNotContainAny("#<>%&:\\?/"),
				validation.StringLenBetween(1, 64),
				validation.StringMatch(regexp.MustCompile("[^ .]$"), "The name cannot end with a period or space."),
			),
		},
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(commonids.ParseResourceGroupIDInsensitively),
		},
	}
	return r.base.Arguments(schema)
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return r.base.Attributes()
}
