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
	NetAppProvisioningStateUpdating  NetAppProvisioningState = "Updating"
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
		string(NetAppProvisioningStateUpdating),
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
		"updating":  NetAppProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetAppProvisioningState(input)
	return &out, nil
}

type QuotaType string

const (
	QuotaTypeDefaultGroupQuota    QuotaType = "DefaultGroupQuota"
	QuotaTypeDefaultUserQuota     QuotaType = "DefaultUserQuota"
	QuotaTypeIndividualGroupQuota QuotaType = "IndividualGroupQuota"
	QuotaTypeIndividualUserQuota  QuotaType = "IndividualUserQuota"
)

func PossibleValuesForQuotaType() []string {
	return []string{
		string(QuotaTypeDefaultGroupQuota),
		string(QuotaTypeDefaultUserQuota),
		string(QuotaTypeIndividualGroupQuota),
		string(QuotaTypeIndividualUserQuota),
	}
}

func (s *QuotaType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQuotaType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQuotaType(input string) (*QuotaType, error) {
	vals := map[string]QuotaType{
		"defaultgroupquota":    QuotaTypeDefaultGroupQuota,
		"defaultuserquota":     QuotaTypeDefaultUserQuota,
		"individualgroupquota": QuotaTypeIndividualGroupQuota,
		"individualuserquota":  QuotaTypeIndividualUserQuota,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QuotaType(input)
	return &out, nil
}
