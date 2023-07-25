// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var optionalToRequiredBaseNode = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     nil,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var optionalToRequiredPasses = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     nil,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var optionalToRequiredViolates = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false, // violation
	Required:    true,  // violation
	Default:     nil,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

func TestOptionalToRequired_Check(t *testing.T) {
	data := optionalToRequired{}
	if res := data.Check(optionalToRequiredBaseNode, optionalToRequiredPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(optionalToRequiredBaseNode, optionalToRequiredViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
}
