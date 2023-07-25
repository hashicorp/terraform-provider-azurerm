// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type BillingAccountCostManagementExportId struct {
	BillingAccountName string
	ExportName         string
}

func NewBillingAccountCostManagementExportID(billingAccountName, exportName string) BillingAccountCostManagementExportId {
	return BillingAccountCostManagementExportId{
		BillingAccountName: billingAccountName,
		ExportName:         exportName,
	}
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
