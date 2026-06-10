package networkconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConnectionUpdateProperties struct {
	DomainName       *string `json:"domainName,omitempty"`
	DomainPassword   *string `json:"domainPassword,omitempty"`
	DomainUsername   *string `json:"domainUsername,omitempty"`
	OrganizationUnit *string `json:"organizationUnit,omitempty"`
	SubnetId         *string `json:"subnetId,omitempty"`
}
