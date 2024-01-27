package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = FacebookChannel{}

type FacebookChannel struct {
	Properties *FacebookChannelProperties `json:"properties,omitempty"`

	// Fields inherited from Channel
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = FacebookChannel{}

func (s FacebookChannel) MarshalJSON() ([]byte, error) {
	type wrapper FacebookChannel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FacebookChannel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FacebookChannel: %+v", err)
	}
	decoded["channelName"] = "FacebookChannel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FacebookChannel: %+v", err)
	}

	return encoded, nil
}
