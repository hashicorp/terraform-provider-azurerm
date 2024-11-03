package organizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseDataPartnerOrganizationPropertiesUpdate struct {
	OrganizationId         *string                          `json:"organizationId,omitempty"`
	OrganizationName       *string                          `json:"organizationName,omitempty"`
	SingleSignOnProperties *LiftrBaseSingleSignOnProperties `json:"singleSignOnProperties,omitempty"`
	WorkspaceId            *string                          `json:"workspaceId,omitempty"`
	WorkspaceName          *string                          `json:"workspaceName,omitempty"`
}
