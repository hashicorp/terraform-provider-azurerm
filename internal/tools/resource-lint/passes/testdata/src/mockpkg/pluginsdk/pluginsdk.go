// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This file is to mock pluginsdk in azurerm

const (
	TypeString = schema.TypeString
	TypeBool   = schema.TypeBool
	TypeInt    = schema.TypeInt
	TypeMap    = schema.TypeMap
	TypeList   = schema.TypeList
)

type (
	Resource     = schema.Resource
	Schema       = schema.Schema
	ResourceData = schema.ResourceData
)
