package api

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiVersionSetContractDetails struct {
	Description       *string           `json:"description,omitempty"`
	Id                *string           `json:"id,omitempty"`
	Name              *string           `json:"name,omitempty"`
	VersionHeaderName *string           `json:"versionHeaderName,omitempty"`
	VersionQueryName  *string           `json:"versionQueryName,omitempty"`
	VersioningScheme  *VersioningScheme `json:"versioningScheme,omitempty"`
}
