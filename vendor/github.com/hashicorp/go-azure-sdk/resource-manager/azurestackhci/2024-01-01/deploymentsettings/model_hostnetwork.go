package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HostNetwork struct {
	EnableStorageAutoIP           *bool              `json:"enableStorageAutoIp,omitempty"`
	Intents                       *[]Intents         `json:"intents,omitempty"`
	StorageConnectivitySwitchless *bool              `json:"storageConnectivitySwitchless,omitempty"`
	StorageNetworks               *[]StorageNetworks `json:"storageNetworks,omitempty"`
}
