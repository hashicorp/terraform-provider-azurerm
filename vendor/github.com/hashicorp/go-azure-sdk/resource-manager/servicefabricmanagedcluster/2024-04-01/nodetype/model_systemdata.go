package nodetype

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SystemData struct {
	CreatedAt          *string `json:"createdAt,omitempty"`
	CreatedBy          *string `json:"createdBy,omitempty"`
	CreatedByType      *string `json:"createdByType,omitempty"`
	LastModifiedAt     *string `json:"lastModifiedAt,omitempty"`
	LastModifiedBy     *string `json:"lastModifiedBy,omitempty"`
	LastModifiedByType *string `json:"lastModifiedByType,omitempty"`
}

func (o *SystemData) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SystemData) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *SystemData) GetLastModifiedAtAsTime() (*time.Time, error) {
	if o.LastModifiedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SystemData) SetLastModifiedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedAt = &formatted
}
