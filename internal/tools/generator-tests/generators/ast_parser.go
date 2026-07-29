package generators

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
)

// InferredIdentity represents the identity properties parsed from the AST
type InferredIdentity struct {
	Properties        []string
	HasSubscriptionID bool
	IsVirtual         bool
}

// InferIdentityProperties parses the given file and attempts to extract the Identity property fields
func InferIdentityProperties(filePath string) (*InferredIdentity, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	importsMap := make(map[string]string)
	for _, imp := range node.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		if imp.Name != nil {
			importsMap[imp.Name.Name] = path
		} else {
			parts := strings.Split(path, "/")
			importsMap[parts[len(parts)-1]] = path
		}
	}

	var identityPkg, identityType string
	var isVirtual bool

	// Walk AST to find Identity() or GenerateIdentitySchema
	ast.Inspect(node, func(n ast.Node) bool {
		// Do not return early, as Identity() and IdentityType() can be separate methods.

		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name != nil && (x.Name.Name == "Identity" || x.Name.Name == "VirtualIdentity") && x.Recv != nil {
				ast.Inspect(x.Body, func(bn ast.Node) bool {
					ret, ok := bn.(*ast.ReturnStmt)
					if ok && len(ret.Results) == 1 {
						if unary, uOk := ret.Results[0].(*ast.UnaryExpr); uOk && unary.Op == token.AND {
							if comp, cOk := unary.X.(*ast.CompositeLit); cOk {
								if sel, sOk := comp.Type.(*ast.SelectorExpr); sOk {
									if ident, iOk := sel.X.(*ast.Ident); iOk {
										identityPkg = ident.Name
										identityType = sel.Sel.Name
										if x.Name.Name == "VirtualIdentity" {
											isVirtual = true
										}
									}
								}
							}
						}
					}
					return true
				})
			} else if x.Name != nil && x.Name.Name == "IdentityType" && x.Recv != nil {
				ast.Inspect(x.Body, func(bn ast.Node) bool {
					ret, ok := bn.(*ast.ReturnStmt)
					if ok && len(ret.Results) == 1 {
						if sel, sOk := ret.Results[0].(*ast.SelectorExpr); sOk {
							if ident, iOk := sel.X.(*ast.Ident); iOk && ident.Name == "pluginsdk" && sel.Sel.Name == "ResourceTypeForIdentityVirtual" {
								isVirtual = true
							}
						}
					}
					return true
				})
			}
		case *ast.CallExpr:
			// Look for Pattern B: pluginsdk.GenerateIdentitySchema(&commonids.PublicIPAddressId{}, ...)
			if sel, ok := x.Fun.(*ast.SelectorExpr); ok {
				if ident, iOk := sel.X.(*ast.Ident); iOk && ident.Name == "pluginsdk" && sel.Sel.Name == "GenerateIdentitySchema" {
					if len(x.Args) >= 1 {
						if unary, uOk := x.Args[0].(*ast.UnaryExpr); uOk && unary.Op == token.AND {
							if comp, cOk := unary.X.(*ast.CompositeLit); cOk {
								if typeSel, sOk := comp.Type.(*ast.SelectorExpr); sOk {
									if typeIdent, tOk := typeSel.X.(*ast.Ident); tOk {
										identityPkg = typeIdent.Name
										identityType = typeSel.Sel.Name
									}
								}
							}
						}
						if len(x.Args) >= 2 {
							if selArg, ok := x.Args[1].(*ast.SelectorExpr); ok {
								if selArg.Sel.Name == "ResourceTypeForIdentityVirtual" {
									isVirtual = true
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	if identityPkg == "" || identityType == "" {
		return nil, fmt.Errorf("could not locate Identity() method or GenerateIdentitySchema call in %s", filePath)
	}

	result, err := resolveIdentityStruct(identityPkg, identityType, importsMap, isVirtual)
	if err == nil && result != nil {
		result.IsVirtual = isVirtual
	}
	return result, err
}

func resolveIdentityStruct(pkgName, typeName string, importsMap map[string]string, isVirtual bool) (*InferredIdentity, error) {
	// Map common identity packages to their relative paths in the vendor folder
	// We assume generator-tests runs from internal/services/<package>
	providerRoot := "../../../"
	var pkgPath string

	if importPath, ok := importsMap[pkgName]; ok {
		if strings.HasPrefix(importPath, "github.com/hashicorp/terraform-provider-azurerm/") {
			pkgPath = filepath.Join(providerRoot, strings.TrimPrefix(importPath, "github.com/hashicorp/terraform-provider-azurerm/"))
		} else {
			pkgPath = filepath.Join(providerRoot, "vendor", importPath)
		}
	} else {
		// Attempt a best-effort fallback to hardcoded paths if imports failed
		switch pkgName {
		case "commonids":
			pkgPath = filepath.Join(providerRoot, "vendor/github.com/hashicorp/go-azure-helpers/resourcemanager/commonids")
		case "resourceids":
			pkgPath = filepath.Join(providerRoot, "vendor/github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids")
		default:
			return nil, fmt.Errorf("unknown identity package: %s (cannot resolve path)", pkgName)
		}
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory %s: %w", pkgPath, err)
	}

	var fields []string
	found := false

	// Search all files in the package for the type definition
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(n ast.Node) bool {
				if found {
					return false
				}
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != typeName {
					return true
				}

				// Check if it's a struct
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					found = true
					for _, field := range structType.Fields.List {
						for _, name := range field.Names {
							fields = append(fields, name.Name)
						}
					}
					return false
				}

				if ident, ok := typeSpec.Type.(*ast.Ident); ok {
					// We need to resolve the alias by recursing
					inferred, err := resolveIdentityStruct(pkgName, ident.Name, importsMap, isVirtual)
					if err == nil {
						fields = inferred.Properties
						if inferred.HasSubscriptionID {
							fields = append(fields, "SubscriptionId") // Will be checked again below
						}
						found = true
					}
					return false
				}

				return true
			})
		}
	}

	if !found {
		return nil, fmt.Errorf("could not find struct definition for %s.%s", pkgName, typeName)
	}

	result := &InferredIdentity{}
	for i, f := range fields {
		snake := toSnakeCase(f)
		if snake == "subscription_id" {
			result.HasSubscriptionID = true
		} else {
			if i == len(fields)-1 && !isVirtual {
				result.Properties = append(result.Properties, "name")
			} else {
				result.Properties = append(result.Properties, snake)
			}
		}
	}

	return result, nil
}

func toSnakeCase(s string) string {
	var res strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				res.WriteByte('_')
			}
			res.WriteRune(unicode.ToLower(r))
		} else {
			res.WriteRune(r)
		}
	}
	// Some structs might end up with `i_d`, fix to `id`
	return strings.ReplaceAll(res.String(), "_i_d", "_id")
}
