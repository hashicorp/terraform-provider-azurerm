package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WatchlistItemId struct {
	SubscriptionId string
	ResourceGroup  string
	WorkspaceName  string
	WatchlistName  string
	Name           string
}

func NewWatchlistItemID(subscriptionId, resourceGroup, workspaceName, watchlistName, name string) WatchlistItemId {
	return WatchlistItemId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		WorkspaceName:  workspaceName,
		WatchlistName:  watchlistName,
		Name:           name,
	}
}

func (id WatchlistItemId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Watchlist Name %q", id.WatchlistName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Watchlist Item", segmentsStr)
}

func (id WatchlistItemId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/watchlists/%s/watchlistItems/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.WatchlistName, id.Name)
}

// WatchlistItemID parses a WatchlistItem ID into an WatchlistItemId struct
func WatchlistItemID(input string) (*WatchlistItemId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WatchlistItemId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.WatchlistName, err = id.PopSegment("watchlists"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("watchlistItems"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
