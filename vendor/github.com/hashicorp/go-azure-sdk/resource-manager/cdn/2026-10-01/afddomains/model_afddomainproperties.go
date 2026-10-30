package afddomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDDomainProperties struct {
	AzureDnsZone                       *ResourceReference          `json:"azureDnsZone,omitempty"`
	DeploymentStatus                   *DeploymentStatus           `json:"deploymentStatus,omitempty"`
	DomainValidationState              *DomainValidationState      `json:"domainValidationState,omitempty"`
	ExtendedProperties                 *map[string]string          `json:"extendedProperties,omitempty"`
	HostName                           string                      `json:"hostName"`
	MtlsSettings                       AFDDomainMtlsParameters     `json:"mtlsSettings"`
	PreValidatedCustomDomainResourceId *ResourceReference          `json:"preValidatedCustomDomainResourceId,omitempty"`
	ProfileName                        *string                     `json:"profileName,omitempty"`
	ProvisioningState                  *AfdProvisioningState       `json:"provisioningState,omitempty"`
	TlsSettings                        *AFDDomainHTTPSParameters   `json:"tlsSettings,omitempty"`
	ValidationProperties               *DomainValidationProperties `json:"validationProperties,omitempty"`
}

var _ json.Unmarshaler = &AFDDomainProperties{}

func (s *AFDDomainProperties) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AzureDnsZone                       *ResourceReference          `json:"azureDnsZone,omitempty"`
		DeploymentStatus                   *DeploymentStatus           `json:"deploymentStatus,omitempty"`
		DomainValidationState              *DomainValidationState      `json:"domainValidationState,omitempty"`
		ExtendedProperties                 *map[string]string          `json:"extendedProperties,omitempty"`
		HostName                           string                      `json:"hostName"`
		PreValidatedCustomDomainResourceId *ResourceReference          `json:"preValidatedCustomDomainResourceId,omitempty"`
		ProfileName                        *string                     `json:"profileName,omitempty"`
		ProvisioningState                  *AfdProvisioningState       `json:"provisioningState,omitempty"`
		TlsSettings                        *AFDDomainHTTPSParameters   `json:"tlsSettings,omitempty"`
		ValidationProperties               *DomainValidationProperties `json:"validationProperties,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AzureDnsZone = decoded.AzureDnsZone
	s.DeploymentStatus = decoded.DeploymentStatus
	s.DomainValidationState = decoded.DomainValidationState
	s.ExtendedProperties = decoded.ExtendedProperties
	s.HostName = decoded.HostName
	s.PreValidatedCustomDomainResourceId = decoded.PreValidatedCustomDomainResourceId
	s.ProfileName = decoded.ProfileName
	s.ProvisioningState = decoded.ProvisioningState
	s.TlsSettings = decoded.TlsSettings
	s.ValidationProperties = decoded.ValidationProperties

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AFDDomainProperties into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["mtlsSettings"]; ok {
		impl, err := UnmarshalAFDDomainMtlsParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'MtlsSettings' for 'AFDDomainProperties': %+v", err)
		}
		s.MtlsSettings = impl
	}

	return nil
}
