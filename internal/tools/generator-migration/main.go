package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/generator-migration/snapshot"
)

var (
	//go:embed migration.gotpl
	tplStr string

	migrationTpl = template.Must(template.New("migration").Parse(tplStr))
)

type Model struct {
	TypeName     string
	ResourceName string
	FromVersion  int
	ToVersion    int
	Schema       string
	MigratorType string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: generator-migration <resource_type>")
	}
	resource := os.Args[1]
	res, ok := provider.AzureProvider().ResourcesMap[resource]
	if !ok {
		log.Fatalf("unknown resource type: %s", resource)
	}
	model := buildModel(resource, res)
	model.ResourceName = resource

	var buf bytes.Buffer
	if err := migrationTpl.Execute(&buf, &model); err != nil {
		log.Fatalf("render code err: %v", err)
	}
	code, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("format code err: %v\n\n%s", err, buf.String())
	}
	fmt.Printf("%s", string(code))
}

func buildModel(rt string, res *schema.Resource) Model {
	var m Model
	m.ResourceName = rt
	m.FromVersion = res.SchemaVersion
	m.ToVersion = res.SchemaVersion + 1

	toUp := true
	var buf []byte
	for _, ch := range strings.TrimPrefix(rt, "azurerm") {
		if ch == '_' {
			toUp = true
			continue
		}
		if toUp {
			ch -= 'a' - 'A'
		}
		buf = append(buf, byte(ch))
		toUp = false
	}
	m.TypeName = string(buf)
	m.MigratorType = fmt.Sprintf("%sV%dToV%d", m.TypeName, m.FromVersion, m.ToVersion)

	m.Schema = snapshot.GenSchemaSnapshot(res)
	return m
}
