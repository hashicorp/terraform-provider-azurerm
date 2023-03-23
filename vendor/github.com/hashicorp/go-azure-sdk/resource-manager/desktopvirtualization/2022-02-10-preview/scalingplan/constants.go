package scalingplan

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DaysOfWeek string

const (
	DaysOfWeekFriday    DaysOfWeek = "Friday"
	DaysOfWeekMonday    DaysOfWeek = "Monday"
	DaysOfWeekSaturday  DaysOfWeek = "Saturday"
	DaysOfWeekSunday    DaysOfWeek = "Sunday"
	DaysOfWeekThursday  DaysOfWeek = "Thursday"
	DaysOfWeekTuesday   DaysOfWeek = "Tuesday"
	DaysOfWeekWednesday DaysOfWeek = "Wednesday"
)

func PossibleValuesForDaysOfWeek() []string {
	return []string{
		string(DaysOfWeekFriday),
		string(DaysOfWeekMonday),
		string(DaysOfWeekSaturday),
		string(DaysOfWeekSunday),
		string(DaysOfWeekThursday),
		string(DaysOfWeekTuesday),
		string(DaysOfWeekWednesday),
	}
}

func parseDaysOfWeek(input string) (*DaysOfWeek, error) {
	vals := map[string]DaysOfWeek{
		"friday":    DaysOfWeekFriday,
		"monday":    DaysOfWeekMonday,
		"saturday":  DaysOfWeekSaturday,
		"sunday":    DaysOfWeekSunday,
		"thursday":  DaysOfWeekThursday,
		"tuesday":   DaysOfWeekTuesday,
		"wednesday": DaysOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DaysOfWeek(input)
	return &out, nil
}

type ScalingHostPoolType string

const (
	ScalingHostPoolTypePooled ScalingHostPoolType = "Pooled"
)

func PossibleValuesForScalingHostPoolType() []string {
	return []string{
		string(ScalingHostPoolTypePooled),
	}
}

func parseScalingHostPoolType(input string) (*ScalingHostPoolType, error) {
	vals := map[string]ScalingHostPoolType{
		"pooled": ScalingHostPoolTypePooled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScalingHostPoolType(input)
	return &out, nil
}

type SessionHostLoadBalancingAlgorithm string

const (
	SessionHostLoadBalancingAlgorithmBreadthFirst SessionHostLoadBalancingAlgorithm = "BreadthFirst"
	SessionHostLoadBalancingAlgorithmDepthFirst   SessionHostLoadBalancingAlgorithm = "DepthFirst"
)

func PossibleValuesForSessionHostLoadBalancingAlgorithm() []string {
	return []string{
		string(SessionHostLoadBalancingAlgorithmBreadthFirst),
		string(SessionHostLoadBalancingAlgorithmDepthFirst),
	}
}

func parseSessionHostLoadBalancingAlgorithm(input string) (*SessionHostLoadBalancingAlgorithm, error) {
	vals := map[string]SessionHostLoadBalancingAlgorithm{
		"breadthfirst": SessionHostLoadBalancingAlgorithmBreadthFirst,
		"depthfirst":   SessionHostLoadBalancingAlgorithmDepthFirst,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SessionHostLoadBalancingAlgorithm(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type StopHostsWhen string

const (
	StopHostsWhenZeroActiveSessions StopHostsWhen = "ZeroActiveSessions"
	StopHostsWhenZeroSessions       StopHostsWhen = "ZeroSessions"
)

func PossibleValuesForStopHostsWhen() []string {
	return []string{
		string(StopHostsWhenZeroActiveSessions),
		string(StopHostsWhenZeroSessions),
	}
}

func parseStopHostsWhen(input string) (*StopHostsWhen, error) {
	vals := map[string]StopHostsWhen{
		"zeroactivesessions": StopHostsWhenZeroActiveSessions,
		"zerosessions":       StopHostsWhenZeroSessions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StopHostsWhen(input)
	return &out, nil
}
