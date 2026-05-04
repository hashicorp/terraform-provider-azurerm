package factories

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FactoryRepoConfiguration = FactoryGitHubConfiguration{}

type FactoryGitHubConfiguration struct {
	ClientId     *string             `json:"clientId,omitempty"`
	ClientSecret *GitHubClientSecret `json:"clientSecret,omitempty"`
	HostName     *string             `json:"hostName,omitempty"`

	// Fields inherited from FactoryRepoConfiguration

	AccountName         string  `json:"accountName"`
	CollaborationBranch string  `json:"collaborationBranch"`
	DisablePublish      *bool   `json:"disablePublish,omitempty"`
	LastCommitId        *string `json:"lastCommitId,omitempty"`
	RepositoryName      string  `json:"repositoryName"`
	RootFolder          string  `json:"rootFolder"`
	Type                string  `json:"type"`
}

func (s FactoryGitHubConfiguration) FactoryRepoConfiguration() BaseFactoryRepoConfigurationImpl {
	return BaseFactoryRepoConfigurationImpl{
		AccountName:         s.AccountName,
		CollaborationBranch: s.CollaborationBranch,
		DisablePublish:      s.DisablePublish,
		LastCommitId:        s.LastCommitId,
		RepositoryName:      s.RepositoryName,
		RootFolder:          s.RootFolder,
		Type:                s.Type,
	}
}

var _ json.Marshaler = FactoryGitHubConfiguration{}

func (s FactoryGitHubConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper FactoryGitHubConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling FactoryGitHubConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling FactoryGitHubConfiguration: %+v", err)
	}

	decoded["type"] = "FactoryGitHubConfiguration"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling FactoryGitHubConfiguration: %+v", err)
	}

	return encoded, nil
}
