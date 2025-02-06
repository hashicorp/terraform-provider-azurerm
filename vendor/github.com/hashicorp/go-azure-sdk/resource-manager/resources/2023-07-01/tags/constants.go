package tags

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagsPatchOperation string

const (
	TagsPatchOperationDelete  TagsPatchOperation = "Delete"
	TagsPatchOperationMerge   TagsPatchOperation = "Merge"
	TagsPatchOperationReplace TagsPatchOperation = "Replace"
)

func PossibleValuesForTagsPatchOperation() []string {
	return []string{
		string(TagsPatchOperationDelete),
		string(TagsPatchOperationMerge),
		string(TagsPatchOperationReplace),
	}
}

func (s *TagsPatchOperation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTagsPatchOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTagsPatchOperation(input string) (*TagsPatchOperation, error) {
	vals := map[string]TagsPatchOperation{
		"delete":  TagsPatchOperationDelete,
		"merge":   TagsPatchOperationMerge,
		"replace": TagsPatchOperationReplace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TagsPatchOperation(input)
	return &out, nil
}
