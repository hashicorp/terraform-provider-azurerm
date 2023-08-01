package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllocationState string

const (
	AllocationStateResizing AllocationState = "Resizing"
	AllocationStateSteady   AllocationState = "Steady"
)

func PossibleValuesForAllocationState() []string {
	return []string{
		string(AllocationStateResizing),
		string(AllocationStateSteady),
	}
}

func (s *AllocationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllocationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllocationState(input string) (*AllocationState, error) {
	vals := map[string]AllocationState{
		"resizing": AllocationStateResizing,
		"steady":   AllocationStateSteady,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllocationState(input)
	return &out, nil
}

type ApplicationSharingPolicy string

const (
	ApplicationSharingPolicyPersonal ApplicationSharingPolicy = "Personal"
	ApplicationSharingPolicyShared   ApplicationSharingPolicy = "Shared"
)

func PossibleValuesForApplicationSharingPolicy() []string {
	return []string{
		string(ApplicationSharingPolicyPersonal),
		string(ApplicationSharingPolicyShared),
	}
}

func (s *ApplicationSharingPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationSharingPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationSharingPolicy(input string) (*ApplicationSharingPolicy, error) {
	vals := map[string]ApplicationSharingPolicy{
		"personal": ApplicationSharingPolicyPersonal,
		"shared":   ApplicationSharingPolicyShared,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationSharingPolicy(input)
	return &out, nil
}

type Autosave string

const (
	AutosaveLocal  Autosave = "Local"
	AutosaveNone   Autosave = "None"
	AutosaveRemote Autosave = "Remote"
)

func PossibleValuesForAutosave() []string {
	return []string{
		string(AutosaveLocal),
		string(AutosaveNone),
		string(AutosaveRemote),
	}
}

func (s *Autosave) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutosave(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutosave(input string) (*Autosave, error) {
	vals := map[string]Autosave{
		"local":  AutosaveLocal,
		"none":   AutosaveNone,
		"remote": AutosaveRemote,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Autosave(input)
	return &out, nil
}

type Caching string

const (
	CachingNone      Caching = "None"
	CachingReadOnly  Caching = "ReadOnly"
	CachingReadWrite Caching = "ReadWrite"
)

func PossibleValuesForCaching() []string {
	return []string{
		string(CachingNone),
		string(CachingReadOnly),
		string(CachingReadWrite),
	}
}

func (s *Caching) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCaching(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCaching(input string) (*Caching, error) {
	vals := map[string]Caching{
		"none":      CachingNone,
		"readonly":  CachingReadOnly,
		"readwrite": CachingReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Caching(input)
	return &out, nil
}

type ClusterPurpose string

const (
	ClusterPurposeDenseProd ClusterPurpose = "DenseProd"
	ClusterPurposeDevTest   ClusterPurpose = "DevTest"
	ClusterPurposeFastProd  ClusterPurpose = "FastProd"
)

func PossibleValuesForClusterPurpose() []string {
	return []string{
		string(ClusterPurposeDenseProd),
		string(ClusterPurposeDevTest),
		string(ClusterPurposeFastProd),
	}
}

func (s *ClusterPurpose) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterPurpose(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterPurpose(input string) (*ClusterPurpose, error) {
	vals := map[string]ClusterPurpose{
		"denseprod": ClusterPurposeDenseProd,
		"devtest":   ClusterPurposeDevTest,
		"fastprod":  ClusterPurposeFastProd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterPurpose(input)
	return &out, nil
}

type ComputeInstanceAuthorizationType string

const (
	ComputeInstanceAuthorizationTypePersonal ComputeInstanceAuthorizationType = "personal"
)

func PossibleValuesForComputeInstanceAuthorizationType() []string {
	return []string{
		string(ComputeInstanceAuthorizationTypePersonal),
	}
}

func (s *ComputeInstanceAuthorizationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeInstanceAuthorizationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeInstanceAuthorizationType(input string) (*ComputeInstanceAuthorizationType, error) {
	vals := map[string]ComputeInstanceAuthorizationType{
		"personal": ComputeInstanceAuthorizationTypePersonal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeInstanceAuthorizationType(input)
	return &out, nil
}

type ComputeInstanceState string

const (
	ComputeInstanceStateCreateFailed    ComputeInstanceState = "CreateFailed"
	ComputeInstanceStateCreating        ComputeInstanceState = "Creating"
	ComputeInstanceStateDeleting        ComputeInstanceState = "Deleting"
	ComputeInstanceStateJobRunning      ComputeInstanceState = "JobRunning"
	ComputeInstanceStateRestarting      ComputeInstanceState = "Restarting"
	ComputeInstanceStateRunning         ComputeInstanceState = "Running"
	ComputeInstanceStateSettingUp       ComputeInstanceState = "SettingUp"
	ComputeInstanceStateSetupFailed     ComputeInstanceState = "SetupFailed"
	ComputeInstanceStateStarting        ComputeInstanceState = "Starting"
	ComputeInstanceStateStopped         ComputeInstanceState = "Stopped"
	ComputeInstanceStateStopping        ComputeInstanceState = "Stopping"
	ComputeInstanceStateUnknown         ComputeInstanceState = "Unknown"
	ComputeInstanceStateUnusable        ComputeInstanceState = "Unusable"
	ComputeInstanceStateUserSettingUp   ComputeInstanceState = "UserSettingUp"
	ComputeInstanceStateUserSetupFailed ComputeInstanceState = "UserSetupFailed"
)

func PossibleValuesForComputeInstanceState() []string {
	return []string{
		string(ComputeInstanceStateCreateFailed),
		string(ComputeInstanceStateCreating),
		string(ComputeInstanceStateDeleting),
		string(ComputeInstanceStateJobRunning),
		string(ComputeInstanceStateRestarting),
		string(ComputeInstanceStateRunning),
		string(ComputeInstanceStateSettingUp),
		string(ComputeInstanceStateSetupFailed),
		string(ComputeInstanceStateStarting),
		string(ComputeInstanceStateStopped),
		string(ComputeInstanceStateStopping),
		string(ComputeInstanceStateUnknown),
		string(ComputeInstanceStateUnusable),
		string(ComputeInstanceStateUserSettingUp),
		string(ComputeInstanceStateUserSetupFailed),
	}
}

func (s *ComputeInstanceState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeInstanceState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeInstanceState(input string) (*ComputeInstanceState, error) {
	vals := map[string]ComputeInstanceState{
		"createfailed":    ComputeInstanceStateCreateFailed,
		"creating":        ComputeInstanceStateCreating,
		"deleting":        ComputeInstanceStateDeleting,
		"jobrunning":      ComputeInstanceStateJobRunning,
		"restarting":      ComputeInstanceStateRestarting,
		"running":         ComputeInstanceStateRunning,
		"settingup":       ComputeInstanceStateSettingUp,
		"setupfailed":     ComputeInstanceStateSetupFailed,
		"starting":        ComputeInstanceStateStarting,
		"stopped":         ComputeInstanceStateStopped,
		"stopping":        ComputeInstanceStateStopping,
		"unknown":         ComputeInstanceStateUnknown,
		"unusable":        ComputeInstanceStateUnusable,
		"usersettingup":   ComputeInstanceStateUserSettingUp,
		"usersetupfailed": ComputeInstanceStateUserSetupFailed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeInstanceState(input)
	return &out, nil
}

type ComputePowerAction string

const (
	ComputePowerActionStart ComputePowerAction = "Start"
	ComputePowerActionStop  ComputePowerAction = "Stop"
)

func PossibleValuesForComputePowerAction() []string {
	return []string{
		string(ComputePowerActionStart),
		string(ComputePowerActionStop),
	}
}

func (s *ComputePowerAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputePowerAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputePowerAction(input string) (*ComputePowerAction, error) {
	vals := map[string]ComputePowerAction{
		"start": ComputePowerActionStart,
		"stop":  ComputePowerActionStop,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputePowerAction(input)
	return &out, nil
}

type ComputeType string

const (
	ComputeTypeAKS               ComputeType = "AKS"
	ComputeTypeAmlCompute        ComputeType = "AmlCompute"
	ComputeTypeComputeInstance   ComputeType = "ComputeInstance"
	ComputeTypeDataFactory       ComputeType = "DataFactory"
	ComputeTypeDataLakeAnalytics ComputeType = "DataLakeAnalytics"
	ComputeTypeDatabricks        ComputeType = "Databricks"
	ComputeTypeHDInsight         ComputeType = "HDInsight"
	ComputeTypeKubernetes        ComputeType = "Kubernetes"
	ComputeTypeSynapseSpark      ComputeType = "SynapseSpark"
	ComputeTypeVirtualMachine    ComputeType = "VirtualMachine"
)

func PossibleValuesForComputeType() []string {
	return []string{
		string(ComputeTypeAKS),
		string(ComputeTypeAmlCompute),
		string(ComputeTypeComputeInstance),
		string(ComputeTypeDataFactory),
		string(ComputeTypeDataLakeAnalytics),
		string(ComputeTypeDatabricks),
		string(ComputeTypeHDInsight),
		string(ComputeTypeKubernetes),
		string(ComputeTypeSynapseSpark),
		string(ComputeTypeVirtualMachine),
	}
}

func (s *ComputeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeType(input string) (*ComputeType, error) {
	vals := map[string]ComputeType{
		"aks":               ComputeTypeAKS,
		"amlcompute":        ComputeTypeAmlCompute,
		"computeinstance":   ComputeTypeComputeInstance,
		"datafactory":       ComputeTypeDataFactory,
		"datalakeanalytics": ComputeTypeDataLakeAnalytics,
		"databricks":        ComputeTypeDatabricks,
		"hdinsight":         ComputeTypeHDInsight,
		"kubernetes":        ComputeTypeKubernetes,
		"synapsespark":      ComputeTypeSynapseSpark,
		"virtualmachine":    ComputeTypeVirtualMachine,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeType(input)
	return &out, nil
}

type EnvironmentVariableType string

const (
	EnvironmentVariableTypeLocal EnvironmentVariableType = "local"
)

func PossibleValuesForEnvironmentVariableType() []string {
	return []string{
		string(EnvironmentVariableTypeLocal),
	}
}

func (s *EnvironmentVariableType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnvironmentVariableType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnvironmentVariableType(input string) (*EnvironmentVariableType, error) {
	vals := map[string]EnvironmentVariableType{
		"local": EnvironmentVariableTypeLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EnvironmentVariableType(input)
	return &out, nil
}

type ImageType string

const (
	ImageTypeAzureml ImageType = "azureml"
	ImageTypeDocker  ImageType = "docker"
)

func PossibleValuesForImageType() []string {
	return []string{
		string(ImageTypeAzureml),
		string(ImageTypeDocker),
	}
}

func (s *ImageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseImageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseImageType(input string) (*ImageType, error) {
	vals := map[string]ImageType{
		"azureml": ImageTypeAzureml,
		"docker":  ImageTypeDocker,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ImageType(input)
	return &out, nil
}

type LoadBalancerType string

const (
	LoadBalancerTypeInternalLoadBalancer LoadBalancerType = "InternalLoadBalancer"
	LoadBalancerTypePublicIP             LoadBalancerType = "PublicIp"
)

func PossibleValuesForLoadBalancerType() []string {
	return []string{
		string(LoadBalancerTypeInternalLoadBalancer),
		string(LoadBalancerTypePublicIP),
	}
}

func (s *LoadBalancerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLoadBalancerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLoadBalancerType(input string) (*LoadBalancerType, error) {
	vals := map[string]LoadBalancerType{
		"internalloadbalancer": LoadBalancerTypeInternalLoadBalancer,
		"publicip":             LoadBalancerTypePublicIP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LoadBalancerType(input)
	return &out, nil
}

type MountAction string

const (
	MountActionMount   MountAction = "Mount"
	MountActionUnmount MountAction = "Unmount"
)

func PossibleValuesForMountAction() []string {
	return []string{
		string(MountActionMount),
		string(MountActionUnmount),
	}
}

func (s *MountAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMountAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMountAction(input string) (*MountAction, error) {
	vals := map[string]MountAction{
		"mount":   MountActionMount,
		"unmount": MountActionUnmount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MountAction(input)
	return &out, nil
}

type MountState string

const (
	MountStateMountFailed      MountState = "MountFailed"
	MountStateMountRequested   MountState = "MountRequested"
	MountStateMounted          MountState = "Mounted"
	MountStateUnmountFailed    MountState = "UnmountFailed"
	MountStateUnmountRequested MountState = "UnmountRequested"
	MountStateUnmounted        MountState = "Unmounted"
)

func PossibleValuesForMountState() []string {
	return []string{
		string(MountStateMountFailed),
		string(MountStateMountRequested),
		string(MountStateMounted),
		string(MountStateUnmountFailed),
		string(MountStateUnmountRequested),
		string(MountStateUnmounted),
	}
}

func (s *MountState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMountState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMountState(input string) (*MountState, error) {
	vals := map[string]MountState{
		"mountfailed":      MountStateMountFailed,
		"mountrequested":   MountStateMountRequested,
		"mounted":          MountStateMounted,
		"unmountfailed":    MountStateUnmountFailed,
		"unmountrequested": MountStateUnmountRequested,
		"unmounted":        MountStateUnmounted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MountState(input)
	return &out, nil
}

type Network string

const (
	NetworkBridge Network = "Bridge"
	NetworkHost   Network = "Host"
)

func PossibleValuesForNetwork() []string {
	return []string{
		string(NetworkBridge),
		string(NetworkHost),
	}
}

func (s *Network) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetwork(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetwork(input string) (*Network, error) {
	vals := map[string]Network{
		"bridge": NetworkBridge,
		"host":   NetworkHost,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Network(input)
	return &out, nil
}

type NodeState string

const (
	NodeStateIdle      NodeState = "idle"
	NodeStateLeaving   NodeState = "leaving"
	NodeStatePreempted NodeState = "preempted"
	NodeStatePreparing NodeState = "preparing"
	NodeStateRunning   NodeState = "running"
	NodeStateUnusable  NodeState = "unusable"
)

func PossibleValuesForNodeState() []string {
	return []string{
		string(NodeStateIdle),
		string(NodeStateLeaving),
		string(NodeStatePreempted),
		string(NodeStatePreparing),
		string(NodeStateRunning),
		string(NodeStateUnusable),
	}
}

func (s *NodeState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNodeState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNodeState(input string) (*NodeState, error) {
	vals := map[string]NodeState{
		"idle":      NodeStateIdle,
		"leaving":   NodeStateLeaving,
		"preempted": NodeStatePreempted,
		"preparing": NodeStatePreparing,
		"running":   NodeStateRunning,
		"unusable":  NodeStateUnusable,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NodeState(input)
	return &out, nil
}

type OperationName string

const (
	OperationNameCreate  OperationName = "Create"
	OperationNameDelete  OperationName = "Delete"
	OperationNameReimage OperationName = "Reimage"
	OperationNameRestart OperationName = "Restart"
	OperationNameStart   OperationName = "Start"
	OperationNameStop    OperationName = "Stop"
)

func PossibleValuesForOperationName() []string {
	return []string{
		string(OperationNameCreate),
		string(OperationNameDelete),
		string(OperationNameReimage),
		string(OperationNameRestart),
		string(OperationNameStart),
		string(OperationNameStop),
	}
}

func (s *OperationName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationName(input string) (*OperationName, error) {
	vals := map[string]OperationName{
		"create":  OperationNameCreate,
		"delete":  OperationNameDelete,
		"reimage": OperationNameReimage,
		"restart": OperationNameRestart,
		"start":   OperationNameStart,
		"stop":    OperationNameStop,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationName(input)
	return &out, nil
}

type OperationStatus string

const (
	OperationStatusCreateFailed  OperationStatus = "CreateFailed"
	OperationStatusDeleteFailed  OperationStatus = "DeleteFailed"
	OperationStatusInProgress    OperationStatus = "InProgress"
	OperationStatusReimageFailed OperationStatus = "ReimageFailed"
	OperationStatusRestartFailed OperationStatus = "RestartFailed"
	OperationStatusStartFailed   OperationStatus = "StartFailed"
	OperationStatusStopFailed    OperationStatus = "StopFailed"
	OperationStatusSucceeded     OperationStatus = "Succeeded"
)

func PossibleValuesForOperationStatus() []string {
	return []string{
		string(OperationStatusCreateFailed),
		string(OperationStatusDeleteFailed),
		string(OperationStatusInProgress),
		string(OperationStatusReimageFailed),
		string(OperationStatusRestartFailed),
		string(OperationStatusStartFailed),
		string(OperationStatusStopFailed),
		string(OperationStatusSucceeded),
	}
}

func (s *OperationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationStatus(input string) (*OperationStatus, error) {
	vals := map[string]OperationStatus{
		"createfailed":  OperationStatusCreateFailed,
		"deletefailed":  OperationStatusDeleteFailed,
		"inprogress":    OperationStatusInProgress,
		"reimagefailed": OperationStatusReimageFailed,
		"restartfailed": OperationStatusRestartFailed,
		"startfailed":   OperationStatusStartFailed,
		"stopfailed":    OperationStatusStopFailed,
		"succeeded":     OperationStatusSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationStatus(input)
	return &out, nil
}

type OperationTrigger string

const (
	OperationTriggerIdleShutdown OperationTrigger = "IdleShutdown"
	OperationTriggerSchedule     OperationTrigger = "Schedule"
	OperationTriggerUser         OperationTrigger = "User"
)

func PossibleValuesForOperationTrigger() []string {
	return []string{
		string(OperationTriggerIdleShutdown),
		string(OperationTriggerSchedule),
		string(OperationTriggerUser),
	}
}

func (s *OperationTrigger) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperationTrigger(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperationTrigger(input string) (*OperationTrigger, error) {
	vals := map[string]OperationTrigger{
		"idleshutdown": OperationTriggerIdleShutdown,
		"schedule":     OperationTriggerSchedule,
		"user":         OperationTriggerUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationTrigger(input)
	return &out, nil
}

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeWindows OsType = "Windows"
)

func PossibleValuesForOsType() []string {
	return []string{
		string(OsTypeLinux),
		string(OsTypeWindows),
	}
}

func (s *OsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsType(input string) (*OsType, error) {
	vals := map[string]OsType{
		"linux":   OsTypeLinux,
		"windows": OsTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsType(input)
	return &out, nil
}

type Protocol string

const (
	ProtocolHTTP Protocol = "http"
	ProtocolTcp  Protocol = "tcp"
	ProtocolUdp  Protocol = "udp"
)

func PossibleValuesForProtocol() []string {
	return []string{
		string(ProtocolHTTP),
		string(ProtocolTcp),
		string(ProtocolUdp),
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
		"http": ProtocolHTTP,
		"tcp":  ProtocolTcp,
		"udp":  ProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Protocol(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
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
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"unknown":   ProvisioningStateUnknown,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ProvisioningStatus string

const (
	ProvisioningStatusCompleted    ProvisioningStatus = "Completed"
	ProvisioningStatusFailed       ProvisioningStatus = "Failed"
	ProvisioningStatusProvisioning ProvisioningStatus = "Provisioning"
)

func PossibleValuesForProvisioningStatus() []string {
	return []string{
		string(ProvisioningStatusCompleted),
		string(ProvisioningStatusFailed),
		string(ProvisioningStatusProvisioning),
	}
}

func (s *ProvisioningStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningStatus(input string) (*ProvisioningStatus, error) {
	vals := map[string]ProvisioningStatus{
		"completed":    ProvisioningStatusCompleted,
		"failed":       ProvisioningStatusFailed,
		"provisioning": ProvisioningStatusProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningStatus(input)
	return &out, nil
}

type RecurrenceFrequency string

const (
	RecurrenceFrequencyDay    RecurrenceFrequency = "Day"
	RecurrenceFrequencyHour   RecurrenceFrequency = "Hour"
	RecurrenceFrequencyMinute RecurrenceFrequency = "Minute"
	RecurrenceFrequencyMonth  RecurrenceFrequency = "Month"
	RecurrenceFrequencyWeek   RecurrenceFrequency = "Week"
)

func PossibleValuesForRecurrenceFrequency() []string {
	return []string{
		string(RecurrenceFrequencyDay),
		string(RecurrenceFrequencyHour),
		string(RecurrenceFrequencyMinute),
		string(RecurrenceFrequencyMonth),
		string(RecurrenceFrequencyWeek),
	}
}

func (s *RecurrenceFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecurrenceFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecurrenceFrequency(input string) (*RecurrenceFrequency, error) {
	vals := map[string]RecurrenceFrequency{
		"day":    RecurrenceFrequencyDay,
		"hour":   RecurrenceFrequencyHour,
		"minute": RecurrenceFrequencyMinute,
		"month":  RecurrenceFrequencyMonth,
		"week":   RecurrenceFrequencyWeek,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecurrenceFrequency(input)
	return &out, nil
}

type RemoteLoginPortPublicAccess string

const (
	RemoteLoginPortPublicAccessDisabled     RemoteLoginPortPublicAccess = "Disabled"
	RemoteLoginPortPublicAccessEnabled      RemoteLoginPortPublicAccess = "Enabled"
	RemoteLoginPortPublicAccessNotSpecified RemoteLoginPortPublicAccess = "NotSpecified"
)

func PossibleValuesForRemoteLoginPortPublicAccess() []string {
	return []string{
		string(RemoteLoginPortPublicAccessDisabled),
		string(RemoteLoginPortPublicAccessEnabled),
		string(RemoteLoginPortPublicAccessNotSpecified),
	}
}

func (s *RemoteLoginPortPublicAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRemoteLoginPortPublicAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRemoteLoginPortPublicAccess(input string) (*RemoteLoginPortPublicAccess, error) {
	vals := map[string]RemoteLoginPortPublicAccess{
		"disabled":     RemoteLoginPortPublicAccessDisabled,
		"enabled":      RemoteLoginPortPublicAccessEnabled,
		"notspecified": RemoteLoginPortPublicAccessNotSpecified,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RemoteLoginPortPublicAccess(input)
	return &out, nil
}

type ScheduleProvisioningState string

const (
	ScheduleProvisioningStateCompleted    ScheduleProvisioningState = "Completed"
	ScheduleProvisioningStateFailed       ScheduleProvisioningState = "Failed"
	ScheduleProvisioningStateProvisioning ScheduleProvisioningState = "Provisioning"
)

func PossibleValuesForScheduleProvisioningState() []string {
	return []string{
		string(ScheduleProvisioningStateCompleted),
		string(ScheduleProvisioningStateFailed),
		string(ScheduleProvisioningStateProvisioning),
	}
}

func (s *ScheduleProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleProvisioningState(input string) (*ScheduleProvisioningState, error) {
	vals := map[string]ScheduleProvisioningState{
		"completed":    ScheduleProvisioningStateCompleted,
		"failed":       ScheduleProvisioningStateFailed,
		"provisioning": ScheduleProvisioningStateProvisioning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleProvisioningState(input)
	return &out, nil
}

type ScheduleStatus string

const (
	ScheduleStatusDisabled ScheduleStatus = "Disabled"
	ScheduleStatusEnabled  ScheduleStatus = "Enabled"
)

func PossibleValuesForScheduleStatus() []string {
	return []string{
		string(ScheduleStatusDisabled),
		string(ScheduleStatusEnabled),
	}
}

func (s *ScheduleStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleStatus(input string) (*ScheduleStatus, error) {
	vals := map[string]ScheduleStatus{
		"disabled": ScheduleStatusDisabled,
		"enabled":  ScheduleStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleStatus(input)
	return &out, nil
}

type SkuTier string

const (
	SkuTierBasic    SkuTier = "Basic"
	SkuTierFree     SkuTier = "Free"
	SkuTierPremium  SkuTier = "Premium"
	SkuTierStandard SkuTier = "Standard"
)

func PossibleValuesForSkuTier() []string {
	return []string{
		string(SkuTierBasic),
		string(SkuTierFree),
		string(SkuTierPremium),
		string(SkuTierStandard),
	}
}

func (s *SkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuTier(input string) (*SkuTier, error) {
	vals := map[string]SkuTier{
		"basic":    SkuTierBasic,
		"free":     SkuTierFree,
		"premium":  SkuTierPremium,
		"standard": SkuTierStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuTier(input)
	return &out, nil
}

type SourceType string

const (
	SourceTypeDataset   SourceType = "Dataset"
	SourceTypeDatastore SourceType = "Datastore"
	SourceTypeURI       SourceType = "URI"
)

func PossibleValuesForSourceType() []string {
	return []string{
		string(SourceTypeDataset),
		string(SourceTypeDatastore),
		string(SourceTypeURI),
	}
}

func (s *SourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSourceType(input string) (*SourceType, error) {
	vals := map[string]SourceType{
		"dataset":   SourceTypeDataset,
		"datastore": SourceTypeDatastore,
		"uri":       SourceTypeURI,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SourceType(input)
	return &out, nil
}

type SshPublicAccess string

const (
	SshPublicAccessDisabled SshPublicAccess = "Disabled"
	SshPublicAccessEnabled  SshPublicAccess = "Enabled"
)

func PossibleValuesForSshPublicAccess() []string {
	return []string{
		string(SshPublicAccessDisabled),
		string(SshPublicAccessEnabled),
	}
}

func (s *SshPublicAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSshPublicAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSshPublicAccess(input string) (*SshPublicAccess, error) {
	vals := map[string]SshPublicAccess{
		"disabled": SshPublicAccessDisabled,
		"enabled":  SshPublicAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SshPublicAccess(input)
	return &out, nil
}

type SslConfigStatus string

const (
	SslConfigStatusAuto     SslConfigStatus = "Auto"
	SslConfigStatusDisabled SslConfigStatus = "Disabled"
	SslConfigStatusEnabled  SslConfigStatus = "Enabled"
)

func PossibleValuesForSslConfigStatus() []string {
	return []string{
		string(SslConfigStatusAuto),
		string(SslConfigStatusDisabled),
		string(SslConfigStatusEnabled),
	}
}

func (s *SslConfigStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSslConfigStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSslConfigStatus(input string) (*SslConfigStatus, error) {
	vals := map[string]SslConfigStatus{
		"auto":     SslConfigStatusAuto,
		"disabled": SslConfigStatusDisabled,
		"enabled":  SslConfigStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SslConfigStatus(input)
	return &out, nil
}

type StorageAccountType string

const (
	StorageAccountTypePremiumLRS  StorageAccountType = "Premium_LRS"
	StorageAccountTypeStandardLRS StorageAccountType = "Standard_LRS"
)

func PossibleValuesForStorageAccountType() []string {
	return []string{
		string(StorageAccountTypePremiumLRS),
		string(StorageAccountTypeStandardLRS),
	}
}

func (s *StorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAccountType(input string) (*StorageAccountType, error) {
	vals := map[string]StorageAccountType{
		"premium_lrs":  StorageAccountTypePremiumLRS,
		"standard_lrs": StorageAccountTypeStandardLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountType(input)
	return &out, nil
}

type TriggerType string

const (
	TriggerTypeCron       TriggerType = "Cron"
	TriggerTypeRecurrence TriggerType = "Recurrence"
)

func PossibleValuesForTriggerType() []string {
	return []string{
		string(TriggerTypeCron),
		string(TriggerTypeRecurrence),
	}
}

func (s *TriggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggerType(input string) (*TriggerType, error) {
	vals := map[string]TriggerType{
		"cron":       TriggerTypeCron,
		"recurrence": TriggerTypeRecurrence,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerType(input)
	return &out, nil
}

type UnderlyingResourceAction string

const (
	UnderlyingResourceActionDelete UnderlyingResourceAction = "Delete"
	UnderlyingResourceActionDetach UnderlyingResourceAction = "Detach"
)

func PossibleValuesForUnderlyingResourceAction() []string {
	return []string{
		string(UnderlyingResourceActionDelete),
		string(UnderlyingResourceActionDetach),
	}
}

func (s *UnderlyingResourceAction) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUnderlyingResourceAction(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUnderlyingResourceAction(input string) (*UnderlyingResourceAction, error) {
	vals := map[string]UnderlyingResourceAction{
		"delete": UnderlyingResourceActionDelete,
		"detach": UnderlyingResourceActionDetach,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UnderlyingResourceAction(input)
	return &out, nil
}

type VMPriority string

const (
	VMPriorityDedicated   VMPriority = "Dedicated"
	VMPriorityLowPriority VMPriority = "LowPriority"
)

func PossibleValuesForVMPriority() []string {
	return []string{
		string(VMPriorityDedicated),
		string(VMPriorityLowPriority),
	}
}

func (s *VMPriority) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMPriority(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMPriority(input string) (*VMPriority, error) {
	vals := map[string]VMPriority{
		"dedicated":   VMPriorityDedicated,
		"lowpriority": VMPriorityLowPriority,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMPriority(input)
	return &out, nil
}

type VolumeDefinitionType string

const (
	VolumeDefinitionTypeBind   VolumeDefinitionType = "bind"
	VolumeDefinitionTypeNpipe  VolumeDefinitionType = "npipe"
	VolumeDefinitionTypeTmpfs  VolumeDefinitionType = "tmpfs"
	VolumeDefinitionTypeVolume VolumeDefinitionType = "volume"
)

func PossibleValuesForVolumeDefinitionType() []string {
	return []string{
		string(VolumeDefinitionTypeBind),
		string(VolumeDefinitionTypeNpipe),
		string(VolumeDefinitionTypeTmpfs),
		string(VolumeDefinitionTypeVolume),
	}
}

func (s *VolumeDefinitionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVolumeDefinitionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVolumeDefinitionType(input string) (*VolumeDefinitionType, error) {
	vals := map[string]VolumeDefinitionType{
		"bind":   VolumeDefinitionTypeBind,
		"npipe":  VolumeDefinitionTypeNpipe,
		"tmpfs":  VolumeDefinitionTypeTmpfs,
		"volume": VolumeDefinitionTypeVolume,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VolumeDefinitionType(input)
	return &out, nil
}

type WeekDay string

const (
	WeekDayFriday    WeekDay = "Friday"
	WeekDayMonday    WeekDay = "Monday"
	WeekDaySaturday  WeekDay = "Saturday"
	WeekDaySunday    WeekDay = "Sunday"
	WeekDayThursday  WeekDay = "Thursday"
	WeekDayTuesday   WeekDay = "Tuesday"
	WeekDayWednesday WeekDay = "Wednesday"
)

func PossibleValuesForWeekDay() []string {
	return []string{
		string(WeekDayFriday),
		string(WeekDayMonday),
		string(WeekDaySaturday),
		string(WeekDaySunday),
		string(WeekDayThursday),
		string(WeekDayTuesday),
		string(WeekDayWednesday),
	}
}

func (s *WeekDay) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWeekDay(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWeekDay(input string) (*WeekDay, error) {
	vals := map[string]WeekDay{
		"friday":    WeekDayFriday,
		"monday":    WeekDayMonday,
		"saturday":  WeekDaySaturday,
		"sunday":    WeekDaySunday,
		"thursday":  WeekDayThursday,
		"tuesday":   WeekDayTuesday,
		"wednesday": WeekDayWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WeekDay(input)
	return &out, nil
}
