package replicationpolicies

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SetMultiVmSyncStatus string

const (
	SetMultiVmSyncStatusDisable SetMultiVmSyncStatus = "Disable"
	SetMultiVmSyncStatusEnable  SetMultiVmSyncStatus = "Enable"
)

func PossibleValuesForSetMultiVmSyncStatus() []string {
	return []string{
		string(SetMultiVmSyncStatusDisable),
		string(SetMultiVmSyncStatusEnable),
	}
}

func parseSetMultiVmSyncStatus(input string) (*SetMultiVmSyncStatus, error) {
	vals := map[string]SetMultiVmSyncStatus{
		"disable": SetMultiVmSyncStatusDisable,
		"enable":  SetMultiVmSyncStatusEnable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SetMultiVmSyncStatus(input)
	return &out, nil
}
