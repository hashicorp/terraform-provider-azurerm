package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateProperties interface {
}

// RawModeOfTransitImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCertificatePropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalCertificatePropertiesImplementation(input []byte) (CertificateProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CertificateProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["type"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "ContentCertificate") {
		var out ContentCertificateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ContentCertificateProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KeyVaultCertificate") {
		var out KeyVaultCertificateProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeyVaultCertificateProperties: %+v", err)
		}
		return out, nil
	}

	out := RawCertificatePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
