package serviceresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MongoDbProgress interface {
	MongoDbProgress() BaseMongoDbProgressImpl
}

var _ MongoDbProgress = BaseMongoDbProgressImpl{}

type BaseMongoDbProgressImpl struct {
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

func (s BaseMongoDbProgressImpl) MongoDbProgress() BaseMongoDbProgressImpl {
	return s
}

var _ MongoDbProgress = RawMongoDbProgressImpl{}

// RawMongoDbProgressImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawMongoDbProgressImpl struct {
	mongoDbProgress BaseMongoDbProgressImpl
	Type            string
	Values          map[string]interface{}
}

func (s RawMongoDbProgressImpl) MongoDbProgress() BaseMongoDbProgressImpl {
	return s.mongoDbProgress
}

func UnmarshalMongoDbProgressImplementation(input []byte) (MongoDbProgress, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling MongoDbProgress into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["resultType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "Collection") {
		var out MongoDbCollectionProgress
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbCollectionProgress: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Database") {
		var out MongoDbDatabaseProgress
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbDatabaseProgress: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Migration") {
		var out MongoDbMigrationProgress
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MongoDbMigrationProgress: %+v", err)
		}
		return out, nil
	}

	var parent BaseMongoDbProgressImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseMongoDbProgressImpl: %+v", err)
	}

	return RawMongoDbProgressImpl{
		mongoDbProgress: parent,
		Type:            value,
		Values:          temp,
	}, nil

}
