package appservicecertificateorders

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateOrderActionType string

const (
	CertificateOrderActionTypeCertificateExpirationWarning CertificateOrderActionType = "CertificateExpirationWarning"
	CertificateOrderActionTypeCertificateExpired           CertificateOrderActionType = "CertificateExpired"
	CertificateOrderActionTypeCertificateIssued            CertificateOrderActionType = "CertificateIssued"
	CertificateOrderActionTypeCertificateOrderCanceled     CertificateOrderActionType = "CertificateOrderCanceled"
	CertificateOrderActionTypeCertificateOrderCreated      CertificateOrderActionType = "CertificateOrderCreated"
	CertificateOrderActionTypeCertificateRevoked           CertificateOrderActionType = "CertificateRevoked"
	CertificateOrderActionTypeDomainValidationComplete     CertificateOrderActionType = "DomainValidationComplete"
	CertificateOrderActionTypeFraudCleared                 CertificateOrderActionType = "FraudCleared"
	CertificateOrderActionTypeFraudDetected                CertificateOrderActionType = "FraudDetected"
	CertificateOrderActionTypeFraudDocumentationRequired   CertificateOrderActionType = "FraudDocumentationRequired"
	CertificateOrderActionTypeOrgNameChange                CertificateOrderActionType = "OrgNameChange"
	CertificateOrderActionTypeOrgValidationComplete        CertificateOrderActionType = "OrgValidationComplete"
	CertificateOrderActionTypeSanDrop                      CertificateOrderActionType = "SanDrop"
	CertificateOrderActionTypeUnknown                      CertificateOrderActionType = "Unknown"
)

func PossibleValuesForCertificateOrderActionType() []string {
	return []string{
		string(CertificateOrderActionTypeCertificateExpirationWarning),
		string(CertificateOrderActionTypeCertificateExpired),
		string(CertificateOrderActionTypeCertificateIssued),
		string(CertificateOrderActionTypeCertificateOrderCanceled),
		string(CertificateOrderActionTypeCertificateOrderCreated),
		string(CertificateOrderActionTypeCertificateRevoked),
		string(CertificateOrderActionTypeDomainValidationComplete),
		string(CertificateOrderActionTypeFraudCleared),
		string(CertificateOrderActionTypeFraudDetected),
		string(CertificateOrderActionTypeFraudDocumentationRequired),
		string(CertificateOrderActionTypeOrgNameChange),
		string(CertificateOrderActionTypeOrgValidationComplete),
		string(CertificateOrderActionTypeSanDrop),
		string(CertificateOrderActionTypeUnknown),
	}
}

