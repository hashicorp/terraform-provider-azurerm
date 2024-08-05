// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

// FileForResource for typed sdk resource, the file is terraform-provider-azurerm/internal/sdk/wrapper_resource.go
func FileForResource(funcs ...interface{}) (file string) {
	for _, fn := range funcs {
		if file, _ = util.FuncFileLine(fn); file != "" {
			return file
		}
	}
	return
}

type Resource struct {
	FilePath     string
	ResourceType string // azurerm_xxx

	// one of Schema or SDKResource must use
	Schema      *schema.Resource `json:"-"`
	SDKResource sdk.Resource     `json:"-"`

	PossibleValues map[string][]string // possible values for key(property path)
}

func ResourceForSDKType(res sdk.Resource) *schema.Resource {
	r := sdk.NewResourceWrapper(res)
	ins, _ := r.Resource()
	return ins
}

// NewResourceByTyped NewResource ...
// r is Schema.Resource or Typed SDK Resource
func NewResourceByTyped(r sdk.Resource) *Resource {
	s := &Resource{}
	s.SDKResource = r
	s.Schema = ResourceForSDKType(r)
	s.ResourceType = r.ResourceType()
	s.Init()
	return s
}

func NewResourceByUntyped(r *schema.Resource, rType string) *Resource {
	s := &Resource{}
	s.Schema = r
	s.ResourceType = rType
	s.Init()
	return s
}

func NewResource(r interface{}, rType string) *Resource {
	switch ins := r.(type) {
	case sdk.Resource:
		return NewResourceByTyped(ins)
	case *schema.Resource:
		return NewResourceByUntyped(ins, rType)
	}
	return nil
}

func (r *Resource) Init() {
	if r.SDKResource != nil {
		// SDKResource is a type of interface, have to get the real
		// vd := reflect.ValueOf(r.SDKResource).Interface()
		// vd = reflect.ValueOf(vd).MethodByName("Arguments")
		// this is not work if Read() defined in other file
		r.FilePath = FileForResource(r.SDKResource.Read().Func)
	} else {
		r.FilePath = FileForResource(r.Schema.Read, r.Schema.ReadContext) //nolint:staticcheck
	}
	r.PossibleValues = map[string][]string{}
	r.FindAllInSlicePropByMonkey()
}

func (r *Resource) FilePathRel() string {
	if idx := strings.Index(r.FilePath, "internal"); idx > 0 {
		return "./" + r.FilePath[idx:]
	}
	return r.FilePath
}

func (r *Resource) HasPathFor(path []string) bool {
	sch := r.Schema.Schema
	for _, part := range path {
		ele, ok := sch[part]
		if !ok {
			return false
		}
		if ele.Elem != nil {
			if res, ok := ele.Elem.(*schema.Resource); ok {
				sch = res.Schema
			} else {
				sch = nil
			}
		}
	}
	return true
}

func (r *Resource) IsDeprecated() bool {
	if r.Schema != nil {
		return r.Schema.DeprecationMessage != ""
	} else if _, ok := r.SDKResource.(sdk.ResourceWithDeprecationReplacedBy); ok {
		return true
	}

	return false
}
