package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPropertiesForPatch struct {
	AdministratorLogin         *string                    `json:"administratorLogin,omitempty"`
	AdministratorLoginPassword *string                    `json:"administratorLoginPassword,omitempty"`
	AuthConfig                 *AuthConfigForPatch        `json:"authConfig,omitempty"`
	AvailabilityZone           *string                    `json:"availabilityZone,omitempty"`
	Backup                     *BackupForPatch            `json:"backup,omitempty"`
	Cluster                    *Cluster                   `json:"cluster,omitempty"`
	CreateMode                 *CreateModeForPatch        `json:"createMode,omitempty"`
	DataEncryption             *DataEncryption            `json:"dataEncryption,omitempty"`
	HighAvailability           *HighAvailabilityForPatch  `json:"highAvailability,omitempty"`
	MaintenanceWindow          *MaintenanceWindowForPatch `json:"maintenanceWindow,omitempty"`
	Network                    *Network                   `json:"network,omitempty"`
	Replica                    *Replica                   `json:"replica,omitempty"`
	ReplicationRole            *ReplicationRole           `json:"replicationRole,omitempty"`
	Storage                    *Storage                   `json:"storage,omitempty"`
	Version                    *PostgresMajorVersion      `json:"version,omitempty"`
}
