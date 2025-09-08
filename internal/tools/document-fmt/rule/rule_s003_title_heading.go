package rule

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/template"
	log "github.com/sirupsen/logrus"
)

type S003 struct{}

var _ Rule = S003{}

func (r S003) ID() string {
	return "G003"
}

func (r S003) Name() string {
	return "Title Heading"
}

func (r S003) Description() string {
	return "validates that the title section heading of the documentation page is correct"
}

func (r S003) Run(d *data.TerraformNodeData, fix bool) []error {
	if SkipRule(d.Type, d.Name, r.ID()) {
		return nil
	}

	if !d.Document.Exists {
		return nil
	}

	errs := make([]error, 0)

	logWithFields := log.WithFields(log.Fields{
		"rule": IdAndName(r),
		"type": d.Type,
		"name": d.Name,
	})

	var section *markdown.TitleSection
	for _, sec := range d.Document.Sections {
		if sec, ok := sec.(*markdown.TitleSection); ok {
			section = sec
			break
		}
	}

	if section == nil {
		section = &markdown.TitleSection{}

		errs = append(errs, fmt.Errorf("%s: missing Title section", IdAndName(r)))
		if !fix {
			return errs
		}

		content, err := template.Render(d, section.Template())
		if err != nil {
			logWithFields.Errorf("%s: failed to render template: %+v", IdAndName(r), err)
		}

		d.Document.HasChange = true
		section.SetContent(content)

		sections, err := markdown.InsertAfterSection(section, d.Document.Sections, &markdown.FrontMatterSection{})
		if err != nil {
			logWithFields.Errorf("%s: failed to insert new templated section: %+v", IdAndName(r), err)
		}

		d.Document.Sections = sections
		return errs
	}

	content := section.GetContent()
	heading := section.GetHeading()

	headingUpdated := false
	if heading.Level != 1 {
		errs = append(errs, fmt.Errorf("%s: expected heading level to be 1, got %d", IdAndName(r), heading.Level))

		if fix {
			heading.Level = 1
			section.SetHeading(heading.String())
			headingUpdated = true
		}
	}

	if expected := r.headingText(d); heading.Text != expected {
		errs = append(errs, fmt.Errorf("%s: expected heading text to be `%s`, got `%s`", IdAndName(r), expected, heading.Text))

		if fix {
			heading.Text = expected
			section.SetHeading(heading.String())
			headingUpdated = true
		}
	}

	if headingUpdated {
		for idx, l := range content {
			if section.Match(l) {
				content[idx] = heading.String()
				d.Document.HasChange = true
			}
		}
	}

	return errs
}

func (r S003) headingText(d *data.TerraformNodeData) string {
	if d.Type == data.ResourceTypeData {
		return fmt.Sprintf("Data Source: %s", d.Name)
	}

	return d.Name
}
