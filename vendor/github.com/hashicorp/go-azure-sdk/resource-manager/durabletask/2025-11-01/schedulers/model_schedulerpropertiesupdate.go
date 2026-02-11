package schedulers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SchedulerPropertiesUpdate struct {
	Endpoint          *string             `json:"endpoint,omitempty"`
	IPAllowlist       *[]string           `json:"ipAllowlist,omitempty"`
	ProvisioningState *ProvisioningState  `json:"provisioningState,omitempty"`
	Sku               *SchedulerSkuUpdate `json:"sku,omitempty"`
}
