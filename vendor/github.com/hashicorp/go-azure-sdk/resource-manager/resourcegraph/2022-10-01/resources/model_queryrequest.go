package resources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryRequest struct {
	Facets           *[]FacetRequest      `json:"facets,omitempty"`
	ManagementGroups *[]string            `json:"managementGroups,omitempty"`
	Options          *QueryRequestOptions `json:"options,omitempty"`
	Query            string               `json:"query"`
	Subscriptions    *[]string            `json:"subscriptions,omitempty"`
}
