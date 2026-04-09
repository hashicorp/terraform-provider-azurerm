package synchronizationsetting

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SynchronizationSetting interface {
	SynchronizationSetting() BaseSynchronizationSettingImpl
}

var _ SynchronizationSetting = BaseSynchronizationSettingImpl{}

type BaseSynchronizationSettingImpl struct {
	Id   *string                    `json:"id,omitempty"`
	Kind SynchronizationSettingKind `json:"kind"`
	Name *string                    `json:"name,omitempty"`
	Type *string                    `json:"type,omitempty"`
}

func (s BaseSynchronizationSettingImpl) SynchronizationSetting() BaseSynchronizationSettingImpl {
	return s
}

var _ SynchronizationSetting = RawSynchronizationSettingImpl{}

// RawSynchronizationSettingImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSynchronizationSettingImpl struct {
	synchronizationSetting BaseSynchronizationSettingImpl
	Type                   string
	Values                 map[string]interface{}
}

func (s RawSynchronizationSettingImpl) SynchronizationSetting() BaseSynchronizationSettingImpl {
	return s.synchronizationSetting
}

func UnmarshalSynchronizationSettingImplementation(input []byte) (SynchronizationSetting, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SynchronizationSetting into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ScheduleBased") {
		var out ScheduledSynchronizationSetting
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ScheduledSynchronizationSetting: %+v", err)
		}
		return out, nil
	}

	var parent BaseSynchronizationSettingImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSynchronizationSettingImpl: %+v", err)
	}

	return RawSynchronizationSettingImpl{
		synchronizationSetting: parent,
		Type:                   value,
		Values:                 temp,
	}, nil

}
