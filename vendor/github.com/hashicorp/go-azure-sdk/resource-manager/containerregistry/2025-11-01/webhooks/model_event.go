package webhooks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Event struct {
	EventRequestMessage  *EventRequestMessage  `json:"eventRequestMessage,omitempty"`
	EventResponseMessage *EventResponseMessage `json:"eventResponseMessage,omitempty"`
	Id                   *string               `json:"id,omitempty"`
}
