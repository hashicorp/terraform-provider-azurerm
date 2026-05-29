package rules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogRules struct {
	FilteringTags        *[]FilteringTag `json:"filteringTags,omitempty"`
	SendAadLogs          *bool           `json:"sendAadLogs,omitempty"`
	SendResourceLogs     *bool           `json:"sendResourceLogs,omitempty"`
	SendSubscriptionLogs *bool           `json:"sendSubscriptionLogs,omitempty"`
}
