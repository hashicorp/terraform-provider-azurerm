package network

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// setUnexportedField sets an unexported field using reflection + unsafe.
// This is test-only so we can inject RawConfig into ResourceData (to simulate unknowns during plan).
func setUnexportedField(t *testing.T, target any, fieldName string, value any) {
	t.Helper()

	v := reflect.ValueOf(target).Elem()
	f := v.FieldByName(fieldName)
	if !f.IsValid() {
		t.Fatalf("field %q not found on %T", fieldName, target)
	}

	reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).
		Elem().
		Set(reflect.ValueOf(value))
}

func TestVirtualNetwork_AddressSpaceDiffSuppress_DoesNotPanicWhenIpAddressPoolUnknown(t *testing.T) {
	schemaMap := resourceVirtualNetworkSchema()
	addressSpaceSchema := schemaMap["address_space"]
	if addressSpaceSchema == nil || addressSpaceSchema.DiffSuppressFunc == nil {
		t.Fatalf("expected address_space DiffSuppressFunc to be set")
	}

	// Minimal ResourceData required by schema initialization
	d := pluginsdk.TestResourceDataRaw(t, schemaMap, map[string]interface{}{
		"name":                "test-vnet",
		"resource_group_name": "rg",
		"location":            "westeurope",
		"address_space":       []interface{}{"10.0.0.0/16"},
	})

	// Simulate ip_address_pool being present but unknown during plan.
	rawCfg := cty.ObjectVal(map[string]cty.Value{
		"ip_address_pool": cty.UnknownVal(cty.List(cty.DynamicPseudoType)),
	})

	diff := &terraform.InstanceDiff{RawConfig: rawCfg}

	// ResourceData has an unexported "diff" field used by GetRawConfig().
	setUnexportedField(t, d, "diff", diff)

	// Before the fix this path can panic; after the fix it must not.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("DiffSuppressFunc panicked: %v", r)
		}
	}()

	suppressed := addressSpaceSchema.DiffSuppressFunc("address_space", "10.0.0.0/16", "10.0.0.0/16", d)
	if !suppressed {
		t.Fatalf("expected diff to be suppressed when ip_address_pool is unknown-but-present")
	}
}
