package networkconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProperties struct {
	DomainJoinType              DomainJoinType     `json:"domainJoinType"`
	DomainName                  *string            `json:"domainName,omitempty"`
	DomainPassword              *string            `json:"domainPassword,omitempty"`
	DomainUsername              *string            `json:"domainUsername,omitempty"`
	HealthCheckStatus           *HealthCheckStatus `json:"healthCheckStatus,omitempty"`
	NetworkingResourceGroupName *string            `json:"networkingResourceGroupName,omitempty"`
	OrganizationUnit            *string            `json:"organizationUnit,omitempty"`
	ProvisioningState           *ProvisioningState `json:"provisioningState,omitempty"`
	SubnetId                    *string            `json:"subnetId,omitempty"`
}
