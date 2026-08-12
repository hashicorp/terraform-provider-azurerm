// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package gen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
)

// fileData is the data passed to the file templates. It embeds the ResourceIR
// and adds pre-rendered Go source fragments so the templates stay declarative.
type fileData struct {
	*ir.ResourceIR
	Package       string
	ArgsJoined    string
	Arguments     string
	Attributes    string
	ModelStructs  string
	ExpandFlatten string
	CreatePayload string
	FlattenMethod string
	UpdateBody    string
}

// Generate renders the resource Go source (schema, model structs, CRUD and
// expand/flatten) for a resolved ResourceIR. The returned source is not yet
// formatted; callers should run goimports/gofmt.
func Generate(res *ir.ResourceIR) (resourceGo string, err error) {
	data := fileData{
		ResourceIR:    res,
		Package:       res.ServicePackage,
		ArgsJoined:    strings.Join(res.IDNewArgs, ", "),
		Arguments:     RenderArguments(res),
		Attributes:    RenderAttributes(res),
		ModelStructs:  RenderModelStructs(res),
		ExpandFlatten: RenderExpandFlatten(res),
		CreatePayload: RenderCreatePayload(res),
		FlattenMethod: RenderFlattenMethod(res),
		UpdateBody:    RenderUpdateBody(res),
	}

	resourceGo, err = execTemplate("resource", resourceTemplate, data)
	if err != nil {
		return "", err
	}
	return resourceGo, nil
}

// GenerateList renders the list resource source for a resolved ResourceIR. It
// requires the resource to be listable (see ir.ResourceIR.IsListable).
func GenerateList(res *ir.ResourceIR) (string, error) {
	if !res.IsListable {
		return "", fmt.Errorf("resource %q has no subscription/resource-group list operations", res.TerraformType)
	}
	return RenderListResource(res), nil
}

// GenerateListTest renders the acceptance test for a resolved ResourceIR's list
// resource. It requires the resource to be listable (see ir.ResourceIR.IsListable).
func GenerateListTest(res *ir.ResourceIR) (string, error) {
	if !res.IsListable {
		return "", fmt.Errorf("resource %q has no subscription/resource-group list operations", res.TerraformType)
	}
	return RenderListTest(res), nil
}

// GenerateDataSource renders the data source source for a resolved ResourceIR.
func GenerateDataSource(res *ir.ResourceIR) (string, error) {
	return RenderDataSource(res), nil
}

func execTemplate(name, tpl string, data fileData) (string, error) {
	t, err := template.New(name).Parse(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// fileHeader is the license header prepended to every generated Go file so the
// output matches the repository's copyright convention.
const fileHeader = "// Copyright IBM Corp. 2014, 2025\n" +
	"// SPDX-License-Identifier: MPL-2.0\n\n"

const resourceTemplate = fileHeader + `package {{.Package}}

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"{{.SDKImportPath}}"
	"github.com/{{.ProviderGithubOrg}}/terraform-provider-{{.ProviderName}}/internal/sdk"
	"github.com/{{.ProviderGithubOrg}}/terraform-provider-{{.ProviderName}}/internal/tf/pluginsdk"
	"github.com/{{.ProviderGithubOrg}}/terraform-provider-{{.ProviderName}}/internal/tf/validation"
)

var _ sdk.Resource{{if .Updatable}}WithUpdate{{end}} = {{.Name}}Resource{}

type {{.Name}}Resource struct{}

{{.ModelStructs}}

func (r {{.Name}}Resource) ModelObject() interface{} {
	return &{{.ModelStructName}}{}
}

func (r {{.Name}}Resource) ResourceType() string {
	return "{{.TerraformType}}"
}

func (r {{.Name}}Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return {{.SDKPackage}}.{{.IDValidateFunc}}
}

func (r {{.Name}}Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
{{.Arguments}}
	}
}

func (r {{.Name}}Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
{{.Attributes}}
	}
}

func (r {{.Name}}Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.{{.ServiceName}}.{{.ClientField}}
{{if .HasSubscription}}			subscriptionId := metadata.Client.Account.SubscriptionId
{{end}}
			var config {{.ModelStructName}}
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := {{.SDKPackage}}.{{.IDNewFunc}}({{.ArgsJoined}})

			existing, err := client.{{.ReadOp}}(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			{{.CreatePayload}}
{{if .CreateLRO}}
			if err := client.{{.CreateOp}}ThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
{{else}}
			if _, err := client.{{.CreateOp}}(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
{{end}}
			metadata.SetID(id)
			return nil
		},
	}
}

func (r {{.Name}}Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.{{.ServiceName}}.{{.ClientField}}

			id, err := {{.SDKPackage}}.{{.IDParseFunc}}(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.{{.ReadOp}}(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			return r.flatten(metadata, id, model)
		},
	}
}
{{if .Updatable}}
func (r {{.Name}}Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			{{.UpdateBody}}
		},
	}
}
{{end}}
func (r {{.Name}}Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.{{.ServiceName}}.{{.ClientField}}

			id, err := {{.SDKPackage}}.{{.IDParseFunc}}(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
{{if .DeleteLRO}}
			if err := client.{{.DeleteOp}}ThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
{{else}}
			if _, err := client.{{.DeleteOp}}(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
{{end}}
			return nil
		},
	}
}

{{.FlattenMethod}}

{{.ExpandFlatten}}`
