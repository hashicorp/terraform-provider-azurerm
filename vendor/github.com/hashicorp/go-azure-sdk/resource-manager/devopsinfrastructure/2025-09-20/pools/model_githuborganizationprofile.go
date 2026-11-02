package pools

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OrganizationProfile = GitHubOrganizationProfile{}

type GitHubOrganizationProfile struct {
	Organizations []GitHubOrganization `json:"organizations"`

	// Fields inherited from OrganizationProfile

	Kind string `json:"kind"`
}

func (s GitHubOrganizationProfile) OrganizationProfile() BaseOrganizationProfileImpl {
	return BaseOrganizationProfileImpl{
		Kind: s.Kind,
	}
}

var _ json.Marshaler = GitHubOrganizationProfile{}

func (s GitHubOrganizationProfile) MarshalJSON() ([]byte, error) {
	type wrapper GitHubOrganizationProfile
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GitHubOrganizationProfile: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GitHubOrganizationProfile: %+v", err)
	}

	decoded["kind"] = "GitHub"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GitHubOrganizationProfile: %+v", err)
	}

	return encoded, nil
}
