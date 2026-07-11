// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package list_upgrade

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

// resourceFuncBody returns the body of the inner `Func:` closure of a typed
// resource method (Create/Read/Update/Delete), which all follow the shape
// `return sdk.ResourceFunc{ Timeout: ..., Func: func(ctx, metadata) error { ... } }`.
func resourceFuncBody(fd *ast.FuncDecl) *ast.FuncLit {
	var out *ast.FuncLit
	ast.Inspect(fd, func(n ast.Node) bool {
		kv, ok := n.(*ast.KeyValueExpr)
		if !ok {
			return true
		}
		if key, ok := kv.Key.(*ast.Ident); ok && key.Name == "Func" {
			if fl, ok := kv.Value.(*ast.FuncLit); ok {
				out = fl
				return false
			}
		}
		return true
	})
	return out
}

// flattenRegion describes the contiguous run of statements inside Read that build
// and encode the resource state, i.e. the code that should move into flatten.
type flattenRegion struct {
	stateVar  string    // the model variable name, usually "state"
	modelVar  string    // the SDK model variable name, usually "model"
	startPos  token.Pos // start of the `state := X{` statement
	encodePos token.Pos // start of the terminal encode statement
	endPos    token.Pos // end of the terminal encode statement (region end)
}

// findFlattenRegion locates the state-building region within the Read closure.
// It recognises the canonical shape: a `state := <ModelStruct>{...}` assignment
// followed (eventually) by a terminal `return metadata.Encode(&state)` (or the
// `if err := metadata.Encode(&state); err != nil { return err }` / `return nil`
// variant). It returns an error describing why extraction is not possible when
// the shape is not recognised.
func (r *Resource) findFlattenRegion() (*flattenRegion, error) {
	read := r.methods["Read"]
	if read == nil {
		return nil, fmt.Errorf("resource has no Read method")
	}
	body := resourceFuncBody(read)
	if body == nil || body.Body == nil {
		return nil, fmt.Errorf("could not locate the Read closure body")
	}

	stmts := body.Body.List
	startIdx := -1
	reg := &flattenRegion{modelVar: "model"}
	for i, st := range stmts {
		as, ok := st.(*ast.AssignStmt)
		if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
			continue
		}
		lhs, ok := as.Lhs[0].(*ast.Ident)
		if !ok {
			continue
		}
		if name := compositeTypeName(as.Rhs[0]); name != "" && name == r.ModelStruct {
			startIdx = i
			reg.stateVar = lhs.Name
			reg.startPos = as.Pos()
			break
		}
	}
	if startIdx == -1 {
		return nil, fmt.Errorf("could not find a `state := %s{...}` assignment in Read; the resource may already delegate to a flatten function or use a non-standard shape", r.ModelStruct)
	}

	// Find the terminal encode statement at or after the state assignment.
	for i := startIdx; i < len(stmts); i++ {
		if pos, end, ok := encodeStatement(stmts[i], reg.stateVar); ok {
			reg.encodePos = pos
			reg.endPos = end
			// The if-form is followed by a `return nil`; fold it into the region.
			if _, isIf := stmts[i].(*ast.IfStmt); isIf && i+1 < len(stmts) {
				if isReturnNil(stmts[i+1]) {
					reg.endPos = stmts[i+1].End()
				}
			}
			// Detect the SDK model variable name from an `if <mv> := resp.Model`
			// init within the region, if present.
			if mv := modelVarInRegion(stmts[startIdx : i+1]); mv != "" {
				reg.modelVar = mv
			}
			return reg, nil
		}
	}
	return nil, fmt.Errorf("could not find a terminal `metadata.Encode(&%s)` in Read", reg.stateVar)
}

// encodeStatement reports whether st is a terminal encode of &stateVar, either
// `return metadata.Encode(&state)` or `if err := metadata.Encode(&state); err != nil { ... }`.
func encodeStatement(st ast.Stmt, stateVar string) (start, end token.Pos, ok bool) {
	switch s := st.(type) {
	case *ast.ReturnStmt:
		if len(s.Results) == 1 && isEncodeCall(s.Results[0], stateVar) {
			return s.Pos(), s.End(), true
		}
	case *ast.IfStmt:
		if as, ok := s.Init.(*ast.AssignStmt); ok && len(as.Rhs) == 1 && isEncodeCall(as.Rhs[0], stateVar) {
			return s.Pos(), s.End(), true
		}
	}
	return 0, 0, false
}

// isEncodeCall reports whether expr is `metadata.Encode(&stateVar)`.
func isEncodeCall(expr ast.Expr, stateVar string) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok || len(call.Args) != 1 {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Encode" {
		return false
	}
	if x, ok := sel.X.(*ast.Ident); !ok || x.Name != "metadata" {
		return false
	}
	unary, ok := call.Args[0].(*ast.UnaryExpr)
	if !ok || unary.Op != token.AND {
		return false
	}
	id, ok := unary.X.(*ast.Ident)
	return ok && id.Name == stateVar
}

