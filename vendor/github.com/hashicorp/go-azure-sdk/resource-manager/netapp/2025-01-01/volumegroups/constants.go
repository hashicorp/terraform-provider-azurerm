package volumegroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationType string

const (
	ApplicationTypeORACLE          ApplicationType = "ORACLE"
	ApplicationTypeSAPNegativeHANA ApplicationType = "SAP-HANA"
)

func PossibleValuesForApplicationType() []string {
	return []string{
		string(ApplicationTypeORACLE),
		string(ApplicationTypeSAPNegativeHANA),
	}
}

func (s *ApplicationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationType(input string) (*ApplicationType, error) {
	vals := map[string]ApplicationType{
		"oracle":   ApplicationTypeORACLE,
		"sap-hana": ApplicationTypeSAPNegativeHANA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationType(input)
	return &out, nil
}

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

func (s *AvsDataStore) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAvsDataStore(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *ChownMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseChownMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type CoolAccessRetrievalPolicy string

const (
	CoolAccessRetrievalPolicyDefault CoolAccessRetrievalPolicy = "Default"
	CoolAccessRetrievalPolicyNever   CoolAccessRetrievalPolicy = "Never"
	CoolAccessRetrievalPolicyOnRead  CoolAccessRetrievalPolicy = "OnRead"
)

func PossibleValuesForCoolAccessRetrievalPolicy() []string {
	return []string{
		string(CoolAccessRetrievalPolicyDefault),
		string(CoolAccessRetrievalPolicyNever),
		string(CoolAccessRetrievalPolicyOnRead),
	}
}

func (s *CoolAccessRetrievalPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCoolAccessRetrievalPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCoolAccessRetrievalPolicy(input string) (*CoolAccessRetrievalPolicy, error) {
	vals := map[string]CoolAccessRetrievalPolicy{
		"default": CoolAccessRetrievalPolicyDefault,
		"never":   CoolAccessRetrievalPolicyNever,
		"onread":  CoolAccessRetrievalPolicyOnRead,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CoolAccessRetrievalPolicy(input)
	return &out, nil
}

type CoolAccessTieringPolicy string

const (
	CoolAccessTieringPolicyAuto         CoolAccessTieringPolicy = "Auto"
	CoolAccessTieringPolicySnapshotOnly CoolAccessTieringPolicy = "SnapshotOnly"
)

func PossibleValuesForCoolAccessTieringPolicy() []string {
	return []string{
		string(CoolAccessTieringPolicyAuto),
		string(CoolAccessTieringPolicySnapshotOnly),
	}
}

func (s *CoolAccessTieringPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCoolAccessTieringPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCoolAccessTieringPolicy(input string) (*CoolAccessTieringPolicy, error) {
	vals := map[string]CoolAccessTieringPolicy{
		"auto":         CoolAccessTieringPolicyAuto,
		"snapshotonly": CoolAccessTieringPolicySnapshotOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CoolAccessTieringPolicy(input)
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

func (s *EnableSubvolumes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnableSubvolumes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type EncryptionKeySource string

const (
	EncryptionKeySourceMicrosoftPointKeyVault EncryptionKeySource = "Microsoft.KeyVault"
	EncryptionKeySourceMicrosoftPointNetApp   EncryptionKeySource = "Microsoft.NetApp"
)

func PossibleValuesForEncryptionKeySource() []string {
	return []string{
		string(EncryptionKeySourceMicrosoftPointKeyVault),
		string(EncryptionKeySourceMicrosoftPointNetApp),
	}
}

func (s *EncryptionKeySource) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionKeySource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionKeySource(input string) (*EncryptionKeySource, error) {
	vals := map[string]EncryptionKeySource{
		"microsoft.keyvault": EncryptionKeySourceMicrosoftPointKeyVault,
		"microsoft.netapp":   EncryptionKeySourceMicrosoftPointNetApp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionKeySource(input)
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

func (s *EndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type FileAccessLogs string

const (
	FileAccessLogsDisabled FileAccessLogs = "Disabled"
	FileAccessLogsEnabled  FileAccessLogs = "Enabled"
)

func PossibleValuesForFileAccessLogs() []string {
	return []string{
		string(FileAccessLogsDisabled),
		string(FileAccessLogsEnabled),
	}
}

func (s *FileAccessLogs) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFileAccessLogs(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFileAccessLogs(input string) (*FileAccessLogs, error) {
	vals := map[string]FileAccessLogs{
		"disabled": FileAccessLogsDisabled,
		"enabled":  FileAccessLogsEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FileAccessLogs(input)
	return &out, nil
}

type NetworkFeatures string

const (
	NetworkFeaturesBasic         NetworkFeatures = "Basic"
	NetworkFeaturesBasicStandard NetworkFeatures = "Basic_Standard"
	NetworkFeaturesStandard      NetworkFeatures = "Standard"
	NetworkFeaturesStandardBasic NetworkFeatures = "Standard_Basic"
)

func PossibleValuesForNetworkFeatures() []string {
	return []string{
		string(NetworkFeaturesBasic),
		string(NetworkFeaturesBasicStandard),
		string(NetworkFeaturesStandard),
		string(NetworkFeaturesStandardBasic),
	}
}

func (s *NetworkFeatures) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkFeatures(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkFeatures(input string) (*NetworkFeatures, error) {
	vals := map[string]NetworkFeatures{
		"basic":          NetworkFeaturesBasic,
		"basic_standard": NetworkFeaturesBasicStandard,
		"standard":       NetworkFeaturesStandard,
		"standard_basic": NetworkFeaturesStandardBasic,
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

func (s *ReplicationSchedule) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationSchedule(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type ReplicationType string

const (
	ReplicationTypeCrossRegionReplication ReplicationType = "CrossRegionReplication"
	ReplicationTypeCrossZoneReplication   ReplicationType = "CrossZoneReplication"
)

func PossibleValuesForReplicationType() []string {
	return []string{
		string(ReplicationTypeCrossRegionReplication),
		string(ReplicationTypeCrossZoneReplication),
	}
}

func (s *ReplicationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationType(input string) (*ReplicationType, error) {
	vals := map[string]ReplicationType{
		"crossregionreplication": ReplicationTypeCrossRegionReplication,
		"crosszonereplication":   ReplicationTypeCrossZoneReplication,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationType(input)
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

func (s *SecurityStyle) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityStyle(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

func (s *ServiceLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
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

type SmbAccessBasedEnumeration string

const (
	SmbAccessBasedEnumerationDisabled SmbAccessBasedEnumeration = "Disabled"
	SmbAccessBasedEnumerationEnabled  SmbAccessBasedEnumeration = "Enabled"
)

func PossibleValuesForSmbAccessBasedEnumeration() []string {
	return []string{
		string(SmbAccessBasedEnumerationDisabled),
		string(SmbAccessBasedEnumerationEnabled),
	}
}

func (s *SmbAccessBasedEnumeration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSmbAccessBasedEnumeration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSmbAccessBasedEnumeration(input string) (*SmbAccessBasedEnumeration, error) {
	vals := map[string]SmbAccessBasedEnumeration{
		"disabled": SmbAccessBasedEnumerationDisabled,
		"enabled":  SmbAccessBasedEnumerationEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SmbAccessBasedEnumeration(input)
	return &out, nil
}

type SmbNonBrowsable string

const (
	SmbNonBrowsableDisabled SmbNonBrowsable = "Disabled"
	SmbNonBrowsableEnabled  SmbNonBrowsable = "Enabled"
)

func PossibleValuesForSmbNonBrowsable() []string {
	return []string{
		string(SmbNonBrowsableDisabled),
		string(SmbNonBrowsableEnabled),
	}
}

func (s *SmbNonBrowsable) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSmbNonBrowsable(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSmbNonBrowsable(input string) (*SmbNonBrowsable, error) {
	vals := map[string]SmbNonBrowsable{
		"disabled": SmbNonBrowsableDisabled,
		"enabled":  SmbNonBrowsableEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SmbNonBrowsable(input)
	return &out, nil
}

type VolumeStorageToNetworkProximity string

const (
	VolumeStorageToNetworkProximityAcrossTTwo VolumeStorageToNetworkProximity = "AcrossT2"
	VolumeStorageToNetworkProximityDefault    VolumeStorageToNetworkProximity = "Default"
	VolumeStorageToNetworkProximityTOne       VolumeStorageToNetworkProximity = "T1"
	VolumeStorageToNetworkProximityTTwo       VolumeStorageToNetworkProximity = "T2"
)

func PossibleValuesForVolumeStorageToNetworkProximity() []string {
	return []string{
		string(VolumeStorageToNetworkProximityAcrossTTwo),
		string(VolumeStorageToNetworkProximityDefault),
		string(VolumeStorageToNetworkProximityTOne),
		string(VolumeStorageToNetworkProximityTTwo),
	}
}

func (s *VolumeStorageToNetworkProximity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVolumeStorageToNetworkProximity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVolumeStorageToNetworkProximity(input string) (*VolumeStorageToNetworkProximity, error) {
	vals := map[string]VolumeStorageToNetworkProximity{
		"acrosst2": VolumeStorageToNetworkProximityAcrossTTwo,
		"default":  VolumeStorageToNetworkProximityDefault,
		"t1":       VolumeStorageToNetworkProximityTOne,
		"t2":       VolumeStorageToNetworkProximityTTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VolumeStorageToNetworkProximity(input)
	return &out, nil
}
