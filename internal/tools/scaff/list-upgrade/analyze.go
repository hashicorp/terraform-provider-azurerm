package source

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// Kind classifies how a resource is implemented.
type Kind int

const (
	KindUnknown Kind = iota
	// KindTyped is a resource built on the internal/sdk typed wrapper
	// (type XResource struct{} with ResourceType/ModelObject/Create/Read/...).
	KindTyped
	// KindUntyped is a native Plugin SDK resource (func resourceX() *pluginsdk.Resource).
	KindUntyped
)

func (k Kind) String() string {
	switch k {
	case KindTyped:
		return "typed"
	case KindUntyped:
		return "untyped"
	default:
		return "unknown"
	}
}

// Resource is a parsed resource source file together with the metadata and AST
// handles needed to analyse and rewrite it.
type Resource struct {
	Path    string
	Package string
	Kind    Kind

	// Typed resource identity.
	StructName    string // e.g. "WorkspaceResource"
	ModelStruct   string // e.g. "WorkspaceResourceModel"
	TerraformType string // e.g. "azurerm_monitor_workspace"
	HasUpdate     bool

	// SDK + resource ID.
	SDKPackage     string // e.g. "azuremonitorworkspaces"
	SDKImportPath  string // full import path
	IDTypeName     string // e.g. "AccountId"
	IDBase         string // e.g. "Account"
	IDParseFunc    string // e.g. "ParseAccountID"
	IDValidateFunc string // e.g. "ValidateAccountID"

	// Client accessor, parsed from metadata.Client.<Service>.<Field>.
	ServiceName string // e.g. "Monitor"
	ClientField string // e.g. "WorkspacesClient"

	// Feature detection.
	HasIdentity bool
	HasFlatten  bool

	// AST handles.
	fset *token.FileSet
	file *ast.File
	src  []byte

	// Located declarations (typed).
	assertDecl *ast.GenDecl               // var (…) block or single var holding sdk.Resource* assertions
	methods    map[string]*ast.FuncDecl   // resource methods keyed by name
	imports    map[string]*ast.ImportSpec // package name -> import spec
}

// Analyze parses the file at path and derives everything the upgrader needs to
// reason about it. It never mutates the file.
func Analyze(path string) (*Resource, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", path, err)
	}

	r := &Resource{
		Path:    path,
		Package: file.Name.Name,
		Kind:    KindUnknown,
		fset:    fset,
		file:    file,
		src:     src,
		methods: map[string]*ast.FuncDecl{},
		imports: map[string]*ast.ImportSpec{},
	}

	r.collectImports()

	if structName := r.findTypedResourceStruct(); structName != "" {
		r.Kind = KindTyped
		r.StructName = structName
		r.collectMethods(structName)
		r.analyzeTyped()
		return r, nil
	}

	if r.isUntyped() {
		r.Kind = KindUntyped
		return r, nil
	}

	return r, nil
}

// collectImports indexes imports by their effective package name so ID/SDK
// references can be resolved back to an import path.
func (r *Resource) collectImports() {
	for _, imp := range r.file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		name := ""
		if imp.Name != nil {
			name = imp.Name.Name
		} else {
			seg := path
			if i := strings.LastIndex(seg, "/"); i >= 0 {
				seg = seg[i+1:]
			}
			name = seg
		}
		r.imports[name] = imp
	}
}

// findTypedResourceStruct returns the name of the resource struct backing a
// typed resource, identified by an `var _ sdk.Resource... = XResource{}`
// interface assertion. The GenDecl holding the assertion is retained for later
// identity-assertion edits.
func (r *Resource) findTypedResourceStruct() string {
	for _, decl := range r.file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "_" {
				continue
			}
			if !isSDKResourceInterface(vs.Type) {
				continue
			}
			if len(vs.Values) == 1 {
				if name := compositeTypeName(vs.Values[0]); name != "" {
					r.assertDecl = gd
					return name
				}
			}
		}
	}
	return ""
}

// isSDKResourceInterface reports whether expr is a selector of the form
// sdk.Resource or sdk.ResourceWith*.
func isSDKResourceInterface(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok || pkg.Name != "sdk" {
		return false
	}
	return sel.Sel.Name == "Resource" || strings.HasPrefix(sel.Sel.Name, "ResourceWith")
}

// compositeTypeName returns the type name of a composite literal like
// XResource{}, dereferencing a leading &.
func compositeTypeName(expr ast.Expr) string {
	if u, ok := expr.(*ast.UnaryExpr); ok {
		expr = u.X
	}
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return ""
	}
	if id, ok := cl.Type.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// collectMethods indexes all methods whose receiver base type matches struct.
func (r *Resource) collectMethods(structName string) {
	for _, decl := range r.file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil || len(fd.Recv.List) != 1 {
			continue
		}
		if receiverTypeName(fd.Recv.List[0].Type) == structName {
			r.methods[fd.Name.Name] = fd
		}
	}
}

// receiverTypeName returns the base type name of a receiver, dereferencing a
// pointer receiver.
func receiverTypeName(expr ast.Expr) string {
	if star, ok := expr.(*ast.StarExpr); ok {
		expr = star.X
	}
	if id, ok := expr.(*ast.Ident); ok {
		return id.Name
	}
	return ""
}

