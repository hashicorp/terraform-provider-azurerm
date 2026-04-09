package connectionmonitors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectionMonitorEndpointFilterItemType string

const (
	ConnectionMonitorEndpointFilterItemTypeAgentAddress ConnectionMonitorEndpointFilterItemType = "AgentAddress"
)

func PossibleValuesForConnectionMonitorEndpointFilterItemType() []string {
	return []string{
		string(ConnectionMonitorEndpointFilterItemTypeAgentAddress),
	}
}

func (s *ConnectionMonitorEndpointFilterItemType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionMonitorEndpointFilterItemType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionMonitorEndpointFilterItemType(input string) (*ConnectionMonitorEndpointFilterItemType, error) {
	vals := map[string]ConnectionMonitorEndpointFilterItemType{
		"agentaddress": ConnectionMonitorEndpointFilterItemTypeAgentAddress,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionMonitorEndpointFilterItemType(input)
	return &out, nil
}

type ConnectionMonitorEndpointFilterType string

const (
	ConnectionMonitorEndpointFilterTypeInclude ConnectionMonitorEndpointFilterType = "Include"
)

func PossibleValuesForConnectionMonitorEndpointFilterType() []string {
	return []string{
		string(ConnectionMonitorEndpointFilterTypeInclude),
	}
}

func (s *ConnectionMonitorEndpointFilterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionMonitorEndpointFilterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionMonitorEndpointFilterType(input string) (*ConnectionMonitorEndpointFilterType, error) {
	vals := map[string]ConnectionMonitorEndpointFilterType{
		"include": ConnectionMonitorEndpointFilterTypeInclude,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionMonitorEndpointFilterType(input)
	return &out, nil
}

type ConnectionMonitorTestConfigurationProtocol string

const (
	ConnectionMonitorTestConfigurationProtocolHTTP ConnectionMonitorTestConfigurationProtocol = "Http"
	ConnectionMonitorTestConfigurationProtocolIcmp ConnectionMonitorTestConfigurationProtocol = "Icmp"
	ConnectionMonitorTestConfigurationProtocolTcp  ConnectionMonitorTestConfigurationProtocol = "Tcp"
)

func PossibleValuesForConnectionMonitorTestConfigurationProtocol() []string {
	return []string{
		string(ConnectionMonitorTestConfigurationProtocolHTTP),
		string(ConnectionMonitorTestConfigurationProtocolIcmp),
		string(ConnectionMonitorTestConfigurationProtocolTcp),
	}
}

func (s *ConnectionMonitorTestConfigurationProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionMonitorTestConfigurationProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionMonitorTestConfigurationProtocol(input string) (*ConnectionMonitorTestConfigurationProtocol, error) {
	vals := map[string]ConnectionMonitorTestConfigurationProtocol{
		"http": ConnectionMonitorTestConfigurationProtocolHTTP,
		"icmp": ConnectionMonitorTestConfigurationProtocolIcmp,
		"tcp":  ConnectionMonitorTestConfigurationProtocolTcp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionMonitorTestConfigurationProtocol(input)
	return &out, nil
}

type ConnectionMonitorType string

const (
	ConnectionMonitorTypeMultiEndpoint           ConnectionMonitorType = "MultiEndpoint"
	ConnectionMonitorTypeSingleSourceDestination ConnectionMonitorType = "SingleSourceDestination"
)

func PossibleValuesForConnectionMonitorType() []string {
	return []string{
		string(ConnectionMonitorTypeMultiEndpoint),
		string(ConnectionMonitorTypeSingleSourceDestination),
	}
}

func (s *ConnectionMonitorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectionMonitorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectionMonitorType(input string) (*ConnectionMonitorType, error) {
	vals := map[string]ConnectionMonitorType{
		"multiendpoint":           ConnectionMonitorTypeMultiEndpoint,
		"singlesourcedestination": ConnectionMonitorTypeSingleSourceDestination,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectionMonitorType(input)
	return &out, nil
}

type CoverageLevel string

const (
	CoverageLevelAboveAverage CoverageLevel = "AboveAverage"
	CoverageLevelAverage      CoverageLevel = "Average"
	CoverageLevelBelowAverage CoverageLevel = "BelowAverage"
	CoverageLevelDefault      CoverageLevel = "Default"
	CoverageLevelFull         CoverageLevel = "Full"
	CoverageLevelLow          CoverageLevel = "Low"
)

func PossibleValuesForCoverageLevel() []string {
	return []string{
		string(CoverageLevelAboveAverage),
		string(CoverageLevelAverage),
		string(CoverageLevelBelowAverage),
		string(CoverageLevelDefault),
		string(CoverageLevelFull),
		string(CoverageLevelLow),
	}
}

func (s *CoverageLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCoverageLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCoverageLevel(input string) (*CoverageLevel, error) {
	vals := map[string]CoverageLevel{
		"aboveaverage": CoverageLevelAboveAverage,
		"average":      CoverageLevelAverage,
		"belowaverage": CoverageLevelBelowAverage,
		"default":      CoverageLevelDefault,
		"full":         CoverageLevelFull,
		"low":          CoverageLevelLow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CoverageLevel(input)
	return &out, nil
}

type DestinationPortBehavior string

const (
	DestinationPortBehaviorListenIfAvailable DestinationPortBehavior = "ListenIfAvailable"
	DestinationPortBehaviorNone              DestinationPortBehavior = "None"
)

func PossibleValuesForDestinationPortBehavior() []string {
	return []string{
		string(DestinationPortBehaviorListenIfAvailable),
		string(DestinationPortBehaviorNone),
	}
}

func (s *DestinationPortBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDestinationPortBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDestinationPortBehavior(input string) (*DestinationPortBehavior, error) {
	vals := map[string]DestinationPortBehavior{
		"listenifavailable": DestinationPortBehaviorListenIfAvailable,
		"none":              DestinationPortBehaviorNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DestinationPortBehavior(input)
	return &out, nil
}

type EndpointType string

const (
	EndpointTypeAzureArcNetwork     EndpointType = "AzureArcNetwork"
	EndpointTypeAzureArcVM          EndpointType = "AzureArcVM"
	EndpointTypeAzureSubnet         EndpointType = "AzureSubnet"
	EndpointTypeAzureVM             EndpointType = "AzureVM"
	EndpointTypeAzureVMSS           EndpointType = "AzureVMSS"
	EndpointTypeAzureVNet           EndpointType = "AzureVNet"
	EndpointTypeExternalAddress     EndpointType = "ExternalAddress"
	EndpointTypeMMAWorkspaceMachine EndpointType = "MMAWorkspaceMachine"
	EndpointTypeMMAWorkspaceNetwork EndpointType = "MMAWorkspaceNetwork"
)

func PossibleValuesForEndpointType() []string {
	return []string{
		string(EndpointTypeAzureArcNetwork),
		string(EndpointTypeAzureArcVM),
		string(EndpointTypeAzureSubnet),
		string(EndpointTypeAzureVM),
		string(EndpointTypeAzureVMSS),
		string(EndpointTypeAzureVNet),
		string(EndpointTypeExternalAddress),
		string(EndpointTypeMMAWorkspaceMachine),
		string(EndpointTypeMMAWorkspaceNetwork),
	}
}

func (s *EndpointType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEndpointType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEndpointType(input string) (*EndpointType, error) {
	vals := map[string]EndpointType{
		"azurearcnetwork":     EndpointTypeAzureArcNetwork,
		"azurearcvm":          EndpointTypeAzureArcVM,
		"azuresubnet":         EndpointTypeAzureSubnet,
		"azurevm":             EndpointTypeAzureVM,
		"azurevmss":           EndpointTypeAzureVMSS,
		"azurevnet":           EndpointTypeAzureVNet,
		"externaladdress":     EndpointTypeExternalAddress,
		"mmaworkspacemachine": EndpointTypeMMAWorkspaceMachine,
		"mmaworkspacenetwork": EndpointTypeMMAWorkspaceNetwork,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EndpointType(input)
	return &out, nil
}

type HTTPConfigurationMethod string

const (
	HTTPConfigurationMethodGet  HTTPConfigurationMethod = "Get"
	HTTPConfigurationMethodPost HTTPConfigurationMethod = "Post"
)

func PossibleValuesForHTTPConfigurationMethod() []string {
	return []string{
		string(HTTPConfigurationMethodGet),
		string(HTTPConfigurationMethodPost),
	}
}

func (s *HTTPConfigurationMethod) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPConfigurationMethod(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPConfigurationMethod(input string) (*HTTPConfigurationMethod, error) {
	vals := map[string]HTTPConfigurationMethod{
		"get":  HTTPConfigurationMethodGet,
		"post": HTTPConfigurationMethodPost,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPConfigurationMethod(input)
	return &out, nil
}

type OutputType string

const (
	OutputTypeWorkspace OutputType = "Workspace"
)

func PossibleValuesForOutputType() []string {
	return []string{
		string(OutputTypeWorkspace),
	}
}

func (s *OutputType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOutputType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOutputType(input string) (*OutputType, error) {
	vals := map[string]OutputType{
		"workspace": OutputTypeWorkspace,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OutputType(input)
	return &out, nil
}

type PreferredIPVersion string

const (
	PreferredIPVersionIPvFour PreferredIPVersion = "IPv4"
	PreferredIPVersionIPvSix  PreferredIPVersion = "IPv6"
)

func PossibleValuesForPreferredIPVersion() []string {
	return []string{
		string(PreferredIPVersionIPvFour),
		string(PreferredIPVersionIPvSix),
	}
}

func (s *PreferredIPVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePreferredIPVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePreferredIPVersion(input string) (*PreferredIPVersion, error) {
	vals := map[string]PreferredIPVersion{
		"ipv4": PreferredIPVersionIPvFour,
		"ipv6": PreferredIPVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PreferredIPVersion(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
