package managedinstancelongtermretentionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupStorageAccessTier string

const (
	BackupStorageAccessTierArchive BackupStorageAccessTier = "Archive"
	BackupStorageAccessTierHot     BackupStorageAccessTier = "Hot"
)

func PossibleValuesForBackupStorageAccessTier() []string {
	return []string{
		string(BackupStorageAccessTierArchive),
		string(BackupStorageAccessTierHot),
	}
}

func (s *BackupStorageAccessTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupStorageAccessTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupStorageAccessTier(input string) (*BackupStorageAccessTier, error) {
	vals := map[string]BackupStorageAccessTier{
		"archive": BackupStorageAccessTierArchive,
		"hot":     BackupStorageAccessTierHot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupStorageAccessTier(input)
	return &out, nil
}
