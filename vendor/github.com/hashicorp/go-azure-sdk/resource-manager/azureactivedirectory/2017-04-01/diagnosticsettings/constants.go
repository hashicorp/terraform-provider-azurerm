package diagnosticsettings

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Category string

const (
	CategoryAuditLogs  Category = "AuditLogs"
	CategorySignInLogs Category = "SignInLogs"
)

func PossibleValuesForCategory() []string {
	return []string{
		string(CategoryAuditLogs),
		string(CategorySignInLogs),
	}
}

func (s *Category) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCategory(input string) (*Category, error) {
	vals := map[string]Category{
		"auditlogs":  CategoryAuditLogs,
		"signinlogs": CategorySignInLogs,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Category(input)
	return &out, nil
}
