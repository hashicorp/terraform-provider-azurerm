package deploymentscripts

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeploymentScript = AzureCliScript{}

type AzureCliScript struct {
	Properties AzureCliScriptProperties `json:"properties"`

	// Fields inherited from DeploymentScript
	Id         *string                   `json:"id,omitempty"`
	Identity   *identity.UserAssignedMap `json:"identity,omitempty"`
	Location   string                    `json:"location"`
	Name       *string                   `json:"name,omitempty"`
	SystemData *systemdata.SystemData    `json:"systemData,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}

var _ json.Marshaler = AzureCliScript{}

func (s AzureCliScript) MarshalJSON() ([]byte, error) {
	type wrapper AzureCliScript
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureCliScript: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureCliScript: %+v", err)
	}
	decoded["kind"] = "AzureCLI"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureCliScript: %+v", err)
	}

	return encoded, nil
}
