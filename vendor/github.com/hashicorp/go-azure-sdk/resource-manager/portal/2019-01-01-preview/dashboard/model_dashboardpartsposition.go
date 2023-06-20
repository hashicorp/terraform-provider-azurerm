package dashboard

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DashboardPartsPosition struct {
	ColSpan  int64                   `json:"colSpan"`
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
	RowSpan  int64                   `json:"rowSpan"`
	X        int64                   `json:"x"`
	Y        int64                   `json:"y"`
}
