package connectiontype

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionTypeProperties struct {
	CreationTime     *string                     `json:"creationTime,omitempty"`
	Description      *string                     `json:"description,omitempty"`
	FieldDefinitions *map[string]FieldDefinition `json:"fieldDefinitions,omitempty"`
	IsGlobal         *bool                       `json:"isGlobal,omitempty"`
	LastModifiedTime *string                     `json:"lastModifiedTime,omitempty"`
}

func (o *ConnectionTypeProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectionTypeProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *ConnectionTypeProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectionTypeProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
