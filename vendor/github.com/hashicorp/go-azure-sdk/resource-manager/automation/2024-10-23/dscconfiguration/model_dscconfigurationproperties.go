package dscconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscConfigurationProperties struct {
	CreationTime           *string                               `json:"creationTime,omitempty"`
	Description            *string                               `json:"description,omitempty"`
	JobCount               *int64                                `json:"jobCount,omitempty"`
	LastModifiedTime       *string                               `json:"lastModifiedTime,omitempty"`
	LogVerbose             *bool                                 `json:"logVerbose,omitempty"`
	NodeConfigurationCount *int64                                `json:"nodeConfigurationCount,omitempty"`
	Parameters             *map[string]DscConfigurationParameter `json:"parameters,omitempty"`
	ProvisioningState      *DscConfigurationProvisioningState    `json:"provisioningState,omitempty"`
	Source                 *ContentSource                        `json:"source,omitempty"`
	State                  *DscConfigurationState                `json:"state,omitempty"`
}

func (o *DscConfigurationProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DscConfigurationProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *DscConfigurationProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DscConfigurationProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
