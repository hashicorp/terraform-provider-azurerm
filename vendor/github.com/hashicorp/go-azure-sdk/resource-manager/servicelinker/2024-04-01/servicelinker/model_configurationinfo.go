package servicelinker

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationInfo struct {
	Action                               *ActionType             `json:"action,omitempty"`
	AdditionalConfigurations             *map[string]string      `json:"additionalConfigurations,omitempty"`
	AdditionalConnectionStringProperties *map[string]string      `json:"additionalConnectionStringProperties,omitempty"`
	ConfigurationStore                   *ConfigurationStore     `json:"configurationStore,omitempty"`
	CustomizedKeys                       *map[string]string      `json:"customizedKeys,omitempty"`
	DaprProperties                       *DaprProperties         `json:"daprProperties,omitempty"`
	DeleteOrUpdateBehavior               *DeleteOrUpdateBehavior `json:"deleteOrUpdateBehavior,omitempty"`
}
