package frontdoors

import (
	"encoding/json"
	"fmt"
)

type ForwardingConfiguration struct {
	BackendPool          *SubResource                 `json:"backendPool,omitempty"`
	CacheConfiguration   *CacheConfiguration          `json:"cacheConfiguration,omitempty"`
	CustomForwardingPath *string                      `json:"customForwardingPath,omitempty"`
	ForwardingProtocol   *FrontDoorForwardingProtocol `json:"forwardingProtocol,omitempty"`
}

var _ json.Marshaler = ForwardingConfiguration{}

func (s ForwardingConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper ForwardingConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ForwardingConfiguration: %+v", err)
	}

	var decoded map[string]string
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ForwardingConfiguration: %+v", err)
	}
	decoded["OdataType"] = "#Microsoft.Azure.FrontDoor.Models.FrontdoorForwardingConfiguration"

	encoded, err = json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ForwardingConfiguration: %+v", err)
	}

	return encoded, nil
}
