package managedcassandras

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMethod string

const (
	AuthenticationMethodCassandra AuthenticationMethod = "Cassandra"
	AuthenticationMethodLdap      AuthenticationMethod = "Ldap"
	AuthenticationMethodNone      AuthenticationMethod = "None"
)

func PossibleValuesForAuthenticationMethod() []string {
	return []string{
		string(AuthenticationMethodCassandra),
		string(AuthenticationMethodLdap),
		string(AuthenticationMethodNone),
	}
}

func (s *AuthenticationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthenticationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthenticationMethod(input string) (*AuthenticationMethod, error) {
	vals := map[string]AuthenticationMethod{
		"cassandra": AuthenticationMethodCassandra,
		"ldap":      AuthenticationMethodLdap,
		"none":      AuthenticationMethodNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMethod(input)
	return &out, nil
}

type ConnectionState string

const (
	ConnectionStateDatacenterToDatacenterNetworkError           ConnectionState = "DatacenterToDatacenterNetworkError"
	ConnectionStateInternalError                                ConnectionState = "InternalError"
	ConnectionStateInternalOperatorToDataCenterCertificateError ConnectionState = "InternalOperatorToDataCenterCertificateError"
	ConnectionStateOK                                           ConnectionState = "OK"
	ConnectionStateOperatorToDataCenterNetworkError             ConnectionState = "OperatorToDataCenterNetworkError"
	ConnectionStateUnknown                                      ConnectionState = "Unknown"
)

func PossibleValuesForConnectionState() []string {
	return []string{
		string(ConnectionStateDatacenterToDatacenterNetworkError),
		string(ConnectionStateInternalError),
		string(ConnectionStateInternalOperatorToDataCenterCertificateError),
		string(ConnectionStateOK),
		string(ConnectionStateOperatorToDataCenterNetworkError),
		string(ConnectionStateUnknown),
	}
}

func (s *ConnectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionState(input string) (*ConnectionState, error) {
	vals := map[string]ConnectionState{
		"datacentertodatacenternetworkerror":           ConnectionStateDatacenterToDatacenterNetworkError,
		"internalerror":                                ConnectionStateInternalError,
		"internaloperatortodatacentercertificateerror": ConnectionStateInternalOperatorToDataCenterCertificateError,
		"ok":                               ConnectionStateOK,
		"operatortodatacenternetworkerror": ConnectionStateOperatorToDataCenterNetworkError,
		"unknown":                          ConnectionStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionState(input)
	return &out, nil
}

type ManagedCassandraProvisioningState string

const (
	ManagedCassandraProvisioningStateCanceled  ManagedCassandraProvisioningState = "Canceled"
	ManagedCassandraProvisioningStateCreating  ManagedCassandraProvisioningState = "Creating"
	ManagedCassandraProvisioningStateDeleting  ManagedCassandraProvisioningState = "Deleting"
	ManagedCassandraProvisioningStateFailed    ManagedCassandraProvisioningState = "Failed"
	ManagedCassandraProvisioningStateSucceeded ManagedCassandraProvisioningState = "Succeeded"
	ManagedCassandraProvisioningStateUpdating  ManagedCassandraProvisioningState = "Updating"
)

func PossibleValuesForManagedCassandraProvisioningState() []string {
	return []string{
		string(ManagedCassandraProvisioningStateCanceled),
		string(ManagedCassandraProvisioningStateCreating),
		string(ManagedCassandraProvisioningStateDeleting),
		string(ManagedCassandraProvisioningStateFailed),
		string(ManagedCassandraProvisioningStateSucceeded),
		string(ManagedCassandraProvisioningStateUpdating),
	}
}

func (s *ManagedCassandraProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedCassandraProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedCassandraProvisioningState(input string) (*ManagedCassandraProvisioningState, error) {
	vals := map[string]ManagedCassandraProvisioningState{
		"canceled":  ManagedCassandraProvisioningStateCanceled,
		"creating":  ManagedCassandraProvisioningStateCreating,
		"deleting":  ManagedCassandraProvisioningStateDeleting,
		"failed":    ManagedCassandraProvisioningStateFailed,
		"succeeded": ManagedCassandraProvisioningStateSucceeded,
		"updating":  ManagedCassandraProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedCassandraProvisioningState(input)
	return &out, nil
}

type NodeState string

const (
	NodeStateJoining NodeState = "Joining"
	NodeStateLeaving NodeState = "Leaving"
	NodeStateMoving  NodeState = "Moving"
	NodeStateNormal  NodeState = "Normal"
	NodeStateStopped NodeState = "Stopped"
)

func PossibleValuesForNodeState() []string {
	return []string{
		string(NodeStateJoining),
		string(NodeStateLeaving),
		string(NodeStateMoving),
		string(NodeStateNormal),
		string(NodeStateStopped),
	}
}

func (s *NodeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeState(input string) (*NodeState, error) {
	vals := map[string]NodeState{
		"joining": NodeStateJoining,
		"leaving": NodeStateLeaving,
		"moving":  NodeStateMoving,
		"normal":  NodeStateNormal,
		"stopped": NodeStateStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeState(input)
	return &out, nil
}
