package nodereports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscReportResource struct {
	DependsOn         *[]DscReportResourceNavigation `json:"dependsOn,omitempty"`
	DurationInSeconds *float64                       `json:"durationInSeconds,omitempty"`
	Error             *string                        `json:"error,omitempty"`
	ModuleName        *string                        `json:"moduleName,omitempty"`
	ModuleVersion     *string                        `json:"moduleVersion,omitempty"`
	ResourceId        *string                        `json:"resourceId,omitempty"`
	ResourceName      *string                        `json:"resourceName,omitempty"`
	SourceInfo        *string                        `json:"sourceInfo,omitempty"`
	StartDate         *string                        `json:"startDate,omitempty"`
	Status            *string                        `json:"status,omitempty"`
}

func (o *DscReportResource) GetStartDateAsTime() (*time.Time, error) {
	if o.StartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *DscReportResource) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = &formatted
}
