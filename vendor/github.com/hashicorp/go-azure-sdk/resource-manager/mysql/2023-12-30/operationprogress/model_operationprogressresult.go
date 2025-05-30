package operationprogress

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationProgressResult struct {
	EndTime         *string                       `json:"endTime,omitempty"`
	Error           *ErrorDetail                  `json:"error,omitempty"`
	Id              *string                       `json:"id,omitempty"`
	Name            *string                       `json:"name,omitempty"`
	Operations      *[]OperationStatusResult      `json:"operations,omitempty"`
	PercentComplete *float64                      `json:"percentComplete,omitempty"`
	Properties      OperationProgressResponseType `json:"properties"`
	ResourceId      *string                       `json:"resourceId,omitempty"`
	StartTime       *string                       `json:"startTime,omitempty"`
	Status          string                        `json:"status"`
}

func (o *OperationProgressResult) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *OperationProgressResult) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *OperationProgressResult) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *OperationProgressResult) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

var _ json.Unmarshaler = &OperationProgressResult{}

func (s *OperationProgressResult) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EndTime         *string                  `json:"endTime,omitempty"`
		Error           *ErrorDetail             `json:"error,omitempty"`
		Id              *string                  `json:"id,omitempty"`
		Name            *string                  `json:"name,omitempty"`
		Operations      *[]OperationStatusResult `json:"operations,omitempty"`
		PercentComplete *float64                 `json:"percentComplete,omitempty"`
		ResourceId      *string                  `json:"resourceId,omitempty"`
		StartTime       *string                  `json:"startTime,omitempty"`
		Status          string                   `json:"status"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EndTime = decoded.EndTime
	s.Error = decoded.Error
	s.Id = decoded.Id
	s.Name = decoded.Name
	s.Operations = decoded.Operations
	s.PercentComplete = decoded.PercentComplete
	s.ResourceId = decoded.ResourceId
	s.StartTime = decoded.StartTime
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OperationProgressResult into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := UnmarshalOperationProgressResponseTypeImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'OperationProgressResult': %+v", err)
		}
		s.Properties = impl
	}

	return nil
}
