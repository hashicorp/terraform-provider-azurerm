package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPropertiesForUpdate struct {
	AdministratorLoginPassword *string              `json:"administratorLoginPassword,omitempty"`
	AuthConfig                 *AuthConfig          `json:"authConfig"`
	Backup                     *Backup              `json:"backup"`
	CreateMode                 *CreateModeForUpdate `json:"createMode,omitempty"`
	DataEncryption             *DataEncryption      `json:"dataEncryption"`
	HighAvailability           *HighAvailability    `json:"highAvailability"`
	MaintenanceWindow          *MaintenanceWindow   `json:"maintenanceWindow"`
	ReplicationRole            *ReplicationRole     `json:"replicationRole,omitempty"`
	Storage                    *Storage             `json:"storage"`
	Version                    *ServerVersion       `json:"version,omitempty"`
}
