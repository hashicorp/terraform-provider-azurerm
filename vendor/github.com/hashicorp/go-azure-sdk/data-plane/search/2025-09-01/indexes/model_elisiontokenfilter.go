package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = ElisionTokenFilter{}

type ElisionTokenFilter struct {
	Articles *[]string `json:"articles,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s ElisionTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = ElisionTokenFilter{}

func (s ElisionTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper ElisionTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ElisionTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ElisionTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.ElisionTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ElisionTokenFilter: %+v", err)
	}

	return encoded, nil
}
