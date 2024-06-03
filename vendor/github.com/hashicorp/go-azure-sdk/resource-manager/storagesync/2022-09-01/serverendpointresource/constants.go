package serverendpointresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudTieringLowDiskModeState string

const (
	CloudTieringLowDiskModeStateDisabled CloudTieringLowDiskModeState = "Disabled"
	CloudTieringLowDiskModeStateEnabled  CloudTieringLowDiskModeState = "Enabled"
)

func PossibleValuesForCloudTieringLowDiskModeState() []string {
	return []string{
		string(CloudTieringLowDiskModeStateDisabled),
		string(CloudTieringLowDiskModeStateEnabled),
	}
}

func (s *CloudTieringLowDiskModeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCloudTieringLowDiskModeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCloudTieringLowDiskModeState(input string) (*CloudTieringLowDiskModeState, error) {
	vals := map[string]CloudTieringLowDiskModeState{
		"disabled": CloudTieringLowDiskModeStateDisabled,
		"enabled":  CloudTieringLowDiskModeStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CloudTieringLowDiskModeState(input)
	return &out, nil
}

type FeatureStatus string

const (
	FeatureStatusOff FeatureStatus = "off"
	FeatureStatusOn  FeatureStatus = "on"
)

func PossibleValuesForFeatureStatus() []string {
	return []string{
		string(FeatureStatusOff),
		string(FeatureStatusOn),
	}
}

func (s *FeatureStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFeatureStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFeatureStatus(input string) (*FeatureStatus, error) {
	vals := map[string]FeatureStatus{
		"off": FeatureStatusOff,
		"on":  FeatureStatusOn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FeatureStatus(input)
	return &out, nil
}

type InitialDownloadPolicy string

const (
	InitialDownloadPolicyAvoidTieredFiles           InitialDownloadPolicy = "AvoidTieredFiles"
	InitialDownloadPolicyNamespaceOnly              InitialDownloadPolicy = "NamespaceOnly"
	InitialDownloadPolicyNamespaceThenModifiedFiles InitialDownloadPolicy = "NamespaceThenModifiedFiles"
)

func PossibleValuesForInitialDownloadPolicy() []string {
	return []string{
		string(InitialDownloadPolicyAvoidTieredFiles),
		string(InitialDownloadPolicyNamespaceOnly),
		string(InitialDownloadPolicyNamespaceThenModifiedFiles),
	}
}

func (s *InitialDownloadPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInitialDownloadPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInitialDownloadPolicy(input string) (*InitialDownloadPolicy, error) {
	vals := map[string]InitialDownloadPolicy{
		"avoidtieredfiles":           InitialDownloadPolicyAvoidTieredFiles,
		"namespaceonly":              InitialDownloadPolicyNamespaceOnly,
		"namespacethenmodifiedfiles": InitialDownloadPolicyNamespaceThenModifiedFiles,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InitialDownloadPolicy(input)
	return &out, nil
}

type InitialUploadPolicy string

const (
	InitialUploadPolicyMerge               InitialUploadPolicy = "Merge"
	InitialUploadPolicyServerAuthoritative InitialUploadPolicy = "ServerAuthoritative"
)

func PossibleValuesForInitialUploadPolicy() []string {
	return []string{
		string(InitialUploadPolicyMerge),
		string(InitialUploadPolicyServerAuthoritative),
	}
}

func (s *InitialUploadPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInitialUploadPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInitialUploadPolicy(input string) (*InitialUploadPolicy, error) {
	vals := map[string]InitialUploadPolicy{
		"merge":               InitialUploadPolicyMerge,
		"serverauthoritative": InitialUploadPolicyServerAuthoritative,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InitialUploadPolicy(input)
	return &out, nil
}

type LocalCacheMode string

const (
	LocalCacheModeDownloadNewAndModifiedFiles LocalCacheMode = "DownloadNewAndModifiedFiles"
	LocalCacheModeUpdateLocallyCachedFiles    LocalCacheMode = "UpdateLocallyCachedFiles"
)

func PossibleValuesForLocalCacheMode() []string {
	return []string{
		string(LocalCacheModeDownloadNewAndModifiedFiles),
		string(LocalCacheModeUpdateLocallyCachedFiles),
	}
}

func (s *LocalCacheMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalCacheMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalCacheMode(input string) (*LocalCacheMode, error) {
	vals := map[string]LocalCacheMode{
		"downloadnewandmodifiedfiles": LocalCacheModeDownloadNewAndModifiedFiles,
		"updatelocallycachedfiles":    LocalCacheModeUpdateLocallyCachedFiles,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalCacheMode(input)
	return &out, nil
}

type ServerEndpointHealthState string

const (
	ServerEndpointHealthStateError       ServerEndpointHealthState = "Error"
	ServerEndpointHealthStateHealthy     ServerEndpointHealthState = "Healthy"
	ServerEndpointHealthStateUnavailable ServerEndpointHealthState = "Unavailable"
)

func PossibleValuesForServerEndpointHealthState() []string {
	return []string{
		string(ServerEndpointHealthStateError),
		string(ServerEndpointHealthStateHealthy),
		string(ServerEndpointHealthStateUnavailable),
	}
}

func (s *ServerEndpointHealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerEndpointHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerEndpointHealthState(input string) (*ServerEndpointHealthState, error) {
	vals := map[string]ServerEndpointHealthState{
		"error":       ServerEndpointHealthStateError,
		"healthy":     ServerEndpointHealthStateHealthy,
		"unavailable": ServerEndpointHealthStateUnavailable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerEndpointHealthState(input)
	return &out, nil
}

type ServerEndpointOfflineDataTransferState string

const (
	ServerEndpointOfflineDataTransferStateComplete   ServerEndpointOfflineDataTransferState = "Complete"
	ServerEndpointOfflineDataTransferStateInProgress ServerEndpointOfflineDataTransferState = "InProgress"
	ServerEndpointOfflineDataTransferStateNotRunning ServerEndpointOfflineDataTransferState = "NotRunning"
	ServerEndpointOfflineDataTransferStateStopping   ServerEndpointOfflineDataTransferState = "Stopping"
)

func PossibleValuesForServerEndpointOfflineDataTransferState() []string {
	return []string{
		string(ServerEndpointOfflineDataTransferStateComplete),
		string(ServerEndpointOfflineDataTransferStateInProgress),
		string(ServerEndpointOfflineDataTransferStateNotRunning),
		string(ServerEndpointOfflineDataTransferStateStopping),
	}
}

func (s *ServerEndpointOfflineDataTransferState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerEndpointOfflineDataTransferState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerEndpointOfflineDataTransferState(input string) (*ServerEndpointOfflineDataTransferState, error) {
	vals := map[string]ServerEndpointOfflineDataTransferState{
		"complete":   ServerEndpointOfflineDataTransferStateComplete,
		"inprogress": ServerEndpointOfflineDataTransferStateInProgress,
		"notrunning": ServerEndpointOfflineDataTransferStateNotRunning,
		"stopping":   ServerEndpointOfflineDataTransferStateStopping,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerEndpointOfflineDataTransferState(input)
	return &out, nil
}

type ServerEndpointSyncActivityState string

const (
	ServerEndpointSyncActivityStateDownload          ServerEndpointSyncActivityState = "Download"
	ServerEndpointSyncActivityStateUpload            ServerEndpointSyncActivityState = "Upload"
	ServerEndpointSyncActivityStateUploadAndDownload ServerEndpointSyncActivityState = "UploadAndDownload"
)

func PossibleValuesForServerEndpointSyncActivityState() []string {
	return []string{
		string(ServerEndpointSyncActivityStateDownload),
		string(ServerEndpointSyncActivityStateUpload),
		string(ServerEndpointSyncActivityStateUploadAndDownload),
	}
}

func (s *ServerEndpointSyncActivityState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerEndpointSyncActivityState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerEndpointSyncActivityState(input string) (*ServerEndpointSyncActivityState, error) {
	vals := map[string]ServerEndpointSyncActivityState{
		"download":          ServerEndpointSyncActivityStateDownload,
		"upload":            ServerEndpointSyncActivityStateUpload,
		"uploadanddownload": ServerEndpointSyncActivityStateUploadAndDownload,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerEndpointSyncActivityState(input)
	return &out, nil
}

type ServerEndpointSyncMode string

const (
	ServerEndpointSyncModeInitialFullDownload ServerEndpointSyncMode = "InitialFullDownload"
	ServerEndpointSyncModeInitialUpload       ServerEndpointSyncMode = "InitialUpload"
	ServerEndpointSyncModeNamespaceDownload   ServerEndpointSyncMode = "NamespaceDownload"
	ServerEndpointSyncModeRegular             ServerEndpointSyncMode = "Regular"
	ServerEndpointSyncModeSnapshotUpload      ServerEndpointSyncMode = "SnapshotUpload"
)

func PossibleValuesForServerEndpointSyncMode() []string {
	return []string{
		string(ServerEndpointSyncModeInitialFullDownload),
		string(ServerEndpointSyncModeInitialUpload),
		string(ServerEndpointSyncModeNamespaceDownload),
		string(ServerEndpointSyncModeRegular),
		string(ServerEndpointSyncModeSnapshotUpload),
	}
}

func (s *ServerEndpointSyncMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerEndpointSyncMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerEndpointSyncMode(input string) (*ServerEndpointSyncMode, error) {
	vals := map[string]ServerEndpointSyncMode{
		"initialfulldownload": ServerEndpointSyncModeInitialFullDownload,
		"initialupload":       ServerEndpointSyncModeInitialUpload,
		"namespacedownload":   ServerEndpointSyncModeNamespaceDownload,
		"regular":             ServerEndpointSyncModeRegular,
		"snapshotupload":      ServerEndpointSyncModeSnapshotUpload,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerEndpointSyncMode(input)
	return &out, nil
}

type ServerProvisioningStatus string

const (
	ServerProvisioningStatusError                  ServerProvisioningStatus = "Error"
	ServerProvisioningStatusInProgress             ServerProvisioningStatus = "InProgress"
	ServerProvisioningStatusNotStarted             ServerProvisioningStatus = "NotStarted"
	ServerProvisioningStatusReadySyncFunctional    ServerProvisioningStatus = "Ready_SyncFunctional"
	ServerProvisioningStatusReadySyncNotFunctional ServerProvisioningStatus = "Ready_SyncNotFunctional"
)

func PossibleValuesForServerProvisioningStatus() []string {
	return []string{
		string(ServerProvisioningStatusError),
		string(ServerProvisioningStatusInProgress),
		string(ServerProvisioningStatusNotStarted),
		string(ServerProvisioningStatusReadySyncFunctional),
		string(ServerProvisioningStatusReadySyncNotFunctional),
	}
}

func (s *ServerProvisioningStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerProvisioningStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerProvisioningStatus(input string) (*ServerProvisioningStatus, error) {
	vals := map[string]ServerProvisioningStatus{
		"error":                   ServerProvisioningStatusError,
		"inprogress":              ServerProvisioningStatusInProgress,
		"notstarted":              ServerProvisioningStatusNotStarted,
		"ready_syncfunctional":    ServerProvisioningStatusReadySyncFunctional,
		"ready_syncnotfunctional": ServerProvisioningStatusReadySyncNotFunctional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerProvisioningStatus(input)
	return &out, nil
}
