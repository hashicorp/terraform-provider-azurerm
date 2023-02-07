package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataMigrationServiceProperties struct {
	ProvisioningState *ServiceProvisioningState `json:"provisioningState,omitempty"`
	PublicKey         *string                   `json:"publicKey,omitempty"`
	VirtualSubnetId   string                    `json:"virtualSubnetId"`
}
