package databases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Database interface {
	Database() BaseDatabaseImpl
}

var _ Database = BaseDatabaseImpl{}

type BaseDatabaseImpl struct {
	Id       *string `json:"id,omitempty"`
	Kind     Kind    `json:"kind"`
	Location *string `json:"location,omitempty"`
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
}

func (s BaseDatabaseImpl) Database() BaseDatabaseImpl {
	return s
}

var _ Database = RawDatabaseImpl{}

// RawDatabaseImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDatabaseImpl struct {
	database BaseDatabaseImpl
	Type     string
	Values   map[string]interface{}
}

func (s RawDatabaseImpl) Database() BaseDatabaseImpl {
	return s.database
}

func UnmarshalDatabaseImplementation(input []byte) (Database, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Database into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "ReadOnlyFollowing") {
		var out ReadOnlyFollowingDatabase
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ReadOnlyFollowingDatabase: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ReadWrite") {
		var out ReadWriteDatabase
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ReadWriteDatabase: %+v", err)
		}
		return out, nil
	}

	var parent BaseDatabaseImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDatabaseImpl: %+v", err)
	}

	return RawDatabaseImpl{
		database: parent,
		Type:     value,
		Values:   temp,
	}, nil

}
