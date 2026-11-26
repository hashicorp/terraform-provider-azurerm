package hcpopenshiftclusters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterImageRegistryState string

const (
	ClusterImageRegistryStateDisabled ClusterImageRegistryState = "Disabled"
	ClusterImageRegistryStateEnabled  ClusterImageRegistryState = "Enabled"
)

func PossibleValuesForClusterImageRegistryState() []string {
	return []string{
		string(ClusterImageRegistryStateDisabled),
		string(ClusterImageRegistryStateEnabled),
	}
}

func (s *ClusterImageRegistryState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseClusterImageRegistryState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseClusterImageRegistryState(input string) (*ClusterImageRegistryState, error) {
	vals := map[string]ClusterImageRegistryState{
		"disabled": ClusterImageRegistryStateDisabled,
		"enabled":  ClusterImageRegistryStateEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ClusterImageRegistryState(input)
	return &out, nil
}

type ConditionType string

const (
	ConditionTypeAvailable   ConditionType = "Available"
	ConditionTypeDegraded    ConditionType = "Degraded"
	ConditionTypeProgressing ConditionType = "Progressing"
)

func PossibleValuesForConditionType() []string {
	return []string{
		string(ConditionTypeAvailable),
		string(ConditionTypeDegraded),
		string(ConditionTypeProgressing),
	}
}

func (s *ConditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConditionType(input string) (*ConditionType, error) {
	vals := map[string]ConditionType{
		"available":   ConditionTypeAvailable,
		"degraded":    ConditionTypeDegraded,
		"progressing": ConditionTypeProgressing,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConditionType(input)
	return &out, nil
}

type CryptoRestrictions string

const (
	CryptoRestrictionsFIPS CryptoRestrictions = "FIPS"
	CryptoRestrictionsNone CryptoRestrictions = "None"
)

func PossibleValuesForCryptoRestrictions() []string {
	return []string{
		string(CryptoRestrictionsFIPS),
		string(CryptoRestrictionsNone),
	}
}

func (s *CryptoRestrictions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCryptoRestrictions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCryptoRestrictions(input string) (*CryptoRestrictions, error) {
	vals := map[string]CryptoRestrictions{
		"fips": CryptoRestrictionsFIPS,
		"none": CryptoRestrictionsNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CryptoRestrictions(input)
	return &out, nil
}

type CustomerManagedEncryptionType string

const (
	CustomerManagedEncryptionTypeKMS CustomerManagedEncryptionType = "KMS"
)

func PossibleValuesForCustomerManagedEncryptionType() []string {
	return []string{
		string(CustomerManagedEncryptionTypeKMS),
	}
}

func (s *CustomerManagedEncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCustomerManagedEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCustomerManagedEncryptionType(input string) (*CustomerManagedEncryptionType, error) {
	vals := map[string]CustomerManagedEncryptionType{
		"kms": CustomerManagedEncryptionTypeKMS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CustomerManagedEncryptionType(input)
	return &out, nil
}

type EtcdDataEncryptionKeyManagementModeType string

const (
	EtcdDataEncryptionKeyManagementModeTypeCustomerManaged EtcdDataEncryptionKeyManagementModeType = "CustomerManaged"
)

func PossibleValuesForEtcdDataEncryptionKeyManagementModeType() []string {
	return []string{
		string(EtcdDataEncryptionKeyManagementModeTypeCustomerManaged),
	}
}

func (s *EtcdDataEncryptionKeyManagementModeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEtcdDataEncryptionKeyManagementModeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEtcdDataEncryptionKeyManagementModeType(input string) (*EtcdDataEncryptionKeyManagementModeType, error) {
	vals := map[string]EtcdDataEncryptionKeyManagementModeType{
		"customermanaged": EtcdDataEncryptionKeyManagementModeTypeCustomerManaged,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EtcdDataEncryptionKeyManagementModeType(input)
	return &out, nil
}

type IngressType string

const (
	IngressTypeDisabled IngressType = "Disabled"
	IngressTypePrivate  IngressType = "Private"
	IngressTypePublic   IngressType = "Public"
)

func PossibleValuesForIngressType() []string {
	return []string{
		string(IngressTypeDisabled),
		string(IngressTypePrivate),
		string(IngressTypePublic),
	}
}

func (s *IngressType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIngressType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIngressType(input string) (*IngressType, error) {
	vals := map[string]IngressType{
		"disabled": IngressTypeDisabled,
		"private":  IngressTypePrivate,
		"public":   IngressTypePublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IngressType(input)
	return &out, nil
}

type KeyVaultVisibility string

const (
	KeyVaultVisibilityPrivate KeyVaultVisibility = "Private"
	KeyVaultVisibilityPublic  KeyVaultVisibility = "Public"
)

func PossibleValuesForKeyVaultVisibility() []string {
	return []string{
		string(KeyVaultVisibilityPrivate),
		string(KeyVaultVisibilityPublic),
	}
}

func (s *KeyVaultVisibility) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyVaultVisibility(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyVaultVisibility(input string) (*KeyVaultVisibility, error) {
	vals := map[string]KeyVaultVisibility{
		"private": KeyVaultVisibilityPrivate,
		"public":  KeyVaultVisibilityPublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultVisibility(input)
	return &out, nil
}

type KmsKeyVaultType string

const (
	KmsKeyVaultTypeKeyVault   KmsKeyVaultType = "KeyVault"
	KmsKeyVaultTypeManagedHSM KmsKeyVaultType = "ManagedHSM"
)

func PossibleValuesForKmsKeyVaultType() []string {
	return []string{
		string(KmsKeyVaultTypeKeyVault),
		string(KmsKeyVaultTypeManagedHSM),
	}
}

func (s *KmsKeyVaultType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKmsKeyVaultType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKmsKeyVaultType(input string) (*KmsKeyVaultType, error) {
	vals := map[string]KmsKeyVaultType{
		"keyvault":   KmsKeyVaultTypeKeyVault,
		"managedhsm": KmsKeyVaultTypeManagedHSM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KmsKeyVaultType(input)
	return &out, nil
}

type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone                       ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned             ManagedServiceIdentityType = "SystemAssigned"
	ManagedServiceIdentityTypeSystemAssignedUserAssigned ManagedServiceIdentityType = "SystemAssigned,UserAssigned"
	ManagedServiceIdentityTypeUserAssigned               ManagedServiceIdentityType = "UserAssigned"
)

func PossibleValuesForManagedServiceIdentityType() []string {
	return []string{
		string(ManagedServiceIdentityTypeNone),
		string(ManagedServiceIdentityTypeSystemAssigned),
		string(ManagedServiceIdentityTypeSystemAssignedUserAssigned),
		string(ManagedServiceIdentityTypeUserAssigned),
	}
}

func (s *ManagedServiceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedServiceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedServiceIdentityType(input string) (*ManagedServiceIdentityType, error) {
	vals := map[string]ManagedServiceIdentityType{
		"none":                        ManagedServiceIdentityTypeNone,
		"systemassigned":              ManagedServiceIdentityTypeSystemAssigned,
		"systemassigned,userassigned": ManagedServiceIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                ManagedServiceIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServiceIdentityType(input)
	return &out, nil
}

type NetworkType string

const (
	NetworkTypeOVNKubernetes NetworkType = "OVNKubernetes"
	NetworkTypeOther         NetworkType = "Other"
)

func PossibleValuesForNetworkType() []string {
	return []string{
		string(NetworkTypeOVNKubernetes),
		string(NetworkTypeOther),
	}
}

func (s *NetworkType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkType(input string) (*NetworkType, error) {
	vals := map[string]NetworkType{
		"ovnkubernetes": NetworkTypeOVNKubernetes,
		"other":         NetworkTypeOther,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkType(input)
	return &out, nil
}

type OutboundType string

const (
	OutboundTypeLoadBalancer OutboundType = "LoadBalancer"
)

func PossibleValuesForOutboundType() []string {
	return []string{
		string(OutboundTypeLoadBalancer),
	}
}

func (s *OutboundType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOutboundType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOutboundType(input string) (*OutboundType, error) {
	vals := map[string]OutboundType{
		"loadbalancer": OutboundTypeLoadBalancer,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OutboundType(input)
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

type StatusType string

const (
	StatusTypeFalse   StatusType = "False"
	StatusTypeTrue    StatusType = "True"
	StatusTypeUnknown StatusType = "Unknown"
)

func PossibleValuesForStatusType() []string {
	return []string{
		string(StatusTypeFalse),
		string(StatusTypeTrue),
		string(StatusTypeUnknown),
	}
}

func (s *StatusType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusType(input string) (*StatusType, error) {
	vals := map[string]StatusType{
		"false":   StatusTypeFalse,
		"true":    StatusTypeTrue,
		"unknown": StatusTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusType(input)
	return &out, nil
}

type Visibility string

const (
	VisibilityPrivate Visibility = "Private"
	VisibilityPublic  Visibility = "Public"
)

func PossibleValuesForVisibility() []string {
	return []string{
		string(VisibilityPrivate),
		string(VisibilityPublic),
	}
}

func (s *Visibility) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVisibility(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVisibility(input string) (*Visibility, error) {
	vals := map[string]Visibility{
		"private": VisibilityPrivate,
		"public":  VisibilityPublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Visibility(input)
	return &out, nil
}
