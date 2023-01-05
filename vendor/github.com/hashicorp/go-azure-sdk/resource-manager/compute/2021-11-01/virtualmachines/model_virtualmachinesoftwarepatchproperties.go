package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineSoftwarePatchProperties struct {
	ActivityId           *string                     `json:"activityId,omitempty"`
	AssessmentState      *PatchAssessmentState       `json:"assessmentState,omitempty"`
	Classifications      *[]string                   `json:"classifications,omitempty"`
	KbId                 *string                     `json:"kbId,omitempty"`
	LastModifiedDateTime *string                     `json:"lastModifiedDateTime,omitempty"`
	Name                 *string                     `json:"name,omitempty"`
	PatchId              *string                     `json:"patchId,omitempty"`
	PublishedDate        *string                     `json:"publishedDate,omitempty"`
	RebootBehavior       *VMGuestPatchRebootBehavior `json:"rebootBehavior,omitempty"`
	Version              *string                     `json:"version,omitempty"`
}

func (o *VirtualMachineSoftwarePatchProperties) GetLastModifiedDateTimeAsTime() (*time.Time, error) {
	if o.LastModifiedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineSoftwarePatchProperties) SetLastModifiedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDateTime = &formatted
}

func (o *VirtualMachineSoftwarePatchProperties) GetPublishedDateAsTime() (*time.Time, error) {
	if o.PublishedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PublishedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *VirtualMachineSoftwarePatchProperties) SetPublishedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PublishedDate = &formatted
}
