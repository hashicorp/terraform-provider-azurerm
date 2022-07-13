package apps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppTemplate struct {
	Description     *string                 `json:"description,omitempty"`
	Industry        *string                 `json:"industry,omitempty"`
	Locations       *[]AppTemplateLocations `json:"locations,omitempty"`
	ManifestId      *string                 `json:"manifestId,omitempty"`
	ManifestVersion *string                 `json:"manifestVersion,omitempty"`
	Name            *string                 `json:"name,omitempty"`
	Order           *float64                `json:"order,omitempty"`
	Title           *string                 `json:"title,omitempty"`
}
