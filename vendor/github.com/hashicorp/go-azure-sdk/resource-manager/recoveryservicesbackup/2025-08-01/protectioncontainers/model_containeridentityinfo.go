package protectioncontainers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerIdentityInfo struct {
	AadTenantId              *string `json:"aadTenantId,omitempty"`
	Audience                 *string `json:"audience,omitempty"`
	ServicePrincipalClientId *string `json:"servicePrincipalClientId,omitempty"`
	UniqueName               *string `json:"uniqueName,omitempty"`
}
