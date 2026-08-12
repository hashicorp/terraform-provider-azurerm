// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

// Package scaffold implements the schema-customization layer for scaff. A
// mapping file, generated alongside a resource, lets users rename or remove
// schema attributes and adjust their metadata; the customizations are applied
// to the resolved IR before code generation so the same edits survive a
// regeneration.
package scaffold

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

// Mapping is the schema-customization file generated alongside a resource. It
// records the resolution inputs (so `scaff regen` is self-contained) plus a set
// of per-attribute overrides keyed by each property's stable source path.
type Mapping struct {
	// Resolution inputs — mirror the `scaff generate` flags so regeneration can
	// re-resolve the resource from Pandora without any additional arguments.
	ResourceName    string `hcl:"resource_name"`
	ARMType         string `hcl:"arm_type,optional"`
	Service         string `hcl:"service,optional"`
	PandoraResource string `hcl:"pandora_resource,optional"`
	APIVersion      string `hcl:"api_version,optional"`
	GoName          string `hcl:"go_name,optional"`
	ServicePackage  string `hcl:"servicepackage,optional"`
	Provider        string `hcl:"provider,optional"`
	Org             string `hcl:"org,optional"`

	// Which artifacts regeneration should (re)produce.
	GenResource *bool `hcl:"gen_resource,optional"`
	List        *bool `hcl:"list,optional"`
	DataSource  *bool `hcl:"data_source,optional"`

	Attributes []AttributeOverride `hcl:"attribute,block"`
}

// AttributeOverride customizes a single schema property, addressed by its stable
// source path — the JSON path from the model root, e.g.
// "properties.networkProfile.podCidr". Pointer fields distinguish "unset" (leave
// the resolved value untouched) from an explicit override.
type AttributeOverride struct {
	SourcePath  string  `hcl:"source_path,label"`
	TFName      *string `hcl:"tf_name,optional"`
	Remove      *bool   `hcl:"remove,optional"`
	Required    *bool   `hcl:"required,optional"`
	Optional    *bool   `hcl:"optional,optional"`
	Computed    *bool   `hcl:"computed,optional"`
	Sensitive   *bool   `hcl:"sensitive,optional"`
	ForceNew    *bool   `hcl:"force_new,optional"`
	Description *string `hcl:"description,optional"`
}

// Parse reads and decodes a mapping file from disk.
func Parse(path string) (*Mapping, error) {
	var m Mapping
	if err := hclsimple.DecodeFile(path, nil, &m); err != nil {
		return nil, fmt.Errorf("reading mapping file %q: %w", path, err)
	}
	if m.ResourceName == "" {
		return nil, fmt.Errorf("mapping file %q: resource_name is required", path)
	}
	return &m, nil
}
