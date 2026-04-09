package policies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyEvaluatorType string

const (
	PolicyEvaluatorTypeAllowedValuesPolicy PolicyEvaluatorType = "AllowedValuesPolicy"
	PolicyEvaluatorTypeMaxValuePolicy      PolicyEvaluatorType = "MaxValuePolicy"
)

func PossibleValuesForPolicyEvaluatorType() []string {
	return []string{
		string(PolicyEvaluatorTypeAllowedValuesPolicy),
		string(PolicyEvaluatorTypeMaxValuePolicy),
	}
}

func (s *PolicyEvaluatorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyEvaluatorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyEvaluatorType(input string) (*PolicyEvaluatorType, error) {
	vals := map[string]PolicyEvaluatorType{
		"allowedvaluespolicy": PolicyEvaluatorTypeAllowedValuesPolicy,
		"maxvaluepolicy":      PolicyEvaluatorTypeMaxValuePolicy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyEvaluatorType(input)
	return &out, nil
}

type PolicyFactName string

const (
	PolicyFactNameEnvironmentTemplate         PolicyFactName = "EnvironmentTemplate"
	PolicyFactNameGalleryImage                PolicyFactName = "GalleryImage"
	PolicyFactNameLabPremiumVMCount           PolicyFactName = "LabPremiumVmCount"
	PolicyFactNameLabTargetCost               PolicyFactName = "LabTargetCost"
	PolicyFactNameLabVMCount                  PolicyFactName = "LabVmCount"
	PolicyFactNameLabVMSize                   PolicyFactName = "LabVmSize"
	PolicyFactNameScheduleEditPermission      PolicyFactName = "ScheduleEditPermission"
	PolicyFactNameUserOwnedLabPremiumVMCount  PolicyFactName = "UserOwnedLabPremiumVmCount"
	PolicyFactNameUserOwnedLabVMCount         PolicyFactName = "UserOwnedLabVmCount"
	PolicyFactNameUserOwnedLabVMCountInSubnet PolicyFactName = "UserOwnedLabVmCountInSubnet"
)

func PossibleValuesForPolicyFactName() []string {
	return []string{
		string(PolicyFactNameEnvironmentTemplate),
		string(PolicyFactNameGalleryImage),
		string(PolicyFactNameLabPremiumVMCount),
		string(PolicyFactNameLabTargetCost),
		string(PolicyFactNameLabVMCount),
		string(PolicyFactNameLabVMSize),
		string(PolicyFactNameScheduleEditPermission),
		string(PolicyFactNameUserOwnedLabPremiumVMCount),
		string(PolicyFactNameUserOwnedLabVMCount),
		string(PolicyFactNameUserOwnedLabVMCountInSubnet),
	}
}

func (s *PolicyFactName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyFactName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyFactName(input string) (*PolicyFactName, error) {
	vals := map[string]PolicyFactName{
		"environmenttemplate":         PolicyFactNameEnvironmentTemplate,
		"galleryimage":                PolicyFactNameGalleryImage,
		"labpremiumvmcount":           PolicyFactNameLabPremiumVMCount,
		"labtargetcost":               PolicyFactNameLabTargetCost,
		"labvmcount":                  PolicyFactNameLabVMCount,
		"labvmsize":                   PolicyFactNameLabVMSize,
		"scheduleeditpermission":      PolicyFactNameScheduleEditPermission,
		"userownedlabpremiumvmcount":  PolicyFactNameUserOwnedLabPremiumVMCount,
		"userownedlabvmcount":         PolicyFactNameUserOwnedLabVMCount,
		"userownedlabvmcountinsubnet": PolicyFactNameUserOwnedLabVMCountInSubnet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyFactName(input)
	return &out, nil
}

type PolicyStatus string

const (
	PolicyStatusDisabled PolicyStatus = "Disabled"
	PolicyStatusEnabled  PolicyStatus = "Enabled"
)

func PossibleValuesForPolicyStatus() []string {
	return []string{
		string(PolicyStatusDisabled),
		string(PolicyStatusEnabled),
	}
}

func (s *PolicyStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyStatus(input string) (*PolicyStatus, error) {
	vals := map[string]PolicyStatus{
		"disabled": PolicyStatusDisabled,
		"enabled":  PolicyStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyStatus(input)
	return &out, nil
}
