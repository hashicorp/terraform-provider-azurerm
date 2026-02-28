package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointAuthenticationUserAssignedManagedIdentity struct {
	ClientId string  `json:"clientId"`
	Scope    *string `json:"scope,omitempty"`
	TenantId string  `json:"tenantId"`
}
