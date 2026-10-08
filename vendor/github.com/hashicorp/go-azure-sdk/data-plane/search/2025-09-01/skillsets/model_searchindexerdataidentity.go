package skillsets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerDataIdentity interface {
	SearchIndexerDataIdentity() BaseSearchIndexerDataIdentityImpl
}

var _ SearchIndexerDataIdentity = BaseSearchIndexerDataIdentityImpl{}

type BaseSearchIndexerDataIdentityImpl struct {
	OdataType string `json:"@odata.type"`
}

func (s BaseSearchIndexerDataIdentityImpl) SearchIndexerDataIdentity() BaseSearchIndexerDataIdentityImpl {
	return s
}

var _ SearchIndexerDataIdentity = RawSearchIndexerDataIdentityImpl{}

// RawSearchIndexerDataIdentityImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSearchIndexerDataIdentityImpl struct {
	searchIndexerDataIdentity BaseSearchIndexerDataIdentityImpl
	Type                      string
	Values                    map[string]interface{}
}

func (s RawSearchIndexerDataIdentityImpl) SearchIndexerDataIdentity() BaseSearchIndexerDataIdentityImpl {
	return s.searchIndexerDataIdentity
}

func (s RawSearchIndexerDataIdentityImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSearchIndexerDataIdentityImplementation(input []byte) (SearchIndexerDataIdentity, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SearchIndexerDataIdentity into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.DataNoneIdentity") {
		var out SearchIndexerDataNoneIdentity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SearchIndexerDataNoneIdentity: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Azure.Search.DataUserAssignedIdentity") {
		var out SearchIndexerDataUserAssignedIdentity
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SearchIndexerDataUserAssignedIdentity: %+v", err)
		}
		return out, nil
	}

	var parent BaseSearchIndexerDataIdentityImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSearchIndexerDataIdentityImpl: %+v", err)
	}

	return RawSearchIndexerDataIdentityImpl{
		searchIndexerDataIdentity: parent,
		Type:                      value,
		Values:                    temp,
	}, nil

}
