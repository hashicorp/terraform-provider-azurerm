package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TokenFilter = PhoneticTokenFilter{}

type PhoneticTokenFilter struct {
	Encoder *PhoneticEncoder `json:"encoder,omitempty"`
	Replace *bool            `json:"replace,omitempty"`

	// Fields inherited from TokenFilter

	Name      string `json:"name"`
	OdataType string `json:"@odata.type"`
}

func (s PhoneticTokenFilter) TokenFilter() BaseTokenFilterImpl {
	return BaseTokenFilterImpl{
		Name:      s.Name,
		OdataType: s.OdataType,
	}
}

var _ json.Marshaler = PhoneticTokenFilter{}

func (s PhoneticTokenFilter) MarshalJSON() ([]byte, error) {
	type wrapper PhoneticTokenFilter
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PhoneticTokenFilter: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PhoneticTokenFilter: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.PhoneticTokenFilter"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PhoneticTokenFilter: %+v", err)
	}

	return encoded, nil
}
