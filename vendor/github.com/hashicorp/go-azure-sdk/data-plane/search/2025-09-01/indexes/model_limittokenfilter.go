package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = LimitTokenFilter{}

type LimitTokenFilter struct {
	ConsumeAllTokens *bool  `json:"consumeAllTokens,omitempty"`
	MaxTokenCount    *int64 `json:"maxTokenCount,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s LimitTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = LimitTokenFilter{}

func (s LimitTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper LimitTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LimitTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LimitTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.LimitTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LimitTokenFilter: %+v", err)
	}

	return encoded, nil
}
