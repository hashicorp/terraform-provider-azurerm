package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyConfiguration = ContentKeyPolicyUnknownConfiguration{}

type ContentKeyPolicyUnknownConfiguration struct {

	// Fields inherited from ContentKeyPolicyConfiguration
}

var _ json.Marshaler = ContentKeyPolicyUnknownConfiguration{}

func (s ContentKeyPolicyUnknownConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyUnknownConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyUnknownConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyUnknownConfiguration: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyUnknownConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyUnknownConfiguration: %+v", err)
	}

	return encoded, nil
}
