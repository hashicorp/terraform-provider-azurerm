// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

var becomeComputedOnlyOptionalBaseNode = providerschema.SchemaJSON{
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

var becomeComputedOnlyRequiredBaseNode = providerschema.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    true,
	Default:     nil,
	Description: "",
	Computed:    false,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var becomeComputedOnlyOptionalPasses = providerschema.SchemaJSON{
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

var becomeComputedOnlyRequiredPasses = providerschema.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false,
	Required:    true,
	Default:     nil,
	Description: "",
	Computed:    true,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

var becomeComputedOnlyViolates = providerschema.SchemaJSON{
	Type:        "",
	ConfigMode:  "",
	Optional:    false, // violation
	Required:    false, // violation
	Default:     nil,
	Description: "",
	Computed:    true,
	ForceNew:    false,
	Elem:        nil,
	MaxItems:    0,
	MinItems:    0,
}

func TestBecomeComputedOnly_Check(t *testing.T) {
	data := becomeComputedOnly{}
	if res := data.Check(becomeComputedOnlyOptionalBaseNode, becomeComputedOnlyOptionalPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(becomeComputedOnlyRequiredBaseNode, becomeComputedOnlyRequiredPasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
	if res := data.Check(becomeComputedOnlyOptionalBaseNode, becomeComputedOnlyViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
	if res := data.Check(becomeComputedOnlyRequiredBaseNode, becomeComputedOnlyViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}
}
