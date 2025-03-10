package performconnectivitycheck

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionStatus string

const (
	ConnectionStatusConnected    ConnectionStatus = "Connected"
	ConnectionStatusDegraded     ConnectionStatus = "Degraded"
	ConnectionStatusDisconnected ConnectionStatus = "Disconnected"
	ConnectionStatusUnknown      ConnectionStatus = "Unknown"
)

func PossibleValuesForConnectionStatus() []string {
	return []string{
		string(ConnectionStatusConnected),
		string(ConnectionStatusDegraded),
		string(ConnectionStatusDisconnected),
		string(ConnectionStatusUnknown),
	}
}

func (s *ConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionStatus(input string) (*ConnectionStatus, error) {
	vals := map[string]ConnectionStatus{
		"connected":    ConnectionStatusConnected,
		"degraded":     ConnectionStatusDegraded,
		"disconnected": ConnectionStatusDisconnected,
		"unknown":      ConnectionStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionStatus(input)
	return &out, nil
}

type ConnectivityCheckProtocol string

const (
	ConnectivityCheckProtocolHTTP  ConnectivityCheckProtocol = "HTTP"
	ConnectivityCheckProtocolHTTPS ConnectivityCheckProtocol = "HTTPS"
	ConnectivityCheckProtocolTCP   ConnectivityCheckProtocol = "TCP"
)

func PossibleValuesForConnectivityCheckProtocol() []string {
	return []string{
		string(ConnectivityCheckProtocolHTTP),
		string(ConnectivityCheckProtocolHTTPS),
		string(ConnectivityCheckProtocolTCP),
	}
}

func (s *ConnectivityCheckProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectivityCheckProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectivityCheckProtocol(input string) (*ConnectivityCheckProtocol, error) {
	vals := map[string]ConnectivityCheckProtocol{
		"http":  ConnectivityCheckProtocolHTTP,
		"https": ConnectivityCheckProtocolHTTPS,
		"tcp":   ConnectivityCheckProtocolTCP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityCheckProtocol(input)
	return &out, nil
}

type IssueType string

const (
	IssueTypeAgentStopped        IssueType = "AgentStopped"
	IssueTypeDnsResolution       IssueType = "DnsResolution"
	IssueTypeGuestFirewall       IssueType = "GuestFirewall"
	IssueTypeNetworkSecurityRule IssueType = "NetworkSecurityRule"
	IssueTypePlatform            IssueType = "Platform"
	IssueTypePortThrottled       IssueType = "PortThrottled"
	IssueTypeSocketBind          IssueType = "SocketBind"
	IssueTypeUnknown             IssueType = "Unknown"
	IssueTypeUserDefinedRoute    IssueType = "UserDefinedRoute"
)

func PossibleValuesForIssueType() []string {
	return []string{
		string(IssueTypeAgentStopped),
		string(IssueTypeDnsResolution),
		string(IssueTypeGuestFirewall),
		string(IssueTypeNetworkSecurityRule),
		string(IssueTypePlatform),
		string(IssueTypePortThrottled),
		string(IssueTypeSocketBind),
		string(IssueTypeUnknown),
		string(IssueTypeUserDefinedRoute),
	}
}

func (s *IssueType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIssueType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIssueType(input string) (*IssueType, error) {
	vals := map[string]IssueType{
		"agentstopped":        IssueTypeAgentStopped,
		"dnsresolution":       IssueTypeDnsResolution,
		"guestfirewall":       IssueTypeGuestFirewall,
		"networksecurityrule": IssueTypeNetworkSecurityRule,
		"platform":            IssueTypePlatform,
		"portthrottled":       IssueTypePortThrottled,
		"socketbind":          IssueTypeSocketBind,
		"unknown":             IssueTypeUnknown,
		"userdefinedroute":    IssueTypeUserDefinedRoute,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IssueType(input)
	return &out, nil
}

type Method string

const (
	MethodGET  Method = "GET"
	MethodPOST Method = "POST"
)

func PossibleValuesForMethod() []string {
	return []string{
		string(MethodGET),
		string(MethodPOST),
	}
}

func (s *Method) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMethod(input string) (*Method, error) {
	vals := map[string]Method{
		"get":  MethodGET,
		"post": MethodPOST,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Method(input)
	return &out, nil
}

type Origin string

const (
	OriginInbound  Origin = "Inbound"
	OriginLocal    Origin = "Local"
	OriginOutbound Origin = "Outbound"
)

func PossibleValuesForOrigin() []string {
	return []string{
		string(OriginInbound),
		string(OriginLocal),
		string(OriginOutbound),
	}
}

func (s *Origin) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOrigin(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOrigin(input string) (*Origin, error) {
	vals := map[string]Origin{
		"inbound":  OriginInbound,
		"local":    OriginLocal,
		"outbound": OriginOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Origin(input)
	return &out, nil
}

type PreferredIPVersion string

const (
	PreferredIPVersionIPvFour PreferredIPVersion = "IPv4"
)

func PossibleValuesForPreferredIPVersion() []string {
	return []string{
		string(PreferredIPVersionIPvFour),
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
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreferredIPVersion(input)
	return &out, nil
}

type Severity string

const (
	SeverityError   Severity = "Error"
	SeverityWarning Severity = "Warning"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeverityError),
		string(SeverityWarning),
	}
}

func (s *Severity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"error":   SeverityError,
		"warning": SeverityWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}
