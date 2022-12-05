package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyPlayReadyContentKeyLocation = ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader{}

type ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader struct {

	// Fields inherited from ContentKeyPolicyPlayReadyContentKeyLocation
}

var _ json.Marshaler = ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader{}

func (s ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader: %+v", err)
	}

	return encoded, nil
}
