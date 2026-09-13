package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookPropertiesUpdateParameters struct {
	Actions       *[]WebhookAction   `json:"actions,omitempty"`
	CustomHeaders *map[string]string `json:"customHeaders,omitempty"`
	Scope         *string            `json:"scope,omitempty"`
	ServiceUri    *string            `json:"serviceUri,omitempty"`
	Status        *WebhookStatus     `json:"status,omitempty"`
}
