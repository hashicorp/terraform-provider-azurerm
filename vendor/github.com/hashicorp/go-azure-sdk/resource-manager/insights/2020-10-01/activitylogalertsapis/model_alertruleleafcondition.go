package activitylogalertsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleLeafCondition struct {
	ContainsAny *[]string `json:"containsAny,omitempty"`
	Equals      *string   `json:"equals,omitempty"`
	Field       *string   `json:"field,omitempty"`
}
