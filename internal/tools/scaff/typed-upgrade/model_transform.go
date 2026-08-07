// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"regexp"
	"strings"
)

// applyModelTransforms applies the model-field substitution pass to CRUD
// function bodies. It uses the extracted schema fields to:
//   - Replace metadata.ResourceData.Get("field").(TYPE) → model.Field (Create/Update)
//   - Replace metadata.ResourceData.Set("field", value) → state.Field = value (Read)
//   - Inject metadata.Decode at the start of Create/Update bodies
//   - Inject state declaration and metadata.Encode at the start/end of Read body
//
// When a Pandora IR is present (info.PandoraIR != nil), it also:
//   - Updates block model field names to match the IR
//   - Replaces old expand/flatten calls with the IR-generated typed ones
//
// When renames are present (info.Renames is non-empty), it updates CRUD code
// references from old field names to new field names.
func (info *Info) applyModelTransforms() {
	// Apply renames first, before other transformations.
	if len(info.Renames) > 0 {
		if info.CreateBody != "" {
			info.CreateBody = applyFieldRenames(info.CreateBody, info.Renames)
		}
		if info.UpdateBody != "" {
			info.UpdateBody = applyFieldRenames(info.UpdateBody, info.Renames)
		}
		if info.ReadBody != "" {
			info.ReadBody = applyFieldRenames(info.ReadBody, info.Renames)
		}
		if info.DeleteBody != "" {
			info.DeleteBody = applyFieldRenames(info.DeleteBody, info.Renames)
		}
	}

	// When we have an IR, update the block field names so the model struct and
	// the expand/flatten functions agree.
	if info.PandoraIR != nil {
		blocksByTF := buildIRBlockByTFName(info.PandoraIR)
		updateModelFieldsFromIR(info.Fields, blocksByTF)
	}

	fm := buildFieldMap(info.Fields)
	if len(fm) == 0 {
		return
	}

	if info.CreateBody != "" {
		info.CreateBody = applyDecodeTransform(info.CreateBody, info.ModelName, fm)
		if info.PandoraIR != nil {
			info.CreateBody = applyBlockTransforms(info.CreateBody, info)
		}
	}
	if info.UpdateBody != "" {
		info.UpdateBody = applyDecodeTransform(info.UpdateBody, info.ModelName, fm)
		if info.PandoraIR != nil {
			info.UpdateBody = applyBlockTransforms(info.UpdateBody, info)
		}
	}
	if info.ReadBody != "" {
		info.ReadBody = applyEncodeTransform(info.ReadBody, info.ModelName, fm)
		if info.PandoraIR != nil {
			info.ReadBody = applyBlockTransforms(info.ReadBody, info)
		}
	}
}

// buildFieldMap indexes schema fields by their TF name. Only transformable
// scalar fields are included.
func buildFieldMap(fields []*SchemaField) map[string]*SchemaField {
	m := make(map[string]*SchemaField, len(fields))
	for _, f := range fields {
		if isTransformableField(f) {
			m[f.TFName] = f
		}
	}
	return m
}

// isTransformableField reports whether we can safely replace Get/Set calls
// for this field with model field accesses. Nested blocks with IR backing are
// also included so the model carries the typed struct slice.
func isTransformableField(f *SchemaField) bool {
	switch f.Kind {
	case FieldString, FieldBool, FieldInt, FieldFloat, FieldMap:
		return true
	case FieldListObj:
		// Included so model.Field carries the typed block slice; the actual
		// expand/flatten wiring is handled by applyBlockTransforms.
		return true
	}
	return false
}

// applyFieldRenames replaces Get/Set calls from old field names to new field names.
// It handles both quoted string keys like "old_name" and bar d.Get("old_name") patterns.
func applyFieldRenames(body string, renames map[string]string) string {
	// For each rename mapping, replace quoted field name keys.
	// We replace both d.Get("oldName") and d.Set("oldName", ...) patterns.
	for oldName, newName := range renames {
		// Replace "oldName" in Get calls: metadata.ResourceData.Get("oldName")
		oldKey := `"` + oldName + `"`
		newKey := `"` + newName + `"`
		body = strings.ReplaceAll(body, oldKey, newKey)
	}
	return body
}

