package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountBatchConfigurationId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	BatchConfigurationName string
}

func NewIntegrationAccountBatchConfigurationID(subscriptionId, resourceGroup, integrationAccountName, batchConfigurationName string) IntegrationAccountBatchConfigurationId {
	return IntegrationAccountBatchConfigurationId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		BatchConfigurationName: batchConfigurationName,
	}
}

func (id IntegrationAccountBatchConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Batch Configuration Name %q", id.BatchConfigurationName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Batch Configuration", segmentsStr)
}

func (id IntegrationAccountBatchConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/batchConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.BatchConfigurationName)
}

// IntegrationAccountBatchConfigurationID parses a IntegrationAccountBatchConfiguration ID into an IntegrationAccountBatchConfigurationId struct
func IntegrationAccountBatchConfigurationID(input string) (*IntegrationAccountBatchConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountBatchConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.IntegrationAccountName, err = id.PopSegment("integrationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.BatchConfigurationName, err = id.PopSegment("batchConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
