package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyEntry struct {
	ApplicationId *string     `json:"applicationId,omitempty"`
	ObjectId      string      `json:"objectId"`
	Permissions   Permissions `json:"permissions"`
	TenantId      string      `json:"tenantId"`
}