// --- Create / Update (Decode pass) -------------------------------------------

// applyDecodeTransform replaces metadata.ResourceData.Get("field").(type) with
// model.Field for all transformable fields and injects metadata.Decode at the
// start of the body.
func applyDecodeTransform(body, modelName string, fm map[string]*SchemaField) string {
	// Handle tags first, before the general Get pass replaces the inner Get
	// call and breaks the longer tags.Expand pattern.
	if _, hasTags := fm["tags"]; hasTags {
		body = replaceTags_Expand(body)
	}

	// Replace Get calls for each transformable non-tags field.
	for tfName, f := range fm {
		if tfName == "tags" {
			continue
		}
		body = replaceGetCalls(body, tfName, "model."+f.GoField, f.Kind)
	}

	// Remove the old TODO comment block (if present) and inject real Decode.
	body = removeTODOComment(body, modelName)
	decodeHeader := "\tvar model " + modelName + "\n\tif err := metadata.Decode(&model); err != nil {\n\t\treturn err\n\t}\n"
	body = decodeHeader + body

	return body
}

// replaceGetCalls replaces all occurrences of:
//
//	metadata.ResourceData.Get("tfName").(TYPE)
//	metadata.ResourceData.GetOk("tfName")
//	int64(metadata.ResourceData.Get("tfName").(TYPE))
//
// GetOk patterns are rewritten to a non-zero guard:
//
//	if v, ok := metadata.ResourceData.GetOk("tfName"); ok { target = v.(T) }
//	→ if model.GoField != <zero> { target = model.GoField }
func replaceGetCalls(body, tfName, modelExpr string, kind FieldKind) string {
	quoted := regexp.QuoteMeta(tfName)

	// Full pattern: metadata.ResourceData.Get("tfName").(anyType)
	re := regexp.MustCompile(`metadata\.ResourceData\.Get\("` + quoted + `"\)\.\([^)]*\)`)
	body = re.ReplaceAllString(body, modelExpr)

	// GetOk init: `if v, ok := metadata.ResourceData.GetOk("tfName"); ok {`
	// → `if model.Field != "" {`  (string zero)  or `if model.Field {` (bool)
	reGetOkInit := regexp.MustCompile(
		`if \w+, ok := metadata\.ResourceData\.GetOk\("` + quoted + `"\); ok \{`)
	zeroGuard := getOkZeroGuard(modelExpr, kind)
	body = reGetOkInit.ReplaceAllString(body, "if "+zeroGuard+" {")

	// Bare GetOk check: `metadata.ResourceData.GetOk("tfName")`
	reGetOk := regexp.MustCompile(`metadata\.ResourceData\.GetOk\("` + quoted + `"\)`)
	body = reGetOk.ReplaceAllString(body, modelExpr+`, true`)

	// Inside a former GetOk block the captured var (v, scriptURL, etc.) may still
	// be used. Replace `v.(string)` / `localVar.(string)` → modelExpr after the
	// guard was collapsed.  We do this with a broad match on any non-model
	// identifier followed by `.(type)`.
	reVarAssertion := regexp.MustCompile(`\b(?:v|ok|` + regexp.QuoteMeta(strings.TrimPrefix(modelExpr, "model.")) + `)\b\.\([^)]+\)`)
	body = reVarAssertion.ReplaceAllString(body, modelExpr)

	// Also strip a surrounding int64(...) when the model field is already int64.
	if kind == FieldInt {
		reInt := regexp.MustCompile(`int64\(` + regexp.QuoteMeta(modelExpr) + `\)`)
		body = reInt.ReplaceAllString(body, modelExpr)
	}

	return body
}

// getOkZeroGuard returns the guard expression for a GetOk replacement.
func getOkZeroGuard(modelExpr string, kind FieldKind) string {
	switch kind {
	case FieldBool:
		return modelExpr
	case FieldInt, FieldFloat:
		return modelExpr + " != 0"
	default:
		return modelExpr + ` != ""`
	}
}

