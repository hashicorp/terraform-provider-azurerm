package automations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationSource struct {
	EventSource *EventSource         `json:"eventSource,omitempty"`
	RuleSets    *[]AutomationRuleSet `json:"ruleSets,omitempty"`
}
