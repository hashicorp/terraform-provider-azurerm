package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VersionProfile struct {
	ChannelGroup *string `json:"channelGroup,omitempty"`
	Id           string  `json:"id"`
}
