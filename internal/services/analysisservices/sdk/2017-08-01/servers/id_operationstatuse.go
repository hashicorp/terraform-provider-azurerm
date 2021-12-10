package servers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type OperationstatuseId struct {
	SubscriptionId string
	LocationName   string
	Name           string
}

func NewOperationstatuseID(subscriptionId, locationName, name string) OperationstatuseId {
	return OperationstatuseId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		Name:           name,
	}
}

func (id OperationstatuseId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Location Name %q", id.LocationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Operationstatuse", segmentsStr)
}

func (id OperationstatuseId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.AnalysisServices/locations/%s/operationstatuses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.Name)
}

// ParseOperationstatuseID parses a Operationstatuse ID into an OperationstatuseId struct
func ParseOperationstatuseID(input string) (*OperationstatuseId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OperationstatuseId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.LocationName, err = id.PopSegment("locations"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("operationstatuses"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// ParseOperationstatuseIDInsensitively parses an Operationstatuse ID into an OperationstatuseId struct, insensitively
// This should only be used to parse an ID for rewriting to a consistent casing,
// the ParseOperationstatuseID method should be used instead for validation etc.
func ParseOperationstatuseIDInsensitively(input string) (*OperationstatuseId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := OperationstatuseId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	// find the correct casing for the 'locations' segment
	locationsKey := "locations"
	for key := range id.Path {
		if strings.EqualFold(key, locationsKey) {
			locationsKey = key
			break
		}
	}
	if resourceId.LocationName, err = id.PopSegment(locationsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'operationstatuses' segment
	operationstatusesKey := "operationstatuses"
	for key := range id.Path {
		if strings.EqualFold(key, operationstatusesKey) {
			operationstatusesKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(operationstatusesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
