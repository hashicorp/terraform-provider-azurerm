package service

import (
	"encoding/json"
	"fmt"
)

var _ ServicePlacementPolicy = ServicePlacementInvalidDomainPolicy{}

type ServicePlacementInvalidDomainPolicy struct {
	DomainName string `json:"domainName"`

	// Fields inherited from ServicePlacementPolicy
}

var _ json.Marshaler = ServicePlacementInvalidDomainPolicy{}

func (s ServicePlacementInvalidDomainPolicy) MarshalJSON() ([]byte, error) {
	type wrapper ServicePlacementInvalidDomainPolicy
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ServicePlacementInvalidDomainPolicy: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ServicePlacementInvalidDomainPolicy: %+v", err)
	}
	decoded["type"] = "InvalidDomain"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ServicePlacementInvalidDomainPolicy: %+v", err)
	}

	return encoded, nil
}
