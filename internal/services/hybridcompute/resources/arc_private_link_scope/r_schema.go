// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (a Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"location":            commonschema.Location(),
		"public_network_access_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"tags": commonschema.Tags(),
	}
}

func (a Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
