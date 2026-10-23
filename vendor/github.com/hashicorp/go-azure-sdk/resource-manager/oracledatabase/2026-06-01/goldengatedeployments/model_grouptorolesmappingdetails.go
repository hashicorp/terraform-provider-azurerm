package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupToRolesMappingDetails struct {
	AdministratorGroupId *string `json:"administratorGroupId,omitempty"`
	IdentityDomainId     *string `json:"identityDomainId,omitempty"`
	Key                  *string `json:"key,omitempty"`
	OperatorGroupId      *string `json:"operatorGroupId,omitempty"`
	SecurityGroupId      *string `json:"securityGroupId,omitempty"`
	UserGroupId          *string `json:"userGroupId,omitempty"`
}
