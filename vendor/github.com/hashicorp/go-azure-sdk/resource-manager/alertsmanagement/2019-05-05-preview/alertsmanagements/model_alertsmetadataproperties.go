package alertsmanagements

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertsMetaDataProperties interface {
}

// RawAlertsMetaDataPropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawAlertsMetaDataPropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalAlertsMetaDataPropertiesImplementation(input []byte) (AlertsMetaDataProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AlertsMetaDataProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["metadataIdentifier"].(string)
	if !ok {
		return nil, nil
	}

	out := RawAlertsMetaDataPropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
