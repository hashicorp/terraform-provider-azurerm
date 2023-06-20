package contentkeypolicies

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ContentKeyPolicyConfiguration = ContentKeyPolicyPlayReadyConfiguration{}

type ContentKeyPolicyPlayReadyConfiguration struct {
	Licenses           []ContentKeyPolicyPlayReadyLicense `json:"licenses"`
	ResponseCustomData *string                            `json:"responseCustomData,omitempty"`

	// Fields inherited from ContentKeyPolicyConfiguration
}

var _ json.Marshaler = ContentKeyPolicyPlayReadyConfiguration{}

func (s ContentKeyPolicyPlayReadyConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ContentKeyPolicyPlayReadyConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ContentKeyPolicyPlayReadyConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyPlayReadyConfiguration: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.ContentKeyPolicyPlayReadyConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ContentKeyPolicyPlayReadyConfiguration: %+v", err)
	}

	return encoded, nil
}
