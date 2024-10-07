package linkedservices

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedServiceEntityStatus string

const (
	LinkedServiceEntityStatusDeleting            LinkedServiceEntityStatus = "Deleting"
	LinkedServiceEntityStatusProvisioningAccount LinkedServiceEntityStatus = "ProvisioningAccount"
	LinkedServiceEntityStatusSucceeded           LinkedServiceEntityStatus = "Succeeded"
	LinkedServiceEntityStatusUpdating            LinkedServiceEntityStatus = "Updating"
)

func PossibleValuesForLinkedServiceEntityStatus() []string {
	return []string{
		string(LinkedServiceEntityStatusDeleting),
		string(LinkedServiceEntityStatusProvisioningAccount),
		string(LinkedServiceEntityStatusSucceeded),
		string(LinkedServiceEntityStatusUpdating),
	}
}

func (s *LinkedServiceEntityStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinkedServiceEntityStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinkedServiceEntityStatus(input string) (*LinkedServiceEntityStatus, error) {
	vals := map[string]LinkedServiceEntityStatus{
		"deleting":            LinkedServiceEntityStatusDeleting,
		"provisioningaccount": LinkedServiceEntityStatusProvisioningAccount,
		"succeeded":           LinkedServiceEntityStatusSucceeded,
		"updating":            LinkedServiceEntityStatusUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinkedServiceEntityStatus(input)
	return &out, nil
}
