package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Source struct {
	AuthorizedResources *[]string  `json:"authorizedResources,omitempty"`
	DataSourceId        string     `json:"dataSourceId"`
	Query               *string    `json:"query,omitempty"`
	QueryType           *QueryType `json:"queryType,omitempty"`
}
