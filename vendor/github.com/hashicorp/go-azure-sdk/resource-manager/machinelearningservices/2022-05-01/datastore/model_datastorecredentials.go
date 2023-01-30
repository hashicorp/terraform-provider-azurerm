package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatastoreCredentials interface {
}

func unmarshalDatastoreCredentialsImplementation(input []byte) (DatastoreCredentials, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DatastoreCredentials into map[string]interface: %+v", err)
	}

	value, ok := temp["credentialsType"].(string)
	if !ok {
		return nil, nil
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

	type RawDatastoreCredentialsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDatastoreCredentialsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
