package webhook

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebhookCreateOrUpdateParameters struct {
	Name       string                          `json:"name"`
	Properties WebhookCreateOrUpdateProperties `json:"properties"`
}
