package edgeactionversions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EdgeActionIsDefaultVersion string

const (
	EdgeActionIsDefaultVersionFalse EdgeActionIsDefaultVersion = "False"
	EdgeActionIsDefaultVersionTrue  EdgeActionIsDefaultVersion = "True"
)

func PossibleValuesForEdgeActionIsDefaultVersion() []string {
	return []string{
		string(EdgeActionIsDefaultVersionFalse),
		string(EdgeActionIsDefaultVersionTrue),
	}
}

func (s *EdgeActionIsDefaultVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEdgeActionIsDefaultVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEdgeActionIsDefaultVersion(input string) (*EdgeActionIsDefaultVersion, error) {
	vals := map[string]EdgeActionIsDefaultVersion{
		"false": EdgeActionIsDefaultVersionFalse,
		"true":  EdgeActionIsDefaultVersionTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdgeActionIsDefaultVersion(input)
	return &out, nil
}

type EdgeActionVersionDeploymentType string

const (
	EdgeActionVersionDeploymentTypeFile   EdgeActionVersionDeploymentType = "file"
	EdgeActionVersionDeploymentTypeOthers EdgeActionVersionDeploymentType = "others"
	EdgeActionVersionDeploymentTypeZip    EdgeActionVersionDeploymentType = "zip"
)

func PossibleValuesForEdgeActionVersionDeploymentType() []string {
	return []string{
		string(EdgeActionVersionDeploymentTypeFile),
		string(EdgeActionVersionDeploymentTypeOthers),
		string(EdgeActionVersionDeploymentTypeZip),
	}
}

func (s *EdgeActionVersionDeploymentType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEdgeActionVersionDeploymentType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEdgeActionVersionDeploymentType(input string) (*EdgeActionVersionDeploymentType, error) {
	vals := map[string]EdgeActionVersionDeploymentType{
		"file":   EdgeActionVersionDeploymentTypeFile,
		"others": EdgeActionVersionDeploymentTypeOthers,
		"zip":    EdgeActionVersionDeploymentTypeZip,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdgeActionVersionDeploymentType(input)
	return &out, nil
}

type EdgeActionVersionValidationStatus string

const (
	EdgeActionVersionValidationStatusFailed    EdgeActionVersionValidationStatus = "Failed"
	EdgeActionVersionValidationStatusSucceeded EdgeActionVersionValidationStatus = "Succeeded"
)

func PossibleValuesForEdgeActionVersionValidationStatus() []string {
	return []string{
		string(EdgeActionVersionValidationStatusFailed),
		string(EdgeActionVersionValidationStatusSucceeded),
	}
}

func (s *EdgeActionVersionValidationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEdgeActionVersionValidationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEdgeActionVersionValidationStatus(input string) (*EdgeActionVersionValidationStatus, error) {
	vals := map[string]EdgeActionVersionValidationStatus{
		"failed":    EdgeActionVersionValidationStatusFailed,
		"succeeded": EdgeActionVersionValidationStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdgeActionVersionValidationStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpgrading    ProvisioningState = "Upgrading"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpgrading),
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
		"provisioning": ProvisioningStateProvisioning,
		"succeeded":    ProvisioningStateSucceeded,
		"upgrading":    ProvisioningStateUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
