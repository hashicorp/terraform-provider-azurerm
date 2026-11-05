// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_resource_bridge_appliance

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 260),
				validation.StringMatch(regexp.MustCompile(`[^+#%&'?/,%\\]+$`), "any of '+', '#', '%', '&', ''', '?', '/', ',', '%', '&', '\\', are not allowed"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"distro": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(appliances.PossibleValuesForDistro(), false),
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"infrastructure_provider": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(appliances.PossibleValuesForProvider(), false),
		},

		"public_key_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsBase64,
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}
