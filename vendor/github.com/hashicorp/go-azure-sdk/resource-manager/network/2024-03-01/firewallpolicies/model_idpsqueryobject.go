package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IDPSQueryObject struct {
	Filters        *[]FilterItems `json:"filters,omitempty"`
	OrderBy        *OrderBy       `json:"orderBy,omitempty"`
	ResultsPerPage *int64         `json:"resultsPerPage,omitempty"`
	Search         *string        `json:"search,omitempty"`
	Skip           *int64         `json:"skip,omitempty"`
}
