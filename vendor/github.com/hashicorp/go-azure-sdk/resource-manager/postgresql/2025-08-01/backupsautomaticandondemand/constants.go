package backupsautomaticandondemand

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackupType string

const (
	BackupTypeCustomerOnNegativeDemand BackupType = "Customer On-Demand"
	BackupTypeFull                     BackupType = "Full"
)

func PossibleValuesForBackupType() []string {
	return []string{
		string(BackupTypeCustomerOnNegativeDemand),
		string(BackupTypeFull),
	}
}

func (s *BackupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupType(input string) (*BackupType, error) {
	vals := map[string]BackupType{
		"customer on-demand": BackupTypeCustomerOnNegativeDemand,
		"full":               BackupTypeFull,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupType(input)
	return &out, nil
}
