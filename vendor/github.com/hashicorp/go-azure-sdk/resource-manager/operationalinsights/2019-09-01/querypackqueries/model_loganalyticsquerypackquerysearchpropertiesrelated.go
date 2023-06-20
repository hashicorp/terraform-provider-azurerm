package querypackqueries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsQueryPackQuerySearchPropertiesRelated struct {
	Categories    *[]string `json:"categories,omitempty"`
	ResourceTypes *[]string `json:"resourceTypes,omitempty"`
	Solutions     *[]string `json:"solutions,omitempty"`
}
