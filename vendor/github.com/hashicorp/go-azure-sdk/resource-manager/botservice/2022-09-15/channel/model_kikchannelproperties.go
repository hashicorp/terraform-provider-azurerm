package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KikChannelProperties struct {
	ApiKey      *string `json:"apiKey,omitempty"`
	IsEnabled   bool    `json:"isEnabled"`
	IsValidated *bool   `json:"isValidated,omitempty"`
	UserName    string  `json:"userName"`
}
