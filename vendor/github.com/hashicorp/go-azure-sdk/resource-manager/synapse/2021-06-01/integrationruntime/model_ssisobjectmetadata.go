package integrationruntime

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SsisObjectMetadata interface {
	SsisObjectMetadata() BaseSsisObjectMetadataImpl
}

var _ SsisObjectMetadata = BaseSsisObjectMetadataImpl{}

type BaseSsisObjectMetadataImpl struct {
	Description *string                `json:"description,omitempty"`
	Id          *int64                 `json:"id,omitempty"`
	Name        *string                `json:"name,omitempty"`
	Type        SsisObjectMetadataType `json:"type"`
}

func (s BaseSsisObjectMetadataImpl) SsisObjectMetadata() BaseSsisObjectMetadataImpl {
	return s
}

var _ SsisObjectMetadata = RawSsisObjectMetadataImpl{}

// RawSsisObjectMetadataImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSsisObjectMetadataImpl struct {
	ssisObjectMetadata BaseSsisObjectMetadataImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawSsisObjectMetadataImpl) SsisObjectMetadata() BaseSsisObjectMetadataImpl {
	return s.ssisObjectMetadata
}

func (s RawSsisObjectMetadataImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSsisObjectMetadataImplementation(input []byte) (SsisObjectMetadata, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SsisObjectMetadata into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Environment") {
		var out SsisEnvironment
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SsisEnvironment: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Folder") {
		var out SsisFolder
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SsisFolder: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Package") {
		var out SsisPackage
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SsisPackage: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Project") {
		var out SsisProject
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SsisProject: %+v", err)
		}
		return out, nil
	}

	var parent BaseSsisObjectMetadataImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSsisObjectMetadataImpl: %+v", err)
	}

	return RawSsisObjectMetadataImpl{
		ssisObjectMetadata: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
