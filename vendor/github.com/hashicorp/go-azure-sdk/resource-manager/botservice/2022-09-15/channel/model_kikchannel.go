package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = KikChannel{}

type KikChannel struct {
	Properties *KikChannelProperties `json:"properties,omitempty"`

	// Fields inherited from Channel
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = KikChannel{}

func (s KikChannel) MarshalJSON() ([]byte, error) {
	type wrapper KikChannel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KikChannel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KikChannel: %+v", err)
	}
	decoded["channelName"] = "KikChannel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KikChannel: %+v", err)
	}

	return encoded, nil
}
