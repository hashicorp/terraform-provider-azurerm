// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package fabric_capacity

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func TestFabricCapacitySkuNameValidation(t *testing.T) {
	sku := Resource{}.Arguments()["sku"]
	elem, ok := sku.Elem.(*pluginsdk.Resource)
	if !ok {
		t.Fatalf("expected sku element to be a resource")
	}

	validate := elem.Schema["name"].ValidateFunc
	for _, name := range []string{"F0", "F2", "F2048"} {
		if _, errs := validate(name, "sku.0.name"); len(errs) != 0 {
			t.Errorf("expected %q to be valid, got %v", name, errs)
		}
	}

	if _, errs := validate("F1", "sku.0.name"); len(errs) == 0 {
		t.Errorf("expected F1 to be invalid")
	}
}
