package servers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerPropertiesForCreate interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawServerPropertiesForCreateImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalServerPropertiesForCreateImplementation(input []byte) (ServerPropertiesForCreate, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ServerPropertiesForCreate into map[string]interface: %+v", err)
	}

	value, ok := temp["createMode"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Default") {
		var out ServerPropertiesForDefaultCreate
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForDefaultCreate: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GeoRestore") {
		var out ServerPropertiesForGeoRestore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForGeoRestore: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Replica") {
		var out ServerPropertiesForReplica
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForReplica: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "PointInTimeRestore") {
		var out ServerPropertiesForRestore
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ServerPropertiesForRestore: %+v", err)
		}
		return out, nil
	}

	out := RawServerPropertiesForCreateImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
