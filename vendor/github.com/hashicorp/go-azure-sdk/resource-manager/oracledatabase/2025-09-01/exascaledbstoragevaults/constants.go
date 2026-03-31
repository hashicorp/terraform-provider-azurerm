package exascaledbstoragevaults

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

type ExascaleDbStorageVaultLifecycleState string

const (
	ExascaleDbStorageVaultLifecycleStateAvailable    ExascaleDbStorageVaultLifecycleState = "Available"
	ExascaleDbStorageVaultLifecycleStateFailed       ExascaleDbStorageVaultLifecycleState = "Failed"
	ExascaleDbStorageVaultLifecycleStateProvisioning ExascaleDbStorageVaultLifecycleState = "Provisioning"
	ExascaleDbStorageVaultLifecycleStateTerminated   ExascaleDbStorageVaultLifecycleState = "Terminated"
	ExascaleDbStorageVaultLifecycleStateTerminating  ExascaleDbStorageVaultLifecycleState = "Terminating"
	ExascaleDbStorageVaultLifecycleStateUpdating     ExascaleDbStorageVaultLifecycleState = "Updating"
)

func PossibleValuesForExascaleDbStorageVaultLifecycleState() []string {
	return []string{
		string(ExascaleDbStorageVaultLifecycleStateAvailable),
		string(ExascaleDbStorageVaultLifecycleStateFailed),
		string(ExascaleDbStorageVaultLifecycleStateProvisioning),
		string(ExascaleDbStorageVaultLifecycleStateTerminated),
		string(ExascaleDbStorageVaultLifecycleStateTerminating),
		string(ExascaleDbStorageVaultLifecycleStateUpdating),
	}
}

func (s *ExascaleDbStorageVaultLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExascaleDbStorageVaultLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExascaleDbStorageVaultLifecycleState(input string) (*ExascaleDbStorageVaultLifecycleState, error) {
	vals := map[string]ExascaleDbStorageVaultLifecycleState{
		"available":    ExascaleDbStorageVaultLifecycleStateAvailable,
		"failed":       ExascaleDbStorageVaultLifecycleStateFailed,
		"provisioning": ExascaleDbStorageVaultLifecycleStateProvisioning,
		"terminated":   ExascaleDbStorageVaultLifecycleStateTerminated,
		"terminating":  ExascaleDbStorageVaultLifecycleStateTerminating,
		"updating":     ExascaleDbStorageVaultLifecycleStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExascaleDbStorageVaultLifecycleState(input)
	return &out, nil
}

type ShapeAttribute string

const (
	ShapeAttributeBLOCKSTORAGE ShapeAttribute = "BLOCK_STORAGE"
	ShapeAttributeSMARTSTORAGE ShapeAttribute = "SMART_STORAGE"
)

func PossibleValuesForShapeAttribute() []string {
	return []string{
		string(ShapeAttributeBLOCKSTORAGE),
		string(ShapeAttributeSMARTSTORAGE),
	}
}

func (s *ShapeAttribute) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseShapeAttribute(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseShapeAttribute(input string) (*ShapeAttribute, error) {
	vals := map[string]ShapeAttribute{
		"block_storage": ShapeAttributeBLOCKSTORAGE,
		"smart_storage": ShapeAttributeSMARTSTORAGE,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ShapeAttribute(input)
	return &out, nil
}
