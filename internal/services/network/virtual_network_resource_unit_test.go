package network

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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

	d := schema.TestResourceDataRaw(t, schemaMap, map[string]interface{}{
		"name":                "test-vnet",
		"resource_group_name": "rg",
		"location":            "westeurope",
		"address_space":       []interface{}{"10.0.0.0/16"},
	})

	rawCfg := cty.ObjectVal(map[string]cty.Value{
		"ip_address_pool": cty.UnknownVal(cty.List(cty.DynamicPseudoType)),
	})

	diff := &terraform.InstanceDiff{RawConfig: rawCfg}

	setUnexportedField(t, d, "diff", diff)

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
