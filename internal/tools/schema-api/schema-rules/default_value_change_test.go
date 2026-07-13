// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

var defaultValueChangeStringBase = providerschema.SchemaJSON{
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

var defaultValueChangeStringPasses = providerschema.SchemaJSON{
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

var defaultValueChangeStringViolates = providerschema.SchemaJSON{
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

var defaultValueChangeIntBase = providerschema.SchemaJSON{
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

var defaultValueChangeIntPasses = providerschema.SchemaJSON{
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

var defaultValueChangeIntViolates = providerschema.SchemaJSON{
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

var defaultValueChangeFloatBase = providerschema.SchemaJSON{
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

var defaultValueChangeFloatPasses = providerschema.SchemaJSON{
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

var defaultValueChangeFloatViolates = providerschema.SchemaJSON{
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
