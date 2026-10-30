package afddomains

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AFDDomainMtlsParameters interface {
	AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl
}

var _ AFDDomainMtlsParameters = BaseAFDDomainMtlsParametersImpl{}

type BaseAFDDomainMtlsParametersImpl struct {
	Scenario MtlsScenarioType `json:"scenario"`
}

func (s BaseAFDDomainMtlsParametersImpl) AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl {
	return s
}

var _ AFDDomainMtlsParameters = RawAFDDomainMtlsParametersImpl{}

// RawAFDDomainMtlsParametersImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawAFDDomainMtlsParametersImpl struct {
	aFDDomainMtlsParameters BaseAFDDomainMtlsParametersImpl
	Type                    string
	Values                  map[string]interface{}
}

func (s RawAFDDomainMtlsParametersImpl) AFDDomainMtlsParameters() BaseAFDDomainMtlsParametersImpl {
	return s.aFDDomainMtlsParameters
}

func (s RawAFDDomainMtlsParametersImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalAFDDomainMtlsParametersImplementation(input []byte) (AFDDomainMtlsParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AFDDomainMtlsParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["scenario"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ClientCertificateRequiredAndOriginValidates") {
		var out ClientCertificateRequiredAndOriginValidatesAdvancedSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClientCertificateRequiredAndOriginValidatesAdvancedSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ClientCertificateRequiredAndValidated") {
		var out ClientCertificateRequiredAndValidatedAdvancedSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClientCertificateRequiredAndValidatedAdvancedSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ClientCertificateValidatedIfPresented") {
		var out ClientCertificateValidatedIfPresentedAdvancedSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ClientCertificateValidatedIfPresentedAdvancedSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CompleteMtlsPassthroughToOrigin") {
		var out CompleteMtlsPassthroughToOriginAdvancedSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CompleteMtlsPassthroughToOriginAdvancedSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseAFDDomainMtlsParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseAFDDomainMtlsParametersImpl: %+v", err)
	}

	return RawAFDDomainMtlsParametersImpl{
		aFDDomainMtlsParameters: parent,
		Type:                    value,
		Values:                  temp,
	}, nil

}
