package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = DirectLineChannel{}

type DirectLineChannel struct {
	Properties *DirectLineChannelProperties `json:"properties,omitempty"`

	// Fields inherited from Channel
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = DirectLineChannel{}

func (s DirectLineChannel) MarshalJSON() ([]byte, error) {
	type wrapper DirectLineChannel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DirectLineChannel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DirectLineChannel: %+v", err)
	}
	decoded["channelName"] = "DirectLineChannel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DirectLineChannel: %+v", err)
	}

	return encoded, nil
}
