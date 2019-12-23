package tests

import (
	"strings"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
)

func TestValidateArmStorageShareName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
	}
	for _, v := range validNames {
		_, errors := storage.ValidateArmStorageShareName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Share Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"InvalidName1",
		"-invalidname1",
		"invalid_name",
		"invalid!",
		"double-hyphen--invalid",
		"ww",
		strings.Repeat("w", 65),
	}
	for _, v := range invalidNames {
		if _, errors := storage.ValidateArmStorageShareName(v, "name"); len(errors) == 0 {
			t.Fatalf("%q should be an invalid Share Name", v)
		}
	}
}
