package replicationfabrics

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ FabricSpecificDetails = InMageRcmFabricSpecificDetails{}

type InMageRcmFabricSpecificDetails struct {
	AgentDetails               *[]AgentDetails            `json:"agentDetails,omitempty"`
	ControlPlaneUri            *string                    `json:"controlPlaneUri,omitempty"`
	DataPlaneUri               *string                    `json:"dataPlaneUri,omitempty"`
	Dras                       *[]DraDetails              `json:"dras,omitempty"`
	MarsAgents                 *[]MarsAgentDetails        `json:"marsAgents,omitempty"`
	PhysicalSiteId             *string                    `json:"physicalSiteId,omitempty"`
	ProcessServers             *[]ProcessServerDetails    `json:"processServers,omitempty"`
	PushInstallers             *[]PushInstallerDetails    `json:"pushInstallers,omitempty"`
	RcmProxies                 *[]RcmProxyDetails         `json:"rcmProxies,omitempty"`
	ReplicationAgents          *[]ReplicationAgentDetails `json:"replicationAgents,omitempty"`
	ReprotectAgents            *[]ReprotectAgentDetails   `json:"reprotectAgents,omitempty"`
	ServiceContainerId         *string                    `json:"serviceContainerId,omitempty"`
	ServiceEndpoint            *string                    `json:"serviceEndpoint,omitempty"`
	ServiceResourceId          *string                    `json:"serviceResourceId,omitempty"`
	SourceAgentIdentityDetails *IdentityProviderDetails   `json:"sourceAgentIdentityDetails,omitempty"`
	VMwareSiteId               *string                    `json:"vmwareSiteId,omitempty"`

	// Fields inherited from FabricSpecificDetails
}

var _ json.Marshaler = InMageRcmFabricSpecificDetails{}

func (s InMageRcmFabricSpecificDetails) MarshalJSON() ([]byte, error) {
	type wrapper InMageRcmFabricSpecificDetails
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling InMageRcmFabricSpecificDetails: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling InMageRcmFabricSpecificDetails: %+v", err)
	}
	decoded["instanceType"] = "InMageRcm"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling InMageRcmFabricSpecificDetails: %+v", err)
	}

	return encoded, nil
}
