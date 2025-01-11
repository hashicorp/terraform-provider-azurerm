package replicationprotecteditems

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateReplicationProtectedItemInputProperties struct {
	EnableRdpOnTargetOption        *string                                     `json:"enableRdpOnTargetOption,omitempty"`
	LicenseType                    *LicenseType                                `json:"licenseType,omitempty"`
	ProviderSpecificDetails        UpdateReplicationProtectedItemProviderInput `json:"providerSpecificDetails"`
	RecoveryAvailabilitySetId      *string                                     `json:"recoveryAvailabilitySetId,omitempty"`
	RecoveryAzureVMName            *string                                     `json:"recoveryAzureVMName,omitempty"`
	RecoveryAzureVMSize            *string                                     `json:"recoveryAzureVMSize,omitempty"`
	SelectedRecoveryAzureNetworkId *string                                     `json:"selectedRecoveryAzureNetworkId,omitempty"`
	SelectedSourceNicId            *string                                     `json:"selectedSourceNicId,omitempty"`
	SelectedTfoAzureNetworkId      *string                                     `json:"selectedTfoAzureNetworkId,omitempty"`
	VMNics                         *[]VMNicInputDetails                        `json:"vmNics,omitempty"`
}

var _ json.Unmarshaler = &UpdateReplicationProtectedItemInputProperties{}

func (s *UpdateReplicationProtectedItemInputProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		EnableRdpOnTargetOption        *string              `json:"enableRdpOnTargetOption,omitempty"`
		LicenseType                    *LicenseType         `json:"licenseType,omitempty"`
		RecoveryAvailabilitySetId      *string              `json:"recoveryAvailabilitySetId,omitempty"`
		RecoveryAzureVMName            *string              `json:"recoveryAzureVMName,omitempty"`
		RecoveryAzureVMSize            *string              `json:"recoveryAzureVMSize,omitempty"`
		SelectedRecoveryAzureNetworkId *string              `json:"selectedRecoveryAzureNetworkId,omitempty"`
		SelectedSourceNicId            *string              `json:"selectedSourceNicId,omitempty"`
		SelectedTfoAzureNetworkId      *string              `json:"selectedTfoAzureNetworkId,omitempty"`
		VMNics                         *[]VMNicInputDetails `json:"vmNics,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.EnableRdpOnTargetOption = decoded.EnableRdpOnTargetOption
	s.LicenseType = decoded.LicenseType
	s.RecoveryAvailabilitySetId = decoded.RecoveryAvailabilitySetId
	s.RecoveryAzureVMName = decoded.RecoveryAzureVMName
	s.RecoveryAzureVMSize = decoded.RecoveryAzureVMSize
	s.SelectedRecoveryAzureNetworkId = decoded.SelectedRecoveryAzureNetworkId
	s.SelectedSourceNicId = decoded.SelectedSourceNicId
	s.SelectedTfoAzureNetworkId = decoded.SelectedTfoAzureNetworkId
	s.VMNics = decoded.VMNics

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling UpdateReplicationProtectedItemInputProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["providerSpecificDetails"]; ok {
		impl, err := UnmarshalUpdateReplicationProtectedItemProviderInputImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ProviderSpecificDetails' for 'UpdateReplicationProtectedItemInputProperties': %+v", err)
		}
		s.ProviderSpecificDetails = impl
	}

	return nil
}
