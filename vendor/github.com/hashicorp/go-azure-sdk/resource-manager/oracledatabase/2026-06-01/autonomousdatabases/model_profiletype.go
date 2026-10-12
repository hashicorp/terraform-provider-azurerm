package autonomousdatabases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProfileType struct {
	ConsumerGroup     *ConsumerGroup         `json:"consumerGroup,omitempty"`
	DisplayName       string                 `json:"displayName"`
	HostFormat        HostFormatType         `json:"hostFormat"`
	IsRegional        *bool                  `json:"isRegional,omitempty"`
	Protocol          ProtocolType           `json:"protocol"`
	SessionMode       SessionModeType        `json:"sessionMode"`
	SyntaxFormat      SyntaxFormatType       `json:"syntaxFormat"`
	TlsAuthentication *TlsAuthenticationType `json:"tlsAuthentication,omitempty"`
	Value             string                 `json:"value"`
}
