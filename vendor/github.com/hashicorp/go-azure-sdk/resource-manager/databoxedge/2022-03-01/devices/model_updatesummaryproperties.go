package devices

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateSummaryProperties struct {
	DeviceLastScannedDateTime            *string                `json:"deviceLastScannedDateTime,omitempty"`
	DeviceVersionNumber                  *string                `json:"deviceVersionNumber,omitempty"`
	FriendlyDeviceVersionName            *string                `json:"friendlyDeviceVersionName,omitempty"`
	InProgressDownloadJobId              *string                `json:"inProgressDownloadJobId,omitempty"`
	InProgressDownloadJobStartedDateTime *string                `json:"inProgressDownloadJobStartedDateTime,omitempty"`
	InProgressInstallJobId               *string                `json:"inProgressInstallJobId,omitempty"`
	InProgressInstallJobStartedDateTime  *string                `json:"inProgressInstallJobStartedDateTime,omitempty"`
	LastCompletedDownloadJobDateTime     *string                `json:"lastCompletedDownloadJobDateTime,omitempty"`
	LastCompletedDownloadJobId           *string                `json:"lastCompletedDownloadJobId,omitempty"`
	LastCompletedInstallJobDateTime      *string                `json:"lastCompletedInstallJobDateTime,omitempty"`
	LastCompletedInstallJobId            *string                `json:"lastCompletedInstallJobId,omitempty"`
	LastCompletedScanJobDateTime         *string                `json:"lastCompletedScanJobDateTime,omitempty"`
	LastDownloadJobStatus                *JobStatus             `json:"lastDownloadJobStatus,omitempty"`
	LastInstallJobStatus                 *JobStatus             `json:"lastInstallJobStatus,omitempty"`
	LastSuccessfulInstallJobDateTime     *string                `json:"lastSuccessfulInstallJobDateTime,omitempty"`
	LastSuccessfulScanJobTime            *string                `json:"lastSuccessfulScanJobTime,omitempty"`
	OngoingUpdateOperation               *UpdateOperation       `json:"ongoingUpdateOperation,omitempty"`
	RebootBehavior                       *InstallRebootBehavior `json:"rebootBehavior,omitempty"`
	TotalNumberOfUpdatesAvailable        *int64                 `json:"totalNumberOfUpdatesAvailable,omitempty"`
	TotalNumberOfUpdatesPendingDownload  *int64                 `json:"totalNumberOfUpdatesPendingDownload,omitempty"`
	TotalNumberOfUpdatesPendingInstall   *int64                 `json:"totalNumberOfUpdatesPendingInstall,omitempty"`
	TotalTimeInMinutes                   *int64                 `json:"totalTimeInMinutes,omitempty"`
	TotalUpdateSizeInBytes               *float64               `json:"totalUpdateSizeInBytes,omitempty"`
	UpdateTitles                         *[]string              `json:"updateTitles,omitempty"`
	Updates                              *[]UpdateDetails       `json:"updates,omitempty"`
}

func (o *UpdateSummaryProperties) GetDeviceLastScannedDateTimeAsTime() (*time.Time, error) {
	if o.DeviceLastScannedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeviceLastScannedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetDeviceLastScannedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeviceLastScannedDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetInProgressDownloadJobStartedDateTimeAsTime() (*time.Time, error) {
	if o.InProgressDownloadJobStartedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.InProgressDownloadJobStartedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetInProgressDownloadJobStartedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.InProgressDownloadJobStartedDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetInProgressInstallJobStartedDateTimeAsTime() (*time.Time, error) {
	if o.InProgressInstallJobStartedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.InProgressInstallJobStartedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetInProgressInstallJobStartedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.InProgressInstallJobStartedDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetLastCompletedDownloadJobDateTimeAsTime() (*time.Time, error) {
	if o.LastCompletedDownloadJobDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletedDownloadJobDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetLastCompletedDownloadJobDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletedDownloadJobDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetLastCompletedInstallJobDateTimeAsTime() (*time.Time, error) {
	if o.LastCompletedInstallJobDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletedInstallJobDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetLastCompletedInstallJobDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletedInstallJobDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetLastCompletedScanJobDateTimeAsTime() (*time.Time, error) {
	if o.LastCompletedScanJobDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletedScanJobDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetLastCompletedScanJobDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletedScanJobDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetLastSuccessfulInstallJobDateTimeAsTime() (*time.Time, error) {
	if o.LastSuccessfulInstallJobDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSuccessfulInstallJobDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetLastSuccessfulInstallJobDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSuccessfulInstallJobDateTime = &formatted
}

func (o *UpdateSummaryProperties) GetLastSuccessfulScanJobTimeAsTime() (*time.Time, error) {
	if o.LastSuccessfulScanJobTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSuccessfulScanJobTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateSummaryProperties) SetLastSuccessfulScanJobTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSuccessfulScanJobTime = &formatted
}
