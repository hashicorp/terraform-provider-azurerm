package networkwatchers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Access string

const (
	AccessAllow Access = "Allow"
	AccessDeny  Access = "Deny"
)

func PossibleValuesForAccess() []string {
	return []string{
		string(AccessAllow),
		string(AccessDeny),
	}
}

func (s *Access) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccess(input string) (*Access, error) {
	vals := map[string]Access{
		"allow": AccessAllow,
		"deny":  AccessDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Access(input)
	return &out, nil
}

type AssociationType string

const (
	AssociationTypeAssociated AssociationType = "Associated"
	AssociationTypeContains   AssociationType = "Contains"
)

func PossibleValuesForAssociationType() []string {
	return []string{
		string(AssociationTypeAssociated),
		string(AssociationTypeContains),
	}
}

func (s *AssociationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssociationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssociationType(input string) (*AssociationType, error) {
	vals := map[string]AssociationType{
		"associated": AssociationTypeAssociated,
		"contains":   AssociationTypeContains,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssociationType(input)
	return &out, nil
}

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

type Direction string

const (
	DirectionInbound  Direction = "Inbound"
	DirectionOutbound Direction = "Outbound"
)

func PossibleValuesForDirection() []string {
	return []string{
		string(DirectionInbound),
		string(DirectionOutbound),
	}
}

func (s *Direction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDirection(input string) (*Direction, error) {
	vals := map[string]Direction{
		"inbound":  DirectionInbound,
		"outbound": DirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Direction(input)
	return &out, nil
}

type EffectiveSecurityRuleProtocol string

const (
	EffectiveSecurityRuleProtocolAll EffectiveSecurityRuleProtocol = "All"
	EffectiveSecurityRuleProtocolTcp EffectiveSecurityRuleProtocol = "Tcp"
	EffectiveSecurityRuleProtocolUdp EffectiveSecurityRuleProtocol = "Udp"
)

func PossibleValuesForEffectiveSecurityRuleProtocol() []string {
	return []string{
		string(EffectiveSecurityRuleProtocolAll),
		string(EffectiveSecurityRuleProtocolTcp),
		string(EffectiveSecurityRuleProtocolUdp),
	}
}

func (s *EffectiveSecurityRuleProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEffectiveSecurityRuleProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEffectiveSecurityRuleProtocol(input string) (*EffectiveSecurityRuleProtocol, error) {
	vals := map[string]EffectiveSecurityRuleProtocol{
		"all": EffectiveSecurityRuleProtocolAll,
		"tcp": EffectiveSecurityRuleProtocolTcp,
		"udp": EffectiveSecurityRuleProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EffectiveSecurityRuleProtocol(input)
	return &out, nil
}

type FlowLogFormatType string

const (
	FlowLogFormatTypeJSON FlowLogFormatType = "JSON"
)

func PossibleValuesForFlowLogFormatType() []string {
	return []string{
		string(FlowLogFormatTypeJSON),
	}
}

func (s *FlowLogFormatType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFlowLogFormatType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFlowLogFormatType(input string) (*FlowLogFormatType, error) {
	vals := map[string]FlowLogFormatType{
		"json": FlowLogFormatTypeJSON,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FlowLogFormatType(input)
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

type IPFlowProtocol string

const (
	IPFlowProtocolTCP IPFlowProtocol = "TCP"
	IPFlowProtocolUDP IPFlowProtocol = "UDP"
)

func PossibleValuesForIPFlowProtocol() []string {
	return []string{
		string(IPFlowProtocolTCP),
		string(IPFlowProtocolUDP),
	}
}

func (s *IPFlowProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPFlowProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPFlowProtocol(input string) (*IPFlowProtocol, error) {
	vals := map[string]IPFlowProtocol{
		"tcp": IPFlowProtocolTCP,
		"udp": IPFlowProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPFlowProtocol(input)
	return &out, nil
}

type IPVersion string

const (
	IPVersionIPvFour IPVersion = "IPv4"
	IPVersionIPvSix  IPVersion = "IPv6"
)

func PossibleValuesForIPVersion() []string {
	return []string{
		string(IPVersionIPvFour),
		string(IPVersionIPvSix),
	}
}

func (s *IPVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPVersion(input string) (*IPVersion, error) {
	vals := map[string]IPVersion{
		"ipv4": IPVersionIPvFour,
		"ipv6": IPVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPVersion(input)
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

type NextHopType string

const (
	NextHopTypeHyperNetGateway       NextHopType = "HyperNetGateway"
	NextHopTypeInternet              NextHopType = "Internet"
	NextHopTypeNone                  NextHopType = "None"
	NextHopTypeVirtualAppliance      NextHopType = "VirtualAppliance"
	NextHopTypeVirtualNetworkGateway NextHopType = "VirtualNetworkGateway"
	NextHopTypeVnetLocal             NextHopType = "VnetLocal"
)

func PossibleValuesForNextHopType() []string {
	return []string{
		string(NextHopTypeHyperNetGateway),
		string(NextHopTypeInternet),
		string(NextHopTypeNone),
		string(NextHopTypeVirtualAppliance),
		string(NextHopTypeVirtualNetworkGateway),
		string(NextHopTypeVnetLocal),
	}
}

func (s *NextHopType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNextHopType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNextHopType(input string) (*NextHopType, error) {
	vals := map[string]NextHopType{
		"hypernetgateway":       NextHopTypeHyperNetGateway,
		"internet":              NextHopTypeInternet,
		"none":                  NextHopTypeNone,
		"virtualappliance":      NextHopTypeVirtualAppliance,
		"virtualnetworkgateway": NextHopTypeVirtualNetworkGateway,
		"vnetlocal":             NextHopTypeVnetLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NextHopType(input)
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

type SecurityRuleAccess string

const (
	SecurityRuleAccessAllow SecurityRuleAccess = "Allow"
	SecurityRuleAccessDeny  SecurityRuleAccess = "Deny"
)

func PossibleValuesForSecurityRuleAccess() []string {
	return []string{
		string(SecurityRuleAccessAllow),
		string(SecurityRuleAccessDeny),
	}
}

func (s *SecurityRuleAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityRuleAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityRuleAccess(input string) (*SecurityRuleAccess, error) {
	vals := map[string]SecurityRuleAccess{
		"allow": SecurityRuleAccessAllow,
		"deny":  SecurityRuleAccessDeny,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityRuleAccess(input)
	return &out, nil
}

type SecurityRuleDirection string

const (
	SecurityRuleDirectionInbound  SecurityRuleDirection = "Inbound"
	SecurityRuleDirectionOutbound SecurityRuleDirection = "Outbound"
)

func PossibleValuesForSecurityRuleDirection() []string {
	return []string{
		string(SecurityRuleDirectionInbound),
		string(SecurityRuleDirectionOutbound),
	}
}

func (s *SecurityRuleDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityRuleDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityRuleDirection(input string) (*SecurityRuleDirection, error) {
	vals := map[string]SecurityRuleDirection{
		"inbound":  SecurityRuleDirectionInbound,
		"outbound": SecurityRuleDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityRuleDirection(input)
	return &out, nil
}

type SecurityRuleProtocol string

const (
	SecurityRuleProtocolAh   SecurityRuleProtocol = "Ah"
	SecurityRuleProtocolAny  SecurityRuleProtocol = "*"
	SecurityRuleProtocolEsp  SecurityRuleProtocol = "Esp"
	SecurityRuleProtocolIcmp SecurityRuleProtocol = "Icmp"
	SecurityRuleProtocolTcp  SecurityRuleProtocol = "Tcp"
	SecurityRuleProtocolUdp  SecurityRuleProtocol = "Udp"
)

func PossibleValuesForSecurityRuleProtocol() []string {
	return []string{
		string(SecurityRuleProtocolAh),
		string(SecurityRuleProtocolAny),
		string(SecurityRuleProtocolEsp),
		string(SecurityRuleProtocolIcmp),
		string(SecurityRuleProtocolTcp),
		string(SecurityRuleProtocolUdp),
	}
}

func (s *SecurityRuleProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityRuleProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityRuleProtocol(input string) (*SecurityRuleProtocol, error) {
	vals := map[string]SecurityRuleProtocol{
		"ah":   SecurityRuleProtocolAh,
		"*":    SecurityRuleProtocolAny,
		"esp":  SecurityRuleProtocolEsp,
		"icmp": SecurityRuleProtocolIcmp,
		"tcp":  SecurityRuleProtocolTcp,
		"udp":  SecurityRuleProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityRuleProtocol(input)
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

type VerbosityLevel string

const (
	VerbosityLevelFull    VerbosityLevel = "Full"
	VerbosityLevelMinimum VerbosityLevel = "Minimum"
	VerbosityLevelNormal  VerbosityLevel = "Normal"
)

func PossibleValuesForVerbosityLevel() []string {
	return []string{
		string(VerbosityLevelFull),
		string(VerbosityLevelMinimum),
		string(VerbosityLevelNormal),
	}
}

func (s *VerbosityLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVerbosityLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVerbosityLevel(input string) (*VerbosityLevel, error) {
	vals := map[string]VerbosityLevel{
		"full":    VerbosityLevelFull,
		"minimum": VerbosityLevelMinimum,
		"normal":  VerbosityLevelNormal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VerbosityLevel(input)
	return &out, nil
}
