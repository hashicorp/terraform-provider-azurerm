package secrets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecretParameters interface {
	SecretParameters() BaseSecretParametersImpl
}

var _ SecretParameters = BaseSecretParametersImpl{}

type BaseSecretParametersImpl struct {
	Type SecretType `json:"type"`
}

func (s BaseSecretParametersImpl) SecretParameters() BaseSecretParametersImpl {
	return s
}

var _ SecretParameters = RawSecretParametersImpl{}

// RawSecretParametersImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSecretParametersImpl struct {
	secretParameters BaseSecretParametersImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawSecretParametersImpl) SecretParameters() BaseSecretParametersImpl {
	return s.secretParameters
}

func (s RawSecretParametersImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSecretParametersImplementation(input []byte) (SecretParameters, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SecretParameters into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureFirstPartyManagedCertificate") {
		var out AzureFirstPartyManagedCertificateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFirstPartyManagedCertificateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomerCertificate") {
		var out CustomerCertificateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomerCertificateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ManagedCertificate") {
		var out ManagedCertificateParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ManagedCertificateParameters: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "UrlSigningKey") {
		var out URLSigningKeyParameters
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into URLSigningKeyParameters: %+v", err)
		}
		return out, nil
	}

	var parent BaseSecretParametersImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSecretParametersImpl: %+v", err)
	}

	return RawSecretParametersImpl{
		secretParameters: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
