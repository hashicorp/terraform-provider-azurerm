package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendContractProperties struct {
	CircuitBreaker *BackendCircuitBreaker      `json:"circuitBreaker,omitempty"`
	Credentials    *BackendCredentialsContract `json:"credentials,omitempty"`
	Description    *string                     `json:"description,omitempty"`
	Pool           *BackendBaseParametersPool  `json:"pool,omitempty"`
	Properties     *BackendProperties          `json:"properties,omitempty"`
	Protocol       *BackendProtocol            `json:"protocol,omitempty"`
	Proxy          *BackendProxyContract       `json:"proxy,omitempty"`
	ResourceId     *string                     `json:"resourceId,omitempty"`
	Title          *string                     `json:"title,omitempty"`
	Tls            *BackendTlsProperties       `json:"tls,omitempty"`
	Type           *BackendType                `json:"type,omitempty"`
	Url            *string                     `json:"url,omitempty"`
}
