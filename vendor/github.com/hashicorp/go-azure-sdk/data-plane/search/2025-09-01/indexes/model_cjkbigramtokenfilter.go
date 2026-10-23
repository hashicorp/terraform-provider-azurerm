package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = CjkBigramTokenFilter{}

type CjkBigramTokenFilter struct {
	IgnoreScripts  *[]CjkBigramTokenFilterScripts `json:"ignoreScripts,omitempty"`
	OutputUnigrams *bool                          `json:"outputUnigrams,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s CjkBigramTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = CjkBigramTokenFilter{}

func (s CjkBigramTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper CjkBigramTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CjkBigramTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CjkBigramTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.CjkBigramTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CjkBigramTokenFilter: %+v", err)
	}

	return encoded, nil
}
