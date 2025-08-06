// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/md"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/model"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-lint/util"
)

// logic to load schema and markdown to print the diff

type ResourceDiff struct {
	tf *schema.Resource
	md *model.ResourceDoc

	SchemaFile string
	MDFile     string

	Diff []Checker // diff of exist in both md and code
}

func (d *ResourceDiff) Diffs() []Checker {
	return d.Diff
}

func (d *ResourceDiff) ToString() string {
	var bs strings.Builder

	bs.WriteString(
		fmt.Sprintf("%s: %s:1 has %d issue[s]:\n",
			util.Bold(d.tf.ResourceType),
			d.tf.FilePathRel(),
			len(d.Diff),
		),
	)
	file := d.MDFile + ":"
	if idx := strings.Index(file, "website"); idx > 0 {
		file = "./" + file[idx:]
	}
	// read file lines
	fileBuf, _ := os.ReadFile(d.MDFile)
	lines := strings.Split(string(fileBuf), "\n")
	for _, item := range d.Diffs() {
		bs.WriteString(file + item.String() + "\n")
		// print out fixed result
		if lineNum := item.Line(); lineNum > 0 && len(lines) > lineNum {
			line := lines[lineNum]
			if fixed, err := item.Fix(line); err == nil && fixed != "" && fixed != line {
				bs.WriteString("     " + util.IssueLine(line) + "\n")
				bs.WriteString("  => " + util.FixedCode(fixed) + "\n")
			}
		}
	}
	return bs.String()
}

// NewResourceDiff tf is required,
// mdPath is optional, it can be detected by resource name
func NewResourceDiff(tf *schema.Resource) *ResourceDiff {
	r := &ResourceDiff{
		tf: tf,
	}
	// try to detect Markdown path from resource
	// can set it if not a regular MD path
	r.MDFile = md.MDPathFor(tf.ResourceType)
	return r
}

func (r *ResourceDiff) DiffAll() {
	if r.md == nil {
		if r.MDFile == "" {
			r.Diff = append(r.Diff, newDiffWithMessage(fmt.Sprintf("%s has no document", r.tf.ResourceType), r.tf.IsDeprecated()))
			return
		}
		mark := md.MustNewMarkFromFile(r.MDFile)
		r.md = mark.BuildResourceDoc()
	}

	if name := r.md.HasCircularRef(); name != "" {
		r.Diff = append(r.Diff, newCircularRef(name, r.md))
		return
	}

	r.Diff = checkPossibleValues(r.tf, r.md)

	missDiff := crossCheckProperty(r.tf, r.md)
	r.Diff = append(r.Diff, missDiff...)

	timeouts := diffTimeout(r.tf, r.md)
	r.Diff = append(r.Diff, timeouts...)
}
