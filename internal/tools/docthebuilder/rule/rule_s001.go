package rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/markdown"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/template"
	log "github.com/sirupsen/logrus"
)

type S001 struct{}

var _ Rule = S001{}

func (r S001) Name() string {
	return "S001"
}

func (r S001) Description() string {
	return fmt.Sprintf("%s - validates the `API Versions` section", r.Name())
}

func (r S001) Run(d *data.ResourceData, fix bool) []error {
	if !d.Document.Exists {
		return nil
	}

	exists := false
	errs := make([]error, 0)

	if len(d.APIs) == 0 {
		log.WithFields(log.Fields{
			"rule": r.Name(),
			"type": d.Type,
			"name": d.Name,
		}).Debug("resource object contained no APIs, skipping...")
		return nil
	}

	var section *markdown.APISection
	for _, sec := range d.Document.Sections {
		if sec, ok := sec.(*markdown.APISection); ok {
			section, exists = sec, true
		}
	}

	if section == nil {
		section = &markdown.APISection{}
	}

	expected, err := template.Render(d, section.Template())
	if err != nil {
		// This should probably log an error right away, runtime error not validation error
		errs = append(errs, fmt.Errorf("%s - failed to render template: %+v", r.Name(), err))
		return errs
	}

	if !exists {
		errs = append(errs, fmt.Errorf("%s - missing API section", r.Name()))

		if !fix {
			return errs
		}

		section.SetContent(expected)
		d.Document.Sections = append(d.Document.Sections, section)
		d.Document.HasChange = true
	} else {
		expectedStr := strings.Join(expected, "\n")
		currentStr := strings.Join(section.GetContent(), "\n")

		if currentStr != expectedStr {
			errs = append(errs, fmt.Errorf("%s - current section content did not match expected content", r.Name()))

			if fix {
				section.SetContent(expected)
				d.Document.HasChange = true
			}
		}
	}

	return errs
}
