package nodepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodePoolVersionProfile struct {
	ChannelGroup *string `json:"channelGroup,omitempty"`
	Id           string  `json:"id"`
}
