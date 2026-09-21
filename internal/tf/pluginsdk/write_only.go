// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package pluginsdk

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/ctyhelpers"
)

// GetWriteOnly gets a write only attribute, checking that it is of an expected type and subsequently returns it
func GetWriteOnly(d *ResourceData, attributePath string, attributeType cty.Type) (*cty.Value, error) {
	value, diags := d.GetRawConfigAt(ctyhelpers.ConstructCtyPath(attributePath))
	if diags.HasError() {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: %+v", attributePath, diags)
	}

	if !value.Type().Equals(attributeType) {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: value is not of type %v", attributePath, attributeType)
	}
	return pointer.To(value), nil
}

// GetWriteOnlyFromDiff gets a write only attribute from the diff, checking that it is of an expected type and subsequently returns it
func GetWriteOnlyFromDiff(d *ResourceDiff, attributePath string, attributeType cty.Type) (*cty.Value, error) {
	value, diags := d.GetRawConfigAt(ctyhelpers.ConstructCtyPath(attributePath))
	if diags.HasError() {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: %+v", attributePath, diags)
	}

	if !value.Type().Equals(attributeType) {
		return nil, fmt.Errorf("retrieving write-only attribute `%s`: value is not of type %v", attributePath, attributeType)
	}
	return pointer.To(value), nil
}
