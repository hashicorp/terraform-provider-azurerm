package contentkeypolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPolicyPlayReadyContentKeyLocation interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawContentKeyPolicyPlayReadyContentKeyLocationImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalContentKeyPolicyPlayReadyContentKeyLocationImplementation(input []byte) (ContentKeyPolicyPlayReadyContentKeyLocation, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ContentKeyPolicyPlayReadyContentKeyLocation into map[string]interface: %+v", err)
	}

	value, ok := temp["@odata.type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader") {
		var out ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyPlayReadyContentEncryptionKeyFromHeader: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Media.ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier") {
		var out ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentKeyPolicyPlayReadyContentEncryptionKeyFromKeyIdentifier: %+v", err)
		}
		return out, nil
	}

	out := RawContentKeyPolicyPlayReadyContentKeyLocationImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
