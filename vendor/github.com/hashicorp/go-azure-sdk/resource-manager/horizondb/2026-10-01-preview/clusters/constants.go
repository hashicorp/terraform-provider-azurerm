package clusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AiModelManagementStatus string

const (
	AiModelManagementStatusDisabled AiModelManagementStatus = "Disabled"
	AiModelManagementStatusEnabled  AiModelManagementStatus = "Enabled"
)

func PossibleValuesForAiModelManagementStatus() []string {
	return []string{
		string(AiModelManagementStatusDisabled),
		string(AiModelManagementStatusEnabled),
	}
}

func (s *AiModelManagementStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAiModelManagementStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAiModelManagementStatus(input string) (*AiModelManagementStatus, error) {
	vals := map[string]AiModelManagementStatus{
		"disabled": AiModelManagementStatusDisabled,
		"enabled":  AiModelManagementStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AiModelManagementStatus(input)
	return &out, nil
}

type ClusterCreateMode string

const (
	ClusterCreateModeCreate             ClusterCreateMode = "Create"
	ClusterCreateModePointInTimeRestore ClusterCreateMode = "PointInTimeRestore"
)

func PossibleValuesForClusterCreateMode() []string {
	return []string{
		string(ClusterCreateModeCreate),
		string(ClusterCreateModePointInTimeRestore),
	}
}

func (s *ClusterCreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterCreateMode(input string) (*ClusterCreateMode, error) {
	vals := map[string]ClusterCreateMode{
		"create":             ClusterCreateModeCreate,
		"pointintimerestore": ClusterCreateModePointInTimeRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterCreateMode(input)
	return &out, nil
}

type ClusterState string

const (
	ClusterStateDeleting     ClusterState = "Deleting"
	ClusterStateInaccessible ClusterState = "Inaccessible"
	ClusterStateProvisioning ClusterState = "Provisioning"
	ClusterStateReady        ClusterState = "Ready"
	ClusterStateUpdating     ClusterState = "Updating"
)

func PossibleValuesForClusterState() []string {
	return []string{
		string(ClusterStateDeleting),
		string(ClusterStateInaccessible),
		string(ClusterStateProvisioning),
		string(ClusterStateReady),
		string(ClusterStateUpdating),
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
		"deleting":     ClusterStateDeleting,
		"inaccessible": ClusterStateInaccessible,
		"provisioning": ClusterStateProvisioning,
		"ready":        ClusterStateReady,
		"updating":     ClusterStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterState(input)
	return &out, nil
}

type GeoRedundancyState string

const (
	GeoRedundancyStateDisabled GeoRedundancyState = "Disabled"
	GeoRedundancyStateEnabled  GeoRedundancyState = "Enabled"
)

func PossibleValuesForGeoRedundancyState() []string {
	return []string{
		string(GeoRedundancyStateDisabled),
		string(GeoRedundancyStateEnabled),
	}
}

func (s *GeoRedundancyState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGeoRedundancyState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGeoRedundancyState(input string) (*GeoRedundancyState, error) {
	vals := map[string]GeoRedundancyState{
		"disabled": GeoRedundancyStateDisabled,
		"enabled":  GeoRedundancyStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GeoRedundancyState(input)
	return &out, nil
}

type HorizonDbClusterDataEncryptionType string

const (
	HorizonDbClusterDataEncryptionTypeAzureKeyVault HorizonDbClusterDataEncryptionType = "AzureKeyVault"
	HorizonDbClusterDataEncryptionTypeSystemManaged HorizonDbClusterDataEncryptionType = "SystemManaged"
)

func PossibleValuesForHorizonDbClusterDataEncryptionType() []string {
	return []string{
		string(HorizonDbClusterDataEncryptionTypeAzureKeyVault),
		string(HorizonDbClusterDataEncryptionTypeSystemManaged),
	}
}

func (s *HorizonDbClusterDataEncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHorizonDbClusterDataEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHorizonDbClusterDataEncryptionType(input string) (*HorizonDbClusterDataEncryptionType, error) {
	vals := map[string]HorizonDbClusterDataEncryptionType{
		"azurekeyvault": HorizonDbClusterDataEncryptionTypeAzureKeyVault,
		"systemmanaged": HorizonDbClusterDataEncryptionTypeSystemManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HorizonDbClusterDataEncryptionType(input)
	return &out, nil
}

type HorizonDbClusterEncryptionKeyStatus string

const (
	HorizonDbClusterEncryptionKeyStatusInvalid HorizonDbClusterEncryptionKeyStatus = "Invalid"
	HorizonDbClusterEncryptionKeyStatusValid   HorizonDbClusterEncryptionKeyStatus = "Valid"
)

func PossibleValuesForHorizonDbClusterEncryptionKeyStatus() []string {
	return []string{
		string(HorizonDbClusterEncryptionKeyStatusInvalid),
		string(HorizonDbClusterEncryptionKeyStatusValid),
	}
}

func (s *HorizonDbClusterEncryptionKeyStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHorizonDbClusterEncryptionKeyStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHorizonDbClusterEncryptionKeyStatus(input string) (*HorizonDbClusterEncryptionKeyStatus, error) {
	vals := map[string]HorizonDbClusterEncryptionKeyStatus{
		"invalid": HorizonDbClusterEncryptionKeyStatusInvalid,
		"valid":   HorizonDbClusterEncryptionKeyStatusValid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HorizonDbClusterEncryptionKeyStatus(input)
	return &out, nil
}

type PostgreSqlVersion string

const (
	PostgreSqlVersionOneEight PostgreSqlVersion = "18"
	PostgreSqlVersionOneSeven PostgreSqlVersion = "17"
)

func PossibleValuesForPostgreSqlVersion() []string {
	return []string{
		string(PostgreSqlVersionOneEight),
		string(PostgreSqlVersionOneSeven),
	}
}

func (s *PostgreSqlVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePostgreSqlVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePostgreSqlVersion(input string) (*PostgreSqlVersion, error) {
	vals := map[string]PostgreSqlVersion{
		"18": PostgreSqlVersionOneEight,
		"17": PostgreSqlVersionOneSeven,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PostgreSqlVersion(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateInProgress   ProvisioningState = "InProgress"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateProvisioning),
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
		"canceled":     ProvisioningStateCanceled,
		"failed":       ProvisioningStateFailed,
		"inprogress":   ProvisioningStateInProgress,
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicNetworkAccessState string

const (
	PublicNetworkAccessStateDisabled PublicNetworkAccessState = "Disabled"
	PublicNetworkAccessStateEnabled  PublicNetworkAccessState = "Enabled"
)

func PossibleValuesForPublicNetworkAccessState() []string {
	return []string{
		string(PublicNetworkAccessStateDisabled),
		string(PublicNetworkAccessStateEnabled),
	}
}

func (s *PublicNetworkAccessState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccessState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccessState(input string) (*PublicNetworkAccessState, error) {
	vals := map[string]PublicNetworkAccessState{
		"disabled": PublicNetworkAccessStateDisabled,
		"enabled":  PublicNetworkAccessStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccessState(input)
	return &out, nil
}