// isReturnNil reports whether st is `return nil`.
func isReturnNil(st ast.Stmt) bool {
	ret, ok := st.(*ast.ReturnStmt)
	if !ok || len(ret.Results) != 1 {
		return false
	}
	id, ok := ret.Results[0].(*ast.Ident)
	return ok && id.Name == "nil"
}

// modelVarInRegion returns the name bound to resp.Model in an
// `if <mv> := resp.Model; <mv> != nil {` init within the given statements.
func modelVarInRegion(stmts []ast.Stmt) string {
	var out string
	for _, st := range stmts {
		ast.Inspect(st, func(n ast.Node) bool {
			ifs, ok := n.(*ast.IfStmt)
			if !ok || ifs.Init == nil {
				return true
			}
			as, ok := ifs.Init.(*ast.AssignStmt)
			if !ok || as.Tok != token.DEFINE || len(as.Lhs) != 1 || len(as.Rhs) != 1 {
				return true
			}
			if !isRespModel(as.Rhs[0]) {
				return true
			}
			if id, ok := as.Lhs[0].(*ast.Ident); ok {
				out = id.Name
				return false
			}
			return true
		})
		if out != "" {
			return out
		}
	}
	return out
}

// isRespModel reports whether expr is a selector `<x>.Model` (typically resp.Model).
func isRespModel(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && sel.Sel.Name == "Model"
}

// extractFlattenEdits schedules the edits that move the state-building region of
// Read into a dedicated flatten method and rewires Read to call it. When
// withIdentity is true the generated flatten also writes the resource identity
// data before encoding.
func (r *Resource) extractFlattenEdits(e *editor, readModel string, withIdentity bool) error {
	if r.SDKPackage == "" || r.IDTypeName == "" {
		return fmt.Errorf("missing SDK package or ID type; cannot build the flatten signature")
	}
	if readModel == "" {
		return fmt.Errorf("the SDK read model type is required to build the flatten signature")
	}

	reg, err := r.findFlattenRegion()
	if err != nil {
		return err
	}

	buildText := r.src[r.offset(reg.startPos):r.offset(reg.encodePos)]
	terminalText := r.src[r.offset(reg.encodePos):r.offset(reg.endPos)]

	build := rewriteRespModel(string(buildText), reg.modelVar)
	terminal := rewriteRespModel(string(terminalText), reg.modelVar)

	// After rewriting, any remaining `resp.` access refers to something other
	// than the model (e.g. resp.HttpResponse) which would not exist inside
	// flatten, so we bail out rather than emit code that will not compile.
	if strings.Contains(build, "resp.") || strings.Contains(terminal, "resp.") {
		return fmt.Errorf("the flatten region references `resp` beyond `resp.Model`; refactor Read manually (see guide-list-resource.md step 1)")
	}

	var flatten strings.Builder
	fmt.Fprintf(&flatten, "\n\nfunc (r %s) flatten(metadata sdk.ResourceMetaData, id *%s.%s, %s *%s.%s) error {\n",
		r.StructName, r.SDKPackage, r.IDTypeName, reg.modelVar, r.SDKPackage, readModel)
	flatten.WriteString(build)
	if !strings.HasSuffix(build, "\n") {
		flatten.WriteString("\n")
	}
	if withIdentity {
		flatten.WriteString("\nif err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {\nreturn err\n}\n\n")
	}
	flatten.WriteString(terminal)
	flatten.WriteString("\n}\n")

	// Insert the flatten method immediately after the Read method.
	e.insert(r.offset(r.methods["Read"].End()), flatten.String())

	// Replace the region in Read with a delegating call.
	e.replace(r.offset(reg.startPos), r.offset(reg.endPos), "return r.flatten(metadata, id, resp.Model)")

	return nil
}

// rewriteRespModel rewrites references to resp.Model so they use the flatten
// method's model parameter. The `if <mv> := resp.Model; <mv> != nil {` init form
// is collapsed to `if <mv> != nil {`, and any bare `resp.Model` becomes `<mv>`.
func rewriteRespModel(text, modelVar string) string {
	// Collapse the init-guard form first.
	text = collapseModelGuard(text, modelVar)
	// Replace any remaining resp.Model references with the parameter.
	text = strings.ReplaceAll(text, "resp.Model", modelVar)
	return text
}

// collapseModelGuard rewrites `if <mv> := resp.Model; <mv> != nil {` to
// `if <mv> != nil {` for the detected model variable name.
func collapseModelGuard(text, modelVar string) string {
	needle := fmt.Sprintf("if %s := resp.Model; ", modelVar)
	return strings.ReplaceAll(text, needle, "if ")
}
