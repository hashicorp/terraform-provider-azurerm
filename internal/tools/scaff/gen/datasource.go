package gen

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// RenderDataSource renders a typed (internal/sdk) data source. The nested block
// structs are shared with the resource (defined in the resource file), so only
// the top-level {Name}DataSourceModel struct is declared here.
func RenderDataSource(res *ir.ResourceIR) string {
	dsName := res.Name + "DataSource"
	dsModel := res.Name + "DataSourceModel"

	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", res.ServicePackage)

	sb.WriteString("import (\n")
	sb.WriteString("\"context\"\n\"fmt\"\n\"time\"\n\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/lang/pointer\"\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/lang/response\"\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema\"\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/resourcemanager/location\"\n")
	fmt.Fprintf(&sb, "%q\n", res.SDKImportPath)
	fmt.Fprintf(&sb, "\"github.com/%s/terraform-provider-%s/internal/sdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	fmt.Fprintf(&sb, "\"github.com/%s/terraform-provider-%s/internal/tf/pluginsdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	sb.WriteString(")\n\n")

	fmt.Fprintf(&sb, "type %s struct{}\n\n", dsName)
	fmt.Fprintf(&sb, "var _ sdk.DataSource = %s{}\n\n", dsName)

	// Top-level model struct only; nested block structs are shared with the resource.
	fmt.Fprintf(&sb, "type %s struct {\n", dsModel)
	for _, p := range res.TopLevel {
		fmt.Fprintf(&sb, "\t%s %s `tfschema:%q`\n", p.GoField, p.GoType, p.TFName)
	}
	sb.WriteString("}\n\n")

	fmt.Fprintf(&sb, "func (r %s) ModelObject() interface{} {\nreturn &%s{}\n}\n\n", dsName, dsModel)
	fmt.Fprintf(&sb, "func (r %s) ResourceType() string {\nreturn %q\n}\n\n", dsName, res.TerraformType)

	// Arguments - the identifying inputs.
	fmt.Fprintf(&sb, "func (r %s) Arguments() map[string]*pluginsdk.Schema {\n", dsName)
	sb.WriteString("return map[string]*pluginsdk.Schema{\n")
	sb.WriteString("\"name\": {\nType: pluginsdk.TypeString,\nRequired: true,\n},\n")
	if res.HasResourceGroup {
		sb.WriteString("\"resource_group_name\": commonschema.ResourceGroupNameForDataSource(),\n")
	}
	sb.WriteString("}\n}\n\n")

	// Attributes - everything else, all computed.
	blocks := blockByName(res)
	fmt.Fprintf(&sb, "func (r %s) Attributes() map[string]*pluginsdk.Schema {\n", dsName)
	sb.WriteString("return map[string]*pluginsdk.Schema{\n")
	for _, p := range res.TopLevel {
		switch p.TFName {
		case "name", "resource_group_name":
			continue
		case "location":
			sb.WriteString("\"location\": commonschema.LocationComputed(),\n")
		case "tags":
			sb.WriteString("\"tags\": commonschema.TagsDataSource(),\n")
		default:
			renderSchemaProp(&sb, p, blocks, true)
		}
	}
	sb.WriteString("}\n}\n\n")

	renderDataSourceRead(&sb, res, dsName, dsModel)
	return sb.String()
}

func renderDataSourceRead(sb *strings.Builder, res *ir.ResourceIR, dsName, dsModel string) {
	fmt.Fprintf(sb, "func (r %s) Read() sdk.ResourceFunc {\n", dsName)
	sb.WriteString("return sdk.ResourceFunc{\nTimeout: 5 * time.Minute,\n")
	sb.WriteString("Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {\n")
	fmt.Fprintf(sb, "client := metadata.Client.%s.%s\n", res.ServiceName, res.ClientField)
	if res.HasSubscription {
		sb.WriteString("subscriptionId := metadata.Client.Account.SubscriptionId\n")
	}
	sb.WriteString("\n")
	fmt.Fprintf(sb, "var state %s\n", dsModel)
	sb.WriteString("if err := metadata.Decode(&state); err != nil {\nreturn err\n}\n\n")

	// The New*ID args are computed against the create model var "config"; the
	// data source decodes into "state".
	args := make([]string, len(res.IDNewArgs))
	for i, a := range res.IDNewArgs {
		args[i] = strings.ReplaceAll(a, "config.", "state.")
	}
	fmt.Fprintf(sb, "id := %s.%s(%s)\n\n", res.SDKPackage, res.IDNewFunc, strings.Join(args, ", "))

	fmt.Fprintf(sb, "existing, err := client.%s(ctx, id)\n", res.ReadOp)
	sb.WriteString("if err != nil {\nif response.WasNotFound(existing.HttpResponse) {\nreturn fmt.Errorf(\"%s was not found\", id)\n}\nreturn fmt.Errorf(\"retrieving %s: %+v\", id, err)\n}\n\n")

	if res.IDNameSegment != "" {
		fmt.Fprintf(sb, "state.Name = id.%s\n", res.IDNameSegment)
	}
	if res.HasResourceGroup {
		sb.WriteString("state.ResourceGroup = id.ResourceGroupName\n")
	}
	sb.WriteString("\n")

	sb.WriteString("if model := existing.Model; model != nil {\n")
	for _, p := range res.TopLevel {
		if p.UnderProperties {
			continue
		}
		switch p.TFName {
		case "location":
			sb.WriteString("state.Location = location.Normalize(model.Location)\n")
		case "tags":
			sb.WriteString("state.Tags = pointer.From(model.Tags)\n")
		}
	}
	if res.PropertiesModel != "" {
		sb.WriteString("if props := model.Properties; props != nil {\n")
		for _, p := range res.TopLevel {
			if !p.UnderProperties {
				continue
			}
			fmt.Fprintf(sb, "state.%s = %s\n", p.GoField, flattenFieldExpr(res, p, "props"))
		}
		sb.WriteString("}\n")
	}
	sb.WriteString("}\n\n")

	sb.WriteString("metadata.SetID(id)\n")
	sb.WriteString("return metadata.Encode(&state)\n")
	sb.WriteString("},\n}\n}\n")
}
