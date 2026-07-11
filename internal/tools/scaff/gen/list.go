// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gen

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// RenderListResource renders the full source of a list resource file, choosing
// the subscription/resource-group scoped shape or the parent-scoped shape based
// on the operations discovered for the resource.
func RenderListResource(res *ir.ResourceIR) string {
	var body string
	switch {
	case res.Untyped:
		body = renderUntypedList(res)
	case res.ListBySubscriptionOp != "" || res.ListByResourceGroupOp != "":
		body = renderSubscriptionScopedList(res)
	default:
		body = renderParentScopedList(res)
	}
	return fileHeader + body
}

// renderUntypedList renders a list resource for a native Plugin SDK (untyped)
// resource. It wraps the resource's own constructor and flatten function rather
// than the internal/sdk typed wrapper, choosing the parent-scoped shape for
// child resources and the subscription/resource-group shape otherwise.
func renderUntypedList(res *ir.ResourceIR) string {
	if res.ListByParentOp != "" {
		return renderUntypedParentScopedList(res)
	}
	return renderUntypedSubscriptionScopedList(res)
}

// renderUntypedListImports writes the import block shared by the untyped list
// resource shapes.
func renderUntypedListImports(sb *strings.Builder, res *ir.ResourceIR, parentScoped bool) {
	sb.WriteString("import (\n")
	sb.WriteString("\"context\"\n\"fmt\"\n\n")
	if parentScoped {
		sb.WriteString("\"github.com/hashicorp/go-azure-helpers/framework/typehelpers\"\n")
	}
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/lang/pointer\"\n")

	// Resource-manager + commonids imports, de-duplicated: the resource id and
	// parent id may share a package (e.g. both in commonids for subnet), and the
	// id may live in the SDK package itself.
	const commonidsPath = "github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	emitted := map[string]bool{}
	emit := func(path string) {
		if path == "" || emitted[path] {
			return
		}
		emitted[path] = true
		fmt.Fprintf(sb, "%q\n", path)
	}
	if !parentScoped {
		// Subscription/resource-group lists build a commonids.NewSubscriptionID.
		emit(commonidsPath)
	}
	emit(res.SDKImportPath)
	emit(res.IDImportPath)
	if parentScoped {
		emit(res.ParentImportPath)
	}
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/list\"\n")
	if parentScoped {
		sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/list/schema\"\n")
	}
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/resource\"\n")
	if parentScoped {
		sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/schema/validator\"\n")
		sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/types\"\n")
	}
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-sdk/v2/terraform\"\n")
	fmt.Fprintf(sb, "\"github.com/%s/terraform-provider-%s/internal/sdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	fmt.Fprintf(sb, "\"github.com/%s/terraform-provider-%s/internal/tf/pluginsdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	sb.WriteString(")\n\n")
}

// renderUntypedPreamble writes the interface assertion, Metadata and
// ResourceFunc shared by the untyped list resource shapes.
func renderUntypedPreamble(sb *strings.Builder, res *ir.ResourceIR) {
	fmt.Fprintf(sb, "var _ sdk.FrameworkListWrappedResource = new(%sListResource)\n\n", res.Name)
	fmt.Fprintf(sb, "func (%sListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {\n", res.Name)
	fmt.Fprintf(sb, "response.TypeName = %q\n}\n\n", res.TerraformType)
	fmt.Fprintf(sb, "func (%sListResource) ResourceFunc() *pluginsdk.Resource {\n", res.Name)
	fmt.Fprintf(sb, "return %s()\n}\n\n", res.ConstructorFunc)
}

