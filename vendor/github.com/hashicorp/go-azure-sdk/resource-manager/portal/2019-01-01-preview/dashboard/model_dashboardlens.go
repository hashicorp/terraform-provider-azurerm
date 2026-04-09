package dashboard

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardLens struct {
	Metadata *map[string]interface{}   `json:"metadata,omitempty"`
	Order    int64                     `json:"order"`
	Parts    map[string]DashboardParts `json:"parts"`
}
