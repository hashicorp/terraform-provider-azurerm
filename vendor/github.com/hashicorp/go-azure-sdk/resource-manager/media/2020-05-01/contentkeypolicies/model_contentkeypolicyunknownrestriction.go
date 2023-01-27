package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyRestriction = ContentKeyPolicyUnknownRestriction{}

type ContentKeyPolicyUnknownRestriction struct {

	// Fields inherited from ContentKeyPolicyRestriction
}

var _ json.Marshaler = ContentKeyPolicyUnknownRestriction{}

func (s ContentKeyPolicyUnknownRestriction) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyUnknownRestriction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyUnknownRestriction: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyUnknownRestriction: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyUnknownRestriction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyUnknownRestriction: %+v", err)
	}

	return encoded, nil
}
