package integrationruntime

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedIntegrationRuntimeError struct {
	Code       *string   `json:"code,omitempty"`
	Message    *string   `json:"message,omitempty"`
	Parameters *[]string `json:"parameters,omitempty"`
	Time       *string   `json:"time,omitempty"`
}

func (o *ManagedIntegrationRuntimeError) GetTimeAsTime() (*time.Time, error) {
	if o.Time == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Time, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagedIntegrationRuntimeError) SetTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Time = &formatted
}
