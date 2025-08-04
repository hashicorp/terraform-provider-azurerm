// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonschema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ZoneSingleRequired returns the schema used when a single Zone must be specified
func ZoneSingleRequired() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

// ZoneSingleRequiredForceNew returns the schema used when a single Zone must be specified but cannot be changed
func ZoneSingleRequiredForceNew() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

// ZoneSingleOptional returns the schema used when a single Zone can be specified
func ZoneSingleOptional() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

// ZoneSingleOptionalForceNew returns the schema used when a single Zone can be specified but cannot be changed
func ZoneSingleOptionalForceNew() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

// ZoneSingleComputed returns the schema used when a single Zones can be returned
func ZoneSingleComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

// ZoneSingleOptionalComputed returns the schema used when a single Zone can be specified or a single Zone is returned when omitted
func ZoneSingleOptionalComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
}
