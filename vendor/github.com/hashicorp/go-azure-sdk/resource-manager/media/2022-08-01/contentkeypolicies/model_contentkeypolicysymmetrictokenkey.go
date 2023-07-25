package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyRestrictionTokenKey = ContentKeyPolicySymmetricTokenKey{}

type ContentKeyPolicySymmetricTokenKey struct {
	KeyValue string `json:"keyValue"`

	// Fields inherited from ContentKeyPolicyRestrictionTokenKey
}

var _ json.Marshaler = ContentKeyPolicySymmetricTokenKey{}

func (s ContentKeyPolicySymmetricTokenKey) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicySymmetricTokenKey
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicySymmetricTokenKey: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicySymmetricTokenKey: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicySymmetricTokenKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicySymmetricTokenKey: %+v", err)
	}

	return encoded, nil
}
