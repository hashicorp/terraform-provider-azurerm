// Copyright IBM Corp. 2014, 2026
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
			Description: "The Name which should be used for this Arc Resource Bridge Appliance.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				validation.StringLenBetween(1, 260),
				validation.StringMatch(regexp.MustCompile(`[^+#%&'?/,%\\]+$`), "any of '+', '#', '%', '&', ''', '?', '/', ',', '%', '&', '\\', are not allowed"),
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"distro": {
			Description:  "A supported Fabric/Infrastructure for this Arc Resource Bridge Appliance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(appliances.PossibleValuesForDistro(), false),
		},

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"infrastructure_provider": {
			Description:  "The infrastructure provider about the connected Arc Resource Bridge Appliance.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(appliances.PossibleValuesForProvider(), false),
		},

		"public_key_base64": {
			Description:  "The `public_key_base64` is an RSA public key in PKCS1 format encoded in base64.",
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
