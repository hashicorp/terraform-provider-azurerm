package deploymentscripts

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentScript interface {
	DeploymentScript() BaseDeploymentScriptImpl
}

var _ DeploymentScript = BaseDeploymentScriptImpl{}

type BaseDeploymentScriptImpl struct {
	Id         *string                   `json:"id,omitempty"`
	Identity   *identity.UserAssignedMap `json:"identity,omitempty"`
	Kind       ScriptType                `json:"kind"`
	Location   string                    `json:"location"`
	Name       *string                   `json:"name,omitempty"`
	SystemData *systemdata.SystemData    `json:"systemData,omitempty"`
	Tags       *map[string]string        `json:"tags,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}

func (s BaseDeploymentScriptImpl) DeploymentScript() BaseDeploymentScriptImpl {
	return s
}

var _ DeploymentScript = RawDeploymentScriptImpl{}

// RawDeploymentScriptImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawDeploymentScriptImpl struct {
	deploymentScript BaseDeploymentScriptImpl
	Type             string
	Values           map[string]interface{}
}

func (s RawDeploymentScriptImpl) DeploymentScript() BaseDeploymentScriptImpl {
	return s.deploymentScript
}

func UnmarshalDeploymentScriptImplementation(input []byte) (DeploymentScript, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling DeploymentScript into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["kind"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AzureCLI") {
		var out AzureCliScript
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureCliScript: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AzurePowerShell") {
		var out AzurePowerShellScript
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzurePowerShellScript: %+v", err)
		}
		return out, nil
	}

	var parent BaseDeploymentScriptImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseDeploymentScriptImpl: %+v", err)
	}

	return RawDeploymentScriptImpl{
		deploymentScript: parent,
		Type:             value,
		Values:           temp,
	}, nil

}
