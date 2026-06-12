// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package compute

import "testing"

func TestDedicatedHostSkuNameValidation(t *testing.T) {
	schema := resourceDedicatedHost().Schema["sku_name"]

	if _, errors := schema.ValidateFunc("NVadsA10v5_Type1", "sku_name"); len(errors) != 0 {
		t.Fatalf("expected NVadsA10v5_Type1 to be a valid Dedicated Host SKU, got %d errors: %+v", len(errors), errors)
	}
}
