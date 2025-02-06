package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiPortalApiTryOutEnabledState string

const (
	ApiPortalApiTryOutEnabledStateDisabled ApiPortalApiTryOutEnabledState = "Disabled"
	ApiPortalApiTryOutEnabledStateEnabled  ApiPortalApiTryOutEnabledState = "Enabled"
)

func PossibleValuesForApiPortalApiTryOutEnabledState() []string {
	return []string{
		string(ApiPortalApiTryOutEnabledStateDisabled),
		string(ApiPortalApiTryOutEnabledStateEnabled),
	}
}

func (s *ApiPortalApiTryOutEnabledState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiPortalApiTryOutEnabledState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiPortalApiTryOutEnabledState(input string) (*ApiPortalApiTryOutEnabledState, error) {
	vals := map[string]ApiPortalApiTryOutEnabledState{
		"disabled": ApiPortalApiTryOutEnabledStateDisabled,
		"enabled":  ApiPortalApiTryOutEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiPortalApiTryOutEnabledState(input)
	return &out, nil
}

type ApiPortalProvisioningState string

const (
	ApiPortalProvisioningStateCreating  ApiPortalProvisioningState = "Creating"
	ApiPortalProvisioningStateDeleting  ApiPortalProvisioningState = "Deleting"
	ApiPortalProvisioningStateFailed    ApiPortalProvisioningState = "Failed"
	ApiPortalProvisioningStateSucceeded ApiPortalProvisioningState = "Succeeded"
	ApiPortalProvisioningStateUpdating  ApiPortalProvisioningState = "Updating"
)

func PossibleValuesForApiPortalProvisioningState() []string {
	return []string{
		string(ApiPortalProvisioningStateCreating),
		string(ApiPortalProvisioningStateDeleting),
		string(ApiPortalProvisioningStateFailed),
		string(ApiPortalProvisioningStateSucceeded),
		string(ApiPortalProvisioningStateUpdating),
	}
}

func (s *ApiPortalProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApiPortalProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApiPortalProvisioningState(input string) (*ApiPortalProvisioningState, error) {
	vals := map[string]ApiPortalProvisioningState{
		"creating":  ApiPortalProvisioningStateCreating,
		"deleting":  ApiPortalProvisioningStateDeleting,
		"failed":    ApiPortalProvisioningStateFailed,
		"succeeded": ApiPortalProvisioningStateSucceeded,
		"updating":  ApiPortalProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApiPortalProvisioningState(input)
	return &out, nil
}

type ApmProvisioningState string

const (
	ApmProvisioningStateCanceled  ApmProvisioningState = "Canceled"
	ApmProvisioningStateCreating  ApmProvisioningState = "Creating"
	ApmProvisioningStateDeleting  ApmProvisioningState = "Deleting"
	ApmProvisioningStateFailed    ApmProvisioningState = "Failed"
	ApmProvisioningStateSucceeded ApmProvisioningState = "Succeeded"
	ApmProvisioningStateUpdating  ApmProvisioningState = "Updating"
)

func PossibleValuesForApmProvisioningState() []string {
	return []string{
		string(ApmProvisioningStateCanceled),
		string(ApmProvisioningStateCreating),
		string(ApmProvisioningStateDeleting),
		string(ApmProvisioningStateFailed),
		string(ApmProvisioningStateSucceeded),
		string(ApmProvisioningStateUpdating),
	}
}

func (s *ApmProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApmProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApmProvisioningState(input string) (*ApmProvisioningState, error) {
	vals := map[string]ApmProvisioningState{
		"canceled":  ApmProvisioningStateCanceled,
		"creating":  ApmProvisioningStateCreating,
		"deleting":  ApmProvisioningStateDeleting,
		"failed":    ApmProvisioningStateFailed,
		"succeeded": ApmProvisioningStateSucceeded,
		"updating":  ApmProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApmProvisioningState(input)
	return &out, nil
}

type ApmType string

const (
	ApmTypeAppDynamics         ApmType = "AppDynamics"
	ApmTypeApplicationInsights ApmType = "ApplicationInsights"
	ApmTypeDynatrace           ApmType = "Dynatrace"
	ApmTypeElasticAPM          ApmType = "ElasticAPM"
	ApmTypeNewRelic            ApmType = "NewRelic"
)

func PossibleValuesForApmType() []string {
	return []string{
		string(ApmTypeAppDynamics),
		string(ApmTypeApplicationInsights),
		string(ApmTypeDynatrace),
		string(ApmTypeElasticAPM),
		string(ApmTypeNewRelic),
	}
}

func (s *ApmType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApmType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApmType(input string) (*ApmType, error) {
	vals := map[string]ApmType{
		"appdynamics":         ApmTypeAppDynamics,
		"applicationinsights": ApmTypeApplicationInsights,
		"dynatrace":           ApmTypeDynatrace,
		"elasticapm":          ApmTypeElasticAPM,
		"newrelic":            ApmTypeNewRelic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApmType(input)
	return &out, nil
}

type AppResourceProvisioningState string

const (
	AppResourceProvisioningStateCreating  AppResourceProvisioningState = "Creating"
	AppResourceProvisioningStateDeleting  AppResourceProvisioningState = "Deleting"
	AppResourceProvisioningStateFailed    AppResourceProvisioningState = "Failed"
	AppResourceProvisioningStateSucceeded AppResourceProvisioningState = "Succeeded"
	AppResourceProvisioningStateUpdating  AppResourceProvisioningState = "Updating"
)

func PossibleValuesForAppResourceProvisioningState() []string {
	return []string{
		string(AppResourceProvisioningStateCreating),
		string(AppResourceProvisioningStateDeleting),
		string(AppResourceProvisioningStateFailed),
		string(AppResourceProvisioningStateSucceeded),
		string(AppResourceProvisioningStateUpdating),
	}
}

func (s *AppResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAppResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAppResourceProvisioningState(input string) (*AppResourceProvisioningState, error) {
	vals := map[string]AppResourceProvisioningState{
		"creating":  AppResourceProvisioningStateCreating,
		"deleting":  AppResourceProvisioningStateDeleting,
		"failed":    AppResourceProvisioningStateFailed,
		"succeeded": AppResourceProvisioningStateSucceeded,
		"updating":  AppResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AppResourceProvisioningState(input)
	return &out, nil
}

type ApplicationAcceleratorProvisioningState string

const (
	ApplicationAcceleratorProvisioningStateCanceled  ApplicationAcceleratorProvisioningState = "Canceled"
	ApplicationAcceleratorProvisioningStateCreating  ApplicationAcceleratorProvisioningState = "Creating"
	ApplicationAcceleratorProvisioningStateDeleting  ApplicationAcceleratorProvisioningState = "Deleting"
	ApplicationAcceleratorProvisioningStateFailed    ApplicationAcceleratorProvisioningState = "Failed"
	ApplicationAcceleratorProvisioningStateSucceeded ApplicationAcceleratorProvisioningState = "Succeeded"
	ApplicationAcceleratorProvisioningStateUpdating  ApplicationAcceleratorProvisioningState = "Updating"
)

func PossibleValuesForApplicationAcceleratorProvisioningState() []string {
	return []string{
		string(ApplicationAcceleratorProvisioningStateCanceled),
		string(ApplicationAcceleratorProvisioningStateCreating),
		string(ApplicationAcceleratorProvisioningStateDeleting),
		string(ApplicationAcceleratorProvisioningStateFailed),
		string(ApplicationAcceleratorProvisioningStateSucceeded),
		string(ApplicationAcceleratorProvisioningStateUpdating),
	}
}

func (s *ApplicationAcceleratorProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationAcceleratorProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationAcceleratorProvisioningState(input string) (*ApplicationAcceleratorProvisioningState, error) {
	vals := map[string]ApplicationAcceleratorProvisioningState{
		"canceled":  ApplicationAcceleratorProvisioningStateCanceled,
		"creating":  ApplicationAcceleratorProvisioningStateCreating,
		"deleting":  ApplicationAcceleratorProvisioningStateDeleting,
		"failed":    ApplicationAcceleratorProvisioningStateFailed,
		"succeeded": ApplicationAcceleratorProvisioningStateSucceeded,
		"updating":  ApplicationAcceleratorProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationAcceleratorProvisioningState(input)
	return &out, nil
}

type ApplicationLiveViewProvisioningState string

const (
	ApplicationLiveViewProvisioningStateCanceled  ApplicationLiveViewProvisioningState = "Canceled"
	ApplicationLiveViewProvisioningStateCreating  ApplicationLiveViewProvisioningState = "Creating"
	ApplicationLiveViewProvisioningStateDeleting  ApplicationLiveViewProvisioningState = "Deleting"
	ApplicationLiveViewProvisioningStateFailed    ApplicationLiveViewProvisioningState = "Failed"
	ApplicationLiveViewProvisioningStateSucceeded ApplicationLiveViewProvisioningState = "Succeeded"
	ApplicationLiveViewProvisioningStateUpdating  ApplicationLiveViewProvisioningState = "Updating"
)

func PossibleValuesForApplicationLiveViewProvisioningState() []string {
	return []string{
		string(ApplicationLiveViewProvisioningStateCanceled),
		string(ApplicationLiveViewProvisioningStateCreating),
		string(ApplicationLiveViewProvisioningStateDeleting),
		string(ApplicationLiveViewProvisioningStateFailed),
		string(ApplicationLiveViewProvisioningStateSucceeded),
		string(ApplicationLiveViewProvisioningStateUpdating),
	}
}

func (s *ApplicationLiveViewProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseApplicationLiveViewProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseApplicationLiveViewProvisioningState(input string) (*ApplicationLiveViewProvisioningState, error) {
	vals := map[string]ApplicationLiveViewProvisioningState{
		"canceled":  ApplicationLiveViewProvisioningStateCanceled,
		"creating":  ApplicationLiveViewProvisioningStateCreating,
		"deleting":  ApplicationLiveViewProvisioningStateDeleting,
		"failed":    ApplicationLiveViewProvisioningStateFailed,
		"succeeded": ApplicationLiveViewProvisioningStateSucceeded,
		"updating":  ApplicationLiveViewProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ApplicationLiveViewProvisioningState(input)
	return &out, nil
}

type BackendProtocol string

const (
	BackendProtocolDefault BackendProtocol = "Default"
	BackendProtocolGRPC    BackendProtocol = "GRPC"
)

func PossibleValuesForBackendProtocol() []string {
	return []string{
		string(BackendProtocolDefault),
		string(BackendProtocolGRPC),
	}
}

func (s *BackendProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackendProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackendProtocol(input string) (*BackendProtocol, error) {
	vals := map[string]BackendProtocol{
		"default": BackendProtocolDefault,
		"grpc":    BackendProtocolGRPC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackendProtocol(input)
	return &out, nil
}

type BindingType string

const (
	BindingTypeApacheSkyWalking    BindingType = "ApacheSkyWalking"
	BindingTypeAppDynamics         BindingType = "AppDynamics"
	BindingTypeApplicationInsights BindingType = "ApplicationInsights"
	BindingTypeCACertificates      BindingType = "CACertificates"
	BindingTypeDynatrace           BindingType = "Dynatrace"
	BindingTypeElasticAPM          BindingType = "ElasticAPM"
	BindingTypeNewRelic            BindingType = "NewRelic"
)

func PossibleValuesForBindingType() []string {
	return []string{
		string(BindingTypeApacheSkyWalking),
		string(BindingTypeAppDynamics),
		string(BindingTypeApplicationInsights),
		string(BindingTypeCACertificates),
		string(BindingTypeDynatrace),
		string(BindingTypeElasticAPM),
		string(BindingTypeNewRelic),
	}
}

func (s *BindingType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBindingType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBindingType(input string) (*BindingType, error) {
	vals := map[string]BindingType{
		"apacheskywalking":    BindingTypeApacheSkyWalking,
		"appdynamics":         BindingTypeAppDynamics,
		"applicationinsights": BindingTypeApplicationInsights,
		"cacertificates":      BindingTypeCACertificates,
		"dynatrace":           BindingTypeDynatrace,
		"elasticapm":          BindingTypeElasticAPM,
		"newrelic":            BindingTypeNewRelic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BindingType(input)
	return &out, nil
}

type BuildProvisioningState string

const (
	BuildProvisioningStateCreating  BuildProvisioningState = "Creating"
	BuildProvisioningStateDeleting  BuildProvisioningState = "Deleting"
	BuildProvisioningStateFailed    BuildProvisioningState = "Failed"
	BuildProvisioningStateSucceeded BuildProvisioningState = "Succeeded"
	BuildProvisioningStateUpdating  BuildProvisioningState = "Updating"
)

func PossibleValuesForBuildProvisioningState() []string {
	return []string{
		string(BuildProvisioningStateCreating),
		string(BuildProvisioningStateDeleting),
		string(BuildProvisioningStateFailed),
		string(BuildProvisioningStateSucceeded),
		string(BuildProvisioningStateUpdating),
	}
}

func (s *BuildProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBuildProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBuildProvisioningState(input string) (*BuildProvisioningState, error) {
	vals := map[string]BuildProvisioningState{
		"creating":  BuildProvisioningStateCreating,
		"deleting":  BuildProvisioningStateDeleting,
		"failed":    BuildProvisioningStateFailed,
		"succeeded": BuildProvisioningStateSucceeded,
		"updating":  BuildProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BuildProvisioningState(input)
	return &out, nil
}

type BuildResultProvisioningState string

const (
	BuildResultProvisioningStateBuilding  BuildResultProvisioningState = "Building"
	BuildResultProvisioningStateDeleting  BuildResultProvisioningState = "Deleting"
	BuildResultProvisioningStateFailed    BuildResultProvisioningState = "Failed"
	BuildResultProvisioningStateQueuing   BuildResultProvisioningState = "Queuing"
	BuildResultProvisioningStateSucceeded BuildResultProvisioningState = "Succeeded"
)

func PossibleValuesForBuildResultProvisioningState() []string {
	return []string{
		string(BuildResultProvisioningStateBuilding),
		string(BuildResultProvisioningStateDeleting),
		string(BuildResultProvisioningStateFailed),
		string(BuildResultProvisioningStateQueuing),
		string(BuildResultProvisioningStateSucceeded),
	}
}

func (s *BuildResultProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBuildResultProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBuildResultProvisioningState(input string) (*BuildResultProvisioningState, error) {
	vals := map[string]BuildResultProvisioningState{
		"building":  BuildResultProvisioningStateBuilding,
		"deleting":  BuildResultProvisioningStateDeleting,
		"failed":    BuildResultProvisioningStateFailed,
		"queuing":   BuildResultProvisioningStateQueuing,
		"succeeded": BuildResultProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BuildResultProvisioningState(input)
	return &out, nil
}

type BuildServiceProvisioningState string

const (
	BuildServiceProvisioningStateCreating  BuildServiceProvisioningState = "Creating"
	BuildServiceProvisioningStateDeleting  BuildServiceProvisioningState = "Deleting"
	BuildServiceProvisioningStateFailed    BuildServiceProvisioningState = "Failed"
	BuildServiceProvisioningStateSucceeded BuildServiceProvisioningState = "Succeeded"
	BuildServiceProvisioningStateUpdating  BuildServiceProvisioningState = "Updating"
)

func PossibleValuesForBuildServiceProvisioningState() []string {
	return []string{
		string(BuildServiceProvisioningStateCreating),
		string(BuildServiceProvisioningStateDeleting),
		string(BuildServiceProvisioningStateFailed),
		string(BuildServiceProvisioningStateSucceeded),
		string(BuildServiceProvisioningStateUpdating),
	}
}

func (s *BuildServiceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBuildServiceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBuildServiceProvisioningState(input string) (*BuildServiceProvisioningState, error) {
	vals := map[string]BuildServiceProvisioningState{
		"creating":  BuildServiceProvisioningStateCreating,
		"deleting":  BuildServiceProvisioningStateDeleting,
		"failed":    BuildServiceProvisioningStateFailed,
		"succeeded": BuildServiceProvisioningStateSucceeded,
		"updating":  BuildServiceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BuildServiceProvisioningState(input)
	return &out, nil
}

type BuilderProvisioningState string

const (
	BuilderProvisioningStateCreating  BuilderProvisioningState = "Creating"
	BuilderProvisioningStateDeleting  BuilderProvisioningState = "Deleting"
	BuilderProvisioningStateFailed    BuilderProvisioningState = "Failed"
	BuilderProvisioningStateSucceeded BuilderProvisioningState = "Succeeded"
	BuilderProvisioningStateUpdating  BuilderProvisioningState = "Updating"
)

func PossibleValuesForBuilderProvisioningState() []string {
	return []string{
		string(BuilderProvisioningStateCreating),
		string(BuilderProvisioningStateDeleting),
		string(BuilderProvisioningStateFailed),
		string(BuilderProvisioningStateSucceeded),
		string(BuilderProvisioningStateUpdating),
	}
}

func (s *BuilderProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBuilderProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBuilderProvisioningState(input string) (*BuilderProvisioningState, error) {
	vals := map[string]BuilderProvisioningState{
		"creating":  BuilderProvisioningStateCreating,
		"deleting":  BuilderProvisioningStateDeleting,
		"failed":    BuilderProvisioningStateFailed,
		"succeeded": BuilderProvisioningStateSucceeded,
		"updating":  BuilderProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BuilderProvisioningState(input)
	return &out, nil
}

type BuildpackBindingProvisioningState string

const (
	BuildpackBindingProvisioningStateCreating  BuildpackBindingProvisioningState = "Creating"
	BuildpackBindingProvisioningStateDeleting  BuildpackBindingProvisioningState = "Deleting"
	BuildpackBindingProvisioningStateFailed    BuildpackBindingProvisioningState = "Failed"
	BuildpackBindingProvisioningStateSucceeded BuildpackBindingProvisioningState = "Succeeded"
	BuildpackBindingProvisioningStateUpdating  BuildpackBindingProvisioningState = "Updating"
)

func PossibleValuesForBuildpackBindingProvisioningState() []string {
	return []string{
		string(BuildpackBindingProvisioningStateCreating),
		string(BuildpackBindingProvisioningStateDeleting),
		string(BuildpackBindingProvisioningStateFailed),
		string(BuildpackBindingProvisioningStateSucceeded),
		string(BuildpackBindingProvisioningStateUpdating),
	}
}

func (s *BuildpackBindingProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBuildpackBindingProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBuildpackBindingProvisioningState(input string) (*BuildpackBindingProvisioningState, error) {
	vals := map[string]BuildpackBindingProvisioningState{
		"creating":  BuildpackBindingProvisioningStateCreating,
		"deleting":  BuildpackBindingProvisioningStateDeleting,
		"failed":    BuildpackBindingProvisioningStateFailed,
		"succeeded": BuildpackBindingProvisioningStateSucceeded,
		"updating":  BuildpackBindingProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BuildpackBindingProvisioningState(input)
	return &out, nil
}

type CertificateResourceProvisioningState string

const (
	CertificateResourceProvisioningStateCreating  CertificateResourceProvisioningState = "Creating"
	CertificateResourceProvisioningStateDeleting  CertificateResourceProvisioningState = "Deleting"
	CertificateResourceProvisioningStateFailed    CertificateResourceProvisioningState = "Failed"
	CertificateResourceProvisioningStateSucceeded CertificateResourceProvisioningState = "Succeeded"
	CertificateResourceProvisioningStateUpdating  CertificateResourceProvisioningState = "Updating"
)

func PossibleValuesForCertificateResourceProvisioningState() []string {
	return []string{
		string(CertificateResourceProvisioningStateCreating),
		string(CertificateResourceProvisioningStateDeleting),
		string(CertificateResourceProvisioningStateFailed),
		string(CertificateResourceProvisioningStateSucceeded),
		string(CertificateResourceProvisioningStateUpdating),
	}
}

func (s *CertificateResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateResourceProvisioningState(input string) (*CertificateResourceProvisioningState, error) {
	vals := map[string]CertificateResourceProvisioningState{
		"creating":  CertificateResourceProvisioningStateCreating,
		"deleting":  CertificateResourceProvisioningStateDeleting,
		"failed":    CertificateResourceProvisioningStateFailed,
		"succeeded": CertificateResourceProvisioningStateSucceeded,
		"updating":  CertificateResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateResourceProvisioningState(input)
	return &out, nil
}

type ConfigServerEnabledState string

const (
	ConfigServerEnabledStateDisabled ConfigServerEnabledState = "Disabled"
	ConfigServerEnabledStateEnabled  ConfigServerEnabledState = "Enabled"
)

func PossibleValuesForConfigServerEnabledState() []string {
	return []string{
		string(ConfigServerEnabledStateDisabled),
		string(ConfigServerEnabledStateEnabled),
	}
}

func (s *ConfigServerEnabledState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigServerEnabledState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigServerEnabledState(input string) (*ConfigServerEnabledState, error) {
	vals := map[string]ConfigServerEnabledState{
		"disabled": ConfigServerEnabledStateDisabled,
		"enabled":  ConfigServerEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigServerEnabledState(input)
	return &out, nil
}

type ConfigServerState string

const (
	ConfigServerStateDeleted      ConfigServerState = "Deleted"
	ConfigServerStateFailed       ConfigServerState = "Failed"
	ConfigServerStateNotAvailable ConfigServerState = "NotAvailable"
	ConfigServerStateSucceeded    ConfigServerState = "Succeeded"
	ConfigServerStateUpdating     ConfigServerState = "Updating"
)

func PossibleValuesForConfigServerState() []string {
	return []string{
		string(ConfigServerStateDeleted),
		string(ConfigServerStateFailed),
		string(ConfigServerStateNotAvailable),
		string(ConfigServerStateSucceeded),
		string(ConfigServerStateUpdating),
	}
}

func (s *ConfigServerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigServerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigServerState(input string) (*ConfigServerState, error) {
	vals := map[string]ConfigServerState{
		"deleted":      ConfigServerStateDeleted,
		"failed":       ConfigServerStateFailed,
		"notavailable": ConfigServerStateNotAvailable,
		"succeeded":    ConfigServerStateSucceeded,
		"updating":     ConfigServerStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigServerState(input)
	return &out, nil
}

type ConfigurationServiceGeneration string

const (
	ConfigurationServiceGenerationGenOne ConfigurationServiceGeneration = "Gen1"
	ConfigurationServiceGenerationGenTwo ConfigurationServiceGeneration = "Gen2"
)

func PossibleValuesForConfigurationServiceGeneration() []string {
	return []string{
		string(ConfigurationServiceGenerationGenOne),
		string(ConfigurationServiceGenerationGenTwo),
	}
}

func (s *ConfigurationServiceGeneration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationServiceGeneration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationServiceGeneration(input string) (*ConfigurationServiceGeneration, error) {
	vals := map[string]ConfigurationServiceGeneration{
		"gen1": ConfigurationServiceGenerationGenOne,
		"gen2": ConfigurationServiceGenerationGenTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationServiceGeneration(input)
	return &out, nil
}

type ConfigurationServiceProvisioningState string

const (
	ConfigurationServiceProvisioningStateCreating  ConfigurationServiceProvisioningState = "Creating"
	ConfigurationServiceProvisioningStateDeleting  ConfigurationServiceProvisioningState = "Deleting"
	ConfigurationServiceProvisioningStateFailed    ConfigurationServiceProvisioningState = "Failed"
	ConfigurationServiceProvisioningStateSucceeded ConfigurationServiceProvisioningState = "Succeeded"
	ConfigurationServiceProvisioningStateUpdating  ConfigurationServiceProvisioningState = "Updating"
)

func PossibleValuesForConfigurationServiceProvisioningState() []string {
	return []string{
		string(ConfigurationServiceProvisioningStateCreating),
		string(ConfigurationServiceProvisioningStateDeleting),
		string(ConfigurationServiceProvisioningStateFailed),
		string(ConfigurationServiceProvisioningStateSucceeded),
		string(ConfigurationServiceProvisioningStateUpdating),
	}
}

func (s *ConfigurationServiceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfigurationServiceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfigurationServiceProvisioningState(input string) (*ConfigurationServiceProvisioningState, error) {
	vals := map[string]ConfigurationServiceProvisioningState{
		"creating":  ConfigurationServiceProvisioningStateCreating,
		"deleting":  ConfigurationServiceProvisioningStateDeleting,
		"failed":    ConfigurationServiceProvisioningStateFailed,
		"succeeded": ConfigurationServiceProvisioningStateSucceeded,
		"updating":  ConfigurationServiceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfigurationServiceProvisioningState(input)
	return &out, nil
}

type ContainerRegistryProvisioningState string

const (
	ContainerRegistryProvisioningStateCanceled  ContainerRegistryProvisioningState = "Canceled"
	ContainerRegistryProvisioningStateCreating  ContainerRegistryProvisioningState = "Creating"
	ContainerRegistryProvisioningStateDeleting  ContainerRegistryProvisioningState = "Deleting"
	ContainerRegistryProvisioningStateFailed    ContainerRegistryProvisioningState = "Failed"
	ContainerRegistryProvisioningStateSucceeded ContainerRegistryProvisioningState = "Succeeded"
	ContainerRegistryProvisioningStateUpdating  ContainerRegistryProvisioningState = "Updating"
)

func PossibleValuesForContainerRegistryProvisioningState() []string {
	return []string{
		string(ContainerRegistryProvisioningStateCanceled),
		string(ContainerRegistryProvisioningStateCreating),
		string(ContainerRegistryProvisioningStateDeleting),
		string(ContainerRegistryProvisioningStateFailed),
		string(ContainerRegistryProvisioningStateSucceeded),
		string(ContainerRegistryProvisioningStateUpdating),
	}
}

func (s *ContainerRegistryProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContainerRegistryProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContainerRegistryProvisioningState(input string) (*ContainerRegistryProvisioningState, error) {
	vals := map[string]ContainerRegistryProvisioningState{
		"canceled":  ContainerRegistryProvisioningStateCanceled,
		"creating":  ContainerRegistryProvisioningStateCreating,
		"deleting":  ContainerRegistryProvisioningStateDeleting,
		"failed":    ContainerRegistryProvisioningStateFailed,
		"succeeded": ContainerRegistryProvisioningStateSucceeded,
		"updating":  ContainerRegistryProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContainerRegistryProvisioningState(input)
	return &out, nil
}

type CustomDomainResourceProvisioningState string

const (
	CustomDomainResourceProvisioningStateCreating  CustomDomainResourceProvisioningState = "Creating"
	CustomDomainResourceProvisioningStateDeleting  CustomDomainResourceProvisioningState = "Deleting"
	CustomDomainResourceProvisioningStateFailed    CustomDomainResourceProvisioningState = "Failed"
	CustomDomainResourceProvisioningStateSucceeded CustomDomainResourceProvisioningState = "Succeeded"
	CustomDomainResourceProvisioningStateUpdating  CustomDomainResourceProvisioningState = "Updating"
)

func PossibleValuesForCustomDomainResourceProvisioningState() []string {
	return []string{
		string(CustomDomainResourceProvisioningStateCreating),
		string(CustomDomainResourceProvisioningStateDeleting),
		string(CustomDomainResourceProvisioningStateFailed),
		string(CustomDomainResourceProvisioningStateSucceeded),
		string(CustomDomainResourceProvisioningStateUpdating),
	}
}

func (s *CustomDomainResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomDomainResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomDomainResourceProvisioningState(input string) (*CustomDomainResourceProvisioningState, error) {
	vals := map[string]CustomDomainResourceProvisioningState{
		"creating":  CustomDomainResourceProvisioningStateCreating,
		"deleting":  CustomDomainResourceProvisioningStateDeleting,
		"failed":    CustomDomainResourceProvisioningStateFailed,
		"succeeded": CustomDomainResourceProvisioningStateSucceeded,
		"updating":  CustomDomainResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomDomainResourceProvisioningState(input)
	return &out, nil
}

type CustomizedAcceleratorProvisioningState string

const (
	CustomizedAcceleratorProvisioningStateCanceled  CustomizedAcceleratorProvisioningState = "Canceled"
	CustomizedAcceleratorProvisioningStateCreating  CustomizedAcceleratorProvisioningState = "Creating"
	CustomizedAcceleratorProvisioningStateDeleting  CustomizedAcceleratorProvisioningState = "Deleting"
	CustomizedAcceleratorProvisioningStateFailed    CustomizedAcceleratorProvisioningState = "Failed"
	CustomizedAcceleratorProvisioningStateSucceeded CustomizedAcceleratorProvisioningState = "Succeeded"
	CustomizedAcceleratorProvisioningStateUpdating  CustomizedAcceleratorProvisioningState = "Updating"
)

func PossibleValuesForCustomizedAcceleratorProvisioningState() []string {
	return []string{
		string(CustomizedAcceleratorProvisioningStateCanceled),
		string(CustomizedAcceleratorProvisioningStateCreating),
		string(CustomizedAcceleratorProvisioningStateDeleting),
		string(CustomizedAcceleratorProvisioningStateFailed),
		string(CustomizedAcceleratorProvisioningStateSucceeded),
		string(CustomizedAcceleratorProvisioningStateUpdating),
	}
}

func (s *CustomizedAcceleratorProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomizedAcceleratorProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomizedAcceleratorProvisioningState(input string) (*CustomizedAcceleratorProvisioningState, error) {
	vals := map[string]CustomizedAcceleratorProvisioningState{
		"canceled":  CustomizedAcceleratorProvisioningStateCanceled,
		"creating":  CustomizedAcceleratorProvisioningStateCreating,
		"deleting":  CustomizedAcceleratorProvisioningStateDeleting,
		"failed":    CustomizedAcceleratorProvisioningStateFailed,
		"succeeded": CustomizedAcceleratorProvisioningStateSucceeded,
		"updating":  CustomizedAcceleratorProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomizedAcceleratorProvisioningState(input)
	return &out, nil
}

type CustomizedAcceleratorType string

const (
	CustomizedAcceleratorTypeAccelerator CustomizedAcceleratorType = "Accelerator"
	CustomizedAcceleratorTypeFragment    CustomizedAcceleratorType = "Fragment"
)

func PossibleValuesForCustomizedAcceleratorType() []string {
	return []string{
		string(CustomizedAcceleratorTypeAccelerator),
		string(CustomizedAcceleratorTypeFragment),
	}
}

func (s *CustomizedAcceleratorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomizedAcceleratorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomizedAcceleratorType(input string) (*CustomizedAcceleratorType, error) {
	vals := map[string]CustomizedAcceleratorType{
		"accelerator": CustomizedAcceleratorTypeAccelerator,
		"fragment":    CustomizedAcceleratorTypeFragment,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomizedAcceleratorType(input)
	return &out, nil
}

type CustomizedAcceleratorValidateResultState string

const (
	CustomizedAcceleratorValidateResultStateInvalid CustomizedAcceleratorValidateResultState = "Invalid"
	CustomizedAcceleratorValidateResultStateValid   CustomizedAcceleratorValidateResultState = "Valid"
)

func PossibleValuesForCustomizedAcceleratorValidateResultState() []string {
	return []string{
		string(CustomizedAcceleratorValidateResultStateInvalid),
		string(CustomizedAcceleratorValidateResultStateValid),
	}
}

func (s *CustomizedAcceleratorValidateResultState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomizedAcceleratorValidateResultState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomizedAcceleratorValidateResultState(input string) (*CustomizedAcceleratorValidateResultState, error) {
	vals := map[string]CustomizedAcceleratorValidateResultState{
		"invalid": CustomizedAcceleratorValidateResultStateInvalid,
		"valid":   CustomizedAcceleratorValidateResultStateValid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomizedAcceleratorValidateResultState(input)
	return &out, nil
}

type DeploymentResourceProvisioningState string

const (
	DeploymentResourceProvisioningStateCreating  DeploymentResourceProvisioningState = "Creating"
	DeploymentResourceProvisioningStateFailed    DeploymentResourceProvisioningState = "Failed"
	DeploymentResourceProvisioningStateSucceeded DeploymentResourceProvisioningState = "Succeeded"
	DeploymentResourceProvisioningStateUpdating  DeploymentResourceProvisioningState = "Updating"
)

func PossibleValuesForDeploymentResourceProvisioningState() []string {
	return []string{
		string(DeploymentResourceProvisioningStateCreating),
		string(DeploymentResourceProvisioningStateFailed),
		string(DeploymentResourceProvisioningStateSucceeded),
		string(DeploymentResourceProvisioningStateUpdating),
	}
}

func (s *DeploymentResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentResourceProvisioningState(input string) (*DeploymentResourceProvisioningState, error) {
	vals := map[string]DeploymentResourceProvisioningState{
		"creating":  DeploymentResourceProvisioningStateCreating,
		"failed":    DeploymentResourceProvisioningStateFailed,
		"succeeded": DeploymentResourceProvisioningStateSucceeded,
		"updating":  DeploymentResourceProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentResourceProvisioningState(input)
	return &out, nil
}

type DeploymentResourceStatus string

const (
	DeploymentResourceStatusRunning DeploymentResourceStatus = "Running"
	DeploymentResourceStatusStopped DeploymentResourceStatus = "Stopped"
)

func PossibleValuesForDeploymentResourceStatus() []string {
	return []string{
		string(DeploymentResourceStatusRunning),
		string(DeploymentResourceStatusStopped),
	}
}

func (s *DeploymentResourceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeploymentResourceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeploymentResourceStatus(input string) (*DeploymentResourceStatus, error) {
	vals := map[string]DeploymentResourceStatus{
		"running": DeploymentResourceStatusRunning,
		"stopped": DeploymentResourceStatusStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeploymentResourceStatus(input)
	return &out, nil
}

type DevToolPortalFeatureState string

const (
	DevToolPortalFeatureStateDisabled DevToolPortalFeatureState = "Disabled"
	DevToolPortalFeatureStateEnabled  DevToolPortalFeatureState = "Enabled"
)

func PossibleValuesForDevToolPortalFeatureState() []string {
	return []string{
		string(DevToolPortalFeatureStateDisabled),
		string(DevToolPortalFeatureStateEnabled),
	}
}

func (s *DevToolPortalFeatureState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDevToolPortalFeatureState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDevToolPortalFeatureState(input string) (*DevToolPortalFeatureState, error) {
	vals := map[string]DevToolPortalFeatureState{
		"disabled": DevToolPortalFeatureStateDisabled,
		"enabled":  DevToolPortalFeatureStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DevToolPortalFeatureState(input)
	return &out, nil
}

type DevToolPortalProvisioningState string

const (
	DevToolPortalProvisioningStateCanceled  DevToolPortalProvisioningState = "Canceled"
	DevToolPortalProvisioningStateCreating  DevToolPortalProvisioningState = "Creating"
	DevToolPortalProvisioningStateDeleting  DevToolPortalProvisioningState = "Deleting"
	DevToolPortalProvisioningStateFailed    DevToolPortalProvisioningState = "Failed"
	DevToolPortalProvisioningStateSucceeded DevToolPortalProvisioningState = "Succeeded"
	DevToolPortalProvisioningStateUpdating  DevToolPortalProvisioningState = "Updating"
)

func PossibleValuesForDevToolPortalProvisioningState() []string {
	return []string{
		string(DevToolPortalProvisioningStateCanceled),
		string(DevToolPortalProvisioningStateCreating),
		string(DevToolPortalProvisioningStateDeleting),
		string(DevToolPortalProvisioningStateFailed),
		string(DevToolPortalProvisioningStateSucceeded),
		string(DevToolPortalProvisioningStateUpdating),
	}
}

func (s *DevToolPortalProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDevToolPortalProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDevToolPortalProvisioningState(input string) (*DevToolPortalProvisioningState, error) {
	vals := map[string]DevToolPortalProvisioningState{
		"canceled":  DevToolPortalProvisioningStateCanceled,
		"creating":  DevToolPortalProvisioningStateCreating,
		"deleting":  DevToolPortalProvisioningStateDeleting,
		"failed":    DevToolPortalProvisioningStateFailed,
		"succeeded": DevToolPortalProvisioningStateSucceeded,
		"updating":  DevToolPortalProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DevToolPortalProvisioningState(input)
	return &out, nil
}

type EurekaServerEnabledState string

const (
	EurekaServerEnabledStateDisabled EurekaServerEnabledState = "Disabled"
	EurekaServerEnabledStateEnabled  EurekaServerEnabledState = "Enabled"
)

func PossibleValuesForEurekaServerEnabledState() []string {
	return []string{
		string(EurekaServerEnabledStateDisabled),
		string(EurekaServerEnabledStateEnabled),
	}
}

func (s *EurekaServerEnabledState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEurekaServerEnabledState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEurekaServerEnabledState(input string) (*EurekaServerEnabledState, error) {
	vals := map[string]EurekaServerEnabledState{
		"disabled": EurekaServerEnabledStateDisabled,
		"enabled":  EurekaServerEnabledStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EurekaServerEnabledState(input)
	return &out, nil
}

type EurekaServerState string

const (
	EurekaServerStateCanceled  EurekaServerState = "Canceled"
	EurekaServerStateFailed    EurekaServerState = "Failed"
	EurekaServerStateSucceeded EurekaServerState = "Succeeded"
	EurekaServerStateUpdating  EurekaServerState = "Updating"
)

func PossibleValuesForEurekaServerState() []string {
	return []string{
		string(EurekaServerStateCanceled),
		string(EurekaServerStateFailed),
		string(EurekaServerStateSucceeded),
		string(EurekaServerStateUpdating),
	}
}

func (s *EurekaServerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEurekaServerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEurekaServerState(input string) (*EurekaServerState, error) {
	vals := map[string]EurekaServerState{
		"canceled":  EurekaServerStateCanceled,
		"failed":    EurekaServerStateFailed,
		"succeeded": EurekaServerStateSucceeded,
		"updating":  EurekaServerStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EurekaServerState(input)
	return &out, nil
}

type Frequency string

const (
	FrequencyWeekly Frequency = "Weekly"
)

func PossibleValuesForFrequency() []string {
	return []string{
		string(FrequencyWeekly),
	}
}

func (s *Frequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFrequency(input string) (*Frequency, error) {
	vals := map[string]Frequency{
		"weekly": FrequencyWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Frequency(input)
	return &out, nil
}

type GatewayCertificateVerification string

const (
	GatewayCertificateVerificationDisabled GatewayCertificateVerification = "Disabled"
	GatewayCertificateVerificationEnabled  GatewayCertificateVerification = "Enabled"
)

func PossibleValuesForGatewayCertificateVerification() []string {
	return []string{
		string(GatewayCertificateVerificationDisabled),
		string(GatewayCertificateVerificationEnabled),
	}
}

func (s *GatewayCertificateVerification) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayCertificateVerification(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayCertificateVerification(input string) (*GatewayCertificateVerification, error) {
	vals := map[string]GatewayCertificateVerification{
		"disabled": GatewayCertificateVerificationDisabled,
		"enabled":  GatewayCertificateVerificationEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayCertificateVerification(input)
	return &out, nil
}

type GatewayProvisioningState string

const (
	GatewayProvisioningStateCreating  GatewayProvisioningState = "Creating"
	GatewayProvisioningStateDeleting  GatewayProvisioningState = "Deleting"
	GatewayProvisioningStateFailed    GatewayProvisioningState = "Failed"
	GatewayProvisioningStateSucceeded GatewayProvisioningState = "Succeeded"
	GatewayProvisioningStateUpdating  GatewayProvisioningState = "Updating"
)

func PossibleValuesForGatewayProvisioningState() []string {
	return []string{
		string(GatewayProvisioningStateCreating),
		string(GatewayProvisioningStateDeleting),
		string(GatewayProvisioningStateFailed),
		string(GatewayProvisioningStateSucceeded),
		string(GatewayProvisioningStateUpdating),
	}
}

func (s *GatewayProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayProvisioningState(input string) (*GatewayProvisioningState, error) {
	vals := map[string]GatewayProvisioningState{
		"creating":  GatewayProvisioningStateCreating,
		"deleting":  GatewayProvisioningStateDeleting,
		"failed":    GatewayProvisioningStateFailed,
		"succeeded": GatewayProvisioningStateSucceeded,
		"updating":  GatewayProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayProvisioningState(input)
	return &out, nil
}

type GatewayRouteConfigProtocol string

const (
	GatewayRouteConfigProtocolHTTP  GatewayRouteConfigProtocol = "HTTP"
	GatewayRouteConfigProtocolHTTPS GatewayRouteConfigProtocol = "HTTPS"
)

func PossibleValuesForGatewayRouteConfigProtocol() []string {
	return []string{
		string(GatewayRouteConfigProtocolHTTP),
		string(GatewayRouteConfigProtocolHTTPS),
	}
}

func (s *GatewayRouteConfigProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGatewayRouteConfigProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGatewayRouteConfigProtocol(input string) (*GatewayRouteConfigProtocol, error) {
	vals := map[string]GatewayRouteConfigProtocol{
		"http":  GatewayRouteConfigProtocolHTTP,
		"https": GatewayRouteConfigProtocolHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GatewayRouteConfigProtocol(input)
	return &out, nil
}

type GitImplementation string

const (
	GitImplementationGoNegativegit GitImplementation = "go-git"
	GitImplementationLibgitTwo     GitImplementation = "libgit2"
)

func PossibleValuesForGitImplementation() []string {
	return []string{
		string(GitImplementationGoNegativegit),
		string(GitImplementationLibgitTwo),
	}
}

func (s *GitImplementation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGitImplementation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGitImplementation(input string) (*GitImplementation, error) {
	vals := map[string]GitImplementation{
		"go-git":  GitImplementationGoNegativegit,
		"libgit2": GitImplementationLibgitTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GitImplementation(input)
	return &out, nil
}

type HTTPSchemeType string

const (
	HTTPSchemeTypeHTTP  HTTPSchemeType = "HTTP"
	HTTPSchemeTypeHTTPS HTTPSchemeType = "HTTPS"
)

func PossibleValuesForHTTPSchemeType() []string {
	return []string{
		string(HTTPSchemeTypeHTTP),
		string(HTTPSchemeTypeHTTPS),
	}
}

func (s *HTTPSchemeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHTTPSchemeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHTTPSchemeType(input string) (*HTTPSchemeType, error) {
	vals := map[string]HTTPSchemeType{
		"http":  HTTPSchemeTypeHTTP,
		"https": HTTPSchemeTypeHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HTTPSchemeType(input)
	return &out, nil
}

type KPackBuildStageProvisioningState string

const (
	KPackBuildStageProvisioningStateFailed     KPackBuildStageProvisioningState = "Failed"
	KPackBuildStageProvisioningStateNotStarted KPackBuildStageProvisioningState = "NotStarted"
	KPackBuildStageProvisioningStateRunning    KPackBuildStageProvisioningState = "Running"
	KPackBuildStageProvisioningStateSucceeded  KPackBuildStageProvisioningState = "Succeeded"
)

func PossibleValuesForKPackBuildStageProvisioningState() []string {
	return []string{
		string(KPackBuildStageProvisioningStateFailed),
		string(KPackBuildStageProvisioningStateNotStarted),
		string(KPackBuildStageProvisioningStateRunning),
		string(KPackBuildStageProvisioningStateSucceeded),
	}
}

func (s *KPackBuildStageProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKPackBuildStageProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKPackBuildStageProvisioningState(input string) (*KPackBuildStageProvisioningState, error) {
	vals := map[string]KPackBuildStageProvisioningState{
		"failed":     KPackBuildStageProvisioningStateFailed,
		"notstarted": KPackBuildStageProvisioningStateNotStarted,
		"running":    KPackBuildStageProvisioningStateRunning,
		"succeeded":  KPackBuildStageProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KPackBuildStageProvisioningState(input)
	return &out, nil
}

type KeyVaultCertificateAutoSync string

const (
	KeyVaultCertificateAutoSyncDisabled KeyVaultCertificateAutoSync = "Disabled"
	KeyVaultCertificateAutoSyncEnabled  KeyVaultCertificateAutoSync = "Enabled"
)

func PossibleValuesForKeyVaultCertificateAutoSync() []string {
	return []string{
		string(KeyVaultCertificateAutoSyncDisabled),
		string(KeyVaultCertificateAutoSyncEnabled),
	}
}

func (s *KeyVaultCertificateAutoSync) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyVaultCertificateAutoSync(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyVaultCertificateAutoSync(input string) (*KeyVaultCertificateAutoSync, error) {
	vals := map[string]KeyVaultCertificateAutoSync{
		"disabled": KeyVaultCertificateAutoSyncDisabled,
		"enabled":  KeyVaultCertificateAutoSyncEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultCertificateAutoSync(input)
	return &out, nil
}

type MonitoringSettingState string

const (
	MonitoringSettingStateFailed       MonitoringSettingState = "Failed"
	MonitoringSettingStateNotAvailable MonitoringSettingState = "NotAvailable"
	MonitoringSettingStateSucceeded    MonitoringSettingState = "Succeeded"
	MonitoringSettingStateUpdating     MonitoringSettingState = "Updating"
)

func PossibleValuesForMonitoringSettingState() []string {
	return []string{
		string(MonitoringSettingStateFailed),
		string(MonitoringSettingStateNotAvailable),
		string(MonitoringSettingStateSucceeded),
		string(MonitoringSettingStateUpdating),
	}
}

func (s *MonitoringSettingState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitoringSettingState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitoringSettingState(input string) (*MonitoringSettingState, error) {
	vals := map[string]MonitoringSettingState{
		"failed":       MonitoringSettingStateFailed,
		"notavailable": MonitoringSettingStateNotAvailable,
		"succeeded":    MonitoringSettingStateSucceeded,
		"updating":     MonitoringSettingStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitoringSettingState(input)
	return &out, nil
}

type PowerState string

const (
	PowerStateRunning PowerState = "Running"
	PowerStateStopped PowerState = "Stopped"
)

func PossibleValuesForPowerState() []string {
	return []string{
		string(PowerStateRunning),
		string(PowerStateStopped),
	}
}

func (s *PowerState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePowerState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePowerState(input string) (*PowerState, error) {
	vals := map[string]PowerState{
		"running": PowerStateRunning,
		"stopped": PowerStateStopped,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PowerState(input)
	return &out, nil
}

type PredefinedAcceleratorProvisioningState string

const (
	PredefinedAcceleratorProvisioningStateCanceled  PredefinedAcceleratorProvisioningState = "Canceled"
	PredefinedAcceleratorProvisioningStateCreating  PredefinedAcceleratorProvisioningState = "Creating"
	PredefinedAcceleratorProvisioningStateFailed    PredefinedAcceleratorProvisioningState = "Failed"
	PredefinedAcceleratorProvisioningStateSucceeded PredefinedAcceleratorProvisioningState = "Succeeded"
	PredefinedAcceleratorProvisioningStateUpdating  PredefinedAcceleratorProvisioningState = "Updating"
)

func PossibleValuesForPredefinedAcceleratorProvisioningState() []string {
	return []string{
		string(PredefinedAcceleratorProvisioningStateCanceled),
		string(PredefinedAcceleratorProvisioningStateCreating),
		string(PredefinedAcceleratorProvisioningStateFailed),
		string(PredefinedAcceleratorProvisioningStateSucceeded),
		string(PredefinedAcceleratorProvisioningStateUpdating),
	}
}

func (s *PredefinedAcceleratorProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePredefinedAcceleratorProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePredefinedAcceleratorProvisioningState(input string) (*PredefinedAcceleratorProvisioningState, error) {
	vals := map[string]PredefinedAcceleratorProvisioningState{
		"canceled":  PredefinedAcceleratorProvisioningStateCanceled,
		"creating":  PredefinedAcceleratorProvisioningStateCreating,
		"failed":    PredefinedAcceleratorProvisioningStateFailed,
		"succeeded": PredefinedAcceleratorProvisioningStateSucceeded,
		"updating":  PredefinedAcceleratorProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PredefinedAcceleratorProvisioningState(input)
	return &out, nil
}

type PredefinedAcceleratorState string

const (
	PredefinedAcceleratorStateDisabled PredefinedAcceleratorState = "Disabled"
	PredefinedAcceleratorStateEnabled  PredefinedAcceleratorState = "Enabled"
)

func PossibleValuesForPredefinedAcceleratorState() []string {
	return []string{
		string(PredefinedAcceleratorStateDisabled),
		string(PredefinedAcceleratorStateEnabled),
	}
}

func (s *PredefinedAcceleratorState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePredefinedAcceleratorState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePredefinedAcceleratorState(input string) (*PredefinedAcceleratorState, error) {
	vals := map[string]PredefinedAcceleratorState{
		"disabled": PredefinedAcceleratorStateDisabled,
		"enabled":  PredefinedAcceleratorStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PredefinedAcceleratorState(input)
	return &out, nil
}

type ProbeActionType string

const (
	ProbeActionTypeExecAction      ProbeActionType = "ExecAction"
	ProbeActionTypeHTTPGetAction   ProbeActionType = "HTTPGetAction"
	ProbeActionTypeTCPSocketAction ProbeActionType = "TCPSocketAction"
)

func PossibleValuesForProbeActionType() []string {
	return []string{
		string(ProbeActionTypeExecAction),
		string(ProbeActionTypeHTTPGetAction),
		string(ProbeActionTypeTCPSocketAction),
	}
}

func (s *ProbeActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProbeActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProbeActionType(input string) (*ProbeActionType, error) {
	vals := map[string]ProbeActionType{
		"execaction":      ProbeActionTypeExecAction,
		"httpgetaction":   ProbeActionTypeHTTPGetAction,
		"tcpsocketaction": ProbeActionTypeTCPSocketAction,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProbeActionType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCreating   ProvisioningState = "Creating"
	ProvisioningStateDeleted    ProvisioningState = "Deleted"
	ProvisioningStateDeleting   ProvisioningState = "Deleting"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateMoveFailed ProvisioningState = "MoveFailed"
	ProvisioningStateMoved      ProvisioningState = "Moved"
	ProvisioningStateMoving     ProvisioningState = "Moving"
	ProvisioningStateStarting   ProvisioningState = "Starting"
	ProvisioningStateStopping   ProvisioningState = "Stopping"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
	ProvisioningStateUpdating   ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoveFailed),
		string(ProvisioningStateMoved),
		string(ProvisioningStateMoving),
		string(ProvisioningStateStarting),
		string(ProvisioningStateStopping),
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
		"creating":   ProvisioningStateCreating,
		"deleted":    ProvisioningStateDeleted,
		"deleting":   ProvisioningStateDeleting,
		"failed":     ProvisioningStateFailed,
		"movefailed": ProvisioningStateMoveFailed,
		"moved":      ProvisioningStateMoved,
		"moving":     ProvisioningStateMoving,
		"starting":   ProvisioningStateStarting,
		"stopping":   ProvisioningStateStopping,
		"succeeded":  ProvisioningStateSucceeded,
		"updating":   ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ResourceSkuRestrictionsReasonCode string

const (
	ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription ResourceSkuRestrictionsReasonCode = "NotAvailableForSubscription"
	ResourceSkuRestrictionsReasonCodeQuotaId                     ResourceSkuRestrictionsReasonCode = "QuotaId"
)

func PossibleValuesForResourceSkuRestrictionsReasonCode() []string {
	return []string{
		string(ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription),
		string(ResourceSkuRestrictionsReasonCodeQuotaId),
	}
}

func (s *ResourceSkuRestrictionsReasonCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceSkuRestrictionsReasonCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceSkuRestrictionsReasonCode(input string) (*ResourceSkuRestrictionsReasonCode, error) {
	vals := map[string]ResourceSkuRestrictionsReasonCode{
		"notavailableforsubscription": ResourceSkuRestrictionsReasonCodeNotAvailableForSubscription,
		"quotaid":                     ResourceSkuRestrictionsReasonCodeQuotaId,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuRestrictionsReasonCode(input)
	return &out, nil
}

type ResourceSkuRestrictionsType string

const (
	ResourceSkuRestrictionsTypeLocation ResourceSkuRestrictionsType = "Location"
	ResourceSkuRestrictionsTypeZone     ResourceSkuRestrictionsType = "Zone"
)

func PossibleValuesForResourceSkuRestrictionsType() []string {
	return []string{
		string(ResourceSkuRestrictionsTypeLocation),
		string(ResourceSkuRestrictionsTypeZone),
	}
}

func (s *ResourceSkuRestrictionsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceSkuRestrictionsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceSkuRestrictionsType(input string) (*ResourceSkuRestrictionsType, error) {
	vals := map[string]ResourceSkuRestrictionsType{
		"location": ResourceSkuRestrictionsTypeLocation,
		"zone":     ResourceSkuRestrictionsTypeZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceSkuRestrictionsType(input)
	return &out, nil
}

type ServiceRegistryProvisioningState string

const (
	ServiceRegistryProvisioningStateCreating  ServiceRegistryProvisioningState = "Creating"
	ServiceRegistryProvisioningStateDeleting  ServiceRegistryProvisioningState = "Deleting"
	ServiceRegistryProvisioningStateFailed    ServiceRegistryProvisioningState = "Failed"
	ServiceRegistryProvisioningStateSucceeded ServiceRegistryProvisioningState = "Succeeded"
	ServiceRegistryProvisioningStateUpdating  ServiceRegistryProvisioningState = "Updating"
)

func PossibleValuesForServiceRegistryProvisioningState() []string {
	return []string{
		string(ServiceRegistryProvisioningStateCreating),
		string(ServiceRegistryProvisioningStateDeleting),
		string(ServiceRegistryProvisioningStateFailed),
		string(ServiceRegistryProvisioningStateSucceeded),
		string(ServiceRegistryProvisioningStateUpdating),
	}
}

func (s *ServiceRegistryProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceRegistryProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceRegistryProvisioningState(input string) (*ServiceRegistryProvisioningState, error) {
	vals := map[string]ServiceRegistryProvisioningState{
		"creating":  ServiceRegistryProvisioningStateCreating,
		"deleting":  ServiceRegistryProvisioningStateDeleting,
		"failed":    ServiceRegistryProvisioningStateFailed,
		"succeeded": ServiceRegistryProvisioningStateSucceeded,
		"updating":  ServiceRegistryProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceRegistryProvisioningState(input)
	return &out, nil
}

type SessionAffinity string

const (
	SessionAffinityCookie SessionAffinity = "Cookie"
	SessionAffinityNone   SessionAffinity = "None"
)

func PossibleValuesForSessionAffinity() []string {
	return []string{
		string(SessionAffinityCookie),
		string(SessionAffinityNone),
	}
}

func (s *SessionAffinity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSessionAffinity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSessionAffinity(input string) (*SessionAffinity, error) {
	vals := map[string]SessionAffinity{
		"cookie": SessionAffinityCookie,
		"none":   SessionAffinityNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SessionAffinity(input)
	return &out, nil
}

type SkuScaleType string

const (
	SkuScaleTypeAutomatic SkuScaleType = "Automatic"
	SkuScaleTypeManual    SkuScaleType = "Manual"
	SkuScaleTypeNone      SkuScaleType = "None"
)

func PossibleValuesForSkuScaleType() []string {
	return []string{
		string(SkuScaleTypeAutomatic),
		string(SkuScaleTypeManual),
		string(SkuScaleTypeNone),
	}
}

func (s *SkuScaleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSkuScaleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSkuScaleType(input string) (*SkuScaleType, error) {
	vals := map[string]SkuScaleType{
		"automatic": SkuScaleTypeAutomatic,
		"manual":    SkuScaleTypeManual,
		"none":      SkuScaleTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SkuScaleType(input)
	return &out, nil
}

type StorageType string

const (
	StorageTypeStorageAccount StorageType = "StorageAccount"
)

func PossibleValuesForStorageType() []string {
	return []string{
		string(StorageTypeStorageAccount),
	}
}

func (s *StorageType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageType(input string) (*StorageType, error) {
	vals := map[string]StorageType{
		"storageaccount": StorageTypeStorageAccount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageType(input)
	return &out, nil
}

type SupportedRuntimePlatform string

const (
	SupportedRuntimePlatformJava         SupportedRuntimePlatform = "Java"
	SupportedRuntimePlatformPointNETCore SupportedRuntimePlatform = ".NET Core"
)

func PossibleValuesForSupportedRuntimePlatform() []string {
	return []string{
		string(SupportedRuntimePlatformJava),
		string(SupportedRuntimePlatformPointNETCore),
	}
}

func (s *SupportedRuntimePlatform) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSupportedRuntimePlatform(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSupportedRuntimePlatform(input string) (*SupportedRuntimePlatform, error) {
	vals := map[string]SupportedRuntimePlatform{
		"java":      SupportedRuntimePlatformJava,
		".net core": SupportedRuntimePlatformPointNETCore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SupportedRuntimePlatform(input)
	return &out, nil
}

type SupportedRuntimeValue string

const (
	SupportedRuntimeValueJavaEight       SupportedRuntimeValue = "Java_8"
	SupportedRuntimeValueJavaOneOne      SupportedRuntimeValue = "Java_11"
	SupportedRuntimeValueJavaOneSeven    SupportedRuntimeValue = "Java_17"
	SupportedRuntimeValueJavaTwoOne      SupportedRuntimeValue = "Java_21"
	SupportedRuntimeValueNetCoreThreeOne SupportedRuntimeValue = "NetCore_31"
)

func PossibleValuesForSupportedRuntimeValue() []string {
	return []string{
		string(SupportedRuntimeValueJavaEight),
		string(SupportedRuntimeValueJavaOneOne),
		string(SupportedRuntimeValueJavaOneSeven),
		string(SupportedRuntimeValueJavaTwoOne),
		string(SupportedRuntimeValueNetCoreThreeOne),
	}
}

func (s *SupportedRuntimeValue) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSupportedRuntimeValue(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSupportedRuntimeValue(input string) (*SupportedRuntimeValue, error) {
	vals := map[string]SupportedRuntimeValue{
		"java_8":     SupportedRuntimeValueJavaEight,
		"java_11":    SupportedRuntimeValueJavaOneOne,
		"java_17":    SupportedRuntimeValueJavaOneSeven,
		"java_21":    SupportedRuntimeValueJavaTwoOne,
		"netcore_31": SupportedRuntimeValueNetCoreThreeOne,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SupportedRuntimeValue(input)
	return &out, nil
}

type TestEndpointAuthState string

const (
	TestEndpointAuthStateDisabled TestEndpointAuthState = "Disabled"
	TestEndpointAuthStateEnabled  TestEndpointAuthState = "Enabled"
)

func PossibleValuesForTestEndpointAuthState() []string {
	return []string{
		string(TestEndpointAuthStateDisabled),
		string(TestEndpointAuthStateEnabled),
	}
}

func (s *TestEndpointAuthState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTestEndpointAuthState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTestEndpointAuthState(input string) (*TestEndpointAuthState, error) {
	vals := map[string]TestEndpointAuthState{
		"disabled": TestEndpointAuthStateDisabled,
		"enabled":  TestEndpointAuthStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TestEndpointAuthState(input)
	return &out, nil
}

type TestKeyType string

const (
	TestKeyTypePrimary   TestKeyType = "Primary"
	TestKeyTypeSecondary TestKeyType = "Secondary"
)

func PossibleValuesForTestKeyType() []string {
	return []string{
		string(TestKeyTypePrimary),
		string(TestKeyTypeSecondary),
	}
}

func (s *TestKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTestKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTestKeyType(input string) (*TestKeyType, error) {
	vals := map[string]TestKeyType{
		"primary":   TestKeyTypePrimary,
		"secondary": TestKeyTypeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TestKeyType(input)
	return &out, nil
}

type TrafficDirection string

const (
	TrafficDirectionInbound  TrafficDirection = "Inbound"
	TrafficDirectionOutbound TrafficDirection = "Outbound"
)

func PossibleValuesForTrafficDirection() []string {
	return []string{
		string(TrafficDirectionInbound),
		string(TrafficDirectionOutbound),
	}
}

func (s *TrafficDirection) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTrafficDirection(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTrafficDirection(input string) (*TrafficDirection, error) {
	vals := map[string]TrafficDirection{
		"inbound":  TrafficDirectionInbound,
		"outbound": TrafficDirectionOutbound,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TrafficDirection(input)
	return &out, nil
}

type TriggeredBuildResultProvisioningState string

const (
	TriggeredBuildResultProvisioningStateBuilding  TriggeredBuildResultProvisioningState = "Building"
	TriggeredBuildResultProvisioningStateCanceled  TriggeredBuildResultProvisioningState = "Canceled"
	TriggeredBuildResultProvisioningStateDeleting  TriggeredBuildResultProvisioningState = "Deleting"
	TriggeredBuildResultProvisioningStateFailed    TriggeredBuildResultProvisioningState = "Failed"
	TriggeredBuildResultProvisioningStateQueuing   TriggeredBuildResultProvisioningState = "Queuing"
	TriggeredBuildResultProvisioningStateSucceeded TriggeredBuildResultProvisioningState = "Succeeded"
)

func PossibleValuesForTriggeredBuildResultProvisioningState() []string {
	return []string{
		string(TriggeredBuildResultProvisioningStateBuilding),
		string(TriggeredBuildResultProvisioningStateCanceled),
		string(TriggeredBuildResultProvisioningStateDeleting),
		string(TriggeredBuildResultProvisioningStateFailed),
		string(TriggeredBuildResultProvisioningStateQueuing),
		string(TriggeredBuildResultProvisioningStateSucceeded),
	}
}

func (s *TriggeredBuildResultProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggeredBuildResultProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggeredBuildResultProvisioningState(input string) (*TriggeredBuildResultProvisioningState, error) {
	vals := map[string]TriggeredBuildResultProvisioningState{
		"building":  TriggeredBuildResultProvisioningStateBuilding,
		"canceled":  TriggeredBuildResultProvisioningStateCanceled,
		"deleting":  TriggeredBuildResultProvisioningStateDeleting,
		"failed":    TriggeredBuildResultProvisioningStateFailed,
		"queuing":   TriggeredBuildResultProvisioningStateQueuing,
		"succeeded": TriggeredBuildResultProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggeredBuildResultProvisioningState(input)
	return &out, nil
}

type Type string

const (
	TypeAzureFileVolume Type = "AzureFileVolume"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeAzureFileVolume),
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
		"azurefilevolume": TypeAzureFileVolume,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
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
