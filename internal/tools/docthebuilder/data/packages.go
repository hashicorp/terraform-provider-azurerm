package data

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	log "github.com/sirupsen/logrus"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type packageData struct {
	cfg      *packages.Config
	fset     *token.FileSet
	prog     *ssa.Program
	packages map[string][]pkg
}

type pkg struct {
	id  string
	ssa *ssa.Package
	pkg *packages.Package
}

type resourceID struct {
	Type ResourceType
	Name string
}

type sdkMethod struct {
	MethodName string
	APIPath    string
	Pkg        *packages.Package
}

var (
	sdkRegex             = regexp.MustCompile(`github.com/hashicorp/go-azure-sdk/resource-manager`)
	sdkMethodSuffixRegex = regexp.MustCompile(`ThenPoll|Complete|CompleteMatchingPredicate`)
	servicePackageRegex  = regexp.MustCompile(`github.com/hashicorp/terraform-provider-azurerm/internal/services/(\w*)$`)
)

// loadPackages - Loads all service packages
// This is very slow, which is why we do it once, and gather the data for ALL resources regardless of filters applied
func loadPackages(dir string) *packageData {
	cfg := &packages.Config{
		Dir:  dir + "/internal/services",
		Mode: packages.LoadAllSyntax,
	}

	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		log.WithError(err).Fatal("encountered an error loading provider packages")
	}

	fset := pkgs[0].Fset
	prog, ssaPkgs := ssautil.AllPackages(pkgs, ssa.InstantiateGenerics)
	prog.Build()

	if ssaPkgCount, pkgCount := len(ssaPkgs), len(pkgs); ssaPkgCount != pkgCount {
		log.Fatalf("number of SSA packages (%d) did not equal number of Go packages (%d)", ssaPkgCount, pkgCount)
	}

	pkgsMap := make(map[string][]pkg)
	for idx := range ssaPkgs {
		id := pkgs[idx].ID

		if servicePackageRegex.MatchString(id) {
			m := servicePackageRegex.FindStringSubmatch(id)

			// map by service base directory name
			pkgsMap[m[1]] = append(pkgsMap[m[1]], pkg{
				ssa: ssaPkgs[idx],
				pkg: pkgs[idx],
			})
		}
	}

	return &packageData{
		cfg:      cfg,
		fset:     fset,
		prog:     prog,
		packages: pkgsMap,
	}
}

/*
TODO:
- autorest?
- azure-sdk-for-go?
- kermit?
- data plane APIs?
*/

func findRegistrationFuncs(pkg *packages.Package, fn *types.Func, registrationFuncs map[resourceID]*types.Func, t ResourceType) {
	fnDecl := funcToFuncDeclWithPkgs([]*packages.Package{pkg}, fn)
	if fnDecl == nil {
		return
	}

	ast.Inspect(fnDecl.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.BlockStmt:
			return true
		case *ast.ReturnStmt:
			return true
		case *ast.CompositeLit:
			return true
		case *ast.IfStmt:
			return true
		case *ast.ExprStmt:
			return true
		case *ast.AssignStmt: // if feature flagged
			/* assuming 5.0 deprecation guidelines are followed:
			- We should see an AssignStmt for the non-deprecated Resources, which should be a new var, we can use`token.DEFINE` to skip this one
			- The next AssignStmt should live inside a conditional, and be of type `token.ASSIGN`, which is what we're interested in
			*/
			if n.Tok == token.DEFINE {
				return false
			}

			if len(n.Lhs) != 1 || len(n.Rhs) != 1 {
				return true
			}

			lhs, ok := n.Lhs[0].(*ast.IndexExpr)
			if !ok {
				return true
			}

			rhs, ok := n.Rhs[0].(*ast.CallExpr)
			if !ok {
				return true
			}

			lit, ok := lhs.Index.(*ast.BasicLit)
			if !ok {
				return true
			}

			name, _ := strconv.Unquote(lit.Value)

			fn := pkg.TypesInfo.ObjectOf(rhs.Fun.(*ast.Ident))
			if fn == nil {
				return false
				// todo err
			}

			registrationFuncs[resourceID{
				Type: t,
				Name: name,
			}] = fn.(*types.Func)
		case *ast.KeyValueExpr:
			k, ok := n.Key.(*ast.BasicLit)
			if !ok || k.Kind != token.STRING {
				return false
			}

			v, ok := n.Value.(*ast.CallExpr)
			if !ok {
				return false
			}

			fn := pkg.TypesInfo.ObjectOf(v.Fun.(*ast.Ident))
			if fn == nil {
				return false
				// todo err
			}

			name, _ := strconv.Unquote(k.Value)

			registrationFuncs[resourceID{
				Type: t,
				Name: name,
			}] = fn.(*types.Func)
		}
		return false
	})
}

