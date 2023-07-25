package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyConfiguration = ContentKeyPolicyClearKeyConfiguration{}

type ContentKeyPolicyClearKeyConfiguration struct {

	// Fields inherited from ContentKeyPolicyConfiguration
}

var _ json.Marshaler = ContentKeyPolicyClearKeyConfiguration{}

func (s ContentKeyPolicyClearKeyConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyClearKeyConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyClearKeyConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyClearKeyConfiguration: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyClearKeyConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyClearKeyConfiguration: %+v", err)
	}

	return encoded, nil
}
