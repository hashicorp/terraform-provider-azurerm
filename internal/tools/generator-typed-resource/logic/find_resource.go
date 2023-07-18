package logic

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// this file is to find out all untyped resources, name and resource code file path

type untypedResource struct {
	name     string
	path     string
	funcName string // function name of resource definition
	resource *pluginsdk.Resource
}

func allDefineFuncs() map[string]string {
	_, cur, _, _ := runtime.Caller(0)
	dirPath := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(cur)))), "services")
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("[Error] read dir: %s: %v", dirPath, err)
		return nil
	}
	res := map[string]string{}
	for _, ent := range dir {
		if ent.IsDir() {
			part := resourceDefineNames(filepath.Join(dirPath, ent.Name(), "registration.go"))
			for k, v := range part {
				res[k] = v
			}
		}
	}
	return res
}

// load all mappings from registration.go
func resourceDefineNames(file string) map[string]string {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, file, nil, 0)
	if err != nil {
		log.Printf("[Error] ast parse %s: %v", file, err)
		return nil
	}
	defineNames, dsNames := map[string]string{}, map[string]string{}
	ast.Inspect(f, func(node ast.Node) bool {
		if decl, ok := node.(*ast.FuncDecl); ok {
			declName := decl.Name.Name
			if decl.Name.Name != "SupportedResources" && decl.Name.Name != "SupportedDataSources" {
				return false
			}

			ast.Inspect(decl.Body, func(node ast.Node) bool {
				if comp, ok := node.(*ast.CompositeLit); ok {
					for _, elt := range comp.Elts {
						if kv, ok := elt.(*ast.KeyValueExpr); ok {
							if key, ok := kv.Key.(*ast.BasicLit); ok {
								name := basicString(key)
								if val, ok := kv.Value.(*ast.CallExpr); ok {
									if funcName, ok := val.Fun.(*ast.Ident); ok {
										if declName == "SupportedResources" {
											defineNames[name] = funcName.Name
										} else {
											dsNames[name] = funcName.Name
										}
									}
								}
							}
						}
					}
					return false
				}
				return true
			})
		}
		return true
	})
	return defineNames
}

// loop all files
func loadAllUntypedResources() []*untypedResource {
	var rss []*untypedResource
	defineNames := allDefineFuncs()
	for _, reg := range provider.SupportedUntypedServices() {
		for name, r := range reg.SupportedResources() {
			pc := reflect.ValueOf(r.Read).Pointer() //nolint:staticcheck
			file, _ := runtime.FuncForPC(pc).FileLine(pc)
			rss = append(rss, &untypedResource{
				name:     name,
				path:     file,
				funcName: defineNames[name],
				resource: r,
			})
		}
	}
	return rss
}
