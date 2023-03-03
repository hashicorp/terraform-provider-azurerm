package watcher

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WatcherProperties struct {
	CreationTime                *string            `json:"creationTime,omitempty"`
	Description                 *string            `json:"description,omitempty"`
	ExecutionFrequencyInSeconds *int64             `json:"executionFrequencyInSeconds,omitempty"`
	LastModifiedBy              *string            `json:"lastModifiedBy,omitempty"`
	LastModifiedTime            *string            `json:"lastModifiedTime,omitempty"`
	ScriptName                  *string            `json:"scriptName,omitempty"`
	ScriptParameters            *map[string]string `json:"scriptParameters,omitempty"`
	ScriptRunOn                 *string            `json:"scriptRunOn,omitempty"`
	Status                      *string            `json:"status,omitempty"`
}

func (o *WatcherProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WatcherProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *WatcherProperties) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *WatcherProperties) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}