func (s *CertificateOrderActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateOrderActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateOrderActionType(input string) (*CertificateOrderActionType, error) {
	vals := map[string]CertificateOrderActionType{
		"certificateexpirationwarning": CertificateOrderActionTypeCertificateExpirationWarning,
		"certificateexpired":           CertificateOrderActionTypeCertificateExpired,
		"certificateissued":            CertificateOrderActionTypeCertificateIssued,
		"certificateordercanceled":     CertificateOrderActionTypeCertificateOrderCanceled,
		"certificateordercreated":      CertificateOrderActionTypeCertificateOrderCreated,
		"certificaterevoked":           CertificateOrderActionTypeCertificateRevoked,
		"domainvalidationcomplete":     CertificateOrderActionTypeDomainValidationComplete,
		"fraudcleared":                 CertificateOrderActionTypeFraudCleared,
		"frauddetected":                CertificateOrderActionTypeFraudDetected,
		"frauddocumentationrequired":   CertificateOrderActionTypeFraudDocumentationRequired,
		"orgnamechange":                CertificateOrderActionTypeOrgNameChange,
		"orgvalidationcomplete":        CertificateOrderActionTypeOrgValidationComplete,
		"sandrop":                      CertificateOrderActionTypeSanDrop,
		"unknown":                      CertificateOrderActionTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateOrderActionType(input)
	return &out, nil
}

type CertificateOrderStatus string

const (
	CertificateOrderStatusCanceled          CertificateOrderStatus = "Canceled"
	CertificateOrderStatusDenied            CertificateOrderStatus = "Denied"
	CertificateOrderStatusExpired           CertificateOrderStatus = "Expired"
	CertificateOrderStatusIssued            CertificateOrderStatus = "Issued"
	CertificateOrderStatusNotSubmitted      CertificateOrderStatus = "NotSubmitted"
	CertificateOrderStatusPendingRekey      CertificateOrderStatus = "PendingRekey"
	CertificateOrderStatusPendingissuance   CertificateOrderStatus = "Pendingissuance"
	CertificateOrderStatusPendingrevocation CertificateOrderStatus = "Pendingrevocation"
	CertificateOrderStatusRevoked           CertificateOrderStatus = "Revoked"
	CertificateOrderStatusUnused            CertificateOrderStatus = "Unused"
)

func PossibleValuesForCertificateOrderStatus() []string {
	return []string{
		string(CertificateOrderStatusCanceled),
		string(CertificateOrderStatusDenied),
		string(CertificateOrderStatusExpired),
		string(CertificateOrderStatusIssued),
		string(CertificateOrderStatusNotSubmitted),
		string(CertificateOrderStatusPendingRekey),
		string(CertificateOrderStatusPendingissuance),
		string(CertificateOrderStatusPendingrevocation),
		string(CertificateOrderStatusRevoked),
		string(CertificateOrderStatusUnused),
	}
}

func (s *CertificateOrderStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateOrderStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateOrderStatus(input string) (*CertificateOrderStatus, error) {
	vals := map[string]CertificateOrderStatus{
		"canceled":          CertificateOrderStatusCanceled,
		"denied":            CertificateOrderStatusDenied,
		"expired":           CertificateOrderStatusExpired,
		"issued":            CertificateOrderStatusIssued,
		"notsubmitted":      CertificateOrderStatusNotSubmitted,
		"pendingrekey":      CertificateOrderStatusPendingRekey,
		"pendingissuance":   CertificateOrderStatusPendingissuance,
		"pendingrevocation": CertificateOrderStatusPendingrevocation,
		"revoked":           CertificateOrderStatusRevoked,
		"unused":            CertificateOrderStatusUnused,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateOrderStatus(input)
	return &out, nil
}

type CertificateProductType string

const (
	CertificateProductTypeStandardDomainValidatedSsl         CertificateProductType = "StandardDomainValidatedSsl"
	CertificateProductTypeStandardDomainValidatedWildCardSsl CertificateProductType = "StandardDomainValidatedWildCardSsl"
)

func PossibleValuesForCertificateProductType() []string {
	return []string{
		string(CertificateProductTypeStandardDomainValidatedSsl),
		string(CertificateProductTypeStandardDomainValidatedWildCardSsl),
	}
}

func (s *CertificateProductType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCertificateProductType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCertificateProductType(input string) (*CertificateProductType, error) {
	vals := map[string]CertificateProductType{
		"standarddomainvalidatedssl":         CertificateProductTypeStandardDomainValidatedSsl,
		"standarddomainvalidatedwildcardssl": CertificateProductTypeStandardDomainValidatedWildCardSsl,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateProductType(input)
	return &out, nil
}

type KeyVaultSecretStatus string

const (
	KeyVaultSecretStatusAzureServiceUnauthorizedToAccessKeyVault KeyVaultSecretStatus = "AzureServiceUnauthorizedToAccessKeyVault"
	KeyVaultSecretStatusCertificateOrderFailed                   KeyVaultSecretStatus = "CertificateOrderFailed"
	KeyVaultSecretStatusExternalPrivateKey                       KeyVaultSecretStatus = "ExternalPrivateKey"
	KeyVaultSecretStatusInitialized                              KeyVaultSecretStatus = "Initialized"
	KeyVaultSecretStatusKeyVaultDoesNotExist                     KeyVaultSecretStatus = "KeyVaultDoesNotExist"
	KeyVaultSecretStatusKeyVaultSecretDoesNotExist               KeyVaultSecretStatus = "KeyVaultSecretDoesNotExist"
	KeyVaultSecretStatusOperationNotPermittedOnKeyVault          KeyVaultSecretStatus = "OperationNotPermittedOnKeyVault"
	KeyVaultSecretStatusSucceeded                                KeyVaultSecretStatus = "Succeeded"
	KeyVaultSecretStatusUnknown                                  KeyVaultSecretStatus = "Unknown"
	KeyVaultSecretStatusUnknownError                             KeyVaultSecretStatus = "UnknownError"
	KeyVaultSecretStatusWaitingOnCertificateOrder                KeyVaultSecretStatus = "WaitingOnCertificateOrder"
)

func PossibleValuesForKeyVaultSecretStatus() []string {
	return []string{
		string(KeyVaultSecretStatusAzureServiceUnauthorizedToAccessKeyVault),
		string(KeyVaultSecretStatusCertificateOrderFailed),
		string(KeyVaultSecretStatusExternalPrivateKey),
		string(KeyVaultSecretStatusInitialized),
		string(KeyVaultSecretStatusKeyVaultDoesNotExist),
		string(KeyVaultSecretStatusKeyVaultSecretDoesNotExist),
		string(KeyVaultSecretStatusOperationNotPermittedOnKeyVault),
		string(KeyVaultSecretStatusSucceeded),
		string(KeyVaultSecretStatusUnknown),
		string(KeyVaultSecretStatusUnknownError),
		string(KeyVaultSecretStatusWaitingOnCertificateOrder),
	}
}

func (s *KeyVaultSecretStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyVaultSecretStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyVaultSecretStatus(input string) (*KeyVaultSecretStatus, error) {
	vals := map[string]KeyVaultSecretStatus{
		"azureserviceunauthorizedtoaccesskeyvault": KeyVaultSecretStatusAzureServiceUnauthorizedToAccessKeyVault,
		"certificateorderfailed":                   KeyVaultSecretStatusCertificateOrderFailed,
		"externalprivatekey":                       KeyVaultSecretStatusExternalPrivateKey,
		"initialized":                              KeyVaultSecretStatusInitialized,
		"keyvaultdoesnotexist":                     KeyVaultSecretStatusKeyVaultDoesNotExist,
		"keyvaultsecretdoesnotexist":               KeyVaultSecretStatusKeyVaultSecretDoesNotExist,
		"operationnotpermittedonkeyvault":          KeyVaultSecretStatusOperationNotPermittedOnKeyVault,
		"succeeded":                                KeyVaultSecretStatusSucceeded,
		"unknown":                                  KeyVaultSecretStatusUnknown,
		"unknownerror":                             KeyVaultSecretStatusUnknownError,
		"waitingoncertificateorder":                KeyVaultSecretStatusWaitingOnCertificateOrder,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyVaultSecretStatus(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateDeleting   ProvisioningState = "Deleting"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
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
		"canceled":   ProvisioningStateCanceled,
		"deleting":   ProvisioningStateDeleting,
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type ResourceNotRenewableReason string

const (
	ResourceNotRenewableReasonExpirationNotInRenewalTimeRange          ResourceNotRenewableReason = "ExpirationNotInRenewalTimeRange"
	ResourceNotRenewableReasonRegistrationStatusNotSupportedForRenewal ResourceNotRenewableReason = "RegistrationStatusNotSupportedForRenewal"
	ResourceNotRenewableReasonSubscriptionNotActive                    ResourceNotRenewableReason = "SubscriptionNotActive"
)

func PossibleValuesForResourceNotRenewableReason() []string {
	return []string{
		string(ResourceNotRenewableReasonExpirationNotInRenewalTimeRange),
		string(ResourceNotRenewableReasonRegistrationStatusNotSupportedForRenewal),
		string(ResourceNotRenewableReasonSubscriptionNotActive),
	}
}

func (s *ResourceNotRenewableReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourceNotRenewableReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourceNotRenewableReason(input string) (*ResourceNotRenewableReason, error) {
	vals := map[string]ResourceNotRenewableReason{
		"expirationnotinrenewaltimerange":          ResourceNotRenewableReasonExpirationNotInRenewalTimeRange,
		"registrationstatusnotsupportedforrenewal": ResourceNotRenewableReasonRegistrationStatusNotSupportedForRenewal,
		"subscriptionnotactive":                    ResourceNotRenewableReasonSubscriptionNotActive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourceNotRenewableReason(input)
	return &out, nil
}
