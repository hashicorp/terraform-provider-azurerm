package services

import (
	"encoding/json"
	"fmt"
)

var _ ServicePlacementPolicy = ServicePlacementPreferPrimaryDomainPolicy{}

type ServicePlacementPreferPrimaryDomainPolicy struct {
	DomainName string `json:"domainName"`

	// Fields inherited from ServicePlacementPolicy
}

var _ json.Marshaler = ServicePlacementPreferPrimaryDomainPolicy{}

func (s ServicePlacementPreferPrimaryDomainPolicy) MarshalJSON() ([]byte, error) {
	type wrapper ServicePlacementPreferPrimaryDomainPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePlacementPreferPrimaryDomainPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePlacementPreferPrimaryDomainPolicy: %+v", err)
	}
	decoded["type"] = "PreferredPrimaryDomain"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePlacementPreferPrimaryDomainPolicy: %+v", err)
	}

	return encoded, nil
}
