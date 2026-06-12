package autoexportjobs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoExportJobAdminStatus string

const (
	AutoExportJobAdminStatusDisable AutoExportJobAdminStatus = "Disable"
	AutoExportJobAdminStatusEnable  AutoExportJobAdminStatus = "Enable"
)

func PossibleValuesForAutoExportJobAdminStatus() []string {
	return []string{
		string(AutoExportJobAdminStatusDisable),
		string(AutoExportJobAdminStatusEnable),
	}
}

func (s *AutoExportJobAdminStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoExportJobAdminStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoExportJobAdminStatus(input string) (*AutoExportJobAdminStatus, error) {
	vals := map[string]AutoExportJobAdminStatus{
		"disable": AutoExportJobAdminStatusDisable,
		"enable":  AutoExportJobAdminStatusEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoExportJobAdminStatus(input)
	return &out, nil
}

type AutoExportJobProvisioningStateType string

const (
	AutoExportJobProvisioningStateTypeCanceled  AutoExportJobProvisioningStateType = "Canceled"
	AutoExportJobProvisioningStateTypeCreating  AutoExportJobProvisioningStateType = "Creating"
	AutoExportJobProvisioningStateTypeDeleting  AutoExportJobProvisioningStateType = "Deleting"
	AutoExportJobProvisioningStateTypeFailed    AutoExportJobProvisioningStateType = "Failed"
	AutoExportJobProvisioningStateTypeSucceeded AutoExportJobProvisioningStateType = "Succeeded"
	AutoExportJobProvisioningStateTypeUpdating  AutoExportJobProvisioningStateType = "Updating"
)

func PossibleValuesForAutoExportJobProvisioningStateType() []string {
	return []string{
		string(AutoExportJobProvisioningStateTypeCanceled),
		string(AutoExportJobProvisioningStateTypeCreating),
		string(AutoExportJobProvisioningStateTypeDeleting),
		string(AutoExportJobProvisioningStateTypeFailed),
		string(AutoExportJobProvisioningStateTypeSucceeded),
		string(AutoExportJobProvisioningStateTypeUpdating),
	}
}

func (s *AutoExportJobProvisioningStateType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoExportJobProvisioningStateType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoExportJobProvisioningStateType(input string) (*AutoExportJobProvisioningStateType, error) {
	vals := map[string]AutoExportJobProvisioningStateType{
		"canceled":  AutoExportJobProvisioningStateTypeCanceled,
		"creating":  AutoExportJobProvisioningStateTypeCreating,
		"deleting":  AutoExportJobProvisioningStateTypeDeleting,
		"failed":    AutoExportJobProvisioningStateTypeFailed,
		"succeeded": AutoExportJobProvisioningStateTypeSucceeded,
		"updating":  AutoExportJobProvisioningStateTypeUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoExportJobProvisioningStateType(input)
	return &out, nil
}

type AutoExportStatusType string

const (
	AutoExportStatusTypeDisableFailed AutoExportStatusType = "DisableFailed"
	AutoExportStatusTypeDisabled      AutoExportStatusType = "Disabled"
	AutoExportStatusTypeDisabling     AutoExportStatusType = "Disabling"
	AutoExportStatusTypeFailed        AutoExportStatusType = "Failed"
	AutoExportStatusTypeInProgress    AutoExportStatusType = "InProgress"
)

func PossibleValuesForAutoExportStatusType() []string {
	return []string{
		string(AutoExportStatusTypeDisableFailed),
		string(AutoExportStatusTypeDisabled),
		string(AutoExportStatusTypeDisabling),
		string(AutoExportStatusTypeFailed),
		string(AutoExportStatusTypeInProgress),
	}
}

func (s *AutoExportStatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoExportStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoExportStatusType(input string) (*AutoExportStatusType, error) {
	vals := map[string]AutoExportStatusType{
		"disablefailed": AutoExportStatusTypeDisableFailed,
		"disabled":      AutoExportStatusTypeDisabled,
		"disabling":     AutoExportStatusTypeDisabling,
		"failed":        AutoExportStatusTypeFailed,
		"inprogress":    AutoExportStatusTypeInProgress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoExportStatusType(input)
	return &out, nil
}
