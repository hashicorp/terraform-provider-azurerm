package alertsmanagements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsSummary struct {
	Id         *string             `json:"id,omitempty"`
	Name       *string             `json:"name,omitempty"`
	Properties *AlertsSummaryGroup `json:"properties,omitempty"`
	Type       *string             `json:"type,omitempty"`
}
