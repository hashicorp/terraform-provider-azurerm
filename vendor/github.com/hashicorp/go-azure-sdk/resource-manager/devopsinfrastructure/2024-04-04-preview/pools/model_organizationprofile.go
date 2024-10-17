package pools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OrganizationProfile interface {
	OrganizationProfile() BaseOrganizationProfileImpl
}

var _ OrganizationProfile = BaseOrganizationProfileImpl{}

type BaseOrganizationProfileImpl struct {
	Kind string `json:"kind"`
}

func (s BaseOrganizationProfileImpl) OrganizationProfile() BaseOrganizationProfileImpl {
	return s
}

var _ OrganizationProfile = RawOrganizationProfileImpl{}

// RawOrganizationProfileImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawOrganizationProfileImpl struct {
	organizationProfile BaseOrganizationProfileImpl
	Type                string
	Values              map[string]interface{}
}

func (s RawOrganizationProfileImpl) OrganizationProfile() BaseOrganizationProfileImpl {
	return s.organizationProfile
}

func UnmarshalOrganizationProfileImplementation(input []byte) (OrganizationProfile, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling OrganizationProfile into map[string]interface: %+v", err)
	}

	value, ok := temp["kind"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "AzureDevOps") {
		var out AzureDevOpsOrganizationProfile
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureDevOpsOrganizationProfile: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "GitHub") {
		var out GitHubOrganizationProfile
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GitHubOrganizationProfile: %+v", err)
		}
		return out, nil
	}

	var parent BaseOrganizationProfileImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseOrganizationProfileImpl: %+v", err)
	}

	return RawOrganizationProfileImpl{
		organizationProfile: parent,
		Type:                value,
		Values:              temp,
	}, nil

}
