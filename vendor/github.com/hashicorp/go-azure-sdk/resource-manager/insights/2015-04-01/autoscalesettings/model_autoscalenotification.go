package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleNotification struct {
	Email     *EmailNotification     `json:"email,omitempty"`
	Operation OperationType          `json:"operation"`
	WebHooks  *[]WebhookNotification `json:"webhooks,omitempty"`
}
