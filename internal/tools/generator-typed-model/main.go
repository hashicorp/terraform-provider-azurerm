// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: generator-typed-model <resource_type>")
		os.Exit(1)
	}
	rt := flag.Args()[0]
	resource, ok := provider.AzureProvider().ResourcesMap[rt]
	if !ok {
		log.Fatalf("unknown resource type: %s", rt)
	}

	f := jen.NewFile("main")
	modelStmts := modelForSchemaMap(snake2Camel(strings.TrimPrefix(rt, "azurerm_"))+"Model", resource.Schema)
	for _, stmt := range modelStmts {
		stmt := stmt
		f.Add(&stmt)
	}
	fmt.Printf("%#v", f)
}

func snake2Camel(input string) string {
	segs := strings.Split(input, "_")
	var out string
	for _, seg := range segs {
		if seg == "" {
			continue
		}
		out += strings.ToUpper(string(seg[0])) + seg[1:]
	}
	return out
}

func modelForSchemaMap(name string, sm map[string]*schema.Schema) []jen.Statement {
	var out []jen.Statement

	var thisStmt jen.Statement

	fields := []jen.Code{}

	keys := []string{}
	for k := range sm {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		sch := sm[key]
		fieldName := snake2Camel(key)
		tag := map[string]string{"tfschema": key}

		switch sch.Type {
		case schema.TypeBool:
			fields = append(fields, jen.Id(fieldName).Bool().Tag(tag))
		case schema.TypeInt:
			fields = append(fields, jen.Id(fieldName).Int().Tag(tag))
		case schema.TypeString:
			fields = append(fields, jen.Id(fieldName).String().Tag(tag))
		case schema.TypeFloat:
			fields = append(fields, jen.Id(fieldName).Float64().Tag(tag))
		case schema.TypeList,
			schema.TypeSet:
			field := jen.Id(fieldName).Index()

			switch elemSch := sch.Elem.(type) {
			case *schema.Resource:
				typeName := fieldName + "Model"
				out = append(out, modelForSchemaMap(typeName, elemSch.Schema)...)
				fields = append(fields, field.Id(typeName).Tag(tag))
			case *schema.Schema:
				switch elemSch.Type {
				case schema.TypeBool:
					fields = append(fields, field.Bool().Tag(tag))
				case schema.TypeInt:
					fields = append(fields, field.Int().Tag(tag))
				case schema.TypeString:
					fields = append(fields, field.String().Tag(tag))
				case schema.TypeFloat:
					fields = append(fields, field.Float64().Tag(tag))
				default:
					panic(fmt.Errorf("unhandled type: List/Set of Schema of %s", elemSch.Type))
				}
			default:
				panic(fmt.Errorf("unhandled type: List/Set of %t", sch.Elem))
			}
		case schema.TypeMap:
			field := jen.Id(fieldName).Map(jen.String())
			// Map's element must be of type *schema.Schema
			elemSch := sch.Elem.(*schema.Schema)
			switch elemSch.Type {
			case schema.TypeBool:
				fields = append(fields, field.Bool().Tag(tag))
			case schema.TypeInt:
				fields = append(fields, field.Int().Tag(tag))
			case schema.TypeString:
				fields = append(fields, field.String().Tag(tag))
			case schema.TypeFloat:
				fields = append(fields, field.Float64().Tag(tag))
			default:
				panic(fmt.Errorf("unhandled type: Map of %s", elemSch.Type))
			}
		default:
			panic(fmt.Errorf("unhandled type: %s", sch.Type))
		}
	}
	thisStmt = *jen.Type().Id(name).Struct(fields...)

	out = append(out, thisStmt)

	return out
}
