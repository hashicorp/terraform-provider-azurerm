package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainOwnershipValidationResult struct {
	CustomDomainsForTopicSpacesConfiguration *[]CustomDomainConfiguration `json:"customDomainsForTopicSpacesConfiguration,omitempty"`
	CustomDomainsForTopicsConfiguration      *[]CustomDomainConfiguration `json:"customDomainsForTopicsConfiguration,omitempty"`
}
