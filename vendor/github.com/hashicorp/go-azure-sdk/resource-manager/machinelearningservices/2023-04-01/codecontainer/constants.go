package codecontainer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssetProvisioningState string

const (
	AssetProvisioningStateCanceled  AssetProvisioningState = "Canceled"
	AssetProvisioningStateCreating  AssetProvisioningState = "Creating"
	AssetProvisioningStateDeleting  AssetProvisioningState = "Deleting"
	AssetProvisioningStateFailed    AssetProvisioningState = "Failed"
	AssetProvisioningStateSucceeded AssetProvisioningState = "Succeeded"
	AssetProvisioningStateUpdating  AssetProvisioningState = "Updating"
)

func PossibleValuesForAssetProvisioningState() []string {
	return []string{
		string(AssetProvisioningStateCanceled),
		string(AssetProvisioningStateCreating),
		string(AssetProvisioningStateDeleting),
		string(AssetProvisioningStateFailed),
		string(AssetProvisioningStateSucceeded),
		string(AssetProvisioningStateUpdating),
	}
}

func (s *AssetProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssetProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssetProvisioningState(input string) (*AssetProvisioningState, error) {
	vals := map[string]AssetProvisioningState{
		"canceled":  AssetProvisioningStateCanceled,
		"creating":  AssetProvisioningStateCreating,
		"deleting":  AssetProvisioningStateDeleting,
		"failed":    AssetProvisioningStateFailed,
		"succeeded": AssetProvisioningStateSucceeded,
		"updating":  AssetProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssetProvisioningState(input)
	return &out, nil
}
