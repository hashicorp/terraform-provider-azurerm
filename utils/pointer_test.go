package utils_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func TestValue(t *testing.T) {
	s := []struct {
		ptr  *string
		want string
	}{
		{
			nil,
			"",
		},
		{
			utils.Ptr("abc"),
			"abc",
		},
		{
			utils.Ptr(""),
			"",
		},
	}
	for _, arg := range s {
		if arg.want != utils.Value(arg.ptr) {
			t.Fatal(arg)
		}
	}
	v := utils.Value((*map[string]string)(nil))
	if v != nil {
		t.Fatal(v)
	}
}

func TestTryPtr(t *testing.T) {
	args := [][2]interface{}{
		{"abc", "abc"},
		{nil, nil},
		{(map[string]string)(nil), nil},
		{(*int)(nil), nil},
	}
	for _, arg := range args {
		ptr := utils.TryPtr(arg[0])
		//t.Logf("arg: %+v, ptr: %v", arg, ptr)
		if ptr == nil {
			if arg[1] != nil {
				t.Fatal(arg, ptr)
			}
			continue
		}
		if *ptr != arg[1] {
			t.Fatal(arg, ptr)
		}
	}
}
