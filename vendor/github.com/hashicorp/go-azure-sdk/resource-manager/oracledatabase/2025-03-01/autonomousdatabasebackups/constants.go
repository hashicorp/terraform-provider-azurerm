package autonomousdatabasebackups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutonomousDatabaseBackupLifecycleState string

const (
	AutonomousDatabaseBackupLifecycleStateActive   AutonomousDatabaseBackupLifecycleState = "Active"
	AutonomousDatabaseBackupLifecycleStateCreating AutonomousDatabaseBackupLifecycleState = "Creating"
	AutonomousDatabaseBackupLifecycleStateDeleting AutonomousDatabaseBackupLifecycleState = "Deleting"
	AutonomousDatabaseBackupLifecycleStateFailed   AutonomousDatabaseBackupLifecycleState = "Failed"
	AutonomousDatabaseBackupLifecycleStateUpdating AutonomousDatabaseBackupLifecycleState = "Updating"
)

func PossibleValuesForAutonomousDatabaseBackupLifecycleState() []string {
	return []string{
		string(AutonomousDatabaseBackupLifecycleStateActive),
		string(AutonomousDatabaseBackupLifecycleStateCreating),
		string(AutonomousDatabaseBackupLifecycleStateDeleting),
		string(AutonomousDatabaseBackupLifecycleStateFailed),
		string(AutonomousDatabaseBackupLifecycleStateUpdating),
	}
}

func (s *AutonomousDatabaseBackupLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutonomousDatabaseBackupLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutonomousDatabaseBackupLifecycleState(input string) (*AutonomousDatabaseBackupLifecycleState, error) {
	vals := map[string]AutonomousDatabaseBackupLifecycleState{
		"active":   AutonomousDatabaseBackupLifecycleStateActive,
		"creating": AutonomousDatabaseBackupLifecycleStateCreating,
		"deleting": AutonomousDatabaseBackupLifecycleStateDeleting,
		"failed":   AutonomousDatabaseBackupLifecycleStateFailed,
		"updating": AutonomousDatabaseBackupLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutonomousDatabaseBackupLifecycleState(input)
	return &out, nil
}

type AutonomousDatabaseBackupType string

const (
	AutonomousDatabaseBackupTypeFull        AutonomousDatabaseBackupType = "Full"
	AutonomousDatabaseBackupTypeIncremental AutonomousDatabaseBackupType = "Incremental"
	AutonomousDatabaseBackupTypeLongTerm    AutonomousDatabaseBackupType = "LongTerm"
)

func PossibleValuesForAutonomousDatabaseBackupType() []string {
	return []string{
		string(AutonomousDatabaseBackupTypeFull),
		string(AutonomousDatabaseBackupTypeIncremental),
		string(AutonomousDatabaseBackupTypeLongTerm),
	}
}

func (s *AutonomousDatabaseBackupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutonomousDatabaseBackupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutonomousDatabaseBackupType(input string) (*AutonomousDatabaseBackupType, error) {
	vals := map[string]AutonomousDatabaseBackupType{
		"full":        AutonomousDatabaseBackupTypeFull,
		"incremental": AutonomousDatabaseBackupTypeIncremental,
		"longterm":    AutonomousDatabaseBackupTypeLongTerm,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutonomousDatabaseBackupType(input)
	return &out, nil
}

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}
