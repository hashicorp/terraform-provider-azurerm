package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeadLetterDestination = StorageBlobDeadLetterDestination{}

type StorageBlobDeadLetterDestination struct {
	Properties *StorageBlobDeadLetterDestinationProperties `json:"properties,omitempty"`

	// Fields inherited from DeadLetterDestination

	EndpointType DeadLetterEndPointType `json:"endpointType"`
}

func (s StorageBlobDeadLetterDestination) DeadLetterDestination() BaseDeadLetterDestinationImpl {
	return BaseDeadLetterDestinationImpl{
		EndpointType: s.EndpointType,
	}
}

var _ json.Marshaler = StorageBlobDeadLetterDestination{}

func (s StorageBlobDeadLetterDestination) MarshalJSON() ([]byte, error) {
	type wrapper StorageBlobDeadLetterDestination
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StorageBlobDeadLetterDestination: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StorageBlobDeadLetterDestination: %+v", err)
	}

	decoded["endpointType"] = "StorageBlob"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StorageBlobDeadLetterDestination: %+v", err)
	}

	return encoded, nil
}
