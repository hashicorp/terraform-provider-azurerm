// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourcegroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceGroupName() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: resourcegroups.ValidateName,
	}
}

func ResourceGroupNameDeprecated() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: resourcegroups.ValidateName,
		Deprecated:   "This field is no longer used and will be removed in the next major version of the Azure Provider",
	}
}

func ResourceGroupNameDeprecatedComputed() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: resourcegroups.ValidateName,
		Deprecated:   "This field is no longer used and will be removed in the next major version of the Azure Provider",
	}
}

func ResourceGroupNameForDataSource() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: resourcegroups.ValidateName,
	}
}

func ResourceGroupNameOptionalComputed() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		ForceNew:     true,
		Optional:     true,
		Computed:     true,
		ValidateFunc: resourcegroups.ValidateName,
	}
}

func ResourceGroupNameOptional() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: resourcegroups.ValidateName,
	}
}

func ResourceGroupNameSetOptional() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: resourcegroups.ValidateName,
		},
	}
}
