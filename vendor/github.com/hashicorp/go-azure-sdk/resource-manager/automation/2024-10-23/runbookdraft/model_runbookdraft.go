package runbookdraft

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunbookDraft struct {
	CreationTime     *string                      `json:"creationTime,omitempty"`
	DraftContentLink *ContentLink                 `json:"draftContentLink,omitempty"`
	InEdit           *bool                        `json:"inEdit,omitempty"`
	LastModifiedTime *string                      `json:"lastModifiedTime,omitempty"`
	OutputTypes      *[]string                    `json:"outputTypes,omitempty"`
	Parameters       *map[string]RunbookParameter `json:"parameters,omitempty"`
}

func (o *RunbookDraft) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunbookDraft) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *RunbookDraft) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunbookDraft) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
