package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IngressSettings struct {
	BackendProtocol      *BackendProtocol           `json:"backendProtocol,omitempty"`
	ClientAuth           *IngressSettingsClientAuth `json:"clientAuth,omitempty"`
	ReadTimeoutInSeconds *int64                     `json:"readTimeoutInSeconds,omitempty"`
	SendTimeoutInSeconds *int64                     `json:"sendTimeoutInSeconds,omitempty"`
	SessionAffinity      *SessionAffinity           `json:"sessionAffinity,omitempty"`
	SessionCookieMaxAge  *int64                     `json:"sessionCookieMaxAge,omitempty"`
}
