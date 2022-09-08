package netappaccounts

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActiveDirectoryStatus string

const (
	ActiveDirectoryStatusCreated  ActiveDirectoryStatus = "Created"
	ActiveDirectoryStatusDeleted  ActiveDirectoryStatus = "Deleted"
	ActiveDirectoryStatusError    ActiveDirectoryStatus = "Error"
	ActiveDirectoryStatusInUse    ActiveDirectoryStatus = "InUse"
	ActiveDirectoryStatusUpdating ActiveDirectoryStatus = "Updating"
)

func PossibleValuesForActiveDirectoryStatus() []string {
	return []string{
		string(ActiveDirectoryStatusCreated),
		string(ActiveDirectoryStatusDeleted),
		string(ActiveDirectoryStatusError),
		string(ActiveDirectoryStatusInUse),
		string(ActiveDirectoryStatusUpdating),
	}
}

func parseActiveDirectoryStatus(input string) (*ActiveDirectoryStatus, error) {
	vals := map[string]ActiveDirectoryStatus{
		"created":  ActiveDirectoryStatusCreated,
		"deleted":  ActiveDirectoryStatusDeleted,
		"error":    ActiveDirectoryStatusError,
		"inuse":    ActiveDirectoryStatusInUse,
		"updating": ActiveDirectoryStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActiveDirectoryStatus(input)
	return &out, nil
}
