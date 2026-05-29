package projectresource

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProjectProperties struct {
	CreationTime         *string                   `json:"creationTime,omitempty"`
	DatabasesInfo        *[]DatabaseInfo           `json:"databasesInfo,omitempty"`
	ProvisioningState    *ProjectProvisioningState `json:"provisioningState,omitempty"`
	SourceConnectionInfo ConnectionInfo            `json:"sourceConnectionInfo"`
	SourcePlatform       ProjectSourcePlatform     `json:"sourcePlatform"`
	TargetConnectionInfo ConnectionInfo            `json:"targetConnectionInfo"`
	TargetPlatform       ProjectTargetPlatform     `json:"targetPlatform"`
}

func (o *ProjectProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ProjectProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

var _ json.Unmarshaler = &ProjectProperties{}

func (s *ProjectProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CreationTime      *string                   `json:"creationTime,omitempty"`
		DatabasesInfo     *[]DatabaseInfo           `json:"databasesInfo,omitempty"`
		ProvisioningState *ProjectProvisioningState `json:"provisioningState,omitempty"`
		SourcePlatform    ProjectSourcePlatform     `json:"sourcePlatform"`
		TargetPlatform    ProjectTargetPlatform     `json:"targetPlatform"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CreationTime = decoded.CreationTime
	s.DatabasesInfo = decoded.DatabasesInfo
	s.ProvisioningState = decoded.ProvisioningState
	s.SourcePlatform = decoded.SourcePlatform
	s.TargetPlatform = decoded.TargetPlatform

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ProjectProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["sourceConnectionInfo"]; ok {
		impl, err := UnmarshalConnectionInfoImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SourceConnectionInfo' for 'ProjectProperties': %+v", err)
		}
		s.SourceConnectionInfo = impl
	}

	if v, ok := temp["targetConnectionInfo"]; ok {
		impl, err := UnmarshalConnectionInfoImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TargetConnectionInfo' for 'ProjectProperties': %+v", err)
		}
		s.TargetConnectionInfo = impl
	}

	return nil
}
