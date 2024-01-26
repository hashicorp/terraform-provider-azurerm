package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EmailChannelProperties struct {
	AuthMethod   *EmailChannelAuthMethod `json:"authMethod,omitempty"`
	EmailAddress string                  `json:"emailAddress"`
	IsEnabled    bool                    `json:"isEnabled"`
	MagicCode    *string                 `json:"magicCode,omitempty"`
	Password     *string                 `json:"password,omitempty"`
}
