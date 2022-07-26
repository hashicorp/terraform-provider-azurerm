package applicationinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplateLocalizedGallery struct {
	Galleries    *[]WorkbookTemplateGallery `json:"galleries,omitempty"`
	TemplateData *interface{}               `json:"templateData,omitempty"`
}
