package service

import (
	"encoding/json"
	"fmt"
)

var _ ServicePlacementPolicy = ServicePlacementRequireDomainDistributionPolicy{}

type ServicePlacementRequireDomainDistributionPolicy struct {
	DomainName string `json:"domainName"`

	// Fields inherited from ServicePlacementPolicy
}

var _ json.Marshaler = ServicePlacementRequireDomainDistributionPolicy{}

func (s ServicePlacementRequireDomainDistributionPolicy) MarshalJSON() ([]byte, error) {
	type wrapper ServicePlacementRequireDomainDistributionPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePlacementRequireDomainDistributionPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePlacementRequireDomainDistributionPolicy: %+v", err)
	}
	decoded["type"] = "RequiredDomainDistribution"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePlacementRequireDomainDistributionPolicy: %+v", err)
	}

	return encoded, nil
}