// replaceTags_Expand replaces the common untyped tags patterns in Create/Update
// bodies:
//
//	tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{}))
//	tags.Expand(model.Tags)   (if Get was already replaced)
//
// In the typed pattern the SDK typically wants *map[string]string; the
// replacement emits model.Tags with a NOTE comment so the caller can add `&`.
func replaceTags_Expand(body string) string {
	// Pre-replacement form (Get still present).
	re := regexp.MustCompile(`tags\.Expand\(metadata\.ResourceData\.Get\("tags"\)\.\([^)]*\)\)`)
	body = re.ReplaceAllString(body, "model.Tags /* NOTE: SDK may need &model.Tags for *map[string]string */")

	// Post-replacement form (if Get was already substituted).
	re2 := regexp.MustCompile(`tags\.Expand\(model\.Tags\)`)
	body = re2.ReplaceAllString(body, "model.Tags /* NOTE: SDK may need &model.Tags for *map[string]string */")

	return body
}

// removeTODOComment strips the generated TODO block that was inserted as a
// placeholder before the model was properly decoded.
func removeTODOComment(body, modelName string) string {
	body = strings.ReplaceAll(body,
		"\t\t\t// TODO: use a model: var model "+modelName+"\n",
		"")
	body = strings.ReplaceAll(body,
		"\t\t\t// if err := metadata.Decode(&model); err != nil { return err }\n",
		"")
	return body
}

// --- Read (Encode pass) -------------------------------------------------------

// applyEncodeTransform replaces metadata.ResourceData.Set("field", value) with
// state.Field = value for all transformable fields and wires up Decode/Encode.
func applyEncodeTransform(body, modelName string, fm map[string]*SchemaField) string {
	// Replace Set calls line by line (value expression may contain parens).
	body = replaceSetCalls(body, fm)

	// Replace tags.FlattenAndSet(metadata.ResourceData, expr) → state.Tags = pointer.From(expr)
	if _, hasTags := fm["tags"]; hasTags {
		body = replaceTags_FlattenAndSet(body)
	}

	// Inject the state declaration and wire up Encode at the end of the body.
	body = injectStateAndEncode(body, modelName)

	return body
}

// replaceSetCalls replaces metadata.ResourceData.Set("field", value) with
// state.Field = value for each transformable field. The value expression is
// extracted by counting parentheses so nested calls are preserved correctly.
//
// When the value is a bare selector expression (e.g. props.PrincipalName) it
// is wrapped with pointer.From() so that *T SDK fields are safely dereferenced
// into the plain-type model field. If the SDK field is already a non-pointer
// the compiler will flag the unnecessary From() call, making the issue visible.
func replaceSetCalls(body string, fm map[string]*SchemaField) string {
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		for tfName, f := range fm {
			needle := `metadata.ResourceData.Set("` + tfName + `", `
			idx := strings.Index(line, needle)
			if idx < 0 {
				continue
			}
			valueStart := idx + len(needle)
			value, end := extractArg(line, valueStart)
			if value == "" {
				continue
			}
			indent := line[:len(line)-len(strings.TrimLeft(line, "\t "))]
			suffix := ""
			if end < len(line) {
				tail := strings.TrimSpace(line[end:])
				if tail != "" && tail != ")" {
					suffix = " " + tail
				}
			}
			// Wrap bare selector expressions (e.g. props.Field) with pointer.From()
			// so *T → T assignment compiles. String conversions, function calls,
			// and pointer.To/From wrappers are left as-is.
			rhs := wrapPointerIfBareSelector(value, f.Kind)
			lines[i] = indent + "state." + f.GoField + " = " + rhs + suffix
			break
		}
	}
	return strings.Join(lines, "\n")
}

