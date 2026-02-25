package rule

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

type S005 struct{}

var _ Rule = S005{}

const expectedAttributesHeadingText = "Attributes Reference"

func (r S005) ID() string {
	return "S005"
}

func (r S005) Name() string {
	return "Attributes Heading"
}

func (r S005) Description() string {
	return "validates that the attributes section heading of the documentation page is correct"
}

func (r S005) Run(d *data.TerraformNodeData, fix bool) []error {
	if SkipRule(d.Type, d.Name, r.ID()) {
		return nil
	}

	if !d.Document.Exists {
		return nil
	}

	errs := make([]error, 0)

	var section *markdown.AttributesSection
	for _, sec := range d.Document.Sections {
		if sec, ok := sec.(*markdown.AttributesSection); ok {
			section = sec
			break
		}
	}

	if section == nil {
		return append(errs, fmt.Errorf("%s: missing Attributes section - please add this manually", IdAndName(r)))

		// TODO: implement templating of new section when validating/templating properties has been added.
	}

	content := section.GetContent()
	heading := section.GetHeading()

	headingUpdated := false
	if heading.Level != 2 {
		errs = append(errs, fmt.Errorf("%s: expected heading level to be 2, got %d", IdAndName(r), heading.Level))

		if fix {
			heading.Level = 2
			section.SetHeading(heading.String())
			headingUpdated = true
		}
	}

	if heading.Text != expectedAttributesHeadingText {
		errs = append(errs, fmt.Errorf("%s: expected heading text to be `%s`, got `%s`", IdAndName(r), expectedAttributesHeadingText, heading.Text))

		if fix {
			heading.Text = expectedAttributesHeadingText
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
