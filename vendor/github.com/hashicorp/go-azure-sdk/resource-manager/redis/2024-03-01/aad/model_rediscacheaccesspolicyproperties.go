package aad

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RedisCacheAccessPolicyProperties struct {
	Permissions       string                         `json:"permissions"`
	ProvisioningState *AccessPolicyProvisioningState `json:"provisioningState,omitempty"`
	Type              *AccessPolicyType              `json:"type,omitempty"`
}
