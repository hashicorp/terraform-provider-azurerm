package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuggestRequest struct {
	Filter           *string  `json:"filter,omitempty"`
	Fuzzy            *bool    `json:"fuzzy,omitempty"`
	HighlightPostTag *string  `json:"highlightPostTag,omitempty"`
	HighlightPreTag  *string  `json:"highlightPreTag,omitempty"`
	MinimumCoverage  *float64 `json:"minimumCoverage,omitempty"`
	Orderby          *string  `json:"orderby,omitempty"`
	Search           string   `json:"search"`
	SearchFields     *string  `json:"searchFields,omitempty"`
	Select           *string  `json:"select,omitempty"`
	SuggesterName    string   `json:"suggesterName"`
	Top              *int64   `json:"top,omitempty"`
}
