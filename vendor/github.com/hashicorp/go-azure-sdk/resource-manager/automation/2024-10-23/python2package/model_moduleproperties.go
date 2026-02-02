package python2package

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ModuleProperties struct {
	ActivityCount     *int64                   `json:"activityCount,omitempty"`
	ContentLink       *ContentLink             `json:"contentLink,omitempty"`
	CreationTime      *string                  `json:"creationTime,omitempty"`
	Description       *string                  `json:"description,omitempty"`
	Error             *ModuleErrorInfo         `json:"error,omitempty"`
	IsComposite       *bool                    `json:"isComposite,omitempty"`
	IsGlobal          *bool                    `json:"isGlobal,omitempty"`
	LastModifiedTime  *string                  `json:"lastModifiedTime,omitempty"`
	ProvisioningState *ModuleProvisioningState `json:"provisioningState,omitempty"`
	SizeInBytes       *int64                   `json:"sizeInBytes,omitempty"`
	Version           *string                  `json:"version,omitempty"`
}

func (o *ModuleProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ModuleProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *ModuleProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ModuleProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
