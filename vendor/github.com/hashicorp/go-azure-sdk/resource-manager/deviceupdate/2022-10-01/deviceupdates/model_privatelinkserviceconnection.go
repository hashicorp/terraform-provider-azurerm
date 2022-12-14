package deviceupdates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceConnection struct {
	GroupIds       *[]string `json:"groupIds,omitempty"`
	Name           *string   `json:"name,omitempty"`
	RequestMessage *string   `json:"requestMessage,omitempty"`
}
