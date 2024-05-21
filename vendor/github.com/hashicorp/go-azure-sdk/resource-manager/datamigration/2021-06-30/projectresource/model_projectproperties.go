package projectresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProperties struct {
	CreationTime         *string                   `json:"creationTime,omitempty"`
	DatabasesInfo        *[]DatabaseInfo           `json:"databasesInfo,omitempty"`
	ProvisioningState    *ProjectProvisioningState `json:"provisioningState,omitempty"`
	SourceConnectionInfo *ConnectionInfo           `json:"sourceConnectionInfo,omitempty"`
	SourcePlatform       ProjectSourcePlatform     `json:"sourcePlatform"`
	TargetConnectionInfo *ConnectionInfo           `json:"targetConnectionInfo,omitempty"`
	TargetPlatform       ProjectTargetPlatform     `json:"targetPlatform"`
}

func (o *ProjectProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ProjectProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
