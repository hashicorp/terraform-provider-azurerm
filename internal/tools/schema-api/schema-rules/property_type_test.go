// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package schema_rules

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/providerschema"
)

var propertyTypeBaseNode = providerschema.SchemaJSON{
	Type:        providerschema.SchemaTypeString,
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

var propertyTypePasses = providerschema.SchemaJSON{
	Type:        providerschema.SchemaTypeString,
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

var propertyTypeViolates = providerschema.SchemaJSON{
	Type:        providerschema.SchemaTypeInt,
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

var propertyTypeList = providerschema.SchemaJSON{
	Type:        providerschema.SchemaTypeList,
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

var propertyTypeSet = providerschema.SchemaJSON{
	Type:        providerschema.SchemaTypeSet,
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

func TestPropertyTypeMatches_Check(t *testing.T) {
	data := propertyType{}
	if res := data.Check(propertyTypeBaseNode, propertyTypePasses, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}

	if res := data.Check(propertyTypeBaseNode, propertyTypeViolates, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}

	if res := data.Check(propertyTypeList, propertyTypeSet, ""); res == nil {
		t.Errorf("expected violation, but didn't get one")
	}

	if res := data.Check(propertyTypeSet, propertyTypeList, ""); res != nil {
		t.Errorf("expected no violation, got %+v", res)
	}
}
