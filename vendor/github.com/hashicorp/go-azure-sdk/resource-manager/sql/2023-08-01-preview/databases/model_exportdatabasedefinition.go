package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportDatabaseDefinition struct {
	AdministratorLogin         string                    `json:"administratorLogin"`
	AdministratorLoginPassword string                    `json:"administratorLoginPassword"`
	AuthenticationType         *string                   `json:"authenticationType,omitempty"`
	NetworkIsolation           *NetworkIsolationSettings `json:"networkIsolation,omitempty"`
	StorageKey                 string                    `json:"storageKey"`
	StorageKeyType             StorageKeyType            `json:"storageKeyType"`
	StorageUri                 string                    `json:"storageUri"`
}
