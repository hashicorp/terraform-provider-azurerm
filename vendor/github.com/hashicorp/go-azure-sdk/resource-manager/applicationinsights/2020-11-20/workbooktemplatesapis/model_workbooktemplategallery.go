package workbooktemplatesapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookTemplateGallery struct {
	Category     *string `json:"category,omitempty"`
	Name         *string `json:"name,omitempty"`
	Order        *int64  `json:"order,omitempty"`
	ResourceType *string `json:"resourceType,omitempty"`
	Type         *string `json:"type,omitempty"`
}