func findUntypedSSAFunc(pkg pkg, e ast.Expr) *ssa.Function {
	switch e := e.(type) {
	case *ast.Ident:
		return pkg.ssa.Func(e.Name)
	case *ast.CallExpr:
		fn, ok := e.Fun.(*ast.Ident)
		if !ok {
			return nil
		}
		ssaFn := findResourceFunc(pkg.ssa.Prog, pkg.pkg, pkg.ssa, pkg.ssa.Func(fn.Name))
		return ssaFn
	}

	return nil
}

func findAPIsForUntypedResources(d packageData, s *Service) map[string][]API {
	result := make(map[string][]API)

	servicePackages, ok := d.packages[s.Name]
	// todo: refactor loadPackages to only return the service package and ignore others such as `client`
	servicePackage := servicePackages[0]
	if !ok {
		return nil //err?
	}

	registration := servicePackage.pkg.Types.Scope().Lookup("Registration")
	if registration == nil {
		return nil // err?
	}

	var dsRegistration *types.Func
	var rsRegistration *types.Func
	for m := range registration.Type().(*types.Named).Methods() {
		switch m.Name() {
		case "SupportedDataSources":
			dsRegistration = m
		case "SupportedResources":
			rsRegistration = m
		}
	}

	registrationFuncs := make(map[resourceID]*types.Func)
	// Datasources
	findRegistrationFuncs(servicePackage.pkg, dsRegistration, registrationFuncs, ResourceTypeData)

	// Resources
	findRegistrationFuncs(servicePackage.pkg, rsRegistration, registrationFuncs, ResourceTypeResource)

	for _, fn := range registrationFuncs {
		filenames := make(map[string]struct{})
		resourceFileName := d.fset.Position(d.prog.FuncValue(fn).Pos()).Filename
		filenames[resourceFileName] = struct{}{} // for most resources this is all we need

		fnDecl := funcToFuncDeclWithPkgs([]*packages.Package{servicePackage.pkg}, fn)
		if fnDecl == nil {
			continue
		}

		ast.Inspect(fnDecl, func(n ast.Node) bool {
			switch n := n.(type) {
			case *ast.FuncDecl:
				return true
			case *ast.BlockStmt:
				return true
			case *ast.AssignStmt:
				return true
			case *ast.ReturnStmt:
				return true
			case *ast.UnaryExpr:
				return true
			case *ast.CompositeLit:
				return true
			case *ast.KeyValueExpr:
				var resourceFn *ssa.Function
				// here we can find the fields, including `Read:`, `Create:` etc
				k, ok := n.Key.(*ast.Ident)
				if !ok {
					return false
				}

				switch k.Name {
				case "Create":
					resourceFn = findUntypedSSAFunc(servicePackage, n.Value)
				case "Read":
					resourceFn = findUntypedSSAFunc(servicePackage, n.Value)
				case "Update":
					resourceFn = findUntypedSSAFunc(servicePackage, n.Value)
				case "Delete":
					resourceFn = findUntypedSSAFunc(servicePackage, n.Value)
				}

				if resourceFn != nil {
					filenames[d.fset.Position(resourceFn.Pos()).Filename] = struct{}{}
				}
			}

			return false
		})

		sdkMethods := usedMethods(d.fset, servicePackage.pkg, util.MapKeys2Slice(filenames))
		apis := methodsToAPIs(sdkMethods)

		result[resourceFileName] = apis
	}

	return result
}

// Finds all APIs used by resources and returns a map of files, stores as map[<resource filename>][]API
func findAPIsForTypedResources(d packageData, s *Service) map[string][]API {
	possibleFunctionNames := []string{"Create", "Read", "Update", "Delete"}
	result := make(map[string][]API)

	servicePackages, ok := d.packages[s.Name]
	// todo: refactor loadPackages to only return the service package and ignore others such as `client`
	servicePackage := servicePackages[0]
	if !ok {
		return nil //err?
	}

	ssaPkg := servicePackage.ssa
	scope := ssaPkg.Pkg.Scope()

	for _, scopeName := range scope.Names() {
		obj := scope.Lookup(scopeName)
		typeName, ok := obj.(*types.TypeName)
		if !ok {
			continue
		}

		named, ok := typeName.Type().(*types.Named)
		if !ok {
			continue
		}

		sel := d.prog.MethodSets.MethodSet(named).Lookup(servicePackage.pkg.Types, "ResourceType") // if contains `ResourceType` method then we know we want to parse this file for used sdk methods
		if sel == nil {
			continue
		}

		filenames := make(map[string]struct{})
		fPos := d.fset.Position(d.prog.MethodValue(sel).Pos())
		resourceFileName := fPos.Filename
		filenames[resourceFileName] = struct{}{}

		for _, n := range possibleFunctionNames {
			sel := d.prog.MethodSets.MethodSet(named).Lookup(servicePackage.pkg.Types, n)
			if sel == nil {
				continue
			}

			ssaFn := d.prog.MethodValue(sel)

			rfn := findResourceFunc(d.prog, servicePackage.pkg, servicePackage.ssa, ssaFn)
			if rfn == nil {
				continue //err or debug log?
			}

			// for the SSA func, get file name and add to map, most of the time this is exactly the same as `resourceFileName`
			ssaFnPos := d.fset.Position(rfn.Pos())
			filenames[ssaFnPos.Filename] = struct{}{} // overwrite ok
		}

		sdkMethods := usedMethods(d.fset, servicePackage.pkg, util.MapKeys2Slice(filenames))
		apis := methodsToAPIs(sdkMethods)

		result[resourceFileName] = apis
	}
	return result
}

