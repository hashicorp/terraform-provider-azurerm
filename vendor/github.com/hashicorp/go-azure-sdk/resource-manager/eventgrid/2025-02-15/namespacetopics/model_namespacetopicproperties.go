package namespacetopics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NamespaceTopicProperties struct {
	EventRetentionInDays *int64                           `json:"eventRetentionInDays,omitempty"`
	InputSchema          *EventInputSchema                `json:"inputSchema,omitempty"`
	ProvisioningState    *NamespaceTopicProvisioningState `json:"provisioningState,omitempty"`
	PublisherType        *PublisherType                   `json:"publisherType,omitempty"`
}
