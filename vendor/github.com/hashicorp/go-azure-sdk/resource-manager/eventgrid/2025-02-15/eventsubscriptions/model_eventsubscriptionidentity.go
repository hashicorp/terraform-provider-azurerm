package eventsubscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventSubscriptionIdentity struct {
	Type                 *EventSubscriptionIdentityType `json:"type,omitempty"`
	UserAssignedIdentity *string                        `json:"userAssignedIdentity,omitempty"`
}
