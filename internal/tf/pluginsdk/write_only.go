// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-cty/cty"
)

// GetWriteOnly gets a write only attribute, checking that it is of an expected type and subsequently returns it
func GetWriteOnly(d *ResourceData, name string, attributeType cty.Type) (*cty.Value, error) {
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(name))
	if diags.HasError() {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: %+v", name, diags)
	}

	if !value.Type().Equals(attributeType) {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: value is not of type %v", name, attributeType)
	}
	return pointer.To(value), nil
}

// GetWriteOnlyFromDiff gets a write only attribute from the diff, checking that it is of an expected type and subsequently returns it
func GetWriteOnlyFromDiff(d *ResourceDiff, name string, attributeType cty.Type) (*cty.Value, error) {
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(name))
	if diags.HasError() {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: %+v", name, diags)
	}

	if !value.Type().Equals(attributeType) {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: value is not of type %v", name, attributeType)
	}
	return pointer.To(value), nil
}
