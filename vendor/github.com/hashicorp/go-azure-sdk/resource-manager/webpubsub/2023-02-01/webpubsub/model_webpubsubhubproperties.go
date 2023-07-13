package webpubsub

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WebPubSubHubProperties struct {
	AnonymousConnectPolicy *string          `json:"anonymousConnectPolicy,omitempty"`
	EventHandlers          *[]EventHandler  `json:"eventHandlers,omitempty"`
	EventListeners         *[]EventListener `json:"eventListeners,omitempty"`
}
