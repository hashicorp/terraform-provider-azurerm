// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
	"golang.org/x/tools/go/ast/astutil"
)

type operation string

const (
	Register   operation = "register"
	Unregister operation = "unregister"
)

type resourceType string

const (
	TypeResource   resourceType = "resource"
	TypeDatasource resourceType = "datasource"
)

var resourceNameToAdd string

// UpdateRegistration modifies the `registration.go` file to add (or remove?) items from the relevant blocks when
// a user adds (or removes?) a resource or datasource via `scaff resource` or `scaff datasource`
//
// servicePackagePath = The path from the route of the provider to the service package to which the item is to be registered
// resourceName = the model name of a typed resource, or the snakeCase name of the item
// resource = string value for the item type, one of TypeResource or TypeDatasource
// op = one of Register or Unregister // TODO - removing is a future concern so not yet implemented
// isTyped = true if the resources uses the TypedSDK
func UpdateRegistration(servicePackagePath string, resourceName string, resource resourceType, _ operation, isTyped bool) error {
	resourceNameToAdd = getResourceNameToAdd(resourceName, resource, isTyped)
	fSet := token.NewFileSet()
	regFilePath := fmt.Sprintf("%s/registration.go", strings.TrimSuffix(servicePackagePath, "/"))
	regFile, err := parser.ParseFile(fSet, regFilePath, nil, 0)
	if err != nil {
		return err
	}

	nodeName := normaliseNodeName(resource, isTyped)
	ast.Inspect(regFile, func(node ast.Node) bool {
		fn, ok := node.(*ast.FuncDecl)
		if ok {
			if fn.Name.Name == nodeName {
				ast.Inspect(fn.Body.List[0], func(r ast.Node) bool {
					ret, ok := r.(*ast.ReturnStmt)
					if ok {
						ast.Inspect(ret.Results[0], func(n ast.Node) bool {
							if out, ok := n.(*ast.CompositeLit); ok {
								if isTyped {
									astutil.Apply(out, typedAppendResourceToRegistrationBlock(), nil)
								} else {
									astutil.Apply(out, untypedAppendResourceToRegistrationBlock(), nil)
								}
							}
							return true
						})
					}
					return true
				})
			}
		}
		return true
	})

	outBuf := new(bytes.Buffer)
	if err := format.Node(outBuf, fSet, regFile); err != nil {
		return err
	}

	if err := os.WriteFile(regFilePath, outBuf.Bytes(), 0755); err != nil {
		return err
	}

	return nil
}

func normaliseNodeName(input resourceType, isTyped bool) string {
	switch input {
	case TypeDatasource:
		if isTyped {
			return "DataSources"
		} else {
			return "SupportedDataSources"
		}
	default:
		if isTyped {
			return "Resources"
		}
	}
	return "SupportedResources"
}

func getResourceNameToAdd(name string, res resourceType, isTyped bool) string {
	if strings.EqualFold(string(res), "datasource") {
		if isTyped {
			return name + "DataSource"
		}
	} else {
		if isTyped {
			return name + "Resource"
		}
	}

	return name
}

func newUnTypedASTReturnEntry(key string, value string, pos int) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key: &ast.BasicLit{
			ValuePos: token.Pos(pos),
			Kind:     token.STRING,
			Value:    fmt.Sprintf("%q", strings.Trim(key, "\"")),
		},
		Value: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: value,
			},
		},
	}
}

func untypedAppendResourceToRegistrationBlock() astutil.ApplyFunc {
	config := LoadConfig()
	return func(c *astutil.Cursor) bool {
		m := c.Node()
		if t, ok := m.(*ast.KeyValueExpr); ok {
			alreadyPresent := false
			snakeName := TerraformResourceName(config.ProviderName, resourceNameToAdd)
			if strings.Trim(t.Key.(*ast.BasicLit).Value, "\"") == snakeName {
				alreadyPresent = true
			}
			p := c.Parent().(*ast.CompositeLit)
			if len(p.Elts)-1 == c.Index() && !alreadyPresent {
				c.InsertAfter(newUnTypedASTReturnEntry(snakeName, fmt.Sprintf("resource%s", strcase.ToCamel(resourceNameToAdd)), int(m.(*ast.KeyValueExpr).Value.(*ast.CallExpr).Rparen)+4))
			}
			return false
		}
		return true
	}
}

func newTypedASTReturnEntry(name string, pos int) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: &ast.Ident{
			NamePos: token.Pos(pos),
			Name:    name,
		},
		Rbrace: token.Pos(len(name) + 6),
	}
}

func firstTypedRegistrationBlockEntry(name string, pos int) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: &ast.Ident{
			NamePos: token.Pos(pos),
			Name:    name,
		},
		Rbrace: token.Pos(len(name) + 6),
	}
}

// TODO - sorted list, by Ident.Name

func typedAppendResourceToRegistrationBlock() astutil.ApplyFunc {
	return func(c *astutil.Cursor) bool {
		m := c.Node()
		typedName := strcase.ToCamel(resourceNameToAdd)
		if t, ok := m.(*ast.CompositeLit); ok {
			if _, ok := t.Type.(*ast.ArrayType); ok && t.Elts == nil {
				// No entries in the list
				t.Elts = []ast.Expr{firstTypedRegistrationBlockEntry(typedName, int(t.Lbrace)+4)}
				c.Replace(t)
				return false
			}

			if p, parentOk := c.Parent().(*ast.CompositeLit); parentOk {
				addEntry := true
				if _, ok := t.Type.(*ast.Ident); ok {
					for _, v := range p.Elts {
						if v.(*ast.CompositeLit).Type.(*ast.Ident).Name == typedName {
							addEntry = false
						}
					}
					if len(p.Elts)-1 == c.Index() && addEntry {
						c.InsertAfter(newTypedASTReturnEntry(typedName, int(m.(*ast.CompositeLit).Rbrace)+4))
					}
					return false
				}
			}
		}
		return true
	}
}
