package webhook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookUpdateParameters struct {
	Name       *string                  `json:"name,omitempty"`
	Properties *WebhookUpdateProperties `json:"properties,omitempty"`
}
