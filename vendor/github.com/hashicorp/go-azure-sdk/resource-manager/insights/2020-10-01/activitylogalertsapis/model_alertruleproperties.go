package activitylogalertsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleProperties struct {
	Actions     ActionList              `json:"actions"`
	Condition   AlertRuleAllOfCondition `json:"condition"`
	Description *string                 `json:"description,omitempty"`
	Enabled     *bool                   `json:"enabled,omitempty"`
	Scopes      []string                `json:"scopes"`
}
