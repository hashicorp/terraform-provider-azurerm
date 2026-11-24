package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ CognitiveServicesAccount = DefaultCognitiveServicesAccount{}

type DefaultCognitiveServicesAccount struct {

	// Fields inherited from CognitiveServicesAccount

	Description *string `json:"description,omitempty"`
	OdataType   string  `json:"@odata.type"`
}

func (s DefaultCognitiveServicesAccount) CognitiveServicesAccount() BaseCognitiveServicesAccountImpl {
	return BaseCognitiveServicesAccountImpl{
		Description: s.Description,
		OdataType:   s.OdataType,
	}
}

var _ json.Marshaler = DefaultCognitiveServicesAccount{}

func (s DefaultCognitiveServicesAccount) MarshalJSON() ([]byte, error) {
	type wrapper DefaultCognitiveServicesAccount
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DefaultCognitiveServicesAccount: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DefaultCognitiveServicesAccount: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Azure.Search.DefaultCognitiveServices"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DefaultCognitiveServicesAccount: %+v", err)
	}

	return encoded, nil
}
