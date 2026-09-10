package agentpools

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentPoolMode string

const (
	AgentPoolModeGateway AgentPoolMode = "Gateway"
	AgentPoolModeSystem  AgentPoolMode = "System"
	AgentPoolModeUser    AgentPoolMode = "User"
)

func PossibleValuesForAgentPoolMode() []string {
	return []string{
		string(AgentPoolModeGateway),
		string(AgentPoolModeSystem),
		string(AgentPoolModeUser),
	}
}

func (s *AgentPoolMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentPoolMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentPoolMode(input string) (*AgentPoolMode, error) {
	vals := map[string]AgentPoolMode{
		"gateway": AgentPoolModeGateway,
		"system":  AgentPoolModeSystem,
		"user":    AgentPoolModeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentPoolMode(input)
	return &out, nil
}

type AgentPoolSSHAccess string

const (
	AgentPoolSSHAccessDisabled  AgentPoolSSHAccess = "Disabled"
	AgentPoolSSHAccessLocalUser AgentPoolSSHAccess = "LocalUser"
)

func PossibleValuesForAgentPoolSSHAccess() []string {
	return []string{
		string(AgentPoolSSHAccessDisabled),
		string(AgentPoolSSHAccessLocalUser),
	}
}

func (s *AgentPoolSSHAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentPoolSSHAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentPoolSSHAccess(input string) (*AgentPoolSSHAccess, error) {
	vals := map[string]AgentPoolSSHAccess{
		"disabled":  AgentPoolSSHAccessDisabled,
		"localuser": AgentPoolSSHAccessLocalUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentPoolSSHAccess(input)
	return &out, nil
}

type AgentPoolType string

const (
	AgentPoolTypeAvailabilitySet         AgentPoolType = "AvailabilitySet"
	AgentPoolTypeVirtualMachineScaleSets AgentPoolType = "VirtualMachineScaleSets"
	AgentPoolTypeVirtualMachines         AgentPoolType = "VirtualMachines"
)

func PossibleValuesForAgentPoolType() []string {
	return []string{
		string(AgentPoolTypeAvailabilitySet),
		string(AgentPoolTypeVirtualMachineScaleSets),
		string(AgentPoolTypeVirtualMachines),
	}
}

func (s *AgentPoolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgentPoolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgentPoolType(input string) (*AgentPoolType, error) {
	vals := map[string]AgentPoolType{
		"availabilityset":         AgentPoolTypeAvailabilitySet,
		"virtualmachinescalesets": AgentPoolTypeVirtualMachineScaleSets,
		"virtualmachines":         AgentPoolTypeVirtualMachines,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgentPoolType(input)
	return &out, nil
}

type Code string

const (
	CodeRunning Code = "Running"
	CodeStopped Code = "Stopped"
)

func PossibleValuesForCode() []string {
	return []string{
		string(CodeRunning),
		string(CodeStopped),
	}
}

func (s *Code) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCode(input string) (*Code, error) {
	vals := map[string]Code{
		"running": CodeRunning,
		"stopped": CodeStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Code(input)
	return &out, nil
}

type GPUDriver string

const (
	GPUDriverInstall GPUDriver = "Install"
	GPUDriverNone    GPUDriver = "None"
)

func PossibleValuesForGPUDriver() []string {
	return []string{
		string(GPUDriverInstall),
		string(GPUDriverNone),
	}
}

func (s *GPUDriver) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGPUDriver(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGPUDriver(input string) (*GPUDriver, error) {
	vals := map[string]GPUDriver{
		"install": GPUDriverInstall,
		"none":    GPUDriverNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GPUDriver(input)
	return &out, nil
}

type GPUInstanceProfile string

const (
	GPUInstanceProfileMIGFourg  GPUInstanceProfile = "MIG4g"
	GPUInstanceProfileMIGOneg   GPUInstanceProfile = "MIG1g"
	GPUInstanceProfileMIGSeveng GPUInstanceProfile = "MIG7g"
	GPUInstanceProfileMIGThreeg GPUInstanceProfile = "MIG3g"
	GPUInstanceProfileMIGTwog   GPUInstanceProfile = "MIG2g"
)

func PossibleValuesForGPUInstanceProfile() []string {
	return []string{
		string(GPUInstanceProfileMIGFourg),
		string(GPUInstanceProfileMIGOneg),
		string(GPUInstanceProfileMIGSeveng),
		string(GPUInstanceProfileMIGThreeg),
		string(GPUInstanceProfileMIGTwog),
	}
}

func (s *GPUInstanceProfile) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGPUInstanceProfile(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGPUInstanceProfile(input string) (*GPUInstanceProfile, error) {
	vals := map[string]GPUInstanceProfile{
		"mig4g": GPUInstanceProfileMIGFourg,
		"mig1g": GPUInstanceProfileMIGOneg,
		"mig7g": GPUInstanceProfileMIGSeveng,
		"mig3g": GPUInstanceProfileMIGThreeg,
		"mig2g": GPUInstanceProfileMIGTwog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GPUInstanceProfile(input)
	return &out, nil
}

type KubeletDiskType string

const (
	KubeletDiskTypeOS        KubeletDiskType = "OS"
	KubeletDiskTypeTemporary KubeletDiskType = "Temporary"
)

func PossibleValuesForKubeletDiskType() []string {
	return []string{
		string(KubeletDiskTypeOS),
		string(KubeletDiskTypeTemporary),
	}
}

func (s *KubeletDiskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKubeletDiskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKubeletDiskType(input string) (*KubeletDiskType, error) {
	vals := map[string]KubeletDiskType{
		"os":        KubeletDiskTypeOS,
		"temporary": KubeletDiskTypeTemporary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KubeletDiskType(input)
	return &out, nil
}

type LocalDNSForwardDestination string

const (
	LocalDNSForwardDestinationClusterCoreDNS LocalDNSForwardDestination = "ClusterCoreDNS"
	LocalDNSForwardDestinationVnetDNS        LocalDNSForwardDestination = "VnetDNS"
)

func PossibleValuesForLocalDNSForwardDestination() []string {
	return []string{
		string(LocalDNSForwardDestinationClusterCoreDNS),
		string(LocalDNSForwardDestinationVnetDNS),
	}
}

func (s *LocalDNSForwardDestination) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSForwardDestination(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSForwardDestination(input string) (*LocalDNSForwardDestination, error) {
	vals := map[string]LocalDNSForwardDestination{
		"clustercoredns": LocalDNSForwardDestinationClusterCoreDNS,
		"vnetdns":        LocalDNSForwardDestinationVnetDNS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSForwardDestination(input)
	return &out, nil
}

type LocalDNSForwardPolicy string

const (
	LocalDNSForwardPolicyRandom     LocalDNSForwardPolicy = "Random"
	LocalDNSForwardPolicyRoundRobin LocalDNSForwardPolicy = "RoundRobin"
	LocalDNSForwardPolicySequential LocalDNSForwardPolicy = "Sequential"
)

func PossibleValuesForLocalDNSForwardPolicy() []string {
	return []string{
		string(LocalDNSForwardPolicyRandom),
		string(LocalDNSForwardPolicyRoundRobin),
		string(LocalDNSForwardPolicySequential),
	}
}

func (s *LocalDNSForwardPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSForwardPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSForwardPolicy(input string) (*LocalDNSForwardPolicy, error) {
	vals := map[string]LocalDNSForwardPolicy{
		"random":     LocalDNSForwardPolicyRandom,
		"roundrobin": LocalDNSForwardPolicyRoundRobin,
		"sequential": LocalDNSForwardPolicySequential,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSForwardPolicy(input)
	return &out, nil
}

type LocalDNSMode string

const (
	LocalDNSModeDisabled  LocalDNSMode = "Disabled"
	LocalDNSModePreferred LocalDNSMode = "Preferred"
	LocalDNSModeRequired  LocalDNSMode = "Required"
)

func PossibleValuesForLocalDNSMode() []string {
	return []string{
		string(LocalDNSModeDisabled),
		string(LocalDNSModePreferred),
		string(LocalDNSModeRequired),
	}
}

func (s *LocalDNSMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSMode(input string) (*LocalDNSMode, error) {
	vals := map[string]LocalDNSMode{
		"disabled":  LocalDNSModeDisabled,
		"preferred": LocalDNSModePreferred,
		"required":  LocalDNSModeRequired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSMode(input)
	return &out, nil
}

type LocalDNSProtocol string

const (
	LocalDNSProtocolForceTCP  LocalDNSProtocol = "ForceTCP"
	LocalDNSProtocolPreferUDP LocalDNSProtocol = "PreferUDP"
)

func PossibleValuesForLocalDNSProtocol() []string {
	return []string{
		string(LocalDNSProtocolForceTCP),
		string(LocalDNSProtocolPreferUDP),
	}
}

func (s *LocalDNSProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSProtocol(input string) (*LocalDNSProtocol, error) {
	vals := map[string]LocalDNSProtocol{
		"forcetcp":  LocalDNSProtocolForceTCP,
		"preferudp": LocalDNSProtocolPreferUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSProtocol(input)
	return &out, nil
}

type LocalDNSQueryLogging string

const (
	LocalDNSQueryLoggingError LocalDNSQueryLogging = "Error"
	LocalDNSQueryLoggingLog   LocalDNSQueryLogging = "Log"
)

func PossibleValuesForLocalDNSQueryLogging() []string {
	return []string{
		string(LocalDNSQueryLoggingError),
		string(LocalDNSQueryLoggingLog),
	}
}

func (s *LocalDNSQueryLogging) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSQueryLogging(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSQueryLogging(input string) (*LocalDNSQueryLogging, error) {
	vals := map[string]LocalDNSQueryLogging{
		"error": LocalDNSQueryLoggingError,
		"log":   LocalDNSQueryLoggingLog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSQueryLogging(input)
	return &out, nil
}

type LocalDNSServeStale string

const (
	LocalDNSServeStaleDisable   LocalDNSServeStale = "Disable"
	LocalDNSServeStaleImmediate LocalDNSServeStale = "Immediate"
	LocalDNSServeStaleVerify    LocalDNSServeStale = "Verify"
)

func PossibleValuesForLocalDNSServeStale() []string {
	return []string{
		string(LocalDNSServeStaleDisable),
		string(LocalDNSServeStaleImmediate),
		string(LocalDNSServeStaleVerify),
	}
}

func (s *LocalDNSServeStale) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSServeStale(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSServeStale(input string) (*LocalDNSServeStale, error) {
	vals := map[string]LocalDNSServeStale{
		"disable":   LocalDNSServeStaleDisable,
		"immediate": LocalDNSServeStaleImmediate,
		"verify":    LocalDNSServeStaleVerify,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSServeStale(input)
	return &out, nil
}

type LocalDNSState string

const (
	LocalDNSStateDisabled LocalDNSState = "Disabled"
	LocalDNSStateEnabled  LocalDNSState = "Enabled"
)

func PossibleValuesForLocalDNSState() []string {
	return []string{
		string(LocalDNSStateDisabled),
		string(LocalDNSStateEnabled),
	}
}

func (s *LocalDNSState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalDNSState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalDNSState(input string) (*LocalDNSState, error) {
	vals := map[string]LocalDNSState{
		"disabled": LocalDNSStateDisabled,
		"enabled":  LocalDNSStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalDNSState(input)
	return &out, nil
}

type OSDiskType string

const (
	OSDiskTypeEphemeral OSDiskType = "Ephemeral"
	OSDiskTypeManaged   OSDiskType = "Managed"
)

func PossibleValuesForOSDiskType() []string {
	return []string{
		string(OSDiskTypeEphemeral),
		string(OSDiskTypeManaged),
	}
}

func (s *OSDiskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOSDiskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOSDiskType(input string) (*OSDiskType, error) {
	vals := map[string]OSDiskType{
		"ephemeral": OSDiskTypeEphemeral,
		"managed":   OSDiskTypeManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSDiskType(input)
	return &out, nil
}

type OSSKU string

const (
	OSSKUAzureLinux            OSSKU = "AzureLinux"
	OSSKUAzureLinuxThree       OSSKU = "AzureLinux3"
	OSSKUCBLMariner            OSSKU = "CBLMariner"
	OSSKUUbuntu                OSSKU = "Ubuntu"
	OSSKUUbuntuTwoFourZeroFour OSSKU = "Ubuntu2404"
	OSSKUUbuntuTwoTwoZeroFour  OSSKU = "Ubuntu2204"
	OSSKUWindowsTwoZeroOneNine OSSKU = "Windows2019"
	OSSKUWindowsTwoZeroTwoTwo  OSSKU = "Windows2022"
)

func PossibleValuesForOSSKU() []string {
	return []string{
		string(OSSKUAzureLinux),
		string(OSSKUAzureLinuxThree),
		string(OSSKUCBLMariner),
		string(OSSKUUbuntu),
		string(OSSKUUbuntuTwoFourZeroFour),
		string(OSSKUUbuntuTwoTwoZeroFour),
		string(OSSKUWindowsTwoZeroOneNine),
		string(OSSKUWindowsTwoZeroTwoTwo),
	}
}

func (s *OSSKU) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOSSKU(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOSSKU(input string) (*OSSKU, error) {
	vals := map[string]OSSKU{
		"azurelinux":  OSSKUAzureLinux,
		"azurelinux3": OSSKUAzureLinuxThree,
		"cblmariner":  OSSKUCBLMariner,
		"ubuntu":      OSSKUUbuntu,
		"ubuntu2404":  OSSKUUbuntuTwoFourZeroFour,
		"ubuntu2204":  OSSKUUbuntuTwoTwoZeroFour,
		"windows2019": OSSKUWindowsTwoZeroOneNine,
		"windows2022": OSSKUWindowsTwoZeroTwoTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSSKU(input)
	return &out, nil
}

type OSType string

const (
	OSTypeLinux   OSType = "Linux"
	OSTypeWindows OSType = "Windows"
)

func PossibleValuesForOSType() []string {
	return []string{
		string(OSTypeLinux),
		string(OSTypeWindows),
	}
}

func (s *OSType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOSType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOSType(input string) (*OSType, error) {
	vals := map[string]OSType{
		"linux":   OSTypeLinux,
		"windows": OSTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OSType(input)
	return &out, nil
}

type PodIPAllocationMode string

const (
	PodIPAllocationModeDynamicIndividual PodIPAllocationMode = "DynamicIndividual"
	PodIPAllocationModeStaticBlock       PodIPAllocationMode = "StaticBlock"
)

func PossibleValuesForPodIPAllocationMode() []string {
	return []string{
		string(PodIPAllocationModeDynamicIndividual),
		string(PodIPAllocationModeStaticBlock),
	}
}

func (s *PodIPAllocationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePodIPAllocationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePodIPAllocationMode(input string) (*PodIPAllocationMode, error) {
	vals := map[string]PodIPAllocationMode{
		"dynamicindividual": PodIPAllocationModeDynamicIndividual,
		"staticblock":       PodIPAllocationModeStaticBlock,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PodIPAllocationMode(input)
	return &out, nil
}

type Protocol string

const (
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolTCP),
		string(ProtocolUDP),
	}
}

func (s *Protocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProtocol(input string) (*Protocol, error) {
	vals := map[string]Protocol{
		"tcp": ProtocolTCP,
		"udp": ProtocolUDP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Protocol(input)
	return &out, nil
}

type ScaleDownMode string

const (
	ScaleDownModeDeallocate ScaleDownMode = "Deallocate"
	ScaleDownModeDelete     ScaleDownMode = "Delete"
)

func PossibleValuesForScaleDownMode() []string {
	return []string{
		string(ScaleDownModeDeallocate),
		string(ScaleDownModeDelete),
	}
}

func (s *ScaleDownMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleDownMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleDownMode(input string) (*ScaleDownMode, error) {
	vals := map[string]ScaleDownMode{
		"deallocate": ScaleDownModeDeallocate,
		"delete":     ScaleDownModeDelete,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleDownMode(input)
	return &out, nil
}

type ScaleSetEvictionPolicy string

const (
	ScaleSetEvictionPolicyDeallocate ScaleSetEvictionPolicy = "Deallocate"
	ScaleSetEvictionPolicyDelete     ScaleSetEvictionPolicy = "Delete"
)

func PossibleValuesForScaleSetEvictionPolicy() []string {
	return []string{
		string(ScaleSetEvictionPolicyDeallocate),
		string(ScaleSetEvictionPolicyDelete),
	}
}

func (s *ScaleSetEvictionPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleSetEvictionPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleSetEvictionPolicy(input string) (*ScaleSetEvictionPolicy, error) {
	vals := map[string]ScaleSetEvictionPolicy{
		"deallocate": ScaleSetEvictionPolicyDeallocate,
		"delete":     ScaleSetEvictionPolicyDelete,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleSetEvictionPolicy(input)
	return &out, nil
}

type ScaleSetPriority string

const (
	ScaleSetPriorityRegular ScaleSetPriority = "Regular"
	ScaleSetPrioritySpot    ScaleSetPriority = "Spot"
)

func PossibleValuesForScaleSetPriority() []string {
	return []string{
		string(ScaleSetPriorityRegular),
		string(ScaleSetPrioritySpot),
	}
}

func (s *ScaleSetPriority) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScaleSetPriority(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScaleSetPriority(input string) (*ScaleSetPriority, error) {
	vals := map[string]ScaleSetPriority{
		"regular": ScaleSetPriorityRegular,
		"spot":    ScaleSetPrioritySpot,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScaleSetPriority(input)
	return &out, nil
}

type UndrainableNodeBehavior string

const (
	UndrainableNodeBehaviorCordon   UndrainableNodeBehavior = "Cordon"
	UndrainableNodeBehaviorSchedule UndrainableNodeBehavior = "Schedule"
)

func PossibleValuesForUndrainableNodeBehavior() []string {
	return []string{
		string(UndrainableNodeBehaviorCordon),
		string(UndrainableNodeBehaviorSchedule),
	}
}

func (s *UndrainableNodeBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUndrainableNodeBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUndrainableNodeBehavior(input string) (*UndrainableNodeBehavior, error) {
	vals := map[string]UndrainableNodeBehavior{
		"cordon":   UndrainableNodeBehaviorCordon,
		"schedule": UndrainableNodeBehaviorSchedule,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UndrainableNodeBehavior(input)
	return &out, nil
}

type WorkloadRuntime string

const (
	WorkloadRuntimeKataVMIsolation WorkloadRuntime = "KataVmIsolation"
	WorkloadRuntimeOCIContainer    WorkloadRuntime = "OCIContainer"
	WorkloadRuntimeWasmWasi        WorkloadRuntime = "WasmWasi"
)

func PossibleValuesForWorkloadRuntime() []string {
	return []string{
		string(WorkloadRuntimeKataVMIsolation),
		string(WorkloadRuntimeOCIContainer),
		string(WorkloadRuntimeWasmWasi),
	}
}

func (s *WorkloadRuntime) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWorkloadRuntime(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWorkloadRuntime(input string) (*WorkloadRuntime, error) {
	vals := map[string]WorkloadRuntime{
		"katavmisolation": WorkloadRuntimeKataVMIsolation,
		"ocicontainer":    WorkloadRuntimeOCIContainer,
		"wasmwasi":        WorkloadRuntimeWasmWasi,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadRuntime(input)
	return &out, nil
}
