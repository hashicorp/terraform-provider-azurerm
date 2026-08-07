// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package typed_upgrade

import (
	"fmt"
	"strings"
)

const fileHeader = "// Copyright IBM Corp. 2014, 2025\n// SPDX-License-Identifier: MPL-2.0\n\n"

// Generate renders the typed resource Go source for the given Info.
// The returned source is not yet formatted; callers should run goimports.
func Generate(info *Info) (string, error) {
	var sb strings.Builder

	sb.WriteString(fileHeader)
	fmt.Fprintf(&sb, "package %s\n\n", info.Package)
	renderImports(&sb, info)
	sb.WriteString("\n")
	renderInterfaceAssertions(&sb, info)
	sb.WriteString("\n\n")
	renderStructAndMethods(&sb, info)
	renderModelStructs(&sb, info)
	renderSchemaMethods(&sb, info)
	renderCRUDMethods(&sb, info)
	renderAdditionalInterfaces(&sb, info)

	// Append typed expand/flatten helpers generated from the Pandora IR.
	if info.ExpandFlattenSrc != "" {
		sb.WriteString(info.ExpandFlattenSrc)
	}

	return sb.String(), nil
}

func renderImports(sb *strings.Builder, info *Info) {
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\t\"fmt\"\n")
	sb.WriteString("\t\"time\"\n\n")
	for path, alias := range info.Imports {
		if shouldDropImport(path) {
			continue
		}
		if alias != "" {
			fmt.Fprintf(sb, "\t%s %q\n", alias, path)
		} else {
			fmt.Fprintf(sb, "\t%q\n", path)
		}
	}
	sb.WriteString("\t\"github.com/hashicorp/terraform-provider-azurerm/internal/sdk\"\n")
	sb.WriteString(")\n")
}

// shouldDropImport reports whether an import from the original untyped resource
// is superseded by the typed SDK wrapper.
func shouldDropImport(path string) bool {
	switch path {
	case "github.com/hashicorp/terraform-provider-azurerm/internal/timeouts",
		"github.com/hashicorp/terraform-provider-azurerm/internal/clients",
		"context",
		"time",
		"fmt":
		return true
	}
	return false
}

func renderInterfaceAssertions(sb *strings.Builder, info *Info) {
	sb.WriteString("var (\n")
	fmt.Fprintf(sb, "\t_ sdk.Resource = %s{}\n", info.StructName)
	if info.HasUpdate {
		fmt.Fprintf(sb, "\t_ sdk.ResourceWithUpdate = %s{}\n", info.StructName)
	}
	if info.HasCustomizeDiff {
		fmt.Fprintf(sb, "\t_ sdk.ResourceWithCustomizeDiff = %s{}\n", info.StructName)
	}
	if info.HasSchemaVersion && info.SchemaVersion > 0 {
		fmt.Fprintf(sb, "\t_ sdk.ResourceWithStateMigration = %s{}\n", info.StructName)
	}
	sb.WriteString(")")
}

func renderStructAndMethods(sb *strings.Builder, info *Info) {
	fmt.Fprintf(sb, "type %s struct{}\n\n", info.StructName)

	fmt.Fprintf(sb, "func (r %s) ResourceType() string {\n", info.StructName)
	fmt.Fprintf(sb, "\treturn %q\n}\n\n", info.TerraformType)

	fmt.Fprintf(sb, "func (r %s) ModelObject() interface{} {\n", info.StructName)
	fmt.Fprintf(sb, "\treturn &%s{}\n}\n\n", info.ModelName)

	fmt.Fprintf(sb, "func (r %s) IDValidationFunc() pluginsdk.SchemaValidateFunc {\n", info.StructName)
	fmt.Fprintf(sb, "\treturn %s\n}\n\n", idValidateExpr(info))
}

func idValidateExpr(info *Info) string {
	if info.IDValidateFunc == "" {
		return "nil // TODO: set the correct ID validation function"
	}
	pkg := info.IDPackage
	if pkg == "" {
		pkg = info.SDKPackage
	}
	if pkg == "" {
		return info.IDValidateFunc
	}
	return pkg + "." + info.IDValidateFunc
}

