package backupshorttermretentionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiffBackupIntervalInHours float64

const (
	DiffBackupIntervalInHoursOneTwo  DiffBackupIntervalInHours = 12
	DiffBackupIntervalInHoursTwoFour DiffBackupIntervalInHours = 24
)

func PossibleValuesForDiffBackupIntervalInHours() []float64 {
	return []float64{
		float64(DiffBackupIntervalInHoursOneTwo),
		float64(DiffBackupIntervalInHoursTwoFour),
	}
}

type ImmutabilityStatus string

const (
	ImmutabilityStatusDisabled ImmutabilityStatus = "Disabled"
	ImmutabilityStatusEnabled  ImmutabilityStatus = "Enabled"
	ImmutabilityStatusLocked   ImmutabilityStatus = "Locked"
)

func PossibleValuesForImmutabilityStatus() []string {
	return []string{
		string(ImmutabilityStatusDisabled),
		string(ImmutabilityStatusEnabled),
		string(ImmutabilityStatusLocked),
	}
}

func (s *ImmutabilityStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImmutabilityStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImmutabilityStatus(input string) (*ImmutabilityStatus, error) {
	vals := map[string]ImmutabilityStatus{
		"disabled": ImmutabilityStatusDisabled,
		"enabled":  ImmutabilityStatusEnabled,
		"locked":   ImmutabilityStatusLocked,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImmutabilityStatus(input)
	return &out, nil
}
