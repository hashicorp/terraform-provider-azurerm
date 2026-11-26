package clusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProperties struct {
	AiModelManagement               *AiModelManagementStatus        `json:"aiModelManagement,omitempty"`
	Backup                          *BackupProperties               `json:"backup,omitempty"`
	CreationMode                    *ClusterCreateMode              `json:"creationMode,omitempty"`
	CreationSourceClusterLocation   *string                         `json:"creationSourceClusterLocation,omitempty"`
	CreationSourceClusterResourceId *string                         `json:"creationSourceClusterResourceId,omitempty"`
	DataEncryption                  *HorizonDbClusterDataEncryption `json:"dataEncryption,omitempty"`
	Network                         *Network                        `json:"network,omitempty"`
	PointInTimeUtc                  *string                         `json:"pointInTimeUtc,omitempty"`
	PostgreSqlVersion               *PostgreSqlVersion              `json:"postgreSqlVersion,omitempty"`
	ProvisioningState               ProvisioningState               `json:"provisioningState"`
	ReadWritePool                   ReadWritePoolProperties         `json:"readWritePool"`
	State                           ClusterState                    `json:"state"`
}

func (o *ClusterProperties) GetPointInTimeUtcAsTime() (*time.Time, error) {
	if o.PointInTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PointInTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ClusterProperties) SetPointInTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PointInTimeUtc = &formatted
}
