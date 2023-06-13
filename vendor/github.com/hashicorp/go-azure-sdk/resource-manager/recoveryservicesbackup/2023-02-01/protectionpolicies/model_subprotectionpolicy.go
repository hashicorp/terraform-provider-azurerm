package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubProtectionPolicy struct {
	PolicyType      *PolicyType               `json:"policyType,omitempty"`
	RetentionPolicy RetentionPolicy           `json:"retentionPolicy"`
	SchedulePolicy  SchedulePolicy            `json:"schedulePolicy"`
	TieringPolicy   *map[string]TieringPolicy `json:"tieringPolicy,omitempty"`
}

var _ json.Unmarshaler = &SubProtectionPolicy{}

func (s *SubProtectionPolicy) UnmarshalJSON(bytes []byte) error {
	type alias SubProtectionPolicy
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into SubProtectionPolicy: %+v", err)
	}

	s.PolicyType = decoded.PolicyType
	s.TieringPolicy = decoded.TieringPolicy

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SubProtectionPolicy into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["retentionPolicy"]; ok {
		impl, err := unmarshalRetentionPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RetentionPolicy' for 'SubProtectionPolicy': %+v", err)
		}
		s.RetentionPolicy = impl
	}

	if v, ok := temp["schedulePolicy"]; ok {
		impl, err := unmarshalSchedulePolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SchedulePolicy' for 'SubProtectionPolicy': %+v", err)
		}
		s.SchedulePolicy = impl
	}
	return nil
}
