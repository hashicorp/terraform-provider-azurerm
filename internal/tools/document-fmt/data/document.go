// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package data

import (
	"fmt"
)

var documentFilePattern = "%s/website/docs/%s/%s.html.markdown"

func expectedDocumentationPath(providerDir string, resourceName string, resourceType ResourceType) string {
	return fmt.Sprintf(documentFilePattern, providerDir, ResourceTypeToDocumentationSubPath[resourceType], resourceName)
}
