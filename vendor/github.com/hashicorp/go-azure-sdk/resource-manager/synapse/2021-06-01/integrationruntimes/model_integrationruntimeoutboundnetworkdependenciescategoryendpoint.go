package integrationruntimes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationRuntimeOutboundNetworkDependenciesCategoryEndpoint struct {
	Category  *string                                                  `json:"category,omitempty"`
	Endpoints *[]IntegrationRuntimeOutboundNetworkDependenciesEndpoint `json:"endpoints,omitempty"`
}
