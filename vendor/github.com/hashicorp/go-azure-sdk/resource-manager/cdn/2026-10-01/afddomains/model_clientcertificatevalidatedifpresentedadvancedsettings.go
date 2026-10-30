package afddomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AFDDomainMtlsParameters = ClientCertificateValidatedIfPresentedAdvancedSettings{}

type ClientCertificateValidatedIfPresentedAdvancedSettings struct {
	AllowedFqdns               *[]string                               `json:"allowedFqdns,omitempty"`
	CertificateRevocationCheck *CertificateRevocationCheckEnabledState `json:"certificateRevocationCheck,omitempty"`
	Secrets                    []ResourceReference                     `json:"secrets"`

	// Fields inherited from AFDDomainMtlsParameters

	Scenario MtlsScenarioType `json:"scenario"`
}

func (s ClientCertificateValidatedIfPresentedAdvancedSettings) AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl {
	return BaseAFDDomainMtlsParametersImpl{
		Scenario: s.Scenario,
	}
}

var _ json.Marshaler = ClientCertificateValidatedIfPresentedAdvancedSettings{}

func (s ClientCertificateValidatedIfPresentedAdvancedSettings) MarshalJSON() ([]byte, error) {
	type wrapper ClientCertificateValidatedIfPresentedAdvancedSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClientCertificateValidatedIfPresentedAdvancedSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClientCertificateValidatedIfPresentedAdvancedSettings: %+v", err)
	}

	decoded["scenario"] = "ClientCertificateValidatedIfPresented"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClientCertificateValidatedIfPresentedAdvancedSettings: %+v", err)
	}

	return encoded, nil
}
