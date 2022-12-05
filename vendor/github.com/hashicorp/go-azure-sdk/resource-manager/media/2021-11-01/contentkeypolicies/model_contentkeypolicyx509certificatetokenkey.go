package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyRestrictionTokenKey = ContentKeyPolicyX509CertificateTokenKey{}

type ContentKeyPolicyX509CertificateTokenKey struct {
	RawBody string `json:"rawBody"`

	// Fields inherited from ContentKeyPolicyRestrictionTokenKey
}

var _ json.Marshaler = ContentKeyPolicyX509CertificateTokenKey{}

func (s ContentKeyPolicyX509CertificateTokenKey) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyX509CertificateTokenKey
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyX509CertificateTokenKey: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyX509CertificateTokenKey: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyX509CertificateTokenKey"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyX509CertificateTokenKey: %+v", err)
	}

	return encoded, nil
}
