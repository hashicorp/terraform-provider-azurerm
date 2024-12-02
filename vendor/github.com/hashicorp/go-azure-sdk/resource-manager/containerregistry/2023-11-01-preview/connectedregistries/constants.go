package connectedregistries

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActivationStatus string

const (
	ActivationStatusActive   ActivationStatus = "Active"
	ActivationStatusInactive ActivationStatus = "Inactive"
)

func PossibleValuesForActivationStatus() []string {
	return []string{
		string(ActivationStatusActive),
		string(ActivationStatusInactive),
	}
}

func (s *ActivationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActivationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActivationStatus(input string) (*ActivationStatus, error) {
	vals := map[string]ActivationStatus{
		"active":   ActivationStatusActive,
		"inactive": ActivationStatusInactive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActivationStatus(input)
	return &out, nil
}

type AuditLogStatus string

const (
	AuditLogStatusDisabled AuditLogStatus = "Disabled"
	AuditLogStatusEnabled  AuditLogStatus = "Enabled"
)

func PossibleValuesForAuditLogStatus() []string {
	return []string{
		string(AuditLogStatusDisabled),
		string(AuditLogStatusEnabled),
	}
}

func (s *AuditLogStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuditLogStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuditLogStatus(input string) (*AuditLogStatus, error) {
	vals := map[string]AuditLogStatus{
		"disabled": AuditLogStatusDisabled,
		"enabled":  AuditLogStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuditLogStatus(input)
	return &out, nil
}

type CertificateType string

const (
	CertificateTypeLocalDirectory CertificateType = "LocalDirectory"
)

func PossibleValuesForCertificateType() []string {
	return []string{
		string(CertificateTypeLocalDirectory),
	}
}

func (s *CertificateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateType(input string) (*CertificateType, error) {
	vals := map[string]CertificateType{
		"localdirectory": CertificateTypeLocalDirectory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateType(input)
	return &out, nil
}

type ConnectedRegistryMode string

const (
	ConnectedRegistryModeMirror    ConnectedRegistryMode = "Mirror"
	ConnectedRegistryModeReadOnly  ConnectedRegistryMode = "ReadOnly"
	ConnectedRegistryModeReadWrite ConnectedRegistryMode = "ReadWrite"
	ConnectedRegistryModeRegistry  ConnectedRegistryMode = "Registry"
)

func PossibleValuesForConnectedRegistryMode() []string {
	return []string{
		string(ConnectedRegistryModeMirror),
		string(ConnectedRegistryModeReadOnly),
		string(ConnectedRegistryModeReadWrite),
		string(ConnectedRegistryModeRegistry),
	}
}

func (s *ConnectedRegistryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectedRegistryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectedRegistryMode(input string) (*ConnectedRegistryMode, error) {
	vals := map[string]ConnectedRegistryMode{
		"mirror":    ConnectedRegistryModeMirror,
		"readonly":  ConnectedRegistryModeReadOnly,
		"readwrite": ConnectedRegistryModeReadWrite,
		"registry":  ConnectedRegistryModeRegistry,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectedRegistryMode(input)
	return &out, nil
}

type ConnectionState string

const (
	ConnectionStateOffline   ConnectionState = "Offline"
	ConnectionStateOnline    ConnectionState = "Online"
	ConnectionStateSyncing   ConnectionState = "Syncing"
	ConnectionStateUnhealthy ConnectionState = "Unhealthy"
)

func PossibleValuesForConnectionState() []string {
	return []string{
		string(ConnectionStateOffline),
		string(ConnectionStateOnline),
		string(ConnectionStateSyncing),
		string(ConnectionStateUnhealthy),
	}
}

func (s *ConnectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionState(input string) (*ConnectionState, error) {
	vals := map[string]ConnectionState{
		"offline":   ConnectionStateOffline,
		"online":    ConnectionStateOnline,
		"syncing":   ConnectionStateSyncing,
		"unhealthy": ConnectionStateUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionState(input)
	return &out, nil
}

type LogLevel string

const (
	LogLevelDebug       LogLevel = "Debug"
	LogLevelError       LogLevel = "Error"
	LogLevelInformation LogLevel = "Information"
	LogLevelNone        LogLevel = "None"
	LogLevelWarning     LogLevel = "Warning"
)

func PossibleValuesForLogLevel() []string {
	return []string{
		string(LogLevelDebug),
		string(LogLevelError),
		string(LogLevelInformation),
		string(LogLevelNone),
		string(LogLevelWarning),
	}
}

func (s *LogLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLogLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLogLevel(input string) (*LogLevel, error) {
	vals := map[string]LogLevel{
		"debug":       LogLevelDebug,
		"error":       LogLevelError,
		"information": LogLevelInformation,
		"none":        LogLevelNone,
		"warning":     LogLevelWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LogLevel(input)
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

type TlsStatus string

const (
	TlsStatusDisabled TlsStatus = "Disabled"
	TlsStatusEnabled  TlsStatus = "Enabled"
)

func PossibleValuesForTlsStatus() []string {
	return []string{
		string(TlsStatusDisabled),
		string(TlsStatusEnabled),
	}
}

func (s *TlsStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsStatus(input string) (*TlsStatus, error) {
	vals := map[string]TlsStatus{
		"disabled": TlsStatusDisabled,
		"enabled":  TlsStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsStatus(input)
	return &out, nil
}
