package tenantconfigurationsyncstate

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TenantConfigurationSyncStateContractProperties struct {
	Branch                  *string `json:"branch,omitempty"`
	CommitId                *string `json:"commitId,omitempty"`
	ConfigurationChangeDate *string `json:"configurationChangeDate,omitempty"`
	IsExport                *bool   `json:"isExport,omitempty"`
	IsGitEnabled            *bool   `json:"isGitEnabled,omitempty"`
	IsSynced                *bool   `json:"isSynced,omitempty"`
	LastOperationId         *string `json:"lastOperationId,omitempty"`
	SyncDate                *string `json:"syncDate,omitempty"`
}

func (o *TenantConfigurationSyncStateContractProperties) GetConfigurationChangeDateAsTime() (*time.Time, error) {
	if o.ConfigurationChangeDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ConfigurationChangeDate, "2006-01-02T15:04:05Z07:00")
}

func (o *TenantConfigurationSyncStateContractProperties) SetConfigurationChangeDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ConfigurationChangeDate = &formatted
}

func (o *TenantConfigurationSyncStateContractProperties) GetSyncDateAsTime() (*time.Time, error) {
	if o.SyncDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SyncDate, "2006-01-02T15:04:05Z07:00")
}

func (o *TenantConfigurationSyncStateContractProperties) SetSyncDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SyncDate = &formatted
}
