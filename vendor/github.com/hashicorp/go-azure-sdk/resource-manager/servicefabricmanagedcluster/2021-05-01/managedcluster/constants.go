package managedcluster

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

type AddonFeatures string

const (
	AddonFeaturesBackupRestoreService   AddonFeatures = "BackupRestoreService"
	AddonFeaturesDnsService             AddonFeatures = "DnsService"
	AddonFeaturesResourceMonitorService AddonFeatures = "ResourceMonitorService"
)

func PossibleValuesForAddonFeatures() []string {
	return []string{
		string(AddonFeaturesBackupRestoreService),
		string(AddonFeaturesDnsService),
		string(AddonFeaturesResourceMonitorService),
	}
}

func (s *AddonFeatures) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAddonFeatures(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAddonFeatures(input string) (*AddonFeatures, error) {
	vals := map[string]AddonFeatures{
		"backuprestoreservice":   AddonFeaturesBackupRestoreService,
		"dnsservice":             AddonFeaturesDnsService,
		"resourcemonitorservice": AddonFeaturesResourceMonitorService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AddonFeatures(input)
	return &out, nil
}

type ClusterState string

const (
	ClusterStateBaselineUpgrade ClusterState = "BaselineUpgrade"
	ClusterStateDeploying       ClusterState = "Deploying"
	ClusterStateReady           ClusterState = "Ready"
	ClusterStateUpgradeFailed   ClusterState = "UpgradeFailed"
	ClusterStateUpgrading       ClusterState = "Upgrading"
	ClusterStateWaitingForNodes ClusterState = "WaitingForNodes"
)

func PossibleValuesForClusterState() []string {
	return []string{
		string(ClusterStateBaselineUpgrade),
		string(ClusterStateDeploying),
		string(ClusterStateReady),
		string(ClusterStateUpgradeFailed),
		string(ClusterStateUpgrading),
		string(ClusterStateWaitingForNodes),
	}
}

func (s *ClusterState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterState(input string) (*ClusterState, error) {
	vals := map[string]ClusterState{
		"baselineupgrade": ClusterStateBaselineUpgrade,
		"deploying":       ClusterStateDeploying,
		"ready":           ClusterStateReady,
		"upgradefailed":   ClusterStateUpgradeFailed,
		"upgrading":       ClusterStateUpgrading,
		"waitingfornodes": ClusterStateWaitingForNodes,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterState(input)
	return &out, nil
}

type ClusterUpgradeCadence string

const (
	ClusterUpgradeCadenceWaveOne  ClusterUpgradeCadence = "Wave1"
	ClusterUpgradeCadenceWaveTwo  ClusterUpgradeCadence = "Wave2"
	ClusterUpgradeCadenceWaveZero ClusterUpgradeCadence = "Wave0"
)

func PossibleValuesForClusterUpgradeCadence() []string {
	return []string{
		string(ClusterUpgradeCadenceWaveOne),
		string(ClusterUpgradeCadenceWaveTwo),
		string(ClusterUpgradeCadenceWaveZero),
	}
}

func (s *ClusterUpgradeCadence) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterUpgradeCadence(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterUpgradeCadence(input string) (*ClusterUpgradeCadence, error) {
	vals := map[string]ClusterUpgradeCadence{
		"wave1": ClusterUpgradeCadenceWaveOne,
		"wave2": ClusterUpgradeCadenceWaveTwo,
		"wave0": ClusterUpgradeCadenceWaveZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterUpgradeCadence(input)
	return &out, nil
}

type ClusterUpgradeMode string

const (
	ClusterUpgradeModeAutomatic ClusterUpgradeMode = "Automatic"
	ClusterUpgradeModeManual    ClusterUpgradeMode = "Manual"
)

func PossibleValuesForClusterUpgradeMode() []string {
	return []string{
		string(ClusterUpgradeModeAutomatic),
		string(ClusterUpgradeModeManual),
	}
}

func (s *ClusterUpgradeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterUpgradeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterUpgradeMode(input string) (*ClusterUpgradeMode, error) {
	vals := map[string]ClusterUpgradeMode{
		"automatic": ClusterUpgradeModeAutomatic,
		"manual":    ClusterUpgradeModeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterUpgradeMode(input)
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

type ProbeProtocol string

const (
	ProbeProtocolHTTP  ProbeProtocol = "http"
	ProbeProtocolHTTPS ProbeProtocol = "https"
	ProbeProtocolTcp   ProbeProtocol = "tcp"
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

type Protocol string

const (
	ProtocolTcp Protocol = "tcp"
	ProtocolUdp Protocol = "udp"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolTcp),
		string(ProtocolUdp),
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
		"tcp": ProtocolTcp,
		"udp": ProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Protocol(input)
	return &out, nil
}

type SkuName string

const (
	SkuNameBasic    SkuName = "Basic"
	SkuNameStandard SkuName = "Standard"
)

func PossibleValuesForSkuName() []string {
	return []string{
		string(SkuNameBasic),
		string(SkuNameStandard),
	}
}

func (s *SkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuName(input string) (*SkuName, error) {
	vals := map[string]SkuName{
		"basic":    SkuNameBasic,
		"standard": SkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuName(input)
	return &out, nil
}
