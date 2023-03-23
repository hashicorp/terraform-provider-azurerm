package datastores

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreProperties struct {
	DiskPoolVolume    *DiskPoolVolume             `json:"diskPoolVolume,omitempty"`
	NetAppVolume      *NetAppVolume               `json:"netAppVolume,omitempty"`
	ProvisioningState *DatastoreProvisioningState `json:"provisioningState,omitempty"`
	Status            *DatastoreStatus            `json:"status,omitempty"`
}
