package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyRestriction = ContentKeyPolicyOpenRestriction{}

type ContentKeyPolicyOpenRestriction struct {

	// Fields inherited from ContentKeyPolicyRestriction
}

var _ json.Marshaler = ContentKeyPolicyOpenRestriction{}

func (s ContentKeyPolicyOpenRestriction) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyOpenRestriction
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyOpenRestriction: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyOpenRestriction: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyOpenRestriction"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyOpenRestriction: %+v", err)
	}

	return encoded, nil
}
