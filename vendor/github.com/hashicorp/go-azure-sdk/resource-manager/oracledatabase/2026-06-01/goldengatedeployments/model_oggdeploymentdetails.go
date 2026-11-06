package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OggDeploymentDetails struct {
	AdminPassword       *string                     `json:"adminPassword,omitempty"`
	AdminUsername       *string                     `json:"adminUsername,omitempty"`
	Certificate         *string                     `json:"certificate,omitempty"`
	CredentialStore     *CredentialType             `json:"credentialStore,omitempty"`
	DeploymentName      string                      `json:"deploymentName"`
	GroupToRolesMapping *GroupToRolesMappingDetails `json:"groupToRolesMapping,omitempty"`
	OggVersion          *string                     `json:"oggVersion,omitempty"`
	PasswordSecretId    *string                     `json:"passwordSecretId,omitempty"`
}
