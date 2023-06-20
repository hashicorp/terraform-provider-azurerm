package replicationpolicies

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SetMultiVMSyncStatus string

const (
	SetMultiVMSyncStatusDisable SetMultiVMSyncStatus = "Disable"
	SetMultiVMSyncStatusEnable  SetMultiVMSyncStatus = "Enable"
)

func PossibleValuesForSetMultiVMSyncStatus() []string {
	return []string{
		string(SetMultiVMSyncStatusDisable),
		string(SetMultiVMSyncStatusEnable),
	}
}

func parseSetMultiVMSyncStatus(input string) (*SetMultiVMSyncStatus, error) {
	vals := map[string]SetMultiVMSyncStatus{
		"disable": SetMultiVMSyncStatusDisable,
		"enable":  SetMultiVMSyncStatusEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SetMultiVMSyncStatus(input)
	return &out, nil
}
