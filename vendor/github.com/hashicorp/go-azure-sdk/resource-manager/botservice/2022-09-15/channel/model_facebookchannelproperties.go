package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FacebookChannelProperties struct {
	AppId       string          `json:"appId"`
	AppSecret   *string         `json:"appSecret,omitempty"`
	CallbackURL *string         `json:"callbackUrl,omitempty"`
	IsEnabled   bool            `json:"isEnabled"`
	Pages       *[]FacebookPage `json:"pages,omitempty"`
	VerifyToken *string         `json:"verifyToken,omitempty"`
}
