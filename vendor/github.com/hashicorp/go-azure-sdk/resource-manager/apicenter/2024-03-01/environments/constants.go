package environments

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentKind string

const (
	EnvironmentKindDevelopment EnvironmentKind = "development"
	EnvironmentKindProduction  EnvironmentKind = "production"
	EnvironmentKindStaging     EnvironmentKind = "staging"
	EnvironmentKindTesting     EnvironmentKind = "testing"
)

func PossibleValuesForEnvironmentKind() []string {
	return []string{
		string(EnvironmentKindDevelopment),
		string(EnvironmentKindProduction),
		string(EnvironmentKindStaging),
		string(EnvironmentKindTesting),
	}
}

func (s *EnvironmentKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnvironmentKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnvironmentKind(input string) (*EnvironmentKind, error) {
	vals := map[string]EnvironmentKind{
		"development": EnvironmentKindDevelopment,
		"production":  EnvironmentKindProduction,
		"staging":     EnvironmentKindStaging,
		"testing":     EnvironmentKindTesting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnvironmentKind(input)
	return &out, nil
}

type EnvironmentServerType string

const (
	EnvironmentServerTypeAWSAPIGateway         EnvironmentServerType = "AWS API Gateway"
	EnvironmentServerTypeApigeeAPIManagement   EnvironmentServerType = "Apigee API Management"
	EnvironmentServerTypeAzureAPIManagement    EnvironmentServerType = "Azure API Management"
	EnvironmentServerTypeAzureComputeService   EnvironmentServerType = "Azure compute service"
	EnvironmentServerTypeKongAPIGateway        EnvironmentServerType = "Kong API Gateway"
	EnvironmentServerTypeKubernetes            EnvironmentServerType = "Kubernetes"
	EnvironmentServerTypeMuleSoftAPIManagement EnvironmentServerType = "MuleSoft API Management"
)

func PossibleValuesForEnvironmentServerType() []string {
	return []string{
		string(EnvironmentServerTypeAWSAPIGateway),
		string(EnvironmentServerTypeApigeeAPIManagement),
		string(EnvironmentServerTypeAzureAPIManagement),
		string(EnvironmentServerTypeAzureComputeService),
		string(EnvironmentServerTypeKongAPIGateway),
		string(EnvironmentServerTypeKubernetes),
		string(EnvironmentServerTypeMuleSoftAPIManagement),
	}
}

func (s *EnvironmentServerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnvironmentServerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnvironmentServerType(input string) (*EnvironmentServerType, error) {
	vals := map[string]EnvironmentServerType{
		"aws api gateway":         EnvironmentServerTypeAWSAPIGateway,
		"apigee api management":   EnvironmentServerTypeApigeeAPIManagement,
		"azure api management":    EnvironmentServerTypeAzureAPIManagement,
		"azure compute service":   EnvironmentServerTypeAzureComputeService,
		"kong api gateway":        EnvironmentServerTypeKongAPIGateway,
		"kubernetes":              EnvironmentServerTypeKubernetes,
		"mulesoft api management": EnvironmentServerTypeMuleSoftAPIManagement,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnvironmentServerType(input)
	return &out, nil
}
