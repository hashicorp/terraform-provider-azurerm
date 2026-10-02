package managedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MHSMGeoReplicatedRegion struct {
	IsPrimary         *bool                                  `json:"isPrimary,omitempty"`
	Name              *string                                `json:"name,omitempty"`
	ProvisioningState *GeoReplicationRegionProvisioningState `json:"provisioningState,omitempty"`
}
