package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Datastore interface {
}

func unmarshalDatastoreImplementation(input []byte) (Datastore, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Datastore into map[string]interface: %+v", err)
	}

	value, ok := temp["datastoreType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureBlob") {
		var out AzureBlobDatastore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureBlobDatastore: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeGen1") {
		var out AzureDataLakeGen1Datastore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeGen1Datastore: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureDataLakeGen2") {
		var out AzureDataLakeGen2Datastore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDataLakeGen2Datastore: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzureFile") {
		var out AzureFileDatastore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureFileDatastore: %+v", err)
		}
		return out, nil
	}

	type RawDatastoreImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDatastoreImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
