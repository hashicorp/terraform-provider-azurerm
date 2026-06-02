package nodetype

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Access string

const (
	AccessAllow Access = "allow"
	AccessDeny  Access = "deny"
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

type Direction string

const (
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
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

type DiskType string

const (
	DiskTypePremiumLRS     DiskType = "Premium_LRS"
	DiskTypeStandardLRS    DiskType = "Standard_LRS"
	DiskTypeStandardSSDLRS DiskType = "StandardSSD_LRS"
)

func PossibleValuesForDiskType() []string {
	return []string{
		string(DiskTypePremiumLRS),
		string(DiskTypeStandardLRS),
		string(DiskTypeStandardSSDLRS),
	}
}

func (s *DiskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskType(input string) (*DiskType, error) {
	vals := map[string]DiskType{
		"premium_lrs":     DiskTypePremiumLRS,
		"standard_lrs":    DiskTypeStandardLRS,
		"standardssd_lrs": DiskTypeStandardSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskType(input)
	return &out, nil
}

type EvictionPolicyType string

const (
	EvictionPolicyTypeDeallocate EvictionPolicyType = "Deallocate"
	EvictionPolicyTypeDelete     EvictionPolicyType = "Delete"
)

func PossibleValuesForEvictionPolicyType() []string {
	return []string{
		string(EvictionPolicyTypeDeallocate),
		string(EvictionPolicyTypeDelete),
	}
}

func (s *EvictionPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEvictionPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEvictionPolicyType(input string) (*EvictionPolicyType, error) {
	vals := map[string]EvictionPolicyType{
		"deallocate": EvictionPolicyTypeDeallocate,
		"delete":     EvictionPolicyTypeDelete,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EvictionPolicyType(input)
	return &out, nil
}

type IPAddressType string

const (
	IPAddressTypeIPvFour IPAddressType = "IPv4"
	IPAddressTypeIPvSix  IPAddressType = "IPv6"
)

func PossibleValuesForIPAddressType() []string {
	return []string{
		string(IPAddressTypeIPvFour),
		string(IPAddressTypeIPvSix),
	}
}

func (s *IPAddressType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPAddressType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPAddressType(input string) (*IPAddressType, error) {
	vals := map[string]IPAddressType{
		"ipv4": IPAddressTypeIPvFour,
		"ipv6": IPAddressTypeIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPAddressType(input)
	return &out, nil
}

type ManagedResourceProvisioningState string

const (
	ManagedResourceProvisioningStateCanceled  ManagedResourceProvisioningState = "Canceled"
	ManagedResourceProvisioningStateCreated   ManagedResourceProvisioningState = "Created"
	ManagedResourceProvisioningStateCreating  ManagedResourceProvisioningState = "Creating"
	ManagedResourceProvisioningStateDeleted   ManagedResourceProvisioningState = "Deleted"
	ManagedResourceProvisioningStateDeleting  ManagedResourceProvisioningState = "Deleting"
	ManagedResourceProvisioningStateFailed    ManagedResourceProvisioningState = "Failed"
	ManagedResourceProvisioningStateNone      ManagedResourceProvisioningState = "None"
	ManagedResourceProvisioningStateOther     ManagedResourceProvisioningState = "Other"
	ManagedResourceProvisioningStateSucceeded ManagedResourceProvisioningState = "Succeeded"
	ManagedResourceProvisioningStateUpdating  ManagedResourceProvisioningState = "Updating"
)

func PossibleValuesForManagedResourceProvisioningState() []string {
	return []string{
		string(ManagedResourceProvisioningStateCanceled),
		string(ManagedResourceProvisioningStateCreated),
		string(ManagedResourceProvisioningStateCreating),
		string(ManagedResourceProvisioningStateDeleted),
		string(ManagedResourceProvisioningStateDeleting),
		string(ManagedResourceProvisioningStateFailed),
		string(ManagedResourceProvisioningStateNone),
		string(ManagedResourceProvisioningStateOther),
		string(ManagedResourceProvisioningStateSucceeded),
		string(ManagedResourceProvisioningStateUpdating),
	}
}

func (s *ManagedResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedResourceProvisioningState(input string) (*ManagedResourceProvisioningState, error) {
	vals := map[string]ManagedResourceProvisioningState{
		"canceled":  ManagedResourceProvisioningStateCanceled,
		"created":   ManagedResourceProvisioningStateCreated,
		"creating":  ManagedResourceProvisioningStateCreating,
		"deleted":   ManagedResourceProvisioningStateDeleted,
		"deleting":  ManagedResourceProvisioningStateDeleting,
		"failed":    ManagedResourceProvisioningStateFailed,
		"none":      ManagedResourceProvisioningStateNone,
		"other":     ManagedResourceProvisioningStateOther,
		"succeeded": ManagedResourceProvisioningStateSucceeded,
		"updating":  ManagedResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedResourceProvisioningState(input)
	return &out, nil
}

type NodeTypeSkuScaleType string

const (
	NodeTypeSkuScaleTypeAutomatic NodeTypeSkuScaleType = "Automatic"
	NodeTypeSkuScaleTypeManual    NodeTypeSkuScaleType = "Manual"
	NodeTypeSkuScaleTypeNone      NodeTypeSkuScaleType = "None"
)

func PossibleValuesForNodeTypeSkuScaleType() []string {
	return []string{
		string(NodeTypeSkuScaleTypeAutomatic),
		string(NodeTypeSkuScaleTypeManual),
		string(NodeTypeSkuScaleTypeNone),
	}
}

func (s *NodeTypeSkuScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeTypeSkuScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeTypeSkuScaleType(input string) (*NodeTypeSkuScaleType, error) {
	vals := map[string]NodeTypeSkuScaleType{
		"automatic": NodeTypeSkuScaleTypeAutomatic,
		"manual":    NodeTypeSkuScaleTypeManual,
		"none":      NodeTypeSkuScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeTypeSkuScaleType(input)
	return &out, nil
}

type NsgProtocol string

const (
	NsgProtocolAh    NsgProtocol = "ah"
	NsgProtocolEsp   NsgProtocol = "esp"
	NsgProtocolHTTP  NsgProtocol = "http"
	NsgProtocolHTTPS NsgProtocol = "https"
	NsgProtocolIcmp  NsgProtocol = "icmp"
	NsgProtocolTcp   NsgProtocol = "tcp"
	NsgProtocolUdp   NsgProtocol = "udp"
)

func PossibleValuesForNsgProtocol() []string {
	return []string{
		string(NsgProtocolAh),
		string(NsgProtocolEsp),
		string(NsgProtocolHTTP),
		string(NsgProtocolHTTPS),
		string(NsgProtocolIcmp),
		string(NsgProtocolTcp),
		string(NsgProtocolUdp),
	}
}

func (s *NsgProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNsgProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNsgProtocol(input string) (*NsgProtocol, error) {
	vals := map[string]NsgProtocol{
		"ah":    NsgProtocolAh,
		"esp":   NsgProtocolEsp,
		"http":  NsgProtocolHTTP,
		"https": NsgProtocolHTTPS,
		"icmp":  NsgProtocolIcmp,
		"tcp":   NsgProtocolTcp,
		"udp":   NsgProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NsgProtocol(input)
	return &out, nil
}

type PrivateIPAddressVersion string

const (
	PrivateIPAddressVersionIPvFour PrivateIPAddressVersion = "IPv4"
	PrivateIPAddressVersionIPvSix  PrivateIPAddressVersion = "IPv6"
)

func PossibleValuesForPrivateIPAddressVersion() []string {
	return []string{
		string(PrivateIPAddressVersionIPvFour),
		string(PrivateIPAddressVersionIPvSix),
	}
}

func (s *PrivateIPAddressVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateIPAddressVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateIPAddressVersion(input string) (*PrivateIPAddressVersion, error) {
	vals := map[string]PrivateIPAddressVersion{
		"ipv4": PrivateIPAddressVersionIPvFour,
		"ipv6": PrivateIPAddressVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateIPAddressVersion(input)
	return &out, nil
}

type PublicIPAddressVersion string

const (
	PublicIPAddressVersionIPvFour PublicIPAddressVersion = "IPv4"
	PublicIPAddressVersionIPvSix  PublicIPAddressVersion = "IPv6"
)

func PossibleValuesForPublicIPAddressVersion() []string {
	return []string{
		string(PublicIPAddressVersionIPvFour),
		string(PublicIPAddressVersionIPvSix),
	}
}

func (s *PublicIPAddressVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressVersion(input string) (*PublicIPAddressVersion, error) {
	vals := map[string]PublicIPAddressVersion{
		"ipv4": PublicIPAddressVersionIPvFour,
		"ipv6": PublicIPAddressVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressVersion(input)
	return &out, nil
}

type SecurityType string

const (
	SecurityTypeStandard      SecurityType = "Standard"
	SecurityTypeTrustedLaunch SecurityType = "TrustedLaunch"
)

func PossibleValuesForSecurityType() []string {
	return []string{
		string(SecurityTypeStandard),
		string(SecurityTypeTrustedLaunch),
	}
}

func (s *SecurityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityType(input string) (*SecurityType, error) {
	vals := map[string]SecurityType{
		"standard":      SecurityTypeStandard,
		"trustedlaunch": SecurityTypeTrustedLaunch,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityType(input)
	return &out, nil
}

type UpdateType string

const (
	UpdateTypeByUpgradeDomain UpdateType = "ByUpgradeDomain"
	UpdateTypeDefault         UpdateType = "Default"
)

func PossibleValuesForUpdateType() []string {
	return []string{
		string(UpdateTypeByUpgradeDomain),
		string(UpdateTypeDefault),
	}
}

func (s *UpdateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpdateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpdateType(input string) (*UpdateType, error) {
	vals := map[string]UpdateType{
		"byupgradedomain": UpdateTypeByUpgradeDomain,
		"default":         UpdateTypeDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpdateType(input)
	return &out, nil
}

type VMSSExtensionSetupOrder string

const (
	VMSSExtensionSetupOrderBeforeSFRuntime VMSSExtensionSetupOrder = "BeforeSFRuntime"
)

func PossibleValuesForVMSSExtensionSetupOrder() []string {
	return []string{
		string(VMSSExtensionSetupOrderBeforeSFRuntime),
	}
}

func (s *VMSSExtensionSetupOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMSSExtensionSetupOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMSSExtensionSetupOrder(input string) (*VMSSExtensionSetupOrder, error) {
	vals := map[string]VMSSExtensionSetupOrder{
		"beforesfruntime": VMSSExtensionSetupOrderBeforeSFRuntime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMSSExtensionSetupOrder(input)
	return &out, nil
}

type VMSetupAction string

const (
	VMSetupActionEnableContainers VMSetupAction = "EnableContainers"
	VMSetupActionEnableHyperV     VMSetupAction = "EnableHyperV"
)

func PossibleValuesForVMSetupAction() []string {
	return []string{
		string(VMSetupActionEnableContainers),
		string(VMSetupActionEnableHyperV),
	}
}

func (s *VMSetupAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMSetupAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMSetupAction(input string) (*VMSetupAction, error) {
	vals := map[string]VMSetupAction{
		"enablecontainers": VMSetupActionEnableContainers,
		"enablehyperv":     VMSetupActionEnableHyperV,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMSetupAction(input)
	return &out, nil
}
