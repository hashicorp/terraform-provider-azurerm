package sqlvirtualmachinegroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterConfiguration string

const (
	ClusterConfigurationDomainful ClusterConfiguration = "Domainful"
)

func PossibleValuesForClusterConfiguration() []string {
	return []string{
		string(ClusterConfigurationDomainful),
	}
}

func (s *ClusterConfiguration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterConfiguration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterConfiguration(input string) (*ClusterConfiguration, error) {
	vals := map[string]ClusterConfiguration{
		"domainful": ClusterConfigurationDomainful,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterConfiguration(input)
	return &out, nil
}

type ClusterManagerType string

const (
	ClusterManagerTypeWSFC ClusterManagerType = "WSFC"
)

func PossibleValuesForClusterManagerType() []string {
	return []string{
		string(ClusterManagerTypeWSFC),
	}
}

func (s *ClusterManagerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterManagerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterManagerType(input string) (*ClusterManagerType, error) {
	vals := map[string]ClusterManagerType{
		"wsfc": ClusterManagerTypeWSFC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterManagerType(input)
	return &out, nil
}

type ClusterSubnetType string

const (
	ClusterSubnetTypeMultiSubnet  ClusterSubnetType = "MultiSubnet"
	ClusterSubnetTypeSingleSubnet ClusterSubnetType = "SingleSubnet"
)

func PossibleValuesForClusterSubnetType() []string {
	return []string{
		string(ClusterSubnetTypeMultiSubnet),
		string(ClusterSubnetTypeSingleSubnet),
	}
}

func (s *ClusterSubnetType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterSubnetType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterSubnetType(input string) (*ClusterSubnetType, error) {
	vals := map[string]ClusterSubnetType{
		"multisubnet":  ClusterSubnetTypeMultiSubnet,
		"singlesubnet": ClusterSubnetTypeSingleSubnet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterSubnetType(input)
	return &out, nil
}

type ScaleType string

const (
	ScaleTypeHA ScaleType = "HA"
)

func PossibleValuesForScaleType() []string {
	return []string{
		string(ScaleTypeHA),
	}
}

func (s *ScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleType(input string) (*ScaleType, error) {
	vals := map[string]ScaleType{
		"ha": ScaleTypeHA,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleType(input)
	return &out, nil
}

type SqlVMGroupImageSku string

const (
	SqlVMGroupImageSkuDeveloper  SqlVMGroupImageSku = "Developer"
	SqlVMGroupImageSkuEnterprise SqlVMGroupImageSku = "Enterprise"
)

func PossibleValuesForSqlVMGroupImageSku() []string {
	return []string{
		string(SqlVMGroupImageSkuDeveloper),
		string(SqlVMGroupImageSkuEnterprise),
	}
}

func (s *SqlVMGroupImageSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlVMGroupImageSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlVMGroupImageSku(input string) (*SqlVMGroupImageSku, error) {
	vals := map[string]SqlVMGroupImageSku{
		"developer":  SqlVMGroupImageSkuDeveloper,
		"enterprise": SqlVMGroupImageSkuEnterprise,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlVMGroupImageSku(input)
	return &out, nil
}
