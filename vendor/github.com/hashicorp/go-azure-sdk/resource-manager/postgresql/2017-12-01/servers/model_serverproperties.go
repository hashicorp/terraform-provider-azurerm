package servers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerProperties struct {
	AdministratorLogin         *string                            `json:"administratorLogin,omitempty"`
	ByokEnforcement            *string                            `json:"byokEnforcement,omitempty"`
	EarliestRestoreDate        *string                            `json:"earliestRestoreDate,omitempty"`
	FullyQualifiedDomainName   *string                            `json:"fullyQualifiedDomainName,omitempty"`
	InfrastructureEncryption   *InfrastructureEncryption          `json:"infrastructureEncryption,omitempty"`
	MasterServerId             *string                            `json:"masterServerId,omitempty"`
	MinimalTlsVersion          *MinimalTlsVersionEnum             `json:"minimalTlsVersion,omitempty"`
	PrivateEndpointConnections *[]ServerPrivateEndpointConnection `json:"privateEndpointConnections,omitempty"`
	PublicNetworkAccess        *PublicNetworkAccessEnum           `json:"publicNetworkAccess,omitempty"`
	ReplicaCapacity            *int64                             `json:"replicaCapacity,omitempty"`
	ReplicationRole            *string                            `json:"replicationRole,omitempty"`
	SslEnforcement             *SslEnforcementEnum                `json:"sslEnforcement,omitempty"`
	StorageProfile             *StorageProfile                    `json:"storageProfile,omitempty"`
	UserVisibleState           *ServerState                       `json:"userVisibleState,omitempty"`
	Version                    *ServerVersion                     `json:"version,omitempty"`
}

func (o *ServerProperties) GetEarliestRestoreDateAsTime() (*time.Time, error) {
	if o.EarliestRestoreDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EarliestRestoreDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerProperties) SetEarliestRestoreDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EarliestRestoreDate = &formatted
}
