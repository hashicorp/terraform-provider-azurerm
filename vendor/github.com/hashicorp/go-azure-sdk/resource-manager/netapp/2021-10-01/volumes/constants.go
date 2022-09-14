package volumes

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvsDataStore string

const (
	AvsDataStoreDisabled AvsDataStore = "Disabled"
	AvsDataStoreEnabled  AvsDataStore = "Enabled"
)

func PossibleValuesForAvsDataStore() []string {
	return []string{
		string(AvsDataStoreDisabled),
		string(AvsDataStoreEnabled),
	}
}

func parseAvsDataStore(input string) (*AvsDataStore, error) {
	vals := map[string]AvsDataStore{
		"disabled": AvsDataStoreDisabled,
		"enabled":  AvsDataStoreEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AvsDataStore(input)
	return &out, nil
}

type ChownMode string

const (
	ChownModeRestricted   ChownMode = "Restricted"
	ChownModeUnrestricted ChownMode = "Unrestricted"
)

func PossibleValuesForChownMode() []string {
	return []string{
		string(ChownModeRestricted),
		string(ChownModeUnrestricted),
	}
}

func parseChownMode(input string) (*ChownMode, error) {
	vals := map[string]ChownMode{
		"restricted":   ChownModeRestricted,
		"unrestricted": ChownModeUnrestricted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ChownMode(input)
	return &out, nil
}

type EnableSubvolumes string

const (
	EnableSubvolumesDisabled EnableSubvolumes = "Disabled"
	EnableSubvolumesEnabled  EnableSubvolumes = "Enabled"
)

func PossibleValuesForEnableSubvolumes() []string {
	return []string{
		string(EnableSubvolumesDisabled),
		string(EnableSubvolumesEnabled),
	}
}

func parseEnableSubvolumes(input string) (*EnableSubvolumes, error) {
	vals := map[string]EnableSubvolumes{
		"disabled": EnableSubvolumesDisabled,
		"enabled":  EnableSubvolumesEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnableSubvolumes(input)
	return &out, nil
}

type EndpointType string

const (
	EndpointTypeDst EndpointType = "dst"
	EndpointTypeSrc EndpointType = "src"
)

func PossibleValuesForEndpointType() []string {
	return []string{
		string(EndpointTypeDst),
		string(EndpointTypeSrc),
	}
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"dst": EndpointTypeDst,
		"src": EndpointTypeSrc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointType(input)
	return &out, nil
}

type NetworkFeatures string

const (
	NetworkFeaturesBasic    NetworkFeatures = "Basic"
	NetworkFeaturesStandard NetworkFeatures = "Standard"
)

func PossibleValuesForNetworkFeatures() []string {
	return []string{
		string(NetworkFeaturesBasic),
		string(NetworkFeaturesStandard),
	}
}

func parseNetworkFeatures(input string) (*NetworkFeatures, error) {
	vals := map[string]NetworkFeatures{
		"basic":    NetworkFeaturesBasic,
		"standard": NetworkFeaturesStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkFeatures(input)
	return &out, nil
}

type ReplicationSchedule string

const (
	ReplicationScheduleDaily           ReplicationSchedule = "daily"
	ReplicationScheduleHourly          ReplicationSchedule = "hourly"
	ReplicationScheduleOneZerominutely ReplicationSchedule = "_10minutely"
)

func PossibleValuesForReplicationSchedule() []string {
	return []string{
		string(ReplicationScheduleDaily),
		string(ReplicationScheduleHourly),
		string(ReplicationScheduleOneZerominutely),
	}
}

func parseReplicationSchedule(input string) (*ReplicationSchedule, error) {
	vals := map[string]ReplicationSchedule{
		"daily":       ReplicationScheduleDaily,
		"hourly":      ReplicationScheduleHourly,
		"_10minutely": ReplicationScheduleOneZerominutely,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationSchedule(input)
	return &out, nil
}

type SecurityStyle string

const (
	SecurityStyleNtfs SecurityStyle = "ntfs"
	SecurityStyleUnix SecurityStyle = "unix"
)

func PossibleValuesForSecurityStyle() []string {
	return []string{
		string(SecurityStyleNtfs),
		string(SecurityStyleUnix),
	}
}

func parseSecurityStyle(input string) (*SecurityStyle, error) {
	vals := map[string]SecurityStyle{
		"ntfs": SecurityStyleNtfs,
		"unix": SecurityStyleUnix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityStyle(input)
	return &out, nil
}

type ServiceLevel string

const (
	ServiceLevelPremium     ServiceLevel = "Premium"
	ServiceLevelStandard    ServiceLevel = "Standard"
	ServiceLevelStandardZRS ServiceLevel = "StandardZRS"
	ServiceLevelUltra       ServiceLevel = "Ultra"
)

func PossibleValuesForServiceLevel() []string {
	return []string{
		string(ServiceLevelPremium),
		string(ServiceLevelStandard),
		string(ServiceLevelStandardZRS),
		string(ServiceLevelUltra),
	}
}

func parseServiceLevel(input string) (*ServiceLevel, error) {
	vals := map[string]ServiceLevel{
		"premium":     ServiceLevelPremium,
		"standard":    ServiceLevelStandard,
		"standardzrs": ServiceLevelStandardZRS,
		"ultra":       ServiceLevelUltra,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceLevel(input)
	return &out, nil
}

type VolumeStorageToNetworkProximity string

const (
	VolumeStorageToNetworkProximityDefault VolumeStorageToNetworkProximity = "Default"
	VolumeStorageToNetworkProximityTOne    VolumeStorageToNetworkProximity = "T1"
	VolumeStorageToNetworkProximityTTwo    VolumeStorageToNetworkProximity = "T2"
)

func PossibleValuesForVolumeStorageToNetworkProximity() []string {
	return []string{
		string(VolumeStorageToNetworkProximityDefault),
		string(VolumeStorageToNetworkProximityTOne),
		string(VolumeStorageToNetworkProximityTTwo),
	}
}

func parseVolumeStorageToNetworkProximity(input string) (*VolumeStorageToNetworkProximity, error) {
	vals := map[string]VolumeStorageToNetworkProximity{
		"default": VolumeStorageToNetworkProximityDefault,
		"t1":      VolumeStorageToNetworkProximityTOne,
		"t2":      VolumeStorageToNetworkProximityTTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VolumeStorageToNetworkProximity(input)
	return &out, nil
}
