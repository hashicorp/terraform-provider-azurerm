package virtualwans

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationMethod string

const (
	AuthenticationMethodEAPMSCHAPvTwo AuthenticationMethod = "EAPMSCHAPv2"
	AuthenticationMethodEAPTLS        AuthenticationMethod = "EAPTLS"
)

func PossibleValuesForAuthenticationMethod() []string {
	return []string{
		string(AuthenticationMethodEAPMSCHAPvTwo),
		string(AuthenticationMethodEAPTLS),
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
		"eapmschapv2": AuthenticationMethodEAPMSCHAPvTwo,
		"eaptls":      AuthenticationMethodEAPTLS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationMethod(input)
	return &out, nil
}

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

type DhGroup string

const (
	DhGroupDHGroupOne              DhGroup = "DHGroup1"
	DhGroupDHGroupOneFour          DhGroup = "DHGroup14"
	DhGroupDHGroupTwo              DhGroup = "DHGroup2"
	DhGroupDHGroupTwoFour          DhGroup = "DHGroup24"
	DhGroupDHGroupTwoZeroFourEight DhGroup = "DHGroup2048"
	DhGroupECPThreeEightFour       DhGroup = "ECP384"
	DhGroupECPTwoFiveSix           DhGroup = "ECP256"
	DhGroupNone                    DhGroup = "None"
)

func PossibleValuesForDhGroup() []string {
	return []string{
		string(DhGroupDHGroupOne),
		string(DhGroupDHGroupOneFour),
		string(DhGroupDHGroupTwo),
		string(DhGroupDHGroupTwoFour),
		string(DhGroupDHGroupTwoZeroFourEight),
		string(DhGroupECPThreeEightFour),
		string(DhGroupECPTwoFiveSix),
		string(DhGroupNone),
	}
}

func (s *DhGroup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDhGroup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDhGroup(input string) (*DhGroup, error) {
	vals := map[string]DhGroup{
		"dhgroup1":    DhGroupDHGroupOne,
		"dhgroup14":   DhGroupDHGroupOneFour,
		"dhgroup2":    DhGroupDHGroupTwo,
		"dhgroup24":   DhGroupDHGroupTwoFour,
		"dhgroup2048": DhGroupDHGroupTwoZeroFourEight,
		"ecp384":      DhGroupECPThreeEightFour,
		"ecp256":      DhGroupECPTwoFiveSix,
		"none":        DhGroupNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DhGroup(input)
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

type HubBgpConnectionStatus string

const (
	HubBgpConnectionStatusConnected    HubBgpConnectionStatus = "Connected"
	HubBgpConnectionStatusConnecting   HubBgpConnectionStatus = "Connecting"
	HubBgpConnectionStatusNotConnected HubBgpConnectionStatus = "NotConnected"
	HubBgpConnectionStatusUnknown      HubBgpConnectionStatus = "Unknown"
)

func PossibleValuesForHubBgpConnectionStatus() []string {
	return []string{
		string(HubBgpConnectionStatusConnected),
		string(HubBgpConnectionStatusConnecting),
		string(HubBgpConnectionStatusNotConnected),
		string(HubBgpConnectionStatusUnknown),
	}
}

func (s *HubBgpConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHubBgpConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHubBgpConnectionStatus(input string) (*HubBgpConnectionStatus, error) {
	vals := map[string]HubBgpConnectionStatus{
		"connected":    HubBgpConnectionStatusConnected,
		"connecting":   HubBgpConnectionStatusConnecting,
		"notconnected": HubBgpConnectionStatusNotConnected,
		"unknown":      HubBgpConnectionStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HubBgpConnectionStatus(input)
	return &out, nil
}

type HubRoutingPreference string

const (
	HubRoutingPreferenceASPath       HubRoutingPreference = "ASPath"
	HubRoutingPreferenceExpressRoute HubRoutingPreference = "ExpressRoute"
	HubRoutingPreferenceVpnGateway   HubRoutingPreference = "VpnGateway"
)

func PossibleValuesForHubRoutingPreference() []string {
	return []string{
		string(HubRoutingPreferenceASPath),
		string(HubRoutingPreferenceExpressRoute),
		string(HubRoutingPreferenceVpnGateway),
	}
}

func (s *HubRoutingPreference) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHubRoutingPreference(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHubRoutingPreference(input string) (*HubRoutingPreference, error) {
	vals := map[string]HubRoutingPreference{
		"aspath":       HubRoutingPreferenceASPath,
		"expressroute": HubRoutingPreferenceExpressRoute,
		"vpngateway":   HubRoutingPreferenceVpnGateway,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HubRoutingPreference(input)
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

type IPsecEncryption string

const (
	IPsecEncryptionAESOneNineTwo     IPsecEncryption = "AES192"
	IPsecEncryptionAESOneTwoEight    IPsecEncryption = "AES128"
	IPsecEncryptionAESTwoFiveSix     IPsecEncryption = "AES256"
	IPsecEncryptionDES               IPsecEncryption = "DES"
	IPsecEncryptionDESThree          IPsecEncryption = "DES3"
	IPsecEncryptionGCMAESOneNineTwo  IPsecEncryption = "GCMAES192"
	IPsecEncryptionGCMAESOneTwoEight IPsecEncryption = "GCMAES128"
	IPsecEncryptionGCMAESTwoFiveSix  IPsecEncryption = "GCMAES256"
	IPsecEncryptionNone              IPsecEncryption = "None"
)

func PossibleValuesForIPsecEncryption() []string {
	return []string{
		string(IPsecEncryptionAESOneNineTwo),
		string(IPsecEncryptionAESOneTwoEight),
		string(IPsecEncryptionAESTwoFiveSix),
		string(IPsecEncryptionDES),
		string(IPsecEncryptionDESThree),
		string(IPsecEncryptionGCMAESOneNineTwo),
		string(IPsecEncryptionGCMAESOneTwoEight),
		string(IPsecEncryptionGCMAESTwoFiveSix),
		string(IPsecEncryptionNone),
	}
}

func (s *IPsecEncryption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPsecEncryption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPsecEncryption(input string) (*IPsecEncryption, error) {
	vals := map[string]IPsecEncryption{
		"aes192":    IPsecEncryptionAESOneNineTwo,
		"aes128":    IPsecEncryptionAESOneTwoEight,
		"aes256":    IPsecEncryptionAESTwoFiveSix,
		"des":       IPsecEncryptionDES,
		"des3":      IPsecEncryptionDESThree,
		"gcmaes192": IPsecEncryptionGCMAESOneNineTwo,
		"gcmaes128": IPsecEncryptionGCMAESOneTwoEight,
		"gcmaes256": IPsecEncryptionGCMAESTwoFiveSix,
		"none":      IPsecEncryptionNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPsecEncryption(input)
	return &out, nil
}

type IPsecIntegrity string

const (
	IPsecIntegrityGCMAESOneNineTwo  IPsecIntegrity = "GCMAES192"
	IPsecIntegrityGCMAESOneTwoEight IPsecIntegrity = "GCMAES128"
	IPsecIntegrityGCMAESTwoFiveSix  IPsecIntegrity = "GCMAES256"
	IPsecIntegrityMDFive            IPsecIntegrity = "MD5"
	IPsecIntegritySHAOne            IPsecIntegrity = "SHA1"
	IPsecIntegritySHATwoFiveSix     IPsecIntegrity = "SHA256"
)

func PossibleValuesForIPsecIntegrity() []string {
	return []string{
		string(IPsecIntegrityGCMAESOneNineTwo),
		string(IPsecIntegrityGCMAESOneTwoEight),
		string(IPsecIntegrityGCMAESTwoFiveSix),
		string(IPsecIntegrityMDFive),
		string(IPsecIntegritySHAOne),
		string(IPsecIntegritySHATwoFiveSix),
	}
}

func (s *IPsecIntegrity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPsecIntegrity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPsecIntegrity(input string) (*IPsecIntegrity, error) {
	vals := map[string]IPsecIntegrity{
		"gcmaes192": IPsecIntegrityGCMAESOneNineTwo,
		"gcmaes128": IPsecIntegrityGCMAESOneTwoEight,
		"gcmaes256": IPsecIntegrityGCMAESTwoFiveSix,
		"md5":       IPsecIntegrityMDFive,
		"sha1":      IPsecIntegritySHAOne,
		"sha256":    IPsecIntegritySHATwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPsecIntegrity(input)
	return &out, nil
}

type IkeEncryption string

const (
	IkeEncryptionAESOneNineTwo     IkeEncryption = "AES192"
	IkeEncryptionAESOneTwoEight    IkeEncryption = "AES128"
	IkeEncryptionAESTwoFiveSix     IkeEncryption = "AES256"
	IkeEncryptionDES               IkeEncryption = "DES"
	IkeEncryptionDESThree          IkeEncryption = "DES3"
	IkeEncryptionGCMAESOneTwoEight IkeEncryption = "GCMAES128"
	IkeEncryptionGCMAESTwoFiveSix  IkeEncryption = "GCMAES256"
)

func PossibleValuesForIkeEncryption() []string {
	return []string{
		string(IkeEncryptionAESOneNineTwo),
		string(IkeEncryptionAESOneTwoEight),
		string(IkeEncryptionAESTwoFiveSix),
		string(IkeEncryptionDES),
		string(IkeEncryptionDESThree),
		string(IkeEncryptionGCMAESOneTwoEight),
		string(IkeEncryptionGCMAESTwoFiveSix),
	}
}

func (s *IkeEncryption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIkeEncryption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIkeEncryption(input string) (*IkeEncryption, error) {
	vals := map[string]IkeEncryption{
		"aes192":    IkeEncryptionAESOneNineTwo,
		"aes128":    IkeEncryptionAESOneTwoEight,
		"aes256":    IkeEncryptionAESTwoFiveSix,
		"des":       IkeEncryptionDES,
		"des3":      IkeEncryptionDESThree,
		"gcmaes128": IkeEncryptionGCMAESOneTwoEight,
		"gcmaes256": IkeEncryptionGCMAESTwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IkeEncryption(input)
	return &out, nil
}

type IkeIntegrity string

const (
	IkeIntegrityGCMAESOneTwoEight IkeIntegrity = "GCMAES128"
	IkeIntegrityGCMAESTwoFiveSix  IkeIntegrity = "GCMAES256"
	IkeIntegrityMDFive            IkeIntegrity = "MD5"
	IkeIntegritySHAOne            IkeIntegrity = "SHA1"
	IkeIntegritySHAThreeEightFour IkeIntegrity = "SHA384"
	IkeIntegritySHATwoFiveSix     IkeIntegrity = "SHA256"
)

func PossibleValuesForIkeIntegrity() []string {
	return []string{
		string(IkeIntegrityGCMAESOneTwoEight),
		string(IkeIntegrityGCMAESTwoFiveSix),
		string(IkeIntegrityMDFive),
		string(IkeIntegritySHAOne),
		string(IkeIntegritySHAThreeEightFour),
		string(IkeIntegritySHATwoFiveSix),
	}
}

func (s *IkeIntegrity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIkeIntegrity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIkeIntegrity(input string) (*IkeIntegrity, error) {
	vals := map[string]IkeIntegrity{
		"gcmaes128": IkeIntegrityGCMAESOneTwoEight,
		"gcmaes256": IkeIntegrityGCMAESTwoFiveSix,
		"md5":       IkeIntegrityMDFive,
		"sha1":      IkeIntegritySHAOne,
		"sha384":    IkeIntegritySHAThreeEightFour,
		"sha256":    IkeIntegritySHATwoFiveSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IkeIntegrity(input)
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

type NextStep string

const (
	NextStepContinue  NextStep = "Continue"
	NextStepTerminate NextStep = "Terminate"
	NextStepUnknown   NextStep = "Unknown"
)

func PossibleValuesForNextStep() []string {
	return []string{
		string(NextStepContinue),
		string(NextStepTerminate),
		string(NextStepUnknown),
	}
}

func (s *NextStep) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNextStep(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNextStep(input string) (*NextStep, error) {
	vals := map[string]NextStep{
		"continue":  NextStepContinue,
		"terminate": NextStepTerminate,
		"unknown":   NextStepUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NextStep(input)
	return &out, nil
}

type OfficeTrafficCategory string

const (
	OfficeTrafficCategoryAll              OfficeTrafficCategory = "All"
	OfficeTrafficCategoryNone             OfficeTrafficCategory = "None"
	OfficeTrafficCategoryOptimize         OfficeTrafficCategory = "Optimize"
	OfficeTrafficCategoryOptimizeAndAllow OfficeTrafficCategory = "OptimizeAndAllow"
)

func PossibleValuesForOfficeTrafficCategory() []string {
	return []string{
		string(OfficeTrafficCategoryAll),
		string(OfficeTrafficCategoryNone),
		string(OfficeTrafficCategoryOptimize),
		string(OfficeTrafficCategoryOptimizeAndAllow),
	}
}

func (s *OfficeTrafficCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOfficeTrafficCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOfficeTrafficCategory(input string) (*OfficeTrafficCategory, error) {
	vals := map[string]OfficeTrafficCategory{
		"all":              OfficeTrafficCategoryAll,
		"none":             OfficeTrafficCategoryNone,
		"optimize":         OfficeTrafficCategoryOptimize,
		"optimizeandallow": OfficeTrafficCategoryOptimizeAndAllow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OfficeTrafficCategory(input)
	return &out, nil
}

type PfsGroup string

const (
	PfsGroupECPThreeEightFour   PfsGroup = "ECP384"
	PfsGroupECPTwoFiveSix       PfsGroup = "ECP256"
	PfsGroupNone                PfsGroup = "None"
	PfsGroupPFSMM               PfsGroup = "PFSMM"
	PfsGroupPFSOne              PfsGroup = "PFS1"
	PfsGroupPFSOneFour          PfsGroup = "PFS14"
	PfsGroupPFSTwo              PfsGroup = "PFS2"
	PfsGroupPFSTwoFour          PfsGroup = "PFS24"
	PfsGroupPFSTwoZeroFourEight PfsGroup = "PFS2048"
)

func PossibleValuesForPfsGroup() []string {
	return []string{
		string(PfsGroupECPThreeEightFour),
		string(PfsGroupECPTwoFiveSix),
		string(PfsGroupNone),
		string(PfsGroupPFSMM),
		string(PfsGroupPFSOne),
		string(PfsGroupPFSOneFour),
		string(PfsGroupPFSTwo),
		string(PfsGroupPFSTwoFour),
		string(PfsGroupPFSTwoZeroFourEight),
	}
}

func (s *PfsGroup) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePfsGroup(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePfsGroup(input string) (*PfsGroup, error) {
	vals := map[string]PfsGroup{
		"ecp384":  PfsGroupECPThreeEightFour,
		"ecp256":  PfsGroupECPTwoFiveSix,
		"none":    PfsGroupNone,
		"pfsmm":   PfsGroupPFSMM,
		"pfs1":    PfsGroupPFSOne,
		"pfs14":   PfsGroupPFSOneFour,
		"pfs2":    PfsGroupPFSTwo,
		"pfs24":   PfsGroupPFSTwoFour,
		"pfs2048": PfsGroupPFSTwoZeroFourEight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PfsGroup(input)
	return &out, nil
}

type PreferredRoutingGateway string

const (
	PreferredRoutingGatewayExpressRoute PreferredRoutingGateway = "ExpressRoute"
	PreferredRoutingGatewayNone         PreferredRoutingGateway = "None"
	PreferredRoutingGatewayVpnGateway   PreferredRoutingGateway = "VpnGateway"
)

func PossibleValuesForPreferredRoutingGateway() []string {
	return []string{
		string(PreferredRoutingGatewayExpressRoute),
		string(PreferredRoutingGatewayNone),
		string(PreferredRoutingGatewayVpnGateway),
	}
}

func (s *PreferredRoutingGateway) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreferredRoutingGateway(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreferredRoutingGateway(input string) (*PreferredRoutingGateway, error) {
	vals := map[string]PreferredRoutingGateway{
		"expressroute": PreferredRoutingGatewayExpressRoute,
		"none":         PreferredRoutingGatewayNone,
		"vpngateway":   PreferredRoutingGatewayVpnGateway,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreferredRoutingGateway(input)
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

type RouteMapActionType string

const (
	RouteMapActionTypeAdd     RouteMapActionType = "Add"
	RouteMapActionTypeDrop    RouteMapActionType = "Drop"
	RouteMapActionTypeRemove  RouteMapActionType = "Remove"
	RouteMapActionTypeReplace RouteMapActionType = "Replace"
	RouteMapActionTypeUnknown RouteMapActionType = "Unknown"
)

func PossibleValuesForRouteMapActionType() []string {
	return []string{
		string(RouteMapActionTypeAdd),
		string(RouteMapActionTypeDrop),
		string(RouteMapActionTypeRemove),
		string(RouteMapActionTypeReplace),
		string(RouteMapActionTypeUnknown),
	}
}

func (s *RouteMapActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRouteMapActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRouteMapActionType(input string) (*RouteMapActionType, error) {
	vals := map[string]RouteMapActionType{
		"add":     RouteMapActionTypeAdd,
		"drop":    RouteMapActionTypeDrop,
		"remove":  RouteMapActionTypeRemove,
		"replace": RouteMapActionTypeReplace,
		"unknown": RouteMapActionTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RouteMapActionType(input)
	return &out, nil
}

type RouteMapMatchCondition string

const (
	RouteMapMatchConditionContains    RouteMapMatchCondition = "Contains"
	RouteMapMatchConditionEquals      RouteMapMatchCondition = "Equals"
	RouteMapMatchConditionNotContains RouteMapMatchCondition = "NotContains"
	RouteMapMatchConditionNotEquals   RouteMapMatchCondition = "NotEquals"
	RouteMapMatchConditionUnknown     RouteMapMatchCondition = "Unknown"
)

func PossibleValuesForRouteMapMatchCondition() []string {
	return []string{
		string(RouteMapMatchConditionContains),
		string(RouteMapMatchConditionEquals),
		string(RouteMapMatchConditionNotContains),
		string(RouteMapMatchConditionNotEquals),
		string(RouteMapMatchConditionUnknown),
	}
}

func (s *RouteMapMatchCondition) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRouteMapMatchCondition(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRouteMapMatchCondition(input string) (*RouteMapMatchCondition, error) {
	vals := map[string]RouteMapMatchCondition{
		"contains":    RouteMapMatchConditionContains,
		"equals":      RouteMapMatchConditionEquals,
		"notcontains": RouteMapMatchConditionNotContains,
		"notequals":   RouteMapMatchConditionNotEquals,
		"unknown":     RouteMapMatchConditionUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RouteMapMatchCondition(input)
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

type RoutingState string

const (
	RoutingStateFailed       RoutingState = "Failed"
	RoutingStateNone         RoutingState = "None"
	RoutingStateProvisioned  RoutingState = "Provisioned"
	RoutingStateProvisioning RoutingState = "Provisioning"
)

func PossibleValuesForRoutingState() []string {
	return []string{
		string(RoutingStateFailed),
		string(RoutingStateNone),
		string(RoutingStateProvisioned),
		string(RoutingStateProvisioning),
	}
}

func (s *RoutingState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRoutingState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRoutingState(input string) (*RoutingState, error) {
	vals := map[string]RoutingState{
		"failed":       RoutingStateFailed,
		"none":         RoutingStateNone,
		"provisioned":  RoutingStateProvisioned,
		"provisioning": RoutingStateProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RoutingState(input)
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

type VirtualNetworkGatewayConnectionProtocol string

const (
	VirtualNetworkGatewayConnectionProtocolIKEvOne VirtualNetworkGatewayConnectionProtocol = "IKEv1"
	VirtualNetworkGatewayConnectionProtocolIKEvTwo VirtualNetworkGatewayConnectionProtocol = "IKEv2"
)

func PossibleValuesForVirtualNetworkGatewayConnectionProtocol() []string {
	return []string{
		string(VirtualNetworkGatewayConnectionProtocolIKEvOne),
		string(VirtualNetworkGatewayConnectionProtocolIKEvTwo),
	}
}

func (s *VirtualNetworkGatewayConnectionProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewayConnectionProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewayConnectionProtocol(input string) (*VirtualNetworkGatewayConnectionProtocol, error) {
	vals := map[string]VirtualNetworkGatewayConnectionProtocol{
		"ikev1": VirtualNetworkGatewayConnectionProtocolIKEvOne,
		"ikev2": VirtualNetworkGatewayConnectionProtocolIKEvTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewayConnectionProtocol(input)
	return &out, nil
}

type VirtualNetworkPrivateEndpointNetworkPolicies string

const (
	VirtualNetworkPrivateEndpointNetworkPoliciesDisabled VirtualNetworkPrivateEndpointNetworkPolicies = "Disabled"
	VirtualNetworkPrivateEndpointNetworkPoliciesEnabled  VirtualNetworkPrivateEndpointNetworkPolicies = "Enabled"
)

func PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies() []string {
	return []string{
		string(VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
		string(VirtualNetworkPrivateEndpointNetworkPoliciesEnabled),
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
		"disabled": VirtualNetworkPrivateEndpointNetworkPoliciesDisabled,
		"enabled":  VirtualNetworkPrivateEndpointNetworkPoliciesEnabled,
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

type VirtualWanSecurityProviderType string

const (
	VirtualWanSecurityProviderTypeExternal VirtualWanSecurityProviderType = "External"
	VirtualWanSecurityProviderTypeNative   VirtualWanSecurityProviderType = "Native"
)

func PossibleValuesForVirtualWanSecurityProviderType() []string {
	return []string{
		string(VirtualWanSecurityProviderTypeExternal),
		string(VirtualWanSecurityProviderTypeNative),
	}
}

func (s *VirtualWanSecurityProviderType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualWanSecurityProviderType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualWanSecurityProviderType(input string) (*VirtualWanSecurityProviderType, error) {
	vals := map[string]VirtualWanSecurityProviderType{
		"external": VirtualWanSecurityProviderTypeExternal,
		"native":   VirtualWanSecurityProviderTypeNative,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualWanSecurityProviderType(input)
	return &out, nil
}

type VnetLocalRouteOverrideCriteria string

const (
	VnetLocalRouteOverrideCriteriaContains VnetLocalRouteOverrideCriteria = "Contains"
	VnetLocalRouteOverrideCriteriaEqual    VnetLocalRouteOverrideCriteria = "Equal"
)

func PossibleValuesForVnetLocalRouteOverrideCriteria() []string {
	return []string{
		string(VnetLocalRouteOverrideCriteriaContains),
		string(VnetLocalRouteOverrideCriteriaEqual),
	}
}

func (s *VnetLocalRouteOverrideCriteria) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVnetLocalRouteOverrideCriteria(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVnetLocalRouteOverrideCriteria(input string) (*VnetLocalRouteOverrideCriteria, error) {
	vals := map[string]VnetLocalRouteOverrideCriteria{
		"contains": VnetLocalRouteOverrideCriteriaContains,
		"equal":    VnetLocalRouteOverrideCriteriaEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VnetLocalRouteOverrideCriteria(input)
	return &out, nil
}

type VpnAuthenticationType string

const (
	VpnAuthenticationTypeAAD         VpnAuthenticationType = "AAD"
	VpnAuthenticationTypeCertificate VpnAuthenticationType = "Certificate"
	VpnAuthenticationTypeRadius      VpnAuthenticationType = "Radius"
)

func PossibleValuesForVpnAuthenticationType() []string {
	return []string{
		string(VpnAuthenticationTypeAAD),
		string(VpnAuthenticationTypeCertificate),
		string(VpnAuthenticationTypeRadius),
	}
}

func (s *VpnAuthenticationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnAuthenticationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnAuthenticationType(input string) (*VpnAuthenticationType, error) {
	vals := map[string]VpnAuthenticationType{
		"aad":         VpnAuthenticationTypeAAD,
		"certificate": VpnAuthenticationTypeCertificate,
		"radius":      VpnAuthenticationTypeRadius,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnAuthenticationType(input)
	return &out, nil
}

type VpnConnectionStatus string

const (
	VpnConnectionStatusConnected    VpnConnectionStatus = "Connected"
	VpnConnectionStatusConnecting   VpnConnectionStatus = "Connecting"
	VpnConnectionStatusNotConnected VpnConnectionStatus = "NotConnected"
	VpnConnectionStatusUnknown      VpnConnectionStatus = "Unknown"
)

func PossibleValuesForVpnConnectionStatus() []string {
	return []string{
		string(VpnConnectionStatusConnected),
		string(VpnConnectionStatusConnecting),
		string(VpnConnectionStatusNotConnected),
		string(VpnConnectionStatusUnknown),
	}
}

func (s *VpnConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnConnectionStatus(input string) (*VpnConnectionStatus, error) {
	vals := map[string]VpnConnectionStatus{
		"connected":    VpnConnectionStatusConnected,
		"connecting":   VpnConnectionStatusConnecting,
		"notconnected": VpnConnectionStatusNotConnected,
		"unknown":      VpnConnectionStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnConnectionStatus(input)
	return &out, nil
}

type VpnGatewayTunnelingProtocol string

const (
	VpnGatewayTunnelingProtocolIkeVTwo VpnGatewayTunnelingProtocol = "IkeV2"
	VpnGatewayTunnelingProtocolOpenVPN VpnGatewayTunnelingProtocol = "OpenVPN"
)

func PossibleValuesForVpnGatewayTunnelingProtocol() []string {
	return []string{
		string(VpnGatewayTunnelingProtocolIkeVTwo),
		string(VpnGatewayTunnelingProtocolOpenVPN),
	}
}

func (s *VpnGatewayTunnelingProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnGatewayTunnelingProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnGatewayTunnelingProtocol(input string) (*VpnGatewayTunnelingProtocol, error) {
	vals := map[string]VpnGatewayTunnelingProtocol{
		"ikev2":   VpnGatewayTunnelingProtocolIkeVTwo,
		"openvpn": VpnGatewayTunnelingProtocolOpenVPN,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnGatewayTunnelingProtocol(input)
	return &out, nil
}

type VpnLinkConnectionMode string

const (
	VpnLinkConnectionModeDefault       VpnLinkConnectionMode = "Default"
	VpnLinkConnectionModeInitiatorOnly VpnLinkConnectionMode = "InitiatorOnly"
	VpnLinkConnectionModeResponderOnly VpnLinkConnectionMode = "ResponderOnly"
)

func PossibleValuesForVpnLinkConnectionMode() []string {
	return []string{
		string(VpnLinkConnectionModeDefault),
		string(VpnLinkConnectionModeInitiatorOnly),
		string(VpnLinkConnectionModeResponderOnly),
	}
}

func (s *VpnLinkConnectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnLinkConnectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnLinkConnectionMode(input string) (*VpnLinkConnectionMode, error) {
	vals := map[string]VpnLinkConnectionMode{
		"default":       VpnLinkConnectionModeDefault,
		"initiatoronly": VpnLinkConnectionModeInitiatorOnly,
		"responderonly": VpnLinkConnectionModeResponderOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnLinkConnectionMode(input)
	return &out, nil
}

type VpnNatRuleMode string

const (
	VpnNatRuleModeEgressSnat  VpnNatRuleMode = "EgressSnat"
	VpnNatRuleModeIngressSnat VpnNatRuleMode = "IngressSnat"
)

func PossibleValuesForVpnNatRuleMode() []string {
	return []string{
		string(VpnNatRuleModeEgressSnat),
		string(VpnNatRuleModeIngressSnat),
	}
}

func (s *VpnNatRuleMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnNatRuleMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnNatRuleMode(input string) (*VpnNatRuleMode, error) {
	vals := map[string]VpnNatRuleMode{
		"egresssnat":  VpnNatRuleModeEgressSnat,
		"ingresssnat": VpnNatRuleModeIngressSnat,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnNatRuleMode(input)
	return &out, nil
}

type VpnNatRuleType string

const (
	VpnNatRuleTypeDynamic VpnNatRuleType = "Dynamic"
	VpnNatRuleTypeStatic  VpnNatRuleType = "Static"
)

func PossibleValuesForVpnNatRuleType() []string {
	return []string{
		string(VpnNatRuleTypeDynamic),
		string(VpnNatRuleTypeStatic),
	}
}

func (s *VpnNatRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnNatRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnNatRuleType(input string) (*VpnNatRuleType, error) {
	vals := map[string]VpnNatRuleType{
		"dynamic": VpnNatRuleTypeDynamic,
		"static":  VpnNatRuleTypeStatic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnNatRuleType(input)
	return &out, nil
}

type VpnPolicyMemberAttributeType string

const (
	VpnPolicyMemberAttributeTypeAADGroupId         VpnPolicyMemberAttributeType = "AADGroupId"
	VpnPolicyMemberAttributeTypeCertificateGroupId VpnPolicyMemberAttributeType = "CertificateGroupId"
	VpnPolicyMemberAttributeTypeRadiusAzureGroupId VpnPolicyMemberAttributeType = "RadiusAzureGroupId"
)

func PossibleValuesForVpnPolicyMemberAttributeType() []string {
	return []string{
		string(VpnPolicyMemberAttributeTypeAADGroupId),
		string(VpnPolicyMemberAttributeTypeCertificateGroupId),
		string(VpnPolicyMemberAttributeTypeRadiusAzureGroupId),
	}
}

func (s *VpnPolicyMemberAttributeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnPolicyMemberAttributeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnPolicyMemberAttributeType(input string) (*VpnPolicyMemberAttributeType, error) {
	vals := map[string]VpnPolicyMemberAttributeType{
		"aadgroupid":         VpnPolicyMemberAttributeTypeAADGroupId,
		"certificategroupid": VpnPolicyMemberAttributeTypeCertificateGroupId,
		"radiusazuregroupid": VpnPolicyMemberAttributeTypeRadiusAzureGroupId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnPolicyMemberAttributeType(input)
	return &out, nil
}