// wrapPointerIfBareSelector wraps value with pointer.From() when it looks like
// a bare selector expression (pkg.Field or var.Field) that might be a *T in
// the SDK. Calls, type conversions, and already-wrapped expressions are left
// alone.
func wrapPointerIfBareSelector(value string, kind FieldKind) string {
	trimmed := strings.TrimSpace(value)
	// Already has pointer.From / pointer.To — leave as-is.
	if strings.HasPrefix(trimmed, "pointer.") {
		return value
	}
	// Type conversion: string(...), int64(...) etc. — leave as-is.
	if strings.Contains(trimmed, "(") {
		return value
	}
	// Plain identifier (no dot) — leave as-is (it's a local var).
	if !strings.Contains(trimmed, ".") {
		return value
	}
	// Bare selector: X.Y — potentially a *T pointer field from the SDK.
	// Only wrap for plain scalar model types (not map or slice).
	switch kind {
	case FieldString, FieldBool, FieldInt, FieldFloat:
		return "pointer.From(" + trimmed + ")"
	}
	return value
}

// extractArg extracts a function argument starting at pos in line by counting
// parentheses depth. It returns the argument string and the position after the
// closing delimiter (comma or paren of the outer call). Returns ("", 0) on failure.
func extractArg(line string, pos int) (string, int) {
	depth := 0
	start := pos
	for i := pos; i < len(line); i++ {
		ch := line[i]
		switch ch {
		case '(':
			depth++
		case ')':
			if depth == 0 {
				// Closing paren of the outer Set call.
				return strings.TrimSpace(line[start:i]), i + 1
			}
			depth--
		}
	}
	return "", 0
}

// replaceTags_FlattenAndSet replaces:
//
//	if err := tags.FlattenAndSet(metadata.ResourceData, expr); err != nil { ... }
//
// and the simpler:
//
//	tags.FlattenAndSet(metadata.ResourceData, expr)
//
// with state.Tags = pointer.From(expr).
func replaceTags_FlattenAndSet(body string) string {
	// if err := tags.FlattenAndSet(metadata.ResourceData, expr); err != nil { return err }
	re := regexp.MustCompile(`(?m)^\s*if err := tags\.FlattenAndSet\(metadata\.ResourceData, ([^)]+)\); err != nil \{[^}]*\}`)
	body = re.ReplaceAllStringFunc(body, func(m string) string {
		inner := re.FindStringSubmatch(m)
		if len(inner) < 2 {
			return m
		}
		indent := m[:len(m)-len(strings.TrimLeft(m, "\t "))]
		return indent + "state.Tags = pointer.From(" + strings.TrimSpace(inner[1]) + ")"
	})

	// Bare: tags.FlattenAndSet(metadata.ResourceData, expr)
	re2 := regexp.MustCompile(`tags\.FlattenAndSet\(metadata\.ResourceData, ([^)]+)\)`)
	body = re2.ReplaceAllStringFunc(body, func(m string) string {
		inner := re2.FindStringSubmatch(m)
		if len(inner) < 2 {
			return m
		}
		return "state.Tags = pointer.From(" + strings.TrimSpace(inner[1]) + ")"
	})
	return body
}

// injectStateAndEncode adds:
//  1. `state := ModelName{}` just before the first state.Field assignment
//     (or at the very start if none found).
//  2. Replaces the last `return nil` in the body with `return metadata.Encode(&state)`.
func injectStateAndEncode(body, modelName string) string {
	lines := strings.Split(body, "\n")

	// Find the first line that writes to the state variable.
	firstStateIdx := -1
	for i, l := range lines {
		if strings.Contains(l, "state.") {
			firstStateIdx = i
			break
		}
	}

	// Inject state declaration.
	decl := "\tstate := " + modelName + "{}"
	if firstStateIdx > 0 {
		// Insert one line before the first state assignment.
		before := lines[:firstStateIdx]
		after := lines[firstStateIdx:]
		lines = append(append([]string{}, before...), append([]string{decl}, after...)...)
	} else if firstStateIdx == 0 {
		lines = append([]string{decl}, lines...)
	} else {
		// No state assignments found — still inject so the Encode is correct.
		lines = append([]string{decl}, lines...)
	}

	// Replace the last `return nil` with `return metadata.Encode(&state)`.
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == "return nil" {
			indent := lines[i][:len(lines[i])-len(strings.TrimLeft(lines[i], "\t "))]
			lines[i] = indent + "return metadata.Encode(&state)"
			break
		}
	}

	return strings.Join(lines, "\n")
}
