package connectionanalyzers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionAnalyzerEndpointType string

const (
	ConnectionAnalyzerEndpointTypeApplicationGateway ConnectionAnalyzerEndpointType = "ApplicationGateway"
	ConnectionAnalyzerEndpointTypeBastionHost        ConnectionAnalyzerEndpointType = "BastionHost"
	ConnectionAnalyzerEndpointTypeExternalAddress    ConnectionAnalyzerEndpointType = "ExternalAddress"
	ConnectionAnalyzerEndpointTypeVM                 ConnectionAnalyzerEndpointType = "VM"
	ConnectionAnalyzerEndpointTypeVMSS               ConnectionAnalyzerEndpointType = "VMSS"
)

func PossibleValuesForConnectionAnalyzerEndpointType() []string {
	return []string{
		string(ConnectionAnalyzerEndpointTypeApplicationGateway),
		string(ConnectionAnalyzerEndpointTypeBastionHost),
		string(ConnectionAnalyzerEndpointTypeExternalAddress),
		string(ConnectionAnalyzerEndpointTypeVM),
		string(ConnectionAnalyzerEndpointTypeVMSS),
	}
}

func (s *ConnectionAnalyzerEndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionAnalyzerEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionAnalyzerEndpointType(input string) (*ConnectionAnalyzerEndpointType, error) {
	vals := map[string]ConnectionAnalyzerEndpointType{
		"applicationgateway": ConnectionAnalyzerEndpointTypeApplicationGateway,
		"bastionhost":        ConnectionAnalyzerEndpointTypeBastionHost,
		"externaladdress":    ConnectionAnalyzerEndpointTypeExternalAddress,
		"vm":                 ConnectionAnalyzerEndpointTypeVM,
		"vmss":               ConnectionAnalyzerEndpointTypeVMSS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionAnalyzerEndpointType(input)
	return &out, nil
}

type ConnectionAnalyzerStatus string

const (
	ConnectionAnalyzerStatusFailed     ConnectionAnalyzerStatus = "Failed"
	ConnectionAnalyzerStatusNotStarted ConnectionAnalyzerStatus = "NotStarted"
	ConnectionAnalyzerStatusRunning    ConnectionAnalyzerStatus = "Running"
	ConnectionAnalyzerStatusSucceeded  ConnectionAnalyzerStatus = "Succeeded"
)

func PossibleValuesForConnectionAnalyzerStatus() []string {
	return []string{
		string(ConnectionAnalyzerStatusFailed),
		string(ConnectionAnalyzerStatusNotStarted),
		string(ConnectionAnalyzerStatusRunning),
		string(ConnectionAnalyzerStatusSucceeded),
	}
}

func (s *ConnectionAnalyzerStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionAnalyzerStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionAnalyzerStatus(input string) (*ConnectionAnalyzerStatus, error) {
	vals := map[string]ConnectionAnalyzerStatus{
		"failed":     ConnectionAnalyzerStatusFailed,
		"notstarted": ConnectionAnalyzerStatusNotStarted,
		"running":    ConnectionAnalyzerStatusRunning,
		"succeeded":  ConnectionAnalyzerStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionAnalyzerStatus(input)
	return &out, nil
}

type DiagnosticOperation string

const (
	DiagnosticOperationConnectivityCheck      DiagnosticOperation = "ConnectivityCheck"
	DiagnosticOperationExpressRouteDiagnostic DiagnosticOperation = "ExpressRouteDiagnostic"
	DiagnosticOperationNSG                    DiagnosticOperation = "NSG"
	DiagnosticOperationNextHop                DiagnosticOperation = "NextHop"
	DiagnosticOperationPortScan               DiagnosticOperation = "PortScan"
)

func PossibleValuesForDiagnosticOperation() []string {
	return []string{
		string(DiagnosticOperationConnectivityCheck),
		string(DiagnosticOperationExpressRouteDiagnostic),
		string(DiagnosticOperationNSG),
		string(DiagnosticOperationNextHop),
		string(DiagnosticOperationPortScan),
	}
}

func (s *DiagnosticOperation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiagnosticOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiagnosticOperation(input string) (*DiagnosticOperation, error) {
	vals := map[string]DiagnosticOperation{
		"connectivitycheck":      DiagnosticOperationConnectivityCheck,
		"expressroutediagnostic": DiagnosticOperationExpressRouteDiagnostic,
		"nsg":                    DiagnosticOperationNSG,
		"nexthop":                DiagnosticOperationNextHop,
		"portscan":               DiagnosticOperationPortScan,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiagnosticOperation(input)
	return &out, nil
}

type HTTPMethod string

const (
	HTTPMethodGet HTTPMethod = "Get"
)

func PossibleValuesForHTTPMethod() []string {
	return []string{
		string(HTTPMethodGet),
	}
}

func (s *HTTPMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPMethod(input string) (*HTTPMethod, error) {
	vals := map[string]HTTPMethod{
		"get": HTTPMethodGet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPMethod(input)
	return &out, nil
}

type PreferredIPVersion string

const (
	PreferredIPVersionIPvFour PreferredIPVersion = "IPv4"
	PreferredIPVersionIPvSix  PreferredIPVersion = "IPv6"
)

func PossibleValuesForPreferredIPVersion() []string {
	return []string{
		string(PreferredIPVersionIPvFour),
		string(PreferredIPVersionIPvSix),
	}
}

func (s *PreferredIPVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreferredIPVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreferredIPVersion(input string) (*PreferredIPVersion, error) {
	vals := map[string]PreferredIPVersion{
		"ipv4": PreferredIPVersionIPvFour,
		"ipv6": PreferredIPVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreferredIPVersion(input)
	return &out, nil
}

type Protocol string

const (
	ProtocolHTTP  Protocol = "Http"
	ProtocolHTTPS Protocol = "Https"
	ProtocolIcmp  Protocol = "Icmp"
	ProtocolTcp   Protocol = "Tcp"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolHTTP),
		string(ProtocolHTTPS),
		string(ProtocolIcmp),
		string(ProtocolTcp),
	}
}

func (s *Protocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProtocol(input string) (*Protocol, error) {
	vals := map[string]Protocol{
		"http":  ProtocolHTTP,
		"https": ProtocolHTTPS,
		"icmp":  ProtocolIcmp,
		"tcp":   ProtocolTcp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Protocol(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
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
