package integrationaccountmaps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountMapProperties struct {
	ChangedTime      *string                                          `json:"changedTime,omitempty"`
	Content          *string                                          `json:"content,omitempty"`
	ContentLink      *ContentLink                                     `json:"contentLink,omitempty"`
	ContentType      *string                                          `json:"contentType,omitempty"`
	CreatedTime      *string                                          `json:"createdTime,omitempty"`
	MapType          MapType                                          `json:"mapType"`
	Metadata         *interface{}                                     `json:"metadata,omitempty"`
	ParametersSchema *IntegrationAccountMapPropertiesParametersSchema `json:"parametersSchema,omitempty"`
}

func (o *IntegrationAccountMapProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountMapProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *IntegrationAccountMapProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountMapProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
