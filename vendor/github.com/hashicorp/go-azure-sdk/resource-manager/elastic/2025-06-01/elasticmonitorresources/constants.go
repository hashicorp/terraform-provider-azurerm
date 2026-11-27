package elasticmonitorresources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConfigurationType string

const (
	ConfigurationTypeGeneralPurpose ConfigurationType = "GeneralPurpose"
	ConfigurationTypeNotApplicable  ConfigurationType = "NotApplicable"
	ConfigurationTypeTimeSeries     ConfigurationType = "TimeSeries"
	ConfigurationTypeVector         ConfigurationType = "Vector"
)

func PossibleValuesForConfigurationType() []string {
	return []string{
		string(ConfigurationTypeGeneralPurpose),
		string(ConfigurationTypeNotApplicable),
		string(ConfigurationTypeTimeSeries),
		string(ConfigurationTypeVector),
	}
}

func (s *ConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationType(input string) (*ConfigurationType, error) {
	vals := map[string]ConfigurationType{
		"generalpurpose": ConfigurationTypeGeneralPurpose,
		"notapplicable":  ConfigurationTypeNotApplicable,
		"timeseries":     ConfigurationTypeTimeSeries,
		"vector":         ConfigurationTypeVector,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationType(input)
	return &out, nil
}

type ElasticDeploymentStatus string

const (
	ElasticDeploymentStatusHealthy   ElasticDeploymentStatus = "Healthy"
	ElasticDeploymentStatusUnhealthy ElasticDeploymentStatus = "Unhealthy"
)

func PossibleValuesForElasticDeploymentStatus() []string {
	return []string{
		string(ElasticDeploymentStatusHealthy),
		string(ElasticDeploymentStatusUnhealthy),
	}
}

func (s *ElasticDeploymentStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseElasticDeploymentStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseElasticDeploymentStatus(input string) (*ElasticDeploymentStatus, error) {
	vals := map[string]ElasticDeploymentStatus{
		"healthy":   ElasticDeploymentStatusHealthy,
		"unhealthy": ElasticDeploymentStatusUnhealthy,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ElasticDeploymentStatus(input)
	return &out, nil
}

type HostingType string

const (
	HostingTypeHosted     HostingType = "Hosted"
	HostingTypeServerless HostingType = "Serverless"
)

func PossibleValuesForHostingType() []string {
	return []string{
		string(HostingTypeHosted),
		string(HostingTypeServerless),
	}
}

func (s *HostingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostingType(input string) (*HostingType, error) {
	vals := map[string]HostingType{
		"hosted":     HostingTypeHosted,
		"serverless": HostingTypeServerless,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostingType(input)
	return &out, nil
}

type LiftrResourceCategories string

const (
	LiftrResourceCategoriesMonitorLogs LiftrResourceCategories = "MonitorLogs"
	LiftrResourceCategoriesUnknown     LiftrResourceCategories = "Unknown"
)

func PossibleValuesForLiftrResourceCategories() []string {
	return []string{
		string(LiftrResourceCategoriesMonitorLogs),
		string(LiftrResourceCategoriesUnknown),
	}
}

func (s *LiftrResourceCategories) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLiftrResourceCategories(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLiftrResourceCategories(input string) (*LiftrResourceCategories, error) {
	vals := map[string]LiftrResourceCategories{
		"monitorlogs": LiftrResourceCategoriesMonitorLogs,
		"unknown":     LiftrResourceCategoriesUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LiftrResourceCategories(input)
	return &out, nil
}

type MonitoringStatus string

const (
	MonitoringStatusDisabled MonitoringStatus = "Disabled"
	MonitoringStatusEnabled  MonitoringStatus = "Enabled"
)

func PossibleValuesForMonitoringStatus() []string {
	return []string{
		string(MonitoringStatusDisabled),
		string(MonitoringStatusEnabled),
	}
}

func (s *MonitoringStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitoringStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitoringStatus(input string) (*MonitoringStatus, error) {
	vals := map[string]MonitoringStatus{
		"disabled": MonitoringStatusDisabled,
		"enabled":  MonitoringStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitoringStatus(input)
	return &out, nil
}

type OperationName string

const (
	OperationNameAdd    OperationName = "Add"
	OperationNameDelete OperationName = "Delete"
)

func PossibleValuesForOperationName() []string {
	return []string{
		string(OperationNameAdd),
		string(OperationNameDelete),
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
		"add":    OperationNameAdd,
		"delete": OperationNameDelete,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperationName(input)
	return &out, nil
}

type ProjectType string

const (
	ProjectTypeElasticsearch ProjectType = "Elasticsearch"
	ProjectTypeNotApplicable ProjectType = "NotApplicable"
	ProjectTypeObservability ProjectType = "Observability"
	ProjectTypeSecurity      ProjectType = "Security"
)

func PossibleValuesForProjectType() []string {
	return []string{
		string(ProjectTypeElasticsearch),
		string(ProjectTypeNotApplicable),
		string(ProjectTypeObservability),
		string(ProjectTypeSecurity),
	}
}

func (s *ProjectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProjectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProjectType(input string) (*ProjectType, error) {
	vals := map[string]ProjectType{
		"elasticsearch": ProjectTypeElasticsearch,
		"notapplicable": ProjectTypeNotApplicable,
		"observability": ProjectTypeObservability,
		"security":      ProjectTypeSecurity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProjectType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateCreating     ProvisioningState = "Creating"
	ProvisioningStateDeleted      ProvisioningState = "Deleted"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateNotSpecified ProvisioningState = "NotSpecified"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNotSpecified),
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
		"accepted":     ProvisioningStateAccepted,
		"canceled":     ProvisioningStateCanceled,
		"creating":     ProvisioningStateCreating,
		"deleted":      ProvisioningStateDeleted,
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"notspecified": ProvisioningStateNotSpecified,
		"succeeded":    ProvisioningStateSucceeded,
		"updating":     ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type SendingLogs string

const (
	SendingLogsFalse SendingLogs = "False"
	SendingLogsTrue  SendingLogs = "True"
)

func PossibleValuesForSendingLogs() []string {
	return []string{
		string(SendingLogsFalse),
		string(SendingLogsTrue),
	}
}

func (s *SendingLogs) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSendingLogs(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSendingLogs(input string) (*SendingLogs, error) {
	vals := map[string]SendingLogs{
		"false": SendingLogsFalse,
		"true":  SendingLogsTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SendingLogs(input)
	return &out, nil
}

type Type string

const (
	TypeAzurePrivateEndpoint Type = "azure_private_endpoint"
	TypeIP                   Type = "ip"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeAzurePrivateEndpoint),
		string(TypeIP),
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
		"azure_private_endpoint": TypeAzurePrivateEndpoint,
		"ip":                     TypeIP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}
