package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LineRegistration struct {
	ChannelAccessToken *string `json:"channelAccessToken,omitempty"`
	ChannelSecret      *string `json:"channelSecret,omitempty"`
	GeneratedId        *string `json:"generatedId,omitempty"`
}
