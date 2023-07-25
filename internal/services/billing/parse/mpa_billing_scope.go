// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

var _ resourceids.Id = MicrosoftPartnerAccountBillingScopeId{}

type MicrosoftPartnerAccountBillingScopeId struct {
	BillingAccountName string
	CustomerName       string
}

func NewMPABillingScopeID(billingAccountName, customerName string) MicrosoftPartnerAccountBillingScopeId {
	return MicrosoftPartnerAccountBillingScopeId{
		BillingAccountName: billingAccountName,
		CustomerName:       customerName,
	}
}

func (id MicrosoftPartnerAccountBillingScopeId) String() string {
	segments := []string{
		fmt.Sprintf("Customer Name %q", id.CustomerName),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "M P A Billing Scope", segmentsStr)
}

func (id MicrosoftPartnerAccountBillingScopeId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/customers/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.CustomerName)
}

// MicrosoftPartnerAccountBillingScopeID parses a MPABillingScope ID into an MicrosoftPartnerAccountBillingScopeId struct
func MicrosoftPartnerAccountBillingScopeID(input string) (*MicrosoftPartnerAccountBillingScopeId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := MicrosoftPartnerAccountBillingScopeId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.CustomerName, err = id.PopSegment("customers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
