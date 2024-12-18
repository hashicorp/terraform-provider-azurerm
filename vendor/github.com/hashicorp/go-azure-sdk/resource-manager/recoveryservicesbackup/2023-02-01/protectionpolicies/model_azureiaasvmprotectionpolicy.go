package protectionpolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ProtectionPolicy = AzureIaaSVMProtectionPolicy{}

type AzureIaaSVMProtectionPolicy struct {
	InstantRPDetails              *InstantRPAdditionalDetails `json:"instantRPDetails,omitempty"`
	InstantRpRetentionRangeInDays *int64                      `json:"instantRpRetentionRangeInDays,omitempty"`
	PolicyType                    *IAASVMPolicyType           `json:"policyType,omitempty"`
	RetentionPolicy               RetentionPolicy             `json:"retentionPolicy"`
	SchedulePolicy                SchedulePolicy              `json:"schedulePolicy"`
	TieringPolicy                 *map[string]TieringPolicy   `json:"tieringPolicy,omitempty"`
	TimeZone                      *string                     `json:"timeZone,omitempty"`

	// Fields inherited from ProtectionPolicy

	BackupManagementType           string    `json:"backupManagementType"`
	ProtectedItemsCount            *int64    `json:"protectedItemsCount,omitempty"`
	ResourceGuardOperationRequests *[]string `json:"resourceGuardOperationRequests,omitempty"`
}

func (s AzureIaaSVMProtectionPolicy) ProtectionPolicy() BaseProtectionPolicyImpl {
	return BaseProtectionPolicyImpl{
		BackupManagementType:           s.BackupManagementType,
		ProtectedItemsCount:            s.ProtectedItemsCount,
		ResourceGuardOperationRequests: s.ResourceGuardOperationRequests,
	}
}

var _ json.Marshaler = AzureIaaSVMProtectionPolicy{}

func (s AzureIaaSVMProtectionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper AzureIaaSVMProtectionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureIaaSVMProtectionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureIaaSVMProtectionPolicy: %+v", err)
	}

	decoded["backupManagementType"] = "AzureIaasVM"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureIaaSVMProtectionPolicy: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureIaaSVMProtectionPolicy{}

func (s *AzureIaaSVMProtectionPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		InstantRPDetails               *InstantRPAdditionalDetails `json:"instantRPDetails,omitempty"`
		InstantRpRetentionRangeInDays  *int64                      `json:"instantRpRetentionRangeInDays,omitempty"`
		PolicyType                     *IAASVMPolicyType           `json:"policyType,omitempty"`
		TieringPolicy                  *map[string]TieringPolicy   `json:"tieringPolicy,omitempty"`
		TimeZone                       *string                     `json:"timeZone,omitempty"`
		BackupManagementType           string                      `json:"backupManagementType"`
		ProtectedItemsCount            *int64                      `json:"protectedItemsCount,omitempty"`
		ResourceGuardOperationRequests *[]string                   `json:"resourceGuardOperationRequests,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.InstantRPDetails = decoded.InstantRPDetails
	s.InstantRpRetentionRangeInDays = decoded.InstantRpRetentionRangeInDays
	s.PolicyType = decoded.PolicyType
	s.TieringPolicy = decoded.TieringPolicy
	s.TimeZone = decoded.TimeZone
	s.BackupManagementType = decoded.BackupManagementType
	s.ProtectedItemsCount = decoded.ProtectedItemsCount
	s.ResourceGuardOperationRequests = decoded.ResourceGuardOperationRequests

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureIaaSVMProtectionPolicy into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["retentionPolicy"]; ok {
		impl, err := UnmarshalRetentionPolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RetentionPolicy' for 'AzureIaaSVMProtectionPolicy': %+v", err)
		}
		s.RetentionPolicy = impl
	}

	if v, ok := temp["schedulePolicy"]; ok {
		impl, err := UnmarshalSchedulePolicyImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'SchedulePolicy' for 'AzureIaaSVMProtectionPolicy': %+v", err)
		}
		s.SchedulePolicy = impl
	}

	return nil
}
