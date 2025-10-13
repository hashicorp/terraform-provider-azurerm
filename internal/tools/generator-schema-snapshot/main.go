// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dave/jennifer/jen"
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

	f := jen.NewFile("main")
	f.ImportName("github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk", "")

	f.Var().Id("_").Op("=").Add(SchemaMap(res.Schema))

	fmt.Printf("%#v", f)
}

func ResourceValue(res *pluginsdk.Resource) jen.Dict {
	return jen.Dict{
		jen.Id("Schema"): SchemaMap(res.Schema),
	}
}

func SchemaMap(m map[string]*pluginsdk.Schema) *jen.Statement {
	dict := jen.Dict{}
	for k, v := range m {
		dict[jen.Lit(k)] = jen.Values(SchemaValue(v))
	}
	return jen.Map(jen.String()).Op("*").Qual(SchemaPath, "Schema").Values(dict)
}

func SchemaValue(sch *pluginsdk.Schema) jen.Dict {
	out := jen.Dict{}

	var t jen.Code
	switch sch.Type {
	case pluginsdk.TypeBool:
		t = jen.Qual(SchemaPath, "TypeBool")
	case pluginsdk.TypeInt:
		t = jen.Qual(SchemaPath, "TypeInt")
	case pluginsdk.TypeFloat:
		t = jen.Qual(SchemaPath, "TypeFloat")
	case pluginsdk.TypeString:
		t = jen.Qual(SchemaPath, "TypeString")
	case pluginsdk.TypeList:
		t = jen.Qual(SchemaPath, "TypeList")
	case pluginsdk.TypeMap:
		t = jen.Qual(SchemaPath, "TypeMap")
	case pluginsdk.TypeSet:
		t = jen.Qual(SchemaPath, "TypeSet")
	}
	out[jen.Id("Type")] = t

	if sch.Required {
		out[jen.Id("Required")] = jen.True()
	}

	if sch.Optional {
		out[jen.Id("Optional")] = jen.True()
	}

	if sch.Computed {
		out[jen.Id("Computed")] = jen.True()
	}

	switch sch.ConfigMode {
	case pluginsdk.SchemaConfigModeAttr:
		out[jen.Id("ConfigMode")] = jen.Qual(SchemaPath, "SchemaConfigModeAttr")
	case pluginsdk.SchemaConfigModeBlock:
		out[jen.Id("ConfigMode")] = jen.Qual(SchemaPath, "SchemaConfigModeBlock")
	}

	switch sch := sch.Elem.(type) {
	case *pluginsdk.Schema:
		out[jen.Id("Elem")] = jen.Op("&").Qual(SchemaPath, "Schema").Values(SchemaValue(sch))
	case *pluginsdk.Resource:
		out[jen.Id("Elem")] = jen.Op("&").Qual(SchemaPath, "Resource").Values(ResourceValue(sch))
	}

	if sch.Set != nil {
		out[jen.Id("Set")] = jen.Id("TODO")
	}

	return out
}
