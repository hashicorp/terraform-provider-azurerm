// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
)

// RegisterTypedResource ensures the new typed struct is wired into the service
// package's registration.go:
//   - Added to the Resources() []sdk.Resource slice (creating the method if
//     absent).
//   - Removed from the SupportedResources() map (leaving an empty map when
//     it becomes empty — do not delete the method).
//
// It returns the new source and whether anything changed, and does not write
// to disk.
func RegisterTypedResource(path, structName, terraformType string) (newSrc []byte, changed bool, err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("reading %q: %w", path, err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, false, fmt.Errorf("parsing %q: %w", path, err)
	}

	e := &regEditor{src: src, fset: fset}

	if c, err := e.addToResources(file, structName); err != nil {
		return nil, false, fmt.Errorf("adding to Resources(): %w", err)
	} else {
		changed = changed || c
	}

	if c, err := e.removeFromSupportedResources(file, terraformType); err != nil {
		return nil, false, fmt.Errorf("removing from SupportedResources(): %w", err)
	} else {
		changed = changed || c
	}

	if !changed {
		return src, false, nil
	}

	out, err := e.bytes()
	if err != nil {
		return nil, false, err
	}
	formatted, err := format.Source(out)
	if err != nil {
		return out, true, nil // return unformatted but don't fail
	}
	return formatted, true, nil
}

// regEditor accumulates byte-range edits against the source buffer.
type regEditor struct {
	src   []byte
	fset  *token.FileSet
	edits []regEdit
}

type regEdit struct {
	start, end int
	text       string
}

func (e *regEditor) insert(at int, text string) {
	e.edits = append(e.edits, regEdit{start: at, end: at, text: text})
}

func (e *regEditor) replace(start, end int, text string) {
	e.edits = append(e.edits, regEdit{start: start, end: end, text: text})
}

func (e *regEditor) offset(pos token.Pos) int {
	return e.fset.Position(pos).Offset
}

func (e *regEditor) bytes() ([]byte, error) {
	sorted := make([]regEdit, len(e.edits))
	copy(sorted, e.edits)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].start != sorted[j].start {
			return sorted[i].start < sorted[j].start
		}
		return sorted[i].end < sorted[j].end
	})
	var out []byte
	prev, lastEnd := 0, 0
	for _, ed := range sorted {
		if ed.start < lastEnd {
			return nil, fmt.Errorf("overlapping edits at offset %d", ed.start)
		}
		out = append(out, e.src[prev:ed.start]...)
		out = append(out, ed.text...)
		prev = ed.end
		lastEnd = ed.end
	}
	out = append(out, e.src[prev:]...)
	return out, nil
}

// addToResources ensures structName{} is present in the Resources() method,
// creating the method when it does not yet exist.
func (e *regEditor) addToResources(file *ast.File, structName string) (bool, error) {
	const (
		methodName = "Resources"
		recvType   = "Registration"
		sliceType  = "[]sdk.Resource"
	)

	method := findRegistrationMethod(file, methodName)
	entry := structName + "{},\n"

	if method == nil {
		text := fmt.Sprintf("\n\nfunc (r %s) %s() %s {\n\treturn %s{\n%s}\n}\n",
			recvType, methodName, sliceType, sliceType, "\t\t"+entry)
		e.insert(len(e.src), text)
		return true, nil
	}

	lit := returnedSliceLiteral(method)
	if lit == nil {
		return false, fmt.Errorf("could not find the returned slice literal in %s()", methodName)
	}
	for _, el := range lit.Elts {
		if compositeTypeName(el) == structName {
			return false, nil // already registered
		}
	}

	insertAt := sliceInsertOffset(e.fset, e.src, lit, structName)
	e.insert(insertAt, "\t\t"+entry)
	return true, nil
}

// removeFromSupportedResources removes the terraformType key from
// SupportedResources() and ensures the interface assertion is preserved.
func (e *regEditor) removeFromSupportedResources(file *ast.File, terraformType string) (bool, error) {
	method := findRegistrationMethod(file, "SupportedResources")
	if method == nil {
		return false, nil
	}
	lit := returnedMapLiteral(method)
	if lit == nil {
		return false, nil
	}
	for _, el := range lit.Elts {
		kv, ok := el.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		keyLit, ok := kv.Key.(*ast.BasicLit)
		if !ok || keyLit.Kind != token.STRING {
			continue
		}
		if strings.Trim(keyLit.Value, `"`) != terraformType {
			continue
		}
		// Delete from the start of the line to the end of the entry (inclusive).
		start := lineStart(e.fset, e.src, kv.Key.Pos())
		end := e.offset(kv.End())
		// Consume the trailing comma and newline if present.
		for end < len(e.src) && (e.src[end] == ',' || e.src[end] == '\n') {
			end++
		}
		e.replace(start, end, "")
		return true, nil
	}
	return false, nil
}

// findRegistrationMethod returns the named method on the Registration receiver.
func findRegistrationMethod(file *ast.File, name string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Recv == nil || len(fd.Recv.List) != 1 {
			continue
		}
		if fd.Name.Name != name {
			continue
		}
		if recvName(fd.Recv.List[0].Type) == "Registration" {
			return fd
		}
	}
	return nil
}

// recvName returns the identifier name for a receiver type expression.
func recvName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		if id, ok := e.X.(*ast.Ident); ok {
			return id.Name
		}
	}
	return ""
}

// returnedSliceLiteral finds the []sdk.Resource{...} literal returned by method.
func returnedSliceLiteral(method *ast.FuncDecl) *ast.CompositeLit {
	var out *ast.CompositeLit
	ast.Inspect(method, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		if cl, ok := ret.Results[0].(*ast.CompositeLit); ok {
			out = cl
			return false
		}
		return true
	})
	return out
}

// returnedMapLiteral finds the map literal returned by method.
func returnedMapLiteral(method *ast.FuncDecl) *ast.CompositeLit {
	var out *ast.CompositeLit
	ast.Inspect(method, func(n ast.Node) bool {
		ret, ok := n.(*ast.ReturnStmt)
		if !ok || len(ret.Results) != 1 {
			return true
		}
		if cl, ok := ret.Results[0].(*ast.CompositeLit); ok {
			out = cl
			return false
		}
		return true
	})
	return out
}

// sliceInsertOffset returns the byte offset at which to insert a new slice
// element (in alphabetical order, or before the closing brace).
func sliceInsertOffset(fset *token.FileSet, src []byte, lit *ast.CompositeLit, newName string) int {
	for _, el := range lit.Elts {
		if compositeTypeName(el) > newName {
			return lineStart(fset, src, el.Pos())
		}
	}
	return lineStart(fset, src, lit.Rbrace)
}

// lineStart returns the byte offset of the first character of the line
// containing pos.
func lineStart(fset *token.FileSet, src []byte, pos token.Pos) int {
	off := fset.Position(pos).Offset
	for off > 0 && src[off-1] != '\n' {
		off--
	}
	return off
}