func findResourceFunc(prog *ssa.Program, pkg *packages.Package, ssaPkg *ssa.Package, fn *ssa.Function) *ssa.Function {
	switch len(fn.AnonFuncs) {
	case 0:
		return findResourceFuncDigDeeper(prog, pkg, ssaPkg, fn.Object().(*types.Func))
	case 1:
		return fn.AnonFuncs[0]
	default:
		log.WithFields(log.Fields{
			"count":    len(fn.AnonFuncs),
			"function": fn.Name(),
			"package":  fn.Pkg.String(),
		}).Debug("unexpected number of Anonymous Functions")
	}

	return nil
}

func findResourceFuncDigDeeper(prog *ssa.Program, pkg *packages.Package, ssaPkg *ssa.Package, fn *types.Func) *ssa.Function {
	debugLog := func(t any, msg string) {
		log.WithFields(log.Fields{
			"type":     reflect.TypeOf(t),
			"function": fn.Name(),
			"package":  fn.Pkg().Path(),
		}).Debug(msg)
	}

	fnDecl := funcToFuncDeclWithPkgs([]*packages.Package{pkg}, fn)

	r := fnDecl.Body.List[0].(*ast.ReturnStmt).Results[0]
	c, ok := r.(*ast.CallExpr)
	if !ok {
		debugLog(c, "expected return value to be of type `*ast.CallExpr`")
		return nil
	}

	switch f := c.Fun.(type) {
	case *ast.SelectorExpr:
		ssaFn := prog.LookupMethod(pkg.TypesInfo.TypeOf(f.X).(*types.Named), pkg.Types, f.Sel.Name)
		return findResourceFunc(prog, pkg, ssaPkg, ssaFn)
	case *ast.Ident:
		fnObj := pkg.TypesInfo.ObjectOf(f).(*types.Func)
		return findResourceFunc(prog, pkg, ssaPkg, ssaPkg.Func(fnObj.Name()))
	default:
		debugLog(f, "unexpected type")
		return nil
	}
}

func findMethodByName(t *types.Named, name string) *types.Func {
	for m := range t.Methods() {
		if m.Name() == name {
			return m
		}
	}

	return nil
}