// renderUntypedSubscriptionScopedList renders an untyped list resource for a
// top-level resource listed by subscription and/or resource group.
func renderUntypedSubscriptionScopedList(res *ir.ResourceIR) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", res.ServicePackage)
	renderUntypedListImports(&sb, res, false)

	fmt.Fprintf(&sb, "type %sListResource struct{}\n\n", res.Name)
	renderUntypedPreamble(&sb, res)

	fmt.Fprintf(&sb, "func (%sListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {\n", res.Name)
	fmt.Fprintf(&sb, "client := metadata.Client.%s.%s\n\n", res.ServiceName, res.ClientField)
	sb.WriteString("var data sdk.DefaultListModel\n")
	sb.WriteString("diags := request.Config.Get(ctx, &data)\n")
	sb.WriteString("if diags.HasError() {\nstream.Results = list.ListResultsStreamDiagnostics(diags)\nreturn\n}\n\n")
	fmt.Fprintf(&sb, "var results []%s.%s\n", res.SDKPackage, res.ReadModel)
	sb.WriteString("subscriptionID := metadata.SubscriptionId\n")
	sb.WriteString("if !data.SubscriptionId.IsNull() {\nsubscriptionID = data.SubscriptionId.ValueString()\n}\n\n")

	rgCall := func() {
		fmt.Fprintf(&sb, "resp, err := client.%sComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))\n", res.ListByResourceGroupOp)
		fmt.Fprintf(&sb, "if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, %q, err)\nreturn\n}\n", "listing "+res.TerraformType)
		sb.WriteString("results = resp.Items\n")
	}
	subCall := func() {
		fmt.Fprintf(&sb, "resp, err := client.%sComplete(ctx, commonids.NewSubscriptionID(subscriptionID))\n", res.ListBySubscriptionOp)
		fmt.Fprintf(&sb, "if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, %q, err)\nreturn\n}\n", "listing "+res.TerraformType)
		sb.WriteString("results = resp.Items\n")
	}
	switch {
	case res.ListByResourceGroupOp != "" && res.ListBySubscriptionOp != "":
		sb.WriteString("switch {\ncase !data.ResourceGroupName.IsNull():\n")
		rgCall()
		sb.WriteString("default:\n")
		subCall()
		sb.WriteString("}\n\n")
	case res.ListByResourceGroupOp != "":
		rgCall()
		sb.WriteString("\n")
	default:
		subCall()
		sb.WriteString("\n")
	}

	renderUntypedListStream(&sb, res, "results")
	sb.WriteString("}\n")
	return sb.String()
}

// renderUntypedParentScopedList renders an untyped list resource for a child
// resource listed under a parent (a required parent ID is supplied via the list
// config); it issues a single parent-scoped List call.
func renderUntypedParentScopedList(res *ir.ResourceIR) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", res.ServicePackage)
	renderUntypedListImports(&sb, res, true)

	fmt.Fprintf(&sb, "type %sListResource struct{}\n\n", res.Name)
	fmt.Fprintf(&sb, "type %sListModel struct {\n%s types.String `tfsdk:%q`\n}\n\n", res.Name, res.ParentIDType, res.ParentListAttr)
	renderUntypedPreamble(&sb, res)

	fmt.Fprintf(&sb, "func (%sListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {\n", res.Name)
	sb.WriteString("response.Schema = schema.Schema{\nAttributes: map[string]schema.Attribute{\n")
	fmt.Fprintf(&sb, "%q: schema.StringAttribute{\nRequired: true,\nValidators: []validator.String{\n", res.ParentListAttr)
	fmt.Fprintf(&sb, "typehelpers.WrappedStringValidator{Func: %s.%s},\n", res.ParentPackage, res.ParentValidateFunc)
	sb.WriteString("},\n},\n},\n}\n}\n\n")

	fmt.Fprintf(&sb, "func (%sListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {\n", res.Name)
	fmt.Fprintf(&sb, "client := metadata.Client.%s.%s\n\n", res.ServiceName, res.ClientField)
	fmt.Fprintf(&sb, "var data %sListModel\n", res.Name)
	sb.WriteString("diags := request.Config.Get(ctx, &data)\n")
	sb.WriteString("if diags.HasError() {\nstream.Results = list.ListResultsStreamDiagnostics(diags)\nreturn\n}\n\n")
	fmt.Fprintf(&sb, "parentID, err := %s.%s(data.%s.ValueString())\n", res.ParentPackage, res.ParentParseFunc, res.ParentIDType)
	fmt.Fprintf(&sb, "if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, %q, err)\nreturn\n}\n\n", "parsing parent ID for "+res.TerraformType)
	fmt.Fprintf(&sb, "resp, err := client.%sComplete(ctx, *parentID)\n", res.ListByParentOp)
	fmt.Fprintf(&sb, "if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, %q, err)\nreturn\n}\n\n", "listing "+res.TerraformType)

	renderUntypedListStream(&sb, res, "resp.Items")
	sb.WriteString("}\n")
	return sb.String()
}

