package budgets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Notification struct {
	ContactEmails []string       `json:"contactEmails"`
	ContactGroups *[]string      `json:"contactGroups,omitempty"`
	ContactRoles  *[]string      `json:"contactRoles,omitempty"`
	Enabled       bool           `json:"enabled"`
	Locale        *CultureCode   `json:"locale,omitempty"`
	Operator      OperatorType   `json:"operator"`
	Threshold     float64        `json:"threshold"`
	ThresholdType *ThresholdType `json:"thresholdType,omitempty"`
}
