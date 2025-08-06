package imagedefinitions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoImageBuildStatus string

const (
	AutoImageBuildStatusDisabled AutoImageBuildStatus = "Disabled"
	AutoImageBuildStatusEnabled  AutoImageBuildStatus = "Enabled"
)

func PossibleValuesForAutoImageBuildStatus() []string {
	return []string{
		string(AutoImageBuildStatusDisabled),
		string(AutoImageBuildStatusEnabled),
	}
}

func (s *AutoImageBuildStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoImageBuildStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoImageBuildStatus(input string) (*AutoImageBuildStatus, error) {
	vals := map[string]AutoImageBuildStatus{
		"disabled": AutoImageBuildStatusDisabled,
		"enabled":  AutoImageBuildStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoImageBuildStatus(input)
	return &out, nil
}

type CatalogResourceValidationStatus string

const (
	CatalogResourceValidationStatusFailed    CatalogResourceValidationStatus = "Failed"
	CatalogResourceValidationStatusPending   CatalogResourceValidationStatus = "Pending"
	CatalogResourceValidationStatusSucceeded CatalogResourceValidationStatus = "Succeeded"
	CatalogResourceValidationStatusUnknown   CatalogResourceValidationStatus = "Unknown"
)

func PossibleValuesForCatalogResourceValidationStatus() []string {
	return []string{
		string(CatalogResourceValidationStatusFailed),
		string(CatalogResourceValidationStatusPending),
		string(CatalogResourceValidationStatusSucceeded),
		string(CatalogResourceValidationStatusUnknown),
	}
}

func (s *CatalogResourceValidationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogResourceValidationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogResourceValidationStatus(input string) (*CatalogResourceValidationStatus, error) {
	vals := map[string]CatalogResourceValidationStatus{
		"failed":    CatalogResourceValidationStatusFailed,
		"pending":   CatalogResourceValidationStatusPending,
		"succeeded": CatalogResourceValidationStatusSucceeded,
		"unknown":   CatalogResourceValidationStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogResourceValidationStatus(input)
	return &out, nil
}

type ImageDefinitionBuildStatus string

const (
	ImageDefinitionBuildStatusCancelled        ImageDefinitionBuildStatus = "Cancelled"
	ImageDefinitionBuildStatusFailed           ImageDefinitionBuildStatus = "Failed"
	ImageDefinitionBuildStatusRunning          ImageDefinitionBuildStatus = "Running"
	ImageDefinitionBuildStatusSucceeded        ImageDefinitionBuildStatus = "Succeeded"
	ImageDefinitionBuildStatusTimedOut         ImageDefinitionBuildStatus = "TimedOut"
	ImageDefinitionBuildStatusValidationFailed ImageDefinitionBuildStatus = "ValidationFailed"
)

func PossibleValuesForImageDefinitionBuildStatus() []string {
	return []string{
		string(ImageDefinitionBuildStatusCancelled),
		string(ImageDefinitionBuildStatusFailed),
		string(ImageDefinitionBuildStatusRunning),
		string(ImageDefinitionBuildStatusSucceeded),
		string(ImageDefinitionBuildStatusTimedOut),
		string(ImageDefinitionBuildStatusValidationFailed),
	}
}

func (s *ImageDefinitionBuildStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImageDefinitionBuildStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImageDefinitionBuildStatus(input string) (*ImageDefinitionBuildStatus, error) {
	vals := map[string]ImageDefinitionBuildStatus{
		"cancelled":        ImageDefinitionBuildStatusCancelled,
		"failed":           ImageDefinitionBuildStatusFailed,
		"running":          ImageDefinitionBuildStatusRunning,
		"succeeded":        ImageDefinitionBuildStatusSucceeded,
		"timedout":         ImageDefinitionBuildStatusTimedOut,
		"validationfailed": ImageDefinitionBuildStatusValidationFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImageDefinitionBuildStatus(input)
	return &out, nil
}

type ImageValidationStatus string

const (
	ImageValidationStatusFailed    ImageValidationStatus = "Failed"
	ImageValidationStatusPending   ImageValidationStatus = "Pending"
	ImageValidationStatusSucceeded ImageValidationStatus = "Succeeded"
	ImageValidationStatusTimedOut  ImageValidationStatus = "TimedOut"
	ImageValidationStatusUnknown   ImageValidationStatus = "Unknown"
)

func PossibleValuesForImageValidationStatus() []string {
	return []string{
		string(ImageValidationStatusFailed),
		string(ImageValidationStatusPending),
		string(ImageValidationStatusSucceeded),
		string(ImageValidationStatusTimedOut),
		string(ImageValidationStatusUnknown),
	}
}

func (s *ImageValidationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImageValidationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImageValidationStatus(input string) (*ImageValidationStatus, error) {
	vals := map[string]ImageValidationStatus{
		"failed":    ImageValidationStatusFailed,
		"pending":   ImageValidationStatusPending,
		"succeeded": ImageValidationStatusSucceeded,
		"timedout":  ImageValidationStatusTimedOut,
		"unknown":   ImageValidationStatusUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImageValidationStatus(input)
	return &out, nil
}
