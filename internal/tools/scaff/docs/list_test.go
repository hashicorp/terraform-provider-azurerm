// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package docs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleListResource = `package example

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type ExampleThingListResource struct{}

func (ExampleThingListResource) Metadata(_ context.Context, _ resource.MetadataRequest, response *resource.MetadataResponse) {
	response.TypeName = "azurerm_example_thing"
}

func (ExampleThingListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, response *list.ListResourceSchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resource_group_name": schema.StringAttribute{
				Optional: true,
			},
			"subscription_id": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}
`

func TestGenerateListDoc(t *testing.T) {
	root := t.TempDir()
	svcDir := filepath.Join(root, "internal", "services", "example")
	if err := os.MkdirAll(svcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	listFile := filepath.Join(svcDir, "example_thing_resource_list.go")
	if err := os.WriteFile(listFile, []byte(sampleListResource), 0o644); err != nil {
		t.Fatal(err)
	}

	// First run writes the doc.
	outPath, written, err := GenerateListDoc(listFile, false)
	if err != nil {
		t.Fatalf("GenerateListDoc: %v", err)
	}
	if !written {
		t.Fatalf("expected doc to be written on first run")
	}
	wantPath := filepath.Join(root, "website", "docs", "list-resources", "example_thing.html.markdown")
	if outPath != wantPath {
		t.Fatalf("output path = %q, want %q", outPath, wantPath)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(got)
	for _, want := range []string{
		`page_title: "Azure Resource Manager: azurerm_example_thing"`,
		"# List resource: azurerm_example_thing",
		"### List all Example Things in the subscription",
		"### List all Example Things in a Resource Group",
		"resource_group_name = \"example-rg\"",
		"* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.",
		"* `resource_group_name` - (Optional) The name of the Resource Group to query.",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("generated doc missing %q\n---\n%s", want, content)
		}
	}

	// Second run without overwrite preserves the existing doc.
	if _, written, err = GenerateListDoc(listFile, false); err != nil {
		t.Fatalf("GenerateListDoc (second run): %v", err)
	}
	if written {
		t.Fatalf("expected existing doc to be preserved without overwrite")
	}

	// With overwrite the doc is rewritten.
	if _, written, err = GenerateListDoc(listFile, true); err != nil {
		t.Fatalf("GenerateListDoc (overwrite): %v", err)
	}
	if !written {
		t.Fatalf("expected doc to be rewritten with overwrite")
	}
}
