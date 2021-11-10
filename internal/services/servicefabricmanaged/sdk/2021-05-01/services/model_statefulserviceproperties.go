package services

import (
	"encoding/json"
	"fmt"
)

type StatefulServiceProperties struct {
	HasPersistedState          *bool   `json:"hasPersistedState,omitempty"`
	MinReplicaSetSize          *int64  `json:"minReplicaSetSize,omitempty"`
	QuorumLossWaitDuration     *string `json:"quorumLossWaitDuration,omitempty"`
	ReplicaRestartWaitDuration *string `json:"replicaRestartWaitDuration,omitempty"`
	ServicePlacementTimeLimit  *string `json:"servicePlacementTimeLimit,omitempty"`
	StandByReplicaKeepDuration *string `json:"standByReplicaKeepDuration,omitempty"`
	TargetReplicaSetSize       *int64  `json:"targetReplicaSetSize,omitempty"`
}

var _ json.Marshaler = StatefulServiceProperties{}

func (s StatefulServiceProperties) MarshalJSON() ([]byte, error) {
	type wrapper StatefulServiceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StatefulServiceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StatefulServiceProperties: %+v", err)
	}
	decoded["serviceKind"] = "Stateful"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StatefulServiceProperties: %+v", err)
	}

	return encoded, nil
}
