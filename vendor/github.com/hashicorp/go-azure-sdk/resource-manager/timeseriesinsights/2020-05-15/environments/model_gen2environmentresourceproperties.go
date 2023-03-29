package environments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Gen2EnvironmentResourceProperties struct {
	CreationTime           *string                           `json:"creationTime,omitempty"`
	DataAccessFqdn         *string                           `json:"dataAccessFqdn,omitempty"`
	DataAccessId           *string                           `json:"dataAccessId,omitempty"`
	ProvisioningState      *ProvisioningState                `json:"provisioningState,omitempty"`
	Status                 *EnvironmentStatus                `json:"status,omitempty"`
	StorageConfiguration   Gen2StorageConfigurationOutput    `json:"storageConfiguration"`
	TimeSeriesIdProperties []TimeSeriesIdProperty            `json:"timeSeriesIdProperties"`
	WarmStoreConfiguration *WarmStoreConfigurationProperties `json:"warmStoreConfiguration,omitempty"`
}

func (o *Gen2EnvironmentResourceProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *Gen2EnvironmentResourceProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}
