package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureBackupRestoreRequest interface {
	AzureBackupRestoreRequest() BaseAzureBackupRestoreRequestImpl
}

var _ AzureBackupRestoreRequest = BaseAzureBackupRestoreRequestImpl{}

type BaseAzureBackupRestoreRequestImpl struct {
	IdentityDetails                *IdentityDetails      `json:"identityDetails,omitempty"`
	ObjectType                     string                `json:"objectType"`
	ResourceGuardOperationRequests *[]string             `json:"resourceGuardOperationRequests,omitempty"`
	RestoreTargetInfo              RestoreTargetInfoBase `json:"restoreTargetInfo"`
	SourceDataStoreType            SourceDataStoreType   `json:"sourceDataStoreType"`
	SourceResourceId               *string               `json:"sourceResourceId,omitempty"`
}

func (s BaseAzureBackupRestoreRequestImpl) AzureBackupRestoreRequest() BaseAzureBackupRestoreRequestImpl {
	return s
}

var _ AzureBackupRestoreRequest = RawAzureBackupRestoreRequestImpl{}

// RawAzureBackupRestoreRequestImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAzureBackupRestoreRequestImpl struct {
	azureBackupRestoreRequest BaseAzureBackupRestoreRequestImpl
	Type                      string
	Values                    map[string]interface{}
}

func (s RawAzureBackupRestoreRequestImpl) AzureBackupRestoreRequest() BaseAzureBackupRestoreRequestImpl {
	return s.azureBackupRestoreRequest
}

var _ json.Unmarshaler = &BaseAzureBackupRestoreRequestImpl{}

func (s *BaseAzureBackupRestoreRequestImpl) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		IdentityDetails                *IdentityDetails    `json:"identityDetails,omitempty"`
		ObjectType                     string              `json:"objectType"`
		ResourceGuardOperationRequests *[]string           `json:"resourceGuardOperationRequests,omitempty"`
		SourceDataStoreType            SourceDataStoreType `json:"sourceDataStoreType"`
		SourceResourceId               *string             `json:"sourceResourceId,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.IdentityDetails = decoded.IdentityDetails
	s.ObjectType = decoded.ObjectType
	s.ResourceGuardOperationRequests = decoded.ResourceGuardOperationRequests
	s.SourceDataStoreType = decoded.SourceDataStoreType
	s.SourceResourceId = decoded.SourceResourceId

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BaseAzureBackupRestoreRequestImpl into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["restoreTargetInfo"]; ok {
		impl, err := UnmarshalRestoreTargetInfoBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'RestoreTargetInfo' for 'BaseAzureBackupRestoreRequestImpl': %+v", err)
		}
		s.RestoreTargetInfo = impl
	}

	return nil
}

func UnmarshalAzureBackupRestoreRequestImplementation(input []byte) (AzureBackupRestoreRequest, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureBackupRestoreRequest into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureBackupRecoveryPointBasedRestoreRequest") {
		var out AzureBackupRecoveryPointBasedRestoreRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRecoveryPointBasedRestoreRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBackupRecoveryTimeBasedRestoreRequest") {
		var out AzureBackupRecoveryTimeBasedRestoreRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRecoveryTimeBasedRestoreRequest: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureBackupRestoreWithRehydrationRequest") {
		var out AzureBackupRestoreWithRehydrationRequest
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBackupRestoreWithRehydrationRequest: %+v", err)
		}
		return out, nil
	}

	var parent BaseAzureBackupRestoreRequestImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAzureBackupRestoreRequestImpl: %+v", err)
	}

	return RawAzureBackupRestoreRequestImpl{
		azureBackupRestoreRequest: parent,
		Type:                      value,
		Values:                    temp,
	}, nil

}
