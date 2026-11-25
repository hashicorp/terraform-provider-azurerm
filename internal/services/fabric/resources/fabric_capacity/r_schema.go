// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fabric_capacity

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for the Fabric Capacity.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z]([a-z\d]{2,62})$`),
				"`name` must be between 3 and 63 characters. It can contain only lowercase letters and numbers. It must start with a lowercase letter.",
			),
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "The name of the SKU to use for the Fabric Capacity.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"F0",
							"F2",
							"F4",
							"F8",
							"F16",
							"F32",
							"F64",
							"F128",
							"F256",
							"F512",
							"F1024",
							"F2048",
							"F4096",
							"F8192",
						}, false),
					},
					"tier": {
						Description:  "The tier of the SKU to use for the Fabric Capacity.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(fabriccapacities.PossibleValuesForRpSkuTier(), false),
					},
				},
			},
		},

		"administration_members": {
			Description: "An array of administrator user identities. The member must be an Entra user or a Service Principal.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
