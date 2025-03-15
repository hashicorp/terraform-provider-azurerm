package gatewaylistdebugcredentials

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayListDebugCredentialsContractPurpose string

const (
	GatewayListDebugCredentialsContractPurposeTracing GatewayListDebugCredentialsContractPurpose = "tracing"
)

func PossibleValuesForGatewayListDebugCredentialsContractPurpose() []string {
	return []string{
		string(GatewayListDebugCredentialsContractPurposeTracing),
	}
}

func (s *GatewayListDebugCredentialsContractPurpose) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayListDebugCredentialsContractPurpose(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayListDebugCredentialsContractPurpose(input string) (*GatewayListDebugCredentialsContractPurpose, error) {
	vals := map[string]GatewayListDebugCredentialsContractPurpose{
		"tracing": GatewayListDebugCredentialsContractPurposeTracing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayListDebugCredentialsContractPurpose(input)
	return &out, nil
}
