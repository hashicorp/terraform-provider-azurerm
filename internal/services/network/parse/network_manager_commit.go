package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type ManagerCommitId struct {
	SubscriptionId     string
	ResourceGroup      string
	NetworkManagerName string
	Location           string
	ScopeAccess        string
}

func NewNetworkManagerCommitID(subscriptionId string, resourceGroup string, networkManagerName string, location string, scopeAccess string) *ManagerCommitId {
	return &ManagerCommitId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		NetworkManagerName: networkManagerName,
		Location:           azure.NormalizeLocation(location),
		ScopeAccess:        scopeAccess,
	}
}

func NetworkManagerCommitID(networkManagerCommitId string) (*ManagerCommitId, error) {
	v := strings.Split(networkManagerCommitId, "|")
	if len(v) != 3 {
		return nil, fmt.Errorf("expected the network manager commit ID to be in format `{networkManagerId}/commit|{location}|{scopeAccess}`, but got %d segments ", len(v))
	}

	networkManagerId := strings.TrimSuffix(v[0], "/commit")
	managerId, err := NetworkManagerID(networkManagerId)
	if err != nil {
		return nil, err
	}

	if v[1] == "" {
		return nil, fmt.Errorf("expected location in network manager commit ID with format `{networkManagerId}/commit|{location}|{scopeAccess}`, but got %s in %s", v[1], networkManagerCommitId)
	}
	normalizedLocation := azure.NormalizeLocation(v[1])

	if v[2] == "" {
		return nil, fmt.Errorf("expected scopeAccess in network manager commit ID with format `{networkManagerId}/commit|{location}|{scopeAccess} to be one of the [Connectivity, SecurityAdmin]`, but got %s in %s", v[2], networkManagerCommitId)
	}
	scopeAccess := v[2]
	networkManagerCommit := NewNetworkManagerCommitID(managerId.SubscriptionId, managerId.ResourceGroup, managerId.Name, normalizedLocation, scopeAccess)
	return networkManagerCommit, nil
}

func (id ManagerCommitId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/commit|%s|%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.Location, id.ScopeAccess)
}

func (id ManagerCommitId) String() string {
	segments := []string{
		fmt.Sprintf("Scope Access %q", id.ScopeAccess),
		fmt.Sprintf("Location %q", id.Location),
		fmt.Sprintf("Network Manager %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Commit", segmentsStr)
}
