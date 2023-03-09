package vaultcertificates

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResourceCertificateDetails interface {
}

func unmarshalResourceCertificateDetailsImplementation(input []byte) (ResourceCertificateDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling ResourceCertificateDetails into map[string]interface: %+v", err)
	}

	value, ok := temp["authType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureActiveDirectory") {
		var out ResourceCertificateAndAadDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResourceCertificateAndAadDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AccessControlService") {
		var out ResourceCertificateAndAcsDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ResourceCertificateAndAcsDetails: %+v", err)
		}
		return out, nil
	}

	type RawResourceCertificateDetailsImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawResourceCertificateDetailsImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
