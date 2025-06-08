package namedvalue

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyVaultRefreshState string

const (
	KeyVaultRefreshStateFalse KeyVaultRefreshState = "false"
	KeyVaultRefreshStateTrue  KeyVaultRefreshState = "true"
)

func PossibleValuesForKeyVaultRefreshState() []string {
	return []string{
		string(KeyVaultRefreshStateFalse),
		string(KeyVaultRefreshStateTrue),
	}
}

func (s *KeyVaultRefreshState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyVaultRefreshState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyVaultRefreshState(input string) (*KeyVaultRefreshState, error) {
	vals := map[string]KeyVaultRefreshState{
		"false": KeyVaultRefreshStateFalse,
		"true":  KeyVaultRefreshStateTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultRefreshState(input)
	return &out, nil
}
