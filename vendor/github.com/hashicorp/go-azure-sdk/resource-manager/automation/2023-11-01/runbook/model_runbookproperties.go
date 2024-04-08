package runbook

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookProperties struct {
	CreationTime       *string                      `json:"creationTime,omitempty"`
	Description        *string                      `json:"description,omitempty"`
	Draft              *RunbookDraft                `json:"draft,omitempty"`
	JobCount           *int64                       `json:"jobCount,omitempty"`
	LastModifiedBy     *string                      `json:"lastModifiedBy,omitempty"`
	LastModifiedTime   *string                      `json:"lastModifiedTime,omitempty"`
	LogActivityTrace   *int64                       `json:"logActivityTrace,omitempty"`
	LogProgress        *bool                        `json:"logProgress,omitempty"`
	LogVerbose         *bool                        `json:"logVerbose,omitempty"`
	OutputTypes        *[]string                    `json:"outputTypes,omitempty"`
	Parameters         *map[string]RunbookParameter `json:"parameters,omitempty"`
	ProvisioningState  *RunbookProvisioningState    `json:"provisioningState,omitempty"`
	PublishContentLink *ContentLink                 `json:"publishContentLink,omitempty"`
	RunbookType        *RunbookTypeEnum             `json:"runbookType,omitempty"`
	State              *RunbookState                `json:"state,omitempty"`
}

func (o *RunbookProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunbookProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *RunbookProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunbookProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
