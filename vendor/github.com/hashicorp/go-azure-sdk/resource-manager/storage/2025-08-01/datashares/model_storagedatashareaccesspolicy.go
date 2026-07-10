package datashares

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageDataShareAccessPolicy struct {
	Permission  StorageDataShareAccessPolicyPermission `json:"permission"`
	PrincipalId string                                 `json:"principalId"`
	TenantId    string                                 `json:"tenantId"`
}
