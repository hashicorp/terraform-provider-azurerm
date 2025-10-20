package brokerlistener

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerProtocolType string

const (
	BrokerProtocolTypeMqtt       BrokerProtocolType = "Mqtt"
	BrokerProtocolTypeWebSockets BrokerProtocolType = "WebSockets"
)

func PossibleValuesForBrokerProtocolType() []string {
	return []string{
		string(BrokerProtocolTypeMqtt),
		string(BrokerProtocolTypeWebSockets),
	}
}

func (s *BrokerProtocolType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBrokerProtocolType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBrokerProtocolType(input string) (*BrokerProtocolType, error) {
	vals := map[string]BrokerProtocolType{
		"mqtt":       BrokerProtocolTypeMqtt,
		"websockets": BrokerProtocolTypeWebSockets,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BrokerProtocolType(input)
	return &out, nil
}

type CertManagerIssuerKind string

const (
	CertManagerIssuerKindClusterIssuer CertManagerIssuerKind = "ClusterIssuer"
	CertManagerIssuerKindIssuer        CertManagerIssuerKind = "Issuer"
)

func PossibleValuesForCertManagerIssuerKind() []string {
	return []string{
		string(CertManagerIssuerKindClusterIssuer),
		string(CertManagerIssuerKindIssuer),
	}
}

func (s *CertManagerIssuerKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertManagerIssuerKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertManagerIssuerKind(input string) (*CertManagerIssuerKind, error) {
	vals := map[string]CertManagerIssuerKind{
		"clusterissuer": CertManagerIssuerKindClusterIssuer,
		"issuer":        CertManagerIssuerKindIssuer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertManagerIssuerKind(input)
	return &out, nil
}

type ExtendedLocationType string

const (
	ExtendedLocationTypeCustomLocation ExtendedLocationType = "CustomLocation"
)

func PossibleValuesForExtendedLocationType() []string {
	return []string{
		string(ExtendedLocationTypeCustomLocation),
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
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExtendedLocationType(input)
	return &out, nil
}

type PrivateKeyAlgorithm string

const (
	PrivateKeyAlgorithmEcFiveTwoOne         PrivateKeyAlgorithm = "Ec521"
	PrivateKeyAlgorithmEcThreeEightFour     PrivateKeyAlgorithm = "Ec384"
	PrivateKeyAlgorithmEcTwoFiveSix         PrivateKeyAlgorithm = "Ec256"
	PrivateKeyAlgorithmEdTwoFiveFiveOneNine PrivateKeyAlgorithm = "Ed25519"
	PrivateKeyAlgorithmRsaEightOneNineTwo   PrivateKeyAlgorithm = "Rsa8192"
	PrivateKeyAlgorithmRsaFourZeroNineSix   PrivateKeyAlgorithm = "Rsa4096"
	PrivateKeyAlgorithmRsaTwoZeroFourEight  PrivateKeyAlgorithm = "Rsa2048"
)

func PossibleValuesForPrivateKeyAlgorithm() []string {
	return []string{
		string(PrivateKeyAlgorithmEcFiveTwoOne),
		string(PrivateKeyAlgorithmEcThreeEightFour),
		string(PrivateKeyAlgorithmEcTwoFiveSix),
		string(PrivateKeyAlgorithmEdTwoFiveFiveOneNine),
		string(PrivateKeyAlgorithmRsaEightOneNineTwo),
		string(PrivateKeyAlgorithmRsaFourZeroNineSix),
		string(PrivateKeyAlgorithmRsaTwoZeroFourEight),
	}
}

func (s *PrivateKeyAlgorithm) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateKeyAlgorithm(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateKeyAlgorithm(input string) (*PrivateKeyAlgorithm, error) {
	vals := map[string]PrivateKeyAlgorithm{
		"ec521":   PrivateKeyAlgorithmEcFiveTwoOne,
		"ec384":   PrivateKeyAlgorithmEcThreeEightFour,
		"ec256":   PrivateKeyAlgorithmEcTwoFiveSix,
		"ed25519": PrivateKeyAlgorithmEdTwoFiveFiveOneNine,
		"rsa8192": PrivateKeyAlgorithmRsaEightOneNineTwo,
		"rsa4096": PrivateKeyAlgorithmRsaFourZeroNineSix,
		"rsa2048": PrivateKeyAlgorithmRsaTwoZeroFourEight,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateKeyAlgorithm(input)
	return &out, nil
}

type PrivateKeyRotationPolicy string

const (
	PrivateKeyRotationPolicyAlways PrivateKeyRotationPolicy = "Always"
	PrivateKeyRotationPolicyNever  PrivateKeyRotationPolicy = "Never"
)

func PossibleValuesForPrivateKeyRotationPolicy() []string {
	return []string{
		string(PrivateKeyRotationPolicyAlways),
		string(PrivateKeyRotationPolicyNever),
	}
}

func (s *PrivateKeyRotationPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrivateKeyRotationPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrivateKeyRotationPolicy(input string) (*PrivateKeyRotationPolicy, error) {
	vals := map[string]PrivateKeyRotationPolicy{
		"always": PrivateKeyRotationPolicyAlways,
		"never":  PrivateKeyRotationPolicyNever,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrivateKeyRotationPolicy(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted     ProvisioningState = "Accepted"
	ProvisioningStateCanceled     ProvisioningState = "Canceled"
	ProvisioningStateDeleting     ProvisioningState = "Deleting"
	ProvisioningStateFailed       ProvisioningState = "Failed"
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded    ProvisioningState = "Succeeded"
	ProvisioningStateUpdating     ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateProvisioning),
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
		"deleting":     ProvisioningStateDeleting,
		"failed":       ProvisioningStateFailed,
		"provisioning": ProvisioningStateProvisioning,
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

type ServiceType string

const (
	ServiceTypeClusterIP    ServiceType = "ClusterIp"
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
	ServiceTypeNodePort     ServiceType = "NodePort"
)

func PossibleValuesForServiceType() []string {
	return []string{
		string(ServiceTypeClusterIP),
		string(ServiceTypeLoadBalancer),
		string(ServiceTypeNodePort),
	}
}

func (s *ServiceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServiceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServiceType(input string) (*ServiceType, error) {
	vals := map[string]ServiceType{
		"clusterip":    ServiceTypeClusterIP,
		"loadbalancer": ServiceTypeLoadBalancer,
		"nodeport":     ServiceTypeNodePort,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServiceType(input)
	return &out, nil
}

type TlsCertMethodMode string

const (
	TlsCertMethodModeAutomatic TlsCertMethodMode = "Automatic"
	TlsCertMethodModeManual    TlsCertMethodMode = "Manual"
)

func PossibleValuesForTlsCertMethodMode() []string {
	return []string{
		string(TlsCertMethodModeAutomatic),
		string(TlsCertMethodModeManual),
	}
}

func (s *TlsCertMethodMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTlsCertMethodMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTlsCertMethodMode(input string) (*TlsCertMethodMode, error) {
	vals := map[string]TlsCertMethodMode{
		"automatic": TlsCertMethodModeAutomatic,
		"manual":    TlsCertMethodModeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TlsCertMethodMode(input)
	return &out, nil
}
