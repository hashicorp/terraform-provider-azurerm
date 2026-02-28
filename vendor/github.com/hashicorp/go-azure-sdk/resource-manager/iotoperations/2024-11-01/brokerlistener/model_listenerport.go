package brokerlistener

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListenerPort struct {
	AuthenticationRef *string             `json:"authenticationRef,omitempty"`
	AuthorizationRef  *string             `json:"authorizationRef,omitempty"`
	NodePort          *int64              `json:"nodePort,omitempty"`
	Port              int64               `json:"port"`
	Protocol          *BrokerProtocolType `json:"protocol,omitempty"`
	Tls               *TlsCertMethod      `json:"tls,omitempty"`
}
