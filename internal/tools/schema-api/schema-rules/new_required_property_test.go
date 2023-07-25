// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/schema-api/providerjson"
)

var newRequiredPropertyBase = providerjson.SchemaJSON{
	Type:        "", // empty here indicates this doesn't exist in the base resource
	ConfigMode:  "",
	Optional:    false,
	Required:    false,
	Default:     nil,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var newRequiredPropertyPasses = providerjson.SchemaJSON{
	Type:        providerjson.SchemaTypeString,
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

var newRequiredPropertyViolates = providerjson.SchemaJSON{
	Type:        providerjson.SchemaTypeString,
	ConfigMode:  "",
	Optional:    false,
	Required:    true,
	Default:     nil,
	Description: "",
	Computed:    false, // violation
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

func TestNewRequiredProperty_Check(t *testing.T) {
	data := newRequiredPropertyExistingResource{}
	if res := data.Check(newRequiredPropertyBase, newRequiredPropertyPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(newRequiredPropertyBase, newRequiredPropertyViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
}
