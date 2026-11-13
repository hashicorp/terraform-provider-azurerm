package afddomains

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ AFDDomainMtlsParameters = ClientCertificateRequiredAndOriginValidatesAdvancedSettings{}

type ClientCertificateRequiredAndOriginValidatesAdvancedSettings struct {

	// Fields inherited from AFDDomainMtlsParameters

	Scenario MtlsScenarioType `json:"scenario"`
}

func (s ClientCertificateRequiredAndOriginValidatesAdvancedSettings) AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl {
	return BaseAFDDomainMtlsParametersImpl{
		Scenario: s.Scenario,
	}
}

var _ json.Marshaler = ClientCertificateRequiredAndOriginValidatesAdvancedSettings{}

func (s ClientCertificateRequiredAndOriginValidatesAdvancedSettings) MarshalJSON() ([]byte, error) {
	type wrapper ClientCertificateRequiredAndOriginValidatesAdvancedSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ClientCertificateRequiredAndOriginValidatesAdvancedSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ClientCertificateRequiredAndOriginValidatesAdvancedSettings: %+v", err)
	}

	decoded["scenario"] = "ClientCertificateRequiredAndOriginValidates"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ClientCertificateRequiredAndOriginValidatesAdvancedSettings: %+v", err)
	}

	return encoded, nil
}
