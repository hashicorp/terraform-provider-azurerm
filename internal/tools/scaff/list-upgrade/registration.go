package source

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

const (
	registrationStruct   = "Registration"
	listResourcesMethod  = "ListResources"
	listResourcesSlice   = "[]sdk.FrameworkListWrappedResource"
	listResourcesReceipt = "r"
)

// RegisterListResource ensures listStruct (e.g. "WorkspaceListResource") is
// registered in the service package's registration.go ListResources() method,
// creating the method when it does not yet exist. It returns the new source and
// whether anything changed; it does not write to disk.
func RegisterListResource(path, listStruct string) (newSrc []byte, changed bool, err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("reading %q: %w", path, err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, false, fmt.Errorf("parsing %q: %w", path, err)
	}

	e := newEditor(src)
	method := findRegistrationMethod(file, listResourcesMethod)

	// No ListResources() yet: append the whole method.
	if method == nil {
		text := fmt.Sprintf("\n\nfunc (%s %s) %s() %s {\n\treturn %s{\n%s{},\n}\n}\n",
			listResourcesReceipt, registrationStruct, listResourcesMethod, listResourcesSlice, listResourcesSlice, listStruct)
		e.insert(len(src), text)
		out, err := e.bytes()
		return out, true, err
	}

	lit := returnedSliceLiteral(method)
	if lit == nil {
		return nil, false, fmt.Errorf("could not find the returned slice literal in %s()", listResourcesMethod)
	}

	// Already registered: nothing to do.
	for _, el := range lit.Elts {
		if compositeTypeName(el) == listStruct {
			return src, false, nil
		}
	}

	entry := listStruct + "{},\n"

	switch {
	case len(lit.Elts) == 0:
		// Empty slice literal `{}` -> populate it.
		lb := fset.Position(lit.Lbrace).Offset
		rb := fset.Position(lit.Rbrace).Offset
		e.replace(lb, rb+1, fmt.Sprintf("{\n%s}", entry))
	default:
		// Insert in alphabetical order to match the file convention. gofmt
		// (run by the caller) normalises the indentation afterwards.
		insertAt := -1
		for _, el := range lit.Elts {
			if compositeTypeName(el) > listStruct {
				insertAt = lineStartOffset(fset, src, el.Pos())
				break
			}
		}
		if insertAt == -1 {
			insertAt = lineStartOffset(fset, src, lit.Rbrace)
		}
		e.insert(insertAt, entry)
	}

	out, err := e.bytes()
	return out, true, err
}

// findRegistrationMethod returns the named method on the Registration receiver.
func findRegistrationMethod(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil || len(fd.Recv.List) != 1 {
			continue
		}
		if fd.Name.Name == name && receiverTypeName(fd.Recv.List[0].Type) == registrationStruct {
			return fd
		}
	}
	return nil
}

// returnedSliceLiteral returns the composite literal from the method's
// `return []sdk.FrameworkListWrappedResource{...}` statement.
func returnedSliceLiteral(fd *ast.FuncDecl) *ast.CompositeLit {
	var out *ast.CompositeLit
	ast.Inspect(fd, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		if cl, ok := ret.Results[0].(*ast.CompositeLit); ok {
			if _, ok := cl.Type.(*ast.ArrayType); ok {
				out = cl
				return false
			}
		}
		return true
	})
	return out
}

// lineStartOffset returns the byte offset of the first character of the line
// containing pos.
func lineStartOffset(fset *token.FileSet, src []byte, pos token.Pos) int {
	off := fset.Position(pos).Offset
	for off > 0 && src[off-1] != '\n' {
		off--
	}
	return off
}
