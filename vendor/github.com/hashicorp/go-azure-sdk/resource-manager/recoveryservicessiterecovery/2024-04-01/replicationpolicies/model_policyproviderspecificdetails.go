package replicationpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyProviderSpecificDetails interface {
	PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl
}

var _ PolicyProviderSpecificDetails = BasePolicyProviderSpecificDetailsImpl{}

type BasePolicyProviderSpecificDetailsImpl struct {
	InstanceType string `json:"instanceType"`
}

func (s BasePolicyProviderSpecificDetailsImpl) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return s
}

var _ PolicyProviderSpecificDetails = RawPolicyProviderSpecificDetailsImpl{}

// RawPolicyProviderSpecificDetailsImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawPolicyProviderSpecificDetailsImpl struct {
	policyProviderSpecificDetails BasePolicyProviderSpecificDetailsImpl
	Type                          string
	Values                        map[string]interface{}
}

func (s RawPolicyProviderSpecificDetailsImpl) PolicyProviderSpecificDetails() BasePolicyProviderSpecificDetailsImpl {
	return s.policyProviderSpecificDetails
}

func UnmarshalPolicyProviderSpecificDetailsImplementation(input []byte) (PolicyProviderSpecificDetails, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling PolicyProviderSpecificDetails into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["instanceType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "A2A") {
		var out A2APolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into A2APolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaAzure") {
		var out HyperVReplicaAzurePolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaAzurePolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplicaBasePolicyDetails") {
		var out HyperVReplicaBasePolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaBasePolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplica2012R2") {
		var out HyperVReplicaBluePolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaBluePolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "HyperVReplica2012") {
		var out HyperVReplicaPolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into HyperVReplicaPolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageAzureV2") {
		var out InMageAzureV2PolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageAzureV2PolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageBasePolicyDetails") {
		var out InMageBasePolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageBasePolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMage") {
		var out InMagePolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMagePolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcmFailback") {
		var out InMageRcmFailbackPolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmFailbackPolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "InMageRcm") {
		var out InMageRcmPolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into InMageRcmPolicyDetails: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "VMwareCbt") {
		var out VMwareCbtPolicyDetails
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into VMwareCbtPolicyDetails: %+v", err)
		}
		return out, nil
	}

	var parent BasePolicyProviderSpecificDetailsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BasePolicyProviderSpecificDetailsImpl: %+v", err)
	}

	return RawPolicyProviderSpecificDetailsImpl{
		policyProviderSpecificDetails: parent,
		Type:                          value,
		Values:                        temp,
	}, nil

}
