package skillsets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CognitiveServicesAccount interface {
	CognitiveServicesAccount() BaseCognitiveServicesAccountImpl
}

var _ CognitiveServicesAccount = BaseCognitiveServicesAccountImpl{}

type BaseCognitiveServicesAccountImpl struct {
	Description *string `json:"description,omitempty"`
	OdataType   string  `json:"@odata.type"`
}

func (s BaseCognitiveServicesAccountImpl) CognitiveServicesAccount() BaseCognitiveServicesAccountImpl {
	return s
}

var _ CognitiveServicesAccount = RawCognitiveServicesAccountImpl{}

// RawCognitiveServicesAccountImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawCognitiveServicesAccountImpl struct {
	cognitiveServicesAccount BaseCognitiveServicesAccountImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawCognitiveServicesAccountImpl) CognitiveServicesAccount() BaseCognitiveServicesAccountImpl {
	return s.cognitiveServicesAccount
}

func (s RawCognitiveServicesAccountImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalCognitiveServicesAccountImplementation(input []byte) (CognitiveServicesAccount, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CognitiveServicesAccount into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.CognitiveServicesByKey") {
		var out CognitiveServicesAccountKey
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CognitiveServicesAccountKey: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.DefaultCognitiveServices") {
		var out DefaultCognitiveServicesAccount
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DefaultCognitiveServicesAccount: %+v", err)
		}
		return out, nil
	}

	var parent BaseCognitiveServicesAccountImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCognitiveServicesAccountImpl: %+v", err)
	}

	return RawCognitiveServicesAccountImpl{
		cognitiveServicesAccount: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
