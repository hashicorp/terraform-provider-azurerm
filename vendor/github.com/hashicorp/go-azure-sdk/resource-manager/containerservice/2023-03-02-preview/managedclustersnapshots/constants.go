package managedclustersnapshots

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerSku string

const (
	LoadBalancerSkuBasic    LoadBalancerSku = "basic"
	LoadBalancerSkuStandard LoadBalancerSku = "standard"
)

func PossibleValuesForLoadBalancerSku() []string {
	return []string{
		string(LoadBalancerSkuBasic),
		string(LoadBalancerSkuStandard),
	}
}

func parseLoadBalancerSku(input string) (*LoadBalancerSku, error) {
	vals := map[string]LoadBalancerSku{
		"basic":    LoadBalancerSkuBasic,
		"standard": LoadBalancerSkuStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerSku(input)
	return &out, nil
}

type ManagedClusterSKUName string

const (
	ManagedClusterSKUNameBase ManagedClusterSKUName = "Base"
)

func PossibleValuesForManagedClusterSKUName() []string {
	return []string{
		string(ManagedClusterSKUNameBase),
	}
}

func parseManagedClusterSKUName(input string) (*ManagedClusterSKUName, error) {
	vals := map[string]ManagedClusterSKUName{
		"base": ManagedClusterSKUNameBase,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedClusterSKUName(input)
	return &out, nil
}

type ManagedClusterSKUTier string

const (
	ManagedClusterSKUTierFree     ManagedClusterSKUTier = "Free"
	ManagedClusterSKUTierPremium  ManagedClusterSKUTier = "Premium"
	ManagedClusterSKUTierStandard ManagedClusterSKUTier = "Standard"
)

func PossibleValuesForManagedClusterSKUTier() []string {
	return []string{
		string(ManagedClusterSKUTierFree),
		string(ManagedClusterSKUTierPremium),
		string(ManagedClusterSKUTierStandard),
	}
}

func parseManagedClusterSKUTier(input string) (*ManagedClusterSKUTier, error) {
	vals := map[string]ManagedClusterSKUTier{
		"free":     ManagedClusterSKUTierFree,
		"premium":  ManagedClusterSKUTierPremium,
		"standard": ManagedClusterSKUTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedClusterSKUTier(input)
	return &out, nil
}

type NetworkMode string

const (
	NetworkModeBridge      NetworkMode = "bridge"
	NetworkModeTransparent NetworkMode = "transparent"
)

func PossibleValuesForNetworkMode() []string {
	return []string{
		string(NetworkModeBridge),
		string(NetworkModeTransparent),
	}
}

func parseNetworkMode(input string) (*NetworkMode, error) {
	vals := map[string]NetworkMode{
		"bridge":      NetworkModeBridge,
		"transparent": NetworkModeTransparent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkMode(input)
	return &out, nil
}

type NetworkPlugin string

const (
	NetworkPluginAzure   NetworkPlugin = "azure"
	NetworkPluginKubenet NetworkPlugin = "kubenet"
	NetworkPluginNone    NetworkPlugin = "none"
)

func PossibleValuesForNetworkPlugin() []string {
	return []string{
		string(NetworkPluginAzure),
		string(NetworkPluginKubenet),
		string(NetworkPluginNone),
	}
}

func parseNetworkPlugin(input string) (*NetworkPlugin, error) {
	vals := map[string]NetworkPlugin{
		"azure":   NetworkPluginAzure,
		"kubenet": NetworkPluginKubenet,
		"none":    NetworkPluginNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkPlugin(input)
	return &out, nil
}

type NetworkPluginMode string

const (
	NetworkPluginModeOverlay NetworkPluginMode = "overlay"
)

func PossibleValuesForNetworkPluginMode() []string {
	return []string{
		string(NetworkPluginModeOverlay),
	}
}

func parseNetworkPluginMode(input string) (*NetworkPluginMode, error) {
	vals := map[string]NetworkPluginMode{
		"overlay": NetworkPluginModeOverlay,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkPluginMode(input)
	return &out, nil
}

type NetworkPolicy string

const (
	NetworkPolicyAzure  NetworkPolicy = "azure"
	NetworkPolicyCalico NetworkPolicy = "calico"
	NetworkPolicyCilium NetworkPolicy = "cilium"
)

func PossibleValuesForNetworkPolicy() []string {
	return []string{
		string(NetworkPolicyAzure),
		string(NetworkPolicyCalico),
		string(NetworkPolicyCilium),
	}
}

func parseNetworkPolicy(input string) (*NetworkPolicy, error) {
	vals := map[string]NetworkPolicy{
		"azure":  NetworkPolicyAzure,
		"calico": NetworkPolicyCalico,
		"cilium": NetworkPolicyCilium,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkPolicy(input)
	return &out, nil
}

type SnapshotType string

const (
	SnapshotTypeManagedCluster SnapshotType = "ManagedCluster"
	SnapshotTypeNodePool       SnapshotType = "NodePool"
)

func PossibleValuesForSnapshotType() []string {
	return []string{
		string(SnapshotTypeManagedCluster),
		string(SnapshotTypeNodePool),
	}
}

func parseSnapshotType(input string) (*SnapshotType, error) {
	vals := map[string]SnapshotType{
		"managedcluster": SnapshotTypeManagedCluster,
		"nodepool":       SnapshotTypeNodePool,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SnapshotType(input)
	return &out, nil
}