// renderUntypedListStream renders the stream.Results closure for a native
// Plugin SDK list resource, building a ResourceData from the resource's own
// constructor and populating it via its flatten function.
func renderUntypedListStream(sb *strings.Builder, res *ir.ResourceIR, itemsExpr string) {
	idPkg := res.IDPackage
	if idPkg == "" {
		idPkg = res.SDKPackage
	}
	sb.WriteString("stream.Results = func(push func(list.ListResult) bool) {\n")
	fmt.Fprintf(sb, "for _, item := range %s {\n", itemsExpr)
	sb.WriteString("result := request.NewListResult(ctx)\n")
	sb.WriteString("result.DisplayName = pointer.From(item.Name)\n\n")
	fmt.Fprintf(sb, "rd := %s().Data(&terraform.InstanceState{})\n", res.ConstructorFunc)
	fmt.Fprintf(sb, "id, err := %s.%sInsensitively(pointer.From(item.Id))\n", idPkg, res.IDParseFunc)
	fmt.Fprintf(sb, "if err != nil {\nsdk.SetErrorDiagnosticAndPushListResult(result, push, %q, err)\nreturn\n}\n", "parsing "+res.TerraformType+" ID")
	sb.WriteString("rd.SetId(id.ID())\n\n")
	idArg := "id"
	if res.FlattenIDValue {
		idArg = "*id"
	}
	fmt.Fprintf(sb, "if err := %s(rd, %s, &item); err != nil {\nsdk.SetErrorDiagnosticAndPushListResult(result, push, %q, err)\nreturn\n}\n\n", res.FlattenFunc, idArg, "encoding "+res.TerraformType+" resource data")
	sb.WriteString("sdk.EncodeListResult(ctx, rd, &result)\n")
	sb.WriteString("if result.Diagnostics.HasError() {\npush(result)\nreturn\n}\n")
	sb.WriteString("if !push(result) {\nreturn\n}\n")
	sb.WriteString("}\n}\n")
}

// renderSubscriptionScopedList renders a list resource for a top-level resource
// listed by subscription and/or resource group.
func renderSubscriptionScopedList(res *ir.ResourceIR) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "package %s\n\n", res.ServicePackage)

	sb.WriteString("import (\n")
	sb.WriteString("\"context\"\n\"fmt\"\n\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/lang/pointer\"\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids\"\n")
	fmt.Fprintf(&sb, "%q\n", res.SDKImportPath)
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/list\"\n")
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/resource\"\n")
	fmt.Fprintf(&sb, "\"github.com/%s/terraform-provider-%s/internal/sdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	fmt.Fprintf(&sb, "\"github.com/%s/terraform-provider-%s/internal/tf/pluginsdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	sb.WriteString(")\n\n")

	renderListPreamble(&sb, res)

	fmt.Fprintf(&sb, "func (%sListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {\n", res.Name)
	fmt.Fprintf(&sb, "client := metadata.Client.%s.%s\n\n", res.ServiceName, res.ClientField)
	sb.WriteString("var data sdk.DefaultListModel\n")
	sb.WriteString("diags := request.Config.Get(ctx, &data)\n")
	sb.WriteString("if diags.HasError() {\nstream.Results = list.ListResultsStreamDiagnostics(diags)\nreturn\n}\n\n")
	fmt.Fprintf(&sb, "var results []%s.%s\n", res.SDKPackage, res.ReadModel)
	sb.WriteString("subscriptionID := metadata.SubscriptionId\n")
	sb.WriteString("if !data.SubscriptionId.IsNull() {\nsubscriptionID = data.SubscriptionId.ValueString()\n}\n\n")
	fmt.Fprintf(&sb, "resource := %sResource{}\n\n", res.Name)

	rgCall := func() {
		fmt.Fprintf(&sb, "resp, err := client.%sComplete(ctx, commonids.NewResourceGroupID(subscriptionID, data.ResourceGroupName.ValueString()))\n", res.ListByResourceGroupOp)
		sb.WriteString("if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf(\"listing `%s`\", resource.ResourceType()), err)\nreturn\n}\n")
		sb.WriteString("results = resp.Items\n")
	}
	subCall := func() {
		fmt.Fprintf(&sb, "resp, err := client.%sComplete(ctx, commonids.NewSubscriptionID(subscriptionID))\n", res.ListBySubscriptionOp)
		sb.WriteString("if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf(\"listing `%s`\", resource.ResourceType()), err)\nreturn\n}\n")
		sb.WriteString("results = resp.Items\n")
	}

	switch {
	case res.ListByResourceGroupOp != "" && res.ListBySubscriptionOp != "":
		sb.WriteString("switch {\ncase !data.ResourceGroupName.IsNull():\n")
		rgCall()
		sb.WriteString("default:\n")
		subCall()
		sb.WriteString("}\n\n")
	case res.ListByResourceGroupOp != "":
		rgCall()
		sb.WriteString("\n")
	default:
		subCall()
		sb.WriteString("\n")
	}

	renderListStream(&sb, res, "results")
	sb.WriteString("}\n")
	return sb.String()
}

