package namespaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopicsConfiguration struct {
	CustomDomains *[]CustomDomainConfiguration `json:"customDomains,omitempty"`
	Hostname      *string                      `json:"hostname,omitempty"`
}
