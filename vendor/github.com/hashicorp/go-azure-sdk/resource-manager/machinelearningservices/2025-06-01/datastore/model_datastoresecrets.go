package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreSecrets interface {
	DatastoreSecrets() BaseDatastoreSecretsImpl
}

var _ DatastoreSecrets = BaseDatastoreSecretsImpl{}

type BaseDatastoreSecretsImpl struct {
	SecretsType SecretsType `json:"secretsType"`
}

func (s BaseDatastoreSecretsImpl) DatastoreSecrets() BaseDatastoreSecretsImpl {
	return s
}

var _ DatastoreSecrets = RawDatastoreSecretsImpl{}

// RawDatastoreSecretsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDatastoreSecretsImpl struct {
	datastoreSecrets BaseDatastoreSecretsImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawDatastoreSecretsImpl) DatastoreSecrets() BaseDatastoreSecretsImpl {
	return s.datastoreSecrets
}

func UnmarshalDatastoreSecretsImplementation(input []byte) (DatastoreSecrets, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DatastoreSecrets into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["secretsType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AccountKey") {
		var out AccountKeyDatastoreSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AccountKeyDatastoreSecrets: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Certificate") {
		var out CertificateDatastoreSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CertificateDatastoreSecrets: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Sas") {
		var out SasDatastoreSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SasDatastoreSecrets: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServicePrincipal") {
		var out ServicePrincipalDatastoreSecrets
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePrincipalDatastoreSecrets: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatastoreSecretsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatastoreSecretsImpl: %+v", err)
	}

	return RawDatastoreSecretsImpl{
		datastoreSecrets: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
