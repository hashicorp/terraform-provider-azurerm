package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerGroupProfileStub struct {
	ContainerGroupProperties *NGroupContainerGroupProperties `json:"containerGroupProperties,omitempty"`
	NetworkProfile           *NetworkProfile                 `json:"networkProfile,omitempty"`
	Resource                 *ApiEntityReference             `json:"resource,omitempty"`
	Revision                 *int64                          `json:"revision,omitempty"`
	StorageProfile           *StorageProfile                 `json:"storageProfile,omitempty"`
}
