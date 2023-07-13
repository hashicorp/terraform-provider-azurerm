// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var defaultValueChangeStringBase = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    false,
	Default:     "foo",
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeStringPasses = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     "foo",
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeStringViolates = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    true,
	Default:     "bar",
	Description: "",
	Computed:    false, // violation
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeIntBase = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    false,
	Default:     1,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeIntPasses = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     1,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeIntViolates = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    true,
	Default:     2,
	Description: "",
	Computed:    false, // violation
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeFloatBase = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    false,
	Default:     10.0,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeFloatPasses = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    true,
	Required:    false,
	Default:     10.0,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var defaultValueChangeFloatViolates = providerjson.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    true,
	Default:     10.1,
	Description: "",
	Computed:    false, // violation
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

func TestDefaultValueChange_Check(t *testing.T) {
	data := defaultValueChange{}
	if res := data.Check(defaultValueChangeStringBase, defaultValueChangeStringPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(defaultValueChangeStringBase, defaultValueChangeStringViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
	if res := data.Check(defaultValueChangeIntBase, defaultValueChangeIntPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(defaultValueChangeIntBase, defaultValueChangeIntViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
	if res := data.Check(defaultValueChangeFloatBase, defaultValueChangeFloatPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(defaultValueChangeFloatBase, defaultValueChangeFloatViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
}
