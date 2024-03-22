package staticsites

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BuildStatus string

const (
	BuildStatusDeleting             BuildStatus = "Deleting"
	BuildStatusDeploying            BuildStatus = "Deploying"
	BuildStatusDetached             BuildStatus = "Detached"
	BuildStatusFailed               BuildStatus = "Failed"
	BuildStatusReady                BuildStatus = "Ready"
	BuildStatusUploading            BuildStatus = "Uploading"
	BuildStatusWaitingForDeployment BuildStatus = "WaitingForDeployment"
)

func PossibleValuesForBuildStatus() []string {
	return []string{
		string(BuildStatusDeleting),
		string(BuildStatusDeploying),
		string(BuildStatusDetached),
		string(BuildStatusFailed),
		string(BuildStatusReady),
		string(BuildStatusUploading),
		string(BuildStatusWaitingForDeployment),
	}
}

func (s *BuildStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBuildStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBuildStatus(input string) (*BuildStatus, error) {
	vals := map[string]BuildStatus{
		"deleting":             BuildStatusDeleting,
		"deploying":            BuildStatusDeploying,
		"detached":             BuildStatusDetached,
		"failed":               BuildStatusFailed,
		"ready":                BuildStatusReady,
		"uploading":            BuildStatusUploading,
		"waitingfordeployment": BuildStatusWaitingForDeployment,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BuildStatus(input)
	return &out, nil
}

type CustomDomainStatus string

const (
	CustomDomainStatusAdding                    CustomDomainStatus = "Adding"
	CustomDomainStatusDeleting                  CustomDomainStatus = "Deleting"
	CustomDomainStatusFailed                    CustomDomainStatus = "Failed"
	CustomDomainStatusReady                     CustomDomainStatus = "Ready"
	CustomDomainStatusRetrievingValidationToken CustomDomainStatus = "RetrievingValidationToken"
	CustomDomainStatusUnhealthy                 CustomDomainStatus = "Unhealthy"
	CustomDomainStatusValidating                CustomDomainStatus = "Validating"
)

func PossibleValuesForCustomDomainStatus() []string {
	return []string{
		string(CustomDomainStatusAdding),
		string(CustomDomainStatusDeleting),
		string(CustomDomainStatusFailed),
		string(CustomDomainStatusReady),
		string(CustomDomainStatusRetrievingValidationToken),
		string(CustomDomainStatusUnhealthy),
		string(CustomDomainStatusValidating),
	}
}

func (s *CustomDomainStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomDomainStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomDomainStatus(input string) (*CustomDomainStatus, error) {
	vals := map[string]CustomDomainStatus{
		"adding":                    CustomDomainStatusAdding,
		"deleting":                  CustomDomainStatusDeleting,
		"failed":                    CustomDomainStatusFailed,
		"ready":                     CustomDomainStatusReady,
		"retrievingvalidationtoken": CustomDomainStatusRetrievingValidationToken,
		"unhealthy":                 CustomDomainStatusUnhealthy,
		"validating":                CustomDomainStatusValidating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomDomainStatus(input)
	return &out, nil
}

type EnterpriseGradeCdnStatus string

const (
	EnterpriseGradeCdnStatusDisabled  EnterpriseGradeCdnStatus = "Disabled"
	EnterpriseGradeCdnStatusDisabling EnterpriseGradeCdnStatus = "Disabling"
	EnterpriseGradeCdnStatusEnabled   EnterpriseGradeCdnStatus = "Enabled"
	EnterpriseGradeCdnStatusEnabling  EnterpriseGradeCdnStatus = "Enabling"
)

func PossibleValuesForEnterpriseGradeCdnStatus() []string {
	return []string{
		string(EnterpriseGradeCdnStatusDisabled),
		string(EnterpriseGradeCdnStatusDisabling),
		string(EnterpriseGradeCdnStatusEnabled),
		string(EnterpriseGradeCdnStatusEnabling),
	}
}

func (s *EnterpriseGradeCdnStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnterpriseGradeCdnStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnterpriseGradeCdnStatus(input string) (*EnterpriseGradeCdnStatus, error) {
	vals := map[string]EnterpriseGradeCdnStatus{
		"disabled":  EnterpriseGradeCdnStatusDisabled,
		"disabling": EnterpriseGradeCdnStatusDisabling,
		"enabled":   EnterpriseGradeCdnStatusEnabled,
		"enabling":  EnterpriseGradeCdnStatusEnabling,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnterpriseGradeCdnStatus(input)
	return &out, nil
}

type StagingEnvironmentPolicy string

const (
	StagingEnvironmentPolicyDisabled StagingEnvironmentPolicy = "Disabled"
	StagingEnvironmentPolicyEnabled  StagingEnvironmentPolicy = "Enabled"
)

func PossibleValuesForStagingEnvironmentPolicy() []string {
	return []string{
		string(StagingEnvironmentPolicyDisabled),
		string(StagingEnvironmentPolicyEnabled),
	}
}

func (s *StagingEnvironmentPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStagingEnvironmentPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStagingEnvironmentPolicy(input string) (*StagingEnvironmentPolicy, error) {
	vals := map[string]StagingEnvironmentPolicy{
		"disabled": StagingEnvironmentPolicyDisabled,
		"enabled":  StagingEnvironmentPolicyEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StagingEnvironmentPolicy(input)
	return &out, nil
}

type TriggerTypes string

const (
	TriggerTypesHTTPTrigger TriggerTypes = "HttpTrigger"
	TriggerTypesUnknown     TriggerTypes = "Unknown"
)

func PossibleValuesForTriggerTypes() []string {
	return []string{
		string(TriggerTypesHTTPTrigger),
		string(TriggerTypesUnknown),
	}
}

func (s *TriggerTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggerTypes(input string) (*TriggerTypes, error) {
	vals := map[string]TriggerTypes{
		"httptrigger": TriggerTypesHTTPTrigger,
		"unknown":     TriggerTypesUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerTypes(input)
	return &out, nil
}
