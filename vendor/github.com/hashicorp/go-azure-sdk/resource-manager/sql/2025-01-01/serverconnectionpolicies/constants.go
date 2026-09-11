package serverconnectionpolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerConnectionType string

const (
	ServerConnectionTypeDefault  ServerConnectionType = "Default"
	ServerConnectionTypeProxy    ServerConnectionType = "Proxy"
	ServerConnectionTypeRedirect ServerConnectionType = "Redirect"
)

func PossibleValuesForServerConnectionType() []string {
	return []string{
		string(ServerConnectionTypeDefault),
		string(ServerConnectionTypeProxy),
		string(ServerConnectionTypeRedirect),
	}
}

func (s *ServerConnectionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerConnectionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerConnectionType(input string) (*ServerConnectionType, error) {
	vals := map[string]ServerConnectionType{
		"default":  ServerConnectionTypeDefault,
		"proxy":    ServerConnectionTypeProxy,
		"redirect": ServerConnectionTypeRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerConnectionType(input)
	return &out, nil
}
