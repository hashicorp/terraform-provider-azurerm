package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlexaChannelProperties struct {
	AlexaSkillId       string  `json:"alexaSkillId"`
	IsEnabled          bool    `json:"isEnabled"`
	ServiceEndpointUri *string `json:"serviceEndpointUri,omitempty"`
	UrlFragment        *string `json:"urlFragment,omitempty"`
}
