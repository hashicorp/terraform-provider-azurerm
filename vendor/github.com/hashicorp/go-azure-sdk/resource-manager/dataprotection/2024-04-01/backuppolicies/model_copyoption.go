package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CopyOption interface {
	CopyOption() BaseCopyOptionImpl
}

var _ CopyOption = BaseCopyOptionImpl{}

type BaseCopyOptionImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseCopyOptionImpl) CopyOption() BaseCopyOptionImpl {
	return s
}

var _ CopyOption = RawCopyOptionImpl{}

// RawCopyOptionImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawCopyOptionImpl struct {
	copyOption BaseCopyOptionImpl
	Type       string
	Values     map[string]interface{}
}

func (s RawCopyOptionImpl) CopyOption() BaseCopyOptionImpl {
	return s.copyOption
}

func UnmarshalCopyOptionImplementation(input []byte) (CopyOption, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling CopyOption into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "CopyOnExpiryOption") {
		var out CopyOnExpiryOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CopyOnExpiryOption: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "CustomCopyOption") {
		var out CustomCopyOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomCopyOption: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImmediateCopyOption") {
		var out ImmediateCopyOption
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImmediateCopyOption: %+v", err)
		}
		return out, nil
	}

	var parent BaseCopyOptionImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseCopyOptionImpl: %+v", err)
	}

	return RawCopyOptionImpl{
		copyOption: parent,
		Type:       value,
		Values:     temp,
	}, nil

}
