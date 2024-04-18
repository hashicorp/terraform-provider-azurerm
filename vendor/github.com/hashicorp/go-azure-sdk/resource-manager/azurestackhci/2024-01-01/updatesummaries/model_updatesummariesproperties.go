package updatesummaries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateSummariesProperties struct {
	CurrentVersion    *string                         `json:"currentVersion,omitempty"`
	HardwareModel     *string                         `json:"hardwareModel,omitempty"`
	HealthCheckDate   *string                         `json:"healthCheckDate,omitempty"`
	HealthCheckResult *[]PrecheckResult               `json:"healthCheckResult,omitempty"`
	HealthState       *HealthState                    `json:"healthState,omitempty"`
	LastChecked       *string                         `json:"lastChecked,omitempty"`
	LastUpdated       *string                         `json:"lastUpdated,omitempty"`
	OemFamily         *string                         `json:"oemFamily,omitempty"`
	PackageVersions   *[]PackageVersionInfo           `json:"packageVersions,omitempty"`
	ProvisioningState *ProvisioningState              `json:"provisioningState,omitempty"`
	State             *UpdateSummariesPropertiesState `json:"state,omitempty"`
}

func (o *UpdateSummariesProperties) GetHealthCheckDateAsTime() (*time.Time, error) {
	if o.HealthCheckDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.HealthCheckDate, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummariesProperties) SetHealthCheckDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.HealthCheckDate = &formatted
}

func (o *UpdateSummariesProperties) GetLastCheckedAsTime() (*time.Time, error) {
	if o.LastChecked == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastChecked, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummariesProperties) SetLastCheckedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastChecked = &formatted
}

func (o *UpdateSummariesProperties) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummariesProperties) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}
