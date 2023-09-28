package resources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FacetRequestOptions struct {
	Filter    *string         `json:"filter,omitempty"`
	SortBy    *string         `json:"sortBy,omitempty"`
	SortOrder *FacetSortOrder `json:"sortOrder,omitempty"`
	Top       *int64          `json:"$top,omitempty"`
}