func usedMethods(fset *token.FileSet, pkg *packages.Package, fileNames []string) []sdkMethod {
	sdkMethods := make(map[sdkMethod]struct{})
	sdkImportPkgsMap := make(map[*packages.Package]struct{})
	var commonIdsPkg *packages.Package

	for _, impP := range pkg.Imports {
		if commonIdsPkg == nil && strings.Contains(impP.ID, "commonids") { // github.com/hashicorp/go-azure-helpers/resourcemanager/commonids
			commonIdsPkg = impP
		}
		if sdkRegex.MatchString(impP.ID) { // only go-azure-sdk for now
			sdkImportPkgsMap[impP] = struct{}{} // don't want duplicates, overwrite ok
		}
	}

	for _, node := range pkg.Syntax {
		// TODO: pass *ast.File node instead of checking each one in `pkg.Syntax`?
		for _, name := range fileNames {
			if !matchesFile(fset, node, name, false) {
				continue
			}

			ast.Inspect(node, func(n ast.Node) bool {
				apiPath := ""
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				sel, ok := call.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}

				recvIdent, ok := sel.X.(*ast.Ident)
				if !ok {
					return true
				}

				recvObj := pkg.TypesInfo.Uses[recvIdent]
				if !isUnderlyingStruct(recvObj.Type()) {
					return true
				}
				recvType := stripPointer(recvObj.Type()).(*types.Named)

				recvTypePkg := recvType.Obj().Pkg()
				if recvTypePkg == nil {
					return true
				}

				if !sdkRegex.MatchString(recvTypePkg.Path()) {
					return true
				}

				recvTypePos := recvType.Obj().Pos()

				sdkPkg, file := findFileByPosition(util.MapKeys2Slice(sdkImportPkgsMap), recvTypePos)
				if file == nil {
					return true
				}

				sdkFn := findMethodByName(recvType, sdkMethodSuffixRegex.ReplaceAllString(sel.Sel.Name, ""))
				if sdkFn == nil {
					return true
				}

				sdkFnDecl := funcToFuncDeclWithPkgs([]*packages.Package{sdkPkg}, sdkFn)

				ast.Inspect(sdkFnDecl.Body.List[0], func(n ast.Node) bool {
					switch n := n.(type) {
					case *ast.AssignStmt:
						return true
					case *ast.CompositeLit:
						return true
					case *ast.KeyValueExpr:
						if k, ok := n.Key.(*ast.Ident); ok && k.Name == "Path" {
							switch v := n.Value.(type) {
							case *ast.BasicLit:
								apiPath, _ = strconv.Unquote(v.Value)
							case *ast.CallExpr:
								if fn, ok := v.Fun.(*ast.SelectorExpr); ok {
									switch fn.X.(*ast.Ident).Name {
									case "id":
										apiPath = apiPathFromID(sdkPkg, commonIdsPkg, fn)
									case "fmt":
										sel := v.Args[1].(*ast.CallExpr).Fun.(*ast.SelectorExpr)
										apiPath = apiPathFromID(sdkPkg, commonIdsPkg, sel)

										// Some API providers are actually contained in the fmt string, e.g.
										// fmt.Sprintf("%s/providers/Microsoft.Security/defenderForStorageSettings/current", id.ID())
										// thus we add the format string to the API path
										format, _ := strconv.Unquote(v.Args[0].(*ast.BasicLit).Value)
										apiPath = format + apiPath
									}
								}
							}
						}

					}
					return false
				})

				// TODO: Some method calls like `ID` are being caught by this, filter?
				m := sdkMethod{
					APIPath:    apiPath,
					Pkg:        sdkPkg,
					MethodName: sel.Sel.Name,
				}
				sdkMethods[m] = struct{}{}

				return true
			})
		}
	}

	return util.MapKeys2Slice(sdkMethods)
}

func apiPathFromID(pkg *packages.Package, commonIdsPkg *packages.Package, id *ast.SelectorExpr) string {
	obj, ok := pkg.TypesInfo.Uses[id.X.(*ast.Ident)]
	if !ok {
		return ""
	}

	pkgs := []*packages.Package{pkg}
	if commonIdsPkg != nil {
		pkgs = append(pkgs, commonIdsPkg)
	}

	idFn := findMethodByName(obj.Type().(*types.Named), id.Sel.Name)
	idFnDecl := funcToFuncDeclWithPkgs(pkgs, idFn)

	path, _ := strconv.Unquote(idFnDecl.Body.List[0].(*ast.AssignStmt).Rhs[0].(*ast.BasicLit).Value)

	return path
}

func stripPointer(t types.Type) types.Type {
	if ptr, ok := t.(*types.Pointer); ok {
		return stripPointer(ptr.Elem())
	}

	return t
}

func isUnderlyingStruct(t types.Type) bool {
	t = stripPointer(t)

	named, ok := t.(*types.Named)
	if !ok {
		return false
	}

	if _, ok := named.Underlying().(*types.Struct); !ok {
		return false
	}

	return true
}

func funcToFuncDeclWithPkgs(pkgs []*packages.Package, fn *types.Func) *ast.FuncDecl {
	_, file := findFileByPosition(pkgs, fn.Pos())

	if file == nil {
		log.WithFields(log.Fields{
			"function": fn.Name(),
			"scope":    fn.Scope().String(),
		}).Debug("unable to find AST File object for function in provided packages")
	}

	return funcToFuncDeclWithFile(file, fn)
}

func funcToFuncDeclWithFile(file *ast.File, fn *types.Func) *ast.FuncDecl {
	if fn == nil {
		log.Debug("unable to find *ast.FuncDecl, *types.Func was nil")
		return nil
	}

	pos := fn.Pos()
	paths, _ := astutil.PathEnclosingInterval(file, pos, pos)

	return paths[1].(*ast.FuncDecl)
}

func matchesFile(fset *token.FileSet, node *ast.File, targetName string, base bool) bool {
	filePos := fset.File(node.Pos())

	if filePos == nil {
		return false
	}

	name := filePos.Name()
	if base {
		name = filepath.Base(name)
	}

	if name == targetName {
		return true
	}

	return false
}

func findFileByPosition(pkgs []*packages.Package, pos token.Pos) (*packages.Package, *ast.File) {
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			if pos >= file.FileStart && pos <= file.FileEnd {
				return pkg, file
			}
		}
	}

	log.Debug("unable to find *ast.File at provided position in provided packages")
	return nil, nil
}

func findFileByName(pkg *packages.Package, name string) *ast.File {
	for _, file := range pkg.Syntax {
		if matchesFile(pkg.Fset, file, name, true) {
			return file
		}
	}

	return nil
}
