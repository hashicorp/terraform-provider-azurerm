package services

import "strings"

type MoveCost string

const (
	MoveCostHigh   MoveCost = "High"
	MoveCostLow    MoveCost = "Low"
	MoveCostMedium MoveCost = "Medium"
	MoveCostZero   MoveCost = "Zero"
)

func PossibleValuesForMoveCost() []string {
	return []string{
		string(MoveCostHigh),
		string(MoveCostLow),
		string(MoveCostMedium),
		string(MoveCostZero),
	}
}

func parseMoveCost(input string) (*MoveCost, error) {
	vals := map[string]MoveCost{
		"high":   MoveCostHigh,
		"low":    MoveCostLow,
		"medium": MoveCostMedium,
		"zero":   MoveCostZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MoveCost(input)
	return &out, nil
}

type PartitionScheme string

const (
	PartitionSchemeNamed                  PartitionScheme = "Named"
	PartitionSchemeSingleton              PartitionScheme = "Singleton"
	PartitionSchemeUniformIntSixFourRange PartitionScheme = "UniformInt64Range"
)

func PossibleValuesForPartitionScheme() []string {
	return []string{
		string(PartitionSchemeNamed),
		string(PartitionSchemeSingleton),
		string(PartitionSchemeUniformIntSixFourRange),
	}
}

func parsePartitionScheme(input string) (*PartitionScheme, error) {
	vals := map[string]PartitionScheme{
		"named":             PartitionSchemeNamed,
		"singleton":         PartitionSchemeSingleton,
		"uniformint64range": PartitionSchemeUniformIntSixFourRange,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartitionScheme(input)
	return &out, nil
}

type ServiceCorrelationScheme string

const (
	ServiceCorrelationSchemeAlignedAffinity    ServiceCorrelationScheme = "AlignedAffinity"
	ServiceCorrelationSchemeNonAlignedAffinity ServiceCorrelationScheme = "NonAlignedAffinity"
)

func PossibleValuesForServiceCorrelationScheme() []string {
	return []string{
		string(ServiceCorrelationSchemeAlignedAffinity),
		string(ServiceCorrelationSchemeNonAlignedAffinity),
	}
}

func parseServiceCorrelationScheme(input string) (*ServiceCorrelationScheme, error) {
	vals := map[string]ServiceCorrelationScheme{
		"alignedaffinity":    ServiceCorrelationSchemeAlignedAffinity,
		"nonalignedaffinity": ServiceCorrelationSchemeNonAlignedAffinity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceCorrelationScheme(input)
	return &out, nil
}

type ServiceKind string

const (
	ServiceKindStateful  ServiceKind = "Stateful"
	ServiceKindStateless ServiceKind = "Stateless"
)

func PossibleValuesForServiceKind() []string {
	return []string{
		string(ServiceKindStateful),
		string(ServiceKindStateless),
	}
}

func parseServiceKind(input string) (*ServiceKind, error) {
	vals := map[string]ServiceKind{
		"stateful":  ServiceKindStateful,
		"stateless": ServiceKindStateless,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceKind(input)
	return &out, nil
}

type ServiceLoadMetricWeight string

const (
	ServiceLoadMetricWeightHigh   ServiceLoadMetricWeight = "High"
	ServiceLoadMetricWeightLow    ServiceLoadMetricWeight = "Low"
	ServiceLoadMetricWeightMedium ServiceLoadMetricWeight = "Medium"
	ServiceLoadMetricWeightZero   ServiceLoadMetricWeight = "Zero"
)

func PossibleValuesForServiceLoadMetricWeight() []string {
	return []string{
		string(ServiceLoadMetricWeightHigh),
		string(ServiceLoadMetricWeightLow),
		string(ServiceLoadMetricWeightMedium),
		string(ServiceLoadMetricWeightZero),
	}
}

func parseServiceLoadMetricWeight(input string) (*ServiceLoadMetricWeight, error) {
	vals := map[string]ServiceLoadMetricWeight{
		"high":   ServiceLoadMetricWeightHigh,
		"low":    ServiceLoadMetricWeightLow,
		"medium": ServiceLoadMetricWeightMedium,
		"zero":   ServiceLoadMetricWeightZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceLoadMetricWeight(input)
	return &out, nil
}

type ServicePackageActivationMode string

const (
	ServicePackageActivationModeExclusiveProcess ServicePackageActivationMode = "ExclusiveProcess"
	ServicePackageActivationModeSharedProcess    ServicePackageActivationMode = "SharedProcess"
)

func PossibleValuesForServicePackageActivationMode() []string {
	return []string{
		string(ServicePackageActivationModeExclusiveProcess),
		string(ServicePackageActivationModeSharedProcess),
	}
}

func parseServicePackageActivationMode(input string) (*ServicePackageActivationMode, error) {
	vals := map[string]ServicePackageActivationMode{
		"exclusiveprocess": ServicePackageActivationModeExclusiveProcess,
		"sharedprocess":    ServicePackageActivationModeSharedProcess,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServicePackageActivationMode(input)
	return &out, nil
}

type ServicePlacementPolicyType string

const (
	ServicePlacementPolicyTypeInvalidDomain              ServicePlacementPolicyType = "InvalidDomain"
	ServicePlacementPolicyTypeNonPartiallyPlaceService   ServicePlacementPolicyType = "NonPartiallyPlaceService"
	ServicePlacementPolicyTypePreferredPrimaryDomain     ServicePlacementPolicyType = "PreferredPrimaryDomain"
	ServicePlacementPolicyTypeRequiredDomain             ServicePlacementPolicyType = "RequiredDomain"
	ServicePlacementPolicyTypeRequiredDomainDistribution ServicePlacementPolicyType = "RequiredDomainDistribution"
)

func PossibleValuesForServicePlacementPolicyType() []string {
	return []string{
		string(ServicePlacementPolicyTypeInvalidDomain),
		string(ServicePlacementPolicyTypeNonPartiallyPlaceService),
		string(ServicePlacementPolicyTypePreferredPrimaryDomain),
		string(ServicePlacementPolicyTypeRequiredDomain),
		string(ServicePlacementPolicyTypeRequiredDomainDistribution),
	}
}

func parseServicePlacementPolicyType(input string) (*ServicePlacementPolicyType, error) {
	vals := map[string]ServicePlacementPolicyType{
		"invaliddomain":              ServicePlacementPolicyTypeInvalidDomain,
		"nonpartiallyplaceservice":   ServicePlacementPolicyTypeNonPartiallyPlaceService,
		"preferredprimarydomain":     ServicePlacementPolicyTypePreferredPrimaryDomain,
		"requireddomain":             ServicePlacementPolicyTypeRequiredDomain,
		"requireddomaindistribution": ServicePlacementPolicyTypeRequiredDomainDistribution,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServicePlacementPolicyType(input)
	return &out, nil
}

type ServiceScalingMechanismKind string

const (
	ServiceScalingMechanismKindAddRemoveIncrementalNamedPartition ServiceScalingMechanismKind = "AddRemoveIncrementalNamedPartition"
	ServiceScalingMechanismKindScalePartitionInstanceCount        ServiceScalingMechanismKind = "ScalePartitionInstanceCount"
)

func PossibleValuesForServiceScalingMechanismKind() []string {
	return []string{
		string(ServiceScalingMechanismKindAddRemoveIncrementalNamedPartition),
		string(ServiceScalingMechanismKindScalePartitionInstanceCount),
	}
}

func parseServiceScalingMechanismKind(input string) (*ServiceScalingMechanismKind, error) {
	vals := map[string]ServiceScalingMechanismKind{
		"addremoveincrementalnamedpartition": ServiceScalingMechanismKindAddRemoveIncrementalNamedPartition,
		"scalepartitioninstancecount":        ServiceScalingMechanismKindScalePartitionInstanceCount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceScalingMechanismKind(input)
	return &out, nil
}

type ServiceScalingTriggerKind string

const (
	ServiceScalingTriggerKindAveragePartitionLoad ServiceScalingTriggerKind = "AveragePartitionLoad"
	ServiceScalingTriggerKindAverageServiceLoad   ServiceScalingTriggerKind = "AverageServiceLoad"
)

func PossibleValuesForServiceScalingTriggerKind() []string {
	return []string{
		string(ServiceScalingTriggerKindAveragePartitionLoad),
		string(ServiceScalingTriggerKindAverageServiceLoad),
	}
}

func parseServiceScalingTriggerKind(input string) (*ServiceScalingTriggerKind, error) {
	vals := map[string]ServiceScalingTriggerKind{
		"averagepartitionload": ServiceScalingTriggerKindAveragePartitionLoad,
		"averageserviceload":   ServiceScalingTriggerKindAverageServiceLoad,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceScalingTriggerKind(input)
	return &out, nil
}
