package pipelinegroups

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AllowedFormats string

const (
	AllowedFormatsAll                      AllowedFormats = "all"
	AllowedFormatsCefRfcFiveFourTwoFour    AllowedFormats = "cefRfc5424"
	AllowedFormatsCefRfcThreeOneSixFour    AllowedFormats = "cefRfc3164"
	AllowedFormatsRawCef                   AllowedFormats = "rawCef"
	AllowedFormatsSyslogRfcFiveFourTwoFour AllowedFormats = "syslogRfc5424"
	AllowedFormatsSyslogRfcThreeOneSixFour AllowedFormats = "syslogRfc3164"
)

func PossibleValuesForAllowedFormats() []string {
	return []string{
		string(AllowedFormatsAll),
		string(AllowedFormatsCefRfcFiveFourTwoFour),
		string(AllowedFormatsCefRfcThreeOneSixFour),
		string(AllowedFormatsRawCef),
		string(AllowedFormatsSyslogRfcFiveFourTwoFour),
		string(AllowedFormatsSyslogRfcThreeOneSixFour),
	}
}

func (s *AllowedFormats) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAllowedFormats(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAllowedFormats(input string) (*AllowedFormats, error) {
	vals := map[string]AllowedFormats{
		"all":           AllowedFormatsAll,
		"cefrfc5424":    AllowedFormatsCefRfcFiveFourTwoFour,
		"cefrfc3164":    AllowedFormatsCefRfcThreeOneSixFour,
		"rawcef":        AllowedFormatsRawCef,
		"syslogrfc5424": AllowedFormatsSyslogRfcFiveFourTwoFour,
		"syslogrfc3164": AllowedFormatsSyslogRfcThreeOneSixFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AllowedFormats(input)
	return &out, nil
}

type CapabilityOperator string

const (
	CapabilityOperatorDoesNotExist CapabilityOperator = "DoesNotExist"
	CapabilityOperatorExists       CapabilityOperator = "Exists"
	CapabilityOperatorIn           CapabilityOperator = "In"
	CapabilityOperatorNotIn        CapabilityOperator = "NotIn"
)

func PossibleValuesForCapabilityOperator() []string {
	return []string{
		string(CapabilityOperatorDoesNotExist),
		string(CapabilityOperatorExists),
		string(CapabilityOperatorIn),
		string(CapabilityOperatorNotIn),
	}
}

func (s *CapabilityOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCapabilityOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCapabilityOperator(input string) (*CapabilityOperator, error) {
	vals := map[string]CapabilityOperator{
		"doesnotexist": CapabilityOperatorDoesNotExist,
		"exists":       CapabilityOperatorExists,
		"in":           CapabilityOperatorIn,
		"notin":        CapabilityOperatorNotIn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CapabilityOperator(input)
	return &out, nil
}

type CertificateSourceType string

const (
	CertificateSourceTypeKubernetesConfigMap CertificateSourceType = "kubernetesConfigMap"
	CertificateSourceTypeKubernetesSecret    CertificateSourceType = "kubernetesSecret"
)

func PossibleValuesForCertificateSourceType() []string {
	return []string{
		string(CertificateSourceTypeKubernetesConfigMap),
		string(CertificateSourceTypeKubernetesSecret),
	}
}

func (s *CertificateSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateSourceType(input string) (*CertificateSourceType, error) {
	vals := map[string]CertificateSourceType{
		"kubernetesconfigmap": CertificateSourceTypeKubernetesConfigMap,
		"kubernetessecret":    CertificateSourceTypeKubernetesSecret,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateSourceType(input)
	return &out, nil
}

type ExporterType string

const (
	ExporterTypeAzureMonitorWorkspaceLogs ExporterType = "AzureMonitorWorkspaceLogs"
)

func PossibleValuesForExporterType() []string {
	return []string{
		string(ExporterTypeAzureMonitorWorkspaceLogs),
	}
}

func (s *ExporterType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExporterType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExporterType(input string) (*ExporterType, error) {
	vals := map[string]ExporterType{
		"azuremonitorworkspacelogs": ExporterTypeAzureMonitorWorkspaceLogs,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExporterType(input)
	return &out, nil
}

type ExtendedLocationType string

const (
	ExtendedLocationTypeCustomLocation ExtendedLocationType = "CustomLocation"
	ExtendedLocationTypeEdgeZone       ExtendedLocationType = "EdgeZone"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeCustomLocation),
		string(ExtendedLocationTypeEdgeZone),
	}
}

func (s *ExtendedLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExtendedLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExtendedLocationType(input string) (*ExtendedLocationType, error) {
	vals := map[string]ExtendedLocationType{
		"customlocation": ExtendedLocationTypeCustomLocation,
		"edgezone":       ExtendedLocationTypeEdgeZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type PipelineType string

const (
	PipelineTypeLogs PipelineType = "Logs"
)

func PossibleValuesForPipelineType() []string {
	return []string{
		string(PipelineTypeLogs),
	}
}

func (s *PipelineType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePipelineType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePipelineType(input string) (*PipelineType, error) {
	vals := map[string]PipelineType{
		"logs": PipelineTypeLogs,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PipelineType(input)
	return &out, nil
}

type PrivateKeySourceType string

const (
	PrivateKeySourceTypeKubernetesSecret PrivateKeySourceType = "kubernetesSecret"
)

func PossibleValuesForPrivateKeySourceType() []string {
	return []string{
		string(PrivateKeySourceTypeKubernetesSecret),
	}
}

func (s *PrivateKeySourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateKeySourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateKeySourceType(input string) (*PrivateKeySourceType, error) {
	vals := map[string]PrivateKeySourceType{
		"kubernetessecret": PrivateKeySourceTypeKubernetesSecret,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateKeySourceType(input)
	return &out, nil
}

type ProcessorType string

const (
	ProcessorTypeBatch                      ProcessorType = "Batch"
	ProcessorTypeMicrosoftCommonSecurityLog ProcessorType = "MicrosoftCommonSecurityLog"
	ProcessorTypeMicrosoftSyslog            ProcessorType = "MicrosoftSyslog"
	ProcessorTypeTransformLanguage          ProcessorType = "TransformLanguage"
)

func PossibleValuesForProcessorType() []string {
	return []string{
		string(ProcessorTypeBatch),
		string(ProcessorTypeMicrosoftCommonSecurityLog),
		string(ProcessorTypeMicrosoftSyslog),
		string(ProcessorTypeTransformLanguage),
	}
}

func (s *ProcessorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProcessorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProcessorType(input string) (*ProcessorType, error) {
	vals := map[string]ProcessorType{
		"batch":                      ProcessorTypeBatch,
		"microsoftcommonsecuritylog": ProcessorTypeMicrosoftCommonSecurityLog,
		"microsoftsyslog":            ProcessorTypeMicrosoftSyslog,
		"transformlanguage":          ProcessorTypeTransformLanguage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProcessorType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
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
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ReceiverType string

const (
	ReceiverTypeOTLP   ReceiverType = "OTLP"
	ReceiverTypeSyslog ReceiverType = "Syslog"
)

func PossibleValuesForReceiverType() []string {
	return []string{
		string(ReceiverTypeOTLP),
		string(ReceiverTypeSyslog),
	}
}

func (s *ReceiverType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReceiverType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReceiverType(input string) (*ReceiverType, error) {
	vals := map[string]ReceiverType{
		"otlp":   ReceiverTypeOTLP,
		"syslog": ReceiverTypeSyslog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReceiverType(input)
	return &out, nil
}

type TlsMode string

const (
	TlsModeDisabled   TlsMode = "disabled"
	TlsModeMutualTls  TlsMode = "mutualTls"
	TlsModeServerOnly TlsMode = "serverOnly"
)

func PossibleValuesForTlsMode() []string {
	return []string{
		string(TlsModeDisabled),
		string(TlsModeMutualTls),
		string(TlsModeServerOnly),
	}
}

func (s *TlsMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsMode(input string) (*TlsMode, error) {
	vals := map[string]TlsMode{
		"disabled":   TlsModeDisabled,
		"mutualtls":  TlsModeMutualTls,
		"serveronly": TlsModeServerOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsMode(input)
	return &out, nil
}

type TransportProtocol string

const (
	TransportProtocolTcp TransportProtocol = "tcp"
	TransportProtocolUdp TransportProtocol = "udp"
)

func PossibleValuesForTransportProtocol() []string {
	return []string{
		string(TransportProtocolTcp),
		string(TransportProtocolUdp),
	}
}

func (s *TransportProtocol) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTransportProtocol(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTransportProtocol(input string) (*TransportProtocol, error) {
	vals := map[string]TransportProtocol{
		"tcp": TransportProtocolTcp,
		"udp": TransportProtocolUdp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TransportProtocol(input)
	return &out, nil
}
