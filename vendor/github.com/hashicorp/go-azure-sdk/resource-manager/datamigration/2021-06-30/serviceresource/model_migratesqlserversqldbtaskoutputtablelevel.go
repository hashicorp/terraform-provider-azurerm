package serviceresource

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MigrateSqlServerSqlDbTaskOutput = MigrateSqlServerSqlDbTaskOutputTableLevel{}

type MigrateSqlServerSqlDbTaskOutputTableLevel struct {
	EndedOn             *string         `json:"endedOn,omitempty"`
	ErrorPrefix         *string         `json:"errorPrefix,omitempty"`
	ItemsCompletedCount *int64          `json:"itemsCompletedCount,omitempty"`
	ItemsCount          *int64          `json:"itemsCount,omitempty"`
	ObjectName          *string         `json:"objectName,omitempty"`
	ResultPrefix        *string         `json:"resultPrefix,omitempty"`
	StartedOn           *string         `json:"startedOn,omitempty"`
	State               *MigrationState `json:"state,omitempty"`
	StatusMessage       *string         `json:"statusMessage,omitempty"`

	// Fields inherited from MigrateSqlServerSqlDbTaskOutput

	Id         *string `json:"id,omitempty"`
	ResultType string  `json:"resultType"`
}

func (s MigrateSqlServerSqlDbTaskOutputTableLevel) MigrateSqlServerSqlDbTaskOutput() BaseMigrateSqlServerSqlDbTaskOutputImpl {
	return BaseMigrateSqlServerSqlDbTaskOutputImpl{
		Id:         s.Id,
		ResultType: s.ResultType,
	}
}

var _ json.Marshaler = MigrateSqlServerSqlDbTaskOutputTableLevel{}

func (s MigrateSqlServerSqlDbTaskOutputTableLevel) MarshalJSON() ([]byte, error) {
	type wrapper MigrateSqlServerSqlDbTaskOutputTableLevel
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MigrateSqlServerSqlDbTaskOutputTableLevel: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MigrateSqlServerSqlDbTaskOutputTableLevel: %+v", err)
	}

	decoded["resultType"] = "TableLevelOutput"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MigrateSqlServerSqlDbTaskOutputTableLevel: %+v", err)
	}

	return encoded, nil
}
