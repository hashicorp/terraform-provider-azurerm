package services

import (
	"encoding/json"
	"fmt"
)

var _ Partition = NamedPartitionScheme{}

type NamedPartitionScheme struct {
	Names []string `json:"names"`

	// Fields inherited from Partition
}

var _ json.Marshaler = NamedPartitionScheme{}

func (s NamedPartitionScheme) MarshalJSON() ([]byte, error) {
	type wrapper NamedPartitionScheme
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling NamedPartitionScheme: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling NamedPartitionScheme: %+v", err)
	}
	decoded["partitionScheme"] = "Named"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling NamedPartitionScheme: %+v", err)
	}

	return encoded, nil
}
