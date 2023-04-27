package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHubReceiver struct {
	EventHubName         string  `json:"eventHubName"`
	EventHubNameSpace    string  `json:"eventHubNameSpace"`
	Name                 string  `json:"name"`
	SubscriptionId       string  `json:"subscriptionId"`
	TenantId             *string `json:"tenantId,omitempty"`
	UseCommonAlertSchema *bool   `json:"useCommonAlertSchema,omitempty"`
}
