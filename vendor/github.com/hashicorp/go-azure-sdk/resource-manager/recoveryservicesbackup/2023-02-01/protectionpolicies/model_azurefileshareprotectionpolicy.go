package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionPolicy = AzureFileShareProtectionPolicy{}

type AzureFileShareProtectionPolicy struct {
	RetentionPolicy RetentionPolicy `json:"retentionPolicy"`
	SchedulePolicy  SchedulePolicy  `json:"schedulePolicy"`
	TimeZone        *string         `json:"timeZone,omitempty"`
	WorkLoadType    *WorkloadType   `json:"workLoadType,omitempty"`

	// Fields inherited from ProtectionPolicy
	ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
	ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
}

var _ json.Marshaler = AzureFileShareProtectionPolicy{}

func (s AzureFileShareProtectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper AzureFileShareProtectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureFileShareProtectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureFileShareProtectionPolicy: %+v", err)
	}
	decoded["backupManagementType"] = "AzureStorage"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureFileShareProtectionPolicy: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureFileShareProtectionPolicy{}

func (s *AzureFileShareProtectionPolicy) UnmarshalJSON(bytes []byte) error {
	type alias AzureFileShareProtectionPolicy
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AzureFileShareProtectionPolicy: %+v", err)
	}

	s.ProtectedItemsCount = decoded.ProtectedItemsCount
	s.ResourceGuardOperationRequests = decoded.ResourceGuardOperationRequests
	s.TimeZone = decoded.TimeZone
	s.WorkLoadType = decoded.WorkLoadType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureFileShareProtectionPolicy into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["retentionPolicy"]; ok {
		impl, err := unmarshalRetentionPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RetentionPolicy' for 'AzureFileShareProtectionPolicy': %+v", err)
		}
		s.RetentionPolicy = impl
	}

	if v, ok := temp["schedulePolicy"]; ok {
		impl, err := unmarshalSchedulePolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SchedulePolicy' for 'AzureFileShareProtectionPolicy': %+v", err)
		}
		s.SchedulePolicy = impl
	}
	return nil
}
