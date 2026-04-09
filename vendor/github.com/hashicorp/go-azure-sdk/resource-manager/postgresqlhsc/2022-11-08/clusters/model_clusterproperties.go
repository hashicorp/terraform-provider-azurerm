package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	AdministratorLogin              *string                            `json:"administratorLogin,omitempty"`
	AdministratorLoginPassword      *string                            `json:"administratorLoginPassword,omitempty"`
	CitusVersion                    *string                            `json:"citusVersion,omitempty"`
	CoordinatorEnablePublicIPAccess *bool                              `json:"coordinatorEnablePublicIpAccess,omitempty"`
	CoordinatorServerEdition        *string                            `json:"coordinatorServerEdition,omitempty"`
	CoordinatorStorageQuotaInMb     *int64                             `json:"coordinatorStorageQuotaInMb,omitempty"`
	CoordinatorVCores               *int64                             `json:"coordinatorVCores,omitempty"`
	EarliestRestoreTime             *string                            `json:"earliestRestoreTime,omitempty"`
	EnableHa                        *bool                              `json:"enableHa,omitempty"`
	EnableShardsOnCoordinator       *bool                              `json:"enableShardsOnCoordinator,omitempty"`
	MaintenanceWindow               *MaintenanceWindow                 `json:"maintenanceWindow,omitempty"`
	NodeCount                       *int64                             `json:"nodeCount,omitempty"`
	NodeEnablePublicIPAccess        *bool                              `json:"nodeEnablePublicIpAccess,omitempty"`
	NodeServerEdition               *string                            `json:"nodeServerEdition,omitempty"`
	NodeStorageQuotaInMb            *int64                             `json:"nodeStorageQuotaInMb,omitempty"`
	NodeVCores                      *int64                             `json:"nodeVCores,omitempty"`
	PointInTimeUTC                  *string                            `json:"pointInTimeUTC,omitempty"`
	PostgresqlVersion               *string                            `json:"postgresqlVersion,omitempty"`
	PreferredPrimaryZone            *string                            `json:"preferredPrimaryZone,omitempty"`
	PrivateEndpointConnections      *[]SimplePrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	ProvisioningState               *string                            `json:"provisioningState,omitempty"`
	ReadReplicas                    *[]string                          `json:"readReplicas,omitempty"`
	ServerNames                     *[]ServerNameItem                  `json:"serverNames,omitempty"`
	SourceLocation                  *string                            `json:"sourceLocation,omitempty"`
	SourceResourceId                *string                            `json:"sourceResourceId,omitempty"`
	State                           *string                            `json:"state,omitempty"`
}

func (o *ClusterProperties) GetEarliestRestoreTimeAsTime() (*time.Time, error) {
	if o.EarliestRestoreTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestoreTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetEarliestRestoreTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestoreTime = &formatted
}

func (o *ClusterProperties) GetPointInTimeUTCAsTime() (*time.Time, error) {
	if o.PointInTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PointInTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetPointInTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PointInTimeUTC = &formatted
}
