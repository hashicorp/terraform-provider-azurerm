package backupandexport

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupStoreDetails interface {
	BackupStoreDetails() BaseBackupStoreDetailsImpl
}

var _ BackupStoreDetails = BaseBackupStoreDetailsImpl{}

type BaseBackupStoreDetailsImpl struct {
	ObjectType string `json:"objectType"`
}

func (s BaseBackupStoreDetailsImpl) BackupStoreDetails() BaseBackupStoreDetailsImpl {
	return s
}

var _ BackupStoreDetails = RawBackupStoreDetailsImpl{}

// RawBackupStoreDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawBackupStoreDetailsImpl struct {
	backupStoreDetails BaseBackupStoreDetailsImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawBackupStoreDetailsImpl) BackupStoreDetails() BaseBackupStoreDetailsImpl {
	return s.backupStoreDetails
}

func UnmarshalBackupStoreDetailsImplementation(input []byte) (BackupStoreDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BackupStoreDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "FullBackupStoreDetails") {
		var out FullBackupStoreDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FullBackupStoreDetails: %+v", err)
		}
		return out, nil
	}

	var parent BaseBackupStoreDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBackupStoreDetailsImpl: %+v", err)
	}

	return RawBackupStoreDetailsImpl{
		backupStoreDetails: parent,
		Type:               value,
		Values:             temp,
	}, nil

}
