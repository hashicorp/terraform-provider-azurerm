package apiversionset

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiVersionSetContractProperties struct {
	Description       *string          `json:"description,omitempty"`
	DisplayName       string           `json:"displayName"`
	VersionHeaderName *string          `json:"versionHeaderName,omitempty"`
	VersionQueryName  *string          `json:"versionQueryName,omitempty"`
	VersioningScheme  VersioningScheme `json:"versioningScheme"`
}
