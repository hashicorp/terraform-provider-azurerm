package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CognitiveServicesAccount = CognitiveServicesAccountKey{}

type CognitiveServicesAccountKey struct {
	Key string `json:"key"`

	// Fields inherited from CognitiveServicesAccount

	Description *string `json:"description,omitempty"`
	OdataType   string  `json:"@odata.type"`
}

func (s CognitiveServicesAccountKey) CognitiveServicesAccount() BaseCognitiveServicesAccountImpl {
	return BaseCognitiveServicesAccountImpl{
		Description: s.Description,
		OdataType:   s.OdataType,
	}
}

var _ json.Marshaler = CognitiveServicesAccountKey{}

func (s CognitiveServicesAccountKey) MarshalJSON() ([]byte, error) {
	type wrapper CognitiveServicesAccountKey
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CognitiveServicesAccountKey: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CognitiveServicesAccountKey: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.CognitiveServicesByKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CognitiveServicesAccountKey: %+v", err)
	}

	return encoded, nil
}
