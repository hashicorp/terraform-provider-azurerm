package outboundnetworkdependenciesendpoints

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FQDNEndpoint struct {
	DomainName      *string               `json:"domainName,omitempty"`
	EndpointDetails *[]FQDNEndpointDetail `json:"endpointDetails,omitempty"`
}
