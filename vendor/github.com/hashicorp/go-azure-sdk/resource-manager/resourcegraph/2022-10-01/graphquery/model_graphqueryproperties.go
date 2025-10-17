package graphquery

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphQueryProperties struct {
	Description  *string     `json:"description,omitempty"`
	Query        string      `json:"query"`
	ResultKind   *ResultKind `json:"resultKind,omitempty"`
	TimeModified *string     `json:"timeModified,omitempty"`
}

func (o *GraphQueryProperties) GetTimeModifiedAsTime() (*time.Time, error) {
	if o.TimeModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeModified, "2006-01-02T15:04:05Z07:00")
}

func (o *GraphQueryProperties) SetTimeModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeModified = &formatted
}
