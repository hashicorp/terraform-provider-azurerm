package serviceresource

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MongoDbProgress = MongoDbMigrationProgress{}

type MongoDbMigrationProgress struct {
	Databases *map[string]MongoDbDatabaseProgress `json:"databases,omitempty"`

	// Fields inherited from MongoDbProgress

	BytesCopied     int64                   `json:"bytesCopied"`
	DocumentsCopied int64                   `json:"documentsCopied"`
	ElapsedTime     string                  `json:"elapsedTime"`
	Errors          map[string]MongoDbError `json:"errors"`
	EventsPending   int64                   `json:"eventsPending"`
	EventsReplayed  int64                   `json:"eventsReplayed"`
	LastEventTime   *string                 `json:"lastEventTime,omitempty"`
	LastReplayTime  *string                 `json:"lastReplayTime,omitempty"`
	Name            *string                 `json:"name,omitempty"`
	QualifiedName   *string                 `json:"qualifiedName,omitempty"`
	ResultType      ResultType              `json:"resultType"`
	State           MongoDbMigrationState   `json:"state"`
	TotalBytes      int64                   `json:"totalBytes"`
	TotalDocuments  int64                   `json:"totalDocuments"`
}

func (s MongoDbMigrationProgress) MongoDbProgress() BaseMongoDbProgressImpl {
	return BaseMongoDbProgressImpl{
		BytesCopied:     s.BytesCopied,
		DocumentsCopied: s.DocumentsCopied,
		ElapsedTime:     s.ElapsedTime,
		Errors:          s.Errors,
		EventsPending:   s.EventsPending,
		EventsReplayed:  s.EventsReplayed,
		LastEventTime:   s.LastEventTime,
		LastReplayTime:  s.LastReplayTime,
		Name:            s.Name,
		QualifiedName:   s.QualifiedName,
		ResultType:      s.ResultType,
		State:           s.State,
		TotalBytes:      s.TotalBytes,
		TotalDocuments:  s.TotalDocuments,
	}
}

func (o *MongoDbMigrationProgress) GetLastEventTimeAsTime() (*time.Time, error) {
	if o.LastEventTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastEventTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MongoDbMigrationProgress) SetLastEventTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastEventTime = &formatted
}

func (o *MongoDbMigrationProgress) GetLastReplayTimeAsTime() (*time.Time, error) {
	if o.LastReplayTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastReplayTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MongoDbMigrationProgress) SetLastReplayTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastReplayTime = &formatted
}

var _ json.Marshaler = MongoDbMigrationProgress{}

func (s MongoDbMigrationProgress) MarshalJSON() ([]byte, error) {
	type wrapper MongoDbMigrationProgress
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MongoDbMigrationProgress: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MongoDbMigrationProgress: %+v", err)
	}

	decoded["resultType"] = "Migration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MongoDbMigrationProgress: %+v", err)
	}

	return encoded, nil
}