// renderParentScopedList renders a list resource for a child resource listed by
// its parent (a required parent ID is supplied via the list config).
func renderParentScopedList(res *ir.ResourceIR) string {
	var sb strings.Builder
	parentPkg := res.ParentPackage
	if parentPkg == "" {
		parentPkg = res.SDKPackage
	}
	fmt.Fprintf(&sb, "package %s\n\n", res.ServicePackage)

	sb.WriteString("import (\n")
	sb.WriteString("\"context\"\n\"fmt\"\n\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/framework/typehelpers\"\n")
	sb.WriteString("\"github.com/hashicorp/go-azure-helpers/lang/pointer\"\n")
	fmt.Fprintf(&sb, "%q\n", res.SDKImportPath)
	if res.ParentImportPath != "" && res.ParentImportPath != res.SDKImportPath {
		fmt.Fprintf(&sb, "%q\n", res.ParentImportPath)
	}
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/list\"\n")
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/list/schema\"\n")
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/resource\"\n")
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/schema/validator\"\n")
	sb.WriteString("\"github.com/hashicorp/terraform-plugin-framework/types\"\n")
	fmt.Fprintf(&sb, "\"github.com/%s/terraform-provider-%s/internal/sdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	fmt.Fprintf(&sb, "\"github.com/%s/terraform-provider-%s/internal/tf/pluginsdk\"\n", res.ProviderGithubOrg, res.ProviderName)
	sb.WriteString(")\n\n")

	fmt.Fprintf(&sb, "type %sListResource struct{}\n\n", res.Name)
	fmt.Fprintf(&sb, "type %sListModel struct {\n%s types.String `tfsdk:%q`\n}\n\n", res.Name, res.ParentIDType, res.ParentListAttr)
	fmt.Fprintf(&sb, "var _ sdk.FrameworkListWrappedResource = new(%sListResource)\n\n", res.Name)

	fmt.Fprintf(&sb, "func (%sListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {\n", res.Name)
	fmt.Fprintf(&sb, "response.TypeName = %sResource{}.ResourceType()\n}\n\n", res.Name)

	fmt.Fprintf(&sb, "func (%sListResource) ResourceFunc() *pluginsdk.Resource {\n", res.Name)
	fmt.Fprintf(&sb, "return sdk.WrappedResource(%sResource{})\n}\n\n", res.Name)

	fmt.Fprintf(&sb, "func (%sListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {\n", res.Name)
	sb.WriteString("response.Schema = schema.Schema{\n")
	sb.WriteString("Attributes: map[string]schema.Attribute{\n")
	fmt.Fprintf(&sb, "%q: schema.StringAttribute{\n", res.ParentListAttr)
	sb.WriteString("Required: true,\n")
	sb.WriteString("Validators: []validator.String{\n")
	fmt.Fprintf(&sb, "typehelpers.WrappedStringValidator{Func: %s.%s},\n", parentPkg, res.ParentValidateFunc)
	sb.WriteString("},\n},\n},\n}\n}\n\n")

	fmt.Fprintf(&sb, "func (%sListResource) List(ctx context.Context, request list.ListRequest, stream *list.ListResultsStream, metadata sdk.ResourceMetadata) {\n", res.Name)
	fmt.Fprintf(&sb, "client := metadata.Client.%s.%s\n\n", res.ServiceName, res.ClientField)
	fmt.Fprintf(&sb, "var data %sListModel\n", res.Name)
	sb.WriteString("diags := request.Config.Get(ctx, &data)\n")
	sb.WriteString("if diags.HasError() {\nstream.Results = list.ListResultsStreamDiagnostics(diags)\nreturn\n}\n\n")
	fmt.Fprintf(&sb, "resource := %sResource{}\n\n", res.Name)
	fmt.Fprintf(&sb, "parentID, err := %s.%s(data.%s.ValueString())\n", parentPkg, res.ParentParseFunc, res.ParentIDType)
	sb.WriteString("if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf(\"parsing parent ID for `%s`\", resource.ResourceType()), err)\nreturn\n}\n\n")
	fmt.Fprintf(&sb, "resp, err := client.%sComplete(ctx, *parentID)\n", res.ListByParentOp)
	sb.WriteString("if err != nil {\nsdk.SetResponseErrorDiagnostic(stream, fmt.Sprintf(\"listing `%s`\", resource.ResourceType()), err)\nreturn\n}\n\n")

	renderListStream(&sb, res, "resp.Items")
	sb.WriteString("}\n")
	return sb.String()
}

