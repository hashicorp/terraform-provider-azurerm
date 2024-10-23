package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Datastore interface {
	Datastore() BaseDatastoreImpl
}

var _ Datastore = BaseDatastoreImpl{}

type BaseDatastoreImpl struct {
	Credentials   DatastoreCredentials `json:"credentials"`
	DatastoreType DatastoreType        `json:"datastoreType"`
	Description   *string              `json:"description,omitempty"`
	IsDefault     *bool                `json:"isDefault,omitempty"`
	Properties    *map[string]string   `json:"properties,omitempty"`
	Tags          *map[string]string   `json:"tags,omitempty"`
}

func (s BaseDatastoreImpl) Datastore() BaseDatastoreImpl {
	return s
}

var _ Datastore = RawDatastoreImpl{}

// RawDatastoreImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDatastoreImpl struct {
	datastore BaseDatastoreImpl
	Type      string
	Values    map[string]interface{}
}

func (s RawDatastoreImpl) Datastore() BaseDatastoreImpl {
	return s.datastore
}

var _ json.Unmarshaler = &BaseDatastoreImpl{}

func (s *BaseDatastoreImpl) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		DatastoreType DatastoreType      `json:"datastoreType"`
		Description   *string            `json:"description,omitempty"`
		IsDefault     *bool              `json:"isDefault,omitempty"`
		Properties    *map[string]string `json:"properties,omitempty"`
		Tags          *map[string]string `json:"tags,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.DatastoreType = decoded.DatastoreType
	s.Description = decoded.Description
	s.IsDefault = decoded.IsDefault
	s.Properties = decoded.Properties
	s.Tags = decoded.Tags

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BaseDatastoreImpl into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["credentials"]; ok {
		impl, err := UnmarshalDatastoreCredentialsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Credentials' for 'BaseDatastoreImpl': %+v", err)
		}
		s.Credentials = impl
	}

	return nil
}

func UnmarshalDatastoreImplementation(input []byte) (Datastore, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Datastore into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["datastoreType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	if strings.EqualFold(value, "OneLake") {
		var out OneLakeDatastore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OneLakeDatastore: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatastoreImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatastoreImpl: %+v", err)
	}

	return RawDatastoreImpl{
		datastore: parent,
		Type:      value,
		Values:    temp,
	}, nil

}
