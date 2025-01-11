package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = Omnichannel{}

type Omnichannel struct {

	// Fields inherited from Channel

	ChannelName       string  `json:"channelName"`
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

func (s Omnichannel) Channel() BaseChannelImpl {
	return BaseChannelImpl{
		ChannelName:       s.ChannelName,
		Etag:              s.Etag,
		Location:          s.Location,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = Omnichannel{}

func (s Omnichannel) MarshalJSON() ([]byte, error) {
	type wrapper Omnichannel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling Omnichannel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling Omnichannel: %+v", err)
	}

	decoded["channelName"] = "Omnichannel"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling Omnichannel: %+v", err)
	}

	return encoded, nil
}
