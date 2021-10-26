package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type CostManagementExportId struct {
	ResourceId string
	Name       string
}

func CostManagementExportID(input string) (*CostManagementExportId, error) {
	service := CostManagementExportId{}

	if strings.Contains(input, "subscriptions") {
		id, err := azure.ParseAzureResourceID(input)

		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unable to parse Cost Management Export %q: %+v", input, err)
		}

		if id.ResourceGroup == "" {
			service.ResourceId = fmt.Sprintf("/subscriptions/%s", id.SubscriptionID)
		} else {
			service.ResourceId = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", id.SubscriptionID, id.ResourceGroup)
		}

		if service.Name, err = id.PopSegment("exports"); err != nil {
			return nil, err
		}

		if err := id.ValidateNoEmptySegments(input); err != nil {
			return nil, err
		}
	} else {
		id, err := azure.ParseAzureResourceIDWithoutSubscription(input)

		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unable to parse Cost Management Export %q: %+v", input, err)
		}

		if mgtGrp, ok := id.Path["managementGroups"]; ok {
			service.ResourceId = fmt.Sprintf("/providers/%s/managementGroups/%s", id.Provider, mgtGrp)
		}

		if id.Provider == "Microsoft.Billing" {
			if billingAcc, ok := id.Path["billingAccounts"]; ok {
				billingAccounts := fmt.Sprintf("/providers/%s/billingAccounts/%s", id.Provider, billingAcc)
				if len(id.Path) == 2 {
					service.ResourceId = billingAccounts
				}

				if departments, ok := id.Path["departments"]; ok {
					service.ResourceId = fmt.Sprintf("%s/departments/%s", billingAccounts, departments)
				}

				if enrollmentAccounts, ok := id.Path["enrollmentAccounts"]; ok {
					service.ResourceId = fmt.Sprintf("%s/enrollmentAccounts/%s", billingAccounts, enrollmentAccounts)
				}

				if billingProfile, ok := id.Path["billingProfiles"]; ok {
					service.ResourceId = fmt.Sprintf("%s/billingProfiles/%s", billingAccounts, billingProfile)
					if invoiceSections, ok := id.Path["billingProfiles"]; ok {
						service.ResourceId = fmt.Sprintf("%s/billingProfiles/%s/invoiceSections/%s", billingAccounts, billingProfile, invoiceSections)
					}
				}

				if value, ok := id.Path["customers"]; ok {
					service.ResourceId = fmt.Sprintf("%s/customers/%s", billingAccounts, value)
				}
			}
		}

		if service.ResourceId == "" {
			return nil, fmt.Errorf("[ERROR] Unable to parse Cost Management Export %q: %+v", input, err)
		}

		if service.Name, err = id.PopSegment("exports"); err != nil {
			return nil, err
		}
	}

	return &service, nil
}