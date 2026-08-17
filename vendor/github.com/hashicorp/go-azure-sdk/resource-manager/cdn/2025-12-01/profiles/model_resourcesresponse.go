package profiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourcesResponse struct {
	CustomDomains *[]ResourcesResponseCustomDomainsItem `json:"customDomains,omitempty"`
	Endpoints     *[]ResourcesResponseEndpointsItem     `json:"endpoints,omitempty"`
}
