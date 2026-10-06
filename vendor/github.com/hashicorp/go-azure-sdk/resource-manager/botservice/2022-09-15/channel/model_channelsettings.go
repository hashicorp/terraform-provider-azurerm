package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ChannelSettings struct {
	BotIconURL            *string `json:"botIconUrl,omitempty"`
	BotId                 *string `json:"botId,omitempty"`
	ChannelDisplayName    *string `json:"channelDisplayName,omitempty"`
	ChannelId             *string `json:"channelId,omitempty"`
	DisableLocalAuth      *bool   `json:"disableLocalAuth,omitempty"`
	ExtensionKey1         *string `json:"extensionKey1,omitempty"`
	ExtensionKey2         *string `json:"extensionKey2,omitempty"`
	IsEnabled             *bool   `json:"isEnabled,omitempty"`
	RequireTermsAgreement *bool   `json:"requireTermsAgreement,omitempty"`
	Sites                 *[]Site `json:"sites,omitempty"`
}
