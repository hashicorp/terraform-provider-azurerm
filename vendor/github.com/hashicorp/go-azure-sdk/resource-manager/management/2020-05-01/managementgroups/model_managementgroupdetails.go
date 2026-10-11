package managementgroups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupDetails struct {
	Parent      *ParentGroupInfo              `json:"parent,omitempty"`
	Path        *[]ManagementGroupPathElement `json:"path,omitempty"`
	UpdatedBy   *string                       `json:"updatedBy,omitempty"`
	UpdatedTime *string                       `json:"updatedTime,omitempty"`
	Version     *float64                      `json:"version,omitempty"`
}

func (o *ManagementGroupDetails) GetUpdatedTimeAsTime() (*time.Time, error) {
	if o.UpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagementGroupDetails) SetUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedTime = &formatted
}