func renderModelStructs(sb *strings.Builder, info *Info) {
	fmt.Fprintf(sb, "type %s struct {\n", info.ModelName)
	if len(info.Fields) == 0 {
		sb.WriteString("\t// TODO: add model fields\n")
	} else {
		for _, f := range info.Fields {
			goType := f.Kind.GoType(f.BlockName)
			fmt.Fprintf(sb, "\t%s %s `tfschema:%q`\n", f.GoField, goType, f.TFName)
		}
	}
	sb.WriteString("}\n\n")

	// If the Pandora IR is available, emit fully-typed block model structs.
	// Otherwise emit the placeholder stubs for any FieldListObj fields.
	if info.PandoraIR != nil {
		sb.WriteString(renderBlockModelStructs(info))
		return
	}
	seen := map[string]bool{}
	for _, f := range info.Fields {
		if f.Kind != FieldListObj || f.BlockName == "" || seen[f.BlockName] {
			continue
		}
		seen[f.BlockName] = true
		fmt.Fprintf(sb, "type %s struct {\n", f.BlockName)
		sb.WriteString("\t// TODO: add nested block model fields\n")
		sb.WriteString("}\n\n")
	}
}

func renderSchemaMethods(sb *strings.Builder, info *Info) {
	renderSchemaMethod(sb, info, "Arguments", false)
	renderSchemaMethod(sb, info, "Attributes", true)
}

func renderSchemaMethod(sb *strings.Builder, info *Info, name string, attributesOnly bool) {
	fmt.Fprintf(sb, "func (r %s) %s() map[string]*pluginsdk.Schema {\n", info.StructName, name)
	sb.WriteString("\treturn map[string]*pluginsdk.Schema{\n")
	for _, f := range info.Fields {
		if attributesOnly != f.IsAttribute() {
			continue
		}
		fmt.Fprintf(sb, "\t\t%q: %s,\n", f.TFName, f.RawSchema)
	}
	sb.WriteString("\t}\n}\n\n")
}

func renderCRUDMethods(sb *strings.Builder, info *Info) {
	renderResourceFunc(sb, info, "Create", info.CreateTimeout, info.ModelName, info.CreateBody)
	renderResourceFunc(sb, info, "Read", info.ReadTimeout, info.ModelName, info.ReadBody)
	if info.HasUpdate {
		renderResourceFunc(sb, info, "Update", info.UpdateTimeout, info.ModelName, info.UpdateBody)
	}
	renderResourceFunc(sb, info, "Delete", info.DeleteTimeout, info.ModelName, info.DeleteBody)
}

func renderResourceFunc(sb *strings.Builder, info *Info, op, timeout, modelName, body string) {
	fmt.Fprintf(sb, "func (r %s) %s() sdk.ResourceFunc {\n", info.StructName, op)
	sb.WriteString("\treturn sdk.ResourceFunc{\n")
	fmt.Fprintf(sb, "\t\tTimeout: %s,\n", timeout)
	sb.WriteString("\t\tFunc: func(ctx context.Context, metadata sdk.ResourceMetaData) error {\n")
	if body != "" {
		for _, line := range strings.Split(body, "\n") {
			if strings.TrimSpace(line) == "" {
				sb.WriteString("\n")
			} else {
				sb.WriteString("\t\t\t" + strings.TrimPrefix(line, "\t") + "\n")
			}
		}
	}
	sb.WriteString("\t\t},\n")
	sb.WriteString("\t}\n}\n\n")
}

func renderAdditionalInterfaces(sb *strings.Builder, info *Info) {
	if info.HasCustomizeDiff {
		fmt.Fprintf(sb, "func (r %s) CustomizeDiff() sdk.ResourceFunc {\n", info.StructName)
		sb.WriteString("\t// TODO: migrate the original CustomizeDiff logic here\n")
		sb.WriteString("\tpanic(\"not implemented\")\n")
		sb.WriteString("}\n\n")
	}
	if info.HasSchemaVersion && info.SchemaVersion > 0 {
		fmt.Fprintf(sb, "func (r %s) StateUpgraders() sdk.StateUpgradeData {\n", info.StructName)
		sb.WriteString("\t// TODO: carry over state upgraders from the original resource.\n")
		fmt.Fprintf(sb, "\t// SchemaVersion must remain %d to match the original.\n", info.SchemaVersion)
		sb.WriteString("\treturn sdk.StateUpgradeData{\n")
		fmt.Fprintf(sb, "\t\tSchemaVersion: %d,\n", info.SchemaVersion)
		sb.WriteString("\t\tUpgraders: map[int]pluginsdk.StateUpgrade{},\n")
		sb.WriteString("\t}\n}\n\n")
	}
}
