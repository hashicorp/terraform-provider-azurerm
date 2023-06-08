package quotas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaResourceProperties struct {
	Limit             *int64         `json:"limit,omitempty"`
	ProvisioningState *ResourceState `json:"provisioningState,omitempty"`
	Usage             *int64         `json:"usage,omitempty"`
}
