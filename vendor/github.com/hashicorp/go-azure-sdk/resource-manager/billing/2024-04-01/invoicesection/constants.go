package invoicesection

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteInvoiceSectionEligibilityCode string

const (
	DeleteInvoiceSectionEligibilityCodeActiveAzurePlans           DeleteInvoiceSectionEligibilityCode = "ActiveAzurePlans"
	DeleteInvoiceSectionEligibilityCodeActiveBillingSubscriptions DeleteInvoiceSectionEligibilityCode = "ActiveBillingSubscriptions"
	DeleteInvoiceSectionEligibilityCodeLastInvoiceSection         DeleteInvoiceSectionEligibilityCode = "LastInvoiceSection"
	DeleteInvoiceSectionEligibilityCodeOther                      DeleteInvoiceSectionEligibilityCode = "Other"
	DeleteInvoiceSectionEligibilityCodeReservedInstances          DeleteInvoiceSectionEligibilityCode = "ReservedInstances"
)

func PossibleValuesForDeleteInvoiceSectionEligibilityCode() []string {
	return []string{
		string(DeleteInvoiceSectionEligibilityCodeActiveAzurePlans),
		string(DeleteInvoiceSectionEligibilityCodeActiveBillingSubscriptions),
		string(DeleteInvoiceSectionEligibilityCodeLastInvoiceSection),
		string(DeleteInvoiceSectionEligibilityCodeOther),
		string(DeleteInvoiceSectionEligibilityCodeReservedInstances),
	}
}

