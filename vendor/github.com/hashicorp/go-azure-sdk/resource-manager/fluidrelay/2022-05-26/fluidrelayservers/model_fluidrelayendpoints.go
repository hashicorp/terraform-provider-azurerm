package fluidrelayservers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FluidRelayEndpoints struct {
	OrdererEndpoints *[]string `json:"ordererEndpoints,omitempty"`
	ServiceEndpoints *[]string `json:"serviceEndpoints,omitempty"`
	StorageEndpoints *[]string `json:"storageEndpoints,omitempty"`
}
