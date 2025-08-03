package rule

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/data"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/document-fmt/util"
	"github.com/spf13/afero"
)

type G001 struct{}

var _ Rule = G001{}

func (r G001) ID() string {
	return "G001"
}

func (r G001) Name() string {
	return "Files Exist"
}

func (r G001) Description() string {
	return "validates resource code and documentation exists at the expected paths"
}

func (r G001) Run(rd *data.TerraformNodeData, _ bool) []error {
	errs := make([]error, 0)

	if !rd.Document.Exists {
		// Some deprecated resources may no longer have a documentation page
		if rd.Resource.DeprecationMessage == "" {
			errs = append(errs, fmt.Errorf("%s: Documentation for `%s` not found at expected path (%s)", IdAndName(r), rd.Name, rd.Document.Path))
		}
	}

	if !util.FileExists(afero.NewOsFs(), rd.Path) {
		errs = append(errs, fmt.Errorf("%s: %s Code for `%s` not found at expected path (%s)", IdAndName(r), rd.Type, rd.Name, rd.Path))
	}

	return errs
}
