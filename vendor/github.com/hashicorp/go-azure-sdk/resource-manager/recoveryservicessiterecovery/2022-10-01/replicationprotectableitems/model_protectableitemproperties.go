package replicationprotectableitems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProtectableItemProperties struct {
	CustomDetails                 ConfigurationSettings `json:"customDetails"`
	FriendlyName                  *string               `json:"friendlyName,omitempty"`
	ProtectionReadinessErrors     *[]string             `json:"protectionReadinessErrors,omitempty"`
	ProtectionStatus              *string               `json:"protectionStatus,omitempty"`
	RecoveryServicesProviderId    *string               `json:"recoveryServicesProviderId,omitempty"`
	ReplicationProtectedItemId    *string               `json:"replicationProtectedItemId,omitempty"`
	SupportedReplicationProviders *[]string             `json:"supportedReplicationProviders,omitempty"`
}

var _ json.Unmarshaler = &ProtectableItemProperties{}

func (s *ProtectableItemProperties) UnmarshalJSON(bytes []byte) error {
	type alias ProtectableItemProperties
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into ProtectableItemProperties: %+v", err)
	}

	s.FriendlyName = decoded.FriendlyName
	s.ProtectionReadinessErrors = decoded.ProtectionReadinessErrors
	s.ProtectionStatus = decoded.ProtectionStatus
	s.RecoveryServicesProviderId = decoded.RecoveryServicesProviderId
	s.ReplicationProtectedItemId = decoded.ReplicationProtectedItemId
	s.SupportedReplicationProviders = decoded.SupportedReplicationProviders

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling ProtectableItemProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["customDetails"]; ok {
		impl, err := unmarshalConfigurationSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CustomDetails' for 'ProtectableItemProperties': %+v", err)
		}
		s.CustomDetails = impl
	}
	return nil
}
