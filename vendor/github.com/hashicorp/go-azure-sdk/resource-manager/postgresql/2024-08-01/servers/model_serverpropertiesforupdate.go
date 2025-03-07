package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPropertiesForUpdate struct {
	AdministratorLogin         *string              `json:"administratorLogin,omitempty"`
	AdministratorLoginPassword *string              `json:"administratorLoginPassword,omitempty"`
	AuthConfig                 *AuthConfig          `json:"authConfig,omitempty"`
	Backup                     *Backup              `json:"backup,omitempty"`
	CreateMode                 *CreateModeForUpdate `json:"createMode,omitempty"`
	DataEncryption             *DataEncryption      `json:"dataEncryption,omitempty"`
	HighAvailability           *HighAvailability    `json:"highAvailability,omitempty"`
	MaintenanceWindow          *MaintenanceWindow   `json:"maintenanceWindow,omitempty"`
	Network                    *Network             `json:"network,omitempty"`
	Replica                    *Replica             `json:"replica,omitempty"`
	ReplicationRole            *ReplicationRole     `json:"replicationRole,omitempty"`
	Storage                    *Storage             `json:"storage,omitempty"`
	Version                    *ServerVersion       `json:"version,omitempty"`
}
