package eventsubscriptions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeadLetterDestination interface {
}

func unmarshalDeadLetterDestinationImplementation(input []byte) (DeadLetterDestination, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeadLetterDestination into map[string]interface: %+v", err)
	}

	value, ok := temp["endpointType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "StorageBlob") {
		var out StorageBlobDeadLetterDestination
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StorageBlobDeadLetterDestination: %+v", err)
		}
		return out, nil
	}

	type RawDeadLetterDestinationImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawDeadLetterDestinationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
