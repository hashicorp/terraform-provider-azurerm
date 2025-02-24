package networkmanagereffectiveconnectivityconfiguration

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityTopology string

const (
	ConnectivityTopologyHubAndSpoke ConnectivityTopology = "HubAndSpoke"
	ConnectivityTopologyMesh        ConnectivityTopology = "Mesh"
)

func PossibleValuesForConnectivityTopology() []string {
	return []string{
		string(ConnectivityTopologyHubAndSpoke),
		string(ConnectivityTopologyMesh),
	}
}

func (s *ConnectivityTopology) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectivityTopology(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectivityTopology(input string) (*ConnectivityTopology, error) {
	vals := map[string]ConnectivityTopology{
		"hubandspoke": ConnectivityTopologyHubAndSpoke,
		"mesh":        ConnectivityTopologyMesh,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityTopology(input)
	return &out, nil
}

type DeleteExistingPeering string

const (
	DeleteExistingPeeringFalse DeleteExistingPeering = "False"
	DeleteExistingPeeringTrue  DeleteExistingPeering = "True"
)

func PossibleValuesForDeleteExistingPeering() []string {
	return []string{
		string(DeleteExistingPeeringFalse),
		string(DeleteExistingPeeringTrue),
	}
}

func (s *DeleteExistingPeering) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteExistingPeering(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteExistingPeering(input string) (*DeleteExistingPeering, error) {
	vals := map[string]DeleteExistingPeering{
		"false": DeleteExistingPeeringFalse,
		"true":  DeleteExistingPeeringTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteExistingPeering(input)
	return &out, nil
}

type GroupConnectivity string

const (
	GroupConnectivityDirectlyConnected GroupConnectivity = "DirectlyConnected"
	GroupConnectivityNone              GroupConnectivity = "None"
)

func PossibleValuesForGroupConnectivity() []string {
	return []string{
		string(GroupConnectivityDirectlyConnected),
		string(GroupConnectivityNone),
	}
}

func (s *GroupConnectivity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGroupConnectivity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGroupConnectivity(input string) (*GroupConnectivity, error) {
	vals := map[string]GroupConnectivity{
		"directlyconnected": GroupConnectivityDirectlyConnected,
		"none":              GroupConnectivityNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GroupConnectivity(input)
	return &out, nil
}

type GroupMemberType string

const (
	GroupMemberTypeSubnet         GroupMemberType = "Subnet"
	GroupMemberTypeVirtualNetwork GroupMemberType = "VirtualNetwork"
)

func PossibleValuesForGroupMemberType() []string {
	return []string{
		string(GroupMemberTypeSubnet),
		string(GroupMemberTypeVirtualNetwork),
	}
}

func (s *GroupMemberType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGroupMemberType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGroupMemberType(input string) (*GroupMemberType, error) {
	vals := map[string]GroupMemberType{
		"subnet":         GroupMemberTypeSubnet,
		"virtualnetwork": GroupMemberTypeVirtualNetwork,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GroupMemberType(input)
	return &out, nil
}

type IsGlobal string

const (
	IsGlobalFalse IsGlobal = "False"
	IsGlobalTrue  IsGlobal = "True"
)

func PossibleValuesForIsGlobal() []string {
	return []string{
		string(IsGlobalFalse),
		string(IsGlobalTrue),
	}
}

func (s *IsGlobal) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIsGlobal(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIsGlobal(input string) (*IsGlobal, error) {
	vals := map[string]IsGlobal{
		"false": IsGlobalFalse,
		"true":  IsGlobalTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IsGlobal(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type UseHubGateway string

const (
	UseHubGatewayFalse UseHubGateway = "False"
	UseHubGatewayTrue  UseHubGateway = "True"
)

func PossibleValuesForUseHubGateway() []string {
	return []string{
		string(UseHubGatewayFalse),
		string(UseHubGatewayTrue),
	}
}

func (s *UseHubGateway) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUseHubGateway(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUseHubGateway(input string) (*UseHubGateway, error) {
	vals := map[string]UseHubGateway{
		"false": UseHubGatewayFalse,
		"true":  UseHubGatewayTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UseHubGateway(input)
	return &out, nil
}
