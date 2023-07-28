package codeversion

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PendingUploadCredentialDto interface {
}

func unmarshalPendingUploadCredentialDtoImplementation(input []byte) (PendingUploadCredentialDto, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling PendingUploadCredentialDto into map[string]interface: %+v", err)
	}

	value, ok := temp["credentialType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "SAS") {
		var out SASCredentialDto
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SASCredentialDto: %+v", err)
		}
		return out, nil
	}

	type RawPendingUploadCredentialDtoImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawPendingUploadCredentialDtoImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}
