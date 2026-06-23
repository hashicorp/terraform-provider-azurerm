package firewallstatusresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallStatusProperty struct {
	HealthReason           *string                    `json:"healthReason,omitempty"`
	HealthStatus           *HealthStatus              `json:"healthStatus,omitempty"`
	IsPanoramaManaged      *BooleanEnum               `json:"isPanoramaManaged,omitempty"`
	IsStrataCloudManaged   *BooleanEnum               `json:"isStrataCloudManaged,omitempty"`
	PanoramaStatus         *PanoramaStatus            `json:"panoramaStatus,omitempty"`
	ProvisioningState      *ReadOnlyProvisioningState `json:"provisioningState,omitempty"`
	StrataCloudManagerInfo *StrataCloudManagerInfo    `json:"strataCloudManagerInfo,omitempty"`
}
