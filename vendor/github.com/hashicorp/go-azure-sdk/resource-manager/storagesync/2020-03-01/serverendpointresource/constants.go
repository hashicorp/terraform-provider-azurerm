package serverendpointresource

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type ServerEndpointCloudTieringHealthState string

const (
	ServerEndpointCloudTieringHealthStateError   ServerEndpointCloudTieringHealthState = "Error"
	ServerEndpointCloudTieringHealthStateHealthy ServerEndpointCloudTieringHealthState = "Healthy"
)

func PossibleValuesForServerEndpointCloudTieringHealthState() []string {
	return []string{
		string(ServerEndpointCloudTieringHealthStateError),
		string(ServerEndpointCloudTieringHealthStateHealthy),
	}
}

func (s *ServerEndpointCloudTieringHealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerEndpointCloudTieringHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerEndpointCloudTieringHealthState(input string) (*ServerEndpointCloudTieringHealthState, error) {
	vals := map[string]ServerEndpointCloudTieringHealthState{
		"error":   ServerEndpointCloudTieringHealthStateError,
		"healthy": ServerEndpointCloudTieringHealthStateHealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerEndpointCloudTieringHealthState(input)
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

type ServerEndpointSyncHealthState string

const (
	ServerEndpointSyncHealthStateError                                    ServerEndpointSyncHealthState = "Error"
	ServerEndpointSyncHealthStateHealthy                                  ServerEndpointSyncHealthState = "Healthy"
	ServerEndpointSyncHealthStateNoActivity                               ServerEndpointSyncHealthState = "NoActivity"
	ServerEndpointSyncHealthStateSyncBlockedForChangeDetectionPostRestore ServerEndpointSyncHealthState = "SyncBlockedForChangeDetectionPostRestore"
	ServerEndpointSyncHealthStateSyncBlockedForRestore                    ServerEndpointSyncHealthState = "SyncBlockedForRestore"
)

func PossibleValuesForServerEndpointSyncHealthState() []string {
	return []string{
		string(ServerEndpointSyncHealthStateError),
		string(ServerEndpointSyncHealthStateHealthy),
		string(ServerEndpointSyncHealthStateNoActivity),
		string(ServerEndpointSyncHealthStateSyncBlockedForChangeDetectionPostRestore),
		string(ServerEndpointSyncHealthStateSyncBlockedForRestore),
	}
}

func (s *ServerEndpointSyncHealthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerEndpointSyncHealthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerEndpointSyncHealthState(input string) (*ServerEndpointSyncHealthState, error) {
	vals := map[string]ServerEndpointSyncHealthState{
		"error":      ServerEndpointSyncHealthStateError,
		"healthy":    ServerEndpointSyncHealthStateHealthy,
		"noactivity": ServerEndpointSyncHealthStateNoActivity,
		"syncblockedforchangedetectionpostrestore": ServerEndpointSyncHealthStateSyncBlockedForChangeDetectionPostRestore,
		"syncblockedforrestore":                    ServerEndpointSyncHealthStateSyncBlockedForRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerEndpointSyncHealthState(input)
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
