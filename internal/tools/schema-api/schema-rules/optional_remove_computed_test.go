// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var optionalRemoveComputedBase = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     nil,
	Description: "",
	Computed:    true,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var optionalRemoveComputedPasses = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     nil,
	Description: "",
	Computed:    true,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var optionalRemoveComputedViolates = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     nil,
	Description: "",
	Computed:    false, // violation
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

func TestOptionalRemoveComputed_Check(t *testing.T) {
	data := optionalRemoveComputed{}
	if res := data.Check(optionalRemoveComputedBase, optionalRemoveComputedPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(optionalRemoveComputedBase, optionalRemoveComputedViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
}
