package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookProperties struct {
	Actions           []WebhookAction    `json:"actions"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	Scope             *string            `json:"scope,omitempty"`
	Status            *WebhookStatus     `json:"status,omitempty"`
}
