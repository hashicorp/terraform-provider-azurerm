package subscriptions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTenantResponse struct {
	TenantId   *string `json:"tenantId,omitempty"`
	TenantName *string `json:"tenantName,omitempty"`
}
