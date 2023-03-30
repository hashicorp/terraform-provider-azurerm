package liveevents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiveEventEndpoint struct {
	Protocol *string `json:"protocol,omitempty"`
	Url      *string `json:"url,omitempty"`
}
