// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"reflect"
	"runtime"
	"strings"

	gomonkey "github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func patchPossibleValuesFn() {
	gomonkey.ApplyFunc(validation.StringInSlice,
		func(valid []string, ignoreCase bool) schema.SchemaValidateFunc { //nolint:staticcheck
			return func(i interface{}, k string) (warnings []string, errors []error) {
				var res []string // must have a copy
				res = append(res, valid...)
				return res, nil
			}
		})
}

func init() {
	patchPossibleValuesFn()
}

func (r *Resource) FindAllInSlicePropByMonkey() {
	for name, item := range r.Schema.Schema {
		r.InSlicePropByMonkey(name, item)
	}
}

func (r *Resource) InSlicePropByMonkey(name string, item *schema.Schema) {
	if item.ValidateFunc != nil {
		// check if it is StringsInSlice
		pc := reflect.ValueOf(item.ValidateFunc).Pointer()
		fn := runtime.FuncForPC(pc).Name()
		// ValidateFunc may direct use sdk v2.StringInSlice then function name will be patchPossibleValuesFn
		if strings.Contains(fn, "patchPossibleValuesFn") || strings.Contains(fn, "StringInSlice") {
			values, _ := item.ValidateFunc(nil, "")
			r.PossibleValues[name] = values
		}
	}
	switch ele := item.Elem.(type) {
	case *schema.Resource:
		for subName, prop := range ele.Schema {
			r.InSlicePropByMonkey(name+"."+subName, prop)
		}
	case *schema.Schema:
		r.InSlicePropByMonkey(name, ele)
	}
}
