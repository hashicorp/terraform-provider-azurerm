package webpubsub

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventHandler struct {
	Auth             *UpstreamAuthSettings `json:"auth,omitempty"`
	SystemEvents     *[]string             `json:"systemEvents,omitempty"`
	UrlTemplate      string                `json:"urlTemplate"`
	UserEventPattern *string               `json:"userEventPattern,omitempty"`
}
