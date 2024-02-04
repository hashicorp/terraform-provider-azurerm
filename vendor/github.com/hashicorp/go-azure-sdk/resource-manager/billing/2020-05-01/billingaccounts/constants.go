package billingaccounts

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountStatus string

const (
	AccountStatusActive      AccountStatus = "Active"
	AccountStatusDeleted     AccountStatus = "Deleted"
	AccountStatusDisabled    AccountStatus = "Disabled"
	AccountStatusExpired     AccountStatus = "Expired"
	AccountStatusExtended    AccountStatus = "Extended"
	AccountStatusTerminated  AccountStatus = "Terminated"
	AccountStatusTransferred AccountStatus = "Transferred"
)

func PossibleValuesForAccountStatus() []string {
	return []string{
		string(AccountStatusActive),
		string(AccountStatusDeleted),
		string(AccountStatusDisabled),
		string(AccountStatusExpired),
		string(AccountStatusExtended),
		string(AccountStatusTerminated),
		string(AccountStatusTransferred),
	}
}

func (s *AccountStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountStatus(input string) (*AccountStatus, error) {
	vals := map[string]AccountStatus{
		"active":      AccountStatusActive,
		"deleted":     AccountStatusDeleted,
		"disabled":    AccountStatusDisabled,
		"expired":     AccountStatusExpired,
		"extended":    AccountStatusExtended,
		"terminated":  AccountStatusTerminated,
		"transferred": AccountStatusTransferred,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountStatus(input)
	return &out, nil
}

type AccountType string

const (
	AccountTypeEnterprise AccountType = "Enterprise"
	AccountTypeIndividual AccountType = "Individual"
	AccountTypePartner    AccountType = "Partner"
)

func PossibleValuesForAccountType() []string {
	return []string{
		string(AccountTypeEnterprise),
		string(AccountTypeIndividual),
		string(AccountTypePartner),
	}
}

func (s *AccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccountType(input string) (*AccountType, error) {
	vals := map[string]AccountType{
		"enterprise": AccountTypeEnterprise,
		"individual": AccountTypeIndividual,
		"partner":    AccountTypePartner,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccountType(input)
	return &out, nil
}

type AgreementType string

const (
	AgreementTypeEnterpriseAgreement            AgreementType = "EnterpriseAgreement"
	AgreementTypeMicrosoftCustomerAgreement     AgreementType = "MicrosoftCustomerAgreement"
	AgreementTypeMicrosoftOnlineServicesProgram AgreementType = "MicrosoftOnlineServicesProgram"
	AgreementTypeMicrosoftPartnerAgreement      AgreementType = "MicrosoftPartnerAgreement"
)

func PossibleValuesForAgreementType() []string {
	return []string{
		string(AgreementTypeEnterpriseAgreement),
		string(AgreementTypeMicrosoftCustomerAgreement),
		string(AgreementTypeMicrosoftOnlineServicesProgram),
		string(AgreementTypeMicrosoftPartnerAgreement),
	}
}

func (s *AgreementType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAgreementType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAgreementType(input string) (*AgreementType, error) {
	vals := map[string]AgreementType{
		"enterpriseagreement":            AgreementTypeEnterpriseAgreement,
		"microsoftcustomeragreement":     AgreementTypeMicrosoftCustomerAgreement,
		"microsoftonlineservicesprogram": AgreementTypeMicrosoftOnlineServicesProgram,
		"microsoftpartneragreement":      AgreementTypeMicrosoftPartnerAgreement,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AgreementType(input)
	return &out, nil
}

type BillingProfileStatus string

const (
	BillingProfileStatusActive   BillingProfileStatus = "Active"
	BillingProfileStatusDisabled BillingProfileStatus = "Disabled"
	BillingProfileStatusWarned   BillingProfileStatus = "Warned"
)

func PossibleValuesForBillingProfileStatus() []string {
	return []string{
		string(BillingProfileStatusActive),
		string(BillingProfileStatusDisabled),
		string(BillingProfileStatusWarned),
	}
}

func (s *BillingProfileStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBillingProfileStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBillingProfileStatus(input string) (*BillingProfileStatus, error) {
	vals := map[string]BillingProfileStatus{
		"active":   BillingProfileStatusActive,
		"disabled": BillingProfileStatusDisabled,
		"warned":   BillingProfileStatusWarned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingProfileStatus(input)
	return &out, nil
}

type BillingRelationshipType string

const (
	BillingRelationshipTypeCSPPartner       BillingRelationshipType = "CSPPartner"
	BillingRelationshipTypeDirect           BillingRelationshipType = "Direct"
	BillingRelationshipTypeIndirectCustomer BillingRelationshipType = "IndirectCustomer"
	BillingRelationshipTypeIndirectPartner  BillingRelationshipType = "IndirectPartner"
)

func PossibleValuesForBillingRelationshipType() []string {
	return []string{
		string(BillingRelationshipTypeCSPPartner),
		string(BillingRelationshipTypeDirect),
		string(BillingRelationshipTypeIndirectCustomer),
		string(BillingRelationshipTypeIndirectPartner),
	}
}

func (s *BillingRelationshipType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBillingRelationshipType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBillingRelationshipType(input string) (*BillingRelationshipType, error) {
	vals := map[string]BillingRelationshipType{
		"csppartner":       BillingRelationshipTypeCSPPartner,
		"direct":           BillingRelationshipTypeDirect,
		"indirectcustomer": BillingRelationshipTypeIndirectCustomer,
		"indirectpartner":  BillingRelationshipTypeIndirectPartner,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingRelationshipType(input)
	return &out, nil
}

type InvoiceSectionState string

const (
	InvoiceSectionStateActive     InvoiceSectionState = "Active"
	InvoiceSectionStateRestricted InvoiceSectionState = "Restricted"
)

func PossibleValuesForInvoiceSectionState() []string {
	return []string{
		string(InvoiceSectionStateActive),
		string(InvoiceSectionStateRestricted),
	}
}

func (s *InvoiceSectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInvoiceSectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInvoiceSectionState(input string) (*InvoiceSectionState, error) {
	vals := map[string]InvoiceSectionState{
		"active":     InvoiceSectionStateActive,
		"restricted": InvoiceSectionStateRestricted,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InvoiceSectionState(input)
	return &out, nil
}

type SpendingLimit string

const (
	SpendingLimitOff SpendingLimit = "Off"
	SpendingLimitOn  SpendingLimit = "On"
)

func PossibleValuesForSpendingLimit() []string {
	return []string{
		string(SpendingLimitOff),
		string(SpendingLimitOn),
	}
}

func (s *SpendingLimit) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSpendingLimit(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSpendingLimit(input string) (*SpendingLimit, error) {
	vals := map[string]SpendingLimit{
		"off": SpendingLimitOff,
		"on":  SpendingLimitOn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SpendingLimit(input)
	return &out, nil
}

type SpendingLimitForBillingProfile string

const (
	SpendingLimitForBillingProfileOff SpendingLimitForBillingProfile = "Off"
	SpendingLimitForBillingProfileOn  SpendingLimitForBillingProfile = "On"
)

func PossibleValuesForSpendingLimitForBillingProfile() []string {
	return []string{
		string(SpendingLimitForBillingProfileOff),
		string(SpendingLimitForBillingProfileOn),
	}
}

func (s *SpendingLimitForBillingProfile) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSpendingLimitForBillingProfile(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSpendingLimitForBillingProfile(input string) (*SpendingLimitForBillingProfile, error) {
	vals := map[string]SpendingLimitForBillingProfile{
		"off": SpendingLimitForBillingProfileOff,
		"on":  SpendingLimitForBillingProfileOn,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SpendingLimitForBillingProfile(input)
	return &out, nil
}

type StatusReasonCode string

const (
	StatusReasonCodePastDue              StatusReasonCode = "PastDue"
	StatusReasonCodeSpendingLimitExpired StatusReasonCode = "SpendingLimitExpired"
	StatusReasonCodeSpendingLimitReached StatusReasonCode = "SpendingLimitReached"
)

func PossibleValuesForStatusReasonCode() []string {
	return []string{
		string(StatusReasonCodePastDue),
		string(StatusReasonCodeSpendingLimitExpired),
		string(StatusReasonCodeSpendingLimitReached),
	}
}

func (s *StatusReasonCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusReasonCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusReasonCode(input string) (*StatusReasonCode, error) {
	vals := map[string]StatusReasonCode{
		"pastdue":              StatusReasonCodePastDue,
		"spendinglimitexpired": StatusReasonCodeSpendingLimitExpired,
		"spendinglimitreached": StatusReasonCodeSpendingLimitReached,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusReasonCode(input)
	return &out, nil
}

type StatusReasonCodeForBillingProfile string

const (
	StatusReasonCodeForBillingProfilePastDue              StatusReasonCodeForBillingProfile = "PastDue"
	StatusReasonCodeForBillingProfileSpendingLimitExpired StatusReasonCodeForBillingProfile = "SpendingLimitExpired"
	StatusReasonCodeForBillingProfileSpendingLimitReached StatusReasonCodeForBillingProfile = "SpendingLimitReached"
)

func PossibleValuesForStatusReasonCodeForBillingProfile() []string {
	return []string{
		string(StatusReasonCodeForBillingProfilePastDue),
		string(StatusReasonCodeForBillingProfileSpendingLimitExpired),
		string(StatusReasonCodeForBillingProfileSpendingLimitReached),
	}
}

func (s *StatusReasonCodeForBillingProfile) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusReasonCodeForBillingProfile(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusReasonCodeForBillingProfile(input string) (*StatusReasonCodeForBillingProfile, error) {
	vals := map[string]StatusReasonCodeForBillingProfile{
		"pastdue":              StatusReasonCodeForBillingProfilePastDue,
		"spendinglimitexpired": StatusReasonCodeForBillingProfileSpendingLimitExpired,
		"spendinglimitreached": StatusReasonCodeForBillingProfileSpendingLimitReached,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusReasonCodeForBillingProfile(input)
	return &out, nil
}

type TargetCloud string

const (
	TargetCloudUSGov TargetCloud = "USGov"
	TargetCloudUSNat TargetCloud = "USNat"
	TargetCloudUSSec TargetCloud = "USSec"
)

func PossibleValuesForTargetCloud() []string {
	return []string{
		string(TargetCloudUSGov),
		string(TargetCloudUSNat),
		string(TargetCloudUSSec),
	}
}

func (s *TargetCloud) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTargetCloud(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTargetCloud(input string) (*TargetCloud, error) {
	vals := map[string]TargetCloud{
		"usgov": TargetCloudUSGov,
		"usnat": TargetCloudUSNat,
		"ussec": TargetCloudUSSec,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TargetCloud(input)
	return &out, nil
}
