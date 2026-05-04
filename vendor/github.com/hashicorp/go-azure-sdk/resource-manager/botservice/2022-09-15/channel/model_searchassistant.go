package channel

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Channel = SearchAssistant{}

type SearchAssistant struct {

	// Fields inherited from Channel

	ChannelName       string  `json:"channelName"`
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

func (s SearchAssistant) Channel() BaseChannelImpl {
	return BaseChannelImpl{
		ChannelName:       s.ChannelName,
		Etag:              s.Etag,
		Location:          s.Location,
		ProvisioningState: s.ProvisioningState,
	}
}

var _ json.Marshaler = SearchAssistant{}

func (s SearchAssistant) MarshalJSON() ([]byte, error) {
	type wrapper SearchAssistant
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SearchAssistant: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SearchAssistant: %+v", err)
	}

	decoded["channelName"] = "SearchAssistant"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SearchAssistant: %+v", err)
	}

	return encoded, nil
}
