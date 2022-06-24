package operationstatus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationResource struct {
	EndTime    *string               `json:"endTime,omitempty"`
	Error      *Error                `json:"error,omitempty"`
	Id         *string               `json:"id,omitempty"`
	Name       *string               `json:"name,omitempty"`
	Properties OperationExtendedInfo `json:"properties"`
	StartTime  *string               `json:"startTime,omitempty"`
	Status     *string               `json:"status,omitempty"`
}

func (o *OperationResource) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *OperationResource) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *OperationResource) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *OperationResource) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

var _ json.Unmarshaler = &OperationResource{}

func (s *OperationResource) UnmarshalJSON(bytes []byte) error {
	type alias OperationResource
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into OperationResource: %+v", err)
	}

	s.EndTime = decoded.EndTime
	s.Error = decoded.Error
	s.Id = decoded.Id
	s.Name = decoded.Name
	s.StartTime = decoded.StartTime
	s.Status = decoded.Status

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling OperationResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["properties"]; ok {
		impl, err := unmarshalOperationExtendedInfoImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Properties' for 'OperationResource': %+v", err)
		}
		s.Properties = impl
	}
	return nil
}
