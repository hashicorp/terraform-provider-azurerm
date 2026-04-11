package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = KeywordMarkerTokenFilter{}

type KeywordMarkerTokenFilter struct {
	IgnoreCase *bool    `json:"ignoreCase,omitempty"`
	Keywords   []string `json:"keywords"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s KeywordMarkerTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = KeywordMarkerTokenFilter{}

func (s KeywordMarkerTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper KeywordMarkerTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeywordMarkerTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeywordMarkerTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.KeywordMarkerTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeywordMarkerTokenFilter: %+v", err)
	}

	return encoded, nil
}
