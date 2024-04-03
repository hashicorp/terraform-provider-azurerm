package alertsmanagements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsSummaryGroup struct {
	Groupedby        *string                   `json:"groupedby,omitempty"`
	SmartGroupsCount *int64                    `json:"smartGroupsCount,omitempty"`
	Total            *int64                    `json:"total,omitempty"`
	Values           *[]AlertsSummaryGroupItem `json:"values,omitempty"`
}