// renderListPreamble renders the type, interface assertion, Metadata and
// ResourceFunc shared by the subscription-scoped list resource.
func renderListPreamble(sb *strings.Builder, res *ir.ResourceIR) {
	fmt.Fprintf(sb, "type %sListResource struct{}\n\n", res.Name)
	fmt.Fprintf(sb, "var _ sdk.FrameworkListWrappedResource = new(%sListResource)\n\n", res.Name)
	fmt.Fprintf(sb, "func (%sListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {\n", res.Name)
	fmt.Fprintf(sb, "response.TypeName = %sResource{}.ResourceType()\n}\n\n", res.Name)
	fmt.Fprintf(sb, "func (%sListResource) ResourceFunc() *pluginsdk.Resource {\n", res.Name)
	fmt.Fprintf(sb, "return sdk.WrappedResource(%sResource{})\n}\n\n", res.Name)
}

// renderListStream renders the shared stream.Results closure that iterates the
// list results, parses each ID, flattens the item and pushes the result.
func renderListStream(sb *strings.Builder, res *ir.ResourceIR, itemsExpr string) {
	sb.WriteString("stream.Results = func(push func(list.ListResult) bool) {\n")
	fmt.Fprintf(sb, "for _, item := range %s {\n", itemsExpr)
	sb.WriteString("result := request.NewListResult(ctx)\n")
	sb.WriteString("result.DisplayName = pointer.From(item.Name)\n\n")
	fmt.Fprintf(sb, "id, err := %s.%sInsensitively(pointer.From(item.Id))\n", res.SDKPackage, res.IDParseFunc)
	sb.WriteString("if err != nil {\nsdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf(\"parsing %s ID\", resource.ResourceType()), err)\nreturn\n}\n\n")
	sb.WriteString("meta := sdk.NewResourceMetaData(metadata.Client, resource)\n")
	sb.WriteString("meta.SetID(id)\n\n")
	sb.WriteString("if err := resource.flatten(meta, id, &item); err != nil {\nsdk.SetErrorDiagnosticAndPushListResult(result, push, fmt.Sprintf(\"encoding `%s` resource data\", resource.ResourceType()), err)\nreturn\n}\n\n")
	sb.WriteString("sdk.EncodeListResult(ctx, meta.ResourceData, &result)\n")
	sb.WriteString("if result.Diagnostics.HasError() {\npush(result)\nreturn\n}\n")
	sb.WriteString("if !push(result) {\nreturn\n}\n")
	sb.WriteString("}\n}\n")
}

