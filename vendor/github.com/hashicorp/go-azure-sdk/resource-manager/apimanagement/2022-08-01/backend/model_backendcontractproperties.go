package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendContractProperties struct {
	Credentials *BackendCredentialsContract `json:"credentials,omitempty"`
	Description *string                     `json:"description,omitempty"`
	Properties  *BackendProperties          `json:"properties,omitempty"`
	Protocol    BackendProtocol             `json:"protocol"`
	Proxy       *BackendProxyContract       `json:"proxy,omitempty"`
	ResourceId  *string                     `json:"resourceId,omitempty"`
	Title       *string                     `json:"title,omitempty"`
	Tls         *BackendTlsProperties       `json:"tls,omitempty"`
	Url         string                      `json:"url"`
}
