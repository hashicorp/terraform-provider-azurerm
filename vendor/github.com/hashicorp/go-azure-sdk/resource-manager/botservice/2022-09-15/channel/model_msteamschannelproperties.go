package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MsTeamsChannelProperties struct {
	AcceptedTerms         *bool   `json:"acceptedTerms,omitempty"`
	CallingWebhook        *string `json:"callingWebhook,omitempty"`
	DeploymentEnvironment *string `json:"deploymentEnvironment,omitempty"`
	EnableCalling         *bool   `json:"enableCalling,omitempty"`
	IncomingCallRoute     *string `json:"incomingCallRoute,omitempty"`
	IsEnabled             bool    `json:"isEnabled"`
}
