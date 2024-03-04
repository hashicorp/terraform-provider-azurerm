package cluster

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddOnFeatures string

const (
	AddOnFeaturesBackupRestoreService   AddOnFeatures = "BackupRestoreService"
	AddOnFeaturesDnsService             AddOnFeatures = "DnsService"
	AddOnFeaturesRepairManager          AddOnFeatures = "RepairManager"
	AddOnFeaturesResourceMonitorService AddOnFeatures = "ResourceMonitorService"
)

func PossibleValuesForAddOnFeatures() []string {
	return []string{
		string(AddOnFeaturesBackupRestoreService),
		string(AddOnFeaturesDnsService),
		string(AddOnFeaturesRepairManager),
		string(AddOnFeaturesResourceMonitorService),
	}
}

func (s *AddOnFeatures) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAddOnFeatures(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAddOnFeatures(input string) (*AddOnFeatures, error) {
	vals := map[string]AddOnFeatures{
		"backuprestoreservice":   AddOnFeaturesBackupRestoreService,
		"dnsservice":             AddOnFeaturesDnsService,
		"repairmanager":          AddOnFeaturesRepairManager,
		"resourcemonitorservice": AddOnFeaturesResourceMonitorService,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AddOnFeatures(input)
	return &out, nil
}

type ClusterState string

const (
	ClusterStateAutoScale                 ClusterState = "AutoScale"
	ClusterStateBaselineUpgrade           ClusterState = "BaselineUpgrade"
	ClusterStateDeploying                 ClusterState = "Deploying"
	ClusterStateEnforcingClusterVersion   ClusterState = "EnforcingClusterVersion"
	ClusterStateReady                     ClusterState = "Ready"
	ClusterStateUpdatingInfrastructure    ClusterState = "UpdatingInfrastructure"
	ClusterStateUpdatingUserCertificate   ClusterState = "UpdatingUserCertificate"
	ClusterStateUpdatingUserConfiguration ClusterState = "UpdatingUserConfiguration"
	ClusterStateUpgradeServiceUnreachable ClusterState = "UpgradeServiceUnreachable"
	ClusterStateWaitingForNodes           ClusterState = "WaitingForNodes"
)

func PossibleValuesForClusterState() []string {
	return []string{
		string(ClusterStateAutoScale),
		string(ClusterStateBaselineUpgrade),
		string(ClusterStateDeploying),
		string(ClusterStateEnforcingClusterVersion),
		string(ClusterStateReady),
		string(ClusterStateUpdatingInfrastructure),
		string(ClusterStateUpdatingUserCertificate),
		string(ClusterStateUpdatingUserConfiguration),
		string(ClusterStateUpgradeServiceUnreachable),
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
		"autoscale":                 ClusterStateAutoScale,
		"baselineupgrade":           ClusterStateBaselineUpgrade,
		"deploying":                 ClusterStateDeploying,
		"enforcingclusterversion":   ClusterStateEnforcingClusterVersion,
		"ready":                     ClusterStateReady,
		"updatinginfrastructure":    ClusterStateUpdatingInfrastructure,
		"updatingusercertificate":   ClusterStateUpdatingUserCertificate,
		"updatinguserconfiguration": ClusterStateUpdatingUserConfiguration,
		"upgradeserviceunreachable": ClusterStateUpgradeServiceUnreachable,
		"waitingfornodes":           ClusterStateWaitingForNodes,
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

type DurabilityLevel string

const (
	DurabilityLevelBronze DurabilityLevel = "Bronze"
	DurabilityLevelGold   DurabilityLevel = "Gold"
	DurabilityLevelSilver DurabilityLevel = "Silver"
)

func PossibleValuesForDurabilityLevel() []string {
	return []string{
		string(DurabilityLevelBronze),
		string(DurabilityLevelGold),
		string(DurabilityLevelSilver),
	}
}

func (s *DurabilityLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDurabilityLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDurabilityLevel(input string) (*DurabilityLevel, error) {
	vals := map[string]DurabilityLevel{
		"bronze": DurabilityLevelBronze,
		"gold":   DurabilityLevelGold,
		"silver": DurabilityLevelSilver,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DurabilityLevel(input)
	return &out, nil
}

type Environment string

const (
	EnvironmentLinux   Environment = "Linux"
	EnvironmentWindows Environment = "Windows"
)

func PossibleValuesForEnvironment() []string {
	return []string{
		string(EnvironmentLinux),
		string(EnvironmentWindows),
	}
}

func (s *Environment) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnvironment(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnvironment(input string) (*Environment, error) {
	vals := map[string]Environment{
		"linux":   EnvironmentLinux,
		"windows": EnvironmentWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Environment(input)
	return &out, nil
}

type NotificationCategory string

const (
	NotificationCategoryWaveProgress NotificationCategory = "WaveProgress"
)

func PossibleValuesForNotificationCategory() []string {
	return []string{
		string(NotificationCategoryWaveProgress),
	}
}

func (s *NotificationCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotificationCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotificationCategory(input string) (*NotificationCategory, error) {
	vals := map[string]NotificationCategory{
		"waveprogress": NotificationCategoryWaveProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotificationCategory(input)
	return &out, nil
}

type NotificationChannel string

const (
	NotificationChannelEmailSubscription NotificationChannel = "EmailSubscription"
	NotificationChannelEmailUser         NotificationChannel = "EmailUser"
)

func PossibleValuesForNotificationChannel() []string {
	return []string{
		string(NotificationChannelEmailSubscription),
		string(NotificationChannelEmailUser),
	}
}

func (s *NotificationChannel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotificationChannel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotificationChannel(input string) (*NotificationChannel, error) {
	vals := map[string]NotificationChannel{
		"emailsubscription": NotificationChannelEmailSubscription,
		"emailuser":         NotificationChannelEmailUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotificationChannel(input)
	return &out, nil
}

type NotificationLevel string

const (
	NotificationLevelAll      NotificationLevel = "All"
	NotificationLevelCritical NotificationLevel = "Critical"
)

func PossibleValuesForNotificationLevel() []string {
	return []string{
		string(NotificationLevelAll),
		string(NotificationLevelCritical),
	}
}

func (s *NotificationLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNotificationLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNotificationLevel(input string) (*NotificationLevel, error) {
	vals := map[string]NotificationLevel{
		"all":      NotificationLevelAll,
		"critical": NotificationLevelCritical,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NotificationLevel(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
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

type ReliabilityLevel string

const (
	ReliabilityLevelBronze   ReliabilityLevel = "Bronze"
	ReliabilityLevelGold     ReliabilityLevel = "Gold"
	ReliabilityLevelNone     ReliabilityLevel = "None"
	ReliabilityLevelPlatinum ReliabilityLevel = "Platinum"
	ReliabilityLevelSilver   ReliabilityLevel = "Silver"
)

func PossibleValuesForReliabilityLevel() []string {
	return []string{
		string(ReliabilityLevelBronze),
		string(ReliabilityLevelGold),
		string(ReliabilityLevelNone),
		string(ReliabilityLevelPlatinum),
		string(ReliabilityLevelSilver),
	}
}

func (s *ReliabilityLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReliabilityLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReliabilityLevel(input string) (*ReliabilityLevel, error) {
	vals := map[string]ReliabilityLevel{
		"bronze":   ReliabilityLevelBronze,
		"gold":     ReliabilityLevelGold,
		"none":     ReliabilityLevelNone,
		"platinum": ReliabilityLevelPlatinum,
		"silver":   ReliabilityLevelSilver,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReliabilityLevel(input)
	return &out, nil
}

type SfZonalUpgradeMode string

const (
	SfZonalUpgradeModeHierarchical SfZonalUpgradeMode = "Hierarchical"
	SfZonalUpgradeModeParallel     SfZonalUpgradeMode = "Parallel"
)

func PossibleValuesForSfZonalUpgradeMode() []string {
	return []string{
		string(SfZonalUpgradeModeHierarchical),
		string(SfZonalUpgradeModeParallel),
	}
}

func (s *SfZonalUpgradeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSfZonalUpgradeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSfZonalUpgradeMode(input string) (*SfZonalUpgradeMode, error) {
	vals := map[string]SfZonalUpgradeMode{
		"hierarchical": SfZonalUpgradeModeHierarchical,
		"parallel":     SfZonalUpgradeModeParallel,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SfZonalUpgradeMode(input)
	return &out, nil
}

type UpgradeMode string

const (
	UpgradeModeAutomatic UpgradeMode = "Automatic"
	UpgradeModeManual    UpgradeMode = "Manual"
)

func PossibleValuesForUpgradeMode() []string {
	return []string{
		string(UpgradeModeAutomatic),
		string(UpgradeModeManual),
	}
}

func (s *UpgradeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUpgradeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUpgradeMode(input string) (*UpgradeMode, error) {
	vals := map[string]UpgradeMode{
		"automatic": UpgradeModeAutomatic,
		"manual":    UpgradeModeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UpgradeMode(input)
	return &out, nil
}

type VMSSZonalUpgradeMode string

const (
	VMSSZonalUpgradeModeHierarchical VMSSZonalUpgradeMode = "Hierarchical"
	VMSSZonalUpgradeModeParallel     VMSSZonalUpgradeMode = "Parallel"
)

func PossibleValuesForVMSSZonalUpgradeMode() []string {
	return []string{
		string(VMSSZonalUpgradeModeHierarchical),
		string(VMSSZonalUpgradeModeParallel),
	}
}

func (s *VMSSZonalUpgradeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMSSZonalUpgradeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMSSZonalUpgradeMode(input string) (*VMSSZonalUpgradeMode, error) {
	vals := map[string]VMSSZonalUpgradeMode{
		"hierarchical": VMSSZonalUpgradeModeHierarchical,
		"parallel":     VMSSZonalUpgradeModeParallel,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMSSZonalUpgradeMode(input)
	return &out, nil
}

type X509StoreName string

const (
	X509StoreNameAddressBook          X509StoreName = "AddressBook"
	X509StoreNameAuthRoot             X509StoreName = "AuthRoot"
	X509StoreNameCertificateAuthority X509StoreName = "CertificateAuthority"
	X509StoreNameDisallowed           X509StoreName = "Disallowed"
	X509StoreNameMy                   X509StoreName = "My"
	X509StoreNameRoot                 X509StoreName = "Root"
	X509StoreNameTrustedPeople        X509StoreName = "TrustedPeople"
	X509StoreNameTrustedPublisher     X509StoreName = "TrustedPublisher"
)

func PossibleValuesForX509StoreName() []string {
	return []string{
		string(X509StoreNameAddressBook),
		string(X509StoreNameAuthRoot),
		string(X509StoreNameCertificateAuthority),
		string(X509StoreNameDisallowed),
		string(X509StoreNameMy),
		string(X509StoreNameRoot),
		string(X509StoreNameTrustedPeople),
		string(X509StoreNameTrustedPublisher),
	}
}

func (s *X509StoreName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseX509StoreName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseX509StoreName(input string) (*X509StoreName, error) {
	vals := map[string]X509StoreName{
		"addressbook":          X509StoreNameAddressBook,
		"authroot":             X509StoreNameAuthRoot,
		"certificateauthority": X509StoreNameCertificateAuthority,
		"disallowed":           X509StoreNameDisallowed,
		"my":                   X509StoreNameMy,
		"root":                 X509StoreNameRoot,
		"trustedpeople":        X509StoreNameTrustedPeople,
		"trustedpublisher":     X509StoreNameTrustedPublisher,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := X509StoreName(input)
	return &out, nil
}
