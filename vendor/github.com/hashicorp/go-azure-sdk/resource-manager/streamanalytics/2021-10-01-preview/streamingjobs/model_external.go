package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type External struct {
	Container            *string               `json:"container,omitempty"`
	Path                 *string               `json:"path,omitempty"`
	RefreshConfiguration *RefreshConfiguration `json:"refreshConfiguration,omitempty"`
	StorageAccount       *StorageAccount       `json:"storageAccount,omitempty"`
}
