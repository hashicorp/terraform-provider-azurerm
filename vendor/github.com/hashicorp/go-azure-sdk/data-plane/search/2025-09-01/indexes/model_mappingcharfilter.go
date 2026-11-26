package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CharFilter = MappingCharFilter{}

type MappingCharFilter struct {
	Mappings []string `json:"mappings"`

	// Fields inherited from CharFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s MappingCharFilter) CharFilter() BaseCharFilterImpl {
	return BaseCharFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = MappingCharFilter{}

func (s MappingCharFilter) MarshalJSON() ([]byte, error) {
	type wrapper MappingCharFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MappingCharFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MappingCharFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.MappingCharFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MappingCharFilter: %+v", err)
	}

	return encoded, nil
}
