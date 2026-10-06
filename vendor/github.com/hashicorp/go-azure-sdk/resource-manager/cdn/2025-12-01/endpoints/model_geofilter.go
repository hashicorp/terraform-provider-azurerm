package endpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GeoFilter struct {
	Action       GeoFilterActions `json:"action"`
	CountryCodes []string         `json:"countryCodes"`
	RelativePath string           `json:"relativePath"`
}
