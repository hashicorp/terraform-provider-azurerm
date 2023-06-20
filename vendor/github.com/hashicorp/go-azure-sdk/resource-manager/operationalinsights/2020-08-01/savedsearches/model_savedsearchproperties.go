package savedsearches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SavedSearchProperties struct {
	Category           string  `json:"category"`
	DisplayName        string  `json:"displayName"`
	FunctionAlias      *string `json:"functionAlias,omitempty"`
	FunctionParameters *string `json:"functionParameters,omitempty"`
	Query              string  `json:"query"`
	Tags               *[]Tag  `json:"tags,omitempty"`
	Version            *int64  `json:"version,omitempty"`
}
