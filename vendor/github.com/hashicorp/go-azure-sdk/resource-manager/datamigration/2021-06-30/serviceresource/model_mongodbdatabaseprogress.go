package serviceresource

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ MongoDbProgress = MongoDbDatabaseProgress{}

type MongoDbDatabaseProgress struct {
	Collections *map[string]MongoDbProgress `json:"collections,omitempty"`

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

func (s MongoDbDatabaseProgress) MongoDbProgress() BaseMongoDbProgressImpl {
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

func (o *MongoDbDatabaseProgress) GetLastEventTimeAsTime() (*time.Time, error) {
	if o.LastEventTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastEventTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MongoDbDatabaseProgress) SetLastEventTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastEventTime = &formatted
}

func (o *MongoDbDatabaseProgress) GetLastReplayTimeAsTime() (*time.Time, error) {
	if o.LastReplayTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastReplayTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MongoDbDatabaseProgress) SetLastReplayTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastReplayTime = &formatted
}

var _ json.Marshaler = MongoDbDatabaseProgress{}

func (s MongoDbDatabaseProgress) MarshalJSON() ([]byte, error) {
	type wrapper MongoDbDatabaseProgress
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling MongoDbDatabaseProgress: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling MongoDbDatabaseProgress: %+v", err)
	}

	decoded["resultType"] = "Database"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling MongoDbDatabaseProgress: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &MongoDbDatabaseProgress{}

func (s *MongoDbDatabaseProgress) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
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
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.BytesCopied = decoded.BytesCopied
	s.DocumentsCopied = decoded.DocumentsCopied
	s.ElapsedTime = decoded.ElapsedTime
	s.Errors = decoded.Errors
	s.EventsPending = decoded.EventsPending
	s.EventsReplayed = decoded.EventsReplayed
	s.LastEventTime = decoded.LastEventTime
	s.LastReplayTime = decoded.LastReplayTime
	s.Name = decoded.Name
	s.QualifiedName = decoded.QualifiedName
	s.ResultType = decoded.ResultType
	s.State = decoded.State
	s.TotalBytes = decoded.TotalBytes
	s.TotalDocuments = decoded.TotalDocuments

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling MongoDbDatabaseProgress into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["collections"]; ok {
		var dictionaryTemp map[string]json.RawMessage
		if err := json.Unmarshal(v, &dictionaryTemp); err != nil {
			return fmt.Errorf("unmarshaling Collections into dictionary map[string]json.RawMessage: %+v", err)
		}

		output := make(map[string]MongoDbProgress)
		for key, val := range dictionaryTemp {
			impl, err := UnmarshalMongoDbProgressImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling key %q field 'Collections' for 'MongoDbDatabaseProgress': %+v", key, err)
			}
			output[key] = impl
		}
		s.Collections = &output
	}

	return nil
}
