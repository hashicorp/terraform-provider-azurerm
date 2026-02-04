package virtualnetworkpeerings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type SyncRemoteAddressSpace string

const (
	SyncRemoteAddressSpaceTrue SyncRemoteAddressSpace = "true"
)

func PossibleValuesForSyncRemoteAddressSpace() []string {
	return []string{
		string(SyncRemoteAddressSpaceTrue),
	}
}

func (s *SyncRemoteAddressSpace) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSyncRemoteAddressSpace(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSyncRemoteAddressSpace(input string) (*SyncRemoteAddressSpace, error) {
	vals := map[string]SyncRemoteAddressSpace{
		"true": SyncRemoteAddressSpaceTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SyncRemoteAddressSpace(input)
	return &out, nil
}

type VirtualNetworkEncryptionEnforcement string

const (
	VirtualNetworkEncryptionEnforcementAllowUnencrypted VirtualNetworkEncryptionEnforcement = "AllowUnencrypted"
	VirtualNetworkEncryptionEnforcementDropUnencrypted  VirtualNetworkEncryptionEnforcement = "DropUnencrypted"
)

func PossibleValuesForVirtualNetworkEncryptionEnforcement() []string {
	return []string{
		string(VirtualNetworkEncryptionEnforcementAllowUnencrypted),
		string(VirtualNetworkEncryptionEnforcementDropUnencrypted),
	}
}

func (s *VirtualNetworkEncryptionEnforcement) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkEncryptionEnforcement(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkEncryptionEnforcement(input string) (*VirtualNetworkEncryptionEnforcement, error) {
	vals := map[string]VirtualNetworkEncryptionEnforcement{
		"allowunencrypted": VirtualNetworkEncryptionEnforcementAllowUnencrypted,
		"dropunencrypted":  VirtualNetworkEncryptionEnforcementDropUnencrypted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkEncryptionEnforcement(input)
	return &out, nil
}

type VirtualNetworkPeeringLevel string

const (
	VirtualNetworkPeeringLevelFullyInSync             VirtualNetworkPeeringLevel = "FullyInSync"
	VirtualNetworkPeeringLevelLocalAndRemoteNotInSync VirtualNetworkPeeringLevel = "LocalAndRemoteNotInSync"
	VirtualNetworkPeeringLevelLocalNotInSync          VirtualNetworkPeeringLevel = "LocalNotInSync"
	VirtualNetworkPeeringLevelRemoteNotInSync         VirtualNetworkPeeringLevel = "RemoteNotInSync"
)

func PossibleValuesForVirtualNetworkPeeringLevel() []string {
	return []string{
		string(VirtualNetworkPeeringLevelFullyInSync),
		string(VirtualNetworkPeeringLevelLocalAndRemoteNotInSync),
		string(VirtualNetworkPeeringLevelLocalNotInSync),
		string(VirtualNetworkPeeringLevelRemoteNotInSync),
	}
}

func (s *VirtualNetworkPeeringLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkPeeringLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkPeeringLevel(input string) (*VirtualNetworkPeeringLevel, error) {
	vals := map[string]VirtualNetworkPeeringLevel{
		"fullyinsync":             VirtualNetworkPeeringLevelFullyInSync,
		"localandremotenotinsync": VirtualNetworkPeeringLevelLocalAndRemoteNotInSync,
		"localnotinsync":          VirtualNetworkPeeringLevelLocalNotInSync,
		"remotenotinsync":         VirtualNetworkPeeringLevelRemoteNotInSync,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkPeeringLevel(input)
	return &out, nil
}

type VirtualNetworkPeeringState string

const (
	VirtualNetworkPeeringStateConnected    VirtualNetworkPeeringState = "Connected"
	VirtualNetworkPeeringStateDisconnected VirtualNetworkPeeringState = "Disconnected"
	VirtualNetworkPeeringStateInitiated    VirtualNetworkPeeringState = "Initiated"
)

func PossibleValuesForVirtualNetworkPeeringState() []string {
	return []string{
		string(VirtualNetworkPeeringStateConnected),
		string(VirtualNetworkPeeringStateDisconnected),
		string(VirtualNetworkPeeringStateInitiated),
	}
}

func (s *VirtualNetworkPeeringState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualNetworkPeeringState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualNetworkPeeringState(input string) (*VirtualNetworkPeeringState, error) {
	vals := map[string]VirtualNetworkPeeringState{
		"connected":    VirtualNetworkPeeringStateConnected,
		"disconnected": VirtualNetworkPeeringStateDisconnected,
		"initiated":    VirtualNetworkPeeringStateInitiated,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualNetworkPeeringState(input)
	return &out, nil
}
