package schema

import (
	"go/ast"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/bflad/tfproviderlint/helper/terraformtype/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/resource-lint/helper"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

const commonAnalyzerDoc = `Extracts and caches schema information from the commonschema package in vendor directory.

Key Features:
 1. Extracts schema definitions from commonschema functions (e.g., ResourceGroupName())
 2. Uses ExtractSchemaInfoFromMap to parse resource schema maps
 3. Supports commonschema cross-package calls and same-package function call resolution

Example:

	// Function in commonschema package
	func ResourceGroupName() *pluginsdk.Schema {
	    return &pluginsdk.Schema{
	        Type:     pluginsdk.TypeString,
	        Required: true,
	        ForceNew: true,
	    }
	}

	// Using commonschema in a resource
	func (r MyResource) Arguments() map[string]*pluginsdk.Schema {
	    return map[string]*pluginsdk.Schema{
	        "name": {Type: pluginsdk.TypeString, Required: true},
	        "resource_group_name": commonschema.ResourceGroupName(),  // commonschema cross-package call
	        "metadata": metadataSchema(),                             // same-package call
	    }
	}
`

// CommonSchemaInfo stores information about common schema functions.
type CommonSchemaInfo struct {
	// Map of package.FunctionName -> *schema.SchemaInfo
	Functions map[string]*schema.SchemaInfo
}

var CommonAnalyzer = &analysis.Analyzer{
	Name:       "commonschemainfo",
	Doc:        commonAnalyzerDoc,
	Run:        runCommon,
	ResultType: reflect.TypeOf(&CommonSchemaInfo{}),
}

// Global cache for schema info - loaded only once successfully
var (
	globalSchemaInfo *CommonSchemaInfo
	loadMutex        sync.RWMutex
)

func runCommon(pass *analysis.Pass) (interface{}, error) {
	loadMutex.RLock()
	cached := globalSchemaInfo
	loadMutex.RUnlock()

	if cached != nil {
		return cached, nil
	}

	loadMutex.Lock()
	defer loadMutex.Unlock()

	// Double-check: another goroutine might have loaded while we were waiting
	if globalSchemaInfo != nil {
		return globalSchemaInfo, nil
	}

	info := loadSchemaInfo(pass)

	if len(info.Functions) > 0 {
		globalSchemaInfo = info
		return info, nil
	} else {
		// Failure: don't cache, allow retry on next call
		return info, nil
	}
}

func loadSchemaInfo(pass *analysis.Pass) *CommonSchemaInfo {
	info := &CommonSchemaInfo{
		Functions: make(map[string]*schema.SchemaInfo),
	}

	if len(pass.Files) == 0 {
		return info
	}

	// Get the file path from the first file in the package
	filePath := pass.Fset.Position(pass.Files[0].Pos()).Filename
	// These are go local cache files
	if strings.Contains(filePath, "go-build") || strings.Contains(filePath, "AppData") || strings.Contains(filePath, ".test") {
		return info
	}

	// Traverse up to find the directory containing "internal"
	dir := filepath.Dir(filePath)
	foundInternal := false
	for dir != "" && dir != "." && dir != string(filepath.Separator) {
		base := filepath.Base(dir)
		if base == "internal" {
			// Go up one more level to get the repo root
			dir = filepath.Dir(dir)
			foundInternal = true
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return info
		}
		dir = parent
	}

	if !foundInternal {
		return info
	}

	vendorPath := filepath.Join(dir, "vendor", "github.com", "hashicorp", "go-azure-helpers", "resourcemanager", "commonschema")
	if _, err := os.Stat(vendorPath); os.IsNotExist(err) {
		return info
	}

	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  vendorPath,
	}

	// Load commonschema package from vendor
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.Printf("Warning: failed to load commonschema package: %v", err)
	} else if len(pkgs) > 0 {
		parseHelperPackage(pkgs[0], info)
	}

	return info
}

func parseHelperPackage(helperPkg *packages.Package, info *CommonSchemaInfo) {
	// Parse all functions in the package
	for _, file := range helperPkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Body == nil {
				return true
			}

			// Only process exported functions (that return schemas)
			if !funcDecl.Name.IsExported() {
				return true
			}

			// Extract schema info from function body using package's TypesInfo
			schemaInfo := extractSchemaFromFuncReturn(funcDecl, helperPkg.TypesInfo)
			if schemaInfo != nil {
				key := helperPkg.PkgPath + "." + funcDecl.Name.Name
				info.Functions[key] = schemaInfo
			}

			return true
		})
	}
}

// extractSchemaFromFuncReturn extracts schema info from a function's return statement.
// It handles three patterns:
// 1. Direct return: &schema.Schema{...}
// 2. Composite literal: schema.Schema{...}
// 3. Variable reference: return schemaVar (traces to definition)
func extractSchemaFromFuncReturn(funcDecl *ast.FuncDecl, typesInfo *types.Info) *schema.SchemaInfo {
	var returnedSchema *ast.CompositeLit

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) == 0 {
			return true
		}

		switch expr := ret.Results[0].(type) {
		case *ast.UnaryExpr:
			// Handle &schema.Schema{...}
			if cl, ok := expr.X.(*ast.CompositeLit); ok {
				returnedSchema = cl
				return false
			}
		case *ast.CompositeLit:
			// Handle schema.Schema{...}
			returnedSchema = expr
			return false
		case *ast.Ident:
			// Handle return schema (variable reference)
			if compLit := helper.TraceIdentToCompositeLit(typesInfo, expr, funcDecl); compLit != nil {
				if helper.IsSchemaSchema(typesInfo, compLit) {
					returnedSchema = compLit
				}
			}
			return false
		}

		return true
	})

	if returnedSchema != nil && helper.IsSchemaSchema(typesInfo, returnedSchema) {
		return schema.NewSchemaInfo(returnedSchema, typesInfo)
	}

	return nil
}
