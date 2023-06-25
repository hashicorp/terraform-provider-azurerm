// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TagsDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func TagsForceNew() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: tags.Validate,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func Tags() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		ValidateFunc: tags.Validate,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func TagsWithLowerCaseKeys() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeMap,
		Optional:     true,
		ValidateFunc: tags.ValidateHasLowerCaseKeys,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}
