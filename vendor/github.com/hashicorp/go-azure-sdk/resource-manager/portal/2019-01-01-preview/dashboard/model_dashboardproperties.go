package dashboard

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardProperties struct {
	Lenses   *map[string]DashboardLens `json:"lenses,omitempty"`
	Metadata *map[string]interface{}   `json:"metadata,omitempty"`
}
