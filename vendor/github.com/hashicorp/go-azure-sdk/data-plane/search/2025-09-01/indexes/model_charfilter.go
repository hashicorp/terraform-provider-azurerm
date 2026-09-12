package indexes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CharFilter interface {
	CharFilter() BaseCharFilterImpl
}

var _ CharFilter = BaseCharFilterImpl{}

type BaseCharFilterImpl struct {
	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s BaseCharFilterImpl) CharFilter() BaseCharFilterImpl {
	return s
}

var _ CharFilter = RawCharFilterImpl{}

// RawCharFilterImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCharFilterImpl struct {
	charFilter BaseCharFilterImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawCharFilterImpl) CharFilter() BaseCharFilterImpl {
	return s.charFilter
}

func UnmarshalCharFilterImplementation(input []byte) (CharFilter, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CharFilter into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.MappingCharFilter") {
		var out MappingCharFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MappingCharFilter: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.PatternReplaceCharFilter") {
		var out PatternReplaceCharFilter
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PatternReplaceCharFilter: %+v", err)
		}
		return out, nil
	}

	var parent BaseCharFilterImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCharFilterImpl: %+v", err)
	}

	return RawCharFilterImpl{
		charFilter: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
