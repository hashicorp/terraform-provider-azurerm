package connections

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiConnectionDefinitionProperties struct {
	Api                      *ApiReference                 `json:"api,omitempty"`
	ChangedTime              *string                       `json:"changedTime,omitempty"`
	CreatedTime              *string                       `json:"createdTime,omitempty"`
	CustomParameterValues    *map[string]string            `json:"customParameterValues,omitempty"`
	DisplayName              *string                       `json:"displayName,omitempty"`
	NonSecretParameterValues *map[string]string            `json:"nonSecretParameterValues,omitempty"`
	ParameterValues          *map[string]string            `json:"parameterValues,omitempty"`
	Statuses                 *[]ConnectionStatusDefinition `json:"statuses,omitempty"`
	TestLinks                *[]ApiConnectionTestLink      `json:"testLinks,omitempty"`
}

func (o *ApiConnectionDefinitionProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ApiConnectionDefinitionProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *ApiConnectionDefinitionProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ApiConnectionDefinitionProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
