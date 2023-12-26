package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	AssociatedWorkspaces          *[]AssociatedWorkspace         `json:"associatedWorkspaces,omitempty"`
	BillingType                   *BillingType                   `json:"billingType,omitempty"`
	CapacityReservationProperties *CapacityReservationProperties `json:"capacityReservationProperties,omitempty"`
	ClusterId                     *string                        `json:"clusterId,omitempty"`
	CreatedDate                   *string                        `json:"createdDate,omitempty"`
	IsAvailabilityZonesEnabled    *bool                          `json:"isAvailabilityZonesEnabled,omitempty"`
	IsDoubleEncryptionEnabled     *bool                          `json:"isDoubleEncryptionEnabled,omitempty"`
	KeyVaultProperties            *KeyVaultProperties            `json:"keyVaultProperties,omitempty"`
	LastModifiedDate              *string                        `json:"lastModifiedDate,omitempty"`
	ProvisioningState             *ClusterEntityStatus           `json:"provisioningState,omitempty"`
}
