package settings

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Setting interface {
	Setting() BaseSettingImpl
}

var _ Setting = BaseSettingImpl{}

type BaseSettingImpl struct {
	Id         *string                `json:"id,omitempty"`
	Kind       SettingKind            `json:"kind"`
	Name       *string                `json:"name,omitempty"`
	SystemData *systemdata.SystemData `json:"systemData,omitempty"`
	Type       *string                `json:"type,omitempty"`
}

func (s BaseSettingImpl) Setting() BaseSettingImpl {
	return s
}

var _ Setting = RawSettingImpl{}

// RawSettingImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawSettingImpl struct {
	setting BaseSettingImpl
	Type    string
	Values  map[string]interface{}
}

func (s RawSettingImpl) Setting() BaseSettingImpl {
	return s.setting
}

func (s RawSettingImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalSettingImplementation(input []byte) (Setting, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Setting into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AlertSyncSettings") {
		var out AlertSyncSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AlertSyncSettings: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DataExportSettings") {
		var out DataExportSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DataExportSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseSettingImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSettingImpl: %+v", err)
	}

	return RawSettingImpl{
		setting: parent,
		Type:    value,
		Values:  temp,
	}, nil

}