func (s *DeleteInvoiceSectionEligibilityCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteInvoiceSectionEligibilityCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteInvoiceSectionEligibilityCode(input string) (*DeleteInvoiceSectionEligibilityCode, error) {
	vals := map[string]DeleteInvoiceSectionEligibilityCode{
		"activeazureplans":           DeleteInvoiceSectionEligibilityCodeActiveAzurePlans,
		"activebillingsubscriptions": DeleteInvoiceSectionEligibilityCodeActiveBillingSubscriptions,
		"lastinvoicesection":         DeleteInvoiceSectionEligibilityCodeLastInvoiceSection,
		"other":                      DeleteInvoiceSectionEligibilityCodeOther,
		"reservedinstances":          DeleteInvoiceSectionEligibilityCodeReservedInstances,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteInvoiceSectionEligibilityCode(input)
	return &out, nil
}

type DeleteInvoiceSectionEligibilityStatus string

const (
	DeleteInvoiceSectionEligibilityStatusAllowed    DeleteInvoiceSectionEligibilityStatus = "Allowed"
	DeleteInvoiceSectionEligibilityStatusNotAllowed DeleteInvoiceSectionEligibilityStatus = "NotAllowed"
)

func PossibleValuesForDeleteInvoiceSectionEligibilityStatus() []string {
	return []string{
		string(DeleteInvoiceSectionEligibilityStatusAllowed),
		string(DeleteInvoiceSectionEligibilityStatusNotAllowed),
	}
}

func (s *DeleteInvoiceSectionEligibilityStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteInvoiceSectionEligibilityStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteInvoiceSectionEligibilityStatus(input string) (*DeleteInvoiceSectionEligibilityStatus, error) {
	vals := map[string]DeleteInvoiceSectionEligibilityStatus{
		"allowed":    DeleteInvoiceSectionEligibilityStatusAllowed,
		"notallowed": DeleteInvoiceSectionEligibilityStatusNotAllowed,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteInvoiceSectionEligibilityStatus(input)
	return &out, nil
}

type InvoiceSectionState string

const (
	InvoiceSectionStateActive      InvoiceSectionState = "Active"
	InvoiceSectionStateDeleted     InvoiceSectionState = "Deleted"
	InvoiceSectionStateDisabled    InvoiceSectionState = "Disabled"
	InvoiceSectionStateOther       InvoiceSectionState = "Other"
	InvoiceSectionStateRestricted  InvoiceSectionState = "Restricted"
	InvoiceSectionStateUnderReview InvoiceSectionState = "UnderReview"
	InvoiceSectionStateWarned      InvoiceSectionState = "Warned"
)

func PossibleValuesForInvoiceSectionState() []string {
	return []string{
		string(InvoiceSectionStateActive),
		string(InvoiceSectionStateDeleted),
		string(InvoiceSectionStateDisabled),
		string(InvoiceSectionStateOther),
		string(InvoiceSectionStateRestricted),
		string(InvoiceSectionStateUnderReview),
		string(InvoiceSectionStateWarned),
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
		"active":      InvoiceSectionStateActive,
		"deleted":     InvoiceSectionStateDeleted,
		"disabled":    InvoiceSectionStateDisabled,
		"other":       InvoiceSectionStateOther,
		"restricted":  InvoiceSectionStateRestricted,
		"underreview": InvoiceSectionStateUnderReview,
		"warned":      InvoiceSectionStateWarned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InvoiceSectionState(input)
	return &out, nil
}

type InvoiceSectionStateReasonCode string

const (
	InvoiceSectionStateReasonCodeOther                InvoiceSectionStateReasonCode = "Other"
	InvoiceSectionStateReasonCodePastDue              InvoiceSectionStateReasonCode = "PastDue"
	InvoiceSectionStateReasonCodeSpendingLimitExpired InvoiceSectionStateReasonCode = "SpendingLimitExpired"
	InvoiceSectionStateReasonCodeSpendingLimitReached InvoiceSectionStateReasonCode = "SpendingLimitReached"
	InvoiceSectionStateReasonCodeUnusualActivity      InvoiceSectionStateReasonCode = "UnusualActivity"
)

func PossibleValuesForInvoiceSectionStateReasonCode() []string {
	return []string{
		string(InvoiceSectionStateReasonCodeOther),
		string(InvoiceSectionStateReasonCodePastDue),
		string(InvoiceSectionStateReasonCodeSpendingLimitExpired),
		string(InvoiceSectionStateReasonCodeSpendingLimitReached),
		string(InvoiceSectionStateReasonCodeUnusualActivity),
	}
}

func (s *InvoiceSectionStateReasonCode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseInvoiceSectionStateReasonCode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseInvoiceSectionStateReasonCode(input string) (*InvoiceSectionStateReasonCode, error) {
	vals := map[string]InvoiceSectionStateReasonCode{
		"other":                InvoiceSectionStateReasonCodeOther,
		"pastdue":              InvoiceSectionStateReasonCodePastDue,
		"spendinglimitexpired": InvoiceSectionStateReasonCodeSpendingLimitExpired,
		"spendinglimitreached": InvoiceSectionStateReasonCodeSpendingLimitReached,
		"unusualactivity":      InvoiceSectionStateReasonCodeUnusualActivity,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InvoiceSectionStateReasonCode(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled         ProvisioningState = "Canceled"
	ProvisioningStateConfirmedBilling ProvisioningState = "ConfirmedBilling"
	ProvisioningStateCreated          ProvisioningState = "Created"
	ProvisioningStateCreating         ProvisioningState = "Creating"
	ProvisioningStateExpired          ProvisioningState = "Expired"
	ProvisioningStateFailed           ProvisioningState = "Failed"
	ProvisioningStateNew              ProvisioningState = "New"
	ProvisioningStatePending          ProvisioningState = "Pending"
	ProvisioningStatePendingBilling   ProvisioningState = "PendingBilling"
	ProvisioningStateProvisioning     ProvisioningState = "Provisioning"
	ProvisioningStateSucceeded        ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateConfirmedBilling),
		string(ProvisioningStateCreated),
		string(ProvisioningStateCreating),
		string(ProvisioningStateExpired),
		string(ProvisioningStateFailed),
		string(ProvisioningStateNew),
		string(ProvisioningStatePending),
		string(ProvisioningStatePendingBilling),
		string(ProvisioningStateProvisioning),
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
		"canceled":         ProvisioningStateCanceled,
		"confirmedbilling": ProvisioningStateConfirmedBilling,
		"created":          ProvisioningStateCreated,
		"creating":         ProvisioningStateCreating,
		"expired":          ProvisioningStateExpired,
		"failed":           ProvisioningStateFailed,
		"new":              ProvisioningStateNew,
		"pending":          ProvisioningStatePending,
		"pendingbilling":   ProvisioningStatePendingBilling,
		"provisioning":     ProvisioningStateProvisioning,
		"succeeded":        ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
