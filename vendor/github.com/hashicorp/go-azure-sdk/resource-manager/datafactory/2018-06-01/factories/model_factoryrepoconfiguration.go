package factories

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FactoryRepoConfiguration interface {
	FactoryRepoConfiguration() BaseFactoryRepoConfigurationImpl
}

var _ FactoryRepoConfiguration = BaseFactoryRepoConfigurationImpl{}

type BaseFactoryRepoConfigurationImpl struct {
	AccountName         string  `json:"accountName"`
	CollaborationBranch string  `json:"collaborationBranch"`
	DisablePublish      *bool   `json:"disablePublish,omitempty"`
	LastCommitId        *string `json:"lastCommitId,omitempty"`
	RepositoryName      string  `json:"repositoryName"`
	RootFolder          string  `json:"rootFolder"`
	Type                string  `json:"type"`
}

func (s BaseFactoryRepoConfigurationImpl) FactoryRepoConfiguration() BaseFactoryRepoConfigurationImpl {
	return s
}

var _ FactoryRepoConfiguration = RawFactoryRepoConfigurationImpl{}

// RawFactoryRepoConfigurationImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawFactoryRepoConfigurationImpl struct {
	factoryRepoConfiguration BaseFactoryRepoConfigurationImpl
	Type                     string
	Values                   map[string]interface{}
}

func (s RawFactoryRepoConfigurationImpl) FactoryRepoConfiguration() BaseFactoryRepoConfigurationImpl {
	return s.factoryRepoConfiguration
}

func UnmarshalFactoryRepoConfigurationImplementation(input []byte) (FactoryRepoConfiguration, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling FactoryRepoConfiguration into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "FactoryGitHubConfiguration") {
		var out FactoryGitHubConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FactoryGitHubConfiguration: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FactoryVSTSConfiguration") {
		var out FactoryVSTSConfiguration
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FactoryVSTSConfiguration: %+v", err)
		}
		return out, nil
	}

	var parent BaseFactoryRepoConfigurationImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseFactoryRepoConfigurationImpl: %+v", err)
	}

	return RawFactoryRepoConfigurationImpl{
		factoryRepoConfiguration: parent,
		Type:                     value,
		Values:                   temp,
	}, nil

}
