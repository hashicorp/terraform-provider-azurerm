package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionPolicy = MabProtectionPolicy{}

type MabProtectionPolicy struct {
	RetentionPolicy RetentionPolicy `json:"retentionPolicy"`
	SchedulePolicy  SchedulePolicy  `json:"schedulePolicy"`

	// Fields inherited from ProtectionPolicy

	BackupManagementType           string    `json:"backupManagementType"`
	ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
	ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
}

func (s MabProtectionPolicy) ProtectionPolicy() BaseProtectionPolicyImpl {
	return BaseProtectionPolicyImpl{
		BackupManagementType:           s.BackupManagementType,
		ProtectedItemsCount:            s.ProtectedItemsCount,
		ResourceGuardOperationRequests: s.ResourceGuardOperationRequests,
	}
}

var _ json.Marshaler = MabProtectionPolicy{}

func (s MabProtectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper MabProtectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MabProtectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MabProtectionPolicy: %+v", err)
	}

	decoded["backupManagementType"] = "MAB"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MabProtectionPolicy: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &MabProtectionPolicy{}

func (s *MabProtectionPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		BackupManagementType           string    `json:"backupManagementType"`
		ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
		ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.BackupManagementType = decoded.BackupManagementType
	s.ProtectedItemsCount = decoded.ProtectedItemsCount
	s.ResourceGuardOperationRequests = decoded.ResourceGuardOperationRequests

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MabProtectionPolicy into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["retentionPolicy"]; ok {
		impl, err := UnmarshalRetentionPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RetentionPolicy' for 'MabProtectionPolicy': %+v", err)
		}
		s.RetentionPolicy = impl
	}

	if v, ok := temp["schedulePolicy"]; ok {
		impl, err := UnmarshalSchedulePolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SchedulePolicy' for 'MabProtectionPolicy': %+v", err)
		}
		s.SchedulePolicy = impl
	}

	return nil
}
