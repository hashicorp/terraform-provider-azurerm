package loadbalancers

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DdosSettingsProtectionMode string

const (
	DdosSettingsProtectionModeDisabled                DdosSettingsProtectionMode = "Disabled"
	DdosSettingsProtectionModeEnabled                 DdosSettingsProtectionMode = "Enabled"
	DdosSettingsProtectionModeVirtualNetworkInherited DdosSettingsProtectionMode = "VirtualNetworkInherited"
)

func PossibleValuesForDdosSettingsProtectionMode() []string {
	return []string{
		string(DdosSettingsProtectionModeDisabled),
		string(DdosSettingsProtectionModeEnabled),
		string(DdosSettingsProtectionModeVirtualNetworkInherited),
	}
}

func (s *DdosSettingsProtectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDdosSettingsProtectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDdosSettingsProtectionMode(input string) (*DdosSettingsProtectionMode, error) {
	vals := map[string]DdosSettingsProtectionMode{
		"disabled":                DdosSettingsProtectionModeDisabled,
		"enabled":                 DdosSettingsProtectionModeEnabled,
		"virtualnetworkinherited": DdosSettingsProtectionModeVirtualNetworkInherited,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DdosSettingsProtectionMode(input)
	return &out, nil
}

type DeleteOptions string

const (
	DeleteOptionsDelete DeleteOptions = "Delete"
	DeleteOptionsDetach DeleteOptions = "Detach"
)

func PossibleValuesForDeleteOptions() []string {
	return []string{
		string(DeleteOptionsDelete),
		string(DeleteOptionsDetach),
	}
}

func (s *DeleteOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteOptions(input string) (*DeleteOptions, error) {
	vals := map[string]DeleteOptions{
		"delete": DeleteOptionsDelete,
		"detach": DeleteOptionsDetach,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteOptions(input)
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

type GatewayLoadBalancerTunnelInterfaceType string

const (
	GatewayLoadBalancerTunnelInterfaceTypeExternal GatewayLoadBalancerTunnelInterfaceType = "External"
	GatewayLoadBalancerTunnelInterfaceTypeInternal GatewayLoadBalancerTunnelInterfaceType = "Internal"
	GatewayLoadBalancerTunnelInterfaceTypeNone     GatewayLoadBalancerTunnelInterfaceType = "None"
)

func PossibleValuesForGatewayLoadBalancerTunnelInterfaceType() []string {
	return []string{
		string(GatewayLoadBalancerTunnelInterfaceTypeExternal),
		string(GatewayLoadBalancerTunnelInterfaceTypeInternal),
		string(GatewayLoadBalancerTunnelInterfaceTypeNone),
	}
}

func (s *GatewayLoadBalancerTunnelInterfaceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayLoadBalancerTunnelInterfaceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayLoadBalancerTunnelInterfaceType(input string) (*GatewayLoadBalancerTunnelInterfaceType, error) {
	vals := map[string]GatewayLoadBalancerTunnelInterfaceType{
		"external": GatewayLoadBalancerTunnelInterfaceTypeExternal,
		"internal": GatewayLoadBalancerTunnelInterfaceTypeInternal,
		"none":     GatewayLoadBalancerTunnelInterfaceTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayLoadBalancerTunnelInterfaceType(input)
	return &out, nil
}

type GatewayLoadBalancerTunnelProtocol string

const (
	GatewayLoadBalancerTunnelProtocolNative GatewayLoadBalancerTunnelProtocol = "Native"
	GatewayLoadBalancerTunnelProtocolNone   GatewayLoadBalancerTunnelProtocol = "None"
	GatewayLoadBalancerTunnelProtocolVXLAN  GatewayLoadBalancerTunnelProtocol = "VXLAN"
)

func PossibleValuesForGatewayLoadBalancerTunnelProtocol() []string {
	return []string{
		string(GatewayLoadBalancerTunnelProtocolNative),
		string(GatewayLoadBalancerTunnelProtocolNone),
		string(GatewayLoadBalancerTunnelProtocolVXLAN),
	}
}

func (s *GatewayLoadBalancerTunnelProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayLoadBalancerTunnelProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayLoadBalancerTunnelProtocol(input string) (*GatewayLoadBalancerTunnelProtocol, error) {
	vals := map[string]GatewayLoadBalancerTunnelProtocol{
		"native": GatewayLoadBalancerTunnelProtocolNative,
		"none":   GatewayLoadBalancerTunnelProtocolNone,
		"vxlan":  GatewayLoadBalancerTunnelProtocolVXLAN,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayLoadBalancerTunnelProtocol(input)
	return &out, nil
}

type IPAllocationMethod string

const (
	IPAllocationMethodDynamic IPAllocationMethod = "Dynamic"
	IPAllocationMethodStatic  IPAllocationMethod = "Static"
)

func PossibleValuesForIPAllocationMethod() []string {
	return []string{
		string(IPAllocationMethodDynamic),
		string(IPAllocationMethodStatic),
	}
}

func (s *IPAllocationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPAllocationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPAllocationMethod(input string) (*IPAllocationMethod, error) {
	vals := map[string]IPAllocationMethod{
		"dynamic": IPAllocationMethodDynamic,
		"static":  IPAllocationMethodStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPAllocationMethod(input)
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

type LoadBalancerBackendAddressAdminState string

const (
	LoadBalancerBackendAddressAdminStateDown LoadBalancerBackendAddressAdminState = "Down"
	LoadBalancerBackendAddressAdminStateNone LoadBalancerBackendAddressAdminState = "None"
	LoadBalancerBackendAddressAdminStateUp   LoadBalancerBackendAddressAdminState = "Up"
)

func PossibleValuesForLoadBalancerBackendAddressAdminState() []string {
	return []string{
		string(LoadBalancerBackendAddressAdminStateDown),
		string(LoadBalancerBackendAddressAdminStateNone),
		string(LoadBalancerBackendAddressAdminStateUp),
	}
}

func (s *LoadBalancerBackendAddressAdminState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerBackendAddressAdminState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerBackendAddressAdminState(input string) (*LoadBalancerBackendAddressAdminState, error) {
	vals := map[string]LoadBalancerBackendAddressAdminState{
		"down": LoadBalancerBackendAddressAdminStateDown,
		"none": LoadBalancerBackendAddressAdminStateNone,
		"up":   LoadBalancerBackendAddressAdminStateUp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerBackendAddressAdminState(input)
	return &out, nil
}

type LoadBalancerOutboundRuleProtocol string

const (
	LoadBalancerOutboundRuleProtocolAll LoadBalancerOutboundRuleProtocol = "All"
	LoadBalancerOutboundRuleProtocolTcp LoadBalancerOutboundRuleProtocol = "Tcp"
	LoadBalancerOutboundRuleProtocolUdp LoadBalancerOutboundRuleProtocol = "Udp"
)

func PossibleValuesForLoadBalancerOutboundRuleProtocol() []string {
	return []string{
		string(LoadBalancerOutboundRuleProtocolAll),
		string(LoadBalancerOutboundRuleProtocolTcp),
		string(LoadBalancerOutboundRuleProtocolUdp),
	}
}

func (s *LoadBalancerOutboundRuleProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerOutboundRuleProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerOutboundRuleProtocol(input string) (*LoadBalancerOutboundRuleProtocol, error) {
	vals := map[string]LoadBalancerOutboundRuleProtocol{
		"all": LoadBalancerOutboundRuleProtocolAll,
		"tcp": LoadBalancerOutboundRuleProtocolTcp,
		"udp": LoadBalancerOutboundRuleProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerOutboundRuleProtocol(input)
	return &out, nil
}

type LoadBalancerSkuName string

const (
	LoadBalancerSkuNameBasic    LoadBalancerSkuName = "Basic"
	LoadBalancerSkuNameGateway  LoadBalancerSkuName = "Gateway"
	LoadBalancerSkuNameStandard LoadBalancerSkuName = "Standard"
)

func PossibleValuesForLoadBalancerSkuName() []string {
	return []string{
		string(LoadBalancerSkuNameBasic),
		string(LoadBalancerSkuNameGateway),
		string(LoadBalancerSkuNameStandard),
	}
}

func (s *LoadBalancerSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerSkuName(input string) (*LoadBalancerSkuName, error) {
	vals := map[string]LoadBalancerSkuName{
		"basic":    LoadBalancerSkuNameBasic,
		"gateway":  LoadBalancerSkuNameGateway,
		"standard": LoadBalancerSkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerSkuName(input)
	return &out, nil
}

type LoadBalancerSkuTier string

const (
	LoadBalancerSkuTierGlobal   LoadBalancerSkuTier = "Global"
	LoadBalancerSkuTierRegional LoadBalancerSkuTier = "Regional"
)

func PossibleValuesForLoadBalancerSkuTier() []string {
	return []string{
		string(LoadBalancerSkuTierGlobal),
		string(LoadBalancerSkuTierRegional),
	}
}

func (s *LoadBalancerSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerSkuTier(input string) (*LoadBalancerSkuTier, error) {
	vals := map[string]LoadBalancerSkuTier{
		"global":   LoadBalancerSkuTierGlobal,
		"regional": LoadBalancerSkuTierRegional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerSkuTier(input)
	return &out, nil
}

type LoadDistribution string

const (
	LoadDistributionDefault          LoadDistribution = "Default"
	LoadDistributionSourceIP         LoadDistribution = "SourceIP"
	LoadDistributionSourceIPProtocol LoadDistribution = "SourceIPProtocol"
)

func PossibleValuesForLoadDistribution() []string {
	return []string{
		string(LoadDistributionDefault),
		string(LoadDistributionSourceIP),
		string(LoadDistributionSourceIPProtocol),
	}
}

func (s *LoadDistribution) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadDistribution(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadDistribution(input string) (*LoadDistribution, error) {
	vals := map[string]LoadDistribution{
		"default":          LoadDistributionDefault,
		"sourceip":         LoadDistributionSourceIP,
		"sourceipprotocol": LoadDistributionSourceIPProtocol,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadDistribution(input)
	return &out, nil
}

type NatGatewaySkuName string

const (
	NatGatewaySkuNameStandard NatGatewaySkuName = "Standard"
)

func PossibleValuesForNatGatewaySkuName() []string {
	return []string{
		string(NatGatewaySkuNameStandard),
	}
}

func (s *NatGatewaySkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNatGatewaySkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNatGatewaySkuName(input string) (*NatGatewaySkuName, error) {
	vals := map[string]NatGatewaySkuName{
		"standard": NatGatewaySkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NatGatewaySkuName(input)
	return &out, nil
}

type NetworkInterfaceAuxiliaryMode string

const (
	NetworkInterfaceAuxiliaryModeAcceleratedConnections NetworkInterfaceAuxiliaryMode = "AcceleratedConnections"
	NetworkInterfaceAuxiliaryModeFloating               NetworkInterfaceAuxiliaryMode = "Floating"
	NetworkInterfaceAuxiliaryModeMaxConnections         NetworkInterfaceAuxiliaryMode = "MaxConnections"
	NetworkInterfaceAuxiliaryModeNone                   NetworkInterfaceAuxiliaryMode = "None"
)

func PossibleValuesForNetworkInterfaceAuxiliaryMode() []string {
	return []string{
		string(NetworkInterfaceAuxiliaryModeAcceleratedConnections),
		string(NetworkInterfaceAuxiliaryModeFloating),
		string(NetworkInterfaceAuxiliaryModeMaxConnections),
		string(NetworkInterfaceAuxiliaryModeNone),
	}
}

func (s *NetworkInterfaceAuxiliaryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceAuxiliaryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceAuxiliaryMode(input string) (*NetworkInterfaceAuxiliaryMode, error) {
	vals := map[string]NetworkInterfaceAuxiliaryMode{
		"acceleratedconnections": NetworkInterfaceAuxiliaryModeAcceleratedConnections,
		"floating":               NetworkInterfaceAuxiliaryModeFloating,
		"maxconnections":         NetworkInterfaceAuxiliaryModeMaxConnections,
		"none":                   NetworkInterfaceAuxiliaryModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceAuxiliaryMode(input)
	return &out, nil
}

type NetworkInterfaceAuxiliarySku string

const (
	NetworkInterfaceAuxiliarySkuAEight NetworkInterfaceAuxiliarySku = "A8"
	NetworkInterfaceAuxiliarySkuAFour  NetworkInterfaceAuxiliarySku = "A4"
	NetworkInterfaceAuxiliarySkuAOne   NetworkInterfaceAuxiliarySku = "A1"
	NetworkInterfaceAuxiliarySkuATwo   NetworkInterfaceAuxiliarySku = "A2"
	NetworkInterfaceAuxiliarySkuNone   NetworkInterfaceAuxiliarySku = "None"
)

func PossibleValuesForNetworkInterfaceAuxiliarySku() []string {
	return []string{
		string(NetworkInterfaceAuxiliarySkuAEight),
		string(NetworkInterfaceAuxiliarySkuAFour),
		string(NetworkInterfaceAuxiliarySkuAOne),
		string(NetworkInterfaceAuxiliarySkuATwo),
		string(NetworkInterfaceAuxiliarySkuNone),
	}
}

func (s *NetworkInterfaceAuxiliarySku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceAuxiliarySku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceAuxiliarySku(input string) (*NetworkInterfaceAuxiliarySku, error) {
	vals := map[string]NetworkInterfaceAuxiliarySku{
		"a8":   NetworkInterfaceAuxiliarySkuAEight,
		"a4":   NetworkInterfaceAuxiliarySkuAFour,
		"a1":   NetworkInterfaceAuxiliarySkuAOne,
		"a2":   NetworkInterfaceAuxiliarySkuATwo,
		"none": NetworkInterfaceAuxiliarySkuNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceAuxiliarySku(input)
	return &out, nil
}

type NetworkInterfaceMigrationPhase string

const (
	NetworkInterfaceMigrationPhaseAbort     NetworkInterfaceMigrationPhase = "Abort"
	NetworkInterfaceMigrationPhaseCommit    NetworkInterfaceMigrationPhase = "Commit"
	NetworkInterfaceMigrationPhaseCommitted NetworkInterfaceMigrationPhase = "Committed"
	NetworkInterfaceMigrationPhaseNone      NetworkInterfaceMigrationPhase = "None"
	NetworkInterfaceMigrationPhasePrepare   NetworkInterfaceMigrationPhase = "Prepare"
)

func PossibleValuesForNetworkInterfaceMigrationPhase() []string {
	return []string{
		string(NetworkInterfaceMigrationPhaseAbort),
		string(NetworkInterfaceMigrationPhaseCommit),
		string(NetworkInterfaceMigrationPhaseCommitted),
		string(NetworkInterfaceMigrationPhaseNone),
		string(NetworkInterfaceMigrationPhasePrepare),
	}
}

func (s *NetworkInterfaceMigrationPhase) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceMigrationPhase(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceMigrationPhase(input string) (*NetworkInterfaceMigrationPhase, error) {
	vals := map[string]NetworkInterfaceMigrationPhase{
		"abort":     NetworkInterfaceMigrationPhaseAbort,
		"commit":    NetworkInterfaceMigrationPhaseCommit,
		"committed": NetworkInterfaceMigrationPhaseCommitted,
		"none":      NetworkInterfaceMigrationPhaseNone,
		"prepare":   NetworkInterfaceMigrationPhasePrepare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceMigrationPhase(input)
	return &out, nil
}

type NetworkInterfaceNicType string

const (
	NetworkInterfaceNicTypeElastic  NetworkInterfaceNicType = "Elastic"
	NetworkInterfaceNicTypeStandard NetworkInterfaceNicType = "Standard"
)

func PossibleValuesForNetworkInterfaceNicType() []string {
	return []string{
		string(NetworkInterfaceNicTypeElastic),
		string(NetworkInterfaceNicTypeStandard),
	}
}

func (s *NetworkInterfaceNicType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceNicType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceNicType(input string) (*NetworkInterfaceNicType, error) {
	vals := map[string]NetworkInterfaceNicType{
		"elastic":  NetworkInterfaceNicTypeElastic,
		"standard": NetworkInterfaceNicTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceNicType(input)
	return &out, nil
}

type ProbeProtocol string

const (
	ProbeProtocolHTTP  ProbeProtocol = "Http"
	ProbeProtocolHTTPS ProbeProtocol = "Https"
	ProbeProtocolTcp   ProbeProtocol = "Tcp"
)

func PossibleValuesForProbeProtocol() []string {
	return []string{
		string(ProbeProtocolHTTP),
		string(ProbeProtocolHTTPS),
		string(ProbeProtocolTcp),
	}
}

func (s *ProbeProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProbeProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProbeProtocol(input string) (*ProbeProtocol, error) {
	vals := map[string]ProbeProtocol{
		"http":  ProbeProtocolHTTP,
		"https": ProbeProtocolHTTPS,
		"tcp":   ProbeProtocolTcp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProbeProtocol(input)
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

type PublicIPAddressDnsSettingsDomainNameLabelScope string

const (
	PublicIPAddressDnsSettingsDomainNameLabelScopeNoReuse            PublicIPAddressDnsSettingsDomainNameLabelScope = "NoReuse"
	PublicIPAddressDnsSettingsDomainNameLabelScopeResourceGroupReuse PublicIPAddressDnsSettingsDomainNameLabelScope = "ResourceGroupReuse"
	PublicIPAddressDnsSettingsDomainNameLabelScopeSubscriptionReuse  PublicIPAddressDnsSettingsDomainNameLabelScope = "SubscriptionReuse"
	PublicIPAddressDnsSettingsDomainNameLabelScopeTenantReuse        PublicIPAddressDnsSettingsDomainNameLabelScope = "TenantReuse"
)

func PossibleValuesForPublicIPAddressDnsSettingsDomainNameLabelScope() []string {
	return []string{
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeNoReuse),
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeResourceGroupReuse),
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeSubscriptionReuse),
		string(PublicIPAddressDnsSettingsDomainNameLabelScopeTenantReuse),
	}
}

func (s *PublicIPAddressDnsSettingsDomainNameLabelScope) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressDnsSettingsDomainNameLabelScope(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressDnsSettingsDomainNameLabelScope(input string) (*PublicIPAddressDnsSettingsDomainNameLabelScope, error) {
	vals := map[string]PublicIPAddressDnsSettingsDomainNameLabelScope{
		"noreuse":            PublicIPAddressDnsSettingsDomainNameLabelScopeNoReuse,
		"resourcegroupreuse": PublicIPAddressDnsSettingsDomainNameLabelScopeResourceGroupReuse,
		"subscriptionreuse":  PublicIPAddressDnsSettingsDomainNameLabelScopeSubscriptionReuse,
		"tenantreuse":        PublicIPAddressDnsSettingsDomainNameLabelScopeTenantReuse,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressDnsSettingsDomainNameLabelScope(input)
	return &out, nil
}

type PublicIPAddressMigrationPhase string

const (
	PublicIPAddressMigrationPhaseAbort     PublicIPAddressMigrationPhase = "Abort"
	PublicIPAddressMigrationPhaseCommit    PublicIPAddressMigrationPhase = "Commit"
	PublicIPAddressMigrationPhaseCommitted PublicIPAddressMigrationPhase = "Committed"
	PublicIPAddressMigrationPhaseNone      PublicIPAddressMigrationPhase = "None"
	PublicIPAddressMigrationPhasePrepare   PublicIPAddressMigrationPhase = "Prepare"
)

func PossibleValuesForPublicIPAddressMigrationPhase() []string {
	return []string{
		string(PublicIPAddressMigrationPhaseAbort),
		string(PublicIPAddressMigrationPhaseCommit),
		string(PublicIPAddressMigrationPhaseCommitted),
		string(PublicIPAddressMigrationPhaseNone),
		string(PublicIPAddressMigrationPhasePrepare),
	}
}

func (s *PublicIPAddressMigrationPhase) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressMigrationPhase(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressMigrationPhase(input string) (*PublicIPAddressMigrationPhase, error) {
	vals := map[string]PublicIPAddressMigrationPhase{
		"abort":     PublicIPAddressMigrationPhaseAbort,
		"commit":    PublicIPAddressMigrationPhaseCommit,
		"committed": PublicIPAddressMigrationPhaseCommitted,
		"none":      PublicIPAddressMigrationPhaseNone,
		"prepare":   PublicIPAddressMigrationPhasePrepare,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressMigrationPhase(input)
	return &out, nil
}

type PublicIPAddressSkuName string

const (
	PublicIPAddressSkuNameBasic    PublicIPAddressSkuName = "Basic"
	PublicIPAddressSkuNameStandard PublicIPAddressSkuName = "Standard"
)

func PossibleValuesForPublicIPAddressSkuName() []string {
	return []string{
		string(PublicIPAddressSkuNameBasic),
		string(PublicIPAddressSkuNameStandard),
	}
}

func (s *PublicIPAddressSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressSkuName(input string) (*PublicIPAddressSkuName, error) {
	vals := map[string]PublicIPAddressSkuName{
		"basic":    PublicIPAddressSkuNameBasic,
		"standard": PublicIPAddressSkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressSkuName(input)
	return &out, nil
}

type PublicIPAddressSkuTier string

const (
	PublicIPAddressSkuTierGlobal   PublicIPAddressSkuTier = "Global"
	PublicIPAddressSkuTierRegional PublicIPAddressSkuTier = "Regional"
)

func PossibleValuesForPublicIPAddressSkuTier() []string {
	return []string{
		string(PublicIPAddressSkuTierGlobal),
		string(PublicIPAddressSkuTierRegional),
	}
}

func (s *PublicIPAddressSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressSkuTier(input string) (*PublicIPAddressSkuTier, error) {
	vals := map[string]PublicIPAddressSkuTier{
		"global":   PublicIPAddressSkuTierGlobal,
		"regional": PublicIPAddressSkuTierRegional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressSkuTier(input)
	return &out, nil
}

type RouteNextHopType string

const (
	RouteNextHopTypeInternet              RouteNextHopType = "Internet"
	RouteNextHopTypeNone                  RouteNextHopType = "None"
	RouteNextHopTypeVirtualAppliance      RouteNextHopType = "VirtualAppliance"
	RouteNextHopTypeVirtualNetworkGateway RouteNextHopType = "VirtualNetworkGateway"
	RouteNextHopTypeVnetLocal             RouteNextHopType = "VnetLocal"
)

func PossibleValuesForRouteNextHopType() []string {
	return []string{
		string(RouteNextHopTypeInternet),
		string(RouteNextHopTypeNone),
		string(RouteNextHopTypeVirtualAppliance),
		string(RouteNextHopTypeVirtualNetworkGateway),
		string(RouteNextHopTypeVnetLocal),
	}
}

func (s *RouteNextHopType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRouteNextHopType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRouteNextHopType(input string) (*RouteNextHopType, error) {
	vals := map[string]RouteNextHopType{
		"internet":              RouteNextHopTypeInternet,
		"none":                  RouteNextHopTypeNone,
		"virtualappliance":      RouteNextHopTypeVirtualAppliance,
		"virtualnetworkgateway": RouteNextHopTypeVirtualNetworkGateway,
		"vnetlocal":             RouteNextHopTypeVnetLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RouteNextHopType(input)
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

type SyncMode string

const (
	SyncModeAutomatic SyncMode = "Automatic"
	SyncModeManual    SyncMode = "Manual"
)

func PossibleValuesForSyncMode() []string {
	return []string{
		string(SyncModeAutomatic),
		string(SyncModeManual),
	}
}

func (s *SyncMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSyncMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSyncMode(input string) (*SyncMode, error) {
	vals := map[string]SyncMode{
		"automatic": SyncModeAutomatic,
		"manual":    SyncModeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncMode(input)
	return &out, nil
}

type TransportProtocol string

const (
	TransportProtocolAll TransportProtocol = "All"
	TransportProtocolTcp TransportProtocol = "Tcp"
	TransportProtocolUdp TransportProtocol = "Udp"
)

func PossibleValuesForTransportProtocol() []string {
	return []string{
		string(TransportProtocolAll),
		string(TransportProtocolTcp),
		string(TransportProtocolUdp),
	}
}

func (s *TransportProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransportProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransportProtocol(input string) (*TransportProtocol, error) {
	vals := map[string]TransportProtocol{
		"all": TransportProtocolAll,
		"tcp": TransportProtocolTcp,
		"udp": TransportProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransportProtocol(input)
	return &out, nil
}

type VirtualNetworkPrivateEndpointNetworkPolicies string

const (
	VirtualNetworkPrivateEndpointNetworkPoliciesDisabled                    VirtualNetworkPrivateEndpointNetworkPolicies = "Disabled"
	VirtualNetworkPrivateEndpointNetworkPoliciesEnabled                     VirtualNetworkPrivateEndpointNetworkPolicies = "Enabled"
	VirtualNetworkPrivateEndpointNetworkPoliciesNetworkSecurityGroupEnabled VirtualNetworkPrivateEndpointNetworkPolicies = "NetworkSecurityGroupEnabled"
	VirtualNetworkPrivateEndpointNetworkPoliciesRouteTableEnabled           VirtualNetworkPrivateEndpointNetworkPolicies = "RouteTableEnabled"
)

func PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies() []string {
	return []string{
		string(VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
		string(VirtualNetworkPrivateEndpointNetworkPoliciesEnabled),
		string(VirtualNetworkPrivateEndpointNetworkPoliciesNetworkSecurityGroupEnabled),
		string(VirtualNetworkPrivateEndpointNetworkPoliciesRouteTableEnabled),
	}
}

func (s *VirtualNetworkPrivateEndpointNetworkPolicies) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkPrivateEndpointNetworkPolicies(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkPrivateEndpointNetworkPolicies(input string) (*VirtualNetworkPrivateEndpointNetworkPolicies, error) {
	vals := map[string]VirtualNetworkPrivateEndpointNetworkPolicies{
		"disabled":                    VirtualNetworkPrivateEndpointNetworkPoliciesDisabled,
		"enabled":                     VirtualNetworkPrivateEndpointNetworkPoliciesEnabled,
		"networksecuritygroupenabled": VirtualNetworkPrivateEndpointNetworkPoliciesNetworkSecurityGroupEnabled,
		"routetableenabled":           VirtualNetworkPrivateEndpointNetworkPoliciesRouteTableEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkPrivateEndpointNetworkPolicies(input)
	return &out, nil
}

type VirtualNetworkPrivateLinkServiceNetworkPolicies string

const (
	VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled VirtualNetworkPrivateLinkServiceNetworkPolicies = "Disabled"
	VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled  VirtualNetworkPrivateLinkServiceNetworkPolicies = "Enabled"
)

func PossibleValuesForVirtualNetworkPrivateLinkServiceNetworkPolicies() []string {
	return []string{
		string(VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled),
		string(VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled),
	}
}

func (s *VirtualNetworkPrivateLinkServiceNetworkPolicies) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkPrivateLinkServiceNetworkPolicies(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkPrivateLinkServiceNetworkPolicies(input string) (*VirtualNetworkPrivateLinkServiceNetworkPolicies, error) {
	vals := map[string]VirtualNetworkPrivateLinkServiceNetworkPolicies{
		"disabled": VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled,
		"enabled":  VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkPrivateLinkServiceNetworkPolicies(input)
	return &out, nil
}
