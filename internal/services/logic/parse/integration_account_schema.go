package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type IntegrationAccountSchemaId struct {
	SubscriptionId         string
	ResourceGroup          string
	IntegrationAccountName string
	SchemaName             string
}

func NewIntegrationAccountSchemaID(subscriptionId, resourceGroup, integrationAccountName, schemaName string) IntegrationAccountSchemaId {
	return IntegrationAccountSchemaId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		IntegrationAccountName: integrationAccountName,
		SchemaName:             schemaName,
	}
}

func (id IntegrationAccountSchemaId) String() string {
	segments := []string{
		fmt.Sprintf("Schema Name %q", id.SchemaName),
		fmt.Sprintf("Integration Account Name %q", id.IntegrationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Integration Account Schema", segmentsStr)
}

func (id IntegrationAccountSchemaId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/schemas/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.IntegrationAccountName, id.SchemaName)
}

// IntegrationAccountSchemaID parses a IntegrationAccountSchema ID into an IntegrationAccountSchemaId struct
func IntegrationAccountSchemaID(input string) (*IntegrationAccountSchemaId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := IntegrationAccountSchemaId{
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
	if resourceId.SchemaName, err = id.PopSegment("schemas"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
