// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package commands

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/gen"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/ir"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/pandora"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/scaff/scaffold"

	"github.com/mitchellh/cli"
)

type RegenCommand struct {
	Ui cli.Ui
}

var _ cli.Command = RegenCommand{}

type regenData struct {
	File       string
	PandoraURL string
	Provider   string
	Org        string
	OutputPath string
	Overwrite  bool
	Config     *helpers.Configuration
}

func (d *regenData) parseArgs(args []string) error {
	set := flag.NewFlagSet("regen", flag.ExitOnError)
	set.StringVar(&d.File, "file", "", "(Required) path to a .scaffold.hcl schema mapping file")
	set.StringVar(&d.PandoraURL, "pandora-url", "", "(Optional) Pandora Data API base URL; defaults to the configured schema_api_url")
	set.StringVar(&d.OutputPath, "path", "", "(Optional) output directory; defaults to the mapping file's directory")
	set.StringVar(&d.Provider, "provider", "", "(Optional) provider name override, e.g. \"azurerm\"")
	set.StringVar(&d.Org, "org", "", "(Optional) provider GitHub org override, e.g. \"hashicorp\"")
	set.BoolVar(&d.Overwrite, "overwrite", false, "(Required) confirm rewriting the generated files")
	if err := set.Parse(args); err != nil {
		return err
	}
	if d.File == "" {
		return errors.New("required option -file missing (path to a .scaffold.hcl mapping)")
	}
	return nil
}

func (c RegenCommand) Run(args []string) int {
	d := &regenData{Config: config}
	if err := d.parseArgs(args); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	if err := d.regen(c.Ui); err != nil {
		c.Ui.Error(fmt.Sprintf("regenerating resource: %+v", err))
		return 1
	}
	return 0
}

// regen re-resolves the resource described by the mapping file from Pandora,
// applies the mapping's customizations to the schema, and rewrites the
// generated files.
func (d regenData) regen(ui cli.Ui) error {
	if !d.Overwrite {
		return errors.New("regen rewrites the generated files; pass -overwrite to confirm")
	}

	m, err := scaffold.Parse(d.File)
	if err != nil {
		return err
	}

	pandoraURL := d.PandoraURL
	if pandoraURL == "" && d.Config != nil {
		pandoraURL = d.Config.SchemaAPIURL
	}

	provider, org := "", ""
	if d.Config != nil {
		provider, org = d.Config.ProviderName, d.Config.ProviderGithubOrg
	}
	if m.Provider != "" {
		provider = m.Provider
	}
	if m.Org != "" {
		org = m.Org
	}
	if d.Provider != "" {
		provider = d.Provider
	}
	if d.Org != "" {
		org = d.Org
	}

	client := pandora.NewClient(pandoraURL)
	res, err := ir.Resolve(client, ir.Options{
		ARMType:           m.ARMType,
		Service:           m.Service,
		Resource:          m.PandoraResource,
		APIVersion:        m.APIVersion,
		Name:              m.ResourceName,
		GoName:            m.GoName,
		ServicePackage:    m.ServicePackage,
		ProviderName:      provider,
		ProviderGithubOrg: org,
	})
	if err != nil {
		return err
	}

	warnings, err := scaffold.Apply(res, m)
	for _, w := range warnings {
		ui.Warn(w)
	}
	if err != nil {
		return err
	}

	outputDir := d.OutputPath
	if outputDir == "" {
		outputDir = filepath.Dir(d.File)
	}

	// Reuse the generate command's write / format / registration helpers.
	g := GenerateData{Overwrite: true, Config: d.Config}

	if m.GenResource == nil || *m.GenResource {
		resourceGo, err := gen.Generate(res)
		if err != nil {
			return err
		}
		file := filepath.Join(outputDir, fmt.Sprintf("%s_resource.go", m.ResourceName))
		if err := g.writeGenerated(ui, file, resourceGo); err != nil {
			return err
		}
	}

	if m.List != nil && *m.List {
		listGo, err := gen.GenerateList(res)
		if err != nil {
			return err
		}
		listFile := filepath.Join(outputDir, fmt.Sprintf("%s_resource_list.go", m.ResourceName))
		if err := g.writeGenerated(ui, listFile, listGo); err != nil {
			return err
		}

		listTestGo, err := gen.GenerateListTest(res)
		if err != nil {
			return err
		}
		listTestFile := filepath.Join(outputDir, fmt.Sprintf("%s_resource_list_test.go", m.ResourceName))
		if err := g.writeGenerated(ui, listTestFile, listTestGo); err != nil {
			return err
		}

		regPath := filepath.Join(outputDir, "registration.go")
		if err := g.registerListResource(ui, regPath, res.Name+"ListResource"); err != nil {
			return err
		}

		generateListDoc(ui, listFile, true)
	}

	if m.DataSource != nil && *m.DataSource {
		dsGo, err := gen.GenerateDataSource(res)
		if err != nil {
			return err
		}
		dsFile := filepath.Join(outputDir, fmt.Sprintf("%s_data_source.go", m.ResourceName))
		if err := g.writeGenerated(ui, dsFile, dsGo); err != nil {
			return err
		}
	}

	return nil
}

func (c RegenCommand) Help() string {
	return `
Usage: scaff regen -file <path/to/name.scaffold.hcl> -overwrite

Regenerates a resource (and its data source / list, per the mapping) from a
.scaffold.hcl schema-customization file. The resource is re-resolved from the
Pandora Data API using the inputs recorded in the mapping, the mapping's
attribute customizations (rename / remove / metadata) are applied to the schema,
and the generated files are rewritten in place.

Because it rewrites files, -overwrite is required.

Example:
  $ scaff regen -file internal/services/redhatopenshift/redhat_openshift_cluster.scaffold.hcl -overwrite
`
}

func (c RegenCommand) Synopsis() string {
	return "Regenerates a resource from its .scaffold.hcl schema mapping."
}
