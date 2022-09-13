package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkableEnvironmentRequest struct {
	Region        *string `json:"region,omitempty"`
	TenantId      *string `json:"tenantId,omitempty"`
	UserPrincipal *string `json:"userPrincipal,omitempty"`
}
