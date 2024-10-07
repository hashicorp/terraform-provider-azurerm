package packetcaptures

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCaptureTargetType string

const (
	PacketCaptureTargetTypeAzureVM   PacketCaptureTargetType = "AzureVM"
	PacketCaptureTargetTypeAzureVMSS PacketCaptureTargetType = "AzureVMSS"
)

func PossibleValuesForPacketCaptureTargetType() []string {
	return []string{
		string(PacketCaptureTargetTypeAzureVM),
		string(PacketCaptureTargetTypeAzureVMSS),
	}
}

func (s *PacketCaptureTargetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePacketCaptureTargetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePacketCaptureTargetType(input string) (*PacketCaptureTargetType, error) {
	vals := map[string]PacketCaptureTargetType{
		"azurevm":   PacketCaptureTargetTypeAzureVM,
		"azurevmss": PacketCaptureTargetTypeAzureVMSS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PacketCaptureTargetType(input)
	return &out, nil
}

type PcError string

const (
	PcErrorAgentStopped    PcError = "AgentStopped"
	PcErrorCaptureFailed   PcError = "CaptureFailed"
	PcErrorInternalError   PcError = "InternalError"
	PcErrorLocalFileFailed PcError = "LocalFileFailed"
	PcErrorStorageFailed   PcError = "StorageFailed"
)

func PossibleValuesForPcError() []string {
	return []string{
		string(PcErrorAgentStopped),
		string(PcErrorCaptureFailed),
		string(PcErrorInternalError),
		string(PcErrorLocalFileFailed),
		string(PcErrorStorageFailed),
	}
}

func (s *PcError) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePcError(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePcError(input string) (*PcError, error) {
	vals := map[string]PcError{
		"agentstopped":    PcErrorAgentStopped,
		"capturefailed":   PcErrorCaptureFailed,
		"internalerror":   PcErrorInternalError,
		"localfilefailed": PcErrorLocalFileFailed,
		"storagefailed":   PcErrorStorageFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PcError(input)
	return &out, nil
}

type PcProtocol string

const (
	PcProtocolAny PcProtocol = "Any"
	PcProtocolTCP PcProtocol = "TCP"
	PcProtocolUDP PcProtocol = "UDP"
)

func PossibleValuesForPcProtocol() []string {
	return []string{
		string(PcProtocolAny),
		string(PcProtocolTCP),
		string(PcProtocolUDP),
	}
}

func (s *PcProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePcProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePcProtocol(input string) (*PcProtocol, error) {
	vals := map[string]PcProtocol{
		"any": PcProtocolAny,
		"tcp": PcProtocolTCP,
		"udp": PcProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PcProtocol(input)
	return &out, nil
}

type PcStatus string

const (
	PcStatusError      PcStatus = "Error"
	PcStatusNotStarted PcStatus = "NotStarted"
	PcStatusRunning    PcStatus = "Running"
	PcStatusStopped    PcStatus = "Stopped"
	PcStatusUnknown    PcStatus = "Unknown"
)

func PossibleValuesForPcStatus() []string {
	return []string{
		string(PcStatusError),
		string(PcStatusNotStarted),
		string(PcStatusRunning),
		string(PcStatusStopped),
		string(PcStatusUnknown),
	}
}

func (s *PcStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePcStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePcStatus(input string) (*PcStatus, error) {
	vals := map[string]PcStatus{
		"error":      PcStatusError,
		"notstarted": PcStatusNotStarted,
		"running":    PcStatusRunning,
		"stopped":    PcStatusStopped,
		"unknown":    PcStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PcStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
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
