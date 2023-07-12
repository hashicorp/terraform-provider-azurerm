package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomationRunbookReceiver struct {
	AutomationAccountId  string  `json:"automationAccountId"`
	IsGlobalRunbook      bool    `json:"isGlobalRunbook"`
	Name                 *string `json:"name,omitempty"`
	RunbookName          string  `json:"runbookName"`
	ServiceUri           *string `json:"serviceUri,omitempty"`
	UseCommonAlertSchema *bool   `json:"useCommonAlertSchema,omitempty"`
	WebhookResourceId    string  `json:"webhookResourceId"`
}