// analyzeTyped fills in the typed-resource metadata from the located methods and
// interface assertions.
func (r *Resource) analyzeTyped() {
	r.HasUpdate = r.hasAssertion("ResourceWithUpdate") || r.methods["Update"] != nil
	r.HasIdentity = r.hasAssertion("ResourceWithIdentity") || r.methods["Identity"] != nil
	r.HasFlatten = r.methods["flatten"] != nil

	if fd := r.methods["ResourceType"]; fd != nil {
		r.TerraformType = stringReturnValue(fd)
	}
	if fd := r.methods["ModelObject"]; fd != nil {
		r.ModelStruct = modelObjectType(fd)
	}
	if fd := r.methods["IDValidationFunc"]; fd != nil {
		if pkg, fn := selectorReturnValue(fd); fn != "" {
			r.IDValidateFunc = fn
			if r.SDKPackage == "" {
				r.setSDKPackage(pkg)
			}
			r.IDBase = strings.TrimSuffix(strings.TrimPrefix(fn, "Validate"), "ID")
		}
	}

	// The parse call inside Read is the most reliable source of the ID type and
	// SDK package actually used by the resource.
	if fd := r.methods["Read"]; fd != nil {
		if pkg, base := parseIDCall(fd); base != "" {
			r.IDBase = base
			r.setSDKPackage(pkg)
		}
		if svc, field := clientAccessor(fd); field != "" {
			r.ServiceName = svc
			r.ClientField = field
		}
	}
	if r.IDBase != "" {
		r.IDTypeName = r.IDBase + "Id"
		r.IDParseFunc = "Parse" + r.IDBase + "ID"
		if r.IDValidateFunc == "" {
			r.IDValidateFunc = "Validate" + r.IDBase + "ID"
		}
	}
}

// hasAssertion reports whether the resource declares an `sdk.<name>` interface
// assertion (e.g. ResourceWithUpdate, ResourceWithIdentity).
func (r *Resource) hasAssertion(name string) bool {
	for _, decl := range r.file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			continue
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok || len(vs.Names) != 1 || vs.Names[0].Name != "_" {
				continue
			}
			if sel, ok := vs.Type.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "sdk" && sel.Sel.Name == name {
					return true
				}
			}
		}
	}
	return false
}

// setSDKPackage records the SDK package name and resolves its import path.
func (r *Resource) setSDKPackage(pkg string) {
	if pkg == "" {
		return
	}
	r.SDKPackage = pkg
	if imp, ok := r.imports[pkg]; ok {
		r.SDKImportPath = strings.Trim(imp.Path.Value, `"`)
	}
}

// isUntyped reports whether the file declares a native Plugin SDK resource via a
// func returning *pluginsdk.Resource.
func (r *Resource) isUntyped() bool {
	for _, decl := range r.file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv != nil || !strings.HasPrefix(fd.Name.Name, "resource") {
			continue
		}
		if fd.Type.Results == nil || len(fd.Type.Results.List) != 1 {
			continue
		}
		if star, ok := fd.Type.Results.List[0].Type.(*ast.StarExpr); ok {
			if sel, ok := star.X.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == "pluginsdk" && sel.Sel.Name == "Resource" {
					return true
				}
			}
		}
	}
	return false
}

// ----- small AST readers -----

// stringReturnValue returns the string literal from a `return "..."` function.
func stringReturnValue(fd *ast.FuncDecl) string {
	var out string
	ast.Inspect(fd, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		if lit, ok := ret.Results[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			out = strings.Trim(lit.Value, "\"")
			return false
		}
		return true
	})
	return out
}

// modelObjectType returns the type name from a `return &XModel{}` function.
func modelObjectType(fd *ast.FuncDecl) string {
	var out string
	ast.Inspect(fd, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		if name := compositeTypeName(ret.Results[0]); name != "" {
			out = name
			return false
		}
		return true
	})
	return out
}

// selectorReturnValue returns the package and selector from a
// `return pkg.Selector` function (no call), e.g. return pkg.ValidateXID.
func selectorReturnValue(fd *ast.FuncDecl) (pkg, sel string) {
	ast.Inspect(fd, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		if s, ok := ret.Results[0].(*ast.SelectorExpr); ok {
			if id, ok := s.X.(*ast.Ident); ok {
				pkg, sel = id.Name, s.Sel.Name
				return false
			}
		}
		return true
	})
	return pkg, sel
}

// parseIDCall finds a `pkg.Parse<Base>ID(...)` call within fd and returns the
// package and Base.
func parseIDCall(fd *ast.FuncDecl) (pkg, base string) {
	ast.Inspect(fd, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		id, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}
		fn := sel.Sel.Name
		if strings.HasPrefix(fn, "Parse") && strings.HasSuffix(fn, "ID") {
			pkg = id.Name
			base = strings.TrimSuffix(strings.TrimPrefix(fn, "Parse"), "ID")
			return false
		}
		return true
	})
	return pkg, base
}

// clientAccessor finds a `metadata.Client.<Service>.<Field>` selector and
// returns the service and field names.
func clientAccessor(fd *ast.FuncDecl) (service, field string) {
	ast.Inspect(fd, func(n ast.Node) bool {
		// Match the outer selector X.Field where X is metadata.Client.Service.
		outer, ok := n.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		mid, ok := outer.X.(*ast.SelectorExpr) // metadata.Client.Service
		if !ok {
			return true
		}
		inner, ok := mid.X.(*ast.SelectorExpr) // metadata.Client
		if !ok {
			return true
		}
		base, ok := inner.X.(*ast.Ident) // metadata
		if !ok || base.Name != "metadata" || inner.Sel.Name != "Client" {
			return true
		}
		service = mid.Sel.Name
		field = outer.Sel.Name
		return false
	})
	return service, field
}