// listTestData is the data supplied to the list-resource acceptance test template.
type listTestData struct {
	PackageName           string
	TestName              string
	ResourceStruct        string
	TerraformResourceName string
	UseResourceGroup      bool
	Parent                string
	ParentTFName          string
}

// RenderListTest renders the acceptance test for a generated list resource,
// adapted from the first-iteration list_test.go.gotpl template. It assumes the
// resource's own _test.go provides a basic(data) config method.
func RenderListTest(res *ir.ResourceIR) string {
	useRG := res.ListBySubscriptionOp != "" || res.ListByResourceGroupOp != ""
	resourceStruct := res.TestStructName
	if resourceStruct == "" {
		resourceStruct = res.Name + "Resource"
	}
	data := listTestData{
		PackageName:           res.ServicePackage,
		TestName:              res.Name,
		ResourceStruct:        resourceStruct,
		TerraformResourceName: res.TerraformType,
		UseResourceGroup:      useRG,
	}
	if !useRG {
		data.Parent = strings.TrimSuffix(res.ParentIDType, "Id")
		data.ParentTFName = strings.TrimSuffix(res.ParentListAttr, "_id")
	}

	funcs := template.FuncMap{"bt": func() string { return "`" }}
	t := template.Must(template.New("listtest").Funcs(funcs).Parse(listTestTemplate))
	var sb strings.Builder
	if err := t.Execute(&sb, data); err != nil {
		// The template is static and the data is trivial, so execution errors
		// are not expected in practice.
		return ""
	}
	return sb.String()
}

const listTestTemplate = fileHeader + `package {{.PackageName}}_test

import (
	"context"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/querycheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestAcc{{.TestName}}_listBy{{if .UseResourceGroup}}SubscriptionAndRG{{else}}{{.Parent}}ID{{end}}(t *testing.T) {
	data := acceptance.BuildTestData(t, "{{.TerraformResourceName}}", "testlist1")
	r := {{.ResourceStruct}}{}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: r.basic(data),
			},
{{if .UseResourceGroup}}			{
				Query:  true,
				Config: r.basicQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("{{.TerraformResourceName}}.list", 1),
					querycheck.ExpectIdentity(
						"{{.TerraformResourceName}}.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
			{
				Query:  true,
				Config: r.basicQueryByResourceGroupName(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLength("{{.TerraformResourceName}}.list", 1),
					querycheck.ExpectIdentity(
						"{{.TerraformResourceName}}.list",
						map[string]knownvalue.Check{
							"name":                knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":     knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
{{else}}			{
				Query:  true,
				Config: r.basicQuery(data),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectLengthAtLeast("{{.TerraformResourceName}}.list", 1),
					querycheck.ExpectIdentity(
						"{{.TerraformResourceName}}.list",
						map[string]knownvalue.Check{
							"name":                    knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"resource_group_name":     knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"{{.ParentTFName}}_name": knownvalue.StringRegexp(regexp.MustCompile(strconv.Itoa(data.RandomInteger))),
							"subscription_id":         knownvalue.StringExact(data.Subscriptions.Primary),
						},
					),
				},
			},
{{end}}		},
	})
}
{{if .UseResourceGroup}}
func (r {{.ResourceStruct}}) basicQuery() string {
	return {{bt}}
list "{{.TerraformResourceName}}" "list" {
  provider = azurerm
  config {}
}
{{bt}}
}

func (r {{.ResourceStruct}}) basicQueryByResourceGroupName() string {
	return {{bt}}
list "{{.TerraformResourceName}}" "list" {
  provider = azurerm
  config {
    resource_group_name = azurerm_resource_group.test.name
  }
}
{{bt}}
}
{{else}}
func (r {{.ResourceStruct}}) basicQuery(data acceptance.TestData) string {
	return {{bt}}
list "{{.TerraformResourceName}}" "list" {
  provider = azurerm
  config {
    {{.ParentTFName}}_id = azurerm_{{.ParentTFName}}.test.id
  }
}
{{bt}}
}
{{end}}`
