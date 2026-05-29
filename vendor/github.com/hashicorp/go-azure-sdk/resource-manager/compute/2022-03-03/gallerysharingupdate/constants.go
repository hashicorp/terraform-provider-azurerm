package gallerysharingupdate

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharingProfileGroupTypes string

const (
	SharingProfileGroupTypesAADTenants    SharingProfileGroupTypes = "AADTenants"
	SharingProfileGroupTypesSubscriptions SharingProfileGroupTypes = "Subscriptions"
)

func PossibleValuesForSharingProfileGroupTypes() []string {
	return []string{
		string(SharingProfileGroupTypesAADTenants),
		string(SharingProfileGroupTypesSubscriptions),
	}
}

func (s *SharingProfileGroupTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharingProfileGroupTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharingProfileGroupTypes(input string) (*SharingProfileGroupTypes, error) {
	vals := map[string]SharingProfileGroupTypes{
		"aadtenants":    SharingProfileGroupTypesAADTenants,
		"subscriptions": SharingProfileGroupTypesSubscriptions,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharingProfileGroupTypes(input)
	return &out, nil
}

type SharingUpdateOperationTypes string

const (
	SharingUpdateOperationTypesAdd             SharingUpdateOperationTypes = "Add"
	SharingUpdateOperationTypesEnableCommunity SharingUpdateOperationTypes = "EnableCommunity"
	SharingUpdateOperationTypesRemove          SharingUpdateOperationTypes = "Remove"
	SharingUpdateOperationTypesReset           SharingUpdateOperationTypes = "Reset"
)

func PossibleValuesForSharingUpdateOperationTypes() []string {
	return []string{
		string(SharingUpdateOperationTypesAdd),
		string(SharingUpdateOperationTypesEnableCommunity),
		string(SharingUpdateOperationTypesRemove),
		string(SharingUpdateOperationTypesReset),
	}
}

func (s *SharingUpdateOperationTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSharingUpdateOperationTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSharingUpdateOperationTypes(input string) (*SharingUpdateOperationTypes, error) {
	vals := map[string]SharingUpdateOperationTypes{
		"add":             SharingUpdateOperationTypesAdd,
		"enablecommunity": SharingUpdateOperationTypesEnableCommunity,
		"remove":          SharingUpdateOperationTypesRemove,
		"reset":           SharingUpdateOperationTypesReset,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SharingUpdateOperationTypes(input)
	return &out, nil
}
