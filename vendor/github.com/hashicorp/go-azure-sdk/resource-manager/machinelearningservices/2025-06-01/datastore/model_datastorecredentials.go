package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreCredentials interface {
	DatastoreCredentials() BaseDatastoreCredentialsImpl
}

var _ DatastoreCredentials = BaseDatastoreCredentialsImpl{}

type BaseDatastoreCredentialsImpl struct {
	CredentialsType CredentialsType `json:"credentialsType"`
}

func (s BaseDatastoreCredentialsImpl) DatastoreCredentials() BaseDatastoreCredentialsImpl {
	return s
}

var _ DatastoreCredentials = RawDatastoreCredentialsImpl{}

// RawDatastoreCredentialsImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawDatastoreCredentialsImpl struct {
	datastoreCredentials BaseDatastoreCredentialsImpl
	Type                 string
	Values               map[string]interface{}
}

func (s RawDatastoreCredentialsImpl) DatastoreCredentials() BaseDatastoreCredentialsImpl {
	return s.datastoreCredentials
}

func (s RawDatastoreCredentialsImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalDatastoreCredentialsImplementation(input []byte) (DatastoreCredentials, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DatastoreCredentials into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["credentialsType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AccountKey") {
		var out AccountKeyDatastoreCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AccountKeyDatastoreCredentials: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Certificate") {
		var out CertificateDatastoreCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CertificateDatastoreCredentials: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "None") {
		var out NoneDatastoreCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into NoneDatastoreCredentials: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Sas") {
		var out SasDatastoreCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SasDatastoreCredentials: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ServicePrincipal") {
		var out ServicePrincipalDatastoreCredentials
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServicePrincipalDatastoreCredentials: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatastoreCredentialsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatastoreCredentialsImpl: %+v", err)
	}

	return RawDatastoreCredentialsImpl{
		datastoreCredentials: parent,
		Type:                 value,
		Values:               temp,
	}, nil

}
