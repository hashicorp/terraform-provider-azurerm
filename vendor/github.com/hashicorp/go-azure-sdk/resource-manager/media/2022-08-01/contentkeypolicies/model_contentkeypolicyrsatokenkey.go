package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyRestrictionTokenKey = ContentKeyPolicyRsaTokenKey{}

type ContentKeyPolicyRsaTokenKey struct {
	Exponent string `json:"exponent"`
	Modulus  string `json:"modulus"`

	// Fields inherited from ContentKeyPolicyRestrictionTokenKey
}

var _ json.Marshaler = ContentKeyPolicyRsaTokenKey{}

func (s ContentKeyPolicyRsaTokenKey) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyRsaTokenKey
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyRsaTokenKey: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyRsaTokenKey: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyRsaTokenKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyRsaTokenKey: %+v", err)
	}

	return encoded, nil
}
