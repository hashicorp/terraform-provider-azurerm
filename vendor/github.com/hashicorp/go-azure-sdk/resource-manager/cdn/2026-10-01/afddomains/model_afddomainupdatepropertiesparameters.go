package afddomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDDomainUpdatePropertiesParameters struct {
	AzureDnsZone                       *ResourceReference        `json:"azureDnsZone,omitempty"`
	MtlsSettings                       AFDDomainMtlsParameters   `json:"mtlsSettings"`
	PreValidatedCustomDomainResourceId *ResourceReference        `json:"preValidatedCustomDomainResourceId,omitempty"`
	ProfileName                        *string                   `json:"profileName,omitempty"`
	TlsSettings                        *AFDDomainHTTPSParameters `json:"tlsSettings,omitempty"`
}

var _ json.Unmarshaler = &AFDDomainUpdatePropertiesParameters{}

func (s *AFDDomainUpdatePropertiesParameters) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AzureDnsZone                       *ResourceReference        `json:"azureDnsZone,omitempty"`
		PreValidatedCustomDomainResourceId *ResourceReference        `json:"preValidatedCustomDomainResourceId,omitempty"`
		ProfileName                        *string                   `json:"profileName,omitempty"`
		TlsSettings                        *AFDDomainHTTPSParameters `json:"tlsSettings,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AzureDnsZone = decoded.AzureDnsZone
	s.PreValidatedCustomDomainResourceId = decoded.PreValidatedCustomDomainResourceId
	s.ProfileName = decoded.ProfileName
	s.TlsSettings = decoded.TlsSettings

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AFDDomainUpdatePropertiesParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["mtlsSettings"]; ok {
		impl, err := UnmarshalAFDDomainMtlsParametersImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'MtlsSettings' for 'AFDDomainUpdatePropertiesParameters': %+v", err)
		}
		s.MtlsSettings = impl
	}

	return nil
}
