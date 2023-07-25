// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/dave/jennifer/jen"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const (
	SchemaPath = "github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: generator-schema-snapshot <reource_type>")
	}
	rt := os.Args[1]
	res, ok := provider.AzureProvider().ResourcesMap[rt]
	if !ok {
		log.Fatalf("unknown resource type %q", rt)
	}

	f := NewFile("main")
	f.ImportName("github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk", "")

	f.Var().Id("_").Op("=").Add(SchemaMap(res.Schema))

	fmt.Printf("%#v", f)
}

func ResourceValue(res *pluginsdk.Resource) Dict {
	return Dict{
		Id("Schema"): SchemaMap(res.Schema),
	}
}

func SchemaMap(m map[string]*pluginsdk.Schema) *Statement {
	dict := Dict{}
	for k, v := range m {
		dict[Lit(k)] = Values(SchemaValue(v))
	}
	return Map(String()).Op("*").Qual(SchemaPath, "Schema").Values(dict)
}

func SchemaValue(sch *pluginsdk.Schema) Dict {
	out := Dict{}

	var t Code
	switch sch.Type {
	case pluginsdk.TypeBool:
		t = Qual(SchemaPath, "TypeBool")
	case pluginsdk.TypeInt:
		t = Qual(SchemaPath, "TypeInt")
	case pluginsdk.TypeFloat:
		t = Qual(SchemaPath, "TypeFloat")
	case pluginsdk.TypeString:
		t = Qual(SchemaPath, "TypeString")
	case pluginsdk.TypeList:
		t = Qual(SchemaPath, "TypeList")
	case pluginsdk.TypeMap:
		t = Qual(SchemaPath, "TypeMap")
	case pluginsdk.TypeSet:
		t = Qual(SchemaPath, "TypeSet")
	}
	out[Id("Type")] = t

	if sch.Required {
		out[Id("Required")] = True()
	}

	if sch.Optional {
		out[Id("Optional")] = True()
	}

	if sch.Computed {
		out[Id("Computed")] = True()
	}

	if sch.ForceNew {
		out[Id("ForceNew")] = True()
	}

	switch sch.ConfigMode {
	case pluginsdk.SchemaConfigModeAttr:
		out[Id("ConfigMode")] = Qual(SchemaPath, "SchemaConfigModeAttr")
	case pluginsdk.SchemaConfigModeBlock:
		out[Id("ConfigMode")] = Qual(SchemaPath, "SchemaConfigModeBlock")
	}

	switch sch := sch.Elem.(type) {
	case *pluginsdk.Schema:
		out[Id("Elem")] = Op("&").Qual(SchemaPath, "Schema").Values(SchemaValue(sch))
	case *pluginsdk.Resource:
		out[Id("Elem")] = Op("&").Qual(SchemaPath, "Resource").Values(ResourceValue(sch))
	}

	if sch.Set != nil {
		out[Id("Set")] = Id("TODO")
	}

	return out
}
