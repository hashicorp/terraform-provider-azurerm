// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is written by hand, as the ID generator cannot handle IDs without a subscription

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type BillingAccountCostManagementExportId struct {
	BillingAccountName string
	ExportName         string
}

func (id BillingAccountCostManagementExportId) String() string {
	segments := []string{
		fmt.Sprintf("Export Name %q", id.ExportName),
		fmt.Sprintf("Billing Account Name %q", id.BillingAccountName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Billing Account Cost Management Export", segmentsStr)
}

func (id BillingAccountCostManagementExportId) ID() string {
	fmtString := "/providers/Microsoft.Billing/billingAccounts/%s/providers/Microsoft.CostManagement/exports/%s"
	return fmt.Sprintf(fmtString, id.BillingAccountName, id.ExportName)
}

// BillingAccountCostManagementExportID parses a BillingAccountCostManagementExport ID into an BillingAccountCostManagementExportId struct
func BillingAccountCostManagementExportID(input string) (*BillingAccountCostManagementExportId, error) {
	id, err := azure.ParseAzureResourceIDWithoutSubscription(input)
	if err != nil {
		return nil, err
	}

	resourceId := BillingAccountCostManagementExportId{}

	if resourceId.BillingAccountName, err = id.PopSegment("billingAccounts"); err != nil {
		return nil, err
	}
	if resourceId.ExportName, err = id.PopSegment("exports"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
