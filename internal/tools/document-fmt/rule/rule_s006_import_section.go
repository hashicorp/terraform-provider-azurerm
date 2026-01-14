package rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/template"
	log "github.com/sirupsen/logrus"
)

type S006 struct{}

var _ Rule = S006{}

func (r S006) ID() string {
	return "S006"
}

func (r S006) Name() string {
	return "Import Section"
}

func (r S006) Description() string {
	return "validates the import section uses the new format if resource identity is supported"
}

func (r S006) Run(d *data.TerraformNodeData, fix bool) (errs []error) {
	if SkipRule(d.Type, d.Name, r.ID()) {
		return nil
	}

	if !d.Document.Exists || !(d.Type == data.ResourceTypeResource) {
		return nil
	}

	var section *markdown.ImportSection
	var exists bool
	for _, sec := range d.Document.Sections {
		if sec, ok := sec.(*markdown.ImportSection); ok {
			section, exists = sec, true
			break
		}
	}

	if section == nil {
		section = &markdown.ImportSection{}
	}

	// If ResourceIdentity is not supported, return
	// TODO: add validation for import sections that _don't_ support resource identity or are untyped
	// We can probably do so by parsing existing import section and using the import ID string using recaser to find the right Resource ID
	if d.ResourceID == nil {
		return
	}

	logWithFields := log.WithFields(log.Fields{
		"rule": IdAndName(r),
		"type": d.Type,
		"name": d.Name,
	})

	expected, err := template.Render(d, section.Template())
	if err != nil {
		logWithFields.WithError(err).Error("failed to render template: %+v", err)
		return
	}

	if !exists {
		errs = append(errs, fmt.Errorf("%s: missing Import section", IdAndName(r)))

		if fix {
			section.SetContent(expected)
			sections, err := markdown.InsertAfterSection(section, d.Document.Sections, &markdown.TimeoutsSection{})
			if err != nil {
				logWithFields.WithError(err).Error("failed to insert new templated section: %+v", err)
			}
			d.Document.Sections = sections
			d.Document.HasChange = true
		}

		return
	}

	if strings.Join(expected, "\n") != strings.Join(section.GetContent(), "\n") {
		errs = append(errs, fmt.Errorf("%s: current section content did not match expected content", IdAndName(r)))

		if fix {
			section.SetContent(expected)
			d.Document.HasChange = true
		}
	}

	return
}

func (r S006) ResourceIDToIdentitySchema(d data.TerraformNodeData) []string {
	result := make([]string, 0)

	segments := d.ResourceID.Segments()
	numSegments := len(segments)
	for idx, segment := range segments {
		if pluginsdk.SegmentTypeSupported(segment.Type) {
			name := pluginsdk.SegmentName(segment, d.ResourceIdentityType, numSegments, idx)
			result = append(result, name)
		}
	}

	return result
}
