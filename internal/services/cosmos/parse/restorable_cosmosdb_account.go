package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type RestorableCosmosdbAccountId struct {
	SubscriptionId                string
	LocationName                  string
	RestorableDatabaseAccountName string
}

func NewRestorableCosmosdbAccountID(subscriptionId, locationName, restorableDatabaseAccountName string) RestorableCosmosdbAccountId {
	return RestorableCosmosdbAccountId{
		SubscriptionId:                subscriptionId,
		LocationName:                  locationName,
		RestorableDatabaseAccountName: restorableDatabaseAccountName,
	}
}

func (id RestorableCosmosdbAccountId) String() string {
	segments := []string{
		fmt.Sprintf("Restorable Database Account Name %q", id.RestorableDatabaseAccountName),
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Restorable Cosmosdb Account", segmentsStr)
}

func (id RestorableCosmosdbAccountId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DocumentDB/locations/%s/restorableDatabaseAccounts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.RestorableDatabaseAccountName)
}

// RestorableCosmosdbAccountID parses a RestorableCosmosdbAccount ID into an RestorableCosmosdbAccountId struct
func RestorableCosmosdbAccountID(input string) (*RestorableCosmosdbAccountId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RestorableCosmosdbAccountId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.RestorableDatabaseAccountName, err = id.PopSegment("restorableDatabaseAccounts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
