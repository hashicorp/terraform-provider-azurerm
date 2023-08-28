package virtualnetworkgateways

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdminState string

const (
	AdminStateDisabled AdminState = "Disabled"
	AdminStateEnabled  AdminState = "Enabled"
)

func PossibleValuesForAdminState() []string {
	return []string{
		string(AdminStateDisabled),
		string(AdminStateEnabled),
	}
}

func (s *AdminState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdminState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdminState(input string) (*AdminState, error) {
	vals := map[string]AdminState{
		"disabled": AdminStateDisabled,
		"enabled":  AdminStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdminState(input)
	return &out, nil
}

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

type BgpPeerState string

const (
	BgpPeerStateConnected  BgpPeerState = "Connected"
	BgpPeerStateConnecting BgpPeerState = "Connecting"
	BgpPeerStateIdle       BgpPeerState = "Idle"
	BgpPeerStateStopped    BgpPeerState = "Stopped"
	BgpPeerStateUnknown    BgpPeerState = "Unknown"
)

func PossibleValuesForBgpPeerState() []string {
	return []string{
		string(BgpPeerStateConnected),
		string(BgpPeerStateConnecting),
		string(BgpPeerStateIdle),
		string(BgpPeerStateStopped),
		string(BgpPeerStateUnknown),
	}
}

func (s *BgpPeerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBgpPeerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBgpPeerState(input string) (*BgpPeerState, error) {
	vals := map[string]BgpPeerState{
		"connected":  BgpPeerStateConnected,
		"connecting": BgpPeerStateConnecting,
		"idle":       BgpPeerStateIdle,
		"stopped":    BgpPeerStateStopped,
		"unknown":    BgpPeerStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BgpPeerState(input)
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

type ProcessorArchitecture string

const (
	ProcessorArchitectureAmdSixFour ProcessorArchitecture = "Amd64"
	ProcessorArchitectureXEightSix  ProcessorArchitecture = "X86"
)

func PossibleValuesForProcessorArchitecture() []string {
	return []string{
		string(ProcessorArchitectureAmdSixFour),
		string(ProcessorArchitectureXEightSix),
	}
}

func (s *ProcessorArchitecture) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProcessorArchitecture(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProcessorArchitecture(input string) (*ProcessorArchitecture, error) {
	vals := map[string]ProcessorArchitecture{
		"amd64": ProcessorArchitectureAmdSixFour,
		"x86":   ProcessorArchitectureXEightSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProcessorArchitecture(input)
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

type VirtualNetworkGatewayConnectionMode string

const (
	VirtualNetworkGatewayConnectionModeDefault       VirtualNetworkGatewayConnectionMode = "Default"
	VirtualNetworkGatewayConnectionModeInitiatorOnly VirtualNetworkGatewayConnectionMode = "InitiatorOnly"
	VirtualNetworkGatewayConnectionModeResponderOnly VirtualNetworkGatewayConnectionMode = "ResponderOnly"
)

func PossibleValuesForVirtualNetworkGatewayConnectionMode() []string {
	return []string{
		string(VirtualNetworkGatewayConnectionModeDefault),
		string(VirtualNetworkGatewayConnectionModeInitiatorOnly),
		string(VirtualNetworkGatewayConnectionModeResponderOnly),
	}
}

func (s *VirtualNetworkGatewayConnectionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewayConnectionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewayConnectionMode(input string) (*VirtualNetworkGatewayConnectionMode, error) {
	vals := map[string]VirtualNetworkGatewayConnectionMode{
		"default":       VirtualNetworkGatewayConnectionModeDefault,
		"initiatoronly": VirtualNetworkGatewayConnectionModeInitiatorOnly,
		"responderonly": VirtualNetworkGatewayConnectionModeResponderOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewayConnectionMode(input)
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

type VirtualNetworkGatewayConnectionStatus string

const (
	VirtualNetworkGatewayConnectionStatusConnected    VirtualNetworkGatewayConnectionStatus = "Connected"
	VirtualNetworkGatewayConnectionStatusConnecting   VirtualNetworkGatewayConnectionStatus = "Connecting"
	VirtualNetworkGatewayConnectionStatusNotConnected VirtualNetworkGatewayConnectionStatus = "NotConnected"
	VirtualNetworkGatewayConnectionStatusUnknown      VirtualNetworkGatewayConnectionStatus = "Unknown"
)

func PossibleValuesForVirtualNetworkGatewayConnectionStatus() []string {
	return []string{
		string(VirtualNetworkGatewayConnectionStatusConnected),
		string(VirtualNetworkGatewayConnectionStatusConnecting),
		string(VirtualNetworkGatewayConnectionStatusNotConnected),
		string(VirtualNetworkGatewayConnectionStatusUnknown),
	}
}

func (s *VirtualNetworkGatewayConnectionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewayConnectionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewayConnectionStatus(input string) (*VirtualNetworkGatewayConnectionStatus, error) {
	vals := map[string]VirtualNetworkGatewayConnectionStatus{
		"connected":    VirtualNetworkGatewayConnectionStatusConnected,
		"connecting":   VirtualNetworkGatewayConnectionStatusConnecting,
		"notconnected": VirtualNetworkGatewayConnectionStatusNotConnected,
		"unknown":      VirtualNetworkGatewayConnectionStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewayConnectionStatus(input)
	return &out, nil
}

type VirtualNetworkGatewayConnectionType string

const (
	VirtualNetworkGatewayConnectionTypeExpressRoute VirtualNetworkGatewayConnectionType = "ExpressRoute"
	VirtualNetworkGatewayConnectionTypeIPsec        VirtualNetworkGatewayConnectionType = "IPsec"
	VirtualNetworkGatewayConnectionTypeVPNClient    VirtualNetworkGatewayConnectionType = "VPNClient"
	VirtualNetworkGatewayConnectionTypeVnetTwoVnet  VirtualNetworkGatewayConnectionType = "Vnet2Vnet"
)

func PossibleValuesForVirtualNetworkGatewayConnectionType() []string {
	return []string{
		string(VirtualNetworkGatewayConnectionTypeExpressRoute),
		string(VirtualNetworkGatewayConnectionTypeIPsec),
		string(VirtualNetworkGatewayConnectionTypeVPNClient),
		string(VirtualNetworkGatewayConnectionTypeVnetTwoVnet),
	}
}

func (s *VirtualNetworkGatewayConnectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewayConnectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewayConnectionType(input string) (*VirtualNetworkGatewayConnectionType, error) {
	vals := map[string]VirtualNetworkGatewayConnectionType{
		"expressroute": VirtualNetworkGatewayConnectionTypeExpressRoute,
		"ipsec":        VirtualNetworkGatewayConnectionTypeIPsec,
		"vpnclient":    VirtualNetworkGatewayConnectionTypeVPNClient,
		"vnet2vnet":    VirtualNetworkGatewayConnectionTypeVnetTwoVnet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewayConnectionType(input)
	return &out, nil
}

type VirtualNetworkGatewaySkuName string

const (
	VirtualNetworkGatewaySkuNameBasic            VirtualNetworkGatewaySkuName = "Basic"
	VirtualNetworkGatewaySkuNameErGwOneAZ        VirtualNetworkGatewaySkuName = "ErGw1AZ"
	VirtualNetworkGatewaySkuNameErGwThreeAZ      VirtualNetworkGatewaySkuName = "ErGw3AZ"
	VirtualNetworkGatewaySkuNameErGwTwoAZ        VirtualNetworkGatewaySkuName = "ErGw2AZ"
	VirtualNetworkGatewaySkuNameHighPerformance  VirtualNetworkGatewaySkuName = "HighPerformance"
	VirtualNetworkGatewaySkuNameStandard         VirtualNetworkGatewaySkuName = "Standard"
	VirtualNetworkGatewaySkuNameUltraPerformance VirtualNetworkGatewaySkuName = "UltraPerformance"
	VirtualNetworkGatewaySkuNameVpnGwFive        VirtualNetworkGatewaySkuName = "VpnGw5"
	VirtualNetworkGatewaySkuNameVpnGwFiveAZ      VirtualNetworkGatewaySkuName = "VpnGw5AZ"
	VirtualNetworkGatewaySkuNameVpnGwFour        VirtualNetworkGatewaySkuName = "VpnGw4"
	VirtualNetworkGatewaySkuNameVpnGwFourAZ      VirtualNetworkGatewaySkuName = "VpnGw4AZ"
	VirtualNetworkGatewaySkuNameVpnGwOne         VirtualNetworkGatewaySkuName = "VpnGw1"
	VirtualNetworkGatewaySkuNameVpnGwOneAZ       VirtualNetworkGatewaySkuName = "VpnGw1AZ"
	VirtualNetworkGatewaySkuNameVpnGwThree       VirtualNetworkGatewaySkuName = "VpnGw3"
	VirtualNetworkGatewaySkuNameVpnGwThreeAZ     VirtualNetworkGatewaySkuName = "VpnGw3AZ"
	VirtualNetworkGatewaySkuNameVpnGwTwo         VirtualNetworkGatewaySkuName = "VpnGw2"
	VirtualNetworkGatewaySkuNameVpnGwTwoAZ       VirtualNetworkGatewaySkuName = "VpnGw2AZ"
)

func PossibleValuesForVirtualNetworkGatewaySkuName() []string {
	return []string{
		string(VirtualNetworkGatewaySkuNameBasic),
		string(VirtualNetworkGatewaySkuNameErGwOneAZ),
		string(VirtualNetworkGatewaySkuNameErGwThreeAZ),
		string(VirtualNetworkGatewaySkuNameErGwTwoAZ),
		string(VirtualNetworkGatewaySkuNameHighPerformance),
		string(VirtualNetworkGatewaySkuNameStandard),
		string(VirtualNetworkGatewaySkuNameUltraPerformance),
		string(VirtualNetworkGatewaySkuNameVpnGwFive),
		string(VirtualNetworkGatewaySkuNameVpnGwFiveAZ),
		string(VirtualNetworkGatewaySkuNameVpnGwFour),
		string(VirtualNetworkGatewaySkuNameVpnGwFourAZ),
		string(VirtualNetworkGatewaySkuNameVpnGwOne),
		string(VirtualNetworkGatewaySkuNameVpnGwOneAZ),
		string(VirtualNetworkGatewaySkuNameVpnGwThree),
		string(VirtualNetworkGatewaySkuNameVpnGwThreeAZ),
		string(VirtualNetworkGatewaySkuNameVpnGwTwo),
		string(VirtualNetworkGatewaySkuNameVpnGwTwoAZ),
	}
}

func (s *VirtualNetworkGatewaySkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewaySkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewaySkuName(input string) (*VirtualNetworkGatewaySkuName, error) {
	vals := map[string]VirtualNetworkGatewaySkuName{
		"basic":            VirtualNetworkGatewaySkuNameBasic,
		"ergw1az":          VirtualNetworkGatewaySkuNameErGwOneAZ,
		"ergw3az":          VirtualNetworkGatewaySkuNameErGwThreeAZ,
		"ergw2az":          VirtualNetworkGatewaySkuNameErGwTwoAZ,
		"highperformance":  VirtualNetworkGatewaySkuNameHighPerformance,
		"standard":         VirtualNetworkGatewaySkuNameStandard,
		"ultraperformance": VirtualNetworkGatewaySkuNameUltraPerformance,
		"vpngw5":           VirtualNetworkGatewaySkuNameVpnGwFive,
		"vpngw5az":         VirtualNetworkGatewaySkuNameVpnGwFiveAZ,
		"vpngw4":           VirtualNetworkGatewaySkuNameVpnGwFour,
		"vpngw4az":         VirtualNetworkGatewaySkuNameVpnGwFourAZ,
		"vpngw1":           VirtualNetworkGatewaySkuNameVpnGwOne,
		"vpngw1az":         VirtualNetworkGatewaySkuNameVpnGwOneAZ,
		"vpngw3":           VirtualNetworkGatewaySkuNameVpnGwThree,
		"vpngw3az":         VirtualNetworkGatewaySkuNameVpnGwThreeAZ,
		"vpngw2":           VirtualNetworkGatewaySkuNameVpnGwTwo,
		"vpngw2az":         VirtualNetworkGatewaySkuNameVpnGwTwoAZ,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewaySkuName(input)
	return &out, nil
}

type VirtualNetworkGatewaySkuTier string

const (
	VirtualNetworkGatewaySkuTierBasic            VirtualNetworkGatewaySkuTier = "Basic"
	VirtualNetworkGatewaySkuTierErGwOneAZ        VirtualNetworkGatewaySkuTier = "ErGw1AZ"
	VirtualNetworkGatewaySkuTierErGwThreeAZ      VirtualNetworkGatewaySkuTier = "ErGw3AZ"
	VirtualNetworkGatewaySkuTierErGwTwoAZ        VirtualNetworkGatewaySkuTier = "ErGw2AZ"
	VirtualNetworkGatewaySkuTierHighPerformance  VirtualNetworkGatewaySkuTier = "HighPerformance"
	VirtualNetworkGatewaySkuTierStandard         VirtualNetworkGatewaySkuTier = "Standard"
	VirtualNetworkGatewaySkuTierUltraPerformance VirtualNetworkGatewaySkuTier = "UltraPerformance"
	VirtualNetworkGatewaySkuTierVpnGwFive        VirtualNetworkGatewaySkuTier = "VpnGw5"
	VirtualNetworkGatewaySkuTierVpnGwFiveAZ      VirtualNetworkGatewaySkuTier = "VpnGw5AZ"
	VirtualNetworkGatewaySkuTierVpnGwFour        VirtualNetworkGatewaySkuTier = "VpnGw4"
	VirtualNetworkGatewaySkuTierVpnGwFourAZ      VirtualNetworkGatewaySkuTier = "VpnGw4AZ"
	VirtualNetworkGatewaySkuTierVpnGwOne         VirtualNetworkGatewaySkuTier = "VpnGw1"
	VirtualNetworkGatewaySkuTierVpnGwOneAZ       VirtualNetworkGatewaySkuTier = "VpnGw1AZ"
	VirtualNetworkGatewaySkuTierVpnGwThree       VirtualNetworkGatewaySkuTier = "VpnGw3"
	VirtualNetworkGatewaySkuTierVpnGwThreeAZ     VirtualNetworkGatewaySkuTier = "VpnGw3AZ"
	VirtualNetworkGatewaySkuTierVpnGwTwo         VirtualNetworkGatewaySkuTier = "VpnGw2"
	VirtualNetworkGatewaySkuTierVpnGwTwoAZ       VirtualNetworkGatewaySkuTier = "VpnGw2AZ"
)

func PossibleValuesForVirtualNetworkGatewaySkuTier() []string {
	return []string{
		string(VirtualNetworkGatewaySkuTierBasic),
		string(VirtualNetworkGatewaySkuTierErGwOneAZ),
		string(VirtualNetworkGatewaySkuTierErGwThreeAZ),
		string(VirtualNetworkGatewaySkuTierErGwTwoAZ),
		string(VirtualNetworkGatewaySkuTierHighPerformance),
		string(VirtualNetworkGatewaySkuTierStandard),
		string(VirtualNetworkGatewaySkuTierUltraPerformance),
		string(VirtualNetworkGatewaySkuTierVpnGwFive),
		string(VirtualNetworkGatewaySkuTierVpnGwFiveAZ),
		string(VirtualNetworkGatewaySkuTierVpnGwFour),
		string(VirtualNetworkGatewaySkuTierVpnGwFourAZ),
		string(VirtualNetworkGatewaySkuTierVpnGwOne),
		string(VirtualNetworkGatewaySkuTierVpnGwOneAZ),
		string(VirtualNetworkGatewaySkuTierVpnGwThree),
		string(VirtualNetworkGatewaySkuTierVpnGwThreeAZ),
		string(VirtualNetworkGatewaySkuTierVpnGwTwo),
		string(VirtualNetworkGatewaySkuTierVpnGwTwoAZ),
	}
}

func (s *VirtualNetworkGatewaySkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewaySkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewaySkuTier(input string) (*VirtualNetworkGatewaySkuTier, error) {
	vals := map[string]VirtualNetworkGatewaySkuTier{
		"basic":            VirtualNetworkGatewaySkuTierBasic,
		"ergw1az":          VirtualNetworkGatewaySkuTierErGwOneAZ,
		"ergw3az":          VirtualNetworkGatewaySkuTierErGwThreeAZ,
		"ergw2az":          VirtualNetworkGatewaySkuTierErGwTwoAZ,
		"highperformance":  VirtualNetworkGatewaySkuTierHighPerformance,
		"standard":         VirtualNetworkGatewaySkuTierStandard,
		"ultraperformance": VirtualNetworkGatewaySkuTierUltraPerformance,
		"vpngw5":           VirtualNetworkGatewaySkuTierVpnGwFive,
		"vpngw5az":         VirtualNetworkGatewaySkuTierVpnGwFiveAZ,
		"vpngw4":           VirtualNetworkGatewaySkuTierVpnGwFour,
		"vpngw4az":         VirtualNetworkGatewaySkuTierVpnGwFourAZ,
		"vpngw1":           VirtualNetworkGatewaySkuTierVpnGwOne,
		"vpngw1az":         VirtualNetworkGatewaySkuTierVpnGwOneAZ,
		"vpngw3":           VirtualNetworkGatewaySkuTierVpnGwThree,
		"vpngw3az":         VirtualNetworkGatewaySkuTierVpnGwThreeAZ,
		"vpngw2":           VirtualNetworkGatewaySkuTierVpnGwTwo,
		"vpngw2az":         VirtualNetworkGatewaySkuTierVpnGwTwoAZ,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewaySkuTier(input)
	return &out, nil
}

type VirtualNetworkGatewayType string

const (
	VirtualNetworkGatewayTypeExpressRoute VirtualNetworkGatewayType = "ExpressRoute"
	VirtualNetworkGatewayTypeLocalGateway VirtualNetworkGatewayType = "LocalGateway"
	VirtualNetworkGatewayTypeVpn          VirtualNetworkGatewayType = "Vpn"
)

func PossibleValuesForVirtualNetworkGatewayType() []string {
	return []string{
		string(VirtualNetworkGatewayTypeExpressRoute),
		string(VirtualNetworkGatewayTypeLocalGateway),
		string(VirtualNetworkGatewayTypeVpn),
	}
}

func (s *VirtualNetworkGatewayType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkGatewayType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkGatewayType(input string) (*VirtualNetworkGatewayType, error) {
	vals := map[string]VirtualNetworkGatewayType{
		"expressroute": VirtualNetworkGatewayTypeExpressRoute,
		"localgateway": VirtualNetworkGatewayTypeLocalGateway,
		"vpn":          VirtualNetworkGatewayTypeVpn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkGatewayType(input)
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

type VpnClientProtocol string

const (
	VpnClientProtocolIkeVTwo VpnClientProtocol = "IkeV2"
	VpnClientProtocolOpenVPN VpnClientProtocol = "OpenVPN"
	VpnClientProtocolSSTP    VpnClientProtocol = "SSTP"
)

func PossibleValuesForVpnClientProtocol() []string {
	return []string{
		string(VpnClientProtocolIkeVTwo),
		string(VpnClientProtocolOpenVPN),
		string(VpnClientProtocolSSTP),
	}
}

func (s *VpnClientProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnClientProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnClientProtocol(input string) (*VpnClientProtocol, error) {
	vals := map[string]VpnClientProtocol{
		"ikev2":   VpnClientProtocolIkeVTwo,
		"openvpn": VpnClientProtocolOpenVPN,
		"sstp":    VpnClientProtocolSSTP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnClientProtocol(input)
	return &out, nil
}

type VpnGatewayGeneration string

const (
	VpnGatewayGenerationGenerationOne VpnGatewayGeneration = "Generation1"
	VpnGatewayGenerationGenerationTwo VpnGatewayGeneration = "Generation2"
	VpnGatewayGenerationNone          VpnGatewayGeneration = "None"
)

func PossibleValuesForVpnGatewayGeneration() []string {
	return []string{
		string(VpnGatewayGenerationGenerationOne),
		string(VpnGatewayGenerationGenerationTwo),
		string(VpnGatewayGenerationNone),
	}
}

func (s *VpnGatewayGeneration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnGatewayGeneration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnGatewayGeneration(input string) (*VpnGatewayGeneration, error) {
	vals := map[string]VpnGatewayGeneration{
		"generation1": VpnGatewayGenerationGenerationOne,
		"generation2": VpnGatewayGenerationGenerationTwo,
		"none":        VpnGatewayGenerationNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnGatewayGeneration(input)
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

type VpnType string

const (
	VpnTypePolicyBased VpnType = "PolicyBased"
	VpnTypeRouteBased  VpnType = "RouteBased"
)

func PossibleValuesForVpnType() []string {
	return []string{
		string(VpnTypePolicyBased),
		string(VpnTypeRouteBased),
	}
}

func (s *VpnType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVpnType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVpnType(input string) (*VpnType, error) {
	vals := map[string]VpnType{
		"policybased": VpnTypePolicyBased,
		"routebased":  VpnTypeRouteBased,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VpnType(input)
	return &out, nil
}
