package registrymanagement

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryRegionArmDetails struct {
	AcrDetails            *[]AcrDetails            `json:"acrDetails,omitempty"`
	Location              *string                  `json:"location,omitempty"`
	StorageAccountDetails *[]StorageAccountDetails `json:"storageAccountDetails,omitempty"`
}
