package servers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerProperties struct {
	AdministratorLogin         *string            `json:"administratorLogin,omitempty"`
	AdministratorLoginPassword *string            `json:"administratorLoginPassword,omitempty"`
	AvailabilityZone           *string            `json:"availabilityZone,omitempty"`
	Backup                     *Backup            `json:"backup,omitempty"`
	CreateMode                 *CreateMode        `json:"createMode,omitempty"`
	DataEncryption             *DataEncryption    `json:"dataEncryption,omitempty"`
	FullyQualifiedDomainName   *string            `json:"fullyQualifiedDomainName,omitempty"`
	HighAvailability           *HighAvailability  `json:"highAvailability,omitempty"`
	MaintenanceWindow          *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	Network                    *Network           `json:"network,omitempty"`
	ReplicaCapacity            *int64             `json:"replicaCapacity,omitempty"`
	ReplicationRole            *ReplicationRole   `json:"replicationRole,omitempty"`
	RestorePointInTime         *string            `json:"restorePointInTime,omitempty"`
	SourceServerResourceId     *string            `json:"sourceServerResourceId,omitempty"`
	State                      *ServerState       `json:"state,omitempty"`
	Storage                    *Storage           `json:"storage,omitempty"`
	Version                    *ServerVersion     `json:"version,omitempty"`
}

func (o *ServerProperties) GetRestorePointInTimeAsTime() (*time.Time, error) {
	if o.RestorePointInTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RestorePointInTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerProperties) SetRestorePointInTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RestorePointInTime = &formatted
}
