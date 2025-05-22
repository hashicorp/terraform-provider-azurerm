package data

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/util"
	"github.com/spf13/afero"
)

var documentFilePattern = "%s/website/docs/%s/%s.html.markdown"

func expectedDocumentationPath(providerDir string, resourceName string, resourceType ResourceType) string {
	return fmt.Sprintf(documentFilePattern, providerDir, ResourceTypeToDocumentationSubPath[resourceType], resourceName)
}

func getDocumentData(fs afero.Fs, data *ResourceData) {
	data.Document.Exists = util.FileExists(fs, data.Document.Path)

	if data.Document.Exists {
		if err := data.Document.Parse(fs); err != nil {
			data.Errors = append(data.Errors, fmt.Errorf("failed to parse documentation: %+v", err))
		}
	}
}
