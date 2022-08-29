package workbooksapis

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkbookProperties struct {
	Category       string    `json:"category"`
	Description    *string   `json:"description,omitempty"`
	DisplayName    string    `json:"displayName"`
	Revision       *string   `json:"revision,omitempty"`
	SerializedData string    `json:"serializedData"`
	SourceId       *string   `json:"sourceId,omitempty"`
	StorageUri     *string   `json:"storageUri,omitempty"`
	Tags           *[]string `json:"tags,omitempty"`
	TimeModified   *string   `json:"timeModified,omitempty"`
	UserId         *string   `json:"userId,omitempty"`
	Version        *string   `json:"version,omitempty"`
}

func (o *WorkbookProperties) GetTimeModifiedAsTime() (*time.Time, error) {
	if o.TimeModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeModified, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkbookProperties) SetTimeModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeModified = &formatted
}
