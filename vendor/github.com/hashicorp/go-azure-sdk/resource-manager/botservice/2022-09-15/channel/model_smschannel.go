package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = SmsChannel{}

type SmsChannel struct {
	Properties *SmsChannelProperties `json:"properties,omitempty"`

	// Fields inherited from Channel
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

var _ json.Marshaler = SmsChannel{}

func (s SmsChannel) MarshalJSON() ([]byte, error) {
	type wrapper SmsChannel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SmsChannel: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SmsChannel: %+v", err)
	}
	decoded["channelName"] = "SmsChannel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SmsChannel: %+v", err)
	}

	return encoded, nil
}
