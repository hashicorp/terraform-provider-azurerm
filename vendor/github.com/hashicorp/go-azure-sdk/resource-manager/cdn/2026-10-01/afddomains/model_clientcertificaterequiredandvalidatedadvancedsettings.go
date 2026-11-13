package afddomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AFDDomainMtlsParameters = ClientCertificateRequiredAndValidatedAdvancedSettings{}

type ClientCertificateRequiredAndValidatedAdvancedSettings struct {
	AllowedFqdns               *[]string                               `json:"allowedFqdns,omitempty"`
	CertificateRevocationCheck *CertificateRevocationCheckEnabledState `json:"certificateRevocationCheck,omitempty"`
	Secrets                    []ResourceReference                     `json:"secrets"`

	// Fields inherited from AFDDomainMtlsParameters

	Scenario MtlsScenarioType `json:"scenario"`
}

func (s ClientCertificateRequiredAndValidatedAdvancedSettings) AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl {
	return BaseAFDDomainMtlsParametersImpl{
		Scenario: s.Scenario,
	}
}

var _ json.Marshaler = ClientCertificateRequiredAndValidatedAdvancedSettings{}

func (s ClientCertificateRequiredAndValidatedAdvancedSettings) MarshalJSON() ([]byte, error) {
	type wrapper ClientCertificateRequiredAndValidatedAdvancedSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClientCertificateRequiredAndValidatedAdvancedSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClientCertificateRequiredAndValidatedAdvancedSettings: %+v", err)
	}

	decoded["scenario"] = "ClientCertificateRequiredAndValidated"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClientCertificateRequiredAndValidatedAdvancedSettings: %+v", err)
	}

	return encoded, nil
}
