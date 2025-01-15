package databases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessKeyType string

const (
	AccessKeyTypePrimary   AccessKeyType = "Primary"
	AccessKeyTypeSecondary AccessKeyType = "Secondary"
)

func PossibleValuesForAccessKeyType() []string {
	return []string{
		string(AccessKeyTypePrimary),
		string(AccessKeyTypeSecondary),
	}
}

func (s *AccessKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessKeyType(input string) (*AccessKeyType, error) {
	vals := map[string]AccessKeyType{
		"primary":   AccessKeyTypePrimary,
		"secondary": AccessKeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessKeyType(input)
	return &out, nil
}

type AofFrequency string

const (
	AofFrequencyAlways AofFrequency = "always"
	AofFrequencyOnes   AofFrequency = "1s"
)

func PossibleValuesForAofFrequency() []string {
	return []string{
		string(AofFrequencyAlways),
		string(AofFrequencyOnes),
	}
}

func (s *AofFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAofFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAofFrequency(input string) (*AofFrequency, error) {
	vals := map[string]AofFrequency{
		"always": AofFrequencyAlways,
		"1s":     AofFrequencyOnes,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AofFrequency(input)
	return &out, nil
}

type ClusteringPolicy string

const (
	ClusteringPolicyEnterpriseCluster ClusteringPolicy = "EnterpriseCluster"
	ClusteringPolicyOSSCluster        ClusteringPolicy = "OSSCluster"
)

func PossibleValuesForClusteringPolicy() []string {
	return []string{
		string(ClusteringPolicyEnterpriseCluster),
		string(ClusteringPolicyOSSCluster),
	}
}

func (s *ClusteringPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusteringPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusteringPolicy(input string) (*ClusteringPolicy, error) {
	vals := map[string]ClusteringPolicy{
		"enterprisecluster": ClusteringPolicyEnterpriseCluster,
		"osscluster":        ClusteringPolicyOSSCluster,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusteringPolicy(input)
	return &out, nil
}

type DeferUpgradeSetting string

const (
	DeferUpgradeSettingDeferred    DeferUpgradeSetting = "Deferred"
	DeferUpgradeSettingNotDeferred DeferUpgradeSetting = "NotDeferred"
)

func PossibleValuesForDeferUpgradeSetting() []string {
	return []string{
		string(DeferUpgradeSettingDeferred),
		string(DeferUpgradeSettingNotDeferred),
	}
}

func (s *DeferUpgradeSetting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeferUpgradeSetting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeferUpgradeSetting(input string) (*DeferUpgradeSetting, error) {
	vals := map[string]DeferUpgradeSetting{
		"deferred":    DeferUpgradeSettingDeferred,
		"notdeferred": DeferUpgradeSettingNotDeferred,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeferUpgradeSetting(input)
	return &out, nil
}

type EvictionPolicy string

const (
	EvictionPolicyAllKeysLFU     EvictionPolicy = "AllKeysLFU"
	EvictionPolicyAllKeysLRU     EvictionPolicy = "AllKeysLRU"
	EvictionPolicyAllKeysRandom  EvictionPolicy = "AllKeysRandom"
	EvictionPolicyNoEviction     EvictionPolicy = "NoEviction"
	EvictionPolicyVolatileLFU    EvictionPolicy = "VolatileLFU"
	EvictionPolicyVolatileLRU    EvictionPolicy = "VolatileLRU"
	EvictionPolicyVolatileRandom EvictionPolicy = "VolatileRandom"
	EvictionPolicyVolatileTTL    EvictionPolicy = "VolatileTTL"
)

func PossibleValuesForEvictionPolicy() []string {
	return []string{
		string(EvictionPolicyAllKeysLFU),
		string(EvictionPolicyAllKeysLRU),
		string(EvictionPolicyAllKeysRandom),
		string(EvictionPolicyNoEviction),
		string(EvictionPolicyVolatileLFU),
		string(EvictionPolicyVolatileLRU),
		string(EvictionPolicyVolatileRandom),
		string(EvictionPolicyVolatileTTL),
	}
}

func (s *EvictionPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEvictionPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEvictionPolicy(input string) (*EvictionPolicy, error) {
	vals := map[string]EvictionPolicy{
		"allkeyslfu":     EvictionPolicyAllKeysLFU,
		"allkeyslru":     EvictionPolicyAllKeysLRU,
		"allkeysrandom":  EvictionPolicyAllKeysRandom,
		"noeviction":     EvictionPolicyNoEviction,
		"volatilelfu":    EvictionPolicyVolatileLFU,
		"volatilelru":    EvictionPolicyVolatileLRU,
		"volatilerandom": EvictionPolicyVolatileRandom,
		"volatilettl":    EvictionPolicyVolatileTTL,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EvictionPolicy(input)
	return &out, nil
}

type LinkState string

const (
	LinkStateLinkFailed   LinkState = "LinkFailed"
	LinkStateLinked       LinkState = "Linked"
	LinkStateLinking      LinkState = "Linking"
	LinkStateUnlinkFailed LinkState = "UnlinkFailed"
	LinkStateUnlinking    LinkState = "Unlinking"
)

func PossibleValuesForLinkState() []string {
	return []string{
		string(LinkStateLinkFailed),
		string(LinkStateLinked),
		string(LinkStateLinking),
		string(LinkStateUnlinkFailed),
		string(LinkStateUnlinking),
	}
}

func (s *LinkState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinkState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinkState(input string) (*LinkState, error) {
	vals := map[string]LinkState{
		"linkfailed":   LinkStateLinkFailed,
		"linked":       LinkStateLinked,
		"linking":      LinkStateLinking,
		"unlinkfailed": LinkStateUnlinkFailed,
		"unlinking":    LinkStateUnlinking,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinkState(input)
	return &out, nil
}

type Protocol string

const (
	ProtocolEncrypted Protocol = "Encrypted"
	ProtocolPlaintext Protocol = "Plaintext"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolEncrypted),
		string(ProtocolPlaintext),
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
		"encrypted": ProtocolEncrypted,
		"plaintext": ProtocolPlaintext,
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

type RdbFrequency string

const (
	RdbFrequencyOneTwoh RdbFrequency = "12h"
	RdbFrequencyOneh    RdbFrequency = "1h"
	RdbFrequencySixh    RdbFrequency = "6h"
)

func PossibleValuesForRdbFrequency() []string {
	return []string{
		string(RdbFrequencyOneTwoh),
		string(RdbFrequencyOneh),
		string(RdbFrequencySixh),
	}
}

func (s *RdbFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRdbFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRdbFrequency(input string) (*RdbFrequency, error) {
	vals := map[string]RdbFrequency{
		"12h": RdbFrequencyOneTwoh,
		"1h":  RdbFrequencyOneh,
		"6h":  RdbFrequencySixh,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RdbFrequency(input)
	return &out, nil
}

type ResourceState string

const (
	ResourceStateCreateFailed  ResourceState = "CreateFailed"
	ResourceStateCreating      ResourceState = "Creating"
	ResourceStateDeleteFailed  ResourceState = "DeleteFailed"
	ResourceStateDeleting      ResourceState = "Deleting"
	ResourceStateDisableFailed ResourceState = "DisableFailed"
	ResourceStateDisabled      ResourceState = "Disabled"
	ResourceStateDisabling     ResourceState = "Disabling"
	ResourceStateEnableFailed  ResourceState = "EnableFailed"
	ResourceStateEnabling      ResourceState = "Enabling"
	ResourceStateRunning       ResourceState = "Running"
	ResourceStateScaling       ResourceState = "Scaling"
	ResourceStateScalingFailed ResourceState = "ScalingFailed"
	ResourceStateUpdateFailed  ResourceState = "UpdateFailed"
	ResourceStateUpdating      ResourceState = "Updating"
)

func PossibleValuesForResourceState() []string {
	return []string{
		string(ResourceStateCreateFailed),
		string(ResourceStateCreating),
		string(ResourceStateDeleteFailed),
		string(ResourceStateDeleting),
		string(ResourceStateDisableFailed),
		string(ResourceStateDisabled),
		string(ResourceStateDisabling),
		string(ResourceStateEnableFailed),
		string(ResourceStateEnabling),
		string(ResourceStateRunning),
		string(ResourceStateScaling),
		string(ResourceStateScalingFailed),
		string(ResourceStateUpdateFailed),
		string(ResourceStateUpdating),
	}
}

func (s *ResourceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceState(input string) (*ResourceState, error) {
	vals := map[string]ResourceState{
		"createfailed":  ResourceStateCreateFailed,
		"creating":      ResourceStateCreating,
		"deletefailed":  ResourceStateDeleteFailed,
		"deleting":      ResourceStateDeleting,
		"disablefailed": ResourceStateDisableFailed,
		"disabled":      ResourceStateDisabled,
		"disabling":     ResourceStateDisabling,
		"enablefailed":  ResourceStateEnableFailed,
		"enabling":      ResourceStateEnabling,
		"running":       ResourceStateRunning,
		"scaling":       ResourceStateScaling,
		"scalingfailed": ResourceStateScalingFailed,
		"updatefailed":  ResourceStateUpdateFailed,
		"updating":      ResourceStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceState(input)
	return &out, nil
}
