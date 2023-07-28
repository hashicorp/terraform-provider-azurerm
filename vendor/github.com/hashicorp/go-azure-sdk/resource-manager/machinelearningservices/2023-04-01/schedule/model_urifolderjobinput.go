package schedule

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobInput = UriFolderJobInput{}

type UriFolderJobInput struct {
	Mode *InputDeliveryMode `json:"mode,omitempty"`
	Uri  string             `json:"uri"`

	// Fields inherited from JobInput
	Description *string `json:"description,omitempty"`
}

var _ json.Marshaler = UriFolderJobInput{}

func (s UriFolderJobInput) MarshalJSON() ([]byte, error) {
	type wrapper UriFolderJobInput
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling UriFolderJobInput: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling UriFolderJobInput: %+v", err)
	}
	decoded["jobInputType"] = "uri_folder"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling UriFolderJobInput: %+v", err)
	}

	return encoded, nil
}
