package rule

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/markdown"
)

type S004 struct{}

var _ Rule = S004{}

const expectedArgumentsHeadingText = "Arguments Reference"

func (r S004) ID() string {
	return "S004"
}

func (r S004) Name() string {
	return "Arguments Heading"
}

func (r S004) Description() string {
	return "validates that the arguments section heading of the documentation page is correct"
}

func (r S004) Run(d *data.TerraformNodeData, fix bool) []error {
	if SkipRule(d.Type, d.Name, r.ID()) {
		return nil
	}

	if !d.Document.Exists {
		return nil
	}

	errs := make([]error, 0)

	var section *markdown.ArgumentsSection
	for _, sec := range d.Document.Sections {
		if sec, ok := sec.(*markdown.ArgumentsSection); ok {
			section = sec
			break
		}
	}

	if section == nil {
		return append(errs, fmt.Errorf("%s: missing Arguments section - please add this manually", IdAndName(r)))

		// TODO: implement templating of new section when validating/templating of all properties has been added.
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

	if heading.Text != expectedArgumentsHeadingText {
		errs = append(errs, fmt.Errorf("%s: expected heading text to be `%s`, got `%s`", IdAndName(r), expectedArgumentsHeadingText, heading.Text))

		if fix {
			heading.Text = expectedArgumentsHeadingText
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
