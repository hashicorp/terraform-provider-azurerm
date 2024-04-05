package alertsmanagements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsSummaryGroupItem struct {
	Count     *int64                    `json:"count,omitempty"`
	Groupedby *string                   `json:"groupedby,omitempty"`
	Name      *string                   `json:"name,omitempty"`
	Values    *[]AlertsSummaryGroupItem `json:"values,omitempty"`
}
