package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerAttachResponse struct {
	Password     *string `json:"password,omitempty"`
	WebSocketUri *string `json:"webSocketUri,omitempty"`
}
