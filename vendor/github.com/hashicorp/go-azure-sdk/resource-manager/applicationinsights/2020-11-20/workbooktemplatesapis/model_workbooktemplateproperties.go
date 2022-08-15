package workbooktemplatesapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplateProperties struct {
	Author       *string                                        `json:"author,omitempty"`
	Galleries    []WorkbookTemplateGallery                      `json:"galleries"`
	Localized    *map[string][]WorkbookTemplateLocalizedGallery `json:"localized,omitempty"`
	Priority     *int64                                         `json:"priority,omitempty"`
	TemplateData interface{}                                    `json:"templateData"`
}
