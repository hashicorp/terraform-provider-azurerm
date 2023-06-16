package actiongroupsapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookReceiver struct {
	IdentifierUri        *string `json:"identifierUri,omitempty"`
	Name                 string  `json:"name"`
	ObjectId             *string `json:"objectId,omitempty"`
	ServiceUri           string  `json:"serviceUri"`
	TenantId             *string `json:"tenantId,omitempty"`
	UseAadAuth           *bool   `json:"useAadAuth,omitempty"`
	UseCommonAlertSchema *bool   `json:"useCommonAlertSchema,omitempty"`
}
