package databases

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Database = ReadOnlyFollowingDatabase{}

type ReadOnlyFollowingDatabase struct {
	Properties *ReadOnlyFollowingDatabaseProperties `json:"properties,omitempty"`

	// Fields inherited from Database

	Id       *string `json:"id,omitempty"`
	Kind     Kind    `json:"kind"`
	Location *string `json:"location,omitempty"`
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
}

func (s ReadOnlyFollowingDatabase) Database() BaseDatabaseImpl {
	return BaseDatabaseImpl{
		Id:       s.Id,
		Kind:     s.Kind,
		Location: s.Location,
		Name:     s.Name,
		Type:     s.Type,
	}
}

var _ json.Marshaler = ReadOnlyFollowingDatabase{}

func (s ReadOnlyFollowingDatabase) MarshalJSON() ([]byte, error) {
	type wrapper ReadOnlyFollowingDatabase
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ReadOnlyFollowingDatabase: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ReadOnlyFollowingDatabase: %+v", err)
	}

	decoded["kind"] = "ReadOnlyFollowing"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ReadOnlyFollowingDatabase: %+v", err)
	}

	return encoded, nil
}
