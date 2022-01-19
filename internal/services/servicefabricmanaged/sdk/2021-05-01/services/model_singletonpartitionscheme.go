package services

import (
	"encoding/json"
	"fmt"
)

var _ Partition = SingletonPartitionScheme{}

type SingletonPartitionScheme struct {

	// Fields inherited from Partition
}

var _ json.Marshaler = SingletonPartitionScheme{}

func (s SingletonPartitionScheme) MarshalJSON() ([]byte, error) {
	type wrapper SingletonPartitionScheme
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SingletonPartitionScheme: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SingletonPartitionScheme: %+v", err)
	}
	decoded["partitionScheme"] = "Singleton"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SingletonPartitionScheme: %+v", err)
	}

	return encoded, nil
}
