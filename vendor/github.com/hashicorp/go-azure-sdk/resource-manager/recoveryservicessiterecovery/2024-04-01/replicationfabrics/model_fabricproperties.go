package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FabricProperties struct {
	BcdrState                 *string               `json:"bcdrState,omitempty"`
	CustomDetails             FabricSpecificDetails `json:"customDetails"`
	EncryptionDetails         *EncryptionDetails    `json:"encryptionDetails,omitempty"`
	FriendlyName              *string               `json:"friendlyName,omitempty"`
	Health                    *string               `json:"health,omitempty"`
	HealthErrorDetails        *[]HealthError        `json:"healthErrorDetails,omitempty"`
	InternalIdentifier        *string               `json:"internalIdentifier,omitempty"`
	RolloverEncryptionDetails *EncryptionDetails    `json:"rolloverEncryptionDetails,omitempty"`
}

var _ json.Unmarshaler = &FabricProperties{}

func (s *FabricProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		BcdrState                 *string            `json:"bcdrState,omitempty"`
		EncryptionDetails         *EncryptionDetails `json:"encryptionDetails,omitempty"`
		FriendlyName              *string            `json:"friendlyName,omitempty"`
		Health                    *string            `json:"health,omitempty"`
		HealthErrorDetails        *[]HealthError     `json:"healthErrorDetails,omitempty"`
		InternalIdentifier        *string            `json:"internalIdentifier,omitempty"`
		RolloverEncryptionDetails *EncryptionDetails `json:"rolloverEncryptionDetails,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.BcdrState = decoded.BcdrState
	s.EncryptionDetails = decoded.EncryptionDetails
	s.FriendlyName = decoded.FriendlyName
	s.Health = decoded.Health
	s.HealthErrorDetails = decoded.HealthErrorDetails
	s.InternalIdentifier = decoded.InternalIdentifier
	s.RolloverEncryptionDetails = decoded.RolloverEncryptionDetails

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling FabricProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customDetails"]; ok {
		impl, err := UnmarshalFabricSpecificDetailsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomDetails' for 'FabricProperties': %+v", err)
		}
		s.CustomDetails = impl
	}

	return nil
}
