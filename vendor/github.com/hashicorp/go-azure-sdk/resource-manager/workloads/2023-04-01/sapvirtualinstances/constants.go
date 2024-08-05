package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationType string

const (
	ConfigurationTypeCreateAndMount ConfigurationType = "CreateAndMount"
	ConfigurationTypeMount          ConfigurationType = "Mount"
	ConfigurationTypeSkip           ConfigurationType = "Skip"
)

func PossibleValuesForConfigurationType() []string {
	return []string{
		string(ConfigurationTypeCreateAndMount),
		string(ConfigurationTypeMount),
		string(ConfigurationTypeSkip),
	}
}

func (s *ConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationType(input string) (*ConfigurationType, error) {
	vals := map[string]ConfigurationType{
		"createandmount": ConfigurationTypeCreateAndMount,
		"mount":          ConfigurationTypeMount,
		"skip":           ConfigurationTypeSkip,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationType(input)
	return &out, nil
}

type DiskSkuName string

const (
	DiskSkuNamePremiumLRS     DiskSkuName = "Premium_LRS"
	DiskSkuNamePremiumVTwoLRS DiskSkuName = "PremiumV2_LRS"
	DiskSkuNamePremiumZRS     DiskSkuName = "Premium_ZRS"
	DiskSkuNameStandardLRS    DiskSkuName = "Standard_LRS"
	DiskSkuNameStandardSSDLRS DiskSkuName = "StandardSSD_LRS"
	DiskSkuNameStandardSSDZRS DiskSkuName = "StandardSSD_ZRS"
	DiskSkuNameUltraSSDLRS    DiskSkuName = "UltraSSD_LRS"
)

func PossibleValuesForDiskSkuName() []string {
	return []string{
		string(DiskSkuNamePremiumLRS),
		string(DiskSkuNamePremiumVTwoLRS),
		string(DiskSkuNamePremiumZRS),
		string(DiskSkuNameStandardLRS),
		string(DiskSkuNameStandardSSDLRS),
		string(DiskSkuNameStandardSSDZRS),
		string(DiskSkuNameUltraSSDLRS),
	}
}

func (s *DiskSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskSkuName(input string) (*DiskSkuName, error) {
	vals := map[string]DiskSkuName{
		"premium_lrs":     DiskSkuNamePremiumLRS,
		"premiumv2_lrs":   DiskSkuNamePremiumVTwoLRS,
		"premium_zrs":     DiskSkuNamePremiumZRS,
		"standard_lrs":    DiskSkuNameStandardLRS,
		"standardssd_lrs": DiskSkuNameStandardSSDLRS,
		"standardssd_zrs": DiskSkuNameStandardSSDZRS,
		"ultrassd_lrs":    DiskSkuNameUltraSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskSkuName(input)
	return &out, nil
}

type NamingPatternType string

const (
	NamingPatternTypeFullResourceName NamingPatternType = "FullResourceName"
)

func PossibleValuesForNamingPatternType() []string {
	return []string{
		string(NamingPatternTypeFullResourceName),
	}
}

func (s *NamingPatternType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNamingPatternType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNamingPatternType(input string) (*NamingPatternType, error) {
	vals := map[string]NamingPatternType{
		"fullresourcename": NamingPatternTypeFullResourceName,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NamingPatternType(input)
	return &out, nil
}

type OSType string

const (
	OSTypeLinux   OSType = "Linux"
	OSTypeWindows OSType = "Windows"
)

func PossibleValuesForOSType() []string {
	return []string{
		string(OSTypeLinux),
		string(OSTypeWindows),
	}
}

func (s *OSType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOSType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOSType(input string) (*OSType, error) {
	vals := map[string]OSType{
		"linux":   OSTypeLinux,
		"windows": OSTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSType(input)
	return &out, nil
}

type SAPConfigurationType string

const (
	SAPConfigurationTypeDeployment             SAPConfigurationType = "Deployment"
	SAPConfigurationTypeDeploymentWithOSConfig SAPConfigurationType = "DeploymentWithOSConfig"
	SAPConfigurationTypeDiscovery              SAPConfigurationType = "Discovery"
)

func PossibleValuesForSAPConfigurationType() []string {
	return []string{
		string(SAPConfigurationTypeDeployment),
		string(SAPConfigurationTypeDeploymentWithOSConfig),
		string(SAPConfigurationTypeDiscovery),
	}
}

func (s *SAPConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPConfigurationType(input string) (*SAPConfigurationType, error) {
	vals := map[string]SAPConfigurationType{
		"deployment":             SAPConfigurationTypeDeployment,
		"deploymentwithosconfig": SAPConfigurationTypeDeploymentWithOSConfig,
		"discovery":              SAPConfigurationTypeDiscovery,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPConfigurationType(input)
	return &out, nil
}

type SAPDatabaseType string

const (
	SAPDatabaseTypeDBTwo SAPDatabaseType = "DB2"
	SAPDatabaseTypeHANA  SAPDatabaseType = "HANA"
)

func PossibleValuesForSAPDatabaseType() []string {
	return []string{
		string(SAPDatabaseTypeDBTwo),
		string(SAPDatabaseTypeHANA),
	}
}

func (s *SAPDatabaseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPDatabaseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPDatabaseType(input string) (*SAPDatabaseType, error) {
	vals := map[string]SAPDatabaseType{
		"db2":  SAPDatabaseTypeDBTwo,
		"hana": SAPDatabaseTypeHANA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPDatabaseType(input)
	return &out, nil
}

type SAPDeploymentType string

const (
	SAPDeploymentTypeSingleServer SAPDeploymentType = "SingleServer"
	SAPDeploymentTypeThreeTier    SAPDeploymentType = "ThreeTier"
)

func PossibleValuesForSAPDeploymentType() []string {
	return []string{
		string(SAPDeploymentTypeSingleServer),
		string(SAPDeploymentTypeThreeTier),
	}
}

func (s *SAPDeploymentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPDeploymentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPDeploymentType(input string) (*SAPDeploymentType, error) {
	vals := map[string]SAPDeploymentType{
		"singleserver": SAPDeploymentTypeSingleServer,
		"threetier":    SAPDeploymentTypeThreeTier,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPDeploymentType(input)
	return &out, nil
}

type SAPEnvironmentType string

const (
	SAPEnvironmentTypeNonProd SAPEnvironmentType = "NonProd"
	SAPEnvironmentTypeProd    SAPEnvironmentType = "Prod"
)

func PossibleValuesForSAPEnvironmentType() []string {
	return []string{
		string(SAPEnvironmentTypeNonProd),
		string(SAPEnvironmentTypeProd),
	}
}

func (s *SAPEnvironmentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPEnvironmentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPEnvironmentType(input string) (*SAPEnvironmentType, error) {
	vals := map[string]SAPEnvironmentType{
		"nonprod": SAPEnvironmentTypeNonProd,
		"prod":    SAPEnvironmentTypeProd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPEnvironmentType(input)
	return &out, nil
}

type SAPHealthState string

const (
	SAPHealthStateDegraded  SAPHealthState = "Degraded"
	SAPHealthStateHealthy   SAPHealthState = "Healthy"
	SAPHealthStateUnhealthy SAPHealthState = "Unhealthy"
	SAPHealthStateUnknown   SAPHealthState = "Unknown"
)

func PossibleValuesForSAPHealthState() []string {
	return []string{
		string(SAPHealthStateDegraded),
		string(SAPHealthStateHealthy),
		string(SAPHealthStateUnhealthy),
		string(SAPHealthStateUnknown),
	}
}

func (s *SAPHealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPHealthState(input string) (*SAPHealthState, error) {
	vals := map[string]SAPHealthState{
		"degraded":  SAPHealthStateDegraded,
		"healthy":   SAPHealthStateHealthy,
		"unhealthy": SAPHealthStateUnhealthy,
		"unknown":   SAPHealthStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPHealthState(input)
	return &out, nil
}

type SAPHighAvailabilityType string

const (
	SAPHighAvailabilityTypeAvailabilitySet  SAPHighAvailabilityType = "AvailabilitySet"
	SAPHighAvailabilityTypeAvailabilityZone SAPHighAvailabilityType = "AvailabilityZone"
)

func PossibleValuesForSAPHighAvailabilityType() []string {
	return []string{
		string(SAPHighAvailabilityTypeAvailabilitySet),
		string(SAPHighAvailabilityTypeAvailabilityZone),
	}
}

func (s *SAPHighAvailabilityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPHighAvailabilityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPHighAvailabilityType(input string) (*SAPHighAvailabilityType, error) {
	vals := map[string]SAPHighAvailabilityType{
		"availabilityset":  SAPHighAvailabilityTypeAvailabilitySet,
		"availabilityzone": SAPHighAvailabilityTypeAvailabilityZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPHighAvailabilityType(input)
	return &out, nil
}

type SAPProductType string

const (
	SAPProductTypeECC       SAPProductType = "ECC"
	SAPProductTypeOther     SAPProductType = "Other"
	SAPProductTypeSFourHANA SAPProductType = "S4HANA"
)

func PossibleValuesForSAPProductType() []string {
	return []string{
		string(SAPProductTypeECC),
		string(SAPProductTypeOther),
		string(SAPProductTypeSFourHANA),
	}
}

func (s *SAPProductType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPProductType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPProductType(input string) (*SAPProductType, error) {
	vals := map[string]SAPProductType{
		"ecc":    SAPProductTypeECC,
		"other":  SAPProductTypeOther,
		"s4hana": SAPProductTypeSFourHANA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPProductType(input)
	return &out, nil
}

type SAPSoftwareInstallationType string

const (
	SAPSoftwareInstallationTypeExternal                  SAPSoftwareInstallationType = "External"
	SAPSoftwareInstallationTypeSAPInstallWithoutOSConfig SAPSoftwareInstallationType = "SAPInstallWithoutOSConfig"
	SAPSoftwareInstallationTypeServiceInitiated          SAPSoftwareInstallationType = "ServiceInitiated"
)

func PossibleValuesForSAPSoftwareInstallationType() []string {
	return []string{
		string(SAPSoftwareInstallationTypeExternal),
		string(SAPSoftwareInstallationTypeSAPInstallWithoutOSConfig),
		string(SAPSoftwareInstallationTypeServiceInitiated),
	}
}

func (s *SAPSoftwareInstallationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPSoftwareInstallationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPSoftwareInstallationType(input string) (*SAPSoftwareInstallationType, error) {
	vals := map[string]SAPSoftwareInstallationType{
		"external":                  SAPSoftwareInstallationTypeExternal,
		"sapinstallwithoutosconfig": SAPSoftwareInstallationTypeSAPInstallWithoutOSConfig,
		"serviceinitiated":          SAPSoftwareInstallationTypeServiceInitiated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPSoftwareInstallationType(input)
	return &out, nil
}

type SAPVirtualInstanceState string

const (
	SAPVirtualInstanceStateDiscoveryFailed                    SAPVirtualInstanceState = "DiscoveryFailed"
	SAPVirtualInstanceStateDiscoveryInProgress                SAPVirtualInstanceState = "DiscoveryInProgress"
	SAPVirtualInstanceStateDiscoveryPending                   SAPVirtualInstanceState = "DiscoveryPending"
	SAPVirtualInstanceStateInfrastructureDeploymentFailed     SAPVirtualInstanceState = "InfrastructureDeploymentFailed"
	SAPVirtualInstanceStateInfrastructureDeploymentInProgress SAPVirtualInstanceState = "InfrastructureDeploymentInProgress"
	SAPVirtualInstanceStateInfrastructureDeploymentPending    SAPVirtualInstanceState = "InfrastructureDeploymentPending"
	SAPVirtualInstanceStateRegistrationComplete               SAPVirtualInstanceState = "RegistrationComplete"
	SAPVirtualInstanceStateSoftwareDetectionFailed            SAPVirtualInstanceState = "SoftwareDetectionFailed"
	SAPVirtualInstanceStateSoftwareDetectionInProgress        SAPVirtualInstanceState = "SoftwareDetectionInProgress"
	SAPVirtualInstanceStateSoftwareInstallationFailed         SAPVirtualInstanceState = "SoftwareInstallationFailed"
	SAPVirtualInstanceStateSoftwareInstallationInProgress     SAPVirtualInstanceState = "SoftwareInstallationInProgress"
	SAPVirtualInstanceStateSoftwareInstallationPending        SAPVirtualInstanceState = "SoftwareInstallationPending"
)

func PossibleValuesForSAPVirtualInstanceState() []string {
	return []string{
		string(SAPVirtualInstanceStateDiscoveryFailed),
		string(SAPVirtualInstanceStateDiscoveryInProgress),
		string(SAPVirtualInstanceStateDiscoveryPending),
		string(SAPVirtualInstanceStateInfrastructureDeploymentFailed),
		string(SAPVirtualInstanceStateInfrastructureDeploymentInProgress),
		string(SAPVirtualInstanceStateInfrastructureDeploymentPending),
		string(SAPVirtualInstanceStateRegistrationComplete),
		string(SAPVirtualInstanceStateSoftwareDetectionFailed),
		string(SAPVirtualInstanceStateSoftwareDetectionInProgress),
		string(SAPVirtualInstanceStateSoftwareInstallationFailed),
		string(SAPVirtualInstanceStateSoftwareInstallationInProgress),
		string(SAPVirtualInstanceStateSoftwareInstallationPending),
	}
}

func (s *SAPVirtualInstanceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPVirtualInstanceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPVirtualInstanceState(input string) (*SAPVirtualInstanceState, error) {
	vals := map[string]SAPVirtualInstanceState{
		"discoveryfailed":                    SAPVirtualInstanceStateDiscoveryFailed,
		"discoveryinprogress":                SAPVirtualInstanceStateDiscoveryInProgress,
		"discoverypending":                   SAPVirtualInstanceStateDiscoveryPending,
		"infrastructuredeploymentfailed":     SAPVirtualInstanceStateInfrastructureDeploymentFailed,
		"infrastructuredeploymentinprogress": SAPVirtualInstanceStateInfrastructureDeploymentInProgress,
		"infrastructuredeploymentpending":    SAPVirtualInstanceStateInfrastructureDeploymentPending,
		"registrationcomplete":               SAPVirtualInstanceStateRegistrationComplete,
		"softwaredetectionfailed":            SAPVirtualInstanceStateSoftwareDetectionFailed,
		"softwaredetectioninprogress":        SAPVirtualInstanceStateSoftwareDetectionInProgress,
		"softwareinstallationfailed":         SAPVirtualInstanceStateSoftwareInstallationFailed,
		"softwareinstallationinprogress":     SAPVirtualInstanceStateSoftwareInstallationInProgress,
		"softwareinstallationpending":        SAPVirtualInstanceStateSoftwareInstallationPending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPVirtualInstanceState(input)
	return &out, nil
}

type SAPVirtualInstanceStatus string

const (
	SAPVirtualInstanceStatusOffline          SAPVirtualInstanceStatus = "Offline"
	SAPVirtualInstanceStatusPartiallyRunning SAPVirtualInstanceStatus = "PartiallyRunning"
	SAPVirtualInstanceStatusRunning          SAPVirtualInstanceStatus = "Running"
	SAPVirtualInstanceStatusSoftShutdown     SAPVirtualInstanceStatus = "SoftShutdown"
	SAPVirtualInstanceStatusStarting         SAPVirtualInstanceStatus = "Starting"
	SAPVirtualInstanceStatusStopping         SAPVirtualInstanceStatus = "Stopping"
	SAPVirtualInstanceStatusUnavailable      SAPVirtualInstanceStatus = "Unavailable"
)

func PossibleValuesForSAPVirtualInstanceStatus() []string {
	return []string{
		string(SAPVirtualInstanceStatusOffline),
		string(SAPVirtualInstanceStatusPartiallyRunning),
		string(SAPVirtualInstanceStatusRunning),
		string(SAPVirtualInstanceStatusSoftShutdown),
		string(SAPVirtualInstanceStatusStarting),
		string(SAPVirtualInstanceStatusStopping),
		string(SAPVirtualInstanceStatusUnavailable),
	}
}

func (s *SAPVirtualInstanceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSAPVirtualInstanceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSAPVirtualInstanceStatus(input string) (*SAPVirtualInstanceStatus, error) {
	vals := map[string]SAPVirtualInstanceStatus{
		"offline":          SAPVirtualInstanceStatusOffline,
		"partiallyrunning": SAPVirtualInstanceStatusPartiallyRunning,
		"running":          SAPVirtualInstanceStatusRunning,
		"softshutdown":     SAPVirtualInstanceStatusSoftShutdown,
		"starting":         SAPVirtualInstanceStatusStarting,
		"stopping":         SAPVirtualInstanceStatusStopping,
		"unavailable":      SAPVirtualInstanceStatusUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SAPVirtualInstanceStatus(input)
	return &out, nil
}

type SapVirtualInstanceProvisioningState string

const (
	SapVirtualInstanceProvisioningStateCreating  SapVirtualInstanceProvisioningState = "Creating"
	SapVirtualInstanceProvisioningStateDeleting  SapVirtualInstanceProvisioningState = "Deleting"
	SapVirtualInstanceProvisioningStateFailed    SapVirtualInstanceProvisioningState = "Failed"
	SapVirtualInstanceProvisioningStateSucceeded SapVirtualInstanceProvisioningState = "Succeeded"
	SapVirtualInstanceProvisioningStateUpdating  SapVirtualInstanceProvisioningState = "Updating"
)

func PossibleValuesForSapVirtualInstanceProvisioningState() []string {
	return []string{
		string(SapVirtualInstanceProvisioningStateCreating),
		string(SapVirtualInstanceProvisioningStateDeleting),
		string(SapVirtualInstanceProvisioningStateFailed),
		string(SapVirtualInstanceProvisioningStateSucceeded),
		string(SapVirtualInstanceProvisioningStateUpdating),
	}
}

func (s *SapVirtualInstanceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSapVirtualInstanceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSapVirtualInstanceProvisioningState(input string) (*SapVirtualInstanceProvisioningState, error) {
	vals := map[string]SapVirtualInstanceProvisioningState{
		"creating":  SapVirtualInstanceProvisioningStateCreating,
		"deleting":  SapVirtualInstanceProvisioningStateDeleting,
		"failed":    SapVirtualInstanceProvisioningStateFailed,
		"succeeded": SapVirtualInstanceProvisioningStateSucceeded,
		"updating":  SapVirtualInstanceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SapVirtualInstanceProvisioningState(input)
	return &out, nil
}
