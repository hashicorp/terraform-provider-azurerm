package tagrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogRules struct {
	FilteringTags        *[]FilteringTag             `json:"filteringTags,omitempty"`
	SendAadLogs          *SendAadLogsStatus          `json:"sendAadLogs,omitempty"`
	SendActivityLogs     *SendActivityLogsStatus     `json:"sendActivityLogs,omitempty"`
	SendSubscriptionLogs *SendSubscriptionLogsStatus `json:"sendSubscriptionLogs,omitempty"`
}
