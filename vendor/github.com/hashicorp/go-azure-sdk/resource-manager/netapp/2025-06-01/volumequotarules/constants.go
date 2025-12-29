package volumequotarules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetAppProvisioningState string

const (
	NetAppProvisioningStateAccepted  NetAppProvisioningState = "Accepted"
	NetAppProvisioningStateCreating  NetAppProvisioningState = "Creating"
	NetAppProvisioningStateDeleting  NetAppProvisioningState = "Deleting"
	NetAppProvisioningStateFailed    NetAppProvisioningState = "Failed"
	NetAppProvisioningStateMoving    NetAppProvisioningState = "Moving"
	NetAppProvisioningStatePatching  NetAppProvisioningState = "Patching"
	NetAppProvisioningStateSucceeded NetAppProvisioningState = "Succeeded"
)

func PossibleValuesForNetAppProvisioningState() []string {
	return []string{
		string(NetAppProvisioningStateAccepted),
		string(NetAppProvisioningStateCreating),
		string(NetAppProvisioningStateDeleting),
		string(NetAppProvisioningStateFailed),
		string(NetAppProvisioningStateMoving),
		string(NetAppProvisioningStatePatching),
		string(NetAppProvisioningStateSucceeded),
	}
}

func (s *NetAppProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetAppProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetAppProvisioningState(input string) (*NetAppProvisioningState, error) {
	vals := map[string]NetAppProvisioningState{
		"accepted":  NetAppProvisioningStateAccepted,
		"creating":  NetAppProvisioningStateCreating,
		"deleting":  NetAppProvisioningStateDeleting,
		"failed":    NetAppProvisioningStateFailed,
		"moving":    NetAppProvisioningStateMoving,
		"patching":  NetAppProvisioningStatePatching,
		"succeeded": NetAppProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetAppProvisioningState(input)
	return &out, nil
}

type Type string

const (
	TypeDefaultGroupQuota    Type = "DefaultGroupQuota"
	TypeDefaultUserQuota     Type = "DefaultUserQuota"
	TypeIndividualGroupQuota Type = "IndividualGroupQuota"
	TypeIndividualUserQuota  Type = "IndividualUserQuota"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeDefaultGroupQuota),
		string(TypeDefaultUserQuota),
		string(TypeIndividualGroupQuota),
		string(TypeIndividualUserQuota),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"defaultgroupquota":    TypeDefaultGroupQuota,
		"defaultuserquota":     TypeDefaultUserQuota,
		"individualgroupquota": TypeIndividualGroupQuota,
		"individualuserquota":  TypeIndividualUserQuota,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
