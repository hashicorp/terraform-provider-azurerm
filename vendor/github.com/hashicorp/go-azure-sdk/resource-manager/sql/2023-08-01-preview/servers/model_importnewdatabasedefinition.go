package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportNewDatabaseDefinition struct {
	AdministratorLogin         string                    `json:"administratorLogin"`
	AdministratorLoginPassword string                    `json:"administratorLoginPassword"`
	AuthenticationType         *string                   `json:"authenticationType,omitempty"`
	DatabaseName               *string                   `json:"databaseName,omitempty"`
	Edition                    *string                   `json:"edition,omitempty"`
	MaxSizeBytes               *string                   `json:"maxSizeBytes,omitempty"`
	NetworkIsolation           *NetworkIsolationSettings `json:"networkIsolation,omitempty"`
	ServiceObjectiveName       *string                   `json:"serviceObjectiveName,omitempty"`
	StorageKey                 string                    `json:"storageKey"`
	StorageKeyType             StorageKeyType            `json:"storageKeyType"`
	StorageUri                 string                    `json:"storageUri"`
}
