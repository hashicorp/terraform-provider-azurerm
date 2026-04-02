package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OrganizationProfile = AzureDevOpsOrganizationProfile{}

type AzureDevOpsOrganizationProfile struct {
	Alias             *string                       `json:"alias,omitempty"`
	Organizations     []Organization                `json:"organizations"`
	PermissionProfile *AzureDevOpsPermissionProfile `json:"permissionProfile,omitempty"`

	// Fields inherited from OrganizationProfile

	Kind string `json:"kind"`
}

func (s AzureDevOpsOrganizationProfile) OrganizationProfile() BaseOrganizationProfileImpl {
	return BaseOrganizationProfileImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = AzureDevOpsOrganizationProfile{}

func (s AzureDevOpsOrganizationProfile) MarshalJSON() ([]byte, error) {
	type wrapper AzureDevOpsOrganizationProfile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureDevOpsOrganizationProfile: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureDevOpsOrganizationProfile: %+v", err)
	}

	decoded["kind"] = "AzureDevOps"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureDevOpsOrganizationProfile: %+v", err)
	}

	return encoded, nil
}
