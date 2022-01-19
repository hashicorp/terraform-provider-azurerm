package services

import (
	"encoding/json"
	"fmt"
)

var _ ServiceResourceProperties = StatelessServiceProperties{}

type StatelessServiceProperties struct {
	InstanceCount         int64  `json:"instanceCount"`
	MinInstanceCount      *int64 `json:"minInstanceCount,omitempty"`
	MinInstancePercentage *int64 `json:"minInstancePercentage,omitempty"`

	// Fields inherited from ServiceResourceProperties
	CorrelationScheme            *[]ServiceCorrelation         `json:"correlationScheme,omitempty"`
	DefaultMoveCost              *MoveCost                     `json:"defaultMoveCost,omitempty"`
	PartitionDescription         Partition                     `json:"partitionDescription"`
	PlacementConstraints         *string                       `json:"placementConstraints,omitempty"`
	ProvisioningState            *string                       `json:"provisioningState,omitempty"`
	ScalingPolicies              *[]ScalingPolicy              `json:"scalingPolicies,omitempty"`
	ServiceLoadMetrics           *[]ServiceLoadMetric          `json:"serviceLoadMetrics,omitempty"`
	ServicePackageActivationMode *ServicePackageActivationMode `json:"servicePackageActivationMode,omitempty"`
	ServicePlacementPolicies     *[]ServicePlacementPolicy     `json:"servicePlacementPolicies,omitempty"`
	ServiceTypeName              string                        `json:"serviceTypeName"`
}

var _ json.Marshaler = StatelessServiceProperties{}

func (s StatelessServiceProperties) MarshalJSON() ([]byte, error) {
	type wrapper StatelessServiceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StatelessServiceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StatelessServiceProperties: %+v", err)
	}
	decoded["serviceKind"] = "Stateless"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StatelessServiceProperties: %+v", err)
	}

	return encoded, nil
}
